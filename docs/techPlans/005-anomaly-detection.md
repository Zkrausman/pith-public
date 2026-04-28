# Tech Plan: Pith Anomaly Detection & Overseer Integration

## Overview
Implement an intelligent anomaly detection system in Pith that identifies hallucinations, low-quality LLM responses, and security leaks by auditing Overseer's telemetry logs.

## Strategic Goal
Move Pith from a basic token compressor to a proactive quality advisor that ensures the reliability of estate-wide LLM interactions.

## Architecture
1. **Overseer Bridge**: Connect to Overseer's partitioned JSONL logs via DuckDB.
2. **Heuristic Engine**: 
   - **Hallucination Detection**: Identify patterns of high-temperature randomness or repetitive loops.
   - **Quality Guard**: Flag "lazy" or "unhelpful" AI canned responses.
   - **Security Audit**: Scan for leaked credentials or sensitive filesystem paths in responses.
3. **Semantic Judge**: Use a high-fidelity model (Gemini 3 Pro) to perform a secondary "cold" audit on flagged interactions.

## CLI Interface
`pith audit anomalies [--lookback <duration>] [--format json|text]`

## Integration
- Automatically create **Beads** tasks for critical anomalies detected.
- Signal **Squire** to alert the user during the next `brief` execution.

## Verification
- Unit tests for heuristic detection logic.
- Mock integration tests querying local JSONL logs.
