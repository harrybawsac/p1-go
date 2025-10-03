Constitution & Versioning

This short README explains how to update the constitution dates and the
recommended workflow for amendments.

Files involved:

- `.specify/memory/constitution.md` — The canonical constitution document.
- `.specify/templates/*` — Templates that may reference constitution-dependent
  gates and version strings.
- `.specify/tools/update-constitution-dates.sh` — Small helper to update
  ratification and last-amended dates in the constitution file.

Versioning rules (summary):

- MAJOR: Backwards-incompatible governance or principle removals/redefinitions.
- MINOR: New principle/section added or materially expanded guidance.
- PATCH: Clarifications, wording, typo fixes, non-semantic refinements.

How to set dates:

- To set the ratification date:

```bash
.specify/tools/update-constitution-dates.sh --ratified 2025-06-13
```

- To set the last amended date:

```bash
.specify/tools/update-constitution-dates.sh --amended 2025-10-03
```

- To set both at once:

```bash
.specify/tools/update-constitution-dates.sh --ratified 2025-06-13 --amended 2025-10-03
```

Amendment process:

1. Propose amendment in a branch and PR with reasons and intended version bump.
2. Update the constitution file and templates as needed.
3. Run the date update script to set `Last Amended` to the amendment date.
4. Merge the PR; update the plan/template references if the version changed.

Manual follow-ups:

- Review templates in `.specify/templates/` to ensure per-feature placeholders
  remain intentional.
