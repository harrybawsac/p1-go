# p1-go — Meter CLI

Small command-line tool to fetch gas & electricity meter JSON (P1) and persist readings into PostgreSQL.

This README shows how to build, run, test, and schedule the CLI.

## Quickstart

Prerequisites:

- Go 1.20+
- Docker (for running Postgres in tests)

Build:

```bash
go build -o bin/metercli ./cmd/metercli
```

Run once:

```bash
# Run the CLI once (the current implementation executes runOnce)
./bin/metercli
```

## Configuration

The CLI supports two ways to configure runtime settings:

1. JSON config file (preferred for local installs)
2. Environment variables (fallback and common in CI/containers)

## Config file (JSON)

Create a JSON file (example: `config.json`) with the following fields (an example file `config.example.json` is included in the repository):

```json
{
  "meter_endpoint": "http://192.168.101.20/api/v1/data",
  "db_dsn": "postgres://p1:password@127.0.0.1:5432/postgres?sslmode=disable"
}
```

Pass the config file to the CLI using `--config`:

```bash
./bin/metercli --config ./config.json
```

If `--config` is not provided, the CLI will look for `./config.json`.

## Environment variables (fallback)

You can also configure the runtime via environment variables. These are used if present and are also set by the CLI after loading a config file for compatibility with existing code paths:

- `METER_ENDPOINT` — HTTP endpoint of the meter
- `DB_DSN` — Postgres DSN

Example using env vars directly:

```bash
export METER_ENDPOINT="http://192.168.101.20/api/v1/data"
export DB_DSN="postgres://p1:password@127.0.0.1:5432/postgres?sslmode=disable"
./bin/metercli
```

## Migrations

The repository contains migrations under `migrations/`. To create the `p1` schema and tables used by the CLI, apply the SQL migration:

```bash
# Example using psql (replace DSN parts accordingly)
export PGPASSWORD="yourpassword"
psql -h localhost -U youruser -d yourdb -f migrations/001_create_tables.sql
```

## Testing

Unit tests run with the normal `go test` tooling:

```bash
go test ./...
```

## Integration tests

An integration test and a `docker-compose.yml` are included to run a transient Postgres for end-to-end validation.

Start Postgres (docker compose will map host port 5433 -> container 5432 to avoid collisions):

```bash
docker compose up -d
```

Run the integration test (it expects the test DB user/database created by the compose file):

```bash
TEST_DATABASE_DSN="postgres://testuser:testpass@127.0.0.1:5433/testdb?sslmode=disable" \
  go test ./tests/integration -run TestIntegration_InsertReading -v
```

Cleanup the compose resources when finished:

```bash
docker compose down -v
```

## Scheduling (cron)

The CLI is designed to be run periodically (every minute) by an external scheduler like cron or systemd timers. Example cron entry (runs every minute):

```cron
* * * * * /path/to/bin/metercli >> /var/log/metercli.log 2>&1
```

Systemd timer example (place these files under `/etc/systemd/system`) — `metercli.service`:

```
[Unit]
Description=Meter CLI runner

[Service]
Type=oneshot
ExecStart=/usr/local/bin/metercli
```

`metercli.timer`:

```
[Unit]
Description=Run metercli every minute

[Timer]
OnBootSec=1min
OnUnitActiveSec=1min

[Install]
WantedBy=timers.target
```

Enable and start the timer:

```bash
sudo systemctl enable --now metercli.timer
```

## Notes and next steps

- The core fetch/parse/persist flow is implemented in `src/services/*`. `cmd/metercli/main.go` currently calls `runOnce()` as a placeholder. Next tasks include wiring configuration, implementing the scheduler and advisory locking, adding offline buffering, and exposing metrics.
- Do not commit DB credentials. Use environment variables or a secrets manager for production.

## Contributing

Please follow the project's constitution (`.specify/memory/constitution.md`) and tests-first workflow. Open issues or PRs for feature work and link to the spec under `specs/001-build-a-cli/`.

## License

MIT (or your chosen license)
