## First-time setup

Run `make init` before anything else. It installs pre-commit hooks (gofmt, templ fmt, prettier, etc.) and the post-push git hook. Without this, commits will not be formatted correctly.

## Architecture

This is a custom static site generator written in Go. It was built from scratch rather than using an existing framework (Hugo, Jekyll, etc.) because the markdown pipeline requires features no generic tool provides out of the box: extended code fences with filename headers and diff highlighting, directory-tree blocks, LaTeX rendering, dynamic-color images, byline injection, and a footnote sidebar.

The build is two-phase:

1. The Makefile copies static assets and runs `templ generate` to compile `.templ` files into Go code.
2. `go run src/*.go` reads article manifests and markdown, renders every page via templ components, minifies CSS and JS directly into each HTML file, and writes output to `web/`.

Every output HTML file is fully self-contained: page-specific CSS and JS are minified and embedded as `<style>` and `<script>` blocks in `<head>` at build time. This eliminates extra network fetches per page. The trade-off — no cross-page CSS caching — is acceptable because total asset size is small.

Output is pure static HTML in `web/`. There is no runtime server.

## Build system

The `Makefile` is the source of truth for all workflows. Read it before running anything. It is short and self-explanatory.

The same applies to:

- `src/config.go` for environment variables and their defaults
- `src/` for source code structure and which file does what
- `articles/` for content structure

## Development environment

The project uses [devbox](https://www.jetify.com/devbox) and direnv. Entering the directory activates the devbox shell automatically via `.envrc`, which also exports `ENV=development`. This makes draft articles visible in the build. In any other environment, drafts are excluded.

## Article system

Each article consists of a JSON manifest and a markdown file sharing a slug as basename under `articles/`:

```
articles/
  my-article.json
  my-article.md
  my-article.css   (optional)
  my-article.js    (optional)
```

The manifest controls metadata: title, date, author, description, tags, draft status, and which files to use. Articles with `"draft": true` are excluded from production builds. Optional `.css` and `.js` files are inlined into that article's `<head>` only — they are not shared across pages.

Use `scripts/new-article.py` or the `new-article` skill to scaffold new articles. Do not create manifest files by hand.

### Custom markdown features

The markdown pipeline extends standard GFM. When writing article content, these are available:

- **Extended code fences**: ` ```language:filename:diff ` — any combination. Renders a filename header with a file-type icon; `diff` enables green/red line backgrounds.
- **Directory-tree blocks**: fenced block with language `directory-structure`. Indented text is rendered as a nested tree with folder/file icons.
- **LaTeX**: inline `$...$` and display `\[...\]` are pre-processed into elements with `data-latex` attributes; KaTeX renders them client-side.
- **Dynamic-color images**: `![alt](src){.dynamic-colors}` — at runtime, the dominant color is extracted and WCAG contrast is checked against the current background; a hue-shift is applied if contrast falls below 7.0.

## Code style

- **Language:** Go
- **Templating:** [templ](https://templ.dev/) — always use templ components for HTML; never generate HTML strings directly in Go
- **Formatting:** `gofmt` and `templ fmt`, run automatically by the pre-commit hook
- **Naming conventions:** standard Go conventions
- **Error handling:** `if err != nil`
- **Imports:** organized by Go toolchain
- **Dependencies:** Go modules (`go.mod`, `go.sum`) — do not introduce new dependencies without explicit approval

### Asset discipline

Never add new external asset references (`<link href="...">` or `<script src="...">`) for page-specific styles or scripts. All page-specific assets go through the inline pipeline: declare them in the article manifest (`cssFile`, `scriptFile`) or add them to the global asset list in `src/main.go`.

### Mobile-first

Most readers are on mobile. Every layout and styling decision should work on mobile first. Code snippets in articles must be no wider than 35 characters — break lines accordingly.
