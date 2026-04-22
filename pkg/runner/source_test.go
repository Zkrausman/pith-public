package runner

import (
	"testing"
)

func TestDetectSourceGemini(t *testing.T) {
	t.Setenv("GEMINI_CLI", "true")
	t.Setenv("GOOGLE_API_KEY", "")
	t.Setenv("CLAUDE_CODE", "")
	t.Setenv("ANTHROPIC_API_KEY", "")
	if DetectSource() != "gemini" {
		t.Errorf("Expected gemini, got %s", DetectSource())
	}
}

func TestDetectSourceClaude(t *testing.T) {
	t.Setenv("GEMINI_CLI", "")
	t.Setenv("GOOGLE_API_KEY", "")
	t.Setenv("CLAUDE_CODE", "true")
	t.Setenv("ANTHROPIC_API_KEY", "")
	if DetectSource() != "claude" {
		t.Errorf("Expected claude, got %s", DetectSource())
	}
}

func TestDetectSourceUnknown(t *testing.T) {
	t.Setenv("GEMINI_CLI", "")
	t.Setenv("GOOGLE_API_KEY", "")
	t.Setenv("CLAUDE_CODE", "")
	t.Setenv("ANTHROPIC_API_KEY", "")
	if DetectSource() != "unknown" {
		t.Errorf("Expected unknown, got %s", DetectSource())
	}
}
