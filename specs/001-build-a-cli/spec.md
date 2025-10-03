# Feature Specification: Build a CLI to read gas & electricity meters

**Feature Branch**: `001-build-a-cli`  
**Created**: 2025-10-03  
**Status**: Draft  
**Input**: User description: "Build a CLI that can get information from my gas and electricity meter. I want to be able to run the program every minute to get the data. The data should also be stored in a database."

## Execution Flow (main)

```
1. Parse user description from Input
   → If empty: ERROR "No feature description provided"
2. Extract key concepts from description
   → Identify: actors, actions, data, constraints
3. For each unclear aspect:
   → Mark with [NEEDS CLARIFICATION: specific question]
4. Fill User Scenarios & Testing section
   → If no clear user flow: ERROR "Cannot determine user scenarios"
5. Generate Functional Requirements
   → Each requirement must be testable
   → Mark ambiguous requirements
6. Identify Key Entities (if data involved)
7. Run Review Checklist
   → If any [NEEDS CLARIFICATION]: WARN "Spec has uncertainties"
   → If implementation details found: ERROR "Remove tech details"
8. Return: SUCCESS (spec ready for planning)
```

---

## ⚡ Quick Guidelines

- ✅ Focus on WHAT users need and WHY
- ❌ Avoid HOW to implement (no tech stack, APIs, code structure)
- 👥 Written for business stakeholders, not developers

## Constitution Check

This feature interacts with hardware and is scheduled frequently; ensure the
spec includes testability (automated tests or simulated inputs), performance
targets (e.g., acceptable latency and DB ingestion rate), and observability
requirements (logging, error metrics).

### Section Requirements

- Mandatory sections: User Scenarios & Testing, Requirements, Key Entities,
  Performance & Operational Considerations.

## User Scenarios & Testing _(mandatory)_

### Primary User Story

As a homeowner, I want to run a CLI every minute that reads my gas and
electricity meters so that consumption can be recorded and analyzed over time.

### Acceptance Scenarios

1. **Given** the meter is reachable and responding, **When** the CLI runs,
   **Then** it MUST read the latest values for gas and electricity and store a
   timestamped record in the database.
2. **Given** the meter is unreachable or times out, **When** the CLI runs,
   **Then** it MUST log a clear, actionable error and NOT insert corrupted
   data; transient errors MAY be retried with exponential backoff.

### Edge Cases

- Meter returns partial data: the CLI MUST validate fields and only store
  complete, validated records, or store partial records with a 'partial' flag
  and associated rationale.
- Network glitches: the CLI MUST avoid duplicate inserts when retries occur
  (idempotency or deduplication strategy required).
- Clock skew: timestamps MUST be normalized to UTC before storage.

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: The CLI MUST be able to read current gas meter readings.
- **FR-002**: The CLI MUST be able to read current electricity meter readings.
- **FR-003**: The CLI MUST run reliably on a cron-like schedule every minute.
- **FR-004**: Each successful read MUST result in a timestamped record stored
  in the configured database.
- **FR-005**: On failure to read meters, the CLI MUST log the failure and
  surface metrics for monitoring; it MUST NOT write invalid readings to the DB.
- **FR-006**: The system MUST provide idempotency to avoid duplicate entries on
  retries.

_Ambiguities / NEEDS CLARIFICATION:_

- **FR-007**: Meter interface/protocol is unspecified — [NEEDS CLARIFICATION:
  what model(s) and communication protocol(s) (e.g., serial, optical, M-Bus,
  Zigbee, Wi-SUN, HTTP API) will the meters use?]
- **FR-008**: Database choice and retention policy — [NEEDS CLARIFICATION:
  which DB (SQLite/Postgres/etc.), retention period, and indexing requirements?]

### Non-functional Requirements

- **NFR-001 (Performance)**: The CLI read + DB insert cycle SHOULD complete in
  under 5 seconds on typical local hardware. If the meter read takes longer,
  document timeout and retry behavior.
- **NFR-002 (Storage)**: The system MUST support sustained inserts every
  minute without data loss; estimate storage growth based on retention policy.
- **NFR-003 (Observability)**: The program MUST emit structured logs and
  metrics for successes, failures, latencies, and retry counts.

### Key Entities _(include if feature involves data)_

- **Reading**: timestamp (UTC), meter_type (gas|electricity), value, unit,
  status (ok|partial|failed), source_id (if available), correlation_id

---

## Performance & Operational Considerations

- Scheduling: At-scale cron every minute is sensitive to overlaps; ensure
  execution time < schedule period or use a lock mechanism to prevent
  concurrent runs.
- Resilience: Implement retry and backoff for transient network/meter faults.
- Local buffering: Consider a small local queue or durable buffer (e.g., a
  local SQLite DB) to avoid data loss when the central DB is temporarily
  unavailable.

## Review & Acceptance Checklist

- [ ] No implementation details that are prescriptive (e.g., exact library
      choices) remain in the spec unless justified.
- [ ] All FRs have testable acceptance criteria.
- [ ] All [NEEDS CLARIFICATION] items resolved before implementation.

## Execution Status

