package telemetry

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type ExecutionRecord struct {
	ID                int64
	Timestamp         time.Time
	Command           string
	OriginalTokens    int
	CompressedTokens  int
	DurationMs        int64
	ParserUsed        string
	IsPassthrough     bool
}

type Telemetry struct {
	db *sql.DB
}

func NewTelemetry() (*Telemetry, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(home, ".diet")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	dbPath := filepath.Join(dir, "diet.db")
	return NewTelemetryWithPath(dbPath)
}

func NewTelemetryWithPath(dbPath string) (*Telemetry, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	t := &Telemetry{db: db}
	if err := t.init(); err != nil {
		return nil, err
	}
	return t, nil
}

func (t *Telemetry) init() error {
	query := `
	CREATE TABLE IF NOT EXISTS executions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		command TEXT,
		original_tokens INTEGER,
		compressed_tokens INTEGER,
		duration_ms INTEGER,
		parser_used TEXT,
		is_passthrough BOOLEAN
	);`
	_, err := t.db.Exec(query)
	return err
}

func (t *Telemetry) Record(rec ExecutionRecord) error {
	query := `INSERT INTO executions (command, original_tokens, compressed_tokens, duration_ms, parser_used, is_passthrough)
	          VALUES (?, ?, ?, ?, ?, ?)`
	_, err := t.db.Exec(query, rec.Command, rec.OriginalTokens, rec.CompressedTokens, rec.DurationMs, rec.ParserUsed, rec.IsPassthrough)
	return err
}

func (t *Telemetry) Close() error {
	return t.db.Close()
}

func (t *Telemetry) GetUnparsedCommands() ([]struct {
	Pattern         string
	InvocationCount int
	TotalRawTokens  int
}, error) {
	query := `
		SELECT command, COUNT(*), SUM(original_tokens)
		FROM executions
		WHERE is_passthrough = 1
		GROUP BY command
		ORDER BY SUM(original_tokens) DESC
	`
	rows, err := t.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []struct {
		Pattern         string
		InvocationCount int
		TotalRawTokens  int
	}
	for rows.Next() {
		var r struct {
			Pattern         string
			InvocationCount int
			TotalRawTokens  int
		}
		if err := rows.Scan(&r.Pattern, &r.InvocationCount, &r.TotalRawTokens); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}

func (t *Telemetry) GetStats() (totalOriginal, totalCompressed int, err error) {
	query := `SELECT COALESCE(SUM(original_tokens), 0), COALESCE(SUM(compressed_tokens), 0) FROM executions`
	err = t.db.QueryRow(query).Scan(&totalOriginal, &totalCompressed)
	if err == sql.ErrNoRows {
		return 0, 0, nil
	}
	return
}

func (t *Telemetry) GetStatsByCommand() ([]struct {
	Command    string
	Original   int
	Compressed int
}, error) {
	query := `SELECT command, SUM(original_tokens), SUM(compressed_tokens) FROM executions GROUP BY command ORDER BY SUM(original_tokens) DESC`
	rows, err := t.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []struct {
		Command    string
		Original   int
		Compressed int
	}
	for rows.Next() {
		var r struct {
			Command    string
			Original   int
			Compressed int
		}
		if err := rows.Scan(&r.Command, &r.Original, &r.Compressed); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}

func (t *Telemetry) ResetAll() error {
	_, err := t.db.Exec("DELETE FROM executions")
	return err
}

func (t *Telemetry) ResetPassthrough() error {
	_, err := t.db.Exec("DELETE FROM executions WHERE is_passthrough = 1")
	return err
}
