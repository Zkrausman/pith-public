package anomaly

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/duckdb/duckdb-go/v2"
)

type Anomaly struct {
	Timestamp   time.Time
	Project     string
	Reason      string
	Severity    string // warning, critical
	Prompt      string
	Response    string
	Model       string
}

// DetectUnusualConsumption flags executions that exceed the average token usage in Pith's local DB.
func DetectUnusualConsumption(db *sql.DB) ([]Anomaly, error) {
	simpleQuery := `
		SELECT id, original_tokens, (SELECT AVG(original_tokens) FROM executions) as avg_t
		FROM executions
		WHERE original_tokens > (SELECT AVG(original_tokens) * 2 FROM executions)
		AND timestamp > datetime('now', '-24 hours')
		ORDER BY timestamp DESC;
	`

	rows, err := db.Query(simpleQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var anomalies []Anomaly
	for rows.Next() {
		var id int64
		var tokens, avg float64
		if err := rows.Scan(&id, &tokens, &avg); err != nil {
			return nil, err
		}
		anomalies = append(anomalies, Anomaly{
			Reason:   fmt.Sprintf("Token consumption (%.0f) is significantly higher than 24h average (%.0f)", tokens, avg),
			Severity: "warning",
		})
	}

	return anomalies, nil
}

// AuditOverseerLogs scans Overseer's DuckDB for behavioral anomalies.
func AuditOverseerLogs(lookback time.Duration) ([]Anomaly, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	vaultPath := filepath.Join(home, ".overseer", "vault", "**", "*.jsonl")
	// Use / for DuckDB
	vaultPath = strings.ReplaceAll(vaultPath, "\\", "/")

	db, err := sql.Open("duckdb", "")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	viewQuery := fmt.Sprintf("CREATE VIEW logs AS SELECT * FROM read_json_auto('%s', hive_partitioning := true)", vaultPath)
	if _, err := db.Exec(viewQuery); err != nil {
		return nil, fmt.Errorf("failed to load overseer vault: %w", err)
	}

	// 1. Heuristic: Hallucination/Lazy markers
	lazyMarkers := []string{"As an AI model", "I am sorry", "I cannot fulfill", "internal error"}
	var markersQuery strings.Builder
	markersQuery.WriteString("SELECT timestamp, service_name, model_name, prompt, response FROM logs WHERE ")
	for i, m := range lazyMarkers {
		if i > 0 { markersQuery.WriteString(" OR ") }
		markersQuery.WriteString(fmt.Sprintf("response LIKE '%%%s%%'", m))
	}
	markersQuery.WriteString(fmt.Sprintf(" AND timestamp::TIMESTAMP > now() - interval '%d minutes'", int(lookback.Minutes())))

	rows, err := db.Query(markersQuery.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var anomalies []Anomaly
	for rows.Next() {
		var a Anomaly
		var tsStr string
		if err := rows.Scan(&tsStr, &a.Project, &a.Model, &a.Prompt, &a.Response); err != nil {
			continue
		}
		a.Timestamp, _ = time.Parse(time.RFC3339, tsStr)
		a.Reason = "Lazy or unhelpful response marker detected."
		a.Severity = "warning"
		anomalies = append(anomalies, a)
	}

	return anomalies, nil
}
