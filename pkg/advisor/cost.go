package advisor

import (
	"fmt"
)

type SessionEstimate struct {
	RemainingTasks int
	EstimatedCost  float64
	Currency       string
}

// EstimateSessionCost calculates the projected cost based on task complexity.
func EstimateSessionCost(taskCount int, avgTokensPerTask int, pricePerMillion float64) SessionEstimate {
	totalTokens := float64(taskCount * avgTokensPerTask)
	cost := (totalTokens / 1000000.0) * pricePerMillion

	return SessionEstimate{
		RemainingTasks: taskCount,
		EstimatedCost:  cost,
		Currency:       "USD",
	}
}

// FormatReport generates a human-readable cost projection.
func (s SessionEstimate) FormatReport() string {
	return fmt.Sprintf("Projected Session Cost: %.4f %s (for %d remaining tasks)", s.EstimatedCost, s.Currency, s.RemainingTasks)
}
