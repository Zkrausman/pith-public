# Mandatory Project Instructions

This project adheres to the **GraphRAG-Aware Architectural Standard** and the **Beads Closure Mandate**.

## Research Protocol
- AI agents **MUST** use `thneed query` as their primary research tool before proposing any major bug fixes or architectural refactors.
- Consult the `thneed-navigator` skill for specific query depths and impact analysis.

## Documentation Protocol
- **ALL WORK** must be performed on a Beads task.
- Before closing a task, you **MUST** run `thneed beads-scaffold <ID>` to infer technical metadata.
- Task closure notes must follow the `bd-3qwe.1` rubric:
    - `[bd-3qwe self-doc] RATIONALE: ...` (including Problem Context and Why Now)
    - `[bd-3qwe self-doc] EVIDENCE: ...` (including Commit SHA)

## Compliance Enforcement
- This repository uses a Git `pre-commit` hook managed by **Synapse**.
- Non-compliant commits (missing package READMEs or under-documented Beads tasks) will be automatically blocked.
