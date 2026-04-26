# Pith Roadmap

## Phase 1: Knowledge Aggregation (Complete)
- [x] Implement High Fidelity Parsing for source and tests.
- [x] Establish "Reasoning Over Brevity" standard.
- [x] Implement Snag log parser (snag list --json).
- [x] Implement Go coverage report parser (go tool cover).
- [x] Improve command extraction for PowerShell & quoted paths.

## Phase 2: Token Economics (Complete)
- [x] Stabilize token tracking and cost analytics.
- [x] Expand parser coverage for common developer tools.

## Phase 3: Advanced Optimization (Current)
- [ ] **Semantic Compression:** Move beyond heuristic-based truncation to LLM-driven summarization for ultra-high-value outputs.
- [ ] **Multi-Agent Sync:** Centralize Pith telemetry across multiple machines/users for enterprise-wide token savings.
- [ ] **Deep-Dive Search:** Add full-text search across all recorded command executions in the dashboard.

## Phase 4: Predictive Analytics (New)
- [x] **Anomaly Detection**: Use local DuckDB models to flag unusual token usage or execution patterns.
- [ ] **Context ROI Heatmap**: Visualize which injected files consistently result in "Hits" vs. "Misses" to optimize `whet` injection logic.
- [x] **Predictive Cost Advisor**: Estimate the total session cost based on current roadmap complexity.
