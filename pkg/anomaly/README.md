# anomaly package

The `anomaly` package provides logic for detecting behavioral and economic anomalies in LLM interactions. It bridges Pith with Overseer's telemetry stream to identify hallucinations and low-quality responses.

## Architecture

```mermaid
graph TD
    AuditID[pkg/anomaly/anomaly.go:AuditOverseerLogs] --> DuckDBID[DuckDB: read_json_auto]
    AuditID --> HeuristicsID[Heuristic Engine]
    
    HeuristicsID --> LazyID[Lazy Marker Check]
    HeuristicsID --> TokenID[Unusual Token Consumption]
```
