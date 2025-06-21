## First-time setup

- Run `make init` to install all dependencies and pre-commit hooks.

## Guidelines

- Build after every change to sources.
- Never comment code or configuration.
- After you make changes, update the documentation if applicable.

## Build, Lint, and Test

- **Build:** `make all`
- **Format:** Pre-commit hook
- **Serve:** `make serve` (serves on port 8080)
- **Clean:** `make clean`
- **PDF:** `make pdf`
- **Deploy:** `make deploy`
- **Test:** No test suite configured

## Verifying your work

- You can inspect generated files in the `web/` directory

## Code Style

- **Language:** Go
- **Templating:** [templ](https://templ.dev/)
- **Formatting:** Use `gofmt` and `templ fmt`. These are run automatically by a pre-commit hook.
- **Naming Conventions:** Standard Go conventions.
- **Error Handling:** Use `if err != nil`.
- **Imports:** Organized by Go toolchain.
- **Dependencies:** Go modules (`go.mod`, `go.sum`).
- **Structure:** Logic in `src/`, articles in `articles/`, notes in `quick-notes/`, output in `web/`.
- Never comment code or configuration
