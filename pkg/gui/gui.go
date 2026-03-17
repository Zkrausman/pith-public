package gui

import (
	"diet/pkg/telemetry"
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
		totalOrig, totalComp, _ := tel.GetStats()
		byCmd, _ := tel.GetStatsByCommand()
		daily, _ := tel.GetStatsByDay()

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
		unparsed, _ := tel.GetUnparsedCommands()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(unparsed)
	})

	url := fmt.Sprintf("http://localhost:%d", port)
	fmt.Printf("Starting Diet Dashboard at %s\n", url)
	
	// Open browser in background
	go openBrowser(url)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
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
