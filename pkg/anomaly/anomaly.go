package anomaly

import (
	"database/sql"
	"fmt"
)

type Anomaly struct {
	ExecutionID int64
	Reason      string
	Severity    string // warning, critical
}

// DetectUnusualConsumption flags executions that exceed the average token usage.
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
			ExecutionID: id,
			Reason:      fmt.Sprintf("Token consumption (%.0f) is significantly higher than 24h average (%.0f)", tokens, avg),
			Severity:    "warning",
		})
	}

	return anomalies, nil
}
