```markdown
# Research: Implement scheduler that runs meter fetch+persist

Date: 2025-10-03

Decision: Implement an internal scheduler (goroutine + ticker) guarded by a Postgres advisory lock to avoid overlapping runs.

Rationale:

- The repository is a Go CLI service that already includes a Postgres-backed scheduler package and uses `pg_try_advisory_lock` in tests. Using Postgres advisory locks is lightweight and reliable when Postgres is already the persistence layer.
- An internal scheduler simplifies deployment for users who run the CLI under systemd or as a standalone process. It also supports a `--loop` mode for continuous operation.

Alternatives considered:

- System cron / systemd timers: simpler but requires external orchestration and makes testing harder.
- File-based flock: works for single-host but does not work across multiple hosts where concurrency could occur.

Unknowns resolved (assumptions made):

- Frequency: default every minute (60s). Exposed via `--interval` flag in seconds.
- Lock key: use a fixed advisory lock key derived from a small integer (e.g., 42) configurable via code constant or flag.
- Error handling: failures during persist will be recorded to a local JSON-lines buffer and retried by a drain operation.
```
