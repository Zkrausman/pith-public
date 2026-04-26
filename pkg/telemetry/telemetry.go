package telemetry

import (
	"database/sql"
	"fmt"
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
	DB *sql.DB
}

func NewTelemetry(storagePath string) (*Telemetry, error) {
	if storagePath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		storagePath = filepath.Join(home, ".pith")
	}

	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(storagePath, "pith.db")
	return NewTelemetryWithPath(dbPath)
}

func NewTelemetryWithPath(dbPath string) (*Telemetry, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	t := &Telemetry{DB: db}
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

	_, err := t.DB.Exec(query)
	if err != nil {
		return err
	}

	// Migrations
	_, _ = t.DB.Exec("ALTER TABLE executions ADD COLUMN original_content TEXT")
	_, _ = t.DB.Exec("ALTER TABLE executions ADD COLUMN compressed_content TEXT")
	_, _ = t.DB.Exec("ALTER TABLE executions ADD COLUMN source TEXT DEFAULT 'unknown'")

	return nil
}

func (t *Telemetry) Record(rec ExecutionRecord) error {
	query := `INSERT INTO executions (command, original_tokens, compressed_tokens, original_content, compressed_content, duration_ms, parser_used, is_passthrough, source)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := t.DB.Exec(query, rec.Command, rec.OriginalTokens, rec.CompressedTokens, rec.OriginalContent, rec.CompressedContent, rec.DurationMs, rec.ParserUsed, rec.IsPassthrough, rec.Source)
	return err
}

func (t *Telemetry) GetStatsByDay(source string) ([]struct {
	Date       string
	Original   int
	Compressed int
}, error) {
	whereClause := ""
	var args []interface{}
	if source != "" && source != "all" {
		whereClause = "WHERE source = ?"
		args = append(args, source)
	}

	query := fmt.Sprintf(`
	SELECT strftime('%%Y-%%m-%%d', timestamp) as date,
	       SUM(original_tokens),
	       SUM(compressed_tokens)
	FROM executions
	%s
	GROUP BY date
	ORDER BY date ASC`, whereClause)

	rows, err := t.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []struct {
		Date       string
		Original   int
		Compressed int
	}

	for rows.Next() {
		var r struct {
			Date       string
			Original   int
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
	return t.DB.Close()
}

func (t *Telemetry) GetUnparsedCommands(source string) ([]struct {
	Pattern         string
	InvocationCount int
	TotalRawTokens  int
	Source          string
}, error) {
	whereClause := "WHERE is_passthrough = 1"
	var args []interface{}
	if source != "" && source != "all" {
		whereClause += " AND source = ?"
		args = append(args, source)
	}

	query := fmt.Sprintf(`
	SELECT command, COUNT(*), SUM(original_tokens), MAX(source)
	FROM executions
	%s
	GROUP BY command
	ORDER BY SUM(original_tokens) DESC
	`, whereClause)

	rows, err := t.DB.Query(query, args...)
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

func (t *Telemetry) GetStats(source string) (totalOriginal, totalCompressed int, err error) {
	whereClause := ""
	var args []interface{}
	if source != "" && source != "all" {
		whereClause = "WHERE source = ?"
		args = append(args, source)
	}

	query := fmt.Sprintf(`SELECT COALESCE(SUM(original_tokens), 0), COALESCE(SUM(compressed_tokens), 0) FROM executions %s`, whereClause)
	err = t.DB.QueryRow(query, args...).Scan(&totalOriginal, &totalCompressed)
	if err == sql.ErrNoRows {
		return 0, 0, nil
	}
	return
}

func (t *Telemetry) GetStatsByCommand(source string) ([]struct {
	Command    string
	Original   int
	Compressed int
	Source     string
}, error) {
	whereClause := ""
	var args []interface{}
	if source != "" && source != "all" {
		whereClause = "WHERE source = ?"
		args = append(args, source)
	}

	query := fmt.Sprintf(`SELECT command, SUM(original_tokens), SUM(compressed_tokens), MAX(source) FROM executions %s GROUP BY command ORDER BY SUM(original_tokens) DESC`, whereClause)

	rows, err := t.DB.Query(query, args...)
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
	_, err := t.DB.Exec("DELETE FROM executions")
	return err
}

func (t *Telemetry) ResetPassthrough() error {
	_, err := t.DB.Exec("DELETE FROM executions WHERE is_passthrough = 1")
	return err
}

func (t *Telemetry) GetRecentExecutions(limit int, source string) ([]ExecutionRecord, error) {
	whereClause := ""
	var args []interface{}
	if source != "" && source != "all" {
		whereClause = "WHERE source = ?"
		args = append(args, source)
	}
	args = append(args, limit)

	query := fmt.Sprintf(`
	SELECT id, timestamp, command, original_tokens, compressed_tokens, duration_ms, parser_used, is_passthrough, source
	FROM executions
	%s
	ORDER BY timestamp DESC
	LIMIT ?`, whereClause)

	rows, err := t.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ExecutionRecord

	for rows.Next() {
		var r ExecutionRecord
		if err := rows.Scan(&r.ID, &r.Timestamp, &r.Command, &r.OriginalTokens, &r.CompressedTokens, &r.DurationMs, &r.ParserUsed, &r.IsPassthrough, &r.Source); err != nil {
			return nil, err
		}
		results = append(results, r)
	}

	return results, nil
}

func (t *Telemetry) GetSources() ([]string, error) {
	query := `SELECT DISTINCT source FROM executions WHERE source != 'unknown' AND source != '' ORDER BY source ASC`
	rows, err := t.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		results = append(results, s)
	}

	return results, nil
}

func (t *Telemetry) GetExecutionDetails(id int64) (*ExecutionRecord, error) {
	query := `
	SELECT id, timestamp, command, original_tokens, compressed_tokens, original_content, compressed_content, duration_ms, parser_used, is_passthrough
	FROM executions
	WHERE id = ?`

	var r ExecutionRecord
	err := t.DB.QueryRow(query, id).Scan(&r.ID, &r.Timestamp, &r.Command, &r.OriginalTokens, &r.CompressedTokens, &r.OriginalContent, &r.CompressedContent, &r.DurationMs, &r.ParserUsed, &r.IsPassthrough)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (t *Telemetry) SearchExecutions(queryStr string, source string, limit int) ([]ExecutionRecord, error) {
	whereClause := "WHERE (command LIKE ? OR original_content LIKE ? OR compressed_content LIKE ?)"
	args := []interface{}{"%" + queryStr + "%", "%" + queryStr + "%", "%" + queryStr + "%"}
	if source != "" && source != "all" {
		whereClause += " AND source = ?"
		args = append(args, source)
	}
	args = append(args, limit)

	sqlQuery := fmt.Sprintf(`
	SELECT id, timestamp, command, original_tokens, compressed_tokens, duration_ms, parser_used, is_passthrough, source
	FROM executions
	%s
	ORDER BY timestamp DESC
	LIMIT ?`, whereClause)

	rows, err := t.DB.Query(sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ExecutionRecord

	for rows.Next() {
		var r ExecutionRecord
		if err := rows.Scan(&r.ID, &r.Timestamp, &r.Command, &r.OriginalTokens, &r.CompressedTokens, &r.DurationMs, &r.ParserUsed, &r.IsPassthrough, &r.Source); err != nil {
			return nil, err
		}
		results = append(results, r)
	}

	return results, nil
}
