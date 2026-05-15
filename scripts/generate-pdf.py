#!/usr/bin/env python3
#
# PDF GENERATION
# ==============
#
# WeasyPrint converts HTML to PDF but is not a browser:
#   - It does not execute JavaScript
#   - It does not share a CSS cascade with other documents
#   - Relative paths are resolved from the HTML file location, which may be
#     a temporary file with no relationship to the asset directories
#
# For these reasons this script builds a fully self-contained HTML document:
#   - All CSS is inlined — no external stylesheets
#   - CSS variables are replaced with hard-coded values — no cascade needed
#   - Print dimensions are static CSS — no JavaScript needed
#   - Font faces use absolute file:// URIs resolved from SRC_DIR
#   - Font Awesome webfonts are embedded as base64 data URIs
#   - The profile image uses an absolute file:// URI
#   - Resume content is loaded directly from experiences.json
#
# To edit the resume layout/styles: edit the build_html() function below.
# To edit the resume content: edit src/experiences.json.
#
# PAGE SIZING
# ===========
#
# The @page rule declares size: A4 with margin: 1.5cm on all sides.
#
# A4 is 210mm × 297mm (793.7 × 1122.5 CSS px at 96 dpi). Subtracting
# 1.5cm margins on each side gives a content area of ~681px wide. 1.5cm
# is the conventional margin for professional printed documents: narrow
# enough to make good use of the page, wide enough to remain comfortable
# when the document is physically held or bound.
#
# The resume container uses width: 100% so it always fills exactly the
# content area defined by @page, with no risk of overflow or clipping
# regardless of the page size chosen.
#
# Body font-size is 8pt. 8pt is the conventional minimum for comfortable
# reading in print; anything smaller risks being illegible on physical
# paper. Heading sizes are expressed in rem so they scale with this base.

import base64
import json
import os
import sys
import tempfile
from pathlib import Path

from weasyprint import HTML

src_dir = os.environ.get("SRC_DIR")
if not src_dir:
    print("Error: SRC_DIR must be set", file=sys.stderr)
    sys.exit(1)

src = Path(src_dir)

with open(src / "experiences.json") as f:
    data = json.load(f)

def b64(path):
    return base64.b64encode(path.read_bytes()).decode("ascii")

def font_uri(path):
    return path.as_uri()

def entry_html(exp):
    bullet_html = ""
    if exp.get("bulletPoints"):
        intro = ""
        if exp.get("bulletPointsIntro"):
            intro = f'<p class="bullet-intro">{exp["bulletPointsIntro"]}:</p>'
        items = "".join(f"<li>{p}</li>" for p in exp["bulletPoints"])
        bullet_html = f'<div class="bullet-points">{intro}<ul>{items}</ul></div>'

    return f"""
    <div class="entry">
      <div class="header">
        <div class="position">
          <span class="company">{exp["company"]}</span>
          <span class="title">{exp["title"]}</span>
        </div>
        <span class="timeline">{exp["begin"]} - {exp["end"]}</span>
      </div>
      <div class="description">{exp["description"]}</div>
      {bullet_html}
    </div>"""

