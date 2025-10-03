# Tasks: Build a CLI to read gas & electricity meters (detailed)

**Input**: Design documents from `/home/lex/projects/p1-build-a-cli/` (plan.md, spec.md)
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

All paths below are relative to the repository root `/home/lex/projects/p1-go`.

Note on credentials: DO NOT hardcode credentials in code. Use environment
variables or a config file stored outside version control. Example credentials
were provided in the spec for integration testing only.

TDD ORDER: Tests before implementation. Each test task must be created and
left failing before implementing the corresponding feature.

T001 Setup: Project skeleton (already created)
  - Create/verify directories:
    - `/home/lex/projects/p1-go/cmd/metercli`
    - `/home/lex/projects/p1-go/src/services`
    - `/home/lex/projects/p1-go/src/models`
    - `/home/lex/projects/p1-go/tests`
  - Files to check: `go.mod`, `cmd/metercli/main.go`

T002 [P] Setup: CI and test harness
  - Create GitHub Actions workflow to run `go test ./...` and `gofmt`.
  - Ensure workflow runs on PRs and pushes to `001-build-a-cli` and master.
  - Path: `.github/workflows/go-ci.yml`

T003 [P] Setup: Database migration files
  - Add SQL migration file to create schema and tables (use `migrations/001_create_tables.sql`)
  - Contents: the provided `CREATE TABLE p1.meter_readings` and
    `CREATE TABLE p1.external_readings` statements.
  - Path: `/home/lex/projects/p1-go/migrations/001_create_tables.sql`

T004 [P] Tests: Contract JSON samples
  - Add sample JSON responses based on the provided payload into `specs/001-build-a-cli/contracts/`:
    - `/home/lex/projects/p1-go/specs/001-build-a-cli/contracts/meter_sample.json`
    - `/home/lex/projects/p1-go/specs/001-build-a-cli/contracts/meter_sample_external.json`

T005 [P] Tests: Parser unit tests (failing)
  - Create unit tests that exercise `src/services/parser/` parsing functions using the sample JSON files.
  - Tests should assert correct mapping to a `models.Reading` and detect missing/invalid fields.
  - Path: `/home/lex/projects/p1-go/tests/unit/parser_test.go`

T006 [P] Tests: Contract tests for JSON schema
  - Create contract tests that validate the expected fields and types from the meter endpoint (use `tests/contract/test_meter_contract.go`).
  - The tests should read `specs/001-build-a-cli/contracts/meter_sample.json` and fail initially.

T007 [P] Tests: Integration DB test (failing)
  - Create an integration test that runs against a local Postgres instance (Docker) using the migration SQL and asserts a reading can be inserted and retrieved.
  - Path: `/home/lex/projects/p1-go/tests/integration/test_db_insert.go`
  - Use the provided credentials for local testing via environment variables (DB_HOST, DB_NAME, DB_USER, DB_PASS, DB_PORT).

T008 Implement: Parser implementation
  - Implement `src/services/parser` to parse the meter JSON payload into a `models.Reading` and an array of external readings.
  - Implement timezone/UTC normalization for timestamps; convert `gas_timestamp` (251003101003) to UTC (clarify if epoch-like or custom — if ambiguous, store raw for now and document conversion in data-model.md).
  - Path: `/home/lex/projects/p1-go/src/services/parser/parser.go`

T009 Implement: Models and DB mapping
  - Implement `models.Reading` mapping and `models.ExternalReading` if necessary.
  - Implement DB adapter `src/services/db/postgres.go` with `InsertReading(ctx, models.Reading)` that performs an upsert based on `unique_id` to enforce idempotency.
  - Use parameterized queries; do not interpolate credentials.

T010 Implement: Meter HTTP client
  - Implement `src/services/meter/http_client.go` to perform `GET http://192.168.101.20/api/v1/data` with a configurable timeout and basic error handling.

T011 Implement: CLI orchestration
  - Implement `cmd/metercli/main.go` orchestration to call the meter client, parse results, and insert into Postgres using the adapter.
  - Add flags: `--config`, `--dry-run`, `--verbose`.

T012 Implement: Scheduler and locking
  - Add `scripts/run_every_minute.sh` as an example cron wrapper that runs `cmd/metercli` and acquires a lock (using `flock` or Postgres advisory lock) to prevent overlapping executions.

T013 Implement: External readings insert
  - Insert entries from the `external` array into `p1.external_readings` referencing `meter_readings.unique_id`.

T014 Integration: Full end-to-end test (failing until implemented)
  - Create an integration test that spins up Postgres via docker-compose, runs migrations, runs `cmd/metercli` once against a stubbed HTTP server that returns the sample JSON, and asserts DB rows.
  - Path: `/home/lex/projects/p1-go/tests/integration/test_end_to_end.go`

T015 [P] Observability & Metrics
  - Add structured logging (use standard `log` or a structured logger) and expose Prometheus metrics for read success/failure counts and latencies.
  - Path: `src/services/observability/`

T016 [P] Config & Secrets
  - Implement config loader that reads DB credentials from environment or a `.env` file; document example in `specs/001-build-a-cli/quickstart.md`.

