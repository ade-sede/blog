# Blog

## Installation

- [Golang 1.20](https://go.dev/doc/install) or greater
- Make sure [`GOPATH`](https://go.dev/wiki/GOPATH) environment variable is properly set
- [GNU Make](https://www.gnu.org/software/make/)

Install [`templ`](https://github.com/a-h/templ) to `$GOPATH/bin`.  
Install pre-commit hooks for auto-formating.

```bash
make init
```

## Running

```bash
# Generate HTML files
make

# Remove all artifacts (*_templ.go, *.html)
make clean

# Remove all artifacts and re-generate files
make re

# Self explanatory
make format
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

Currently deployed to Cloudfare Pages: [blog.ade-sede.com](https://blog.ade-sede.com)  
⚠️ Cloudflare Pages runners do not set the `GOPATH` variable by default.

## PDF Resume generation

Automation has been removed until I find a better solution.  
To generate a PDF, open `resume-printable.html` and print to PDF yourself.
