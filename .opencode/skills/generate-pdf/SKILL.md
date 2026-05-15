---
name: generate-pdf
description: Use when the user wants to generate or update the resume PDF.
---

# Generate Resume PDF

Your job is to produce an up-to-date `ade-sede.pdf` resume, either locally or as a published GitHub release asset. Do the mechanical work and give the user a direct link or path when done — never leave them to go find it themselves.

## Workflow

### 1. Decide where to generate

Choose based on context. Ask only if you cannot infer the answer.

**Generate locally** if:

- The user is working in a local dev environment
- They just edited the resume source and want to preview the result
- They have not pushed recent changes to `main`

**Trigger GitHub Actions** if:

- The user wants to publish an updated PDF available for download
- They have just pushed changes to `main` and want the release asset updated
- They explicitly mention GitHub, the download link, or the release

If you cannot infer either way, ask once:

> Do you want to generate the PDF locally (`make pdf`) or publish a new version to GitHub so the download link is updated?

### 2. Local generation

Run:

```
make pdf
```

This builds the full site, then runs `scripts/generate-pdf.py` to render `web/resume-printable.html` to PDF via weasyprint.

Output: `src/pdfs/ade-sede.pdf`

Tell the user:

> PDF generated at `src/pdfs/ade-sede.pdf`.

### 3. GitHub Actions generation

Run:

```
gh workflow run pdf.yml --ref main
```

This triggers the `Generate resume PDF` workflow. When it completes, the PDF is uploaded as a release asset at the permanent URL:

`https://github.com/ade-sede/blog/releases/latest/download/ade-sede.pdf`

Tell the user:

> Workflow triggered. Once it completes, the updated PDF will be available at:
> https://github.com/ade-sede/blog/releases/latest/download/ade-sede.pdf

If you want to confirm the run was queued, you can follow up with:

```
gh run list --workflow=pdf.yml --limit=1
```
