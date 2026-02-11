# CLI Framework Best Practices

## Flag Naming
-   Use kebab-case for flag names (e.g., `--max-tokens`, `--include-ext`).
-   Keep names short but descriptive.
-   Avoid single-character flags if clarity is compromised.

## Configuration
-   Prioritize flags over configuration files.
-   Use `pkg/config.LoadConfig` to support standard config paths (`.toolrc`).
-   Allow environment variables for sensitive data (API keys).

## Output
-   Always support structured output (JSON/XML) for machine consumption (Agents/LLMs).
-   Use `pkg/output` writers instead of `fmt.Println` for data output.
-   Write logs and info messages to `stderr` (via `pkg/util.Logger`), keeping `stdout` clean for data.

## Error Handling
-   Return valid exit codes (non-zero on error).
-   Use `pkg/errors` (future) for typed errors.
-   Print user-friendly error messages to `stderr`.

## Logging
-   Control verbosity with a `--verbose` flag.
-   Use `logger.Debug` for troubleshooting info.
-   Use `logger.Info` for general progress updates.