- [ ] User description parsed
- [ ] Key concepts extracted
- [ ] Ambiguities marked
- [ ] User scenarios defined
- [ ] Requirements generated
- [ ] Entities identified
- [ ] Review checklist passed

# Feature Specification: [FEATURE NAME]

**Feature Branch**: `[###-feature-name]`  
**Created**: [DATE]  
**Status**: Draft  
**Input**: User description: "$ARGUMENTS"

## Execution Flow (main)

```
1. Parse user description from Input
   → If empty: ERROR "No feature description provided"
2. Extract key concepts from description
   → Identify: actors, actions, data, constraints
3. For each unclear aspect:
   → Mark with [NEEDS CLARIFICATION: specific question]
4. Fill User Scenarios & Testing section
   → If no clear user flow: ERROR "Cannot determine user scenarios"
5. Generate Functional Requirements
   → Each requirement must be testable
   → Mark ambiguous requirements
6. Identify Key Entities (if data involved)
7. Run Review Checklist
   → If any [NEEDS CLARIFICATION]: WARN "Spec has uncertainties"
   → If implementation details found: ERROR "Remove tech details"
8. Return: SUCCESS (spec ready for planning)
```

---

## ⚡ Quick Guidelines

- ✅ Focus on WHAT users need and WHY
- ❌ Avoid HOW to implement (no tech stack, APIs, code structure)
- 👥 Written for business stakeholders, not developers

## Constitution Check

Before a spec is considered ready, ensure it aligns with the project's
constitution. At minimum the spec MUST include:

- Testability: Each functional requirement MUST map to at least one testable
  acceptance criterion (unit/contract/integration).
- Performance Targets: If the feature has operational cost or latency
  sensitivity, the spec MUST include measurable targets (e.g., p95 latency,
  throughput, memory budget) and how they will be validated.
- UX Consistency: Public behaviors (APIs, CLI flags, UX flows) MUST be
  consistent with existing interfaces or include a documented migration plan
  for compatibility-breaking changes.
- Security & Observability: Security-sensitive features MUST include
  threat-analysis notes and required automated checks; instrumentation plans
  for logging/metrics MUST be specified where applicable.

If any of the above cannot be specified, mark the related items as
[NEEDS CLARIFICATION: ...] and resolve during Phase 0 research.

### Section Requirements

- **Mandatory sections**: Must be completed for every feature
- **Optional sections**: Include only when relevant to the feature
- When a section doesn't apply, remove it entirely (don't leave as "N/A")

### For AI Generation

When creating this spec from a user prompt:

1. **Mark all ambiguities**: Use [NEEDS CLARIFICATION: specific question] for any assumption you'd need to make
2. **Don't guess**: If the prompt doesn't specify something (e.g., "login system" without auth method), mark it
3. **Think like a tester**: Every vague requirement should fail the "testable and unambiguous" checklist item
4. **Common underspecified areas**:
   - User types and permissions
   - Data retention/deletion policies
   - Performance targets and scale
   - Error handling behaviors
   - Integration requirements
   - Security/compliance needs

---

## User Scenarios & Testing _(mandatory)_

### Primary User Story

[Describe the main user journey in plain language]

### Acceptance Scenarios

1. **Given** [initial state], **When** [action], **Then** [expected outcome]
2. **Given** [initial state], **When** [action], **Then** [expected outcome]

### Edge Cases

- What happens when [boundary condition]?
- How does system handle [error scenario]?

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: System MUST [specific capability, e.g., "allow users to create accounts"]
- **FR-002**: System MUST [specific capability, e.g., "validate email addresses"]
- **FR-003**: Users MUST be able to [key interaction, e.g., "reset their password"]
- **FR-004**: System MUST [data requirement, e.g., "persist user preferences"]
- **FR-005**: System MUST [behavior, e.g., "log all security events"]

_Example of marking unclear requirements:_

- **FR-006**: System MUST authenticate users via [NEEDS CLARIFICATION: auth method not specified - email/password, SSO, OAuth?]
- **FR-007**: System MUST retain user data for [NEEDS CLARIFICATION: retention period not specified]

### Key Entities _(include if feature involves data)_

- **[Entity 1]**: [What it represents, key attributes without implementation]
- **[Entity 2]**: [What it represents, relationships to other entities]

---

## Review & Acceptance Checklist

_GATE: Automated checks run during main() execution_

### Content Quality

- [ ] No implementation details (languages, frameworks, APIs)
- [ ] Focused on user value and business needs
- [ ] Written for non-technical stakeholders
- [ ] All mandatory sections completed

### Requirement Completeness

- [ ] No [NEEDS CLARIFICATION] markers remain
- [ ] Requirements are testable and unambiguous
- [ ] Success criteria are measurable
- [ ] Scope is clearly bounded
- [ ] Dependencies and assumptions identified

---

## Execution Status

_Updated by main() during processing_

- [ ] User description parsed
- [ ] Key concepts extracted
- [ ] Ambiguities marked
- [ ] User scenarios defined
- [ ] Requirements generated
- [ ] Entities identified
- [ ] Review checklist passed

---
