# Bead Self-Documentation Rubric v1 (`bd-3qwe.1`)

## Purpose
Define the normative self-documentation contract for active beads so planning artifacts stay self-contained, reviewable, and lintable.

## Required Comment Anchors
All implementation and gate beads must carry two machine-readable comment anchors.

### Rationale Anchor
```text
[bd-3qwe self-doc] RATIONALE
PROBLEM_CONTEXT: ...
WHY_NOW: ...
SCOPE_BOUNDARY: ...
PRIMARY_SURFACES: ...
```

### Evidence Anchor
```text
[bd-3qwe self-doc] EVIDENCE
UNIT_TESTS: ...
INTEGRATION_TESTS: ...
E2E_TESTS: ...
COMMIT: ...
PERFORMANCE_VALIDATION: ...
LOGGING_ARTIFACTS: ...
```

`N/A` is allowed only with explicit reason text.

## Machine-Checkable Rule IDs
1. `SDOC-RUBRIC-001` missing rationale anchor.
2. `SDOC-RUBRIC-002` missing rationale fields.
3. `SDOC-RUBRIC-003` missing evidence anchor.
4. `SDOC-RUBRIC-004` missing evidence fields.
5. `SDOC-RUBRIC-005` malformed exception payload.
