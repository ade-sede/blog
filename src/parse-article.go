package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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

type ArticleManifest struct {
	Title        string `json:"title"`
	Date         string `json:"date"`
	MarkdownFile string `json:"markdownFile, omitempty"`
	CssFile      string `json:"cssFile, omitempty"`
	ScriptFile   string `json:"scriptFile"`
	Description  string `json:"description"`
}

type Article struct {
	ManifestFilename string
	HTMLFilename     string
	StringifiedHTML  string
	Date             time.Time
	FormatedDate     string
	Manifest         *ArticleManifest
}

type filenameTitleTransformer struct{}

func (t *filenameTitleTransformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if cb, ok := n.(*ast.FencedCodeBlock); ok {
			language := string(cb.Language(reader.Source()))
			parts := strings.Split(language, ":")

			if len(parts) == 1 {
				cb.SetAttribute([]byte("language"), []byte(parts[0]))
			}

			if len(parts) == 2 {
				cb.SetAttribute([]byte("language"), []byte(parts[0]))
				cb.SetAttribute([]byte("filename"), []byte(parts[1]))
			}
		}
		return ast.WalkContinue, nil
	})
}

type codeBlockRenderer struct {
	html.Config
}

func (r *codeBlockRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindFencedCodeBlock, r.renderCodeBlock)
}

var (
	multiLineDisplayRegex  = regexp.MustCompile(`(?s)\\\\?\[\s*\n(.*?)\n\s*\\\\?\]`)
	singleLineDisplayRegex = regexp.MustCompile(`\\\\?\[(.*?)\\\\?\]`)
	inlineMathRegex        = regexp.MustCompile(`(?s)\\\\?\((.*?)\\\\?\)`)

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

type DirectoryNode struct {
	Name     string
	IsFolder bool
	Indent   int
	Children []*DirectoryNode
}

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

	currentParent := root
	stack := []*DirectoryNode{root}
	prevIndent := 0

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
			if currentIndent > prevIndent {
				stack = append(stack, currentParent)
				currentParent = stack[len(stack)-1].Children[len(stack[len(stack)-1].Children)-1]
			} else if currentIndent < prevIndent {
				for currentIndent < prevIndent && len(stack) > 1 {
					stack = stack[:len(stack)-1]
					currentParent = stack[len(stack)-1]
					prevIndent -= 2
				}
			}
		}

		currentParent.Children = append(currentParent.Children, entry)
		prevIndent = currentIndent
	}

	return root
}

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

func renderDirectoryStructure(w util.BufWriter, content string) {
	root := buildDirectoryTree(content)
	w.WriteString("<div class=\"directory-tree\">")
	renderDirectoryTreeRecursive(w, root)
	w.WriteString("</div>")
}

func (r *codeBlockRenderer) renderCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.FencedCodeBlock)

	if entering {
		filenameAttr, hasFilename := n.AttributeString("filename")
		languageAttr, hasLanguage := n.AttributeString("language")

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

		iterator, err := lexer.Tokenise(nil, code.String())
		if err != nil {
			w.WriteString("<div class=\"code-content-wrapper\">")
			w.WriteString("<pre><code>")
			w.WriteString(code.String())
			w.WriteString("</code></pre>")
			w.WriteString("<div class=\"code-copy-button\" title=\"Copy code\"><i class=\"fas fa-copy\"></i></div>")
			w.WriteString("</div>")
		} else {
			w.WriteString("<div class=\"highlight code-content-wrapper\">")
			err = formatter.Format(w, style, iterator)
			w.WriteString("<div class=\"code-copy-button\" title=\"Copy code\"><i class=\"fas fa-copy\"></i></div>")
			w.WriteString("</div>")
		}

		w.WriteString("</div>")

		return ast.WalkSkipChildren, nil
	}

	return ast.WalkContinue, nil
}

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

func formatDate(date time.Time) string {
	day := date.Day()
	month := date.Format("January")
	year := date.Year()
	return fmt.Sprintf("%s %d%s %d", month, day, getDateSuffix(day), year)
}

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

func processLatexExpressions(input string) string {
	processed := multiLineDisplayRegex.ReplaceAllStringFunc(input, func(match string) string {
		content := multiLineDisplayRegex.FindStringSubmatch(match)[1]
		return fmt.Sprintf("<div class=\"katex-display\" data-latex=\"%s\"></div>", content)
	})

	processed = singleLineDisplayRegex.ReplaceAllStringFunc(processed, func(match string) string {
		content := singleLineDisplayRegex.FindStringSubmatch(match)[1]
		return fmt.Sprintf("<span class=\"katex-inline\" data-latex=\"%s\"></span>", content)
	})

	processed = inlineMathRegex.ReplaceAllStringFunc(processed, func(match string) string {
		content := inlineMathRegex.FindStringSubmatch(match)[1]
		return fmt.Sprintf("<span class=\"katex-inline\" data-latex=\"%s\"></span>", content)
	})

	return processed
}

func parseArticleMarkdown(filename string) (string, error) {
	var buf bytes.Buffer
	input, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	processedInput := processLatexExpressions(string(input))

	htmlRenderer := renderer.NewRenderer(
		renderer.WithNodeRenderers(
			util.Prioritized(html.NewRenderer(
				html.WithHardWraps(),
				html.WithXHTML(),
				html.WithUnsafe(),
			), 100),
			util.Prioritized(&codeBlockRenderer{}, 80),
		),
	)

	p := goldmark.New(
		goldmark.WithRenderer(htmlRenderer),
		goldmark.WithExtensions(
			extension.GFM,
			extension.NewFootnote(
				extension.WithFootnoteBacklinkTitle("Return to text"),
				extension.WithFootnoteLinkTitle("See footnote"),
				extension.WithFootnoteBacklinkHTML("↩"),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithASTTransformers(
				util.Prioritized(&filenameTitleTransformer{}, 100),
			),
		),
	)

	err = p.Convert([]byte(processedInput), &buf)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func parseArticles(articleDir string) ([]Article, error) {
	files, err := os.ReadDir(articleDir)
	if err != nil {
		log.Fatalf("Error while opening directory '%s': '%v'", articleDir, err)
	}
	articles := make([]Article, 0)
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".json") {
			manifestFilename := filename
			manifestFullPath := articleDir + "/" + manifestFilename
			manifest, err := readArticleManifest(manifestFullPath)
			if err != nil {
				return nil, err
			}
			htmlFilename := strings.TrimSuffix(manifest.MarkdownFile, ".md") + ".html"
			date, err := time.Parse(time.DateOnly, manifest.Date)
			if err != nil {
				return nil, err
			}
			markdownFullPath := articleDir + "/" + manifest.MarkdownFile
			stringifiedHTML, err := parseArticleMarkdown(markdownFullPath)
			if err != nil {
				return nil, err
			}
			article := Article{
				ManifestFilename: manifestFilename,
				HTMLFilename:     htmlFilename,
				Date:             date,
				FormatedDate:     formatDate(date),
				Manifest:         manifest,
				StringifiedHTML:  stringifiedHTML,
			}
			articles = append(articles, article)
		}
	}
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].Date.After(articles[j].Date)
	})
	return articles, nil
}
