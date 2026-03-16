# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.4.9] - 2026-03-16

### Improved
- Overhauled `install` command to copy Diet binary to a permanent location (`~/.diet/bin`) and add that to the system `PATH`.
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
