````markdown
# p1-go — Meter CLI

Lightweight CLI to fetch readings from a P1-compatible meter (JSON endpoint), parse the payload, and persist readings into PostgreSQL.

This README shows how to build, configure, run, test, and troubleshoot the CLI.

## What's included
- Go-based CLI entrypoint: `cmd/metercli`
- Parser for meter JSON payloads: `src/services/parser`
- Postgres persistence adapter with idempotent upsert: `src/services/db`
- File-backed JSON-lines buffer for offline persistence: `src/buffer`
- Scheduler with advisory-lock based single-run semantics: `src/scheduler`
- Implementation scaffolding & specs: `specs/002-implement-scheduler`

## Requirements
- Go 1.20+ (module-aware)
- PostgreSQL for persistence
- Optional: Docker (for running Postgres during integration tests)

## Build

From the repository root:

```bash
go build -o bin/metercli ./cmd/metercli
```

The `build-all.sh` helper will also build artifacts if present:

```bash
./build-all.sh
```

## Configuration

The CLI accepts a JSON configuration file via `--config` (default `./config.json`). Example at `config.example.json`.

Supported fields:

- `meter_endpoint` (string) — HTTP URL to fetch the meter JSON payload.
- `db_dsn` (string) — Postgres DSN. Use the lib/pq key=value form to avoid URL-encoding issues for passwords with special characters. You can also add `options='-c search_path=p1'` if the DB user only has access to the `p1` schema.

Example `config.json`:

```json
{
  "meter_endpoint": "http://192.168.101.20/api/v1/data",
  "db_dsn": "host=127.0.0.1 port=5432 user=p1 password='secret' dbname=postgres sslmode=disable options='-c search_path=p1'"
}
```

Environment variables are supported as a fallback and are also set by the CLI after loading a config file for compatibility with packages that read env vars directly:

- `METER_ENDPOINT`
- `DB_DSN`

Example (environment):

```bash
export METER_ENDPOINT="http://192.168.101.20/api/v1/data"
export DB_DSN="host=127.0.0.1 port=5432 user=p1 password='secret' dbname=postgres sslmode=disable options='-c search_path=p1'"
./bin/metercli
```

## CLI flags

- `--config <path>` — path to JSON config file (default `./config.json`).
- `--loop` — run continuously using the internal scheduler.
- `--interval <seconds>` — interval for scheduler loop (default 60).
- `--drain-buffer` — drain the on-disk buffer (`/tmp/p1-buffer.jsonl` by default) and attempt to persist entries.

Example: run continuously every 60s:

```bash
./bin/metercli --config ./config.json --loop --interval 60
```

Drain buffer manually:

```bash
./bin/metercli --config ./config.json --drain-buffer
```

## Database & Migrations

Schema and tables are defined in `migrations/001_create_tables.sql`. To apply the migration using `psql`:

```bash
# as a user with privileges to create schema/tables
psql "host=127.0.0.1 port=5432 user=postgres dbname=postgres" -f migrations/001_create_tables.sql
```

If your application user only has access to schema `p1`, include `options='-c search_path=p1'` in the DSN or qualify table names in SQL.

Permissions note: the DB role used by the CLI must have permission to INSERT into the `p1` tables and USAGE on the sequences backing the serial columns. If you see `permission denied for sequence ...` grant usage or adjust ownership.

## Buffering and Offline Mode

The CLI buffers failed persistence attempts to `/tmp/p1-buffer.jsonl` (JSON-lines). The buffer supports two operations:

- `Append(v)` — append a JSON entry to the buffer (used when DB insert fails).
- `Drain(ctx, persistFn)` — read all entries and call `persistFn(ctx, json.RawMessage)` for each entry. On success the buffer file is truncated.

Current behavior: `Drain` aborts on the first persist error. Recommendation: extend to partial-drain semantics for robust retrying.

## Scheduling

Two options:

1. External scheduler (cron/systemd)
   - Preferred for production deployments. Example cron runs every minute.

2. Internal scheduler (`--loop`)
   - `cmd/metercli --loop` starts a ticker-based scheduler and uses a Postgres advisory lock to avoid overlapping runs. This is convenient for single-host deployments.

Scheduler uses an advisory lock (pg_try_advisory_lock) to ensure only one runner performs work at a time.

## Testing

Unit, contract, and integration-style tests are in `tests/` and can be run with:

```bash
go test ./...
```

There is an integration docker-compose that runs Postgres (mapped to host port 5433 to avoid collisions). Start it with:

```bash
docker compose up -d
```

Run integration tests:

```bash
TEST_DATABASE_DSN="host=127.0.0.1 port=5433 user=testuser password=testpass dbname=testdb sslmode=disable" \
  go test ./tests/integration -v
```

## Troubleshooting

- Error: `invalid port ":..." after host` — Password contains URL-reserved characters. Use key=value DSN form or percent-encode the password for URL DSNs.
- Error: `permission denied for sequence ...` — Grant USAGE on the sequence or ensure the DB role owns the sequence.
- Error: `DB_DSN not set` — Ensure `--config` is provided or `DB_DSN` env var is set.

## Developer notes

- Follow the project's constitution `.specify/memory/constitution.md` and tests-first workflow when adding features.
- Add contract tests under `tests/contract/` and integration tests under `tests/integration/` before implementing features.

## Contributing

PRs welcome. Please include tests and update `specs/` where appropriate.

## License

MIT

````
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
