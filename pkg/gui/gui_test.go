package gui

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"pith/pkg/config"
	"pith/pkg/telemetry"
	"strings"
	"testing"
)

func TestDashboardHandlers(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "pith-gui-test-*")
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")
	tel, _ := telemetry.NewTelemetryWithPath(dbPath)
	defer tel.Close()

	// Seed some telemetry
	tel.Record(telemetry.ExecutionRecord{
		Command:          "go version",
		OriginalTokens:   10,
		CompressedTokens: 5,
		IsPassthrough:    false,
		Source:           "test",
	})

	cfg := &config.Config{USDPerMillionTokens: 3.0}
	registerHandlers(cfg, tel)

	// Now we can use httptest to hit the DefaultServeMux
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 OK for /, got %d", w.Code)
	}

	req = httptest.NewRequest("GET", "/api/stats", nil)
	w = httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 OK for stats, got %d", w.Code)
	}

	req = httptest.NewRequest("GET", "/api/discover", nil)
	w = httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 OK for discover, got %d", w.Code)
	}

	req = httptest.NewRequest("GET", "/api/recent", nil)
	w = httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 OK for recent, got %d", w.Code)
	}

	req = httptest.NewRequest("GET", "/api/execution?id=1", nil)
	w = httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w, req)
	// It could be 404 if DB is not flushed or ID is wrong, but it shouldn't panic.
	if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
		t.Errorf("Expected 200 or 404 for execution, got %d", w.Code)
	}

	req = httptest.NewRequest("GET", "/api/sources", nil)
	w = httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 OK for sources, got %d", w.Code)
	}

	// Test /api/telemetry/pull
	req = httptest.NewRequest("GET", "/api/telemetry/pull?since_id=0", nil)
	w = httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("Expected forbidden telemetry pull, got %d", w.Code)
	}

	// Dashboard mutation, export, and remote sync must be unavailable without
	// an explicitly designed authenticated service mode.
	req = httptest.NewRequest("POST", "/api/telemetry/push", strings.NewReader("{}"))
	w = httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("Expected forbidden telemetry push, got %d: %s", w.Code, w.Body.String())
	}

	req = httptest.NewRequest("POST", "/api/telemetry/sync", nil)
	w = httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("Expected forbidden telemetry sync, got %d", w.Code)
	}
}
