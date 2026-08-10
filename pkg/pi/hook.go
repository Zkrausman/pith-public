package pi

import (
	"strings"

	"pith/pkg/parser"
)

// HookRequest is the structured, transform-only contract used by Pi hooks.
// Output is supplied after the original command completed; Pith never executes it.
type HookRequest struct {
	Command   string `json:"command"`
	Output    string `json:"output"`
	ExitCode  int    `json:"exitCode"`
	RawBypass bool   `json:"rawBypass"`
}

type HookResponse struct {
	Output      string `json:"output"`
	Parser      string `json:"parser"`
	Passthrough bool   `json:"passthrough"`
}

// OptimizeHook applies Pith's existing parser registry to successful Pi results.
// Failures, diffs, and raw requests remain lossless except mandatory redaction.
func OptimizeHook(req HookRequest) HookResponse {
	cfg := PiConfig{Redact: true, RawBypass: req.RawBypass, Harness: HarnessPi}
	if req.RawBypass || req.ExitCode != 0 || errorMarkerRegex.MatchString(req.Output) || diffMarkerRegex.MatchString(req.Output) {
		return HookResponse{Output: maybeRedact(req.Output, cfg), Passthrough: true}
	}
	parts := strings.Fields(req.Command)
	if len(parts) == 0 || len(req.Output) < cfg.threshold() {
		return HookResponse{Output: maybeRedact(req.Output, cfg), Passthrough: true}
	}
	for _, candidate := range parser.GetAllParsers() {
		if candidate.CanParse(parts[0], parts[1:]) {
			return HookResponse{Output: maybeRedact(candidate.Parse(req.Output), cfg), Parser: candidate.Name()}
		}
	}
	return HookResponse{Output: maybeRedact(req.Output, cfg), Passthrough: true}
}
