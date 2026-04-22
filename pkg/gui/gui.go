package gui

import (
	"pith/pkg/telemetry"
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
)

//go:embed dashboard.html
var staticFiles embed.FS

func StartDashboard(tel *telemetry.Telemetry, port int) error {
	registerHandlers(tel)

	url := fmt.Sprintf("http://localhost:%d", port)
	fmt.Printf("Starting Pith Dashboard at %s\n", url)
	
	// Open browser in background
	go openBrowser(url)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func registerHandlers(tel *telemetry.Telemetry) {
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
			"total_original":   totalOrig,
			"total_compressed": totalComp,
			"by_command":       byCmd,
			"daily":            daily,
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

	http.HandleFunc("/api/recent", func(w http.ResponseWriter, r *http.Request) {
		source := r.URL.Query().Get("source")
		recent, _ := tel.GetRecentExecutions(20, source)
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
