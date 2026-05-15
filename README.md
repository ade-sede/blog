# Blog

Deployed to [blog.ade-sede.dev](https://blog.ade-sede.dev).

## Installation

You can install the dependencies yourself on your host machine.
Or you can use [Devbox](https://www.jetify.com/devbox) to manage your dependencies.

### Managing dependencies yourself

- [Golang 1.20](https://go.dev/doc/install) or greater
- Make sure [`GOPATH`](https://go.dev/wiki/GOPATH) environment variable is properly set
- [GNU Make](https://www.gnu.org/software/make/)
- [Chromium](https://www.chromium.org/getting-involved/download-chromium/) (used for automated PDF resume generation and screenshots)
- [pre-commit](https://pre-commit.com/)

### Managing dependencies through devbox

- [Devbox](https://www.jetify.com/devbox)

Devbox install dependencies through nix.
To have devbox installed dependencies available in your shell, you have two options:

- Step into a devbox managed shell using `devbox shell`
- Run `direnv allow` so that devbox deps are automatically available as soon as you step into this directory

### After installing dev dependencies

Run `make init` to setup pre-commit hooks.

## Running

If you are managing dependencies through devbox, don't forget to step into a devbox managed shell using `devbox shell`.

```bash
# Generate HTML files
make

# Remove all artifacts (*_templ.go, *.html)
make clean

# Remove all artifacts and re-generate files
make re

# Generates HTML files and a PDF version of the resume (available at src/pdfs/ade-sede.pdf)
make pdf

# Serve files on :8080, useful when working on a remote machine
make serve

# Take a screenshot of a page (desktop)
chromium-browser --headless --window-size=1920,1080 --screenshot=desktop.png web/path/to/file

# Take a screenshot of a page (mobile)
chromium-browser --headless --window-size=414,896 --screenshot=mobile.png web/path/to/file
```

## Architecture

### Build pipeline

`main()` drives the build in five steps, each a direct function call:

```
main()
 ├── LoadConfig()            reads env vars + CLI flags → Config
 ├── publishGlobalCSS()      minifies shared CSS → web/css/
 ├── loadExperiencesFromJSON() reads experiences.json → ExperiencesData
 ├── parseArticles()         reads articles/ → []Article
 └── generateAllPages()      renders every Page → web/*.html
      └── for each Page:
           ├── inlineAssets()   minifies page CSS/JS → style/script tag strings
           └── Page.Render()    calls the templ component → writes HTML file
```

Before `main()` runs, the **Makefile** copies static assets (fonts, images,
third-party libs, scripts) from `src/` to `web/`. The Go program only handles
what needs to be generated or minified at build time.

**Loading articles** (`parseArticles` → `article.go`)

For each `.json` manifest found in `ARTICLE_DIR`:

1. `readArticleManifest` unmarshals the JSON into an `ArticleManifest`.
2. The manifest's `markdownFile` path is handed to `parseArticleMarkdown` (`markdown.go`), which:
   - Pre-processes LaTeX expressions with `processLatexExpressions` so KaTeX can render them client-side.
   - Runs Goldmark with two custom AST transformers: `filenameTitleTransformer` (parses the `language:filename:diff` code fence syntax into node attributes) and `tocExtractor` (collects headings into a `[]TOCEntry`).
   - Renders to HTML with two custom node renderers: `headingRenderer` (adds `id` and anchor links) and `codeBlockRenderer` (Chroma syntax highlighting, diff colouring, directory-tree blocks — the last of which delegates to `directorytree.go`).
   - Injects the author byline before the first `<h1>`.
3. The resulting `Article` struct bundles the manifest, rendered HTML, formatted date, and TOC.
4. All articles are sorted newest-first before being returned.

**Building pages** (`buildPages` → `main.go`)

`buildPages` constructs a `[]Page` — one for each of the four static pages
(home, articles, resume, resume-printable) plus one per article. Each `Page`
declares:

- `Filename` — the output path under `web/`.
- `Assets` — a list of `Asset` values describing which CSS/JS files to inline, and whether to load them from `src/` (global scope) or `ARTICLE_DIR` (article scope).
- `Render` — a closure that already has the page's data (articles, experiences, …) captured, and accepts the resolved `styleTags`/`scriptTags` strings to pass into the templ component.

**Rendering pages** (`generateAllPages` → `main.go`)

For each `Page`, `inlineAssets` reads every declared asset, minifies it via
`minify.go`, and wraps it in a `<style>` or `<script>` tag string. Those
strings are passed to `Page.Render`, which calls the corresponding templ
component. The templ component embeds them verbatim inside `<head>` and writes
the full HTML to disk.

**Post-processing** (`postProcessing` → `main.go`)

Runs after all HTML files exist:

- `generateSitemap` (`sitemap.go`) — walks the article list and emits `sitemap.xml`.
- `generatePDF` (`pdf-generator.go`) — optionally launches headless Chromium, navigates to `resume-printable.html`, and prints it to `src/pdfs/ade-sede.pdf`.

### Source file layout

| File               | Responsibility                                                                                                             |
| ------------------ | -------------------------------------------------------------------------------------------------------------------------- |
| `main.go`          | Orchestration; `Page` and `Asset` types; build pipeline                                                                    |
| `config.go`        | `Config` struct and `LoadConfig`                                                                                           |
| `article.go`       | `Article`, `ArticleManifest`, `TOCEntry` types; manifest reading; article collection parsing                               |
| `markdown.go`      | Goldmark pipeline; custom AST transformers and renderers; LaTeX pre-processing; byline injection; footnote post-processing |
| `directorytree.go` | Directory-tree HTML rendering; diff annotation helpers; file-icon lookup tables                                            |
| `minify.go`        | CSS/JS minification wrappers                                                                                               |
| `experiences.go`   | `ExperienceEntry`, `ExperiencesData` types; JSON loading                                                                   |
| `sitemap.go`       | `sitemap.xml` generation                                                                                                   |
| `pdf-generator.go` | Headless-Chromium PDF rendering of the resume                                                                              |
| `*.templ`          | HTML templates (layout, home, articles, article, resume)                                                                   |

### Inline asset strategy

Every page inlines its CSS and JS directly into the HTML `<head>` as minified `<style>` and `<script>` blocks. Assets are split into two scopes:

- **Global** — shared across all pages, loaded from `src/css/` and `src/scripts/`
- **Article** — specific to one article, loaded from the article directory alongside its markdown file

This means each HTML file is fully self-contained: one HTTP request, zero round-trips for styles or scripts. The trade-off is no cross-page caching of CSS, which is acceptable given the small total asset size.

## Writing a new article

### 1. Create the files

All articles live in `articles/`. Each article needs at minimum two files with
a shared basename:

```
articles/
  my-article.json   ← manifest
  my-article.md     ← content
```

If the article needs its own styles or scripts, add them alongside:

```
articles/
  my-article.json
  my-article.md
  my-article.css    ← optional, referenced from the manifest
  my-article.js     ← optional, referenced from the manifest
```

### 2. Write the manifest

The manifest is a JSON file that describes the article's metadata.

Minimal example:

```json
{
  "title": "My article title",
  "date": "2025-06-01",
  "draft": true,
  "tags": ["essay"],
  "author": "Adrien DE SEDE",
  "authorImage": "picture.webp",
  "description": "One or two sentences shown in article cards and meta tags.",
  "markdownFile": "my-article.md"
}
```

With optional per-article assets:

```json
{
  "title": "My article title",
  "date": "2025-06-01",
  "draft": true,
  "tags": ["essay"],
  "author": "Adrien DE SEDE",
  "authorImage": "picture.webp",
  "description": "One or two sentences shown in article cards and meta tags.",
  "markdownFile": "my-article.md",
  "cssFile": "my-article.css",
  "scriptFile": "my-article.js"
}
```

| Field          | Required | Description                                                             |
| -------------- | -------- | ----------------------------------------------------------------------- |
| `title`        | yes      | Displayed in the page `<title>`, article card, and article header       |
| `date`         | yes      | Publication date in `YYYY-MM-DD` format; used for sorting               |
| `draft`        | yes      | Set to `true` while writing; drafts are excluded from production builds |
| `tags`         | no       | Array of strings shown on article cards                                 |
| `author`       | yes      | Displayed in the byline below the article title                         |
| `authorImage`  | yes      | Filename of an image under `web/images/`, used in the byline            |
| `description`  | yes      | Short summary; used in article cards and the `<meta>` description       |
| `markdownFile` | yes      | Filename of the markdown file, relative to `articles/`                  |
| `cssFile`      | no       | Article-specific CSS file, inlined into the page `<head>`               |
| `scriptFile`   | no       | Article-specific JS file, inlined into the page `<head>`                |

### 3. Write the content

The markdown file supports standard CommonMark plus GitHub Flavoured Markdown
extensions. A few custom features are available:

**Code blocks with filenames and diff highlighting**

The language tag accepts up to three colon-separated parts: `language`,
`filename`, and the literal string `diff`:

````
```go:main.go
// shown with a filename header
```

```go:diff
-removed line
+added line
```

```go:main.go:diff
-removed line
+added line
```
````

**Directory tree blocks**

Use the `directory-structure` language tag to render an indented file list as
an interactive tree:

````
```directory-structure
src/
  main.go
  config.go
web/
  index.html
```
````

**LaTeX**

Inline math: `$E = mc^2$`

Display math:

```
\[
  \int_0^1 f(x)\,dx
\]
```

### 4. Preview as a draft

Set `"draft": true` in the manifest, then build with `ENV=development`:

```bash
ENV=development make re
```

Draft articles are included in the build and accessible at
`web/my-article.html`. They are excluded from the sitemap and from production
builds where `ENV` is anything other than `development`.

### 5. Publish

Set `"draft": false` in the manifest and rebuild:

```bash
make re
```

## Deploying to Cloudflare pages

The GUI on the Cloudflare console is pretty self explanatory.
⚠️ Cloudflare Pages runners do not set the `GOPATH` variable by default, don't forget to set it in the pages settings.
