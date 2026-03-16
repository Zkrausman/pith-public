# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
