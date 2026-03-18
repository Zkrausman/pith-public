# Diet: Token-Optimized CLI Proxy

[![Test and Release](https://github.com/Zkrausman/Diet/actions/workflows/test.yml/badge.svg)](https://github.com/Zkrausman/Diet/actions/workflows/test.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/Zkrausman/Diet/blob/main/LICENSE)

`diet` is a high-performance CLI tool that intercepts terminal command outputs, compresses them, and filters out noise to save tokens for LLMs (like Claude Code, Gemini CLI, etc.).

## 🚀 Features
- **Zero-Overhead Proxy:** Fast Go-based binary for near-instant command interception.
- **Smart Parsers:** 15+ specialized optimizers for common developer commands.
- **Composite Commands:** Intelligent handling of shell-joined commands (e.g., `git status & git log`).
- **Escape Hatch (`diet raw`):** Bypass all parsers when you need the exact, unformatted truth.
- **Middle-Out Truncation:** Automatically keep the start and end of massive outputs, removing the redundant middle.
- **Multi-LLM Integration:** Automated hook setup for Gemini CLI, Claude Code, and Codex.
- **Telemetry & Discovery:** Tracks exact token savings and identifies new optimization targets.
- **Interactive Configuration:** Toggle parsers and adjust truncation limits via `diet config`.
- **Interactive Dashboard:** Launch a local web-based analytics dashboard via `diet dashboard`.

## 📦 Installation
1. Download the latest `diet.exe` from [Releases](https://github.com/Zkrausman/Diet/releases).
2. Run the installer (defaults to global install for all CLIs):
   ```bash
   ./diet.exe install
   ```
   *Note: Diet safely **merges** hooks into your existing `settings.json` files and automatically creates a `.bak` backup before any modification. It will NOT overwrite your existing configurations or other hooks.*
3. Restart your terminal.

---

## 🛠️ Power User Features

### 1. Targeted Installation
If you only want to install hooks for a specific CLI or only for the current project:
- **Local Install (current folder):** `diet install --gemini`
- **Global Install (specific CLI):** `diet install --claude --global`
- **Full Global Sync:** Use `diet install` to ensure all system-wide `settings.json` files are up-to-date.

### 2. The Escape Hatch (`diet raw`)
If a parser is being too aggressive and you need to see the raw, bit-for-bit output of a command, prefix it with `raw`:
```bash
diet raw git diff
```
This bypasses all logic and returns the original system output.

### 2. Middle-Out Truncation
When a command returns thousands of lines (like a massive log file), Diet prevents context overflow by keeping the most important parts: the **beginning** (setup/context) and the **end** (errors/results).

**Configurable via `diet config`:**
- `MaxLines`: Total lines allowed before truncation kicks in (default: 500).
- `HeadLines`: Lines to preserve at the top (default: 100).
- `TailLines`: Lines to preserve at the bottom (default: 100).

---

## 🛠️ Supported Parsers
Diet currently optimizes the following commands:
- **Git:** `status`, `log`, `diff`, `branch`, and **composite commands** (e.g., `git status & git log`)
- **Filesystem:** `ls`, `dir`, `find`, `tree`, `du`
- **Text:** `grep`, `rg`, `cat` (JSON/XML/CSS/HTML minification)
- **Infra/Dev:** `docker ps`, `npm list`, `pip list`, `npm test`, `go test`, `pytest`, `env`, `set`

---

## 📊 Compression Examples

### 1. `git status`
**Raw Output (~80 tokens):**
```text
On branch main
Your branch is up to date with 'origin/main'.

Changes to be committed:
  (use "git restore --staged <file>..." to unstage)
	new file:   main.go

Changes not staged for commit:
  (use "git add <file>..." to update what will be committed)
  (use "git restore <file>..." to discard changes in working directory)
	modified:   pkg/parser/git.go
```

**Diet Output (~25 tokens):**
```text
new file:   main.go
modified:   pkg/parser/git.go
```

### 2. `git log`
**Raw Output (~120 tokens):**
```text
commit 3f7d7f7b12345678
Author: Zachary Krausman <zkrausman@gmail.com>
Date:   Sun Mar 15 22:49:38 2026 -0400

    Update: Display changelog during diet update

commit 1cd63311f8330490
Author: Zachary Krausman <zkrausman@gmail.com>
Date:   Sun Mar 15 22:30:00 2026 -0400

    Initial commit
```

**Diet Output (~35 tokens):**
```text
3f7d7f7 | Zachary Krausman | Mar 15 2026 | Update: Display changelog
1cd6331 | Zachary Krausman | Mar 15 2026 | Initial commit
```

### 3. `ls -l`
**Raw Output (~50 tokens):**
```text
total 8
-rw-r--r--  1 user  group  123 Mar 15 22:30 main.go
-rw-r--r--  1 user  group  456 Mar 15 22:31 README.md
```

**Diet Output (~10 tokens):**
```text
main.go README.md
```

---

## 🧪 Testing & Quality
Diet uses [promptfoo](https://promptfoo.dev/) to evaluate the performance and quality of its token-optimization strategies.

- **Status:** **Work In Progress (WIP)** 🛠️
- **Latest Evaluation:** [tests/report.md](tests/report.md)

## 📊 Analytics Dashboard
Diet includes an interactive web dashboard to visualize your token savings and discover new optimization opportunities.

- **Launch:** `diet dashboard`
- **Visualization:** Powered by Chart.js for real-time analytics.
- **Ecosystem:** Diet is **[BuildABoard](https://github.com/Zkrausman/BuildABoard) compatible**. The included `board.json` allows you to seamlessly integrate your Diet analytics into the broader BuildABoard dashboard ecosystem.

---

## ⌨️ CLI Usage
- `diet <command>`: Proxy any command manually.
- `diet config`: Toggle parsers on/off.
- `diet gain`: View total tokens saved.
- `diet discover`: See unparsed commands costing you tokens.
- `diet reset --all`: Clear all telemetry data.
- `diet update`: Auto-update to the latest version.
