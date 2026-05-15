#!/usr/bin/env python3
import os
import sys
from pathlib import Path
from weasyprint import HTML

output_dir = os.environ.get("OUTPUT_DIR")
src_dir = os.environ.get("SRC_DIR")

if not output_dir or not src_dir:
    print("Error: OUTPUT_DIR and SRC_DIR must be set", file=sys.stderr)
    sys.exit(1)

input_path = Path(output_dir) / "resume-printable.html"
output_path = Path(src_dir) / "pdfs" / "ade-sede.pdf"

if not input_path.exists():
    print(f"Error: {input_path} not found", file=sys.stderr)
    sys.exit(1)

output_path.parent.mkdir(parents=True, exist_ok=True)

HTML(filename=str(input_path)).write_pdf(str(output_path))
print(f"PDF written to {output_path}")
