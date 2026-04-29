# Pith: Token-Optimized CLI Proxy
> Skip the peel, get to the pith.

[![Test and Release](https://github.com/Zkrausman/Pith/actions/workflows/test.yml/badge.svg)](https://github.com/Zkrausman/Pith/actions/workflows/test.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/Zkrausman/Pith/blob/main/LICENSE)

`pith` is a high-performance CLI tool that intercepts terminal command outputs, compresses them, and filters out noise to save tokens for LLMs (like Claude Code, Gemini CLI, etc.).

## 🚀 Features
- **Zero-Overhead Proxy:** Fast Go-based binary for near-instant command interception.
- **Smart Parsers:** 20+ specialized optimizers for common developer commands, now featuring **High Fidelity Parsing** for source code and test results.
- **Reasoning Over Brevity Standard:** Built-in mandate that prioritizes comprehensive technical depth over brevity to prevent LLM intelligence loss.
- **Configurable Storage:** Centralized storage for telemetry and configuration at `E:\TheBrain\PithBackup`.
- **Composite Commands:** Intelligent handling of shell-joined commands (e.g., `git status & git log`).
- **Escape Hatch (`pith raw`):** Bypass all parsers when you need the exact, unformatted truth.
- **Middle-Out Truncation:** Automatically keep the start and end of massive outputs, removing the redundant middle.
- **Multi-LLM Integration:** Automated hook setup for Gemini CLI, Claude Code, and Codex.
- **Telemetry & Discovery:** Tracks exact token savings and identifies new optimization targets.
- **Interactive Configuration:** Toggle parsers and adjust truncation limits via `pith config`.
- **Interactive Dashboard:** Launch a local web-based analytics dashboard via `pith dashboard`.

## 📦 Installation
1. Download the latest `pith.exe` from [Releases](https://github.com/Zkrausman/Pith/releases).
2. Run the installer (defaults to global install for all CLIs):
   ```bash
   ./pith.exe install
   ```
   *Note: Pith safely **merges** hooks into your existing `settings.json` files and automatically creates a `.bak` backup before any modification. It will NOT overwrite your existing configurations or other hooks.*
3. Restart your terminal.

---

## 🛠️ Power User Features

### 1. Targeted Installation
If you only want to install hooks for a specific CLI or only for the current project:
- **Local Install (current folder):** `pith install --gemini`
- **Global Install (specific CLI):** `pith install --claude --global`
- **Full Global Sync:** Use `pith install` to ensure all system-wide `settings.json` files are up-to-date.

### 2. The Escape Hatch (`pith raw`)
If a parser is being too aggressive and you need to see the raw, bit-for-bit output of a command, prefix it with `raw`:
```bash
pith raw git diff
```
This bypasses all logic and returns the original system output.

### 2. Middle-Out Truncation
When a command returns thousands of lines (like a massive log file), Pith prevents context overflow by keeping the most important parts: the **beginning** (setup/context) and the **end** (errors/results).

### ⚙️ Interactive Configuration (`pith config`)
Pith features a modern, multi-page TUI built with **Bubble Tea** to manage your optimization engine. Use **`Tab`** to switch between the **Parsers** and **Settings** pages.

#### 1. Parsers Page
Toggle individual command optimizers on or off. If a parser is causing issues with a specific command, you can disable it here without affecting the rest of the system.

#### 2. Settings Page (Output Management)
This page controls **Middle-Out Truncation** and analytics precision.

| Setting | Default | Description |
| :--- | :--- | :--- |
| **`MaxLines`** | `500` | The threshold for truncation. Snipping occurs only if a command's output exceeds this total line count. |
| **`HeadLines`** | `100` | The number of lines to preserve at the **top** of the output. This ensures the LLM understands the command's setup, headers, or initial context. |
| **`TailLines`** | `100` | The number of lines to preserve at the **bottom**. This is critical for capturing final results, exit codes, and error stack traces. |
| **`USD/M Tokens`** | `3.00` | The cost rate (USD per 1 Million Input Tokens) used to calculate financial savings in the dashboard. Adjust this to match your preferred model (e.g., $3.00 for Claude 3.5 Sonnet). |
| **`Heuristic`** | `4.00` | Average characters per token. This is used purely for estimation in analytics and does not affect the actual compression logic. |

---

## 🛠️ Supported Parsers
Pith currently optimizes the following commands:
- **Git:** `status`, `log`, `diff`, `branch`, and **composite commands** (e.g., `git status & git log`)
- **PowerShell:** `Get-Content` (JSON minification & truncation), `Get-ChildItem` (ls/dir minification)
- **Filesystem:** `ls`, `dir`, `find`, `tree`, `du`, `where`
- **Text:** `grep`, `rg`, `cat`, `type` (JSON/XML/CSS/HTML minification)
- **Infra/Dev:** `docker ps`, `npm list`, `pip list`, `npm test`, `go test`, `pytest`, `vitest`, `bd` (Beads), `gh` (GitHub), `env`, `set`

> **Note:** Pith features **cross-platform command matching**, meaning it automatically recognizes commands regardless of their extension (`.exe`, `.cmd`, `.bat`, `.ps1`) or whether they are called via a full path.

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

**Pith Output (~25 tokens):**
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

    Update: Display changelog during pith update

commit 1cd63311f8330490
Author: Zachary Krausman <zkrausman@gmail.com>
Date:   Sun Mar 15 22:30:00 2026 -0400

    Initial commit
```

**Pith Output (~35 tokens):**
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

**Pith Output (~10 tokens):**
```text
main.go README.md
```

---

## 🧪 Testing & Quality
Pith uses [promptfoo](https://promptfoo.dev/) to evaluate the performance and quality of its token-optimization strategies.

- **Status:** **Work In Progress (WIP)** 🛠️
- **Latest Evaluation:** [tests/report.md](tests/report.md)

## 📊 Analytics Dashboard
Pith includes an interactive web dashboard to visualize your token savings and discover new optimization opportunities.

- **Launch:** `pith dashboard`
- **Visualization:** Powered by Chart.js for real-time analytics.
- **Ecosystem:** Pith is **[BuildABoard](https://github.com/Zkrausman/BuildABoard) compatible**. The included `board.json` allows you to seamlessly integrate your Pith analytics into the broader BuildABoard dashboard ecosystem.

---

## ⌨️ CLI Usage
- `pith <command>`: Proxy any command manually.
- `pith config`: Toggle parsers on/off.
- `pith gain`: View total tokens saved.
- `pith discover`: See unparsed commands costing you tokens.
- `pith reset --all`: Clear all telemetry data.
- `pith update`: Auto-update to the latest version.


## Architecture

<!-- mermaid-start -->
```mermaid
graph TD
    e__repos_pith_agents_md["[MODULE] /repos/pith/agents.md [agents.md]"]
    e__repos_pith_changelog_md["[MODULE] /repos/pith/changelog.md [changelog.md]"]
    e__repos_pith_gemini_md["[MODULE] /repos/pith/gemini.md [gemini.md]"]
    e__repos_pith_readme_md["[MODULE] /repos/pith/readme.md [readme.md]"]
    e__repos_pith_main_go["[MODULE] /repos/pith/main.go [main.go]"]
    e__repos_pith_main_go_respondallow["[FUNCTION] respondallow [main.go]"]
    e__repos_pith_main_go_init["[FUNCTION] init [main.go]"]
    e__repos_pith_main_go_main["[FUNCTION] main [main.go]"]
    e__repos_pith_main_go_hookinput["[STRUCT] hookinput [main.go]"]
    e__repos_pith_main_go_hookoutput["[STRUCT] hookoutput [main.go]"]
    e__repos_pith_main_go_toolargs["[STRUCT] toolargs [main.go]"]
    e__repos_pith_main_test_go["[MODULE] /repos/pith/main_test.go [main_test.go]"]
    e__repos_pith_main_test_go_testhookinputschema["[FUNCTION] testhookinputschema [main_test.go]"]
    e__repos_pith_main_test_go_testhookoutputschema["[FUNCTION] testhookoutputschema [main_test.go]"]
    e__repos_pith_pith_parser_generator_skill_md["[MODULE] /repos/pith/pith-parser-generator/skill.md [skill.md]"]
    e__repos_pith_pkg_config_readme_md["[MODULE] /repos/pith/pkg/config/readme.md [readme.md]"]
    e__repos_pith_pkg_config_config_go["[MODULE] /repos/pith/pkg/config/config.go [config.go]"]
    e__repos_pith_pkg_config_config_go_save["[FUNCTION] save [config.go]"]
    e__repos_pith_pkg_config_config_go_interactiveconfig["[FUNCTION] interactiveconfig [config.go]"]
    e__repos_pith_pkg_config_config_go_config["[STRUCT] config [config.go]"]
    e__repos_pith_pkg_config_config_go_getconfigpath["[FUNCTION] getconfigpath [config.go]"]
    e__repos_pith_pkg_config_config_go_loadconfig["[FUNCTION] loadconfig [config.go]"]
    e__repos_pith_pkg_config_config_test_go["[MODULE] /repos/pith/pkg/config/config_test.go [config_test.go]"]
    e__repos_pith_pkg_config_config_test_go_testconfig["[FUNCTION] testconfig [config_test.go]"]
    e__repos_pith_pkg_config_config_test_go_testloadconfiglogic["[FUNCTION] testloadconfiglogic [config_test.go]"]
    e__repos_pith_pkg_config_config_test_go_testconfigsaveload["[FUNCTION] testconfigsaveload [config_test.go]"]
    e__repos_pith_pkg_config_migration_go["[MODULE] /repos/pith/pkg/config/migration.go [migration.go]"]
    e__repos_pith_pkg_config_migration_go_migratestorage["[FUNCTION] migratestorage [migration.go]"]
    e__repos_pith_pkg_config_migration_go_copyfile["[FUNCTION] copyfile [migration.go]"]
    e__repos_pith_pkg_gui_readme_md["[MODULE] /repos/pith/pkg/gui/readme.md [readme.md]"]
    e__repos_pith_pkg_gui_gui_go["[MODULE] /repos/pith/pkg/gui/gui.go [gui.go]"]
    e__repos_pith_pkg_gui_gui_go_route_get_id["[ROUTE] id [gui.go]"]
    e__repos_pith_pkg_gui_gui_go_openbrowser["[FUNCTION] openbrowser [gui.go]"]
    e__repos_pith_pkg_gui_gui_go_startdashboard["[FUNCTION] startdashboard [gui.go]"]
    e__repos_pith_pkg_gui_gui_go_route_get_source["[ROUTE] source [gui.go]"]
    e__repos_pith_pkg_install_readme_md["[MODULE] /repos/pith/pkg/install/readme.md [readme.md]"]
    e__repos_pith_pkg_install_install_go["[MODULE] /repos/pith/pkg/install/install.go [install.go]"]
    e__repos_pith_pkg_install_install_go_installwindows["[FUNCTION] installwindows [install.go]"]
    e__repos_pith_pkg_install_install_go_setuphook["[FUNCTION] setuphook [install.go]"]
    e__repos_pith_pkg_install_install_go_setupcodexhook["[FUNCTION] setupcodexhook [install.go]"]
    e__repos_pith_pkg_install_install_go_hookentry["[STRUCT] hookentry [install.go]"]
    e__repos_pith_pkg_install_install_go_hookgroup["[STRUCT] hookgroup [install.go]"]
    e__repos_pith_pkg_install_install_go_settings["[STRUCT] settings [install.go]"]
    e__repos_pith_pkg_install_install_go_unmarshaljson["[FUNCTION] unmarshaljson [install.go]"]
    e__repos_pith_pkg_install_install_go_marshaljson["[FUNCTION] marshaljson [install.go]"]
    e__repos_pith_pkg_install_install_go_setupgeminihook["[FUNCTION] setupgeminihook [install.go]"]
    e__repos_pith_pkg_install_install_go_setupclaudehook["[FUNCTION] setupclaudehook [install.go]"]
    e__repos_pith_pkg_install_install_go_install["[FUNCTION] install [install.go]"]
    e__repos_pith_pkg_install_install_test_go["[MODULE] /repos/pith/pkg/install/install_test.go [install_test.go]"]
    e__repos_pith_pkg_install_install_test_go_testsetuphooks["[FUNCTION] testsetuphooks [install_test.go]"]
    e__repos_pith_pkg_parser_readme_md["[MODULE] /repos/pith/pkg/parser/readme.md [readme.md]"]
    e__repos_pith_pkg_parser_bd_go["[MODULE] /repos/pith/pkg/parser/bd.go [bd.go]"]
    e__repos_pith_pkg_parser_bd_go_bdparser["[STRUCT] bdparser [bd.go]"]
    e__repos_pith_pkg_parser_bd_go_name["[FUNCTION] name [bd.go]"]
    e__repos_pith_pkg_parser_bd_go_canparse["[FUNCTION] canparse [bd.go]"]
    e__repos_pith_pkg_parser_bd_go_parse["[FUNCTION] parse [bd.go]"]
    e__repos_pith_pkg_parser_chain_go["[MODULE] /repos/pith/pkg/parser/chain.go [chain.go]"]
    e__repos_pith_pkg_parser_chain_go_splitsubcommands["[FUNCTION] splitsubcommands [chain.go]"]
    e__repos_pith_pkg_parser_chain_go_chainparser["[STRUCT] chainparser [chain.go]"]
    e__repos_pith_pkg_parser_chain_go_name["[FUNCTION] name [chain.go]"]
    e__repos_pith_pkg_parser_chain_go_canparse["[FUNCTION] canparse [chain.go]"]
    e__repos_pith_pkg_parser_chain_go_parse["[FUNCTION] parse [chain.go]"]
    e__repos_pith_pkg_parser_fs_go["[MODULE] /repos/pith/pkg/parser/fs.go [fs.go]"]
    e__repos_pith_pkg_parser_fs_go_lsparser["[STRUCT] lsparser [fs.go]"]
    e__repos_pith_pkg_parser_fs_go_name["[FUNCTION] name [fs.go]"]
    e__repos_pith_pkg_parser_fs_go_canparse["[FUNCTION] canparse [fs.go]"]
    e__repos_pith_pkg_parser_fs_go_parse["[FUNCTION] parse [fs.go]"]
    e__repos_pith_pkg_parser_fs_go_findparser["[STRUCT] findparser [fs.go]"]
    e__repos_pith_pkg_parser_fs_go_treeparser["[STRUCT] treeparser [fs.go]"]
    e__repos_pith_pkg_parser_fs_go_duparser["[STRUCT] duparser [fs.go]"]
    e__repos_pith_pkg_parser_fs_test_go["[MODULE] /repos/pith/pkg/parser/fs_test.go [fs_test.go]"]
    e__repos_pith_pkg_parser_fs_test_go_testtreeparser["[FUNCTION] testtreeparser [fs_test.go]"]
    e__repos_pith_pkg_parser_fs_test_go_testduparser["[FUNCTION] testduparser [fs_test.go]"]
    e__repos_pith_pkg_parser_fs_test_go_testlsparser["[FUNCTION] testlsparser [fs_test.go]"]
    e__repos_pith_pkg_parser_fs_test_go_testfindparser["[FUNCTION] testfindparser [fs_test.go]"]
    e__repos_pith_pkg_parser_get_content_test_go["[MODULE] /repos/pith/pkg/parser/get_content_test.go [get_content_test.go]"]
    e__repos_pith_pkg_parser_get_content_test_go_testgetcontentparser["[FUNCTION] testgetcontentparser [get_content_test.go]"]
    e__repos_pith_pkg_parser_git_go["[MODULE] /repos/pith/pkg/parser/git.go [git.go]"]
    e__repos_pith_pkg_parser_git_go_name["[FUNCTION] name [git.go]"]
    e__repos_pith_pkg_parser_git_go_canparse["[FUNCTION] canparse [git.go]"]
    e__repos_pith_pkg_parser_git_go_gitlogparser["[STRUCT] gitlogparser [git.go]"]
    e__repos_pith_pkg_parser_git_go_formatcommit["[FUNCTION] formatcommit [git.go]"]
    e__repos_pith_pkg_parser_git_go_gitdiffparser["[STRUCT] gitdiffparser [git.go]"]
    e__repos_pith_pkg_parser_git_go_gitbranchparser["[STRUCT] gitbranchparser [git.go]"]
    e__repos_pith_pkg_parser_git_go_compositegitparser["[STRUCT] compositegitparser [git.go]"]
    e__repos_pith_pkg_parser_git_go_gitstatusparser["[STRUCT] gitstatusparser [git.go]"]
    e__repos_pith_pkg_parser_git_go_parse["[FUNCTION] parse [git.go]"]
    e__repos_pith_pkg_parser_git_test_go["[MODULE] /repos/pith/pkg/parser/git_test.go [git_test.go]"]
    e__repos_pith_pkg_parser_git_test_go_testgitstatusparser["[FUNCTION] testgitstatusparser [git_test.go]"]
    e__repos_pith_pkg_parser_git_test_go_testgitlogparser["[FUNCTION] testgitlogparser [git_test.go]"]
    e__repos_pith_pkg_parser_git_test_go_testgitdiffparser["[FUNCTION] testgitdiffparser [git_test.go]"]
    e__repos_pith_pkg_parser_git_test_go_testgitbranchparser["[FUNCTION] testgitbranchparser [git_test.go]"]
    e__repos_pith_pkg_parser_git_test_go_testcompositegitparser["[FUNCTION] testcompositegitparser [git_test.go]"]
    e__repos_pith_pkg_parser_github_release_go["[MODULE] /repos/pith/pkg/parser/github_release.go [github_release.go]"]
    e__repos_pith_pkg_parser_github_release_go_githubreleaseparser["[STRUCT] githubreleaseparser [github_release.go]"]
    e__repos_pith_pkg_parser_github_release_go_name["[FUNCTION] name [github_release.go]"]
    e__repos_pith_pkg_parser_github_release_go_canparse["[FUNCTION] canparse [github_release.go]"]
    e__repos_pith_pkg_parser_github_release_go_parse["[FUNCTION] parse [github_release.go]"]
    e__repos_pith_pkg_parser_go_go["[MODULE] /repos/pith/pkg/parser/go.go [go.go]"]
    e__repos_pith_pkg_parser_go_go_name["[FUNCTION] name [go.go]"]
    e__repos_pith_pkg_parser_go_go_canparse["[FUNCTION] canparse [go.go]"]
    e__repos_pith_pkg_parser_go_go_parse["[FUNCTION] parse [go.go]"]
    e__repos_pith_pkg_parser_go_go_goparser["[STRUCT] goparser [go.go]"]
    e__repos_pith_pkg_parser_infra_go["[MODULE] /repos/pith/pkg/parser/infra.go [infra.go]"]
    e__repos_pith_pkg_parser_infra_go_parse["[FUNCTION] parse [infra.go]"]
    e__repos_pith_pkg_parser_infra_go_dockerpsparser["[STRUCT] dockerpsparser [infra.go]"]
    e__repos_pith_pkg_parser_infra_go_dependencyparser["[STRUCT] dependencyparser [infra.go]"]
    e__repos_pith_pkg_parser_infra_go_testparser["[STRUCT] testparser [infra.go]"]
    e__repos_pith_pkg_parser_infra_go_githubparser["[STRUCT] githubparser [infra.go]"]
    e__repos_pith_pkg_parser_infra_go_envparser["[STRUCT] envparser [infra.go]"]
    e__repos_pith_pkg_parser_infra_go_name["[FUNCTION] name [infra.go]"]
    e__repos_pith_pkg_parser_infra_go_canparse["[FUNCTION] canparse [infra.go]"]
    e__repos_pith_pkg_parser_infra_test_go["[MODULE] /repos/pith/pkg/parser/infra_test.go [infra_test.go]"]
    e__repos_pith_pkg_parser_infra_test_go_testenvparser["[FUNCTION] testenvparser [infra_test.go]"]
    e__repos_pith_pkg_parser_infra_test_go_testdockerpsparser["[FUNCTION] testdockerpsparser [infra_test.go]"]
    e__repos_pith_pkg_parser_infra_test_go_testdependencyparser["[FUNCTION] testdependencyparser [infra_test.go]"]
    e__repos_pith_pkg_parser_infra_test_go_testtestparser["[FUNCTION] testtestparser [infra_test.go]"]
    e__repos_pith_pkg_parser_infra_test_go_testgithubparser["[FUNCTION] testgithubparser [infra_test.go]"]
    e__repos_pith_pkg_parser_interface_go["[MODULE] /repos/pith/pkg/parser/interface.go [interface.go]"]
    e__repos_pith_pkg_parser_interface_go_parser["[INTERFACE] parser [interface.go]"]
    e__repos_pith_pkg_parser_interface_go_matchcommand["[FUNCTION] matchcommand [interface.go]"]
    e__repos_pith_pkg_parser_interface_go_getallparsers["[FUNCTION] getallparsers [interface.go]"]
    e__repos_pith_pkg_parser_match_test_go["[MODULE] /repos/pith/pkg/parser/match_test.go [match_test.go]"]
    e__repos_pith_pkg_parser_match_test_go_testmatchcommand["[FUNCTION] testmatchcommand [match_test.go]"]
    e__repos_pith_pkg_parser_new_parsers_test_go["[MODULE] /repos/pith/pkg/parser/new_parsers_test.go [new_parsers_test.go]"]
    e__repos_pith_pkg_parser_new_parsers_test_go_testsourceparser["[FUNCTION] testsourceparser [new_parsers_test.go]"]
    e__repos_pith_pkg_parser_new_parsers_test_go_testgithubreleaseparser["[FUNCTION] testgithubreleaseparser [new_parsers_test.go]"]
    e__repos_pith_pkg_parser_new_parsers_test_go_testchainparser["[FUNCTION] testchainparser [new_parsers_test.go]"]
    e__repos_pith_pkg_parser_new_parsers_test_go_testwebparser["[FUNCTION] testwebparser [new_parsers_test.go]"]
    e__repos_pith_pkg_parser_new_parsers_test_go_testpithparser["[FUNCTION] testpithparser [new_parsers_test.go]"]
    e__repos_pith_pkg_parser_new_parsers_test_go_testgoparser["[FUNCTION] testgoparser [new_parsers_test.go]"]
    e__repos_pith_pkg_parser_npm_go["[MODULE] /repos/pith/pkg/parser/npm.go [npm.go]"]
    e__repos_pith_pkg_parser_npm_go_canparse["[FUNCTION] canparse [npm.go]"]
    e__repos_pith_pkg_parser_npm_go_parse["[FUNCTION] parse [npm.go]"]
    e__repos_pith_pkg_parser_npm_go_npmparser["[STRUCT] npmparser [npm.go]"]
    e__repos_pith_pkg_parser_npm_go_name["[FUNCTION] name [npm.go]"]
    e__repos_pith_pkg_parser_pith_go["[MODULE] /repos/pith/pkg/parser/pith.go [pith.go]"]
    e__repos_pith_pkg_parser_pith_go_name["[FUNCTION] name [pith.go]"]
    e__repos_pith_pkg_parser_pith_go_canparse["[FUNCTION] canparse [pith.go]"]
    e__repos_pith_pkg_parser_pith_go_parse["[FUNCTION] parse [pith.go]"]
    e__repos_pith_pkg_parser_pith_go_pithparser["[STRUCT] pithparser [pith.go]"]
    e__repos_pith_pkg_parser_powershell_go["[MODULE] /repos/pith/pkg/parser/powershell.go [powershell.go]"]
    e__repos_pith_pkg_parser_powershell_go_powershellparser["[STRUCT] powershellparser [powershell.go]"]
    e__repos_pith_pkg_parser_powershell_go_name["[FUNCTION] name [powershell.go]"]
    e__repos_pith_pkg_parser_powershell_go_canparse["[FUNCTION] canparse [powershell.go]"]
    e__repos_pith_pkg_parser_powershell_go_parse["[FUNCTION] parse [powershell.go]"]
    e__repos_pith_pkg_parser_powershell_go_getcontentparser["[STRUCT] getcontentparser [powershell.go]"]
    e__repos_pith_pkg_parser_promptfoo_go["[MODULE] /repos/pith/pkg/parser/promptfoo.go [promptfoo.go]"]
    e__repos_pith_pkg_parser_promptfoo_go_canparse["[FUNCTION] canparse [promptfoo.go]"]
    e__repos_pith_pkg_parser_promptfoo_go_parse["[FUNCTION] parse [promptfoo.go]"]
    e__repos_pith_pkg_parser_promptfoo_go_promptfooparser["[STRUCT] promptfooparser [promptfoo.go]"]
    e__repos_pith_pkg_parser_promptfoo_go_name["[FUNCTION] name [promptfoo.go]"]
    e__repos_pith_pkg_parser_source_go["[MODULE] /repos/pith/pkg/parser/source.go [source.go]"]
    e__repos_pith_pkg_parser_source_go_sourceparser["[STRUCT] sourceparser [source.go]"]
    e__repos_pith_pkg_parser_source_go_name["[FUNCTION] name [source.go]"]
    e__repos_pith_pkg_parser_source_go_canparse["[FUNCTION] canparse [source.go]"]
    e__repos_pith_pkg_parser_source_go_parse["[FUNCTION] parse [source.go]"]
    e__repos_pith_pkg_parser_text_go["[MODULE] /repos/pith/pkg/parser/text.go [text.go]"]
    e__repos_pith_pkg_parser_text_go_minifyparser["[STRUCT] minifyparser [text.go]"]
    e__repos_pith_pkg_parser_text_go_grepparser["[STRUCT] grepparser [text.go]"]
    e__repos_pith_pkg_parser_text_go_name["[FUNCTION] name [text.go]"]
    e__repos_pith_pkg_parser_text_go_canparse["[FUNCTION] canparse [text.go]"]
    e__repos_pith_pkg_parser_text_go_parse["[FUNCTION] parse [text.go]"]
    e__repos_pith_pkg_parser_text_test_go["[MODULE] /repos/pith/pkg/parser/text_test.go [text_test.go]"]
    e__repos_pith_pkg_parser_text_test_go_testgrepparser["[FUNCTION] testgrepparser [text_test.go]"]
    e__repos_pith_pkg_parser_text_test_go_testminifyparser["[FUNCTION] testminifyparser [text_test.go]"]
    e__repos_pith_pkg_parser_thneed_go["[MODULE] /repos/pith/pkg/parser/thneed.go [thneed.go]"]
    e__repos_pith_pkg_parser_thneed_go_parseplain["[FUNCTION] parseplain [thneed.go]"]
    e__repos_pith_pkg_parser_thneed_go_thneedparser["[STRUCT] thneedparser [thneed.go]"]
    e__repos_pith_pkg_parser_thneed_go_name["[FUNCTION] name [thneed.go]"]
    e__repos_pith_pkg_parser_thneed_go_canparse["[FUNCTION] canparse [thneed.go]"]
    e__repos_pith_pkg_parser_thneed_go_parse["[FUNCTION] parse [thneed.go]"]
    e__repos_pith_pkg_parser_thneed_go_parsejson["[FUNCTION] parsejson [thneed.go]"]
    e__repos_pith_pkg_parser_thneed_go_parsejsonobject["[FUNCTION] parsejsonobject [thneed.go]"]
    e__repos_pith_pkg_parser_vitest_go["[MODULE] /repos/pith/pkg/parser/vitest.go [vitest.go]"]
    e__repos_pith_pkg_parser_vitest_go_canparse["[FUNCTION] canparse [vitest.go]"]
    e__repos_pith_pkg_parser_vitest_go_parse["[FUNCTION] parse [vitest.go]"]
    e__repos_pith_pkg_parser_vitest_go_vitestparser["[STRUCT] vitestparser [vitest.go]"]
    e__repos_pith_pkg_parser_vitest_go_name["[FUNCTION] name [vitest.go]"]
    e__repos_pith_pkg_parser_vitest_test_go["[MODULE] /repos/pith/pkg/parser/vitest_test.go [vitest_test.go]"]
    e__repos_pith_pkg_parser_vitest_test_go_testpromptfooparser["[FUNCTION] testpromptfooparser [vitest_test.go]"]
    e__repos_pith_pkg_parser_vitest_test_go_testpowershellparser["[FUNCTION] testpowershellparser [vitest_test.go]"]
    e__repos_pith_pkg_parser_vitest_test_go_testvitestparser["[FUNCTION] testvitestparser [vitest_test.go]"]
    e__repos_pith_pkg_parser_vitest_test_go_testbdparser["[FUNCTION] testbdparser [vitest_test.go]"]
    e__repos_pith_pkg_parser_web_go["[MODULE] /repos/pith/pkg/parser/web.go [web.go]"]
    e__repos_pith_pkg_parser_web_go_webparser["[STRUCT] webparser [web.go]"]
    e__repos_pith_pkg_parser_web_go_name["[FUNCTION] name [web.go]"]
    e__repos_pith_pkg_parser_web_go_canparse["[FUNCTION] canparse [web.go]"]
    e__repos_pith_pkg_parser_web_go_parse["[FUNCTION] parse [web.go]"]
    e__repos_pith_pkg_runner_readme_md["[MODULE] /repos/pith/pkg/runner/readme.md [readme.md]"]
    e__repos_pith_pkg_runner_runner_go["[MODULE] /repos/pith/pkg/runner/runner.go [runner.go]"]
    e__repos_pith_pkg_runner_runner_go_logforsnag["[FUNCTION] logforsnag [runner.go]"]
    e__repos_pith_pkg_runner_runner_go_runwithoptions["[FUNCTION] runwithoptions [runner.go]"]
    e__repos_pith_pkg_runner_runner_go_applymiddleouttruncation["[FUNCTION] applymiddleouttruncation [runner.go]"]
    e__repos_pith_pkg_runner_runner_go_estimatetokens["[FUNCTION] estimatetokens [runner.go]"]
    e__repos_pith_pkg_runner_runner_go_runner["[STRUCT] runner [runner.go]"]
    e__repos_pith_pkg_runner_runner_go_newrunner["[FUNCTION] newrunner [runner.go]"]
    e__repos_pith_pkg_runner_runner_go_detectsource["[FUNCTION] detectsource [runner.go]"]
    e__repos_pith_pkg_runner_runner_go_run["[FUNCTION] run [runner.go]"]
    e__repos_pith_pkg_runner_runner_test_go["[MODULE] /repos/pith/pkg/runner/runner_test.go [runner_test.go]"]
    e__repos_pith_pkg_runner_runner_test_go_testestimatetokens["[FUNCTION] testestimatetokens [runner_test.go]"]
    e__repos_pith_pkg_runner_runner_test_go_testrunner["[FUNCTION] testrunner [runner_test.go]"]
    e__repos_pith_pkg_runner_runner_test_go_testmiddleouttruncation["[FUNCTION] testmiddleouttruncation [runner_test.go]"]
    e__repos_pith_pkg_selfupdate_readme_md["[MODULE] /repos/pith/pkg/selfupdate/readme.md [readme.md]"]
    e__repos_pith_pkg_selfupdate_selfupdate_go["[MODULE] /repos/pith/pkg/selfupdate/selfupdate.go [selfupdate.go]"]
    e__repos_pith_pkg_selfupdate_selfupdate_go_getauthtoken["[FUNCTION] getauthtoken [selfupdate.go]"]
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkandapplyupdate["[FUNCTION] checkandapplyupdate [selfupdate.go]"]
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkforupdatesilent["[FUNCTION] checkforupdatesilent [selfupdate.go]"]
    e__repos_pith_pkg_selfupdate_selfupdate_go_downloadandreplace["[FUNCTION] downloadandreplace [selfupdate.go]"]
    e__repos_pith_pkg_selfupdate_selfupdate_go_release["[STRUCT] release [selfupdate.go]"]
    e__repos_pith_pkg_telemetry_readme_md["[MODULE] /repos/pith/pkg/telemetry/readme.md [readme.md]"]
    e__repos_pith_pkg_telemetry_telemetry_go["[MODULE] /repos/pith/pkg/telemetry/telemetry.go [telemetry.go]"]
    e__repos_pith_pkg_telemetry_telemetry_go_newtelemetry["[FUNCTION] newtelemetry [telemetry.go]"]
    e__repos_pith_pkg_telemetry_telemetry_go_record["[FUNCTION] record [telemetry.go]"]
    e__repos_pith_pkg_telemetry_telemetry_go_getrecentexecutions["[FUNCTION] getrecentexecutions [telemetry.go]"]
    e__repos_pith_pkg_telemetry_telemetry_go_getsources["[FUNCTION] getsources [telemetry.go]"]
    e__repos_pith_pkg_telemetry_telemetry_go_newtelemetrywithpath["[FUNCTION] newtelemetrywithpath [telemetry.go]"]
    e__repos_pith_pkg_telemetry_telemetry_go_init["[FUNCTION] init [telemetry.go]"]
    e__repos_pith_pkg_telemetry_telemetry_go_getstats["[FUNCTION] getstats [telemetry.go]"]
    e__repos_pith_pkg_telemetry_telemetry_go_getstatsbycommand["[FUNCTION] getstatsbycommand [telemetry.go]"]
    e__repos_pith_pkg_telemetry_telemetry_go_resetpassthrough["[FUNCTION] resetpassthrough [telemetry.go]"]
    e__repos_pith_pkg_telemetry_telemetry_go_telemetry["[STRUCT] telemetry [telemetry.go]"]
    e__repos_pith_pkg_telemetry_telemetry_go_close["[FUNCTION] close [telemetry.go]"]
    e__repos_pith_pkg_telemetry_telemetry_go_executionrecord["[STRUCT] executionrecord [telemetry.go]"]
    e__repos_pith_pkg_telemetry_telemetry_go_getstatsbyday["[FUNCTION] getstatsbyday [telemetry.go]"]
    e__repos_pith_pkg_telemetry_telemetry_go_getunparsedcommands["[FUNCTION] getunparsedcommands [telemetry.go]"]
    e__repos_pith_pkg_telemetry_telemetry_go_resetall["[FUNCTION] resetall [telemetry.go]"]
    e__repos_pith_pkg_telemetry_telemetry_go_getexecutiondetails["[FUNCTION] getexecutiondetails [telemetry.go]"]
    e__repos_pith_pkg_telemetry_telemetry_test_go["[MODULE] /repos/pith/pkg/telemetry/telemetry_test.go [telemetry_test.go]"]
    e__repos_pith_pkg_telemetry_telemetry_test_go_testtelemetry["[FUNCTION] testtelemetry [telemetry_test.go]"]
    e__repos_pith_techplanning_roadmap_md["[MODULE] /repos/pith/techplanning/roadmap.md [roadmap.md]"]
    e__repos_pith_techplanning_techplans_v0_11_0_fix_feedback_loop_md["[MODULE] /repos/pith/techplanning/techplans/v0.11.0_fix_feedback_loop.md [v0.11.0_fix_feedback_loop.md]"]
    e__repos_pith_tests_report_md["[MODULE] /repos/pith/tests/report.md [report.md]"]
    beads_pith_bgf["[TASK] pith-bgf [pith]"]
    beads_pith_bhx["[TASK] pith-bhx [pith]"]
    beads_pith_wth["[VERIFIED_SOLUTION] pith-wth [pith]"]
    beads_pith_os4["[TASK] pith-os4 [pith]"]
    beads_pith_4uf["[TASK] pith-4uf [pith]"]
    beads_pith_0d3["[TASK] pith-0d3 [pith]"]
    beads_pith_wg4["[TASK] pith-wg4 [pith]"]
    beads_pith_lcw["[TASK] pith-lcw [pith]"]
    beads_pith_5ju["[TASK] pith-5ju [pith]"]
    beads_pith_bph["[TASK] pith-bph [pith]"]
    beads_pith_41w["[TASK] pith-41w [pith]"]
    pith_pkg_config[["[EXTERNAL] config"]]
    e__repos_pith_main_go -.->|imports|  pith_pkg_config
    pith_pkg_gui[["[EXTERNAL] gui"]]
    e__repos_pith_main_go -.->|imports|  pith_pkg_gui
    pith_pkg_install[["[EXTERNAL] install"]]
    e__repos_pith_main_go -.->|imports|  pith_pkg_install
    pith_pkg_parser[["[EXTERNAL] parser"]]
    e__repos_pith_main_go -.->|imports|  pith_pkg_parser
    pith_pkg_runner[["[EXTERNAL] runner"]]
    e__repos_pith_main_go -.->|imports|  pith_pkg_runner
    pith_pkg_selfupdate[["[EXTERNAL] selfupdate"]]
    e__repos_pith_main_go -.->|imports|  pith_pkg_selfupdate
    pith_pkg_telemetry[["[EXTERNAL] telemetry"]]
    e__repos_pith_main_go -.->|imports|  pith_pkg_telemetry
    encoding_json[["[EXTERNAL] json"]]
    e__repos_pith_main_go -.->|imports|  encoding_json
    fmt[["[EXTERNAL] fmt"]]
    e__repos_pith_main_go -.->|imports|  fmt
    io[["[EXTERNAL] io"]]
    e__repos_pith_main_go -.->|imports|  io
    os[["[EXTERNAL] os"]]
    e__repos_pith_main_go -.->|imports|  os
    strings[["[EXTERNAL] strings"]]
    e__repos_pith_main_go -.->|imports|  strings
    time[["[EXTERNAL] time"]]
    e__repos_pith_main_go -.->|imports|  time
    github_com_spf13_cobra[["[EXTERNAL] cobra"]]
    e__repos_pith_main_go -.->|imports|  github_com_spf13_cobra
    len[["[EXTERNAL] len"]]
    e__repos_pith_main_go -->|calls| len
    help[["[EXTERNAL] help"]]
    e__repos_pith_main_go -->|calls| help
    cmd_help[["[EXTERNAL] cmd.help"]]
    e__repos_pith_main_go -->|calls| cmd_help
    loadconfig[["[EXTERNAL] loadconfig"]]
    e__repos_pith_main_go -->|calls| loadconfig
    config_loadconfig[["[EXTERNAL] config.loadconfig"]]
    e__repos_pith_main_go -->|calls| config_loadconfig
    now[["[EXTERNAL] now"]]
    e__repos_pith_main_go -->|calls| now
    time_now[["[EXTERNAL] time.now"]]
    e__repos_pith_main_go -->|calls| time_now
    unix[["[EXTERNAL] unix"]]
    e__repos_pith_main_go -->|calls| unix
    save[["[EXTERNAL] save"]]
    e__repos_pith_main_go -->|calls| save
    cfg_save[["[EXTERNAL] cfg.save"]]
    e__repos_pith_main_go -->|calls| cfg_save
    checkforupdatesilent[["[EXTERNAL] checkforupdatesilent"]]
    e__repos_pith_main_go -->|calls| checkforupdatesilent
    selfupdate_checkforupdatesilent[["[EXTERNAL] selfupdate.checkforupdatesilent"]]
    e__repos_pith_main_go -->|calls| selfupdate_checkforupdatesilent
    fprintf[["[EXTERNAL] fprintf"]]
    e__repos_pith_main_go -->|calls| fprintf
    fmt_fprintf[["[EXTERNAL] fmt.fprintf"]]
    e__repos_pith_main_go -->|calls| fmt_fprintf
    newtelemetry[["[EXTERNAL] newtelemetry"]]
    e__repos_pith_main_go -->|calls| newtelemetry
    telemetry_newtelemetry[["[EXTERNAL] telemetry.newtelemetry"]]
    e__repos_pith_main_go -->|calls| telemetry_newtelemetry
    close[["[EXTERNAL] close"]]
    e__repos_pith_main_go -->|calls| close
    tel_close[["[EXTERNAL] tel.close"]]
    e__repos_pith_main_go -->|calls| tel_close
    newrunner[["[EXTERNAL] newrunner"]]
    e__repos_pith_main_go -->|calls| newrunner
    runner_newrunner[["[EXTERNAL] runner.newrunner"]]
    e__repos_pith_main_go -->|calls| runner_newrunner
    run[["[EXTERNAL] run"]]
    e__repos_pith_main_go -->|calls| run
    run_run[["[EXTERNAL] run.run"]]
    e__repos_pith_main_go -->|calls| run_run
    getallparsers[["[EXTERNAL] getallparsers"]]
    e__repos_pith_main_go -->|calls| getallparsers
    parser_getallparsers[["[EXTERNAL] parser.getallparsers"]]
    e__repos_pith_main_go -->|calls| parser_getallparsers
    append[["[EXTERNAL] append"]]
    e__repos_pith_main_go -->|calls| append
    name[["[EXTERNAL] name"]]
    e__repos_pith_main_go -->|calls| name
    p_name[["[EXTERNAL] p.name"]]
    e__repos_pith_main_go -->|calls| p_name
    interactiveconfig[["[EXTERNAL] interactiveconfig"]]
    e__repos_pith_main_go -->|calls| interactiveconfig
    cfg_interactiveconfig[["[EXTERNAL] cfg.interactiveconfig"]]
    e__repos_pith_main_go -->|calls| cfg_interactiveconfig
    getstats[["[EXTERNAL] getstats"]]
    e__repos_pith_main_go -->|calls| getstats
    tel_getstats[["[EXTERNAL] tel.getstats"]]
    e__repos_pith_main_go -->|calls| tel_getstats
    printf[["[EXTERNAL] printf"]]
    e__repos_pith_main_go -->|calls| printf
    fmt_printf[["[EXTERNAL] fmt.printf"]]
    e__repos_pith_main_go -->|calls| fmt_printf
    println[["[EXTERNAL] println"]]
    e__repos_pith_main_go -->|calls| println
    fmt_println[["[EXTERNAL] fmt.println"]]
    e__repos_pith_main_go -->|calls| fmt_println
    float64[["[EXTERNAL] float64"]]
    e__repos_pith_main_go -->|calls| float64
    getstatsbycommand[["[EXTERNAL] getstatsbycommand"]]
    e__repos_pith_main_go -->|calls| getstatsbycommand
    tel_getstatsbycommand[["[EXTERNAL] tel.getstatsbycommand"]]
    e__repos_pith_main_go -->|calls| tel_getstatsbycommand
    repeat[["[EXTERNAL] repeat"]]
    e__repos_pith_main_go -->|calls| repeat
    strings_repeat[["[EXTERNAL] strings.repeat"]]
    e__repos_pith_main_go -->|calls| strings_repeat
    getunparsedcommands[["[EXTERNAL] getunparsedcommands"]]
    e__repos_pith_main_go -->|calls| getunparsedcommands
    tel_getunparsedcommands[["[EXTERNAL] tel.getunparsedcommands"]]
    e__repos_pith_main_go -->|calls| tel_getunparsedcommands
    readall[["[EXTERNAL] readall"]]
    e__repos_pith_main_go -->|calls| readall
    io_readall[["[EXTERNAL] io.readall"]]
    e__repos_pith_main_go -->|calls| io_readall
    unmarshal[["[EXTERNAL] unmarshal"]]
    e__repos_pith_main_go -->|calls| unmarshal
    json_unmarshal[["[EXTERNAL] json.unmarshal"]]
    e__repos_pith_main_go -->|calls| json_unmarshal
    respondallow[["[EXTERNAL] respondallow"]]
    e__repos_pith_main_go -->|calls| respondallow
    hasprefix[["[EXTERNAL] hasprefix"]]
    e__repos_pith_main_go -->|calls| hasprefix
    strings_hasprefix[["[EXTERNAL] strings.hasprefix"]]
    e__repos_pith_main_go -->|calls| strings_hasprefix
    trimprefix[["[EXTERNAL] trimprefix"]]
    e__repos_pith_main_go -->|calls| trimprefix
    strings_trimprefix[["[EXTERNAL] strings.trimprefix"]]
    e__repos_pith_main_go -->|calls| strings_trimprefix
    contains[["[EXTERNAL] contains"]]
    e__repos_pith_main_go -->|calls| contains
    strings_contains[["[EXTERNAL] strings.contains"]]
    e__repos_pith_main_go -->|calls| strings_contains
    logforsnag[["[EXTERNAL] logforsnag"]]
    e__repos_pith_main_go -->|calls| logforsnag
    run_logforsnag[["[EXTERNAL] run.logforsnag"]]
    e__repos_pith_main_go -->|calls| run_logforsnag
    fields[["[EXTERNAL] fields"]]
    e__repos_pith_main_go -->|calls| fields
    strings_fields[["[EXTERNAL] strings.fields"]]
    e__repos_pith_main_go -->|calls| strings_fields
    parser_name[["[EXTERNAL] parser.name"]]
    e__repos_pith_main_go -->|calls| parser_name
    canparse[["[EXTERNAL] canparse"]]
    e__repos_pith_main_go -->|calls| canparse
    parser_canparse[["[EXTERNAL] parser.canparse"]]
    e__repos_pith_main_go -->|calls| parser_canparse
    flags[["[EXTERNAL] flags"]]
    e__repos_pith_main_go -->|calls| flags
    cmd_flags[["[EXTERNAL] cmd.flags"]]
    e__repos_pith_main_go -->|calls| cmd_flags
    getstring[["[EXTERNAL] getstring"]]
    e__repos_pith_main_go -->|calls| getstring
    parse[["[EXTERNAL] parse"]]
    e__repos_pith_main_go -->|calls| parse
    p_parse[["[EXTERNAL] p.parse"]]
    e__repos_pith_main_go -->|calls| p_parse
    applymiddleouttruncation[["[EXTERNAL] applymiddleouttruncation"]]
    e__repos_pith_main_go -->|calls| applymiddleouttruncation
    run_applymiddleouttruncation[["[EXTERNAL] run.applymiddleouttruncation"]]
    e__repos_pith_main_go -->|calls| run_applymiddleouttruncation
    record[["[EXTERNAL] record"]]
    e__repos_pith_main_go -->|calls| record
    tel_record[["[EXTERNAL] tel.record"]]
    e__repos_pith_main_go -->|calls| tel_record
    estimatetokens[["[EXTERNAL] estimatetokens"]]
    e__repos_pith_main_go -->|calls| estimatetokens
    runner_estimatetokens[["[EXTERNAL] runner.estimatetokens"]]
    e__repos_pith_main_go -->|calls| runner_estimatetokens
    sprintf[["[EXTERNAL] sprintf"]]
    e__repos_pith_main_go -->|calls| sprintf
    fmt_sprintf[["[EXTERNAL] fmt.sprintf"]]
    e__repos_pith_main_go -->|calls| fmt_sprintf
    newencoder[["[EXTERNAL] newencoder"]]
    e__repos_pith_main_go -->|calls| newencoder
    json_newencoder[["[EXTERNAL] json.newencoder"]]
    e__repos_pith_main_go -->|calls| json_newencoder
    encode[["[EXTERNAL] encode"]]
    e__repos_pith_main_go -->|calls| encode
    e__repos_pith_main_go_respondallow -->|calls| newencoder
    e__repos_pith_main_go_respondallow -->|calls| json_newencoder
    e__repos_pith_main_go_respondallow -->|calls| encode
    install[["[EXTERNAL] install"]]
    e__repos_pith_main_go -->|calls| install
    install_install[["[EXTERNAL] install.install"]]
    e__repos_pith_main_go -->|calls| install_install
    getbool[["[EXTERNAL] getbool"]]
    e__repos_pith_main_go -->|calls| getbool
    changed[["[EXTERNAL] changed"]]
    e__repos_pith_main_go -->|calls| changed
    setupgeminihook[["[EXTERNAL] setupgeminihook"]]
    e__repos_pith_main_go -->|calls| setupgeminihook
    install_setupgeminihook[["[EXTERNAL] install.setupgeminihook"]]
    e__repos_pith_main_go -->|calls| install_setupgeminihook
    setupclaudehook[["[EXTERNAL] setupclaudehook"]]
    e__repos_pith_main_go -->|calls| setupclaudehook
    install_setupclaudehook[["[EXTERNAL] install.setupclaudehook"]]
    e__repos_pith_main_go -->|calls| install_setupclaudehook
    setupcodexhook[["[EXTERNAL] setupcodexhook"]]
    e__repos_pith_main_go -->|calls| setupcodexhook
    install_setupcodexhook[["[EXTERNAL] install.setupcodexhook"]]
    e__repos_pith_main_go -->|calls| install_setupcodexhook
    checkandapplyupdate[["[EXTERNAL] checkandapplyupdate"]]
    e__repos_pith_main_go -->|calls| checkandapplyupdate
    selfupdate_checkandapplyupdate[["[EXTERNAL] selfupdate.checkandapplyupdate"]]
    e__repos_pith_main_go -->|calls| selfupdate_checkandapplyupdate
    resetall[["[EXTERNAL] resetall"]]
    e__repos_pith_main_go -->|calls| resetall
    tel_resetall[["[EXTERNAL] tel.resetall"]]
    e__repos_pith_main_go -->|calls| tel_resetall
    resetpassthrough[["[EXTERNAL] resetpassthrough"]]
    e__repos_pith_main_go -->|calls| resetpassthrough
    tel_resetpassthrough[["[EXTERNAL] tel.resetpassthrough"]]
    e__repos_pith_main_go -->|calls| tel_resetpassthrough
    errorf[["[EXTERNAL] errorf"]]
    e__repos_pith_main_go -->|calls| errorf
    fmt_errorf[["[EXTERNAL] fmt.errorf"]]
    e__repos_pith_main_go -->|calls| fmt_errorf
    runwithoptions[["[EXTERNAL] runwithoptions"]]
    e__repos_pith_main_go -->|calls| runwithoptions
    run_runwithoptions[["[EXTERNAL] run.runwithoptions"]]
    e__repos_pith_main_go -->|calls| run_runwithoptions
    getint[["[EXTERNAL] getint"]]
    e__repos_pith_main_go -->|calls| getint
    startdashboard[["[EXTERNAL] startdashboard"]]
    e__repos_pith_main_go -->|calls| startdashboard
    gui_startdashboard[["[EXTERNAL] gui.startdashboard"]]
    e__repos_pith_main_go -->|calls| gui_startdashboard
    e__repos_pith_main_go_init -->|calls| flags
    resetcmd_flags[["[EXTERNAL] resetcmd.flags"]]
    e__repos_pith_main_go_init -->|calls| resetcmd_flags
    bool[["[EXTERNAL] bool"]]
    e__repos_pith_main_go_init -->|calls| bool
    installcmd_flags[["[EXTERNAL] installcmd.flags"]]
    e__repos_pith_main_go_init -->|calls| installcmd_flags
    boolp[["[EXTERNAL] boolp"]]
    e__repos_pith_main_go_init -->|calls| boolp
    dashboardcmd_flags[["[EXTERNAL] dashboardcmd.flags"]]
    e__repos_pith_main_go_init -->|calls| dashboardcmd_flags
    intp[["[EXTERNAL] intp"]]
    e__repos_pith_main_go_init -->|calls| intp
    hookcmd_flags[["[EXTERNAL] hookcmd.flags"]]
    e__repos_pith_main_go_init -->|calls| hookcmd_flags
    string[["[EXTERNAL] string"]]
    e__repos_pith_main_go_init -->|calls| string
    addcommand[["[EXTERNAL] addcommand"]]
    e__repos_pith_main_go_init -->|calls| addcommand
    rootcmd_addcommand[["[EXTERNAL] rootcmd.addcommand"]]
    e__repos_pith_main_go_init -->|calls| rootcmd_addcommand
    e__repos_pith_main_go_main -->|calls| loadconfig
    e__repos_pith_main_go_main -->|calls| config_loadconfig
    migratestorage[["[EXTERNAL] migratestorage"]]
    e__repos_pith_main_go_main -->|calls| migratestorage
    config_migratestorage[["[EXTERNAL] config.migratestorage"]]
    e__repos_pith_main_go_main -->|calls| config_migratestorage
    execute[["[EXTERNAL] execute"]]
    e__repos_pith_main_go_main -->|calls| execute
    rootcmd_execute[["[EXTERNAL] rootcmd.execute"]]
    e__repos_pith_main_go_main -->|calls| rootcmd_execute
    exit[["[EXTERNAL] exit"]]
    e__repos_pith_main_go_main -->|calls| exit
    os_exit[["[EXTERNAL] os.exit"]]
    e__repos_pith_main_go_main -->|calls| os_exit
    e__repos_pith_main_go ==>|contains| e__repos_pith_main_go_respondallow
    e__repos_pith_main_go ==>|contains| e__repos_pith_main_go_init
    e__repos_pith_main_go ==>|contains| e__repos_pith_main_go_main
    e__repos_pith_main_go ==>|contains| e__repos_pith_main_go_hookinput
    e__repos_pith_main_go ==>|contains| e__repos_pith_main_go_hookoutput
    e__repos_pith_main_go ==>|contains| e__repos_pith_main_go_toolargs
    e__repos_pith_main_test_go -.->|imports|  encoding_json
    testing[["[EXTERNAL] testing"]]
    e__repos_pith_main_test_go -.->|imports|  testing
    e__repos_pith_main_test_go_testhookinputschema -->|calls| unmarshal
    e__repos_pith_main_test_go_testhookinputschema -->|calls| json_unmarshal
    fatal[["[EXTERNAL] fatal"]]
    e__repos_pith_main_test_go_testhookinputschema -->|calls| fatal
    t_fatal[["[EXTERNAL] t.fatal"]]
    e__repos_pith_main_test_go_testhookinputschema -->|calls| t_fatal
    e__repos_pith_main_test_go_testhookinputschema -->|calls| errorf
    t_errorf[["[EXTERNAL] t.errorf"]]
    e__repos_pith_main_test_go_testhookinputschema -->|calls| t_errorf
    error[["[EXTERNAL] error"]]
    e__repos_pith_main_test_go_testhookinputschema -->|calls| error
    t_error[["[EXTERNAL] t.error"]]
    e__repos_pith_main_test_go_testhookinputschema -->|calls| t_error
    marshal[["[EXTERNAL] marshal"]]
    e__repos_pith_main_test_go_testhookoutputschema -->|calls| marshal
    json_marshal[["[EXTERNAL] json.marshal"]]
    e__repos_pith_main_test_go_testhookoutputschema -->|calls| json_marshal
    e__repos_pith_main_test_go_testhookoutputschema -->|calls| fatal
    e__repos_pith_main_test_go_testhookoutputschema -->|calls| t_fatal
    e__repos_pith_main_test_go_testhookoutputschema -->|calls| string
    e__repos_pith_main_test_go_testhookoutputschema -->|calls| errorf
    e__repos_pith_main_test_go_testhookoutputschema -->|calls| t_errorf
    e__repos_pith_main_test_go ==>|contains| e__repos_pith_main_test_go_testhookinputschema
    e__repos_pith_main_test_go ==>|contains| e__repos_pith_main_test_go_testhookoutputschema
    e__repos_pith_pkg_config_config_go -.->|imports|  encoding_json
    e__repos_pith_pkg_config_config_go -.->|imports|  fmt
    e__repos_pith_pkg_config_config_go -.->|imports|  os
    path_filepath[["[EXTERNAL] filepath"]]
    e__repos_pith_pkg_config_config_go -.->|imports|  path_filepath
    sort[["[EXTERNAL] sort"]]
    e__repos_pith_pkg_config_config_go -.->|imports|  sort
    e__repos_pith_pkg_config_config_go -.->|imports|  strings
    github_com_alecaivazis_survey_v2[["[EXTERNAL] v2"]]
    e__repos_pith_pkg_config_config_go -.->|imports|  github_com_alecaivazis_survey_v2
    join[["[EXTERNAL] join"]]
    e__repos_pith_pkg_config_config_go_getconfigpath -->|calls| join
    filepath_join[["[EXTERNAL] filepath.join"]]
    e__repos_pith_pkg_config_config_go_getconfigpath -->|calls| filepath_join
    stat[["[EXTERNAL] stat"]]
    e__repos_pith_pkg_config_config_go_getconfigpath -->|calls| stat
    os_stat[["[EXTERNAL] os.stat"]]
    e__repos_pith_pkg_config_config_go_getconfigpath -->|calls| os_stat
    getenv[["[EXTERNAL] getenv"]]
    e__repos_pith_pkg_config_config_go_getconfigpath -->|calls| getenv
    os_getenv[["[EXTERNAL] os.getenv"]]
    e__repos_pith_pkg_config_config_go_getconfigpath -->|calls| os_getenv
    userhomedir[["[EXTERNAL] userhomedir"]]
    e__repos_pith_pkg_config_config_go_getconfigpath -->|calls| userhomedir
    os_userhomedir[["[EXTERNAL] os.userhomedir"]]
    e__repos_pith_pkg_config_config_go_getconfigpath -->|calls| os_userhomedir
    getconfigpath[["[EXTERNAL] getconfigpath"]]
    e__repos_pith_pkg_config_config_go_loadconfig -->|calls| getconfigpath
    make[["[EXTERNAL] make"]]
    e__repos_pith_pkg_config_config_go_loadconfig -->|calls| make
    dir[["[EXTERNAL] dir"]]
    e__repos_pith_pkg_config_config_go_loadconfig -->|calls| dir
    filepath_dir[["[EXTERNAL] filepath.dir"]]
    e__repos_pith_pkg_config_config_go_loadconfig -->|calls| filepath_dir
    e__repos_pith_pkg_config_config_go_loadconfig -->|calls| join
    e__repos_pith_pkg_config_config_go_loadconfig -->|calls| filepath_join
    e__repos_pith_pkg_config_config_go_loadconfig -->|calls| getenv
    e__repos_pith_pkg_config_config_go_loadconfig -->|calls| os_getenv
    e__repos_pith_pkg_config_config_go_loadconfig -->|calls| stat
    e__repos_pith_pkg_config_config_go_loadconfig -->|calls| os_stat
    e__repos_pith_pkg_config_config_go_loadconfig -->|calls| contains
    e__repos_pith_pkg_config_config_go_loadconfig -->|calls| strings_contains
    readfile[["[EXTERNAL] readfile"]]
    e__repos_pith_pkg_config_config_go_loadconfig -->|calls| readfile
    os_readfile[["[EXTERNAL] os.readfile"]]
    e__repos_pith_pkg_config_config_go_loadconfig -->|calls| os_readfile
    isnotexist[["[EXTERNAL] isnotexist"]]
    e__repos_pith_pkg_config_config_go_loadconfig -->|calls| isnotexist
    os_isnotexist[["[EXTERNAL] os.isnotexist"]]
    e__repos_pith_pkg_config_config_go_loadconfig -->|calls| os_isnotexist
    e__repos_pith_pkg_config_config_go_loadconfig -->|calls| unmarshal
    e__repos_pith_pkg_config_config_go_loadconfig -->|calls| json_unmarshal
    e__repos_pith_pkg_config_config_go_save -->|calls| getconfigpath
    marshalindent[["[EXTERNAL] marshalindent"]]
    e__repos_pith_pkg_config_config_go_save -->|calls| marshalindent
    json_marshalindent[["[EXTERNAL] json.marshalindent"]]
    e__repos_pith_pkg_config_config_go_save -->|calls| json_marshalindent
    writefile[["[EXTERNAL] writefile"]]
    e__repos_pith_pkg_config_config_go_save -->|calls| writefile
    os_writefile[["[EXTERNAL] os.writefile"]]
    e__repos_pith_pkg_config_config_go_save -->|calls| os_writefile
    e__repos_pith_pkg_config_config_go_interactiveconfig -->|calls| strings
    sort_strings[["[EXTERNAL] sort.strings"]]
    e__repos_pith_pkg_config_config_go_interactiveconfig -->|calls| sort_strings
    e__repos_pith_pkg_config_config_go_interactiveconfig -->|calls| append
    askone[["[EXTERNAL] askone"]]
    e__repos_pith_pkg_config_config_go_interactiveconfig -->|calls| askone
    survey_askone[["[EXTERNAL] survey.askone"]]
    e__repos_pith_pkg_config_config_go_interactiveconfig -->|calls| survey_askone
    e__repos_pith_pkg_config_config_go_interactiveconfig -->|calls| make
    e__repos_pith_pkg_config_config_go_interactiveconfig -->|calls| sprintf
    e__repos_pith_pkg_config_config_go_interactiveconfig -->|calls| fmt_sprintf
    ask[["[EXTERNAL] ask"]]
    e__repos_pith_pkg_config_config_go_interactiveconfig -->|calls| ask
    survey_ask[["[EXTERNAL] survey.ask"]]
    e__repos_pith_pkg_config_config_go_interactiveconfig -->|calls| survey_ask
    e__repos_pith_pkg_config_config_go_interactiveconfig -->|calls| save
    c_save[["[EXTERNAL] c.save"]]
    e__repos_pith_pkg_config_config_go_interactiveconfig -->|calls| c_save
    e__repos_pith_pkg_config_config_go ==>|contains| e__repos_pith_pkg_config_config_go_save
    e__repos_pith_pkg_config_config_go ==>|contains| e__repos_pith_pkg_config_config_go_interactiveconfig
    e__repos_pith_pkg_config_config_go ==>|contains| e__repos_pith_pkg_config_config_go_config
    e__repos_pith_pkg_config_config_go ==>|contains| e__repos_pith_pkg_config_config_go_getconfigpath
    e__repos_pith_pkg_config_config_go ==>|contains| e__repos_pith_pkg_config_config_go_loadconfig
    e__repos_pith_pkg_config_config_test_go -.->|imports|  encoding_json
    e__repos_pith_pkg_config_config_test_go -.->|imports|  os
    e__repos_pith_pkg_config_config_test_go -.->|imports|  testing
    mkdirtemp[["[EXTERNAL] mkdirtemp"]]
    e__repos_pith_pkg_config_config_test_go_testconfig -->|calls| mkdirtemp
    os_mkdirtemp[["[EXTERNAL] os.mkdirtemp"]]
    e__repos_pith_pkg_config_config_test_go_testconfig -->|calls| os_mkdirtemp
    e__repos_pith_pkg_config_config_test_go_testconfig -->|calls| fatal
    e__repos_pith_pkg_config_config_test_go_testconfig -->|calls| t_fatal
    removeall[["[EXTERNAL] removeall"]]
    e__repos_pith_pkg_config_config_test_go_testconfig -->|calls| removeall
    os_removeall[["[EXTERNAL] os.removeall"]]
    e__repos_pith_pkg_config_config_test_go_testconfig -->|calls| os_removeall
    e__repos_pith_pkg_config_config_test_go_testloadconfiglogic -->|calls| make
    e__repos_pith_pkg_config_config_test_go_testloadconfiglogic -->|calls| unmarshal
    e__repos_pith_pkg_config_config_test_go_testloadconfiglogic -->|calls| json_unmarshal
    e__repos_pith_pkg_config_config_test_go_testloadconfiglogic -->|calls| fatal
    e__repos_pith_pkg_config_config_test_go_testloadconfiglogic -->|calls| t_fatal
    e__repos_pith_pkg_config_config_test_go_testloadconfiglogic -->|calls| error
    e__repos_pith_pkg_config_config_test_go_testloadconfiglogic -->|calls| t_error
    createtemp[["[EXTERNAL] createtemp"]]
    e__repos_pith_pkg_config_config_test_go_testconfigsaveload -->|calls| createtemp
    os_createtemp[["[EXTERNAL] os.createtemp"]]
    e__repos_pith_pkg_config_config_test_go_testconfigsaveload -->|calls| os_createtemp
    e__repos_pith_pkg_config_config_test_go_testconfigsaveload -->|calls| fatal
    e__repos_pith_pkg_config_config_test_go_testconfigsaveload -->|calls| t_fatal
    remove[["[EXTERNAL] remove"]]
    e__repos_pith_pkg_config_config_test_go_testconfigsaveload -->|calls| remove
    os_remove[["[EXTERNAL] os.remove"]]
    e__repos_pith_pkg_config_config_test_go_testconfigsaveload -->|calls| os_remove
    e__repos_pith_pkg_config_config_test_go_testconfigsaveload -->|calls| name
    tmpfile_name[["[EXTERNAL] tmpfile.name"]]
    e__repos_pith_pkg_config_config_test_go_testconfigsaveload -->|calls| tmpfile_name
    e__repos_pith_pkg_config_config_test_go_testconfigsaveload -->|calls| marshal
    e__repos_pith_pkg_config_config_test_go_testconfigsaveload -->|calls| json_marshal
    e__repos_pith_pkg_config_config_test_go_testconfigsaveload -->|calls| writefile
    e__repos_pith_pkg_config_config_test_go_testconfigsaveload -->|calls| os_writefile
    e__repos_pith_pkg_config_config_test_go_testconfigsaveload -->|calls| readfile
    e__repos_pith_pkg_config_config_test_go_testconfigsaveload -->|calls| os_readfile
    e__repos_pith_pkg_config_config_test_go_testconfigsaveload -->|calls| unmarshal
    e__repos_pith_pkg_config_config_test_go_testconfigsaveload -->|calls| json_unmarshal
    e__repos_pith_pkg_config_config_test_go_testconfigsaveload -->|calls| error
    e__repos_pith_pkg_config_config_test_go_testconfigsaveload -->|calls| t_error
    e__repos_pith_pkg_config_config_test_go ==>|contains| e__repos_pith_pkg_config_config_test_go_testconfig
    e__repos_pith_pkg_config_config_test_go ==>|contains| e__repos_pith_pkg_config_config_test_go_testloadconfiglogic
    e__repos_pith_pkg_config_config_test_go ==>|contains| e__repos_pith_pkg_config_config_test_go_testconfigsaveload
    e__repos_pith_pkg_config_migration_go -.->|imports|  fmt
    e__repos_pith_pkg_config_migration_go -.->|imports|  io
    e__repos_pith_pkg_config_migration_go -.->|imports|  os
    e__repos_pith_pkg_config_migration_go -.->|imports|  path_filepath
    e__repos_pith_pkg_config_migration_go_migratestorage -->|calls| userhomedir
    e__repos_pith_pkg_config_migration_go_migratestorage -->|calls| os_userhomedir
    e__repos_pith_pkg_config_migration_go_migratestorage -->|calls| join
    e__repos_pith_pkg_config_migration_go_migratestorage -->|calls| filepath_join
    e__repos_pith_pkg_config_migration_go_migratestorage -->|calls| stat
    e__repos_pith_pkg_config_migration_go_migratestorage -->|calls| os_stat
    e__repos_pith_pkg_config_migration_go_migratestorage -->|calls| isnotexist
    e__repos_pith_pkg_config_migration_go_migratestorage -->|calls| os_isnotexist
    mkdirall[["[EXTERNAL] mkdirall"]]
    e__repos_pith_pkg_config_migration_go_migratestorage -->|calls| mkdirall
    os_mkdirall[["[EXTERNAL] os.mkdirall"]]
    e__repos_pith_pkg_config_migration_go_migratestorage -->|calls| os_mkdirall
    e__repos_pith_pkg_config_migration_go_migratestorage -->|calls| printf
    e__repos_pith_pkg_config_migration_go_migratestorage -->|calls| fmt_printf
    copyfile[["[EXTERNAL] copyfile"]]
    e__repos_pith_pkg_config_migration_go_migratestorage -->|calls| copyfile
    e__repos_pith_pkg_config_migration_go_migratestorage -->|calls| errorf
    e__repos_pith_pkg_config_migration_go_migratestorage -->|calls| fmt_errorf
    rename[["[EXTERNAL] rename"]]
    e__repos_pith_pkg_config_migration_go_migratestorage -->|calls| rename
    os_rename[["[EXTERNAL] os.rename"]]
    e__repos_pith_pkg_config_migration_go_migratestorage -->|calls| os_rename
    open[["[EXTERNAL] open"]]
    e__repos_pith_pkg_config_migration_go_copyfile -->|calls| open
    os_open[["[EXTERNAL] os.open"]]
    e__repos_pith_pkg_config_migration_go_copyfile -->|calls| os_open
    e__repos_pith_pkg_config_migration_go_copyfile -->|calls| close
    sourcefile_close[["[EXTERNAL] sourcefile.close"]]
    e__repos_pith_pkg_config_migration_go_copyfile -->|calls| sourcefile_close
    create[["[EXTERNAL] create"]]
    e__repos_pith_pkg_config_migration_go_copyfile -->|calls| create
    os_create[["[EXTERNAL] os.create"]]
    e__repos_pith_pkg_config_migration_go_copyfile -->|calls| os_create
    destfile_close[["[EXTERNAL] destfile.close"]]
    e__repos_pith_pkg_config_migration_go_copyfile -->|calls| destfile_close
    copy[["[EXTERNAL] copy"]]
    e__repos_pith_pkg_config_migration_go_copyfile -->|calls| copy
    io_copy[["[EXTERNAL] io.copy"]]
    e__repos_pith_pkg_config_migration_go_copyfile -->|calls| io_copy
    e__repos_pith_pkg_config_migration_go ==>|contains| e__repos_pith_pkg_config_migration_go_migratestorage
    e__repos_pith_pkg_config_migration_go ==>|contains| e__repos_pith_pkg_config_migration_go_copyfile
    e__repos_pith_pkg_gui_gui_go -.->|imports|  pith_pkg_telemetry
    embed[["[EXTERNAL] embed"]]
    e__repos_pith_pkg_gui_gui_go -.->|imports|  embed
    e__repos_pith_pkg_gui_gui_go -.->|imports|  encoding_json
    e__repos_pith_pkg_gui_gui_go -.->|imports|  fmt
    net_http[["[EXTERNAL] http"]]
    e__repos_pith_pkg_gui_gui_go -.->|imports|  net_http
    os_exec[["[EXTERNAL] exec"]]
    e__repos_pith_pkg_gui_gui_go -.->|imports|  os_exec
    runtime[["[EXTERNAL] runtime"]]
    e__repos_pith_pkg_gui_gui_go -.->|imports|  runtime
    handlefunc[["[EXTERNAL] handlefunc"]]
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| handlefunc
    http_handlefunc[["[EXTERNAL] http.handlefunc"]]
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| http_handlefunc
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| readfile
    staticfiles_readfile[["[EXTERNAL] staticfiles.readfile"]]
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| staticfiles_readfile
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| error
    http_error[["[EXTERNAL] http.error"]]
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| http_error
    header[["[EXTERNAL] header"]]
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| header
    w_header[["[EXTERNAL] w.header"]]
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| w_header
    set[["[EXTERNAL] set"]]
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| set
    write[["[EXTERNAL] write"]]
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| write
    w_write[["[EXTERNAL] w.write"]]
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| w_write
    query[["[EXTERNAL] query"]]
    e__repos_pith_pkg_gui_gui_go_route_get_source -->|calls| query
    get[["[EXTERNAL] get"]]
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| get
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| getstats
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| tel_getstats
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| getstatsbycommand
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| tel_getstatsbycommand
    getstatsbyday[["[EXTERNAL] getstatsbyday"]]
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| getstatsbyday
    tel_getstatsbyday[["[EXTERNAL] tel.getstatsbyday"]]
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| tel_getstatsbyday
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| newencoder
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| json_newencoder
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| encode
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| getunparsedcommands
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| tel_getunparsedcommands
    getrecentexecutions[["[EXTERNAL] getrecentexecutions"]]
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| getrecentexecutions
    tel_getrecentexecutions[["[EXTERNAL] tel.getrecentexecutions"]]
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| tel_getrecentexecutions
    e__repos_pith_pkg_gui_gui_go_route_get_id -->|calls| query
    sscanf[["[EXTERNAL] sscanf"]]
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| sscanf
    fmt_sscanf[["[EXTERNAL] fmt.sscanf"]]
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| fmt_sscanf
    getexecutiondetails[["[EXTERNAL] getexecutiondetails"]]
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| getexecutiondetails
    tel_getexecutiondetails[["[EXTERNAL] tel.getexecutiondetails"]]
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| tel_getexecutiondetails
    getsources[["[EXTERNAL] getsources"]]
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| getsources
    tel_getsources[["[EXTERNAL] tel.getsources"]]
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| tel_getsources
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| sprintf
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| fmt_sprintf
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| printf
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| fmt_printf
    openbrowser[["[EXTERNAL] openbrowser"]]
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| openbrowser
    listenandserve[["[EXTERNAL] listenandserve"]]
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| listenandserve
    http_listenandserve[["[EXTERNAL] http.listenandserve"]]
    e__repos_pith_pkg_gui_gui_go_startdashboard -->|calls| http_listenandserve
    command[["[EXTERNAL] command"]]
    e__repos_pith_pkg_gui_gui_go_openbrowser -->|calls| command
    exec_command[["[EXTERNAL] exec.command"]]
    e__repos_pith_pkg_gui_gui_go_openbrowser -->|calls| exec_command
    start[["[EXTERNAL] start"]]
    e__repos_pith_pkg_gui_gui_go_openbrowser -->|calls| start
    e__repos_pith_pkg_gui_gui_go_openbrowser -->|calls| errorf
    e__repos_pith_pkg_gui_gui_go_openbrowser -->|calls| fmt_errorf
    e__repos_pith_pkg_gui_gui_go_openbrowser -->|calls| printf
    e__repos_pith_pkg_gui_gui_go_openbrowser -->|calls| fmt_printf
    e__repos_pith_pkg_gui_gui_go ==>|contains| e__repos_pith_pkg_gui_gui_go_route_get_id
    e__repos_pith_pkg_gui_gui_go ==>|contains| e__repos_pith_pkg_gui_gui_go_openbrowser
    e__repos_pith_pkg_gui_gui_go ==>|contains| e__repos_pith_pkg_gui_gui_go_startdashboard
    e__repos_pith_pkg_gui_gui_go ==>|contains| e__repos_pith_pkg_gui_gui_go_route_get_source
    e__repos_pith_pkg_install_install_go -.->|imports|  encoding_json
    e__repos_pith_pkg_install_install_go -.->|imports|  fmt
    e__repos_pith_pkg_install_install_go -.->|imports|  os
    e__repos_pith_pkg_install_install_go -.->|imports|  os_exec
    e__repos_pith_pkg_install_install_go -.->|imports|  path_filepath
    e__repos_pith_pkg_install_install_go -.->|imports|  runtime
    e__repos_pith_pkg_install_install_go -.->|imports|  strings
    e__repos_pith_pkg_install_install_go_install -->|calls| userhomedir
    e__repos_pith_pkg_install_install_go_install -->|calls| os_userhomedir
    e__repos_pith_pkg_install_install_go_install -->|calls| join
    e__repos_pith_pkg_install_install_go_install -->|calls| filepath_join
    e__repos_pith_pkg_install_install_go_install -->|calls| mkdirall
    e__repos_pith_pkg_install_install_go_install -->|calls| os_mkdirall
    e__repos_pith_pkg_install_install_go_install -->|calls| errorf
    e__repos_pith_pkg_install_install_go_install -->|calls| fmt_errorf
    executable[["[EXTERNAL] executable"]]
    e__repos_pith_pkg_install_install_go_install -->|calls| executable
    os_executable[["[EXTERNAL] os.executable"]]
    e__repos_pith_pkg_install_install_go_install -->|calls| os_executable
    base[["[EXTERNAL] base"]]
    e__repos_pith_pkg_install_install_go_install -->|calls| base
    filepath_base[["[EXTERNAL] filepath.base"]]
    e__repos_pith_pkg_install_install_go_install -->|calls| filepath_base
    abs[["[EXTERNAL] abs"]]
    e__repos_pith_pkg_install_install_go_install -->|calls| abs
    filepath_abs[["[EXTERNAL] filepath.abs"]]
    e__repos_pith_pkg_install_install_go_install -->|calls| filepath_abs
    equalfold[["[EXTERNAL] equalfold"]]
    e__repos_pith_pkg_install_install_go_install -->|calls| equalfold
    strings_equalfold[["[EXTERNAL] strings.equalfold"]]
    e__repos_pith_pkg_install_install_go_install -->|calls| strings_equalfold
    e__repos_pith_pkg_install_install_go_install -->|calls| printf
    e__repos_pith_pkg_install_install_go_install -->|calls| fmt_printf
    e__repos_pith_pkg_install_install_go_install -->|calls| readfile
    e__repos_pith_pkg_install_install_go_install -->|calls| os_readfile
    e__repos_pith_pkg_install_install_go_install -->|calls| writefile
    e__repos_pith_pkg_install_install_go_install -->|calls| os_writefile
    installwindows[["[EXTERNAL] installwindows"]]
    e__repos_pith_pkg_install_install_go_install -->|calls| installwindows
    e__repos_pith_pkg_install_install_go_installwindows -->|calls| command
    e__repos_pith_pkg_install_install_go_installwindows -->|calls| exec_command
    e__repos_pith_pkg_install_install_go_installwindows -->|calls| sprintf
    e__repos_pith_pkg_install_install_go_installwindows -->|calls| fmt_sprintf
    combinedoutput[["[EXTERNAL] combinedoutput"]]
    e__repos_pith_pkg_install_install_go_installwindows -->|calls| combinedoutput
    cmd_combinedoutput[["[EXTERNAL] cmd.combinedoutput"]]
    e__repos_pith_pkg_install_install_go_installwindows -->|calls| cmd_combinedoutput
    e__repos_pith_pkg_install_install_go_installwindows -->|calls| errorf
    e__repos_pith_pkg_install_install_go_installwindows -->|calls| fmt_errorf
    e__repos_pith_pkg_install_install_go_installwindows -->|calls| string
    e__repos_pith_pkg_install_install_go_installwindows -->|calls| println
    e__repos_pith_pkg_install_install_go_installwindows -->|calls| fmt_println
    e__repos_pith_pkg_install_install_go_unmarshaljson -->|calls| unmarshal
    e__repos_pith_pkg_install_install_go_unmarshaljson -->|calls| json_unmarshal
    delete[["[EXTERNAL] delete"]]
    e__repos_pith_pkg_install_install_go_unmarshaljson -->|calls| delete
    e__repos_pith_pkg_install_install_go_marshaljson -->|calls| make
    e__repos_pith_pkg_install_install_go_marshaljson -->|calls| marshal
    e__repos_pith_pkg_install_install_go_marshaljson -->|calls| json_marshal
    e__repos_pith_pkg_install_install_go_setuphook -->|calls| userhomedir
    e__repos_pith_pkg_install_install_go_setuphook -->|calls| os_userhomedir
    e__repos_pith_pkg_install_install_go_setuphook -->|calls| join
    e__repos_pith_pkg_install_install_go_setuphook -->|calls| filepath_join
    e__repos_pith_pkg_install_install_go_setuphook -->|calls| mkdirall
    e__repos_pith_pkg_install_install_go_setuphook -->|calls| os_mkdirall
    e__repos_pith_pkg_install_install_go_setuphook -->|calls| stat
    e__repos_pith_pkg_install_install_go_setuphook -->|calls| os_stat
    e__repos_pith_pkg_install_install_go_setuphook -->|calls| readfile
    e__repos_pith_pkg_install_install_go_setuphook -->|calls| os_readfile
    e__repos_pith_pkg_install_install_go_setuphook -->|calls| writefile
    e__repos_pith_pkg_install_install_go_setuphook -->|calls| os_writefile
    e__repos_pith_pkg_install_install_go_setuphook -->|calls| unmarshal
    e__repos_pith_pkg_install_install_go_setuphook -->|calls| json_unmarshal
    e__repos_pith_pkg_install_install_go_setuphook -->|calls| make
    e__repos_pith_pkg_install_install_go_setuphook -->|calls| sprintf
    e__repos_pith_pkg_install_install_go_setuphook -->|calls| fmt_sprintf
    e__repos_pith_pkg_install_install_go_setuphook -->|calls| append
    e__repos_pith_pkg_install_install_go_setuphook -->|calls| marshalindent
    e__repos_pith_pkg_install_install_go_setuphook -->|calls| json_marshalindent
    setuphook[["[EXTERNAL] setuphook"]]
    e__repos_pith_pkg_install_install_go_setupgeminihook -->|calls| setuphook
    e__repos_pith_pkg_install_install_go_setupgeminihook -->|calls| printf
    e__repos_pith_pkg_install_install_go_setupgeminihook -->|calls| fmt_printf
    e__repos_pith_pkg_install_install_go_setupclaudehook -->|calls| setuphook
    e__repos_pith_pkg_install_install_go_setupclaudehook -->|calls| printf
    e__repos_pith_pkg_install_install_go_setupclaudehook -->|calls| fmt_printf
    e__repos_pith_pkg_install_install_go_setupcodexhook -->|calls| setuphook
    e__repos_pith_pkg_install_install_go_setupcodexhook -->|calls| printf
    e__repos_pith_pkg_install_install_go_setupcodexhook -->|calls| fmt_printf
    e__repos_pith_pkg_install_install_go ==>|contains| e__repos_pith_pkg_install_install_go_installwindows
    e__repos_pith_pkg_install_install_go ==>|contains| e__repos_pith_pkg_install_install_go_setuphook
    e__repos_pith_pkg_install_install_go ==>|contains| e__repos_pith_pkg_install_install_go_setupcodexhook
    e__repos_pith_pkg_install_install_go ==>|contains| e__repos_pith_pkg_install_install_go_hookentry
    e__repos_pith_pkg_install_install_go ==>|contains| e__repos_pith_pkg_install_install_go_hookgroup
    e__repos_pith_pkg_install_install_go ==>|contains| e__repos_pith_pkg_install_install_go_settings
    e__repos_pith_pkg_install_install_go ==>|contains| e__repos_pith_pkg_install_install_go_unmarshaljson
    e__repos_pith_pkg_install_install_go ==>|contains| e__repos_pith_pkg_install_install_go_marshaljson
    e__repos_pith_pkg_install_install_go ==>|contains| e__repos_pith_pkg_install_install_go_setupgeminihook
    e__repos_pith_pkg_install_install_go ==>|contains| e__repos_pith_pkg_install_install_go_setupclaudehook
    e__repos_pith_pkg_install_install_go ==>|contains| e__repos_pith_pkg_install_install_go_install
    e__repos_pith_pkg_install_install_test_go -.->|imports|  os
    e__repos_pith_pkg_install_install_test_go -.->|imports|  testing
    e__repos_pith_pkg_install_install_test_go_testsetuphooks -->|calls| mkdirtemp
    e__repos_pith_pkg_install_install_test_go_testsetuphooks -->|calls| os_mkdirtemp
    e__repos_pith_pkg_install_install_test_go_testsetuphooks -->|calls| fatal
    e__repos_pith_pkg_install_install_test_go_testsetuphooks -->|calls| t_fatal
    e__repos_pith_pkg_install_install_test_go_testsetuphooks -->|calls| removeall
    e__repos_pith_pkg_install_install_test_go_testsetuphooks -->|calls| os_removeall
    getwd[["[EXTERNAL] getwd"]]
    e__repos_pith_pkg_install_install_test_go_testsetuphooks -->|calls| getwd
    os_getwd[["[EXTERNAL] os.getwd"]]
    e__repos_pith_pkg_install_install_test_go_testsetuphooks -->|calls| os_getwd
    chdir[["[EXTERNAL] chdir"]]
    e__repos_pith_pkg_install_install_test_go_testsetuphooks -->|calls| chdir
    os_chdir[["[EXTERNAL] os.chdir"]]
    e__repos_pith_pkg_install_install_test_go_testsetuphooks -->|calls| os_chdir
    e__repos_pith_pkg_install_install_test_go_testsetuphooks -->|calls| setupgeminihook
    e__repos_pith_pkg_install_install_test_go_testsetuphooks -->|calls| errorf
    e__repos_pith_pkg_install_install_test_go_testsetuphooks -->|calls| t_errorf
    e__repos_pith_pkg_install_install_test_go_testsetuphooks -->|calls| stat
    e__repos_pith_pkg_install_install_test_go_testsetuphooks -->|calls| os_stat
    e__repos_pith_pkg_install_install_test_go_testsetuphooks -->|calls| isnotexist
    e__repos_pith_pkg_install_install_test_go_testsetuphooks -->|calls| os_isnotexist
    e__repos_pith_pkg_install_install_test_go_testsetuphooks -->|calls| error
    e__repos_pith_pkg_install_install_test_go_testsetuphooks -->|calls| t_error
    e__repos_pith_pkg_install_install_test_go_testsetuphooks -->|calls| setupclaudehook
    e__repos_pith_pkg_install_install_test_go ==>|contains| e__repos_pith_pkg_install_install_test_go_testsetuphooks
    e__repos_pith_pkg_parser_bd_go -.->|imports|  fmt
    e__repos_pith_pkg_parser_bd_go -.->|imports|  strings
    matchcommand[["[EXTERNAL] matchcommand"]]
    e__repos_pith_pkg_parser_bd_go_canparse -->|calls| matchcommand
    split[["[EXTERNAL] split"]]
    e__repos_pith_pkg_parser_bd_go_parse -->|calls| split
    strings_split[["[EXTERNAL] strings.split"]]
    e__repos_pith_pkg_parser_bd_go_parse -->|calls| strings_split
    trimspace[["[EXTERNAL] trimspace"]]
    e__repos_pith_pkg_parser_bd_go_parse -->|calls| trimspace
    strings_trimspace[["[EXTERNAL] strings.trimspace"]]
    e__repos_pith_pkg_parser_bd_go_parse -->|calls| strings_trimspace
    e__repos_pith_pkg_parser_bd_go_parse -->|calls| hasprefix
    e__repos_pith_pkg_parser_bd_go_parse -->|calls| strings_hasprefix
    e__repos_pith_pkg_parser_bd_go_parse -->|calls| contains
    e__repos_pith_pkg_parser_bd_go_parse -->|calls| strings_contains
    e__repos_pith_pkg_parser_bd_go_parse -->|calls| append
    e__repos_pith_pkg_parser_bd_go_parse -->|calls| len
    e__repos_pith_pkg_parser_bd_go_parse -->|calls| join
    strings_join[["[EXTERNAL] strings.join"]]
    e__repos_pith_pkg_parser_bd_go_parse -->|calls| strings_join
    e__repos_pith_pkg_parser_bd_go_parse -->|calls| sprintf
    e__repos_pith_pkg_parser_bd_go_parse -->|calls| fmt_sprintf
    e__repos_pith_pkg_parser_bd_go ==>|contains| e__repos_pith_pkg_parser_bd_go_bdparser
    e__repos_pith_pkg_parser_bd_go ==>|contains| e__repos_pith_pkg_parser_bd_go_name
    e__repos_pith_pkg_parser_bd_go ==>|contains| e__repos_pith_pkg_parser_bd_go_canparse
    e__repos_pith_pkg_parser_bd_go ==>|contains| e__repos_pith_pkg_parser_bd_go_parse
    e__repos_pith_pkg_parser_chain_go -.->|imports|  strings
    e__repos_pith_pkg_parser_chain_go_canparse -->|calls| join
    e__repos_pith_pkg_parser_chain_go_canparse -->|calls| strings_join
    e__repos_pith_pkg_parser_chain_go_canparse -->|calls| append
    e__repos_pith_pkg_parser_chain_go_canparse -->|calls| contains
    e__repos_pith_pkg_parser_chain_go_canparse -->|calls| strings_contains
    e__repos_pith_pkg_parser_chain_go_splitsubcommands -->|calls| split
    e__repos_pith_pkg_parser_chain_go_splitsubcommands -->|calls| strings_split
    e__repos_pith_pkg_parser_chain_go_splitsubcommands -->|calls| trimspace
    e__repos_pith_pkg_parser_chain_go_splitsubcommands -->|calls| strings_trimspace
    e__repos_pith_pkg_parser_chain_go_splitsubcommands -->|calls| append
    e__repos_pith_pkg_parser_chain_go ==>|contains| e__repos_pith_pkg_parser_chain_go_splitsubcommands
    e__repos_pith_pkg_parser_chain_go ==>|contains| e__repos_pith_pkg_parser_chain_go_chainparser
    e__repos_pith_pkg_parser_chain_go ==>|contains| e__repos_pith_pkg_parser_chain_go_name
    e__repos_pith_pkg_parser_chain_go ==>|contains| e__repos_pith_pkg_parser_chain_go_canparse
    e__repos_pith_pkg_parser_chain_go ==>|contains| e__repos_pith_pkg_parser_chain_go_parse
    e__repos_pith_pkg_parser_fs_go -.->|imports|  fmt
    e__repos_pith_pkg_parser_fs_go -.->|imports|  strings
    e__repos_pith_pkg_parser_fs_go_canparse -->|calls| matchcommand
    e__repos_pith_pkg_parser_fs_go_parse -->|calls| split
    e__repos_pith_pkg_parser_fs_go_parse -->|calls| strings_split
    e__repos_pith_pkg_parser_fs_go_parse -->|calls| trimspace
    e__repos_pith_pkg_parser_fs_go_parse -->|calls| strings_trimspace
    e__repos_pith_pkg_parser_fs_go_parse -->|calls| hasprefix
    e__repos_pith_pkg_parser_fs_go_parse -->|calls| strings_hasprefix
    e__repos_pith_pkg_parser_fs_go_parse -->|calls| fields
    e__repos_pith_pkg_parser_fs_go_parse -->|calls| strings_fields
    e__repos_pith_pkg_parser_fs_go_parse -->|calls| len
    e__repos_pith_pkg_parser_fs_go_parse -->|calls| append
    e__repos_pith_pkg_parser_fs_go_parse -->|calls| sprintf
    e__repos_pith_pkg_parser_fs_go_parse -->|calls| fmt_sprintf
    e__repos_pith_pkg_parser_fs_go_parse -->|calls| join
    e__repos_pith_pkg_parser_fs_go_parse -->|calls| strings_join
    newreplacer[["[EXTERNAL] newreplacer"]]
    e__repos_pith_pkg_parser_fs_go_parse -->|calls| newreplacer
    strings_newreplacer[["[EXTERNAL] strings.newreplacer"]]
    e__repos_pith_pkg_parser_fs_go_parse -->|calls| strings_newreplacer
    replace[["[EXTERNAL] replace"]]
    e__repos_pith_pkg_parser_fs_go_parse -->|calls| replace
    r_replace[["[EXTERNAL] r.replace"]]
    e__repos_pith_pkg_parser_fs_go_parse -->|calls| r_replace
    e__repos_pith_pkg_parser_fs_go ==>|contains| e__repos_pith_pkg_parser_fs_go_lsparser
    e__repos_pith_pkg_parser_fs_go ==>|contains| e__repos_pith_pkg_parser_fs_go_name
    e__repos_pith_pkg_parser_fs_go ==>|contains| e__repos_pith_pkg_parser_fs_go_canparse
    e__repos_pith_pkg_parser_fs_go ==>|contains| e__repos_pith_pkg_parser_fs_go_parse
    e__repos_pith_pkg_parser_fs_go ==>|contains| e__repos_pith_pkg_parser_fs_go_findparser
    e__repos_pith_pkg_parser_fs_go ==>|contains| e__repos_pith_pkg_parser_fs_go_treeparser
    e__repos_pith_pkg_parser_fs_go ==>|contains| e__repos_pith_pkg_parser_fs_go_duparser
    e__repos_pith_pkg_parser_fs_test_go -.->|imports|  strings
    e__repos_pith_pkg_parser_fs_test_go -.->|imports|  testing
    e__repos_pith_pkg_parser_fs_test_go_testlsparser -->|calls| canparse
    p_canparse[["[EXTERNAL] p.canparse"]]
    e__repos_pith_pkg_parser_fs_test_go_testlsparser -->|calls| p_canparse
    e__repos_pith_pkg_parser_fs_test_go_testlsparser -->|calls| error
    e__repos_pith_pkg_parser_fs_test_go_testlsparser -->|calls| t_error
    e__repos_pith_pkg_parser_fs_test_go_testlsparser -->|calls| parse
    e__repos_pith_pkg_parser_fs_test_go_testlsparser -->|calls| p_parse
    e__repos_pith_pkg_parser_fs_test_go_testlsparser -->|calls| contains
    e__repos_pith_pkg_parser_fs_test_go_testlsparser -->|calls| strings_contains
    e__repos_pith_pkg_parser_fs_test_go_testlsparser -->|calls| errorf
    e__repos_pith_pkg_parser_fs_test_go_testlsparser -->|calls| t_errorf
    e__repos_pith_pkg_parser_fs_test_go_testfindparser -->|calls| canparse
    e__repos_pith_pkg_parser_fs_test_go_testfindparser -->|calls| p_canparse
    e__repos_pith_pkg_parser_fs_test_go_testfindparser -->|calls| error
    e__repos_pith_pkg_parser_fs_test_go_testfindparser -->|calls| t_error
    e__repos_pith_pkg_parser_fs_test_go_testfindparser -->|calls| repeat
    e__repos_pith_pkg_parser_fs_test_go_testfindparser -->|calls| strings_repeat
    e__repos_pith_pkg_parser_fs_test_go_testfindparser -->|calls| parse
    e__repos_pith_pkg_parser_fs_test_go_testfindparser -->|calls| p_parse
    e__repos_pith_pkg_parser_fs_test_go_testfindparser -->|calls| contains
    e__repos_pith_pkg_parser_fs_test_go_testfindparser -->|calls| strings_contains
    e__repos_pith_pkg_parser_fs_test_go_testfindparser -->|calls| errorf
    e__repos_pith_pkg_parser_fs_test_go_testfindparser -->|calls| t_errorf
    e__repos_pith_pkg_parser_fs_test_go_testfindparser -->|calls| len
    e__repos_pith_pkg_parser_fs_test_go_testfindparser -->|calls| split
    e__repos_pith_pkg_parser_fs_test_go_testfindparser -->|calls| strings_split
    e__repos_pith_pkg_parser_fs_test_go_testtreeparser -->|calls| canparse
    e__repos_pith_pkg_parser_fs_test_go_testtreeparser -->|calls| p_canparse
    e__repos_pith_pkg_parser_fs_test_go_testtreeparser -->|calls| error
    e__repos_pith_pkg_parser_fs_test_go_testtreeparser -->|calls| t_error
    e__repos_pith_pkg_parser_fs_test_go_testtreeparser -->|calls| parse
    e__repos_pith_pkg_parser_fs_test_go_testtreeparser -->|calls| p_parse
    e__repos_pith_pkg_parser_fs_test_go_testtreeparser -->|calls| contains
    e__repos_pith_pkg_parser_fs_test_go_testtreeparser -->|calls| strings_contains
    e__repos_pith_pkg_parser_fs_test_go_testtreeparser -->|calls| errorf
    e__repos_pith_pkg_parser_fs_test_go_testtreeparser -->|calls| t_errorf
    e__repos_pith_pkg_parser_fs_test_go_testduparser -->|calls| canparse
    e__repos_pith_pkg_parser_fs_test_go_testduparser -->|calls| p_canparse
    e__repos_pith_pkg_parser_fs_test_go_testduparser -->|calls| error
    e__repos_pith_pkg_parser_fs_test_go_testduparser -->|calls| t_error
    e__repos_pith_pkg_parser_fs_test_go_testduparser -->|calls| parse
    e__repos_pith_pkg_parser_fs_test_go_testduparser -->|calls| p_parse
    e__repos_pith_pkg_parser_fs_test_go_testduparser -->|calls| contains
    e__repos_pith_pkg_parser_fs_test_go_testduparser -->|calls| strings_contains
    e__repos_pith_pkg_parser_fs_test_go_testduparser -->|calls| errorf
    e__repos_pith_pkg_parser_fs_test_go_testduparser -->|calls| t_errorf
    e__repos_pith_pkg_parser_fs_test_go ==>|contains| e__repos_pith_pkg_parser_fs_test_go_testtreeparser
    e__repos_pith_pkg_parser_fs_test_go ==>|contains| e__repos_pith_pkg_parser_fs_test_go_testduparser
    e__repos_pith_pkg_parser_fs_test_go ==>|contains| e__repos_pith_pkg_parser_fs_test_go_testlsparser
    e__repos_pith_pkg_parser_fs_test_go ==>|contains| e__repos_pith_pkg_parser_fs_test_go_testfindparser
    e__repos_pith_pkg_parser_get_content_test_go -.->|imports|  strings
    e__repos_pith_pkg_parser_get_content_test_go -.->|imports|  testing
    e__repos_pith_pkg_parser_get_content_test_go_testgetcontentparser -->|calls| canparse
    e__repos_pith_pkg_parser_get_content_test_go_testgetcontentparser -->|calls| p_canparse
    e__repos_pith_pkg_parser_get_content_test_go_testgetcontentparser -->|calls| errorf
    e__repos_pith_pkg_parser_get_content_test_go_testgetcontentparser -->|calls| t_errorf
    e__repos_pith_pkg_parser_get_content_test_go_testgetcontentparser -->|calls| parse
    e__repos_pith_pkg_parser_get_content_test_go_testgetcontentparser -->|calls| p_parse
    e__repos_pith_pkg_parser_get_content_test_go_testgetcontentparser -->|calls| contains
    e__repos_pith_pkg_parser_get_content_test_go_testgetcontentparser -->|calls| strings_contains
    e__repos_pith_pkg_parser_get_content_test_go_testgetcontentparser -->|calls| repeat
    e__repos_pith_pkg_parser_get_content_test_go_testgetcontentparser -->|calls| strings_repeat
    e__repos_pith_pkg_parser_get_content_test_go_testgetcontentparser -->|calls| len
    hassuffix[["[EXTERNAL] hassuffix"]]
    e__repos_pith_pkg_parser_get_content_test_go_testgetcontentparser -->|calls| hassuffix
    strings_hassuffix[["[EXTERNAL] strings.hassuffix"]]
    e__repos_pith_pkg_parser_get_content_test_go_testgetcontentparser -->|calls| strings_hassuffix
    e__repos_pith_pkg_parser_get_content_test_go_testgetcontentparser -->|calls| split
    e__repos_pith_pkg_parser_get_content_test_go_testgetcontentparser -->|calls| strings_split
    e__repos_pith_pkg_parser_get_content_test_go_testgetcontentparser -->|calls| trimspace
    e__repos_pith_pkg_parser_get_content_test_go_testgetcontentparser -->|calls| strings_trimspace
    e__repos_pith_pkg_parser_get_content_test_go ==>|contains| e__repos_pith_pkg_parser_get_content_test_go_testgetcontentparser
    e__repos_pith_pkg_parser_git_go -.->|imports|  fmt
    e__repos_pith_pkg_parser_git_go -.->|imports|  strings
    e__repos_pith_pkg_parser_git_go_canparse -->|calls| len
    e__repos_pith_pkg_parser_git_go_parse -->|calls| split
    e__repos_pith_pkg_parser_git_go_parse -->|calls| strings_split
    e__repos_pith_pkg_parser_git_go_parse -->|calls| trimspace
    e__repos_pith_pkg_parser_git_go_parse -->|calls| strings_trimspace
    e__repos_pith_pkg_parser_git_go_parse -->|calls| hasprefix
    e__repos_pith_pkg_parser_git_go_parse -->|calls| strings_hasprefix
    e__repos_pith_pkg_parser_git_go_parse -->|calls| contains
    e__repos_pith_pkg_parser_git_go_parse -->|calls| strings_contains
    e__repos_pith_pkg_parser_git_go_parse -->|calls| append
    e__repos_pith_pkg_parser_git_go_parse -->|calls| len
    e__repos_pith_pkg_parser_git_go_parse -->|calls| join
    e__repos_pith_pkg_parser_git_go_parse -->|calls| strings_join
    formatcommit[["[EXTERNAL] formatcommit"]]
    e__repos_pith_pkg_parser_git_go_parse -->|calls| formatcommit
    e__repos_pith_pkg_parser_git_go_parse -->|calls| trimprefix
    e__repos_pith_pkg_parser_git_go_parse -->|calls| strings_trimprefix
    index[["[EXTERNAL] index"]]
    e__repos_pith_pkg_parser_git_go_parse -->|calls| index
    strings_index[["[EXTERNAL] strings.index"]]
    e__repos_pith_pkg_parser_git_go_parse -->|calls| strings_index
    e__repos_pith_pkg_parser_git_go_parse -->|calls| fields
    e__repos_pith_pkg_parser_git_go_parse -->|calls| strings_fields
    e__repos_pith_pkg_parser_git_go_parse -->|calls| sprintf
    e__repos_pith_pkg_parser_git_go_parse -->|calls| fmt_sprintf
    e__repos_pith_pkg_parser_git_go_canparse -->|calls| contains
    e__repos_pith_pkg_parser_git_go_canparse -->|calls| strings_contains
    e__repos_pith_pkg_parser_git_go_parse -->|calls| make
    e__repos_pith_pkg_parser_git_go ==>|contains| e__repos_pith_pkg_parser_git_go_name
    e__repos_pith_pkg_parser_git_go ==>|contains| e__repos_pith_pkg_parser_git_go_canparse
    e__repos_pith_pkg_parser_git_go ==>|contains| e__repos_pith_pkg_parser_git_go_gitlogparser
    e__repos_pith_pkg_parser_git_go ==>|contains| e__repos_pith_pkg_parser_git_go_formatcommit
    e__repos_pith_pkg_parser_git_go ==>|contains| e__repos_pith_pkg_parser_git_go_gitdiffparser
    e__repos_pith_pkg_parser_git_go ==>|contains| e__repos_pith_pkg_parser_git_go_gitbranchparser
    e__repos_pith_pkg_parser_git_go ==>|contains| e__repos_pith_pkg_parser_git_go_compositegitparser
    e__repos_pith_pkg_parser_git_go ==>|contains| e__repos_pith_pkg_parser_git_go_gitstatusparser
    e__repos_pith_pkg_parser_git_go ==>|contains| e__repos_pith_pkg_parser_git_go_parse
    e__repos_pith_pkg_parser_git_test_go -.->|imports|  strings
    e__repos_pith_pkg_parser_git_test_go -.->|imports|  testing
    e__repos_pith_pkg_parser_git_test_go_testgitstatusparser -->|calls| parse
    e__repos_pith_pkg_parser_git_test_go_testgitstatusparser -->|calls| p_parse
    e__repos_pith_pkg_parser_git_test_go_testgitstatusparser -->|calls| contains
    e__repos_pith_pkg_parser_git_test_go_testgitstatusparser -->|calls| strings_contains
    e__repos_pith_pkg_parser_git_test_go_testgitstatusparser -->|calls| errorf
    e__repos_pith_pkg_parser_git_test_go_testgitstatusparser -->|calls| t_errorf
    e__repos_pith_pkg_parser_git_test_go_testgitstatusparser -->|calls| error
    e__repos_pith_pkg_parser_git_test_go_testgitstatusparser -->|calls| t_error
    e__repos_pith_pkg_parser_git_test_go_testgitlogparser -->|calls| parse
    e__repos_pith_pkg_parser_git_test_go_testgitlogparser -->|calls| p_parse
    e__repos_pith_pkg_parser_git_test_go_testgitlogparser -->|calls| contains
    e__repos_pith_pkg_parser_git_test_go_testgitlogparser -->|calls| strings_contains
    e__repos_pith_pkg_parser_git_test_go_testgitlogparser -->|calls| errorf
    e__repos_pith_pkg_parser_git_test_go_testgitlogparser -->|calls| t_errorf
    e__repos_pith_pkg_parser_git_test_go_testgitdiffparser -->|calls| parse
    e__repos_pith_pkg_parser_git_test_go_testgitdiffparser -->|calls| p_parse
    e__repos_pith_pkg_parser_git_test_go_testgitdiffparser -->|calls| contains
    e__repos_pith_pkg_parser_git_test_go_testgitdiffparser -->|calls| strings_contains
    e__repos_pith_pkg_parser_git_test_go_testgitdiffparser -->|calls| error
    e__repos_pith_pkg_parser_git_test_go_testgitdiffparser -->|calls| t_error
    e__repos_pith_pkg_parser_git_test_go_testgitbranchparser -->|calls| parse
    e__repos_pith_pkg_parser_git_test_go_testgitbranchparser -->|calls| p_parse
    e__repos_pith_pkg_parser_git_test_go_testgitbranchparser -->|calls| hasprefix
    e__repos_pith_pkg_parser_git_test_go_testgitbranchparser -->|calls| strings_hasprefix
    e__repos_pith_pkg_parser_git_test_go_testgitbranchparser -->|calls| errorf
    e__repos_pith_pkg_parser_git_test_go_testgitbranchparser -->|calls| t_errorf
    e__repos_pith_pkg_parser_git_test_go_testgitbranchparser -->|calls| contains
    e__repos_pith_pkg_parser_git_test_go_testgitbranchparser -->|calls| strings_contains
    e__repos_pith_pkg_parser_git_test_go_testcompositegitparser -->|calls| canparse
    e__repos_pith_pkg_parser_git_test_go_testcompositegitparser -->|calls| p_canparse
    e__repos_pith_pkg_parser_git_test_go_testcompositegitparser -->|calls| error
    e__repos_pith_pkg_parser_git_test_go_testcompositegitparser -->|calls| t_error
    e__repos_pith_pkg_parser_git_test_go_testcompositegitparser -->|calls| parse
    e__repos_pith_pkg_parser_git_test_go_testcompositegitparser -->|calls| p_parse
    e__repos_pith_pkg_parser_git_test_go_testcompositegitparser -->|calls| contains
    e__repos_pith_pkg_parser_git_test_go_testcompositegitparser -->|calls| strings_contains
    e__repos_pith_pkg_parser_git_test_go_testcompositegitparser -->|calls| errorf
    e__repos_pith_pkg_parser_git_test_go_testcompositegitparser -->|calls| t_errorf
    e__repos_pith_pkg_parser_git_test_go ==>|contains| e__repos_pith_pkg_parser_git_test_go_testgitstatusparser
    e__repos_pith_pkg_parser_git_test_go ==>|contains| e__repos_pith_pkg_parser_git_test_go_testgitlogparser
    e__repos_pith_pkg_parser_git_test_go ==>|contains| e__repos_pith_pkg_parser_git_test_go_testgitdiffparser
    e__repos_pith_pkg_parser_git_test_go ==>|contains| e__repos_pith_pkg_parser_git_test_go_testgitbranchparser
    e__repos_pith_pkg_parser_git_test_go ==>|contains| e__repos_pith_pkg_parser_git_test_go_testcompositegitparser
    e__repos_pith_pkg_parser_github_release_go -.->|imports|  fmt
    regexp[["[EXTERNAL] regexp"]]
    e__repos_pith_pkg_parser_github_release_go -.->|imports|  regexp
    e__repos_pith_pkg_parser_github_release_go -.->|imports|  strings
    e__repos_pith_pkg_parser_github_release_go_canparse -->|calls| len
    e__repos_pith_pkg_parser_github_release_go_parse -->|calls| split
    e__repos_pith_pkg_parser_github_release_go_parse -->|calls| strings_split
    mustcompile[["[EXTERNAL] mustcompile"]]
    e__repos_pith_pkg_parser_github_release_go_parse -->|calls| mustcompile
    regexp_mustcompile[["[EXTERNAL] regexp.mustcompile"]]
    e__repos_pith_pkg_parser_github_release_go_parse -->|calls| regexp_mustcompile
    e__repos_pith_pkg_parser_github_release_go_parse -->|calls| trimspace
    e__repos_pith_pkg_parser_github_release_go_parse -->|calls| strings_trimspace
    e__repos_pith_pkg_parser_github_release_go_parse -->|calls| hasprefix
    e__repos_pith_pkg_parser_github_release_go_parse -->|calls| strings_hasprefix
    e__repos_pith_pkg_parser_github_release_go_parse -->|calls| append
    findstringsubmatch[["[EXTERNAL] findstringsubmatch"]]
    e__repos_pith_pkg_parser_github_release_go_parse -->|calls| findstringsubmatch
    reasset_findstringsubmatch[["[EXTERNAL] reasset.findstringsubmatch"]]
    e__repos_pith_pkg_parser_github_release_go_parse -->|calls| reasset_findstringsubmatch
    e__repos_pith_pkg_parser_github_release_go_parse -->|calls| len
    e__repos_pith_pkg_parser_github_release_go_parse -->|calls| sprintf
    e__repos_pith_pkg_parser_github_release_go_parse -->|calls| fmt_sprintf
    replaceallstring[["[EXTERNAL] replaceallstring"]]
    e__repos_pith_pkg_parser_github_release_go_parse -->|calls| replaceallstring
    resha_replaceallstring[["[EXTERNAL] resha.replaceallstring"]]
    e__repos_pith_pkg_parser_github_release_go_parse -->|calls| resha_replaceallstring
    e__repos_pith_pkg_parser_github_release_go_parse -->|calls| join
    e__repos_pith_pkg_parser_github_release_go_parse -->|calls| strings_join
    e__repos_pith_pkg_parser_github_release_go ==>|contains| e__repos_pith_pkg_parser_github_release_go_githubreleaseparser
    e__repos_pith_pkg_parser_github_release_go ==>|contains| e__repos_pith_pkg_parser_github_release_go_name
    e__repos_pith_pkg_parser_github_release_go ==>|contains| e__repos_pith_pkg_parser_github_release_go_canparse
    e__repos_pith_pkg_parser_github_release_go ==>|contains| e__repos_pith_pkg_parser_github_release_go_parse
    e__repos_pith_pkg_parser_go_go -.->|imports|  fmt
    e__repos_pith_pkg_parser_go_go -.->|imports|  strings
    e__repos_pith_pkg_parser_go_go_canparse -->|calls| matchcommand
    e__repos_pith_pkg_parser_go_go_parse -->|calls| split
    e__repos_pith_pkg_parser_go_go_parse -->|calls| strings_split
    e__repos_pith_pkg_parser_go_go_parse -->|calls| trimspace
    e__repos_pith_pkg_parser_go_go_parse -->|calls| strings_trimspace
    e__repos_pith_pkg_parser_go_go_parse -->|calls| contains
    e__repos_pith_pkg_parser_go_go_parse -->|calls| strings_contains
    e__repos_pith_pkg_parser_go_go_parse -->|calls| hasprefix
    e__repos_pith_pkg_parser_go_go_parse -->|calls| strings_hasprefix
    e__repos_pith_pkg_parser_go_go_parse -->|calls| append
    e__repos_pith_pkg_parser_go_go_parse -->|calls| len
    e__repos_pith_pkg_parser_go_go_parse -->|calls| fields
    e__repos_pith_pkg_parser_go_go_parse -->|calls| strings_fields
    e__repos_pith_pkg_parser_go_go_parse -->|calls| join
    e__repos_pith_pkg_parser_go_go_parse -->|calls| strings_join
    e__repos_pith_pkg_parser_go_go_parse -->|calls| sprintf
    e__repos_pith_pkg_parser_go_go_parse -->|calls| fmt_sprintf
    e__repos_pith_pkg_parser_go_go ==>|contains| e__repos_pith_pkg_parser_go_go_name
    e__repos_pith_pkg_parser_go_go ==>|contains| e__repos_pith_pkg_parser_go_go_canparse
    e__repos_pith_pkg_parser_go_go ==>|contains| e__repos_pith_pkg_parser_go_go_parse
    e__repos_pith_pkg_parser_go_go ==>|contains| e__repos_pith_pkg_parser_go_go_goparser
    e__repos_pith_pkg_parser_infra_go -.->|imports|  fmt
    e__repos_pith_pkg_parser_infra_go -.->|imports|  regexp
    e__repos_pith_pkg_parser_infra_go -.->|imports|  strings
    e__repos_pith_pkg_parser_infra_go_canparse -->|calls| matchcommand
    e__repos_pith_pkg_parser_infra_go -->|calls| mustcompile
    e__repos_pith_pkg_parser_infra_go -->|calls| regexp_mustcompile
    e__repos_pith_pkg_parser_infra_go_parse -->|calls| split
    e__repos_pith_pkg_parser_infra_go_parse -->|calls| strings_split
    e__repos_pith_pkg_parser_infra_go_parse -->|calls| contains
    e__repos_pith_pkg_parser_infra_go_parse -->|calls| strings_contains
    splitn[["[EXTERNAL] splitn"]]
    e__repos_pith_pkg_parser_infra_go_parse -->|calls| splitn
    strings_splitn[["[EXTERNAL] strings.splitn"]]
    e__repos_pith_pkg_parser_infra_go_parse -->|calls| strings_splitn
    matchstring[["[EXTERNAL] matchstring"]]
    e__repos_pith_pkg_parser_infra_go_parse -->|calls| matchstring
    skipenvregex_matchstring[["[EXTERNAL] skipenvregex.matchstring"]]
    e__repos_pith_pkg_parser_infra_go_parse -->|calls| skipenvregex_matchstring
    e__repos_pith_pkg_parser_infra_go_parse -->|calls| len
    e__repos_pith_pkg_parser_infra_go_parse -->|calls| append
    e__repos_pith_pkg_parser_infra_go_parse -->|calls| join
    e__repos_pith_pkg_parser_infra_go_parse -->|calls| strings_join
    e__repos_pith_pkg_parser_infra_go_canparse -->|calls| len
    e__repos_pith_pkg_parser_infra_go_parse -->|calls| trimspace
    e__repos_pith_pkg_parser_infra_go_parse -->|calls| strings_trimspace
    e__repos_pith_pkg_parser_infra_go_parse -->|calls| fields
    e__repos_pith_pkg_parser_infra_go_parse -->|calls| strings_fields
    e__repos_pith_pkg_parser_infra_go_parse -->|calls| sprintf
    e__repos_pith_pkg_parser_infra_go_parse -->|calls| fmt_sprintf
    e__repos_pith_pkg_parser_infra_go_parse -->|calls| newreplacer
    e__repos_pith_pkg_parser_infra_go_parse -->|calls| strings_newreplacer
    e__repos_pith_pkg_parser_infra_go_parse -->|calls| replace
    replacer_replace[["[EXTERNAL] replacer.replace"]]
    e__repos_pith_pkg_parser_infra_go_parse -->|calls| replacer_replace
    map[["[EXTERNAL] map"]]
    e__repos_pith_pkg_parser_infra_go_parse -->|calls| map
    strings_map[["[EXTERNAL] strings.map"]]
    e__repos_pith_pkg_parser_infra_go_parse -->|calls| strings_map
    e__repos_pith_pkg_parser_infra_go_parse -->|calls| hasprefix
    e__repos_pith_pkg_parser_infra_go_parse -->|calls| strings_hasprefix
    e__repos_pith_pkg_parser_infra_go_canparse -->|calls| contains
    e__repos_pith_pkg_parser_infra_go_canparse -->|calls| strings_contains
    e__repos_pith_pkg_parser_infra_go_canparse -->|calls| join
    e__repos_pith_pkg_parser_infra_go_canparse -->|calls| strings_join
    e__repos_pith_pkg_parser_infra_go ==>|contains| e__repos_pith_pkg_parser_infra_go_parse
    e__repos_pith_pkg_parser_infra_go ==>|contains| e__repos_pith_pkg_parser_infra_go_dockerpsparser
    e__repos_pith_pkg_parser_infra_go ==>|contains| e__repos_pith_pkg_parser_infra_go_dependencyparser
    e__repos_pith_pkg_parser_infra_go ==>|contains| e__repos_pith_pkg_parser_infra_go_testparser
    e__repos_pith_pkg_parser_infra_go ==>|contains| e__repos_pith_pkg_parser_infra_go_githubparser
    e__repos_pith_pkg_parser_infra_go ==>|contains| e__repos_pith_pkg_parser_infra_go_envparser
    e__repos_pith_pkg_parser_infra_go ==>|contains| e__repos_pith_pkg_parser_infra_go_name
    e__repos_pith_pkg_parser_infra_go ==>|contains| e__repos_pith_pkg_parser_infra_go_canparse
    e__repos_pith_pkg_parser_infra_test_go -.->|imports|  strings
    e__repos_pith_pkg_parser_infra_test_go -.->|imports|  testing
    e__repos_pith_pkg_parser_infra_test_go_testenvparser -->|calls| canparse
    e__repos_pith_pkg_parser_infra_test_go_testenvparser -->|calls| p_canparse
    e__repos_pith_pkg_parser_infra_test_go_testenvparser -->|calls| error
    e__repos_pith_pkg_parser_infra_test_go_testenvparser -->|calls| t_error
    e__repos_pith_pkg_parser_infra_test_go_testenvparser -->|calls| parse
    e__repos_pith_pkg_parser_infra_test_go_testenvparser -->|calls| p_parse
    e__repos_pith_pkg_parser_infra_test_go_testenvparser -->|calls| contains
    e__repos_pith_pkg_parser_infra_test_go_testenvparser -->|calls| strings_contains
    e__repos_pith_pkg_parser_infra_test_go_testdockerpsparser -->|calls| canparse
    e__repos_pith_pkg_parser_infra_test_go_testdockerpsparser -->|calls| p_canparse
    e__repos_pith_pkg_parser_infra_test_go_testdockerpsparser -->|calls| error
    e__repos_pith_pkg_parser_infra_test_go_testdockerpsparser -->|calls| t_error
    e__repos_pith_pkg_parser_infra_test_go_testdockerpsparser -->|calls| parse
    e__repos_pith_pkg_parser_infra_test_go_testdockerpsparser -->|calls| p_parse
    e__repos_pith_pkg_parser_infra_test_go_testdockerpsparser -->|calls| contains
    e__repos_pith_pkg_parser_infra_test_go_testdockerpsparser -->|calls| strings_contains
    e__repos_pith_pkg_parser_infra_test_go_testdockerpsparser -->|calls| errorf
    e__repos_pith_pkg_parser_infra_test_go_testdockerpsparser -->|calls| t_errorf
    e__repos_pith_pkg_parser_infra_test_go_testdependencyparser -->|calls| canparse
    e__repos_pith_pkg_parser_infra_test_go_testdependencyparser -->|calls| p_canparse
    e__repos_pith_pkg_parser_infra_test_go_testdependencyparser -->|calls| error
    e__repos_pith_pkg_parser_infra_test_go_testdependencyparser -->|calls| t_error
    e__repos_pith_pkg_parser_infra_test_go_testdependencyparser -->|calls| parse
    e__repos_pith_pkg_parser_infra_test_go_testdependencyparser -->|calls| p_parse
    e__repos_pith_pkg_parser_infra_test_go_testdependencyparser -->|calls| contains
    e__repos_pith_pkg_parser_infra_test_go_testdependencyparser -->|calls| strings_contains
    e__repos_pith_pkg_parser_infra_test_go_testtestparser -->|calls| canparse
    e__repos_pith_pkg_parser_infra_test_go_testtestparser -->|calls| p_canparse
    e__repos_pith_pkg_parser_infra_test_go_testtestparser -->|calls| error
    e__repos_pith_pkg_parser_infra_test_go_testtestparser -->|calls| t_error
    e__repos_pith_pkg_parser_infra_test_go_testtestparser -->|calls| parse
    e__repos_pith_pkg_parser_infra_test_go_testtestparser -->|calls| p_parse
    e__repos_pith_pkg_parser_infra_test_go_testtestparser -->|calls| contains
    e__repos_pith_pkg_parser_infra_test_go_testtestparser -->|calls| strings_contains
    e__repos_pith_pkg_parser_infra_test_go_testtestparser -->|calls| errorf
    e__repos_pith_pkg_parser_infra_test_go_testtestparser -->|calls| t_errorf
    e__repos_pith_pkg_parser_infra_test_go_testgithubparser -->|calls| canparse
    e__repos_pith_pkg_parser_infra_test_go_testgithubparser -->|calls| p_canparse
    e__repos_pith_pkg_parser_infra_test_go_testgithubparser -->|calls| error
    e__repos_pith_pkg_parser_infra_test_go_testgithubparser -->|calls| t_error
    e__repos_pith_pkg_parser_infra_test_go_testgithubparser -->|calls| parse
    e__repos_pith_pkg_parser_infra_test_go_testgithubparser -->|calls| p_parse
    e__repos_pith_pkg_parser_infra_test_go_testgithubparser -->|calls| contains
    e__repos_pith_pkg_parser_infra_test_go_testgithubparser -->|calls| strings_contains
    e__repos_pith_pkg_parser_infra_test_go_testgithubparser -->|calls| errorf
    e__repos_pith_pkg_parser_infra_test_go_testgithubparser -->|calls| t_errorf
    e__repos_pith_pkg_parser_infra_test_go ==>|contains| e__repos_pith_pkg_parser_infra_test_go_testenvparser
    e__repos_pith_pkg_parser_infra_test_go ==>|contains| e__repos_pith_pkg_parser_infra_test_go_testdockerpsparser
    e__repos_pith_pkg_parser_infra_test_go ==>|contains| e__repos_pith_pkg_parser_infra_test_go_testdependencyparser
    e__repos_pith_pkg_parser_infra_test_go ==>|contains| e__repos_pith_pkg_parser_infra_test_go_testtestparser
    e__repos_pith_pkg_parser_infra_test_go ==>|contains| e__repos_pith_pkg_parser_infra_test_go_testgithubparser
    e__repos_pith_pkg_parser_interface_go -.->|imports|  path_filepath
    e__repos_pith_pkg_parser_interface_go -.->|imports|  strings
    replaceall[["[EXTERNAL] replaceall"]]
    e__repos_pith_pkg_parser_interface_go_matchcommand -->|calls| replaceall
    strings_replaceall[["[EXTERNAL] strings.replaceall"]]
    e__repos_pith_pkg_parser_interface_go_matchcommand -->|calls| strings_replaceall
    tolower[["[EXTERNAL] tolower"]]
    e__repos_pith_pkg_parser_interface_go_matchcommand -->|calls| tolower
    strings_tolower[["[EXTERNAL] strings.tolower"]]
    e__repos_pith_pkg_parser_interface_go_matchcommand -->|calls| strings_tolower
    e__repos_pith_pkg_parser_interface_go_matchcommand -->|calls| base
    e__repos_pith_pkg_parser_interface_go_matchcommand -->|calls| filepath_base
    e__repos_pith_pkg_parser_interface_go_getallparsers -->|calls| append
    e__repos_pith_pkg_parser_interface_go ==>|contains| e__repos_pith_pkg_parser_interface_go_parser
    e__repos_pith_pkg_parser_interface_go ==>|contains| e__repos_pith_pkg_parser_interface_go_matchcommand
    e__repos_pith_pkg_parser_interface_go ==>|contains| e__repos_pith_pkg_parser_interface_go_getallparsers
    e__repos_pith_pkg_parser_match_test_go -.->|imports|  testing
    e__repos_pith_pkg_parser_match_test_go_testmatchcommand -->|calls| run
    t_run[["[EXTERNAL] t.run"]]
    e__repos_pith_pkg_parser_match_test_go_testmatchcommand -->|calls| t_run
    e__repos_pith_pkg_parser_match_test_go_testmatchcommand -->|calls| matchcommand
    e__repos_pith_pkg_parser_match_test_go_testmatchcommand -->|calls| errorf
    e__repos_pith_pkg_parser_match_test_go_testmatchcommand -->|calls| t_errorf
    e__repos_pith_pkg_parser_match_test_go ==>|contains| e__repos_pith_pkg_parser_match_test_go_testmatchcommand
    e__repos_pith_pkg_parser_new_parsers_test_go -.->|imports|  testing
    e__repos_pith_pkg_parser_new_parsers_test_go -.->|imports|  strings
    e__repos_pith_pkg_parser_new_parsers_test_go_testsourceparser -->|calls| parse
    e__repos_pith_pkg_parser_new_parsers_test_go_testsourceparser -->|calls| p_parse
    e__repos_pith_pkg_parser_new_parsers_test_go_testsourceparser -->|calls| contains
    e__repos_pith_pkg_parser_new_parsers_test_go_testsourceparser -->|calls| strings_contains
    e__repos_pith_pkg_parser_new_parsers_test_go_testsourceparser -->|calls| errorf
    e__repos_pith_pkg_parser_new_parsers_test_go_testsourceparser -->|calls| t_errorf
    e__repos_pith_pkg_parser_new_parsers_test_go_testgithubreleaseparser -->|calls| parse
    e__repos_pith_pkg_parser_new_parsers_test_go_testgithubreleaseparser -->|calls| p_parse
    e__repos_pith_pkg_parser_new_parsers_test_go_testgithubreleaseparser -->|calls| contains
    e__repos_pith_pkg_parser_new_parsers_test_go_testgithubreleaseparser -->|calls| strings_contains
    e__repos_pith_pkg_parser_new_parsers_test_go_testgithubreleaseparser -->|calls| errorf
    e__repos_pith_pkg_parser_new_parsers_test_go_testgithubreleaseparser -->|calls| t_errorf
    e__repos_pith_pkg_parser_new_parsers_test_go_testchainparser -->|calls| canparse
    e__repos_pith_pkg_parser_new_parsers_test_go_testchainparser -->|calls| p_canparse
    e__repos_pith_pkg_parser_new_parsers_test_go_testchainparser -->|calls| errorf
    e__repos_pith_pkg_parser_new_parsers_test_go_testchainparser -->|calls| t_errorf
    splitsubcommands[["[EXTERNAL] splitsubcommands"]]
    e__repos_pith_pkg_parser_new_parsers_test_go_testchainparser -->|calls| splitsubcommands
    p_splitsubcommands[["[EXTERNAL] p.splitsubcommands"]]
    e__repos_pith_pkg_parser_new_parsers_test_go_testchainparser -->|calls| p_splitsubcommands
    e__repos_pith_pkg_parser_new_parsers_test_go_testchainparser -->|calls| len
    e__repos_pith_pkg_parser_new_parsers_test_go_testwebparser -->|calls| canparse
    e__repos_pith_pkg_parser_new_parsers_test_go_testwebparser -->|calls| p_canparse
    e__repos_pith_pkg_parser_new_parsers_test_go_testwebparser -->|calls| error
    e__repos_pith_pkg_parser_new_parsers_test_go_testwebparser -->|calls| t_error
    e__repos_pith_pkg_parser_new_parsers_test_go_testwebparser -->|calls| parse
    e__repos_pith_pkg_parser_new_parsers_test_go_testwebparser -->|calls| p_parse
    e__repos_pith_pkg_parser_new_parsers_test_go_testwebparser -->|calls| contains
    e__repos_pith_pkg_parser_new_parsers_test_go_testwebparser -->|calls| strings_contains
    e__repos_pith_pkg_parser_new_parsers_test_go_testwebparser -->|calls| errorf
    e__repos_pith_pkg_parser_new_parsers_test_go_testwebparser -->|calls| t_errorf
    e__repos_pith_pkg_parser_new_parsers_test_go_testpithparser -->|calls| canparse
    e__repos_pith_pkg_parser_new_parsers_test_go_testpithparser -->|calls| p_canparse
    e__repos_pith_pkg_parser_new_parsers_test_go_testpithparser -->|calls| error
    e__repos_pith_pkg_parser_new_parsers_test_go_testpithparser -->|calls| t_error
    e__repos_pith_pkg_parser_new_parsers_test_go_testpithparser -->|calls| parse
    e__repos_pith_pkg_parser_new_parsers_test_go_testpithparser -->|calls| p_parse
    e__repos_pith_pkg_parser_new_parsers_test_go_testpithparser -->|calls| contains
    e__repos_pith_pkg_parser_new_parsers_test_go_testpithparser -->|calls| strings_contains
    e__repos_pith_pkg_parser_new_parsers_test_go_testpithparser -->|calls| errorf
    e__repos_pith_pkg_parser_new_parsers_test_go_testpithparser -->|calls| t_errorf
    e__repos_pith_pkg_parser_new_parsers_test_go_testgoparser -->|calls| canparse
    e__repos_pith_pkg_parser_new_parsers_test_go_testgoparser -->|calls| p_canparse
    e__repos_pith_pkg_parser_new_parsers_test_go_testgoparser -->|calls| error
    e__repos_pith_pkg_parser_new_parsers_test_go_testgoparser -->|calls| t_error
    e__repos_pith_pkg_parser_new_parsers_test_go_testgoparser -->|calls| parse
    e__repos_pith_pkg_parser_new_parsers_test_go_testgoparser -->|calls| p_parse
    e__repos_pith_pkg_parser_new_parsers_test_go_testgoparser -->|calls| contains
    e__repos_pith_pkg_parser_new_parsers_test_go_testgoparser -->|calls| strings_contains
    e__repos_pith_pkg_parser_new_parsers_test_go_testgoparser -->|calls| errorf
    e__repos_pith_pkg_parser_new_parsers_test_go_testgoparser -->|calls| t_errorf
    e__repos_pith_pkg_parser_new_parsers_test_go ==>|contains| e__repos_pith_pkg_parser_new_parsers_test_go_testsourceparser
    e__repos_pith_pkg_parser_new_parsers_test_go ==>|contains| e__repos_pith_pkg_parser_new_parsers_test_go_testgithubreleaseparser
    e__repos_pith_pkg_parser_new_parsers_test_go ==>|contains| e__repos_pith_pkg_parser_new_parsers_test_go_testchainparser
    e__repos_pith_pkg_parser_new_parsers_test_go ==>|contains| e__repos_pith_pkg_parser_new_parsers_test_go_testwebparser
    e__repos_pith_pkg_parser_new_parsers_test_go ==>|contains| e__repos_pith_pkg_parser_new_parsers_test_go_testpithparser
    e__repos_pith_pkg_parser_new_parsers_test_go ==>|contains| e__repos_pith_pkg_parser_new_parsers_test_go_testgoparser
    e__repos_pith_pkg_parser_npm_go -.->|imports|  strings
    e__repos_pith_pkg_parser_npm_go_canparse -->|calls| matchcommand
    e__repos_pith_pkg_parser_npm_go_parse -->|calls| split
    e__repos_pith_pkg_parser_npm_go_parse -->|calls| strings_split
    e__repos_pith_pkg_parser_npm_go_parse -->|calls| trimspace
    e__repos_pith_pkg_parser_npm_go_parse -->|calls| strings_trimspace
    e__repos_pith_pkg_parser_npm_go_parse -->|calls| hasprefix
    e__repos_pith_pkg_parser_npm_go_parse -->|calls| strings_hasprefix
    e__repos_pith_pkg_parser_npm_go_parse -->|calls| contains
    e__repos_pith_pkg_parser_npm_go_parse -->|calls| strings_contains
    e__repos_pith_pkg_parser_npm_go_parse -->|calls| append
    e__repos_pith_pkg_parser_npm_go_parse -->|calls| len
    e__repos_pith_pkg_parser_npm_go_parse -->|calls| join
    e__repos_pith_pkg_parser_npm_go_parse -->|calls| strings_join
    e__repos_pith_pkg_parser_npm_go ==>|contains| e__repos_pith_pkg_parser_npm_go_canparse
    e__repos_pith_pkg_parser_npm_go ==>|contains| e__repos_pith_pkg_parser_npm_go_parse
    e__repos_pith_pkg_parser_npm_go ==>|contains| e__repos_pith_pkg_parser_npm_go_npmparser
    e__repos_pith_pkg_parser_npm_go ==>|contains| e__repos_pith_pkg_parser_npm_go_name
    e__repos_pith_pkg_parser_pith_go -.->|imports|  fmt
    e__repos_pith_pkg_parser_pith_go -.->|imports|  strings
    e__repos_pith_pkg_parser_pith_go_canparse -->|calls| matchcommand
    e__repos_pith_pkg_parser_pith_go_canparse -->|calls| len
    e__repos_pith_pkg_parser_pith_go_parse -->|calls| split
    e__repos_pith_pkg_parser_pith_go_parse -->|calls| strings_split
    e__repos_pith_pkg_parser_pith_go_parse -->|calls| trimspace
    e__repos_pith_pkg_parser_pith_go_parse -->|calls| strings_trimspace
    e__repos_pith_pkg_parser_pith_go_parse -->|calls| hasprefix
    e__repos_pith_pkg_parser_pith_go_parse -->|calls| strings_hasprefix
    e__repos_pith_pkg_parser_pith_go_parse -->|calls| contains
    e__repos_pith_pkg_parser_pith_go_parse -->|calls| strings_contains
    e__repos_pith_pkg_parser_pith_go_parse -->|calls| append
    e__repos_pith_pkg_parser_pith_go_parse -->|calls| len
    e__repos_pith_pkg_parser_pith_go_parse -->|calls| join
    e__repos_pith_pkg_parser_pith_go_parse -->|calls| strings_join
    e__repos_pith_pkg_parser_pith_go_parse -->|calls| sprintf
    e__repos_pith_pkg_parser_pith_go_parse -->|calls| fmt_sprintf
    e__repos_pith_pkg_parser_pith_go ==>|contains| e__repos_pith_pkg_parser_pith_go_name
    e__repos_pith_pkg_parser_pith_go ==>|contains| e__repos_pith_pkg_parser_pith_go_canparse
    e__repos_pith_pkg_parser_pith_go ==>|contains| e__repos_pith_pkg_parser_pith_go_parse
    e__repos_pith_pkg_parser_pith_go ==>|contains| e__repos_pith_pkg_parser_pith_go_pithparser
    e__repos_pith_pkg_parser_powershell_go -.->|imports|  fmt
    e__repos_pith_pkg_parser_powershell_go -.->|imports|  strings
    e__repos_pith_pkg_parser_powershell_go_canparse -->|calls| matchcommand
    e__repos_pith_pkg_parser_powershell_go_parse -->|calls| split
    e__repos_pith_pkg_parser_powershell_go_parse -->|calls| strings_split
    e__repos_pith_pkg_parser_powershell_go_parse -->|calls| trimspace
    e__repos_pith_pkg_parser_powershell_go_parse -->|calls| strings_trimspace
    e__repos_pith_pkg_parser_powershell_go_parse -->|calls| hasprefix
    e__repos_pith_pkg_parser_powershell_go_parse -->|calls| strings_hasprefix
    e__repos_pith_pkg_parser_powershell_go_parse -->|calls| contains
    e__repos_pith_pkg_parser_powershell_go_parse -->|calls| strings_contains
    e__repos_pith_pkg_parser_powershell_go_parse -->|calls| fields
    e__repos_pith_pkg_parser_powershell_go_parse -->|calls| strings_fields
    e__repos_pith_pkg_parser_powershell_go_parse -->|calls| len
    e__repos_pith_pkg_parser_powershell_go_parse -->|calls| append
    e__repos_pith_pkg_parser_powershell_go_parse -->|calls| sprintf
    e__repos_pith_pkg_parser_powershell_go_parse -->|calls| fmt_sprintf
    e__repos_pith_pkg_parser_powershell_go_parse -->|calls| join
    e__repos_pith_pkg_parser_powershell_go_parse -->|calls| strings_join
    e__repos_pith_pkg_parser_powershell_go_canparse -->|calls| join
    e__repos_pith_pkg_parser_powershell_go_canparse -->|calls| strings_join
    e__repos_pith_pkg_parser_powershell_go_canparse -->|calls| append
    e__repos_pith_pkg_parser_powershell_go_canparse -->|calls| contains
    e__repos_pith_pkg_parser_powershell_go_canparse -->|calls| strings_contains
    e__repos_pith_pkg_parser_powershell_go_parse -->|calls| hassuffix
    e__repos_pith_pkg_parser_powershell_go_parse -->|calls| strings_hassuffix
    writebyte[["[EXTERNAL] writebyte"]]
    e__repos_pith_pkg_parser_powershell_go_parse -->|calls| writebyte
    minified_writebyte[["[EXTERNAL] minified.writebyte"]]
    e__repos_pith_pkg_parser_powershell_go_parse -->|calls| minified_writebyte
    e__repos_pith_pkg_parser_powershell_go_parse -->|calls| string
    minified_string[["[EXTERNAL] minified.string"]]
    e__repos_pith_pkg_parser_powershell_go_parse -->|calls| minified_string
    e__repos_pith_pkg_parser_powershell_go ==>|contains| e__repos_pith_pkg_parser_powershell_go_powershellparser
    e__repos_pith_pkg_parser_powershell_go ==>|contains| e__repos_pith_pkg_parser_powershell_go_name
    e__repos_pith_pkg_parser_powershell_go ==>|contains| e__repos_pith_pkg_parser_powershell_go_canparse
    e__repos_pith_pkg_parser_powershell_go ==>|contains| e__repos_pith_pkg_parser_powershell_go_parse
    e__repos_pith_pkg_parser_powershell_go ==>|contains| e__repos_pith_pkg_parser_powershell_go_getcontentparser
    e__repos_pith_pkg_parser_promptfoo_go -.->|imports|  fmt
    e__repos_pith_pkg_parser_promptfoo_go -.->|imports|  strings
    e__repos_pith_pkg_parser_promptfoo_go_canparse -->|calls| matchcommand
    e__repos_pith_pkg_parser_promptfoo_go_canparse -->|calls| len
    e__repos_pith_pkg_parser_promptfoo_go_parse -->|calls| split
    e__repos_pith_pkg_parser_promptfoo_go_parse -->|calls| strings_split
    e__repos_pith_pkg_parser_promptfoo_go_parse -->|calls| trimspace
    e__repos_pith_pkg_parser_promptfoo_go_parse -->|calls| strings_trimspace
    e__repos_pith_pkg_parser_promptfoo_go_parse -->|calls| hasprefix
    e__repos_pith_pkg_parser_promptfoo_go_parse -->|calls| strings_hasprefix
    e__repos_pith_pkg_parser_promptfoo_go_parse -->|calls| contains
    e__repos_pith_pkg_parser_promptfoo_go_parse -->|calls| strings_contains
    e__repos_pith_pkg_parser_promptfoo_go_parse -->|calls| append
    e__repos_pith_pkg_parser_promptfoo_go_parse -->|calls| len
    e__repos_pith_pkg_parser_promptfoo_go_parse -->|calls| join
    e__repos_pith_pkg_parser_promptfoo_go_parse -->|calls| strings_join
    e__repos_pith_pkg_parser_promptfoo_go_parse -->|calls| sprintf
    e__repos_pith_pkg_parser_promptfoo_go_parse -->|calls| fmt_sprintf
    e__repos_pith_pkg_parser_promptfoo_go ==>|contains| e__repos_pith_pkg_parser_promptfoo_go_canparse
    e__repos_pith_pkg_parser_promptfoo_go ==>|contains| e__repos_pith_pkg_parser_promptfoo_go_parse
    e__repos_pith_pkg_parser_promptfoo_go ==>|contains| e__repos_pith_pkg_parser_promptfoo_go_promptfooparser
    e__repos_pith_pkg_parser_promptfoo_go ==>|contains| e__repos_pith_pkg_parser_promptfoo_go_name
    e__repos_pith_pkg_parser_source_go -.->|imports|  regexp
    e__repos_pith_pkg_parser_source_go -.->|imports|  strings
    e__repos_pith_pkg_parser_source_go_parse -->|calls| split
    e__repos_pith_pkg_parser_source_go_parse -->|calls| strings_split
    e__repos_pith_pkg_parser_source_go_parse -->|calls| mustcompile
    e__repos_pith_pkg_parser_source_go_parse -->|calls| regexp_mustcompile
    e__repos_pith_pkg_parser_source_go_parse -->|calls| trimspace
    e__repos_pith_pkg_parser_source_go_parse -->|calls| strings_trimspace
    e__repos_pith_pkg_parser_source_go_parse -->|calls| contains
    e__repos_pith_pkg_parser_source_go_parse -->|calls| strings_contains
    e__repos_pith_pkg_parser_source_go_parse -->|calls| replaceallstring
    reinline_replaceallstring[["[EXTERNAL] reinline.replaceallstring"]]
    e__repos_pith_pkg_parser_source_go_parse -->|calls| reinline_replaceallstring
    respaces_replaceallstring[["[EXTERNAL] respaces.replaceallstring"]]
    e__repos_pith_pkg_parser_source_go_parse -->|calls| respaces_replaceallstring
    e__repos_pith_pkg_parser_source_go_parse -->|calls| append
    e__repos_pith_pkg_parser_source_go_parse -->|calls| join
    e__repos_pith_pkg_parser_source_go_parse -->|calls| strings_join
    e__repos_pith_pkg_parser_source_go ==>|contains| e__repos_pith_pkg_parser_source_go_sourceparser
    e__repos_pith_pkg_parser_source_go ==>|contains| e__repos_pith_pkg_parser_source_go_name
    e__repos_pith_pkg_parser_source_go ==>|contains| e__repos_pith_pkg_parser_source_go_canparse
    e__repos_pith_pkg_parser_source_go ==>|contains| e__repos_pith_pkg_parser_source_go_parse
    e__repos_pith_pkg_parser_text_go -.->|imports|  regexp
    e__repos_pith_pkg_parser_text_go -.->|imports|  strings
    e__repos_pith_pkg_parser_text_go_parse -->|calls| split
    e__repos_pith_pkg_parser_text_go_parse -->|calls| strings_split
    e__repos_pith_pkg_parser_text_go_parse -->|calls| trimspace
    e__repos_pith_pkg_parser_text_go_parse -->|calls| strings_trimspace
    e__repos_pith_pkg_parser_text_go_parse -->|calls| splitn
    e__repos_pith_pkg_parser_text_go_parse -->|calls| strings_splitn
    e__repos_pith_pkg_parser_text_go_parse -->|calls| len
    e__repos_pith_pkg_parser_text_go_parse -->|calls| append
    e__repos_pith_pkg_parser_text_go_parse -->|calls| join
    e__repos_pith_pkg_parser_text_go_parse -->|calls| strings_join
    e__repos_pith_pkg_parser_text_go_canparse -->|calls| len
    e__repos_pith_pkg_parser_text_go_canparse -->|calls| hassuffix
    e__repos_pith_pkg_parser_text_go_canparse -->|calls| strings_hassuffix
    e__repos_pith_pkg_parser_text_go -->|calls| mustcompile
    e__repos_pith_pkg_parser_text_go -->|calls| regexp_mustcompile
    e__repos_pith_pkg_parser_text_go_parse -->|calls| hasprefix
    e__repos_pith_pkg_parser_text_go_parse -->|calls| strings_hasprefix
    e__repos_pith_pkg_parser_text_go_parse -->|calls| replaceallstring
    whitespaceregex_replaceallstring[["[EXTERNAL] whitespaceregex.replaceallstring"]]
    e__repos_pith_pkg_parser_text_go_parse -->|calls| whitespaceregex_replaceallstring
    e__repos_pith_pkg_parser_text_go ==>|contains| e__repos_pith_pkg_parser_text_go_minifyparser
    e__repos_pith_pkg_parser_text_go ==>|contains| e__repos_pith_pkg_parser_text_go_grepparser
    e__repos_pith_pkg_parser_text_go ==>|contains| e__repos_pith_pkg_parser_text_go_name
    e__repos_pith_pkg_parser_text_go ==>|contains| e__repos_pith_pkg_parser_text_go_canparse
    e__repos_pith_pkg_parser_text_go ==>|contains| e__repos_pith_pkg_parser_text_go_parse
    e__repos_pith_pkg_parser_text_test_go -.->|imports|  strings
    e__repos_pith_pkg_parser_text_test_go -.->|imports|  testing
    e__repos_pith_pkg_parser_text_test_go_testgrepparser -->|calls| parse
    e__repos_pith_pkg_parser_text_test_go_testgrepparser -->|calls| p_parse
    e__repos_pith_pkg_parser_text_test_go_testgrepparser -->|calls| contains
    e__repos_pith_pkg_parser_text_test_go_testgrepparser -->|calls| strings_contains
    e__repos_pith_pkg_parser_text_test_go_testgrepparser -->|calls| errorf
    e__repos_pith_pkg_parser_text_test_go_testgrepparser -->|calls| t_errorf
    count[["[EXTERNAL] count"]]
    e__repos_pith_pkg_parser_text_test_go_testgrepparser -->|calls| count
    strings_count[["[EXTERNAL] strings.count"]]
    e__repos_pith_pkg_parser_text_test_go_testgrepparser -->|calls| strings_count
    e__repos_pith_pkg_parser_text_test_go_testgrepparser -->|calls| error
    e__repos_pith_pkg_parser_text_test_go_testgrepparser -->|calls| t_error
    e__repos_pith_pkg_parser_text_test_go_testminifyparser -->|calls| canparse
    e__repos_pith_pkg_parser_text_test_go_testminifyparser -->|calls| p_canparse
    e__repos_pith_pkg_parser_text_test_go_testminifyparser -->|calls| error
    e__repos_pith_pkg_parser_text_test_go_testminifyparser -->|calls| t_error
    e__repos_pith_pkg_parser_text_test_go_testminifyparser -->|calls| parse
    e__repos_pith_pkg_parser_text_test_go_testminifyparser -->|calls| p_parse
    e__repos_pith_pkg_parser_text_test_go_testminifyparser -->|calls| contains
    e__repos_pith_pkg_parser_text_test_go_testminifyparser -->|calls| strings_contains
    e__repos_pith_pkg_parser_text_test_go ==>|contains| e__repos_pith_pkg_parser_text_test_go_testgrepparser
    e__repos_pith_pkg_parser_text_test_go ==>|contains| e__repos_pith_pkg_parser_text_test_go_testminifyparser
    e__repos_pith_pkg_parser_thneed_go -.->|imports|  encoding_json
    e__repos_pith_pkg_parser_thneed_go -.->|imports|  fmt
    e__repos_pith_pkg_parser_thneed_go -.->|imports|  strings
    e__repos_pith_pkg_parser_thneed_go_canparse -->|calls| matchcommand
    e__repos_pith_pkg_parser_thneed_go_parse -->|calls| unmarshal
    e__repos_pith_pkg_parser_thneed_go_parse -->|calls| json_unmarshal
    parsejson[["[EXTERNAL] parsejson"]]
    e__repos_pith_pkg_parser_thneed_go_parse -->|calls| parsejson
    t_parsejson[["[EXTERNAL] t.parsejson"]]
    e__repos_pith_pkg_parser_thneed_go_parse -->|calls| t_parsejson
    e__repos_pith_pkg_parser_thneed_go_parse -->|calls| sprintf
    e__repos_pith_pkg_parser_thneed_go_parse -->|calls| fmt_sprintf
    parsejsonobject[["[EXTERNAL] parsejsonobject"]]
    e__repos_pith_pkg_parser_thneed_go_parse -->|calls| parsejsonobject
    t_parsejsonobject[["[EXTERNAL] t.parsejsonobject"]]
    e__repos_pith_pkg_parser_thneed_go_parse -->|calls| t_parsejsonobject
    parseplain[["[EXTERNAL] parseplain"]]
    e__repos_pith_pkg_parser_thneed_go_parse -->|calls| parseplain
    t_parseplain[["[EXTERNAL] t.parseplain"]]
    e__repos_pith_pkg_parser_thneed_go_parse -->|calls| t_parseplain
    e__repos_pith_pkg_parser_thneed_go_parsejson -->|calls| len
    writestring[["[EXTERNAL] writestring"]]
    e__repos_pith_pkg_parser_thneed_go_parsejson -->|calls| writestring
    sb_writestring[["[EXTERNAL] sb.writestring"]]
    e__repos_pith_pkg_parser_thneed_go_parsejson -->|calls| sb_writestring
    e__repos_pith_pkg_parser_thneed_go_parsejson -->|calls| sprintf
    e__repos_pith_pkg_parser_thneed_go_parsejson -->|calls| fmt_sprintf
    lastindex[["[EXTERNAL] lastindex"]]
    e__repos_pith_pkg_parser_thneed_go_parsejson -->|calls| lastindex
    strings_lastindex[["[EXTERNAL] strings.lastindex"]]
    e__repos_pith_pkg_parser_thneed_go_parsejson -->|calls| strings_lastindex
    e__repos_pith_pkg_parser_thneed_go_parsejson -->|calls| trimspace
    e__repos_pith_pkg_parser_thneed_go_parsejson -->|calls| strings_trimspace
    e__repos_pith_pkg_parser_thneed_go_parsejson -->|calls| replaceall
    e__repos_pith_pkg_parser_thneed_go_parsejson -->|calls| strings_replaceall
    e__repos_pith_pkg_parser_thneed_go_parsejson -->|calls| string
    sb_string[["[EXTERNAL] sb.string"]]
    e__repos_pith_pkg_parser_thneed_go_parsejson -->|calls| sb_string
    e__repos_pith_pkg_parser_thneed_go_parsejsonobject -->|calls| sprintf
    e__repos_pith_pkg_parser_thneed_go_parsejsonobject -->|calls| fmt_sprintf
    e__repos_pith_pkg_parser_thneed_go_parseplain -->|calls| split
    e__repos_pith_pkg_parser_thneed_go_parseplain -->|calls| strings_split
    e__repos_pith_pkg_parser_thneed_go_parseplain -->|calls| trimspace
    e__repos_pith_pkg_parser_thneed_go_parseplain -->|calls| strings_trimspace
    e__repos_pith_pkg_parser_thneed_go_parseplain -->|calls| hasprefix
    e__repos_pith_pkg_parser_thneed_go_parseplain -->|calls| strings_hasprefix
    e__repos_pith_pkg_parser_thneed_go_parseplain -->|calls| append
    e__repos_pith_pkg_parser_thneed_go_parseplain -->|calls| len
    e__repos_pith_pkg_parser_thneed_go_parseplain -->|calls| join
    e__repos_pith_pkg_parser_thneed_go_parseplain -->|calls| strings_join
    e__repos_pith_pkg_parser_thneed_go ==>|contains| e__repos_pith_pkg_parser_thneed_go_parseplain
    e__repos_pith_pkg_parser_thneed_go ==>|contains| e__repos_pith_pkg_parser_thneed_go_thneedparser
    e__repos_pith_pkg_parser_thneed_go ==>|contains| e__repos_pith_pkg_parser_thneed_go_name
    e__repos_pith_pkg_parser_thneed_go ==>|contains| e__repos_pith_pkg_parser_thneed_go_canparse
    e__repos_pith_pkg_parser_thneed_go ==>|contains| e__repos_pith_pkg_parser_thneed_go_parse
    e__repos_pith_pkg_parser_thneed_go ==>|contains| e__repos_pith_pkg_parser_thneed_go_parsejson
    e__repos_pith_pkg_parser_thneed_go ==>|contains| e__repos_pith_pkg_parser_thneed_go_parsejsonobject
    e__repos_pith_pkg_parser_vitest_go -.->|imports|  fmt
    e__repos_pith_pkg_parser_vitest_go -.->|imports|  strings
    e__repos_pith_pkg_parser_vitest_go_canparse -->|calls| matchcommand
    e__repos_pith_pkg_parser_vitest_go_canparse -->|calls| len
    e__repos_pith_pkg_parser_vitest_go_parse -->|calls| split
    e__repos_pith_pkg_parser_vitest_go_parse -->|calls| strings_split
    e__repos_pith_pkg_parser_vitest_go_parse -->|calls| trimspace
    e__repos_pith_pkg_parser_vitest_go_parse -->|calls| strings_trimspace
    e__repos_pith_pkg_parser_vitest_go_parse -->|calls| contains
    e__repos_pith_pkg_parser_vitest_go_parse -->|calls| strings_contains
    e__repos_pith_pkg_parser_vitest_go_parse -->|calls| append
    e__repos_pith_pkg_parser_vitest_go_parse -->|calls| hasprefix
    e__repos_pith_pkg_parser_vitest_go_parse -->|calls| strings_hasprefix
    e__repos_pith_pkg_parser_vitest_go_parse -->|calls| len
    e__repos_pith_pkg_parser_vitest_go_parse -->|calls| join
    e__repos_pith_pkg_parser_vitest_go_parse -->|calls| strings_join
    e__repos_pith_pkg_parser_vitest_go_parse -->|calls| sprintf
    e__repos_pith_pkg_parser_vitest_go_parse -->|calls| fmt_sprintf
    e__repos_pith_pkg_parser_vitest_go ==>|contains| e__repos_pith_pkg_parser_vitest_go_canparse
    e__repos_pith_pkg_parser_vitest_go ==>|contains| e__repos_pith_pkg_parser_vitest_go_parse
    e__repos_pith_pkg_parser_vitest_go ==>|contains| e__repos_pith_pkg_parser_vitest_go_vitestparser
    e__repos_pith_pkg_parser_vitest_go ==>|contains| e__repos_pith_pkg_parser_vitest_go_name
    e__repos_pith_pkg_parser_vitest_test_go -.->|imports|  strings
    e__repos_pith_pkg_parser_vitest_test_go -.->|imports|  testing
    e__repos_pith_pkg_parser_vitest_test_go_testvitestparser -->|calls| canparse
    e__repos_pith_pkg_parser_vitest_test_go_testvitestparser -->|calls| p_canparse
    e__repos_pith_pkg_parser_vitest_test_go_testvitestparser -->|calls| errorf
    e__repos_pith_pkg_parser_vitest_test_go_testvitestparser -->|calls| t_errorf
    e__repos_pith_pkg_parser_vitest_test_go_testvitestparser -->|calls| parse
    e__repos_pith_pkg_parser_vitest_test_go_testvitestparser -->|calls| p_parse
    e__repos_pith_pkg_parser_vitest_test_go_testvitestparser -->|calls| contains
    e__repos_pith_pkg_parser_vitest_test_go_testvitestparser -->|calls| strings_contains
    e__repos_pith_pkg_parser_vitest_test_go_testbdparser -->|calls| canparse
    e__repos_pith_pkg_parser_vitest_test_go_testbdparser -->|calls| p_canparse
    e__repos_pith_pkg_parser_vitest_test_go_testbdparser -->|calls| errorf
    e__repos_pith_pkg_parser_vitest_test_go_testbdparser -->|calls| t_errorf
    e__repos_pith_pkg_parser_vitest_test_go_testbdparser -->|calls| parse
    e__repos_pith_pkg_parser_vitest_test_go_testbdparser -->|calls| p_parse
    e__repos_pith_pkg_parser_vitest_test_go_testbdparser -->|calls| contains
    e__repos_pith_pkg_parser_vitest_test_go_testbdparser -->|calls| strings_contains
    e__repos_pith_pkg_parser_vitest_test_go_testpromptfooparser -->|calls| canparse
    e__repos_pith_pkg_parser_vitest_test_go_testpromptfooparser -->|calls| p_canparse
    e__repos_pith_pkg_parser_vitest_test_go_testpromptfooparser -->|calls| errorf
    e__repos_pith_pkg_parser_vitest_test_go_testpromptfooparser -->|calls| t_errorf
    e__repos_pith_pkg_parser_vitest_test_go_testpromptfooparser -->|calls| parse
    e__repos_pith_pkg_parser_vitest_test_go_testpromptfooparser -->|calls| p_parse
    e__repos_pith_pkg_parser_vitest_test_go_testpromptfooparser -->|calls| contains
    e__repos_pith_pkg_parser_vitest_test_go_testpromptfooparser -->|calls| strings_contains
    e__repos_pith_pkg_parser_vitest_test_go_testpowershellparser -->|calls| canparse
    e__repos_pith_pkg_parser_vitest_test_go_testpowershellparser -->|calls| p_canparse
    e__repos_pith_pkg_parser_vitest_test_go_testpowershellparser -->|calls| errorf
    e__repos_pith_pkg_parser_vitest_test_go_testpowershellparser -->|calls| t_errorf
    e__repos_pith_pkg_parser_vitest_test_go_testpowershellparser -->|calls| parse
    e__repos_pith_pkg_parser_vitest_test_go_testpowershellparser -->|calls| p_parse
    e__repos_pith_pkg_parser_vitest_test_go_testpowershellparser -->|calls| contains
    e__repos_pith_pkg_parser_vitest_test_go_testpowershellparser -->|calls| strings_contains
    e__repos_pith_pkg_parser_vitest_test_go ==>|contains| e__repos_pith_pkg_parser_vitest_test_go_testpromptfooparser
    e__repos_pith_pkg_parser_vitest_test_go ==>|contains| e__repos_pith_pkg_parser_vitest_test_go_testpowershellparser
    e__repos_pith_pkg_parser_vitest_test_go ==>|contains| e__repos_pith_pkg_parser_vitest_test_go_testvitestparser
    e__repos_pith_pkg_parser_vitest_test_go ==>|contains| e__repos_pith_pkg_parser_vitest_test_go_testbdparser
    e__repos_pith_pkg_parser_web_go -.->|imports|  encoding_json
    e__repos_pith_pkg_parser_web_go -.->|imports|  fmt
    e__repos_pith_pkg_parser_web_go -.->|imports|  strings
    e__repos_pith_pkg_parser_web_go_canparse -->|calls| matchcommand
    e__repos_pith_pkg_parser_web_go_parse -->|calls| trimspace
    e__repos_pith_pkg_parser_web_go_parse -->|calls| strings_trimspace
    e__repos_pith_pkg_parser_web_go_parse -->|calls| unmarshal
    e__repos_pith_pkg_parser_web_go_parse -->|calls| json_unmarshal
    e__repos_pith_pkg_parser_web_go_parse -->|calls| marshal
    e__repos_pith_pkg_parser_web_go_parse -->|calls| json_marshal
    e__repos_pith_pkg_parser_web_go_parse -->|calls| len
    e__repos_pith_pkg_parser_web_go_parse -->|calls| make
    e__repos_pith_pkg_parser_web_go_parse -->|calls| append
    e__repos_pith_pkg_parser_web_go_parse -->|calls| sprintf
    e__repos_pith_pkg_parser_web_go_parse -->|calls| fmt_sprintf
    e__repos_pith_pkg_parser_web_go_parse -->|calls| join
    e__repos_pith_pkg_parser_web_go_parse -->|calls| strings_join
    e__repos_pith_pkg_parser_web_go_parse -->|calls| string
    e__repos_pith_pkg_parser_web_go_parse -->|calls| contains
    e__repos_pith_pkg_parser_web_go_parse -->|calls| strings_contains
    e__repos_pith_pkg_parser_web_go_parse -->|calls| index
    e__repos_pith_pkg_parser_web_go_parse -->|calls| strings_index
    e__repos_pith_pkg_parser_web_go ==>|contains| e__repos_pith_pkg_parser_web_go_webparser
    e__repos_pith_pkg_parser_web_go ==>|contains| e__repos_pith_pkg_parser_web_go_name
    e__repos_pith_pkg_parser_web_go ==>|contains| e__repos_pith_pkg_parser_web_go_canparse
    e__repos_pith_pkg_parser_web_go ==>|contains| e__repos_pith_pkg_parser_web_go_parse
    bytes[["[EXTERNAL] bytes"]]
    e__repos_pith_pkg_runner_runner_go -.->|imports|  bytes
    e__repos_pith_pkg_runner_runner_go -.->|imports|  pith_pkg_config
    e__repos_pith_pkg_runner_runner_go -.->|imports|  pith_pkg_parser
    e__repos_pith_pkg_runner_runner_go -.->|imports|  pith_pkg_telemetry
    e__repos_pith_pkg_runner_runner_go -.->|imports|  fmt
    e__repos_pith_pkg_runner_runner_go -.->|imports|  os
    e__repos_pith_pkg_runner_runner_go -.->|imports|  os_exec
    e__repos_pith_pkg_runner_runner_go -.->|imports|  path_filepath
    e__repos_pith_pkg_runner_runner_go -.->|imports|  strings
    e__repos_pith_pkg_runner_runner_go -.->|imports|  time
    unicode_utf8[["[EXTERNAL] utf8"]]
    e__repos_pith_pkg_runner_runner_go -.->|imports|  unicode_utf8
    e__repos_pith_pkg_runner_runner_go_newrunner -->|calls| getallparsers
    e__repos_pith_pkg_runner_runner_go_newrunner -->|calls| parser_getallparsers
    detectsource[["[EXTERNAL] detectsource"]]
    e__repos_pith_pkg_runner_runner_go_newrunner -->|calls| detectsource
    e__repos_pith_pkg_runner_runner_go_detectsource -->|calls| getenv
    e__repos_pith_pkg_runner_runner_go_detectsource -->|calls| os_getenv
    e__repos_pith_pkg_runner_runner_go_run -->|calls| runwithoptions
    r_runwithoptions[["[EXTERNAL] r.runwithoptions"]]
    e__repos_pith_pkg_runner_runner_go_run -->|calls| r_runwithoptions
    e__repos_pith_pkg_runner_runner_go_logforsnag -->|calls| userhomedir
    e__repos_pith_pkg_runner_runner_go_logforsnag -->|calls| os_userhomedir
    e__repos_pith_pkg_runner_runner_go_logforsnag -->|calls| join
    e__repos_pith_pkg_runner_runner_go_logforsnag -->|calls| filepath_join
    e__repos_pith_pkg_runner_runner_go_logforsnag -->|calls| mkdirall
    e__repos_pith_pkg_runner_runner_go_logforsnag -->|calls| os_mkdirall
    openfile[["[EXTERNAL] openfile"]]
    e__repos_pith_pkg_runner_runner_go_logforsnag -->|calls| openfile
    os_openfile[["[EXTERNAL] os.openfile"]]
    e__repos_pith_pkg_runner_runner_go_logforsnag -->|calls| os_openfile
    e__repos_pith_pkg_runner_runner_go_logforsnag -->|calls| close
    f_close[["[EXTERNAL] f.close"]]
    e__repos_pith_pkg_runner_runner_go_logforsnag -->|calls| f_close
    e__repos_pith_pkg_runner_runner_go_logforsnag -->|calls| split
    e__repos_pith_pkg_runner_runner_go_logforsnag -->|calls| strings_split
    e__repos_pith_pkg_runner_runner_go_logforsnag -->|calls| len
    e__repos_pith_pkg_runner_runner_go_logforsnag -->|calls| strings_join
    e__repos_pith_pkg_runner_runner_go_logforsnag -->|calls| sprintf
    e__repos_pith_pkg_runner_runner_go_logforsnag -->|calls| fmt_sprintf
    e__repos_pith_pkg_runner_runner_go_logforsnag -->|calls| hassuffix
    e__repos_pith_pkg_runner_runner_go_logforsnag -->|calls| strings_hassuffix
    e__repos_pith_pkg_runner_runner_go_logforsnag -->|calls| writestring
    f_writestring[["[EXTERNAL] f.writestring"]]
    e__repos_pith_pkg_runner_runner_go_logforsnag -->|calls| f_writestring
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| len
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| errorf
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| fmt_errorf
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| join
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| strings_join
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| now
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| time_now
    containsany[["[EXTERNAL] containsany"]]
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| containsany
    strings_containsany[["[EXTERNAL] strings.containsany"]]
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| strings_containsany
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| command
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| exec_command
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| contains
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| strings_contains
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| fields
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| strings_fields
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| run
    cmd_run[["[EXTERNAL] cmd.run"]]
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| cmd_run
    since[["[EXTERNAL] since"]]
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| since
    time_since[["[EXTERNAL] time.since"]]
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| time_since
    milliseconds[["[EXTERNAL] milliseconds"]]
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| milliseconds
    exitcode[["[EXTERNAL] exitcode"]]
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| exitcode
    exiterr_exitcode[["[EXTERNAL] exiterr.exitcode"]]
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| exiterr_exitcode
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| string
    out_string[["[EXTERNAL] out.string"]]
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| out_string
    stderr_string[["[EXTERNAL] stderr.string"]]
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| stderr_string
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| logforsnag
    r_logforsnag[["[EXTERNAL] r.logforsnag"]]
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| r_logforsnag
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| estimatetokens
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| splitsubcommands
    cp_splitsubcommands[["[EXTERNAL] cp.splitsubcommands"]]
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| cp_splitsubcommands
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| name
    pcandidate_name[["[EXTERNAL] pcandidate.name"]]
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| pcandidate_name
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| canparse
    pcandidate_canparse[["[EXTERNAL] pcandidate.canparse"]]
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| pcandidate_canparse
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| parser_name
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| parser_canparse
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| parse
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| p_parse
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| p_name
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| applymiddleouttruncation
    r_applymiddleouttruncation[["[EXTERNAL] r.applymiddleouttruncation"]]
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| r_applymiddleouttruncation
    print[["[EXTERNAL] print"]]
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| print
    fmt_print[["[EXTERNAL] fmt.print"]]
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| fmt_print
    e__repos_pith_pkg_runner_runner_go_runwithoptions -->|calls| record
    e__repos_pith_pkg_runner_runner_go_applymiddleouttruncation -->|calls| split
    e__repos_pith_pkg_runner_runner_go_applymiddleouttruncation -->|calls| strings_split
    e__repos_pith_pkg_runner_runner_go_applymiddleouttruncation -->|calls| len
    e__repos_pith_pkg_runner_runner_go_applymiddleouttruncation -->|calls| append
    e__repos_pith_pkg_runner_runner_go_applymiddleouttruncation -->|calls| sprintf
    e__repos_pith_pkg_runner_runner_go_applymiddleouttruncation -->|calls| fmt_sprintf
    e__repos_pith_pkg_runner_runner_go_applymiddleouttruncation -->|calls| join
    e__repos_pith_pkg_runner_runner_go_applymiddleouttruncation -->|calls| strings_join
    runecountinstring[["[EXTERNAL] runecountinstring"]]
    e__repos_pith_pkg_runner_runner_go_estimatetokens -->|calls| runecountinstring
    utf8_runecountinstring[["[EXTERNAL] utf8.runecountinstring"]]
    e__repos_pith_pkg_runner_runner_go_estimatetokens -->|calls| utf8_runecountinstring
    e__repos_pith_pkg_runner_runner_go ==>|contains| e__repos_pith_pkg_runner_runner_go_logforsnag
    e__repos_pith_pkg_runner_runner_go ==>|contains| e__repos_pith_pkg_runner_runner_go_runwithoptions
    e__repos_pith_pkg_runner_runner_go ==>|contains| e__repos_pith_pkg_runner_runner_go_applymiddleouttruncation
    e__repos_pith_pkg_runner_runner_go ==>|contains| e__repos_pith_pkg_runner_runner_go_estimatetokens
    e__repos_pith_pkg_runner_runner_go ==>|contains| e__repos_pith_pkg_runner_runner_go_runner
    e__repos_pith_pkg_runner_runner_go ==>|contains| e__repos_pith_pkg_runner_runner_go_newrunner
    e__repos_pith_pkg_runner_runner_go ==>|contains| e__repos_pith_pkg_runner_runner_go_detectsource
    e__repos_pith_pkg_runner_runner_go ==>|contains| e__repos_pith_pkg_runner_runner_go_run
    e__repos_pith_pkg_runner_runner_test_go -.->|imports|  pith_pkg_config
    e__repos_pith_pkg_runner_runner_test_go -.->|imports|  pith_pkg_telemetry
    e__repos_pith_pkg_runner_runner_test_go -.->|imports|  os
    e__repos_pith_pkg_runner_runner_test_go -.->|imports|  path_filepath
    e__repos_pith_pkg_runner_runner_test_go -.->|imports|  strings
    e__repos_pith_pkg_runner_runner_test_go -.->|imports|  testing
    e__repos_pith_pkg_runner_runner_test_go_testestimatetokens -->|calls| estimatetokens
    e__repos_pith_pkg_runner_runner_test_go_testestimatetokens -->|calls| errorf
    e__repos_pith_pkg_runner_runner_test_go_testestimatetokens -->|calls| t_errorf
    e__repos_pith_pkg_runner_runner_test_go_testrunner -->|calls| mkdirtemp
    e__repos_pith_pkg_runner_runner_test_go_testrunner -->|calls| os_mkdirtemp
    e__repos_pith_pkg_runner_runner_test_go_testrunner -->|calls| removeall
    e__repos_pith_pkg_runner_runner_test_go_testrunner -->|calls| os_removeall
    e__repos_pith_pkg_runner_runner_test_go_testrunner -->|calls| join
    e__repos_pith_pkg_runner_runner_test_go_testrunner -->|calls| filepath_join
    newtelemetrywithpath[["[EXTERNAL] newtelemetrywithpath"]]
    e__repos_pith_pkg_runner_runner_test_go_testrunner -->|calls| newtelemetrywithpath
    telemetry_newtelemetrywithpath[["[EXTERNAL] telemetry.newtelemetrywithpath"]]
    e__repos_pith_pkg_runner_runner_test_go_testrunner -->|calls| telemetry_newtelemetrywithpath
    e__repos_pith_pkg_runner_runner_test_go_testrunner -->|calls| close
    e__repos_pith_pkg_runner_runner_test_go_testrunner -->|calls| tel_close
    e__repos_pith_pkg_runner_runner_test_go_testrunner -->|calls| newrunner
    e__repos_pith_pkg_runner_runner_test_go_testrunner -->|calls| len
    e__repos_pith_pkg_runner_runner_test_go_testrunner -->|calls| error
    e__repos_pith_pkg_runner_runner_test_go_testrunner -->|calls| t_error
    e__repos_pith_pkg_runner_runner_test_go_testmiddleouttruncation -->|calls| newrunner
    e__repos_pith_pkg_runner_runner_test_go_testmiddleouttruncation -->|calls| applymiddleouttruncation
    e__repos_pith_pkg_runner_runner_test_go_testmiddleouttruncation -->|calls| run_applymiddleouttruncation
    e__repos_pith_pkg_runner_runner_test_go_testmiddleouttruncation -->|calls| contains
    e__repos_pith_pkg_runner_runner_test_go_testmiddleouttruncation -->|calls| strings_contains
    e__repos_pith_pkg_runner_runner_test_go_testmiddleouttruncation -->|calls| error
    e__repos_pith_pkg_runner_runner_test_go_testmiddleouttruncation -->|calls| t_error
    e__repos_pith_pkg_runner_runner_test_go_testmiddleouttruncation -->|calls| hasprefix
    e__repos_pith_pkg_runner_runner_test_go_testmiddleouttruncation -->|calls| strings_hasprefix
    e__repos_pith_pkg_runner_runner_test_go_testmiddleouttruncation -->|calls| hassuffix
    e__repos_pith_pkg_runner_runner_test_go_testmiddleouttruncation -->|calls| strings_hassuffix
    e__repos_pith_pkg_runner_runner_test_go ==>|contains| e__repos_pith_pkg_runner_runner_test_go_testestimatetokens
    e__repos_pith_pkg_runner_runner_test_go ==>|contains| e__repos_pith_pkg_runner_runner_test_go_testrunner
    e__repos_pith_pkg_runner_runner_test_go ==>|contains| e__repos_pith_pkg_runner_runner_test_go_testmiddleouttruncation
    e__repos_pith_pkg_selfupdate_selfupdate_go -.->|imports|  bytes
    e__repos_pith_pkg_selfupdate_selfupdate_go -.->|imports|  encoding_json
    e__repos_pith_pkg_selfupdate_selfupdate_go -.->|imports|  fmt
    e__repos_pith_pkg_selfupdate_selfupdate_go -.->|imports|  io
    e__repos_pith_pkg_selfupdate_selfupdate_go -.->|imports|  net_http
    e__repos_pith_pkg_selfupdate_selfupdate_go -.->|imports|  os
    e__repos_pith_pkg_selfupdate_selfupdate_go -.->|imports|  os_exec
    e__repos_pith_pkg_selfupdate_selfupdate_go -.->|imports|  runtime
    e__repos_pith_pkg_selfupdate_selfupdate_go -.->|imports|  strings
    e__repos_pith_pkg_selfupdate_selfupdate_go_getauthtoken -->|calls| getenv
    e__repos_pith_pkg_selfupdate_selfupdate_go_getauthtoken -->|calls| os_getenv
    e__repos_pith_pkg_selfupdate_selfupdate_go_getauthtoken -->|calls| command
    e__repos_pith_pkg_selfupdate_selfupdate_go_getauthtoken -->|calls| exec_command
    e__repos_pith_pkg_selfupdate_selfupdate_go_getauthtoken -->|calls| run
    e__repos_pith_pkg_selfupdate_selfupdate_go_getauthtoken -->|calls| cmd_run
    e__repos_pith_pkg_selfupdate_selfupdate_go_getauthtoken -->|calls| trimspace
    e__repos_pith_pkg_selfupdate_selfupdate_go_getauthtoken -->|calls| strings_trimspace
    e__repos_pith_pkg_selfupdate_selfupdate_go_getauthtoken -->|calls| string
    e__repos_pith_pkg_selfupdate_selfupdate_go_getauthtoken -->|calls| out_string
    newrequest[["[EXTERNAL] newrequest"]]
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkandapplyupdate -->|calls| newrequest
    http_newrequest[["[EXTERNAL] http.newrequest"]]
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkandapplyupdate -->|calls| http_newrequest
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkandapplyupdate -->|calls| sprintf
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkandapplyupdate -->|calls| fmt_sprintf
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkandapplyupdate -->|calls| set
    getauthtoken[["[EXTERNAL] getauthtoken"]]
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkandapplyupdate -->|calls| getauthtoken
    do[["[EXTERNAL] do"]]
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkandapplyupdate -->|calls| do
    client_do[["[EXTERNAL] client.do"]]
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkandapplyupdate -->|calls| client_do
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkandapplyupdate -->|calls| close
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkandapplyupdate -->|calls| errorf
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkandapplyupdate -->|calls| fmt_errorf
    newdecoder[["[EXTERNAL] newdecoder"]]
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkandapplyupdate -->|calls| newdecoder
    json_newdecoder[["[EXTERNAL] json.newdecoder"]]
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkandapplyupdate -->|calls| json_newdecoder
    decode[["[EXTERNAL] decode"]]
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkandapplyupdate -->|calls| decode
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkandapplyupdate -->|calls| len
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkandapplyupdate -->|calls| printf
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkandapplyupdate -->|calls| fmt_printf
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkandapplyupdate -->|calls| hassuffix
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkandapplyupdate -->|calls| strings_hassuffix
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkandapplyupdate -->|calls| println
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkandapplyupdate -->|calls| fmt_println
    downloadandreplace[["[EXTERNAL] downloadandreplace"]]
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkandapplyupdate -->|calls| downloadandreplace
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkforupdatesilent -->|calls| newrequest
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkforupdatesilent -->|calls| http_newrequest
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkforupdatesilent -->|calls| sprintf
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkforupdatesilent -->|calls| fmt_sprintf
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkforupdatesilent -->|calls| set
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkforupdatesilent -->|calls| getauthtoken
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkforupdatesilent -->|calls| do
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkforupdatesilent -->|calls| client_do
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkforupdatesilent -->|calls| close
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkforupdatesilent -->|calls| errorf
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkforupdatesilent -->|calls| fmt_errorf
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkforupdatesilent -->|calls| newdecoder
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkforupdatesilent -->|calls| json_newdecoder
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkforupdatesilent -->|calls| decode
    e__repos_pith_pkg_selfupdate_selfupdate_go_checkforupdatesilent -->|calls| len
    e__repos_pith_pkg_selfupdate_selfupdate_go_downloadandreplace -->|calls| executable
    e__repos_pith_pkg_selfupdate_selfupdate_go_downloadandreplace -->|calls| os_executable
    e__repos_pith_pkg_selfupdate_selfupdate_go_downloadandreplace -->|calls| newrequest
    e__repos_pith_pkg_selfupdate_selfupdate_go_downloadandreplace -->|calls| http_newrequest
    e__repos_pith_pkg_selfupdate_selfupdate_go_downloadandreplace -->|calls| set
    e__repos_pith_pkg_selfupdate_selfupdate_go_downloadandreplace -->|calls| do
    e__repos_pith_pkg_selfupdate_selfupdate_go_downloadandreplace -->|calls| client_do
    e__repos_pith_pkg_selfupdate_selfupdate_go_downloadandreplace -->|calls| close
    e__repos_pith_pkg_selfupdate_selfupdate_go_downloadandreplace -->|calls| errorf
    e__repos_pith_pkg_selfupdate_selfupdate_go_downloadandreplace -->|calls| fmt_errorf
    e__repos_pith_pkg_selfupdate_selfupdate_go_downloadandreplace -->|calls| openfile
    e__repos_pith_pkg_selfupdate_selfupdate_go_downloadandreplace -->|calls| os_openfile
    e__repos_pith_pkg_selfupdate_selfupdate_go_downloadandreplace -->|calls| copy
    e__repos_pith_pkg_selfupdate_selfupdate_go_downloadandreplace -->|calls| io_copy
    tmpfile_close[["[EXTERNAL] tmpfile.close"]]
    e__repos_pith_pkg_selfupdate_selfupdate_go_downloadandreplace -->|calls| tmpfile_close
    e__repos_pith_pkg_selfupdate_selfupdate_go_downloadandreplace -->|calls| remove
    e__repos_pith_pkg_selfupdate_selfupdate_go_downloadandreplace -->|calls| os_remove
    e__repos_pith_pkg_selfupdate_selfupdate_go_downloadandreplace -->|calls| rename
    e__repos_pith_pkg_selfupdate_selfupdate_go_downloadandreplace -->|calls| os_rename
    e__repos_pith_pkg_selfupdate_selfupdate_go_downloadandreplace -->|calls| println
    e__repos_pith_pkg_selfupdate_selfupdate_go_downloadandreplace -->|calls| fmt_println
    e__repos_pith_pkg_selfupdate_selfupdate_go ==>|contains| e__repos_pith_pkg_selfupdate_selfupdate_go_getauthtoken
    e__repos_pith_pkg_selfupdate_selfupdate_go ==>|contains| e__repos_pith_pkg_selfupdate_selfupdate_go_checkandapplyupdate
    e__repos_pith_pkg_selfupdate_selfupdate_go ==>|contains| e__repos_pith_pkg_selfupdate_selfupdate_go_checkforupdatesilent
    e__repos_pith_pkg_selfupdate_selfupdate_go ==>|contains| e__repos_pith_pkg_selfupdate_selfupdate_go_downloadandreplace
    e__repos_pith_pkg_selfupdate_selfupdate_go ==>|contains| e__repos_pith_pkg_selfupdate_selfupdate_go_release
    database_sql[["[EXTERNAL] sql"]]
    e__repos_pith_pkg_telemetry_telemetry_go -.->|imports|  database_sql
    e__repos_pith_pkg_telemetry_telemetry_go -.->|imports|  fmt
    e__repos_pith_pkg_telemetry_telemetry_go -.->|imports|  os
    e__repos_pith_pkg_telemetry_telemetry_go -.->|imports|  path_filepath
    e__repos_pith_pkg_telemetry_telemetry_go -.->|imports|  time
    modernc_org_sqlite[["[EXTERNAL] sqlite"]]
    e__repos_pith_pkg_telemetry_telemetry_go -.->|imports|  modernc_org_sqlite
    e__repos_pith_pkg_telemetry_telemetry_go_newtelemetry -->|calls| userhomedir
    e__repos_pith_pkg_telemetry_telemetry_go_newtelemetry -->|calls| os_userhomedir
    e__repos_pith_pkg_telemetry_telemetry_go_newtelemetry -->|calls| join
    e__repos_pith_pkg_telemetry_telemetry_go_newtelemetry -->|calls| filepath_join
    e__repos_pith_pkg_telemetry_telemetry_go_newtelemetry -->|calls| mkdirall
    e__repos_pith_pkg_telemetry_telemetry_go_newtelemetry -->|calls| os_mkdirall
    e__repos_pith_pkg_telemetry_telemetry_go_newtelemetry -->|calls| newtelemetrywithpath
    e__repos_pith_pkg_telemetry_telemetry_go_newtelemetrywithpath -->|calls| open
    sql_open[["[EXTERNAL] sql.open"]]
    e__repos_pith_pkg_telemetry_telemetry_go_newtelemetrywithpath -->|calls| sql_open
    init[["[EXTERNAL] init"]]
    e__repos_pith_pkg_telemetry_telemetry_go_newtelemetrywithpath -->|calls| init
    t_init[["[EXTERNAL] t.init"]]
    e__repos_pith_pkg_telemetry_telemetry_go_newtelemetrywithpath -->|calls| t_init
    exec[["[EXTERNAL] exec"]]
    e__repos_pith_pkg_telemetry_telemetry_go_init -->|calls| exec
    e__repos_pith_pkg_telemetry_telemetry_go_record -->|calls| exec
    e__repos_pith_pkg_telemetry_telemetry_go_getstatsbyday -->|calls| append
    e__repos_pith_pkg_telemetry_telemetry_go_getstatsbyday -->|calls| sprintf
    e__repos_pith_pkg_telemetry_telemetry_go_getstatsbyday -->|calls| fmt_sprintf
    e__repos_pith_pkg_telemetry_telemetry_go_getstatsbyday -->|calls| query
    e__repos_pith_pkg_telemetry_telemetry_go_getstatsbyday -->|calls| close
    rows_close[["[EXTERNAL] rows.close"]]
    e__repos_pith_pkg_telemetry_telemetry_go_getstatsbyday -->|calls| rows_close
    next[["[EXTERNAL] next"]]
    e__repos_pith_pkg_telemetry_telemetry_go_getstatsbyday -->|calls| next
    rows_next[["[EXTERNAL] rows.next"]]
    e__repos_pith_pkg_telemetry_telemetry_go_getstatsbyday -->|calls| rows_next
    scan[["[EXTERNAL] scan"]]
    e__repos_pith_pkg_telemetry_telemetry_go_getstatsbyday -->|calls| scan
    rows_scan[["[EXTERNAL] rows.scan"]]
    e__repos_pith_pkg_telemetry_telemetry_go_getstatsbyday -->|calls| rows_scan
    e__repos_pith_pkg_telemetry_telemetry_go_close -->|calls| close
    e__repos_pith_pkg_telemetry_telemetry_go_getunparsedcommands -->|calls| append
    e__repos_pith_pkg_telemetry_telemetry_go_getunparsedcommands -->|calls| sprintf
    e__repos_pith_pkg_telemetry_telemetry_go_getunparsedcommands -->|calls| fmt_sprintf
    e__repos_pith_pkg_telemetry_telemetry_go_getunparsedcommands -->|calls| query
    e__repos_pith_pkg_telemetry_telemetry_go_getunparsedcommands -->|calls| close
    e__repos_pith_pkg_telemetry_telemetry_go_getunparsedcommands -->|calls| rows_close
    e__repos_pith_pkg_telemetry_telemetry_go_getunparsedcommands -->|calls| next
    e__repos_pith_pkg_telemetry_telemetry_go_getunparsedcommands -->|calls| rows_next
    e__repos_pith_pkg_telemetry_telemetry_go_getunparsedcommands -->|calls| scan
    e__repos_pith_pkg_telemetry_telemetry_go_getunparsedcommands -->|calls| rows_scan
    e__repos_pith_pkg_telemetry_telemetry_go_getstats -->|calls| append
    e__repos_pith_pkg_telemetry_telemetry_go_getstats -->|calls| sprintf
    e__repos_pith_pkg_telemetry_telemetry_go_getstats -->|calls| fmt_sprintf
    queryrow[["[EXTERNAL] queryrow"]]
    e__repos_pith_pkg_telemetry_telemetry_go_getstats -->|calls| queryrow
    e__repos_pith_pkg_telemetry_telemetry_go_getstats -->|calls| scan
    e__repos_pith_pkg_telemetry_telemetry_go_getstatsbycommand -->|calls| append
    e__repos_pith_pkg_telemetry_telemetry_go_getstatsbycommand -->|calls| sprintf
    e__repos_pith_pkg_telemetry_telemetry_go_getstatsbycommand -->|calls| fmt_sprintf
    e__repos_pith_pkg_telemetry_telemetry_go_getstatsbycommand -->|calls| query
    e__repos_pith_pkg_telemetry_telemetry_go_getstatsbycommand -->|calls| close
    e__repos_pith_pkg_telemetry_telemetry_go_getstatsbycommand -->|calls| rows_close
    e__repos_pith_pkg_telemetry_telemetry_go_getstatsbycommand -->|calls| next
    e__repos_pith_pkg_telemetry_telemetry_go_getstatsbycommand -->|calls| rows_next
    e__repos_pith_pkg_telemetry_telemetry_go_getstatsbycommand -->|calls| scan
    e__repos_pith_pkg_telemetry_telemetry_go_getstatsbycommand -->|calls| rows_scan
    e__repos_pith_pkg_telemetry_telemetry_go_resetall -->|calls| exec
    e__repos_pith_pkg_telemetry_telemetry_go_resetpassthrough -->|calls| exec
    e__repos_pith_pkg_telemetry_telemetry_go_getrecentexecutions -->|calls| append
    e__repos_pith_pkg_telemetry_telemetry_go_getrecentexecutions -->|calls| sprintf
    e__repos_pith_pkg_telemetry_telemetry_go_getrecentexecutions -->|calls| fmt_sprintf
    e__repos_pith_pkg_telemetry_telemetry_go_getrecentexecutions -->|calls| query
    e__repos_pith_pkg_telemetry_telemetry_go_getrecentexecutions -->|calls| close
    e__repos_pith_pkg_telemetry_telemetry_go_getrecentexecutions -->|calls| rows_close
    e__repos_pith_pkg_telemetry_telemetry_go_getrecentexecutions -->|calls| next
    e__repos_pith_pkg_telemetry_telemetry_go_getrecentexecutions -->|calls| rows_next
    e__repos_pith_pkg_telemetry_telemetry_go_getrecentexecutions -->|calls| scan
    e__repos_pith_pkg_telemetry_telemetry_go_getrecentexecutions -->|calls| rows_scan
    e__repos_pith_pkg_telemetry_telemetry_go_getsources -->|calls| query
    e__repos_pith_pkg_telemetry_telemetry_go_getsources -->|calls| close
    e__repos_pith_pkg_telemetry_telemetry_go_getsources -->|calls| rows_close
    e__repos_pith_pkg_telemetry_telemetry_go_getsources -->|calls| next
    e__repos_pith_pkg_telemetry_telemetry_go_getsources -->|calls| rows_next
    e__repos_pith_pkg_telemetry_telemetry_go_getsources -->|calls| scan
    e__repos_pith_pkg_telemetry_telemetry_go_getsources -->|calls| rows_scan
    e__repos_pith_pkg_telemetry_telemetry_go_getsources -->|calls| append
    e__repos_pith_pkg_telemetry_telemetry_go_getexecutiondetails -->|calls| queryrow
    e__repos_pith_pkg_telemetry_telemetry_go_getexecutiondetails -->|calls| scan
    e__repos_pith_pkg_telemetry_telemetry_go ==>|contains| e__repos_pith_pkg_telemetry_telemetry_go_newtelemetry
    e__repos_pith_pkg_telemetry_telemetry_go ==>|contains| e__repos_pith_pkg_telemetry_telemetry_go_record
    e__repos_pith_pkg_telemetry_telemetry_go ==>|contains| e__repos_pith_pkg_telemetry_telemetry_go_getrecentexecutions
    e__repos_pith_pkg_telemetry_telemetry_go ==>|contains| e__repos_pith_pkg_telemetry_telemetry_go_getsources
    e__repos_pith_pkg_telemetry_telemetry_go ==>|contains| e__repos_pith_pkg_telemetry_telemetry_go_newtelemetrywithpath
    e__repos_pith_pkg_telemetry_telemetry_go ==>|contains| e__repos_pith_pkg_telemetry_telemetry_go_init
    e__repos_pith_pkg_telemetry_telemetry_go ==>|contains| e__repos_pith_pkg_telemetry_telemetry_go_getstats
    e__repos_pith_pkg_telemetry_telemetry_go ==>|contains| e__repos_pith_pkg_telemetry_telemetry_go_getstatsbycommand
    e__repos_pith_pkg_telemetry_telemetry_go ==>|contains| e__repos_pith_pkg_telemetry_telemetry_go_resetpassthrough
    e__repos_pith_pkg_telemetry_telemetry_go ==>|contains| e__repos_pith_pkg_telemetry_telemetry_go_telemetry
    e__repos_pith_pkg_telemetry_telemetry_go ==>|contains| e__repos_pith_pkg_telemetry_telemetry_go_close
    e__repos_pith_pkg_telemetry_telemetry_go ==>|contains| e__repos_pith_pkg_telemetry_telemetry_go_executionrecord
    e__repos_pith_pkg_telemetry_telemetry_go ==>|contains| e__repos_pith_pkg_telemetry_telemetry_go_getstatsbyday
    e__repos_pith_pkg_telemetry_telemetry_go ==>|contains| e__repos_pith_pkg_telemetry_telemetry_go_getunparsedcommands
    e__repos_pith_pkg_telemetry_telemetry_go ==>|contains| e__repos_pith_pkg_telemetry_telemetry_go_resetall
    e__repos_pith_pkg_telemetry_telemetry_go ==>|contains| e__repos_pith_pkg_telemetry_telemetry_go_getexecutiondetails
    e__repos_pith_pkg_telemetry_telemetry_test_go -.->|imports|  os
    e__repos_pith_pkg_telemetry_telemetry_test_go -.->|imports|  path_filepath
    e__repos_pith_pkg_telemetry_telemetry_test_go -.->|imports|  testing
    e__repos_pith_pkg_telemetry_telemetry_test_go_testtelemetry -->|calls| mkdirtemp
    e__repos_pith_pkg_telemetry_telemetry_test_go_testtelemetry -->|calls| os_mkdirtemp
    e__repos_pith_pkg_telemetry_telemetry_test_go_testtelemetry -->|calls| fatal
    e__repos_pith_pkg_telemetry_telemetry_test_go_testtelemetry -->|calls| t_fatal
    e__repos_pith_pkg_telemetry_telemetry_test_go_testtelemetry -->|calls| removeall
    e__repos_pith_pkg_telemetry_telemetry_test_go_testtelemetry -->|calls| os_removeall
    e__repos_pith_pkg_telemetry_telemetry_test_go_testtelemetry -->|calls| join
    e__repos_pith_pkg_telemetry_telemetry_test_go_testtelemetry -->|calls| filepath_join
    e__repos_pith_pkg_telemetry_telemetry_test_go_testtelemetry -->|calls| newtelemetrywithpath
    fatalf[["[EXTERNAL] fatalf"]]
    e__repos_pith_pkg_telemetry_telemetry_test_go_testtelemetry -->|calls| fatalf
    t_fatalf[["[EXTERNAL] t.fatalf"]]
    e__repos_pith_pkg_telemetry_telemetry_test_go_testtelemetry -->|calls| t_fatalf
    e__repos_pith_pkg_telemetry_telemetry_test_go_testtelemetry -->|calls| close
    e__repos_pith_pkg_telemetry_telemetry_test_go_testtelemetry -->|calls| tel_close
    e__repos_pith_pkg_telemetry_telemetry_test_go_testtelemetry -->|calls| record
    e__repos_pith_pkg_telemetry_telemetry_test_go_testtelemetry -->|calls| tel_record
    e__repos_pith_pkg_telemetry_telemetry_test_go_testtelemetry -->|calls| getstats
    e__repos_pith_pkg_telemetry_telemetry_test_go_testtelemetry -->|calls| tel_getstats
    e__repos_pith_pkg_telemetry_telemetry_test_go_testtelemetry -->|calls| errorf
    e__repos_pith_pkg_telemetry_telemetry_test_go_testtelemetry -->|calls| t_errorf
    e__repos_pith_pkg_telemetry_telemetry_test_go_testtelemetry -->|calls| getstatsbycommand
    e__repos_pith_pkg_telemetry_telemetry_test_go_testtelemetry -->|calls| tel_getstatsbycommand
    e__repos_pith_pkg_telemetry_telemetry_test_go_testtelemetry -->|calls| len
    e__repos_pith_pkg_telemetry_telemetry_test_go_testtelemetry -->|calls| getunparsedcommands
    e__repos_pith_pkg_telemetry_telemetry_test_go_testtelemetry -->|calls| tel_getunparsedcommands
    e__repos_pith_pkg_telemetry_telemetry_test_go_testtelemetry -->|calls| error
    e__repos_pith_pkg_telemetry_telemetry_test_go_testtelemetry -->|calls| t_error
    e__repos_pith_pkg_telemetry_telemetry_test_go_testtelemetry -->|calls| resetpassthrough
    e__repos_pith_pkg_telemetry_telemetry_test_go_testtelemetry -->|calls| tel_resetpassthrough
    e__repos_pith_pkg_telemetry_telemetry_test_go_testtelemetry -->|calls| resetall
    e__repos_pith_pkg_telemetry_telemetry_test_go_testtelemetry -->|calls| tel_resetall
    e__repos_pith_pkg_telemetry_telemetry_test_go ==>|contains| e__repos_pith_pkg_telemetry_telemetry_test_go_testtelemetry
    e__repos_pith_pkg_parser[["[EXTERNAL] parser"]]
    beads_pith_bhx --> e__repos_pith_pkg_parser
    beads_pith_wth --> e__repos_pith_pkg_parser
    beads_pith_os4 --> e__repos_pith_pkg_parser
    beads_pith_4uf --> e__repos_pith_pkg_parser
    beads_pith_0d3 --> e__repos_pith_pkg_parser
    e__repos_pith_pkg_parser_discovery_go[["[EXTERNAL] discovery.go"]]
    beads_pith_0d3 --> e__repos_pith_pkg_parser_discovery_go
    beads_pith_41w --> e__repos_pith_pkg_parser
```
<!-- mermaid-end -->
