package anomaly

import (
	"testing"
	"time"
)

func TestAnomalyStructure(t *testing.T) {
	a := Anomaly{
		Timestamp: time.Now(),
		Reason:    "test reason",
		Severity:  "warning",
	}
	if a.Severity != "warning" {
		t.Error("expected warning severity")
	}
}
