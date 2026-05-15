#!/usr/bin/env python3

import argparse
import datetime
import json
import re
import sys
from pathlib import Path

ARTICLES_DIR = Path(__file__).parent.parent / "articles"

VALID_TAGS = ["essay", "quick note"]
AUTHOR = "Adrien DE SEDE"
AUTHOR_IMAGE = "picture.webp"

SLUG_RE = re.compile(r"^[a-z0-9]+(-[a-z0-9]+)*$")


def validate_slug(slug: str) -> str:
    slug = slug.strip()
    if not slug:
        raise ValueError("Slug cannot be empty.")
    if not SLUG_RE.match(slug):
        raise ValueError(
            f"Invalid slug '{slug}'. "
            "Use only lowercase letters, digits, and hyphens. "
            "No leading or trailing hyphens."
        )
    return slug


def validate_title(title: str) -> str:
    title = title.strip()
    if not title:
        raise ValueError("Title cannot be empty.")
    return title


def validate_description(description: str) -> str:
    description = description.strip()
    if not description:
        raise ValueError("Description cannot be empty.")
    return description


def validate_tags(raw: list[str]) -> list[str]:
    tags = [t.strip() for t in raw if t.strip()]
    if not tags:
        raise ValueError("At least one tag is required.")
    invalid = [t for t in tags if t not in VALID_TAGS]
    if invalid:
        raise ValueError(
            f"Unknown tag(s): {invalid}. "
            f"Valid tags are: {VALID_TAGS}"
        )
    return tags


def check_no_existing_files(slug: str) -> None:
    json_path = ARTICLES_DIR / f"{slug}.json"
    md_path = ARTICLES_DIR / f"{slug}.md"
    conflicts = [p for p in (json_path, md_path) if p.exists()]
    if conflicts:
        raise FileExistsError(
            f"File(s) already exist: {[str(p) for p in conflicts]}"
        )


def build_manifest(slug: str, title: str, description: str, tags: list[str]) -> dict:
    return {
        "title": title,
        "date": datetime.date.today().isoformat(),
        "draft": True,
        "tags": tags,
        "author": AUTHOR,
        "authorImage": AUTHOR_IMAGE,
        "description": description,
        "markdownFile": f"{slug}.md",
    }


def write_files(slug: str, manifest: dict) -> None:
    json_path = ARTICLES_DIR / f"{slug}.json"
    md_path = ARTICLES_DIR / f"{slug}.md"

    json_path.write_text(json.dumps(manifest, indent=2) + "\n")
    md_path.touch()

    print(f"Created: {json_path}")
    print(f"Created: {md_path}")


def prompt(label: str, hint: str = "") -> str:
    suffix = f" ({hint})" if hint else ""
    return input(f"{label}{suffix}: ").strip()


def interactive_mode() -> tuple[str, str, str, list[str]]:
    print("Creating a new article. Press Ctrl+C to abort.\n")

    while True:
        slug = prompt("Slug", "lowercase-hyphenated")
        try:
            slug = validate_slug(slug)
            check_no_existing_files(slug)
            break
        except (ValueError, FileExistsError) as e:
            print(f"Error: {e}")

    while True:
        title = prompt("Title")
        try:
            title = validate_title(title)
            break
        except ValueError as e:
            print(f"Error: {e}")

    while True:
        description = prompt("Description")
        try:
            description = validate_description(description)
            break
        except ValueError as e:
            print(f"Error: {e}")

    print(f"Valid tags: {', '.join(VALID_TAGS)}")
    while True:
        raw = prompt("Tags", "comma-separated")
        try:
            tags = validate_tags(raw.split(","))
            break
        except ValueError as e:
            print(f"Error: {e}")

    return slug, title, description, tags


def cli_mode(args: argparse.Namespace) -> tuple[str, str, str, list[str]]:
    errors = []

    try:
        slug = validate_slug(args.slug)
    except ValueError as e:
        errors.append(str(e))
        slug = ""

    try:
        title = validate_title(args.title)
    except ValueError as e:
        errors.append(str(e))
        title = ""

    try:
        description = validate_description(args.description)
    except ValueError as e:
        errors.append(str(e))
        description = ""

    try:
        tags = validate_tags(args.tags)
    except ValueError as e:
        errors.append(str(e))
        tags = []

    if errors:
        for err in errors:
            print(f"Error: {err}", file=sys.stderr)
        sys.exit(1)

    try:
        check_no_existing_files(slug)
    except FileExistsError as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)

    return slug, title, description, tags


def main() -> None:
    parser = argparse.ArgumentParser(
        description="Create a new article manifest and markdown stub.",
        epilog=(
            "Run without arguments for interactive mode. "
            "Pass all flags for non-interactive (agent) use."
        ),
    )
    parser.add_argument("--slug", help="URL-safe identifier, e.g. my-article")
    parser.add_argument("--title", help="Article title")
    parser.add_argument("--description", help="Short summary shown in article cards")
    parser.add_argument(
        "--tags",
        nargs="+",
        help=f"One or more tags. Valid values: {', '.join(VALID_TAGS)}",
    )

    args = parser.parse_args()

    is_interactive = not any([args.slug, args.title, args.description, args.tags])

    if is_interactive:
        try:
            slug, title, description, tags = interactive_mode()
        except KeyboardInterrupt:
            print("\nAborted.")
            sys.exit(1)
    else:
        missing = [
            name
            for name, val in [
                ("--slug", args.slug),
                ("--title", args.title),
                ("--description", args.description),
                ("--tags", args.tags),
            ]
            if not val
        ]
        if missing:
            print(
                f"Error: CLI mode requires all flags. Missing: {', '.join(missing)}",
                file=sys.stderr,
            )
            sys.exit(1)
        slug, title, description, tags = cli_mode(args)

    manifest = build_manifest(slug, title, description, tags)
    write_files(slug, manifest)


if __name__ == "__main__":
    main()
