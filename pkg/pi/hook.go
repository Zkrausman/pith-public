package pi

import (
	"math"
	"strings"
	"time"

	"pith/pkg/parser"
	"pith/pkg/runner"
	"pith/pkg/telemetry"
)

// HookRequest is the JSON stdin contract for `pith pi transform`. Pith only
// transforms completed output and never executes Command.
type HookRequest struct {
	Command             string          `json:"command"`
	Output              string          `json:"output"`
	ExitCode            int             `json:"exitCode"`
	RawBypass           bool            `json:"rawBypass"`
	ThresholdBytes      int             `json:"thresholdBytes"`
	TelemetryEnabled    bool            `json:"telemetryEnabled"`
	StoragePath         string          `json:"storagePath"`
	Model               string          `json:"model"`
	InputCostPerMillion *float64        `json:"inputCostPerMillion"`
	EnabledParsers      map[string]bool `json:"-"`
}

type HookResponse struct {
	Output      string `json:"output"`
	Parser      string `json:"parser"`
	Passthrough bool   `json:"passthrough"`
}

// OptimizeHook uses the Pith parser registry for safe successful results.
// Raw requests, errors, and diffs are lossless except mandatory redaction.
func OptimizeHook(req HookRequest) HookResponse {
	started := time.Now()
	cfg := PiConfig{ThresholdBytes: req.ThresholdBytes, Redact: true, RawBypass: req.RawBypass, Harness: HarnessPi}
	result := HookResponse{Output: maybeRedact(req.Output, cfg), Passthrough: true}
	if !req.RawBypass && req.ExitCode == 0 && !errorMarkerRegex.MatchString(req.Output) && !diffMarkerRegex.MatchString(req.Output) && len(req.Output) >= cfg.threshold() {
		parts := strings.Fields(req.Command)
		if len(parts) > 0 {
			for _, candidate := range parser.GetAllParsers() {
				enabled, configured := req.EnabledParsers[candidate.Name()]
				if (!configured || enabled) && candidate.CanParse(parts[0], parts[1:]) {
					result = HookResponse{Output: maybeRedact(candidate.Parse(req.Output), cfg), Parser: candidate.Name()}
					break
				}
			}
		}
	}
	if req.TelemetryEnabled {
		if tel, err := telemetry.NewTelemetry(req.StoragePath); err == nil {
			defer tel.Close()
			cost := req.InputCostPerMillion
			if cost != nil && (*cost < 0 || math.IsNaN(*cost) || math.IsInf(*cost, 0)) {
				cost = nil
			}
			_ = tel.Record(telemetry.ExecutionRecord{Command: piTelemetryCommand, OriginalTokens: runner.EstimateTokensWithHeuristic(req.Output, 4), CompressedTokens: runner.EstimateTokensWithHeuristic(result.Output, 4), DurationMs: time.Since(started).Milliseconds(), ParserUsed: result.Parser, IsPassthrough: result.Passthrough, Harness: HarnessPi, Model: req.Model, InputCostPerMillion: cost})
		}
	}
	return result
}
