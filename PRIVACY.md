# Privacy and local diagnostic logs

Pith stores telemetry in its configured local storage directory. Diagnostic command/output logs are **disabled by default**.

To opt in, set `snag_logging: true` in Pith's `config.json`. Pith redacts common API-key, token, secret, password, and bearer-token values before writing, but redaction is heuristic and cannot guarantee removal of every sensitive value. Logs are owner-only, retain at most 1 MiB before one rotated prior file is kept, and can be disabled again by setting `snag_logging` to `false`.

`pith reset --all` removes telemetry and both diagnostic log files. Use a storage directory with appropriate local-disk protections.

## Network behavior

Update checks use anonymous requests to the public GitHub Releases API and never read `GITHUB_TOKEN` or invoke `gh auth`. `pith audit anomalies --diagnose` sends bounded, redacted diagnostic snippets to Gemini only when you also pass `--allow-external-ai`; without that flag no data is sent. The local dashboard bundles its JavaScript, CSS, and fonts and makes no browser requests to third-party CDNs.
