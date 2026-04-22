package gui

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"pith/pkg/telemetry"
	"testing"
)

func TestApiHandlers(t *testing.T) {
	tmpDir := t.TempDir()
	tel, _ := telemetry.NewTelemetry(tmpDir)
	defer tel.Close()

	// Seed some data
	tel.Record(telemetry.ExecutionRecord{Command: "test", ID: 1})

	// Test /api/stats
	req, _ := http.NewRequest("GET", "/api/stats?source=all", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		source := r.URL.Query().Get("source")
		tel.GetStats(source)
		w.WriteHeader(http.StatusOK)
	})
	handler.ServeHTTP(rr, req)

	// Test /api/execution Not Found
	req2, _ := http.NewRequest("GET", "/api/execution?id=999", nil)
	rr2 := httptest.NewRecorder()
	handler2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		var id int64
		_, _ = fmt.Sscanf(idStr, "%d", &id)
		_, err := tel.GetExecutionDetails(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	})
	handler2.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", rr2.Code)
	}
}
