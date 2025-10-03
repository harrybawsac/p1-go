# Implementation Plan: Build a CLI to read gas & electricity meters

**Branch**: `001-build-a-cli` | **Date**: 2025-10-03 | **Spec**: ../001-build-a-cli/spec.md
**Input**: Feature specification from `/specs/001-build-a-cli/spec.md`

## Execution Flow (/plan command scope)

```
1. Load feature spec from Input path
   → If not found: ERROR "No feature spec at {path}"
2. Fill Technical Context (scan for NEEDS CLARIFICATION)
   → Detect Project Type: single / CLI-focused service
   → Set Structure Decision: single project, CLI in `src/cli/`, services in `src/services/`
3. Fill the Constitution Check section based on the content of the constitution document.
4. Evaluate Constitution Check section below
   → If violations exist: Document in Complexity Tracking
   → If no justification possible: ERROR "Simplify approach first"
   → Update Progress Tracking: Initial Constitution Check
5. Execute Phase 0 → research.md
   → If NEEDS CLARIFICATION remain: ERROR "Resolve unknowns"
6. Execute Phase 1 → contracts, data-model.md, quickstart.md, agent-specific file
7. Re-evaluate Constitution Check section
   → If new violations: Refactor design, return to Phase 1
   → Update Progress Tracking: Post-Design Constitution Check
8. Plan Phase 2 → Describe task generation approach (DO NOT create tasks.md)
9. STOP - Ready for /tasks command
```

## Summary

Build a small, reliable CLI that performs an HTTP GET to each meter's JSON
endpoint (gas and electricity), parses readings, and writes timestamped rows to
a PostgreSQL database. The CLI will be scheduled to run every minute (cron or
systemd timer). Key concerns: idempotency, local buffering if DB is down,
observability and retry/backoff behavior.

## Technical Context

**Language/Version**: Go (inferred from repository name `p1-go`) — [NEEDS CLARIFICATION if you prefer another language]
**Primary Dependencies**: HTTP client (std lib), PostgreSQL driver (e.g., lib/pq or pgx), logging/metrics lib
**Storage**: PostgreSQL (primary); local durable buffer (SQLite) optional for resilience
**Testing**: go test (unit), integration tests with a local Postgres instance or testcontainers
**Target Platform**: Linux server (Raspberry Pi or small VPS), cron/systemd timer
**Project Type**: single CLI/service
**Performance Goals**: read + parse + insert < 5s (target), support one-minute cadence without overlap
**Constraints**: Avoid long-running reads that exceed the 60s schedule; implement a locking mechanism to prevent overlapping runs
**Scale/Scope**: Single household meter readings (2 endpoints), ~1,440 records/day per meter (~525,600/year)

## Constitution Check

Initial evaluation:

- Testability: MUST provide unit tests for parsing logic and contract tests for the HTTP schema; integration tests for DB writes.
- Performance: Target cycle <5s; validate with a quick benchmark during Phase 1.
- UX Consistency: CLI flags MUST follow existing conventions (e.g., `--config`, `--dry-run`, `--verbose`).
- Observability: Structured logs, metrics for read successes/failures, and retry counts required.

Status: PASS (no gate blockers identified). See Phase 1 for test fixtures and contract tests.

## Project Structure

### Documentation (this feature)

```
specs/001-build-a-cli/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
└── tasks.md
```

### Source Code (repository root)

```
src/
├── cli/                # CLI entrypoint and command handling
├── services/           # meter client, db client, buffering
├── models/             # data models (Reading)
└── cmd/                # small binaries if needed

tests/
├── contract/
├── integration/
└── unit/
```

**Structure Decision**: Single project with `src/cli` and `src/services`. Use Go modules and place main binary in `cmd/metercli` if preferred.

## Phase 0: Outline & Research

Extract unknowns and research tasks:

- Confirm meter endpoint formats and authentication (if any).
- Best practices for idempotent inserts in Postgres (unique constraints, upserts).
- Locking approach to prevent overlapping scheduled runs on single-host (file locks, flock, or advisory locks in Postgres).

Research outputs (to produce `research.md`): Decision, Rationale, Alternatives.

## Phase 1: Design & Contracts

1. **Data model** (`data-model.md`): `reading` table with columns: id,
   timestamp_utc, meter_type, value, unit, status, source_id, correlation_id
2. **Contracts**: Define the expected JSON schema returned by the meters and
   generate contract tests (one per meter endpoint).
3. **Contract tests**: Create tests that validate parsing logic against sample
   JSON payloads; ensure tests fail before implementation.
4. **Integration tests**: Use a test Postgres instance to assert rows are
   inserted with correct fields.

## Phase 2: Task Planning Approach

- Use `.specify/templates/tasks-template.md` as base.
- Key task groups:
  - Setup: project layout, module init, linting
  - Tests: parsing unit tests, contract tests for JSON schema, integration tests for DB
  - Core: Implement meter HTTP client, parsing, DB adapter
  - Integration: scheduler wrapper, locking, buffering
  - Polish: docs, quickstart, performance tests

Ordering: Contract/tests first, then core implementation (TDD).

## Complexity Tracking

- If meter endpoints require auth or non-trivial handshake, complexity increases and may require secure credential handling.

## Progress Tracking

- [ ] Phase 0: Research complete
- [ ] Phase 1: Design complete
- [ ] Phase 2: Task planning complete

---

_Based on Constitution v1.0.0 - See `/memory/constitution.md`_
