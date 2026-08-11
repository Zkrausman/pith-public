# Agent Onboarding - Pith

Pith is a high-performance Go CLI proxy that compresses command output before returning it to an LLM caller.

## Work Tracking and Knowledge

- **Linear is the authoritative issue tracker.** Create, update, and close Pith work in the `GeneralAiDev` team (AIDEV issues).
- **The LLM Wiki is the durable knowledge base.** Recall relevant context at task start; record meaningful decisions, discoveries, and completions when work ends.
- **GitHub pull requests** are the review and merge workflow.
- Do not use Beads (`bd`) or Dolt for Pith task tracking.

## Landing the Plane

Before ending a code-changing session:

1. Create Linear issues for remaining work and update the completed issue.
2. Run applicable tests, linters, and builds.
3. Commit work and open/update a GitHub PR when appropriate.
4. Push changes: `git pull --rebase`, `git push`, then verify `git status` is up to date with origin.
5. Record the handoff and durable insights in the wiki.

## Build & Test

- Binary: `pith.exe`
- Build: `go build -o pith.exe main.go`
- Test: `go test ./...`

## Conventions

- Parsers implement `pkg/parser/interface.go`.
- Telemetry is stored in `~/.pith/pith.db`.
- Pith integrates with LLM CLIs through hook configuration in their `settings.json` files.
