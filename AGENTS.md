## First-time setup

- Run `make init` to setup pre-commit hooks

## Guidelines

- Build after every change to sources.
- Never comment code or configuration.
- After you make changes, update the documentation if applicable.
- Consider every choice for mobile first - most people read me on mobile
- Considering we target mobile first, code snippets shouldn't be larger than 35 chars. Format accordingly

## Command cheat sheet

- **Build:** `make all`
- **Format:** Pre-commit hook
- **Clean:** `make clean`
- **Clean + Build in one command:** `make re`
- **PDF:** `make pdf`
- **Deploy:** `make deploy`
- **Test:** When making cosmetic changes test the result by taking a screenshot. Don't forget you can also inspect the HTML artficats
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
- **Structure:** Logic in `src/`, articles in `articles/`, output in `web/`.
