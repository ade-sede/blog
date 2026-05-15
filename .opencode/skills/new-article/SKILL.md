---
name: new-article
description: Use when the user wants to create a new blog article, start writing a post, or scaffold an article manifest. Guides the human through title, slug, description, and tag selection, then runs scripts/new-article.py to generate the files.
---

# New Article

Your job is to help the human arrive at a well-formed article manifest and kick off the file creation. You are a collaborator, not a form — engage with the ideas, push back if something is unclear, and do the mechanical work so the human can focus on thinking.

## Workflow

Follow these steps in order. Do not jump ahead.

### 1. Understand the topic

If the human has not already described what the article is about, ask before anything else. You need enough context to help with the title and description. One or two sentences from them is enough to proceed.

### 2. Workshop the title

A good title is specific and honest — it says exactly what the article does or argues, without being clickbait.

- If they already have a title, confirm it and move on.
- If they don't, propose 2–3 options that vary in tone (direct, essay-like, question-form). Keep them short.
- Iterate until they confirm one. Do not proceed without a confirmed title.

### 3. Derive and confirm the slug

Derive the slug from the confirmed title:

- Lowercase
- Spaces replaced with hyphens
- Remove all characters that are not `a-z`, `0-9`, or `-`
- No leading or trailing hyphens

Show the slug to the human before using it. Example:

> Slug: `the-tyranny-of-simple` — does that work, or would you like to change it?

If they want to override it, validate that their version is also lowercase-hyphenated. The script will enforce this too, but catch it early.

### 4. Draft the description

The description appears in article cards and the `<meta>` tag. It should be 1–2 sentences that give a reader enough to decide whether to click.

- Suggest a description based on what the human told you about the topic and title.
- They confirm, edit, or ask for alternatives.
- Do not proceed without a confirmed description.

### 5. Select tags

Valid tags are: `essay`, `quick note`

- **essay** — a long-form piece with an argument or narrative arc
- **quick note** — a short, focused reference or how-to, mostly for personal documentation

Ask which applies. If the article sounds like it could be either, say so and let the human decide. Multiple tags are allowed.

### 6. Run the script

Once slug, title, description, and tags are all confirmed, run the script in CLI mode:

```
python3 scripts/new-article.py \
  --slug <slug> \
  --title "<title>" \
  --description "<description>" \
  --tags <tag1> [<tag2> ...]
```

Example:

```
python3 scripts/new-article.py \
  --slug the-tyranny-of-simple \
  --title "The tyranny of simple" \
  --description "Simplicity is a mantra of the software industry, yet few organisations define what it actually means." \
  --tags essay
```

### 7. Confirm and close

Show the human the contents of the created manifest file. Then remind them:

- The article is created as a **draft** (`"draft": true`). It will not appear in production builds until that is set to `false`.
- To preview locally: `ENV=development make re`, then open `web/<slug>.html`.
- To publish: set `"draft": false` in `articles/<slug>.json` and rebuild.
