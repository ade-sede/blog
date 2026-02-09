package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

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

// ArticleManifest represents the JSON metadata for a blog article.
type ArticleManifest struct {
	Title        string   `json:"title"`
	Date         string   `json:"date"`
	Draft        bool     `json:"draft,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	MarkdownFile string   `json:"markdownFile,omitempty"`
	CssFile      string   `json:"cssFile,omitempty"`
	ScriptFile   string   `json:"scriptFile"`
	Description  string   `json:"description"`
	Author       string   `json:"author"`
	AuthorImage  string   `json:"authorImage"`
}

// TOCEntry represents a single entry in the table of contents,
// extracted from markdown headings.
type TOCEntry struct {
	Level int
	Text  string
	ID    string
}

// Article represents a fully parsed blog article with its HTML content,
// metadata, and table of contents.
type Article struct {
	ManifestFilename string
	HTMLFilename     string
	StringifiedHTML  string
	Date             time.Time
	FormatedDate     string
	Manifest         *ArticleManifest
	TOC              []TOCEntry
}

// filenameTitleTransformer supports 4 formats:
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

type DiffLineType int

const (
	DiffContext DiffLineType = iota
	DiffAddition
	DiffDeletion
)

var (
	multiLineDisplayRegex  = regexp.MustCompile(`(?s)\\\[\s*\n(.*?)\n\s*\\\]`)
	singleLineDisplayRegex = regexp.MustCompile(`\\\[([^\n]*?)\\\]`)
	dollarInlineMathRegex  = regexp.MustCompile(`\$([^$\n]+)\$`)

	filenameIcons = map[string]string{
		"makefile":          "fas fa-cogs",
		"dockerfile":        "fab fa-docker",
		"readme.md":         "fas fa-book",
		"license":           "fas fa-certificate",
		"license.txt":       "fas fa-certificate",
		"license.md":        "fas fa-certificate",
		"go.mod":            "fab fa-golang",
		"go.sum":            "fab fa-golang",
		"package.json":      "fab fa-npm",
		"package-lock.json": "fab fa-npm",
		".gitignore":        "fab fa-git-alt",
		".gitmodules":       "fab fa-git-alt",
		".gitattributes":    "fab fa-git-alt",
	}

	extensionIcons = map[string]string{
		".qml":       "fas fa-file-code",
		".json":      "fas fa-file-code",
		".go":        "fab fa-golang",
		".js":        "fab fa-js",
		".ts":        "fab fa-js-square",
		".html":      "fab fa-html5",
		".htm":       "fab fa-html5",
		".css":       "fab fa-css3-alt",
		".py":        "fab fa-python",
		".rb":        "fas fa-gem",
		".php":       "fab fa-php",
		".java":      "fab fa-java",
		".rs":        "fas fa-gears",
		".c":         "fas fa-code",
		".cpp":       "fas fa-code",
		".h":         "fas fa-code",
		".hpp":       "fas fa-code",
		".cs":        "fas fa-code",
		".swift":     "fas fa-code",
		".kt":        "fas fa-k",
		".kts":       "fas fa-k",
		".yml":       "fas fa-file-code",
		".yaml":      "fas fa-file-code",
		".toml":      "fas fa-cog",
		".ini":       "fas fa-sliders-h",
		".conf":      "fas fa-sliders-h",
		".config":    "fas fa-sliders-h",
		".env":       "fas fa-key",
		".xml":       "fas fa-code",
		".csv":       "fas fa-file-csv",
		".sql":       "fas fa-database",
		".txt":       "fas fa-file-lines",
		".pdf":       "fas fa-file-pdf",
		".doc":       "fas fa-file-word",
		".docx":      "fas fa-file-word",
		".xls":       "fas fa-file-excel",
		".xlsx":      "fas fa-file-excel",
		".ppt":       "fas fa-file-powerpoint",
		".pptx":      "fas fa-file-powerpoint",
		".svg":       "fas fa-bezier-curve",
		".mp3":       "fas fa-file-audio",
		".wav":       "fas fa-file-audio",
		".ogg":       "fas fa-file-audio",
		".flac":      "fas fa-file-audio",
		".mp4":       "fas fa-file-video",
		".mov":       "fas fa-file-video",
		".avi":       "fas fa-file-video",
		".mkv":       "fas fa-file-video",
		".webm":      "fas fa-file-video",
		".vue":       "fab fa-vuejs",
		".svelte":    "fas fa-fire",
		".gradle":    "fab fa-android",
		".xcodeproj": "fab fa-apple",
	}
)

// getFileIcon returns a Font Awesome CSS class for the given filename,
// matching first by exact basename then by file extension.
func getFileIcon(filename string) string {
	basename := strings.ToLower(filepath.Base(filename))
	if icon, exists := filenameIcons[basename]; exists {
		return icon
	}

	ext := strings.ToLower(filepath.Ext(filename))
	if icon, exists := extensionIcons[ext]; exists {
		return icon
	}

	return ""
}

// DirectoryNode represents a file or folder in a parsed directory tree
// used for rendering directory-structure code blocks.
type DirectoryNode struct {
	Name     string
	IsFolder bool
	Indent   int
	Children []*DirectoryNode
}

// buildDirectoryTree parses indented text into a tree of DirectoryNode,
// using indentation levels to determine parent-child relationships.
func buildDirectoryTree(content string) *DirectoryNode {
	var lines []string
	for _, line := range strings.Split(content, "\n") {
		if len(strings.TrimSpace(line)) > 0 {
			lines = append(lines, line)
		}
	}

	root := &DirectoryNode{Children: []*DirectoryNode{}}
	if len(lines) == 0 {
		return root
	}

	type stackEntry struct {
		node   *DirectoryNode
		indent int
	}

	currentParent := root
	stack := []stackEntry{{node: root, indent: -1}}

	for i, line := range lines {
		trimmedLine := strings.TrimLeft(line, " \t")
		currentIndent := len(line) - len(trimmedLine)

		entry := &DirectoryNode{
			Name:     trimmedLine,
			IsFolder: strings.HasSuffix(trimmedLine, "/"),
			Indent:   currentIndent,
			Children: []*DirectoryNode{},
		}

		if i > 0 {
			if currentIndent > stack[len(stack)-1].indent {
				lastChild := currentParent.Children[len(currentParent.Children)-1]
				stack = append(stack, stackEntry{node: currentParent, indent: currentParent.Indent})
				currentParent = lastChild
			} else if currentIndent < stack[len(stack)-1].indent {
				for len(stack) > 1 && stack[len(stack)-1].indent >= currentIndent {
					stack = stack[:len(stack)-1]
					currentParent = stack[len(stack)-1].node
				}
			}
		}

		currentParent.Children = append(currentParent.Children, entry)
	}

	return root
}

// renderDirectoryTreeRecursive writes HTML for a DirectoryNode and its
// children, nesting folders within dir-children divs.
func renderDirectoryTreeRecursive(w util.BufWriter, entry *DirectoryNode) {
	for _, child := range entry.Children {
		if child.IsFolder {
			w.WriteString(fmt.Sprintf("<div class=\"dir-entry dir-folder\"><i class=\"fas fa-folder-open\"></i> %s</div>", child.Name))
		} else {
			fileIcon := getFileIcon(child.Name)
			w.WriteString(fmt.Sprintf("<div class=\"dir-entry dir-file\"><i class=\"%s\"></i> %s</div>", fileIcon, child.Name))
		}

		if len(child.Children) > 0 {
			w.WriteString("<div class=\"dir-children\">")
			renderDirectoryTreeRecursive(w, child)
			w.WriteString("</div>")
		}
	}
}

// renderDirectoryStructure parses indented text into a directory tree
// and writes the complete HTML representation to w.
func renderDirectoryStructure(w util.BufWriter, content string) {
	root := buildDirectoryTree(content)
	w.WriteString("<div class=\"directory-tree\">")
	renderDirectoryTreeRecursive(w, root)
	w.WriteString("</div>")
}

// processDiffContent strips leading +/- markers from diff-formatted code,
// returning the cleaned content and a per-line type classification.
func processDiffContent(content string) (string, []DiffLineType) {
	lines := strings.Split(content, "\n")
	var cleanedLines []string
	var lineTypes []DiffLineType

	for _, line := range lines {
		if len(line) == 0 {
			cleanedLines = append(cleanedLines, line)
			lineTypes = append(lineTypes, DiffContext)
			continue
		}

		if line[0] == '+' {
			cleanedLines = append(cleanedLines, line[1:])
			lineTypes = append(lineTypes, DiffAddition)
		} else if line[0] == '-' {
			cleanedLines = append(cleanedLines, line[1:])
			lineTypes = append(lineTypes, DiffDeletion)
		} else {
			cleanedLines = append(cleanedLines, line)
			lineTypes = append(lineTypes, DiffContext)
		}
	}

	return strings.Join(cleanedLines, "\n"), lineTypes
}

// postProcessDiffHTML wraps Chroma-highlighted HTML lines with diff CSS
// classes (gi for additions, gd for deletions) based on line types.
func postProcessDiffHTML(html string, lineTypes []DiffLineType) string {
	lines := strings.Split(html, "\n")
	var processedLines []string

	for i, line := range lines {
		if i < len(lineTypes) {
			switch lineTypes[i] {
			case DiffAddition:
				if strings.Contains(line, "<span class=\"line\">") {
					processedLines = append(processedLines, strings.Replace(line, "<span class=\"line\">", "<span class=\"line gi\">", 1))
				} else {
					processedLines = append(processedLines, fmt.Sprintf("<span class=\"gi\">%s</span>", line))
				}
			case DiffDeletion:
				if strings.Contains(line, "<span class=\"line\">") {
					processedLines = append(processedLines, strings.Replace(line, "<span class=\"line\">", "<span class=\"line gd\">", 1))
				} else {
					processedLines = append(processedLines, fmt.Sprintf("<span class=\"gd\">%s</span>", line))
				}
			default:
				processedLines = append(processedLines, line)
			}
		} else {
			processedLines = append(processedLines, line)
		}
	}

	return strings.Join(processedLines, "\n")
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

// getDateSuffix returns the English ordinal suffix for a day number
// (e.g. "st", "nd", "rd", "th").
func getDateSuffix(day int) string {
	if day >= 11 && day <= 13 {
		return "th"
	}
	switch day % 10 {
	case 1:
		return "st"
	case 2:
		return "nd"
	case 3:
		return "rd"
	default:
		return "th"
	}
}

// formatDate formats a time.Time into a human-readable string
// like "January 2nd 2025".
func formatDate(date time.Time) string {
	day := date.Day()
	month := date.Format("January")
	year := date.Year()
	return fmt.Sprintf("%s %d%s %d", month, day, getDateSuffix(day), year)
}

// readArticleManifest reads and unmarshals a JSON manifest file into
// an ArticleManifest.
func readArticleManifest(filename string) (*ArticleManifest, error) {
	var manifest ArticleManifest
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(content, &manifest)
	if err != nil {
		return nil, err
	}
	return &manifest, nil
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

	processedInput := processLatexExpressions(string(input))

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

// parseArticles reads all JSON manifests from articleDir, parses their
// corresponding markdown files, and returns the articles sorted by date
// descending. Draft articles are excluded unless env is "development".
func parseArticles(articleDir, env string) ([]Article, error) {
	files, err := os.ReadDir(articleDir)
	if err != nil {
		return nil, fmt.Errorf("error while opening directory '%s': '%w'", articleDir, err)
	}
	articles := make([]Article, 0)
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".json") {
			manifestFilename := filename
			manifestFullPath := filepath.Join(articleDir, manifestFilename)
			manifest, err := readArticleManifest(manifestFullPath)
			if err != nil {
				return nil, err
			}
			if manifest.Draft && env != "development" {
				continue
			}
			htmlFilename := strings.TrimSuffix(manifest.MarkdownFile, ".md") + ".html"
			date, err := time.Parse(time.DateOnly, manifest.Date)
			if err != nil {
				return nil, err
			}
			markdownFullPath := filepath.Join(articleDir, manifest.MarkdownFile)
			formattedDate := formatDate(date)
			stringifiedHTML, toc, err := parseArticleMarkdown(markdownFullPath, formattedDate, manifest.Author, manifest.AuthorImage)
			if err != nil {
				return nil, err
			}
			article := Article{
				ManifestFilename: manifestFilename,
				HTMLFilename:     htmlFilename,
				Date:             date,
				FormatedDate:     formattedDate,
				Manifest:         manifest,
				StringifiedHTML:  stringifiedHTML,
				TOC:              toc,
			}
			articles = append(articles, article)
		}
	}
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].Date.After(articles[j].Date)
	})
	return articles, nil
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

	processedHTML := html

	re := regexp.MustCompile(`<div class="footnotes"([^>]*)>`)
	processedHTML = re.ReplaceAllString(processedHTML, `<div class="footnotes desktop-hidden"$1>`)

	return processedHTML
}