def build_html():
    fa_solid = b64(src / "webfonts" / "fa-solid-900.woff2")
    fa_brands = b64(src / "webfonts" / "fa-brands-400.woff2")

    work_entries = "".join(entry_html(e) for e in data["workExperiences"])
    school_entries = "".join(entry_html(e) for e in data["schoolExperiences"])

    picture = (src / "images" / "picture.webp").as_uri()
    montserrat_regular = font_uri(src / "fonts" / "montserrat" / "Montserrat-Regular.ttf")
    montserrat_bold    = font_uri(src / "fonts" / "montserrat" / "Montserrat-Bold.ttf")
    montserrat_italic  = font_uri(src / "fonts" / "montserrat" / "Montserrat-Italic.ttf")

    return f"""<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8"/>
  <title>Adrien DE SEDE - Resume</title>
  <style>
    @font-face {{
      font-family: "Montserrat";
      src: url("{montserrat_regular}") format("truetype");
      font-weight: normal;
      font-style: normal;
    }}
    @font-face {{
      font-family: "Montserrat";
      src: url("{montserrat_bold}") format("truetype");
      font-weight: bold;
      font-style: normal;
    }}
    @font-face {{
      font-family: "Montserrat";
      src: url("{montserrat_italic}") format("truetype");
      font-weight: normal;
      font-style: italic;
    }}
    @font-face {{
      font-family: "Font Awesome 6 Free";
      src: url("data:font/woff2;base64,{fa_solid}") format("woff2");
      font-weight: 900;
      font-style: normal;
    }}
    @font-face {{
      font-family: "Font Awesome 6 Brands";
      src: url("data:font/woff2;base64,{fa_brands}") format("woff2");
      font-weight: 400;
      font-style: normal;
    }}
    @page {{
      size: A4;
      margin: 1.5cm;
    }}
    * {{ box-sizing: border-box; margin: 0; padding: 0; }}
    html, body {{
      background: white;
      color: black;
      font-family: "Montserrat", sans-serif;
      font-size: 8pt;
      line-height: 1.5;
    }}
    a {{ color: #333; text-decoration: none; }}
    ul {{ padding-left: 1.5em; }}
    hr {{
      height: 1px;
      width: 80%;
      margin: 0.5rem auto;
      border: none;
      background: linear-gradient(to right, transparent, #333, transparent);
    }}
    .fas, .fab {{
      font-style: normal;
      font-variant: normal;
      text-rendering: auto;
      line-height: 1;
      display: inline-block;
    }}
    .fas {{ font-family: "Font Awesome 6 Free"; font-weight: 900; }}
    .fab {{ font-family: "Font Awesome 6 Brands"; font-weight: 400; }}
    .fa-location-dot::before {{ content: "\\f3c5"; }}
    .fa-envelope::before     {{ content: "\\f0e0"; }}
    .fa-github::before       {{ content: "\\f09b"; }}
    .fa-linkedin::before     {{ content: "\\f08c"; }}
    .fa-globe::before        {{ content: "\\f0ac"; }}
    .resume-container {{
      display: flex;
      flex-direction: column;
      width: 100%;
    }}
    .resume {{ display: flex; flex-direction: column; }}
    .resume img {{ border-radius: 50%; width: 80px; height: 80px; }}
    .resume .header {{
      display: flex;
      flex-direction: row;
      padding: 1rem;
      gap: 1rem;
    }}
    .resume .header .left-container {{
      display: flex;
      align-items: center;
      justify-content: center;
      width: 30%;
    }}
    .resume .header .right-container {{
      display: flex;
      flex-direction: column;
      justify-content: center;
      width: 70%;
    }}
    .resume .info {{
      display: flex;
      flex-direction: column;
      align-items: flex-end;
    }}
    .resume .info h1.name {{
      font-weight: bold;
      line-height: 1.3;
      font-size: 1.5rem;
      color: black;
    }}
    .resume .section {{
      display: flex;
      flex-direction: column;
      padding-left: 1rem;
      padding-right: 1rem;
    }}
    .resume .section h1 {{
      font-weight: bold;
      font-size: 1.2rem;
      line-height: 1;
      align-self: center;
      color: black;
      margin: 0.5rem 0;
    }}
    .resume .section .links {{
      display: flex;
      flex-direction: row;
      align-self: center;
      gap: 2rem;
      padding: 0.5rem;
    }}
    .resume .entry {{
      display: flex;
      flex-direction: column;
      margin-bottom: 0.5rem;
    }}
    .resume .entry .bullet-intro {{ font-weight: bold; }}
    .resume .entry .description {{
      margin: 4px 0;
      font-style: italic;
    }}
    .resume .entry .header {{
      display: flex;
      flex-direction: row;
      justify-content: space-between;
      padding-top: 8px;
      padding-bottom: 4px;
    }}
    .resume .entry .header .position {{
      display: flex;
      flex-direction: row;
      align-items: center;
    }}
    .resume .entry .header .position .company {{ font-weight: bold; }}
    .resume .entry .header .position .title::before {{
      content: ", ";
      margin-left: 2px;
    }}
    .resume .entry .header .timeline {{ color: #333; }}
  </style>
</head>
<body>
  <div class="resume-container">
    <div class="resume">
      <div class="header">
        <div class="left-container">
          <img alt="Adrien DE SEDE" width="80" height="80" src="{picture}"/>
        </div>
        <div class="right-container">
          <div class="info">
            <h1 class="name">Adrien DE SEDE</h1>
            <span class="address"><i class="fas fa-location-dot"></i> Lyon, France</span>
            <span class="mail">
              <a href="mailto:contact@ade-sede.dev">
                <i class="fas fa-envelope"></i> contact@ade-sede.dev
              </a>
            </span>
          </div>
        </div>
      </div>
      <hr/>
      <div class="section">
        <h1>WORK HISTORY</h1>
        {work_entries}
      </div>
      <hr/>
      <div class="section">
        <h1>EDUCATION</h1>
        {school_entries}
      </div>
      <hr/>
      <div class="section">
        <div class="links">
          <span><a href="https://github.com/ade-sede"><i class="fab fa-github"></i> github.com/ade-sede</a></span>
          <span><a href="https://www.linkedin.com/in/ade-sede"><i class="fab fa-linkedin"></i> linkedin.com/in/ade-sede</a></span>
          <span><a href="https://blog.ade-sede.dev"><i class="fas fa-globe"></i> blog.ade-sede.dev</a></span>
        </div>
      </div>
    </div>
  </div>
</body>
</html>"""

output_path = src / "pdfs" / "resume.pdf"
output_path.parent.mkdir(parents=True, exist_ok=True)

with tempfile.NamedTemporaryFile(suffix=".html", mode="w", delete=False) as f:
    f.write(build_html())
    tmp_path = f.name

try:
    HTML(filename=tmp_path).write_pdf(str(output_path))
    print(f"PDF written to {output_path}")
finally:
    os.unlink(tmp_path)
