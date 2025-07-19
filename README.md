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

## General working principle

All static assets to be served are in the `web/` directory.

- CSS files
- Images & Icons
- Fonts

HTML files are assembled from templates.
CSS & JS is injected in-line wherever possible.

Everything static, efficiently cached.
Usable on mobile, even thought it isn't glorious.

## Deploying to Cloudflare pages

The GUI on the Cloudflare console is pretty self explanatory.
⚠️ Cloudflare Pages runners do not set the `GOPATH` variable by default, don't forget to set it in the pages settings.
