# Project Instructions for AI Agents

This file provides instructions and context for AI coding agents working on this project.

## Thneed-Aware Architectural Standard
- **Mandate**: All projects within 'E:\Repos' MUST have a README.md in every package containing a Mermaid diagram.
- **Purpose**: High-signal conceptual indexing for Thneed GraphRAG.
- **Skill**: Use 'thneed-navigator' for all codebase research.
- **Diagrams**: Mermaid diagrams MUST use descriptive labels (e.g., NodeID[path/to/file.go]) to optimize Thneed lexical search.


Before any structural code changes or major bug fixes, agents MUST:
1. **Search Thneed**: Use `thneed query --text "<task_description>" --depth 2 --robots` to retrieve task-relevant context.
2. Verify that the context returned by Thneed is sufficient for the task.
3. **Audit Signal**: Look for `--- THNEED-POWERED SEARCH RESULTS ---` in the output to confirm you are using high-signal context.
4. If search relevance is poor, file a task to improve Thneed's indexing for that area.


## Work Tracking and Knowledge

- Use **Linear** (`GeneralAiDev` / AIDEV issues) as the authoritative Pith issue tracker.
- Use the **LLM Wiki** for durable architecture, decisions, discoveries, and handoffs.
- Use GitHub pull requests for review and merge.
- Do not use Beads (`bd`) or Dolt for Pith work.

## Session Completion

Run quality gates, update the Linear issue, push (`git pull --rebase`, `git push`), verify the branch is current, and record the handoff in the wiki.

## Build & Test

_Add your build and test commands here_

```bash
# Example:
# npm install
# npm test
```

## Architecture Overview

_Add a brief overview of your project architecture_

## Conventions & Patterns

## Communication & Reasoning

- **Priority on Reasoning**: Following the insights from the Anthropic April 23 post-mortem, this project prioritizes **comprehensive technical reasoning** over brevity.
- **Override**: The "Minimal Output" and "No Chitchat" global constraints are explicitly relaxed for this repository. Agents SHOULD provide detailed technical rationale and thorough analysis when solving complex problems.
- **Exhaustive Planning**: Always favor explicit, detailed planning (e.g., using Plan Mode or thorough Thought blocks) to ensure architectural integrity is maintained during optimizations.

_Add your project-specific conventions here_
