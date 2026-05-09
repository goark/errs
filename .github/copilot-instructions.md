# Copilot Instructions for github.com/goark/errs

## Project Overview
- This repository provides `github.com/goark/errs`, a Go package for structured error wrapping and multi-error handling.
- Main module target is Go 1.20+.
- There is a nested module under `zapobject/` that should stay in sync with top-level changes when relevant.
- If dependency or Go-version related files are changed, keep `go.mod` (root), `zapobject/go.mod`, and `go.work` consistency in mind.

## Coding Guidelines
- Keep public APIs backward compatible unless explicitly requested.
- Follow idiomatic Go naming and error handling patterns.
- Keep code comments in English.
- Avoid introducing new dependencies unless there is a clear benefit.

## Testing and Validation
- Prefer `task` commands for local validation.
- Primary checks:
  - `task test`
  - `task govulncheck`
- For broad validation, run:
  - `task`
- Before proposing a merge, run at least `task test` and `task govulncheck`. For wide-ranging changes, run `task`.

## CI/Workflow Expectations
- Keep GitHub Actions workflows minimal and readable.
- CI should cover lint, tests, and reachable-vulnerability checks.
- Keep workflow action versions reasonably up to date.
- Keep `ci.yml` split into independent jobs for lint/test and govulncheck where possible for faster feedback.
- Use `golang/govulncheck-action` for govulncheck in CI unless there is a strong reason to change.

## Documentation Expectations
- When behavior or developer flow changes, update `README.md` in the same change.
- Keep usage examples executable and aligned with current API.
- Keep README links to sample programs under `sample/` and `zapobject/` aligned with current behavior.

## API and Compatibility Notes
- Do not introduce new usage of deprecated `Cause()`; prefer `errors.Is`/`errors.As` compatible flows and `Unwrap`/`Unwraps`.
- Treat changes to exported symbols, function signatures, and observable error formatting behavior as potentially breaking.
- Error string and JSON formatting are validated by tests; when output behavior changes, update README examples and test expectations in the same change.
