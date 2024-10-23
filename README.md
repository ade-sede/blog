# Blog

## Installation

- [Golang 1.20](https://go.dev/doc/install) or greater
- Make sure [`GOPATH`](https://go.dev/wiki/GOPATH) environment variable is properly set
- [GNU Make](https://www.gnu.org/software/make/)
- [wkhtmltopdf](https://wkhtmltopdf.org) for PDF generation

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

## PDF Resume

There is an HTML & CSS resume on the about page.  
It's also very handy to be able to download a resume so the following setup exists:

- a 'raw' HTML resume, without any page layout is output to `web/resume-light.html` when running `make`
- there is a little bit of javascript running on that page to make sure the format is print friendly
- we can then run `make generate-pdf` to generate a PDF from this html page using `wkhtmltopdf`
- ⚠️ this is not running on the cloudflare CI because don't want to figure out how to run `wkhtmltopdf`
