# Diet: Token-Optimized CLI Proxy

`diet` is a high-performance CLI tool that intercepts terminal command outputs, compresses them, and filters out noise to save tokens for LLMs (like Claude Code, OpenCode, etc.).

## Features
- **Zero-Overhead Proxy:** Fast Go-based binary for near-instant command interception.
- **Smart Parsers:** Specialized logic for commands like `git status` to strip boilerplate.
- **Telemetry & Discovery:** Tracks token usage and identifies common commands that need dedicated parsers.
- **Interactive Configuration:** Simple CLI for toggling parsers and setting token limits.

## Installation
Build from source:
```bash
go build -o diet
```

## Usage
Proxy any command:
```bash
diet git status
diet whoami
```

Manage Diet:
- `diet config`: Interactive configuration.
- `diet gain`: View token savings report.
- `diet discover`: Identify new parser opportunities.

## Project Structure
- `pkg/telemetry`: SQLite-based logging for token usage and discovery.
- `pkg/parser`: Interface and implementations for command-specific optimizers.
- `pkg/runner`: Core execution engine with token estimation.
- `pkg/config`: Configuration management using JSON and interactive prompts.
