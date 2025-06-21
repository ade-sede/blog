## Guidelines

- Build after every change to sources.
- Never comment code or configuration.

## Build, Lint, and Test

- **Build:** `make all`
- **Format:** `make format`
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
- **Formatting:** Use `gofmt` and `templ fmt`. Run `make format`.
- **Naming Conventions:** Standard Go conventions.
- **Error Handling:** Use `if err != nil`.
- **Imports:** Organized by Go toolchain.
- **Dependencies:** Go modules (`go.mod`, `go.sum`).
- **Structure:** Logic in `src/`, articles in `articles/`, notes in `quick-notes/`, output in `web/`.
- Never comment code or configuration
