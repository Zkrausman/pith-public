---
name: pith-parser-generator
description: Analyzes Pith's discovery queue of unparsed CLI commands and suggests new Go-based Parser implementations to optimize token usage. Use when the user wants to expand Pith's coverage based on real-world usage patterns.
---

# Pith Parser Generator

This skill guides the creation of new parsers for Pith, a token-optimized CLI proxy. Pith records commands that it cannot parse in a "discovery queue". By analyzing this queue, we can identify high-impact commands that would benefit from a dedicated parser.

## Workflow

### 1. Identify Candidates
Run `pith discover` to see the top unparsed commands. Look for commands with high `InvocationCount` and `TotalRawTokens`.

```bash
pith discover
```

### 2. Analyze Raw Output
For a chosen candidate command, run it through Pith with the `raw` command to see exactly what the unparsed output looks like.

```bash
pith raw <command>
```

### 3. Design the Parser
A parser must implement the `Parser` interface in `pkg/parser/interface.go`:

```go
type Parser interface {
	Name() string
	CanParse(cmd string, args []string) bool
	Parse(output string) string
}
```

- **Name()**: Returns a unique string identifying the parser (e.g., "docker-ps").
- **CanParse()**: Returns true if the parser can handle the given command name and arguments.
- **Parse()**: Takes the raw output string and returns a compressed/optimized version.

### 4. Implement the Parser
Create a new file in `pkg/parser/` (e.g., `pkg/parser/docker.go`) and implement your parser(s).

**Best Practices for Compression:**
- Remove decorative ASCII art, borders, and excessive whitespace.
- Truncate long lists to a reasonable limit (e.g., 50 items).
- Filter out redundant columns (e.g., file permissions if only name/size is needed).
- Use `strings.Fields` and `strings.Join` for efficient line processing.

### 5. Register the Parser
Add your new parser to the `GetAllParsers()` function in `pkg/parser/interface.go`.

### 6. Verify and Test
Run the command again (without `raw`) and verify that Pith now uses your new parser and produces optimized output.

```bash
pith <command>
```

Check `pith gain` to see the token savings recorded for your new parser.
