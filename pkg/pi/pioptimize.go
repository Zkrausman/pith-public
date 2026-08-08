package pi

import (
	"fmt"
	"regexp"
	"strings"
)

// PiConfig controls PiOptimize behavior. Zero value is valid.
type PiConfig struct {
	// ThresholdBytes: compress only outputs >= this size. 0 = default 8000.
	ThresholdBytes int
	// TelemetryEnabled controls whether telemetry is recorded for Pi.
	// Must remain false by default per ZAR-110.
	TelemetryEnabled bool
	// Redact controls whether secrets are redacted from compressed output.
	Redact bool
	// RawBypass when true returns output unchanged (explicit raw escape).
	RawBypass bool
}

func (c PiConfig) threshold() int {
	if c.ThresholdBytes <= 0 {
		return 8000
	}
	return c.ThresholdBytes
}

// errorMarkers are checked case-insensitively; if present, output is preserved lossless.
var errorMarkerRegex = regexp.MustCompile(`(?i)\[FAIL\]|FAILED|ERROR|panic|traceback|exception|fatal`)

var diffMarkerRegex = regexp.MustCompile(`(?m)^(?:diff --git|@@ |--- |\+\+\+ )`)

// secretPatterns redacts common credential shapes before persistence/compression.
var secretPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(api[_-]?key\s*[:=]\s*)(["']?)[^"'\s;]+(["']?)`),
	regexp.MustCompile(`(?i)(secret\s*[:=]\s*)(["']?)[^"'\s;]+(["']?)`),
	regexp.MustCompile(`(?i)(password\s*[:=]\s*)(["']?)[^\s"']+(["']?)`),
	regexp.MustCompile(`(?i)(token\s*[:=]\s*)(["']?)[^\s"']+(["']?)`),
	regexp.MustCompile(`(?i)Bearer\s+[A-Za-z0-9_\-\.]+`),
	regexp.MustCompile(`(?i)ghp_[A-Za-z0-9_]+`),
	regexp.MustCompile(`(?i)gho_[A-Za-z0-9_]+`),
	regexp.MustCompile(`(?i)sk-[A-Za-z0-9\-]+`),
}

var privateSourceMarkers = []string{
	"evidence/delivery/",
	".llm-wiki/raw/",
	".llm-wiki/meta/",
	"broker_data",
	"Linear evidence",
}

// PiOptimize is the deterministic, transform-only Pith API for Pi.
// It never spawns or re-runs the original command.
func PiOptimize(command, output string, exitCode int) (string, error) {
	return PiOptimizeWithConfig(command, output, exitCode, PiConfig{})
}

// PiOptimizeWithConfig is the configurable form.
func PiOptimizeWithConfig(command, output string, exitCode int, cfg PiConfig) (string, error) {
	if cfg.RawBypass {
		return output, nil
	}
	if output == "" {
		return "", nil
	}
	if exitCode != 0 {
		return maybeRedact(output, cfg), nil
	}
	if errorMarkerRegex.MatchString(output) {
		return maybeRedact(output, cfg), nil
	}
	if diffMarkerRegex.MatchString(output) {
		return maybeRedact(output, cfg), nil
	}
	if len(output) < cfg.threshold() {
		return maybeRedact(output, cfg), nil
	}
	compressed := compressLargeOutput(output, cfg.threshold())
	return maybeRedact(compressed, cfg), nil
}

// PiRedact redacts known secret patterns from text. Exported for telemetry-safe persistence.
func PiRedact(s string) string {
	return redactSecrets(s)
}

// PiShouldRedact reports whether the output contains likely secrets or sensitive broker evidence.
func PiShouldRedact(s string) bool {
	low := strings.ToLower(s)
	for _, m := range privateSourceMarkers {
		if strings.Contains(low, strings.ToLower(m)) {
			return true
		}
	}
	redacted := redactSecrets(s)
	return redacted != s
}

func maybeRedact(s string, cfg PiConfig) string {
	if cfg.Redact {
		return redactSecrets(s)
	}
	return s
}

func redactSecrets(s string) string {
	out := s
	for _, re := range secretPatterns {
		if strings.Contains(re.String(), "(") && strings.Contains(re.String(), "[") {
			out = re.ReplaceAllString(out, `${1}[REDACTED]`)
		} else {
			out = re.ReplaceAllString(out, "[REDACTED]")
		}
	}
	out = regexp.MustCompile(`(?i)Bearer\s+[A-Za-z0-9_\-\.]+`).ReplaceAllString(out, "Bearer [REDACTED]")
	return out
}

func compressLargeOutput(output string, threshold int) string {
	lines := strings.Split(output, "\n")
	head := 60
	tail := 60
	if len(lines) > 400 {
		head = 80
		tail = 80
	}
	if len(lines) <= head+tail+1 {
		if len(output) <= threshold {
			return output
		}
		keep := threshold - 80
		if keep < 200 {
			keep = 200
		}
		prefix := output[:keep/2]
		suffix := output[len(output)-keep/2:]
		return prefix + "\n... [middle truncated by Pith PiOptimize]\n" + suffix
	}

	middleStart := head
	middleEnd := len(lines) - tail
	hotKeywords := []string{"warn", "info", "test", "ok", "pass"}
	hotSet := make(map[int]bool)
	for i := middleStart; i < middleEnd; i++ {
		low := strings.ToLower(lines[i])
		for _, k := range hotKeywords {
			if strings.Contains(low, k) {
				hotSet[i] = true
				break
			}
		}
	}
	result := make([]string, 0, head+tail+10)
	result = append(result, lines[:head]...)
	keptHot := 0
	for i := middleStart; i < middleEnd && keptHot < 10; i++ {
		if hotSet[i] {
			start := i - 1
			if start < middleStart {
				start = middleStart
			}
			end := i + 1
			if end >= middleEnd {
				end = middleEnd - 1
			}
			for j := start; j <= end; j++ {
				result = append(result, lines[j])
			}
			keptHot++
		}
	}
	removed := (middleEnd - middleStart) - (len(result) - head)
	if removed > 0 {
		result = append(result, fmt.Sprintf("... [%d lines removed by Pith PiOptimize] ...", removed))
	}
	result = append(result, lines[len(lines)-tail:]...)
	return strings.Join(result, "\n")
}
