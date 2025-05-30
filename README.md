# Blog

## Installation

### Managing dependencies yourself

- [Golang 1.20](https://go.dev/doc/install) or greater
- Make sure [`GOPATH`](https://go.dev/wiki/GOPATH) environment variable is properly set
- [GNU Make](https://www.gnu.org/software/make/)
- [Chromium]https://www.chromium.org/getting-involved/download-chromium/) (for automated PDF resume generation)


### Managing dependencies through devbox

- [Devbox](https://www.jetify.com/devbox)

Devbox install dependencies through nix.
After installing devbox, you can start a shell with all the registered dependencies available in your path:

```bash
devbox shell
```

### After installing dev dependencies

```bash
make init
```

It will:
- Install [`templ`](https://github.com/a-h/templ) to `$GOPATH/bin`.  
- Install pre-commit hooks

## Running

If you are managing dependencies through devbox, don't forget to step into a devbox managed shell using `devbox shell`.

```bash
# Generate HTML files
make

# Remove all artifacts (*_templ.go, *.html)
make clean

# Remove all artifacts and re-generate files
make re

# Self explanatory
make format

# Generates HTML files and a PDF version of the resume (available at src/pdfs/ade-sede.pdf)
make pdf
```

## General working principle

All static assets to be served are in the `web/` directory.

- CSS files
- Images & Icons
- Fonts

HTML files are assembled from templates.  
1 CSS file, 1 JS file.

Everything static, very efficiently cached.

## Deployment

Currently deployed to Cloudfare Pages: [blog.ade-sede.dev](https://blog.ade-sede.dev)  
⚠️ Cloudflare Pages runners do not set the `GOPATH` variable by default, don't forget to set it in the pages settings.

## Disclaimer

Lot of AI written stuff in there.
