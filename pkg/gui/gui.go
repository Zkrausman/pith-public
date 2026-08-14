package gui

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"pith/pkg/config"
	"pith/pkg/telemetry"
	"runtime"
	"strings"
	"time"
)

//go:embed dashboard.html static/*
var staticFiles embed.FS

func StartDashboard(cfg *config.Config, tel *telemetry.Telemetry, port int) error {
	registerHandlers(cfg, tel)

	// The dashboard exposes telemetry; never bind it to the network by default.
	url := fmt.Sprintf("http://127.0.0.1:%d", port)
	fmt.Printf("Starting Pith Dashboard at %s (loopback only)\n", url)

	// Open browser in background
	go openBrowser(url)

	return http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), nil)
}

func registerHandlers(cfg *config.Config, tel *telemetry.Telemetry) {
	http.Handle("/static/", http.FileServer(http.FS(staticFiles)))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := staticFiles.ReadFile("dashboard.html")
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write(data)
	})

	http.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		source := r.URL.Query().Get("source")
		totalOrig, totalComp, _ := tel.GetStats(source)
		byCmd, _ := tel.GetStatsByCommand(source)
		daily, _ := tel.GetStatsByDay(source)

		resp := map[string]interface{}{
			"total_original":         totalOrig,
			"total_compressed":       totalComp,
			"by_command":             byCmd,
			"daily":                  daily,
			"usd_per_million_tokens": cfg.USDPerMillionTokens,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	http.HandleFunc("/api/discover", func(w http.ResponseWriter, r *http.Request) {
		source := r.URL.Query().Get("source")
		unparsed, _ := tel.GetUnparsedCommands(source)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(unparsed)
	})

	http.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		source := r.URL.Query().Get("source")
		if query == "" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]interface{}{})
			return
		}

		results, err := tel.SearchExecutions(query, source, 50)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	})

	http.HandleFunc("/api/telemetry/push", func(w http.ResponseWriter, r *http.Request) {
		// Dashboard service-mode mutation is intentionally disabled. Use the
		// explicit CLI import command for local files instead.
		http.Error(w, "telemetry import is unavailable from the dashboard", http.StatusForbidden)
		return
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		err := tel.ImportJSONL(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Import failed: %v", err), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status":"ok"}`)
	})

	http.HandleFunc("/api/telemetry/pull", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "telemetry export is unavailable from the dashboard", http.StatusForbidden)
		return
		sinceIDStr := r.URL.Query().Get("since_id")
		var sinceID int64
		if sinceIDStr != "" {
			fmt.Sscanf(sinceIDStr, "%d", &sinceID)
		}
		w.Header().Set("Content-Type", "application/x-jsonlines")
		err := tel.ExportJSONLSince(w, sinceID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/api/telemetry/sync", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "telemetry sync is unavailable from the dashboard", http.StatusForbidden)
		return
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		serverURL := cfg.SyncServerURL
		if serverURL == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, `{"status":"error","error":"Sync server URL not configured"}`)
			return
		}

		var pushBuf strings.Builder
		if err := tel.ExportJSONL(&pushBuf); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"status":"error","error":"Failed to export local telemetry: %v"}\n`, err)
			return
		}

		importEndpoint := strings.TrimSuffix(serverURL, "/") + "/api/telemetry/push"
		importReq, err := http.NewRequest("POST", importEndpoint, strings.NewReader(pushBuf.String()))
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"status":"error","error":"Failed to prepare push request: %v"}\n`, err)
			return
		}
		importReq.Header.Set("Content-Type", "application/x-jsonlines")

		client := &http.Client{Timeout: 10 * time.Second}
		importResp, err := client.Do(importReq)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"status":"error","error":"Push server connection failed: %v"}\n`, err)
			return
		}
		defer importResp.Body.Close()

		if importResp.StatusCode != http.StatusOK {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"status":"error","error":"Push failed with status %d"}\n`, importResp.StatusCode)
			return
		}

		pullEndpoint := strings.TrimSuffix(serverURL, "/") + "/api/telemetry/pull?since_id=0"
		pullReq, err := http.NewRequest("GET", pullEndpoint, nil)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"status":"error","error":"Failed to prepare pull request: %v"}\n`, err)
			return
		}

		pullResp, err := client.Do(pullReq)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"status":"error","error":"Pull server connection failed: %v"}\n`, err)
			return
		}
		defer pullResp.Body.Close()

		if pullResp.StatusCode != http.StatusOK {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"status":"error","error":"Pull failed with status %d"}\n`, pullResp.StatusCode)
			return
		}

		if err := tel.ImportJSONL(pullResp.Body); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"status":"error","error":"Failed to import telemetry: %v"}\n`, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status":"ok"}`)
	})

	http.HandleFunc("/api/recent", func(w http.ResponseWriter, r *http.Request) {
		source := r.URL.Query().Get("source")
		recent, _ := tel.GetRecentExecutions(50, source)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(recent)
	})

	http.HandleFunc("/api/execution", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		var id int64
		fmt.Sscanf(idStr, "%d", &id)

		execution, err := tel.GetExecutionDetails(id)
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(execution)
	})

	http.HandleFunc("/api/sources", func(w http.ResponseWriter, r *http.Request) {
		sources, _ := tel.GetSources()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sources)
	})
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		fmt.Printf("Failed to open browser: %v\n", err)
	}
}
