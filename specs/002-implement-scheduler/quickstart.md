````markdown
# Quickstart: Scheduler feature (run loop mode)

Prereqs:

- PostgreSQL database and migrations applied (see `/migrations/001_create_tables.sql`).

Run once (single fetch):

1. Build CLI:

   go build -o bin/metercli ./cmd/metercli

2. Create a config file, e.g. `/home/lex/projects/p1-go/config.example.json`:

```json
{
  "meter_endpoint": "http://localhost:8080/meter",
  "db_dsn": "postgres://testuser:testpass@localhost:5433/testdb?sslmode=disable"
}
```

3. Run once:

   ./bin/metercli --config /absolute/path/to/config.json

Run continuous loop (every 60s):

./bin/metercli --config /absolute/path/to/config.json --loop --interval 60

Notes:

- Metrics and observability are intentionally excluded for this feature (per user request).
- Failed writes are buffered locally to `/tmp/p1-buffer.jsonl` and can be drained by running `./bin/metercli --drain-buffer` (subject to implementation).
````
