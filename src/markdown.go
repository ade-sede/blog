package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/alecthomas/chroma"
	chromahtml "github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var (
	multiLineDisplayRegex  = regexp.MustCompile(`(?s)\\\[\s*\n(.*?)\n\s*\\\]`)
	singleLineDisplayRegex = regexp.MustCompile(`\\\[([^\n]*?)\\\]`)
	dollarInlineMathRegex  = regexp.MustCompile(`\$([^$\n]+)\$`)
	dynamicColorImageRegex = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)\{\.dynamic-colors\}`)
)

// filenameTitleTransformer is a Goldmark AST transformer that parses the
// language tag of fenced code blocks into structured attributes.
//
// It supports four formats:
//   - language
//   - language:filename
//   - language:diff
//   - language:filename:diff
type filenameTitleTransformer struct{}

// tocExtractor is a Goldmark AST transformer that walks the document
// and collects heading nodes into a table of contents.
type tocExtractor struct {
	TOC []TOCEntry
}

func (t *filenameTitleTransformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if cb, ok := n.(*ast.FencedCodeBlock); ok {
			language := string(cb.Language(reader.Source()))
			parts := strings.Split(language, ":")

			switch len(parts) {
			case 1:
				cb.SetAttribute([]byte("language"), []byte(parts[0]))
				cb.SetAttribute([]byte("isDiff"), []byte("false"))
			case 2:
				cb.SetAttribute([]byte("language"), []byte(parts[0]))
				if parts[1] == "diff" {
					cb.SetAttribute([]byte("isDiff"), []byte("true"))
				} else {
					cb.SetAttribute([]byte("filename"), []byte(parts[1]))
					cb.SetAttribute([]byte("isDiff"), []byte("false"))
				}
			case 3:
				cb.SetAttribute([]byte("language"), []byte(parts[0]))
				cb.SetAttribute([]byte("filename"), []byte(parts[1]))
				if parts[2] == "diff" {
					cb.SetAttribute([]byte("isDiff"), []byte("true"))
				}
			}
		}
		return ast.WalkContinue, nil
	})
}

func (toc *tocExtractor) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if heading, ok := n.(*ast.Heading); ok {
			var headingText strings.Builder

			ast.Walk(heading, func(child ast.Node, childEntering bool) (ast.WalkStatus, error) {
				if !childEntering {
					return ast.WalkContinue, nil
				}
				if textNode, ok := child.(*ast.Text); ok {
					headingText.Write(textNode.Segment.Value(reader.Source()))
				}
				return ast.WalkContinue, nil
			})

			headingID := ""
			if id, exists := heading.AttributeString("id"); exists {
				headingID = string(id.([]byte))
			}

			if headingText.Len() > 0 && heading.Level > 1 {
				toc.TOC = append(toc.TOC, TOCEntry{
					Level: heading.Level,
					Text:  headingText.String(),
					ID:    headingID,
				})
			}
		}
		return ast.WalkContinue, nil
	})
}

// codeBlockRenderer is a custom Goldmark renderer for fenced code blocks.
// It handles syntax highlighting via Chroma, diff annotations, filename
// headers, and directory-structure rendering.
type codeBlockRenderer struct {
	html.Config
}

func (r *codeBlockRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindFencedCodeBlock, r.renderCodeBlock)
}

// headingRenderer is a custom Goldmark renderer for headings. It adds
// anchor links and id attributes for linkable section headers.
type headingRenderer struct {
	html.Config
}

func (r *headingRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindHeading, r.renderHeading)
}

func (r *headingRenderer) renderHeading(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Heading)
	tag := fmt.Sprintf("h%d", n.Level)

	headingID := ""
	if id, exists := n.AttributeString("id"); exists {
		headingID = string(id.([]byte))
	}

	if entering {
		if headingID != "" {
			fmt.Fprintf(w, "<%s id=\"%s\" class=\"heading-with-anchor\">", tag, headingID)
		} else {
			fmt.Fprintf(w, "<%s>", tag)
		}
	} else {
		if headingID != "" {
			fmt.Fprintf(w, "<a href=\"#%s\" class=\"header-anchor\" title=\"Link to this section\"><i class=\"fas fa-link\"></i></a>", headingID)
		}
		fmt.Fprintf(w, "</%s>", tag)
	}

	return ast.WalkContinue, nil
}

func (r *codeBlockRenderer) renderCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.FencedCodeBlock)

	if entering {
		filenameAttr, hasFilename := n.AttributeString("filename")
		languageAttr, hasLanguage := n.AttributeString("language")
		isDiffAttr, hasDiff := n.AttributeString("isDiff")

		var code strings.Builder
		lines := n.Lines()
		for i := 0; i < lines.Len(); i++ {
			line := lines.At(i)
			code.Write(line.Value(source))
		}

		var lang string
		if hasLanguage {
			lang = string(languageAttr.([]byte))
		} else if language := n.Language(source); language != nil {
			lang = string(language)
		}

		if lang == "directory-structure" {
			renderDirectoryStructure(w, code.String())
			return ast.WalkSkipChildren, nil
		}

		w.WriteString("<div class=\"code-block\">")

		if hasFilename {
			filename := string(filenameAttr.([]byte))
			fileIcon := getFileIcon(filename)
			w.WriteString("<div class=\"code-filename\">")
			w.WriteString(fmt.Sprintf("<span><i class=\"%s\"></i> %s</span>", fileIcon, filename))
			w.WriteString(fmt.Sprintf("<span class=\"copy-filename\" title=\"Copy filename\" onclick=\"navigator.clipboard.writeText('%s')\"><i class=\"fas fa-copy\"></i></span>", filename))
			w.WriteString("</div>")
		}

		isDiff := hasDiff && string(isDiffAttr.([]byte)) == "true"

		var codeContent string
		var diffLineTypes []DiffLineType
		if isDiff {
			codeContent, diffLineTypes = processDiffContent(code.String())
		} else {
			codeContent = code.String()
		}

		lexer := lexers.Get(lang)
		if lexer == nil {
			lexer = lexers.Fallback
		}
		lexer = chroma.Coalesce(lexer)

		theme := "dracula"
		style := styles.Get(theme)
		if style == nil {
			style = styles.Fallback
		}

		formatter := chromahtml.New(
			chromahtml.WithClasses(true),
			chromahtml.WithLineNumbers(false),
		)

		iterator, err := lexer.Tokenise(nil, codeContent)
		if err != nil {
			w.WriteString("<div class=\"code-content-wrapper\">")
			w.WriteString("<pre><code>")
			var content string
			if isDiff {
				content = postProcessDiffHTML(codeContent, diffLineTypes)
			} else {
				content = code.String()
			}
			w.WriteString(content)
			w.WriteString("</code></pre>")
			w.WriteString("<div class=\"code-copy-button\" title=\"Copy code\"><i class=\"fas fa-copy\"></i></div>")
			w.WriteString("</div>")
		} else {
			w.WriteString("<div class=\"highlight code-content-wrapper\">")
			if isDiff {
				var buf strings.Builder
				if err = formatter.Format(&buf, style, iterator); err != nil {
					return ast.WalkStop, err
				}
				processedHTML := postProcessDiffHTML(buf.String(), diffLineTypes)
				w.WriteString(processedHTML)
			} else {
				if err = formatter.Format(w, style, iterator); err != nil {
					return ast.WalkStop, err
				}
			}
			w.WriteString("<div class=\"code-copy-button\" title=\"Copy code\"><i class=\"fas fa-copy\"></i></div>")
			w.WriteString("</div>")
		}

		w.WriteString("</div>")

		return ast.WalkSkipChildren, nil
	}

	return ast.WalkContinue, nil
}

// preprocessDynamicColorImages rewrites markdown image syntax with a
// {.dynamic-colors} attribute suffix into raw HTML img tags, adding the
// required class and crossorigin attributes for the color-shifting mechanism.
//
// Example input:  ![alt text](./images/foo.png){.dynamic-colors}
// Example output: <img src="./images/foo.png" alt="alt text" class="dynamic-colors" crossorigin="anonymous">
func preprocessDynamicColorImages(input string) string {
	return dynamicColorImageRegex.ReplaceAllStringFunc(input, func(match string) string {
		parts := dynamicColorImageRegex.FindStringSubmatch(match)
		alt, src := parts[1], parts[2]
		return fmt.Sprintf(`<img src="%s" alt="%s" class="dynamic-colors" crossorigin="anonymous">`, src, alt)
	})
}

// processLatexExpressions converts LaTeX math notation in markdown to
// HTML elements for client-side KaTeX rendering. Fenced code blocks
// are excluded from processing.
func processLatexExpressions(input string) string {
	codeBlockRegex := regexp.MustCompile("(?s)```.*?```")
	var blocks []string
	placeholder := "\x00CODEBLOCK%d\x00"

	stripped := codeBlockRegex.ReplaceAllStringFunc(input, func(match string) string {
		idx := len(blocks)
		blocks = append(blocks, match)
		return fmt.Sprintf(placeholder, idx)
	})

	processed := multiLineDisplayRegex.ReplaceAllStringFunc(stripped, func(match string) string {
		content := multiLineDisplayRegex.FindStringSubmatch(match)[1]
		return fmt.Sprintf("<div class=\"katex-display\" data-latex=\"%s\"></div>", content)
	})

	processed = singleLineDisplayRegex.ReplaceAllStringFunc(processed, func(match string) string {
		content := singleLineDisplayRegex.FindStringSubmatch(match)[1]
		return fmt.Sprintf("<div class=\"katex-display\" data-latex=\"%s\"></div>", content)
	})

	processed = dollarInlineMathRegex.ReplaceAllStringFunc(processed, func(match string) string {
		content := dollarInlineMathRegex.FindStringSubmatch(match)[1]
		return fmt.Sprintf("<span class=\"katex-inline\" data-latex=\"%s\"></span>", content)
	})

	for i, block := range blocks {
		processed = strings.Replace(processed, fmt.Sprintf(placeholder, i), block, 1)
	}

	return processed
}

// injectBylineBeforeFirstH1 inserts an author byline div immediately
// before the first <h1> tag in the HTML.
func injectBylineBeforeFirstH1(rawHTML string, formattedDate string, author string, authorImage string) (string, error) {
	h1Regex := regexp.MustCompile(`(<h1[^>]*>.*?</h1>)`)
	loc := h1Regex.FindStringIndex(rawHTML)
	if loc == nil {
		return "", fmt.Errorf("no h1 tag found in article HTML")
	}

	bylineHTML := fmt.Sprintf(`<div class="article-byline"><img src="images/%s" alt="%s" class="author-avatar"><div class="byline-content"><div class="date-line"><i class="far fa-calendar-alt"></i> %s</div><div class="author-line">by %s</div></div></div>`, authorImage, author, formattedDate, author)
	result := rawHTML[:loc[0]] + bylineHTML + rawHTML[loc[0]:]

	return result, nil
}

// parseArticleMarkdown converts a markdown file to HTML using Goldmark with
// custom renderers for code blocks and headings. Returns the processed HTML,
// extracted table of contents, and any error.
func parseArticleMarkdown(filename string, formattedDate string, author string, authorImage string) (string, []TOCEntry, error) {
	var buf bytes.Buffer
	input, err := os.ReadFile(filename)
	if err != nil {
		return "", nil, err
	}

	tocExtractor := &tocExtractor{TOC: []TOCEntry{}}

	processedInput := preprocessDynamicColorImages(string(input))
	processedInput = processLatexExpressions(processedInput)

	htmlRenderer := renderer.NewRenderer(
		renderer.WithNodeRenderers(
			util.Prioritized(html.NewRenderer(
				html.WithHardWraps(),
				html.WithXHTML(),
				html.WithUnsafe(),
			), 100),
			util.Prioritized(&codeBlockRenderer{}, 80),
			util.Prioritized(&headingRenderer{}, 70),
		),
	)

	p := goldmark.New(
		goldmark.WithRenderer(htmlRenderer),
		goldmark.WithExtensions(
			extension.GFM,
			extension.NewFootnote(
				extension.WithFootnoteBacklinkTitle("Return to text"),
				extension.WithFootnoteLinkTitle(""),
				extension.WithFootnoteBacklinkHTML("↩"),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithASTTransformers(
				util.Prioritized(&filenameTitleTransformer{}, 100),
				util.Prioritized(tocExtractor, 50),
			),
		),
	)

	err = p.Convert([]byte(processedInput), &buf)
	if err != nil {
		return "", nil, err
	}

	rawHTML := buf.String()
	processedHTML, err := injectBylineBeforeFirstH1(rawHTML, formattedDate, author, authorImage)
	if err != nil {
		return "", nil, err
	}

	return processedHTML, tocExtractor.TOC, nil
}

// hasFootnotes returns true if the HTML contains a footnotes section.
func hasFootnotes(html string) bool {
	return strings.Contains(html, `class="footnotes"`)
}

// processContentForSidebarFootnotes adds the "desktop-hidden" class to
// the inline footnotes div so they are hidden on desktop where the
// sidebar footnotes are shown instead.
func processContentForSidebarFootnotes(html string) string {
	if !hasFootnotes(html) {
		return html
	}

	re := regexp.MustCompile(`<div class="footnotes"([^>]*)>`)
	return re.ReplaceAllString(html, `<div class="footnotes desktop-hidden"$1>`)
}
