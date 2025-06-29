## First-time setup

- Run `make init` to setup pre-commit hooks

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
- **Test:** When making cosmetic changes test the result by taking a screenshot
- **Screenshot:** Store in `web/` directory. Base command: `chromium-browser --headless --disable-gpu --virtual-time-budget=5000 --screenshot=filename.png web/path/to/file`
  - Desktop: `--window-size=1920,1080`
  - Mobile: `--window-size=414,896`
  - Long pages: `--window-size=1920,3000` or `--window-size=1920,5000`
  - Scroll to position: Add `--evaluate-on-load="setTimeout(() => window.scrollTo(0, PIXELS), 3000)"`
- **Full page PDF:** `chromium-browser --headless --disable-gpu --print-to-pdf=output.pdf web/path/to/file`

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
