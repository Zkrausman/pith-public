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
	OriginalContent   string
	CompressedContent string
	DurationMs        int64
	ParserUsed        string
	IsPassthrough     bool
	Source            string
}

type Telemetry struct {
	db *sql.DB
}

func NewTelemetry() (*Telemetry, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(home, ".pith")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	dbPath := filepath.Join(dir, "pith.db")
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
		original_content TEXT,
		compressed_content TEXT,
		duration_ms INTEGER,
		parser_used TEXT,
		is_passthrough BOOLEAN,
		source TEXT DEFAULT 'unknown'
	);`
	_, err := t.db.Exec(query)
	if err != nil {
		return err
	}

	// Migration for existing databases
	_, _ = t.db.Exec("ALTER TABLE executions ADD COLUMN original_content TEXT")
	_, _ = t.db.Exec("ALTER TABLE executions ADD COLUMN compressed_content TEXT")
	_, _ = t.db.Exec("ALTER TABLE executions ADD COLUMN source TEXT DEFAULT 'unknown'")

	return nil
}

func (t *Telemetry) Record(rec ExecutionRecord) error {
	query := `INSERT INTO executions (command, original_tokens, compressed_tokens, original_content, compressed_content, duration_ms, parser_used, is_passthrough, source)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := t.db.Exec(query, rec.Command, rec.OriginalTokens, rec.CompressedTokens, rec.OriginalContent, rec.CompressedContent, rec.DurationMs, rec.ParserUsed, rec.IsPassthrough, rec.Source)
	return err
}

func (t *Telemetry) GetStatsByDay() ([]struct {
	Date     string
	Original int
	Compressed int
}, error) {
	query := `
		SELECT strftime('%Y-%m-%d', timestamp) as date, 
		       SUM(original_tokens), 
		       SUM(compressed_tokens) 
		FROM executions 
		GROUP BY date 
		ORDER BY date ASC`
	rows, err := t.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []struct {
		Date     string
		Original int
		Compressed int
	}
	for rows.Next() {
		var r struct {
			Date     string
			Original int
			Compressed int
		}
		if err := rows.Scan(&r.Date, &r.Original, &r.Compressed); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}

func (t *Telemetry) Close() error {
	return t.db.Close()
}

func (t *Telemetry) GetUnparsedCommands() ([]struct {
	Pattern         string
	InvocationCount int
	TotalRawTokens  int
	Source          string
}, error) {
	query := `
		SELECT command, COUNT(*), SUM(original_tokens), MAX(source)
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
		Source          string
	}
	for rows.Next() {
		var r struct {
			Pattern         string
			InvocationCount int
			TotalRawTokens  int
			Source          string
		}
		if err := rows.Scan(&r.Pattern, &r.InvocationCount, &r.TotalRawTokens, &r.Source); err != nil {
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
	Source     string
}, error) {
	query := `SELECT command, SUM(original_tokens), SUM(compressed_tokens), MAX(source) FROM executions GROUP BY command ORDER BY SUM(original_tokens) DESC`
	rows, err := t.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []struct {
		Command    string
		Original   int
		Compressed int
		Source     string
	}
	for rows.Next() {
		var r struct {
			Command    string
			Original   int
			Compressed int
			Source     string
		}
		if err := rows.Scan(&r.Command, &r.Original, &r.Compressed, &r.Source); err != nil {
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

func (t *Telemetry) GetRecentExecutions(limit int) ([]ExecutionRecord, error) {
	query := `
		SELECT id, timestamp, command, original_tokens, compressed_tokens, duration_ms, parser_used, is_passthrough 
		FROM executions 
		ORDER BY timestamp DESC 
		LIMIT ?`
	rows, err := t.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ExecutionRecord
	for rows.Next() {
		var r ExecutionRecord
		if err := rows.Scan(&r.ID, &r.Timestamp, &r.Command, &r.OriginalTokens, &r.CompressedTokens, &r.DurationMs, &r.ParserUsed, &r.IsPassthrough); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}

func (t *Telemetry) GetExecutionDetails(id int64) (*ExecutionRecord, error) {
	query := `
		SELECT id, timestamp, command, original_tokens, compressed_tokens, original_content, compressed_content, duration_ms, parser_used, is_passthrough 
		FROM executions 
		WHERE id = ?`
	var r ExecutionRecord
	err := t.db.QueryRow(query, id).Scan(&r.ID, &r.Timestamp, &r.Command, &r.OriginalTokens, &r.CompressedTokens, &r.OriginalContent, &r.CompressedContent, &r.DurationMs, &r.ParserUsed, &r.IsPassthrough)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
