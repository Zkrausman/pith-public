# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.7.0] - 2026-03-19

### Added
- **Dashboard v2.0:**
    - **Bento Grid Layout:** Completely overhauled the analytics interface with a modern, responsive grid of cards.
    - **Efficiency Gauge:** Interactive SVG ring showing the overall compression ratio at a glance.
    - **Live Activity Stream:** Real-time ticker showing recently executed commands, including parser status (Hit/Miss) and savings per command.
    - **Pith vs. Raw Deep Dive:** Side-by-side comparison modal that allows users to verify Pith's optimization by viewing the raw output alongside the compressed version.
- **Telemetry Improvements:**
    - Added support for recording full original and compressed output content (required for the Deep Dive feature).
    - Database migration to add `original_content` and `compressed_content` columns to the `executions` table.
- **GUI Endpoints:**
    - New `/api/recent` and `/api/execution` API endpoints to power the live dashboard and detailed views.

## [v0.6.0] - 2026-03-19

### Added
- **New Parsers:**
    - `WebParser`: Optimizes `curl`, `wget`, and `Invoke-WebRequest` by minifying JSON and extracting HTML titles.
    - `PithParser`: Internal optimizer for Pith's own help and log outputs.
    - `PowerShellParser`: Cleans up verbose Windows PowerShell and CMD headers while preserving file details.
    - `GoParser`: Targeted optimization for `go version`, `go build`, and module listings.

## [v0.5.9] - 2026-03-18

### Changed
- **Rebrand:** Renamed the project from **Diet** to **Pith** (skip the peel, get to the pith).
- **Module Update:** Updated Go module name to `pith` and updated all internal imports.
- **Paths:** Changed default configuration and telemetry directory from `~/.diet` to `~/.pith`.
- **Log Files:** Renamed `diet.log` to `pith.log` and `diet.db` to `pith.db`.
- **Hooks:** Updated CLI hooks to use the new `pith` command and `pith-optimizer` name.

## [v0.5.8] - 2026-03-17

### Improved
- **Intelligent Hook Installation:** The `pith install` command now defaults to **global** installation for all supported CLIs (Gemini, Claude, and Codex) if no flags are provided.
- **JSON Configuration Merging:** Installer now safely merges Pith hooks into existing `settings.json` files instead of skipping them. It correctly preserves existing hooks (like `bd prime`) and custom settings (like `security` or `ui`).
- **Safety Backups:** Automatically creates a `settings.json.bak` backup before any modification.
- **Duplicate Prevention:** Improved logic to detect and avoid duplicate hook entries during installation.

## [v0.5.7] - 2026-03-17

### Added
- **Hook Support:**
    - Support for **Claude Code** CLI via `.claude/settings.json` (using `PostToolUse` event).
    - Support for **Codex** (and compatible agents) via `.codex/settings.json` (using `AfterTool` event).
    - New `--claude` and `--codex` flags for `pith install`.
    - Generic `setupHook` internal helper for easier future hook additions.

## [v0.5.6] - 2026-03-16

### Added
- **New Parsers:**
    - `SourceParser`: Optimizes `cat` and `type` output by stripping comments and excessive whitespace while preserving code structure.
    - `GitHubReleaseParser`: Optimizes `gh release view/list` by extracting key metadata (Tag, Title, Assets) and removing redundant SHAs.
    - `ChainParser`: Enables optimization of shell-chained commands (e.g., `cmd1; cmd2 && cmd3`) by recursively matching individual parts.
- New test cases and `promptfoo` evaluations for the new parsers.

### Fixed
- Fixed `dashboard` command being incorrectly proxied as a target command.

## [v0.5.5] - 2026-03-16

### Fixed
- Fixed `pith update` to correctly identify and download platform-specific binaries (e.g., `pith-windows-amd64.exe`) from GitHub releases.

## [v0.5.4] - 2026-03-16

### Added
- **Interactive Dashboard:** Introduced `pith dashboard`, a web-based GUI that provides a visual overview of token savings, command breakdowns, and opportunity discovery using Chart.js.
- New `GetStatsByDay` method in telemetry to support historical analysis in the dashboard.

## [v0.5.3] - 2026-03-16

### Added
- Added a GitHub Actions workflow to automate Go testing on all pushes and pull requests.
- Integrated automated cross-platform releases; pushing a version tag now triggers a build for Windows, Linux, and macOS (amd64 and arm64) and creates a GitHub Release with the binaries attached.

## [v0.5.2] - 2026-03-16

### Improved
- Improved the `gain` and `discover` commands to display helpful empty-state messages when no telemetry data has been recorded yet, rather than printing empty tables.

## [v0.5.1] - 2026-03-16

### Fixed
- Fixed an issue where commands intercepted by the Gemini CLI hook were not being logged to the Snag log. Pith now correctly parses exit codes from the hook payload to ensure accurate behavioral tracking.

## [v0.5.0] - 2026-03-16

### Added
- **Snag Integration:** Pith now actively exports a plain text log file specifically formatted for the [Snag](https://github.com/zkrau/snag) behavioral learning tool. Every intercepted command, its output (truncated to the last 50 lines to prevent bloat), and its exit code are appended to `~/.pith/pith.log` using the `[CMD]` and `[EXIT]` syntax required by Snag's collector.

## [v0.4.9] - 2026-03-16

### Improved
- Overhauled `install` command to copy Pith binary to a permanent location (`~/.pith/bin`) and add that to the system `PATH`.
- Updated CLI hooks to use the permanent binary location, making them resilient to source directory changes.

## [v0.4.8] - 2026-03-16

### Fixed
- Fixed Gemini CLI hook to correctly handle `Output: ` and `Error: ` prefixes in tool responses, ensuring parsers like `git_status` work correctly when called via hooks.

### Changed
- Updated `promptfooconfig.yaml` with more comprehensive evaluation cases for all major parsers.
- Updated `tests/prompts.txt` to be more flexible for different evaluation tasks.

## [v0.4.7] - 2026-03-16

### Fixed
- Improved `GitHubParser` to correctly handle `gh repo list` and other list variants.
- Fixed `Runner` to correctly split command arguments when proxied.
- Updated `Runner` to capture combined stdout/stderr for more reliable parsing.

## [v0.4.6] - 2026-03-16

### Added
- New `GitHubParser` to optimize `gh issue list`, `gh pr list`, and `gh release list` outputs.

### Fixed
- Improved `LsParser` reliability across different `ls -l` and `dir` variants.

## [v0.4.5] - 2026-03-16

### Improved
- Balanced `LsParser` and `MinifyParser` to preserve critical context for LLMs while still saving tokens, specifically addressing issues seen in other CLI optimizers (RTK-AI/rtk#582).
- `LsParser` now keeps file mode and size when available.
- `MinifyParser` now strips comments but preserves line structure.

## [v0.4.4] - 2026-03-16

### Added
- New `CompositeGitParser` to optimize shell-joined Git commands (e.g., `git status & git log`).
- Support for shell execution in the runner, enabling optimization of complex command strings.

### Improved
- Updated `README.md` with new features and updated parser counts.

## [v0.4.3] - 2026-03-16

### Fixed
- Updated `_hook` command to handle the latest Gemini CLI hook payload schema (specifically supporting the `tool_input` object instead of the older `tool_call_request.arguments` JSON string).
- Improved robustness of hook command to support multiple schema versions.

### Added
- Improved `.gitignore` to include `.beads/` and `test_hook_input.json`.

## [v0.4.2] - 2026-03-16

### Added
- Initial support for token-optimized CLI proxying.
- Parsers for `git status`, `git log`, `ls`, and others.
- Middle-out truncation for large outputs.
