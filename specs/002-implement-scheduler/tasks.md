```markdown
# Tasks: Implement scheduler-run orchestration and buffer drain

Phases: Setup, Tests, Core, Integration, Polish

1. Setup: repo readiness [P]

   - Ensure `migrations/` present and DB config example in repo

2. Tests: Parser contract tests (created)

   - Implemented `tests/contract/parse_contract_test.go` which validates `ParseFullReading` accepts an example payload.

3. Core: Implement runOnce orchestration (TDD) [COMPLETED]

- Test: `tests/unit/runonce_no_endpoint_test.go` (added) asserts runOnce returns error when `METER_ENDPOINT` unset.
- Implemented: `src/app/runner.go` and `cmd/metercli` now perform fetch -> parse -> insert and buffer on DB failure.

4. Core: Buffer drain CLI [COMPLETED]

- Test: `tests/integration/drain_buffer_test.go` (added) verifies `buffer.Drain` invokes persist function.
- Implemented: `--drain-buffer` flag in `cmd/metercli` which calls `buffer.Drain` and attempts to persist entries.

5. Integration: Run scheduler loop end-to-end (manual)

   - Manual verification steps in `quickstart.md`.

6. Polish: Update docs and mark tasks complete

Notes:

- Each test must be added before implementing the corresponding code (TDD).
- Parallel tasks [P] may be implemented in any order if independent.
```