T017 Polish: Quickstart & Cron setup
  - Write `/home/lex/projects/p1-go/specs/001-build-a-cli/quickstart.md` with a sample cron entry and instructions to run migrations and run the CLI.

T018 [P] Performance test
  - Add a simple script to measure read+insert times; verify <5s under normal conditions.

T019 Release
  - Bump version, update CHANGELOG, prepare release notes describing DB schema and migration steps.

## Parallel execution examples

Run these in parallel (no file overlap): T004, T005, T006, T015, T016

## Dependency notes

- T004-T007 (tests) MUST be added and failing before T008-T011 (implementation).
- Migrations (T003) MUST be applied to integration test DB before running integration tests (T007, T014).

---

Paths created/modified by tasks:
- `/home/lex/projects/p1-go/specs/001-build-a-cli/contracts/meter_sample.json`
- `/home/lex/projects/p1-go/migrations/001_create_tables.sql`
- `/home/lex/projects/p1-go/src/services/parser/parser.go` (impl)
- `/home/lex/projects/p1-go/src/services/db/postgres.go` (impl)
- `/home/lex/projects/p1-go/cmd/metercli/` (binary)
- `/home/lex/projects/p1-go/tests/` (unit/contract/integration)
# Tasks: Build a CLI to read gas & electricity meters

**Input**: Design documents from `/specs/001-build-a-cli/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Execution Flow (main)

```
1. Load plan.md from feature directory
   → If not found: ERROR "No implementation plan found"
2. Load optional design documents:
   → data-model.md: Extract entities → model tasks
   → contracts/: Each file → contract test task
   → research.md: Extract decisions → setup tasks
3. Generate tasks by category:
   → Setup: project init, dependencies, linting
   → Tests: contract tests, integration tests
   → Core: meter client, DB client, scheduler
   → Integration: scheduler, buffering, observability
   → Polish: unit tests, performance, docs
4. Apply task rules:
   → Different files = mark [P] for parallel
   → Same file = sequential (no [P])
   → Tests before implementation (TDD)
5. Return: SUCCESS (tasks ready for execution)
```

## Format: `[ID] [P?] Description`

- **[P]**: Can run in parallel (different files, no dependencies)

## Phase 1: Setup

- T001 Initialize Go module and project layout (create `cmd/metercli`, `src/cli`, `src/services`, `src/models`)
- T002 [P] Add dependencies: Postgres driver (`pgx`), HTTP client (std lib), logging/metrics (e.g., zerolog, prometheus client)
- T003 Configure linting (golangci-lint) and formatting (gofmt)
- T004 [P] Add GitHub Actions skeleton for running `go test` (unit + integration)

## Phase 2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE IMPLEMENTATION

**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**

- T005 [P] Unit tests: parser for meter JSON (tests in `tests/unit/test_parser.go`) — create example JSON payloads in `specs/001-build-a-cli/contracts/`.
- T006 [P] Contract tests: validate expected JSON schema for gas meter and electricity meter (`tests/contract/test_gas_contract.go`, `tests/contract/test_elec_contract.go`).
- T007 [P] Integration test: failing test to assert DB insert for a sample reading using test Postgres (`tests/integration/test_db_insert.go`).

## Phase 3: Core Implementation (ONLY after tests are failing)

- T008 [P] Implement parser in `src/services/parser/parser.go` to parse JSON payloads into `models/Reading`.
- T009 [P] Implement meter client in `src/services/meter/http_client.go` to GET JSON from configured endpoints with timeout and retry logic.
- T010 [P] Implement DB adapter in `src/services/db/postgres.go` with an `InsertReading(ctx, Reading)` method and idempotency handling (unique constraint or upsert).
- T011 Implement CLI entrypoint `cmd/metercli/main.go` to load config and call a run function that reads both meters and writes to DB.

## Phase 4: Integration

- T012 Scheduler wrapper: Add `scripts/run_every_minute.sh` or provide sample cron entry; implement lock to avoid overlaps (`flock` or Postgres advisory lock) in `src/services/scheduler/lock.go`.
- T013 Local buffering: Implement optional local SQLite buffer in `src/services/buffer` to queue readings when Postgres is unavailable and drain when available.
- T014 Observability: Add structured logging and Prometheus metrics for read successes/failures, latencies, retry counts.

## Phase 5: Polish

- T015 [P] Unit tests coverage: add unit tests for parser edge cases and DB adapter error handling (`tests/unit/`)
- T016 Performance tests: measure read+insert latency; ensure <5s target under expected hardware.
- T017 [P] Documentation: update `specs/001-build-a-cli/quickstart.md` with cron example and config instructions.
- T018 Release: bump version, update CHANGELOG

## Dependencies

- Tests (T005-T007) must be present and failing before core implementation tasks (T008-T011).

## Notes

- Use environment variables or a configuration file for endpoints and DB connection strings.
- For idempotency use a unique constraint on (`timestamp_utc`, `meter_type`, `source_id`) or a deduplication key.
