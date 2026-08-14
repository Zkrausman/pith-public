# Third-party notices

Pith is distributed under the [Apache License 2.0](LICENSE). Its releases include Go code and, on supported platforms, DuckDB native bindings. The release SBOM (`sbom.cdx.json`) is the authoritative versioned inventory for each binary release.

## Direct dependencies

| Component | Version | License |
| --- | --- | --- |
| github.com/AlecAivazis/survey/v2 | v2.3.7 | MIT |
| github.com/charmbracelet/bubbles | v1.0.0 | MIT |
| github.com/charmbracelet/bubbletea | v1.3.10 | MIT |
| github.com/charmbracelet/lipgloss | v1.1.0 | MIT |
| github.com/duckdb/duckdb-go/v2 and duckdb-go-bindings | v2.10502.0 / v0.10502.0 | MIT |
| github.com/spf13/cobra | v1.10.2 | Apache-2.0 |
| modernc.org/sqlite | v1.48.2 | BSD-3-Clause |

Transitive dependency versions and license evidence are recorded in the CycloneDX SBOM attached to each release. License detection is automated evidence, not a substitute for legal review. Before changing DuckDB bindings or distributing to a new platform, review the associated native-artifact notices and licenses.
