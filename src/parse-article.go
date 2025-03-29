package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	MarkdownFile string `json:"markdownFile"`
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

func getFileIcon(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	// File type specific icons
	switch ext {
	case ".qml":
		return "🧩" // QML files
	case ".json":
		return "📋" // JSON files
	case ".xml":
		return "📝" // XML files
	case ".md":
		return "📄" // Markdown files
	case ".go":
		return "📄" // Go files
	case ".js":
		return "📄" // JavaScript files
	case ".html", ".htm":
		return "📄" // HTML files
	case ".css":
		return "📄" // CSS files
	case ".py":
		return "📄" // Python files
	case ".sh", ".bash", ".fish":
		return "📄" // Shell scripts
	default:
		if strings.HasPrefix(ext, ".") {
			return "📄" // Other files with extensions
		}
		return "📄" // Default
	}
}

func renderDirectoryStructure(w util.BufWriter, content string) {
	type DirEntry struct {
		Name     string
		IsFolder bool
		Indent   int
		Children []*DirEntry
	}

	var lines []string
	for _, line := range strings.Split(content, "\n") {
		if len(strings.TrimSpace(line)) > 0 {
			lines = append(lines, line)
		}
	}

	if len(lines) == 0 {
		return
	}

	// Root entry
	root := &DirEntry{
		Children: []*DirEntry{},
	}

	// Build tree
	var currentParent *DirEntry = root
	var stack []*DirEntry = []*DirEntry{root}
	var prevIndent int = 0

	for i, line := range lines {
		trimmedLine := strings.TrimLeft(line, " \t")
		currentIndent := len(line) - len(trimmedLine)

		isFolder := strings.HasSuffix(trimmedLine, "/")
		entry := &DirEntry{
			Name:     trimmedLine,
			IsFolder: isFolder,
			Indent:   currentIndent,
			Children: []*DirEntry{},
		}

		// Determine parent based on indentation
		if i > 0 {
			if currentIndent > prevIndent {
				// This is a child of the previous entry
				stack = append(stack, currentParent)
				currentParent = stack[len(stack)-1].Children[len(stack[len(stack)-1].Children)-1]
			} else if currentIndent < prevIndent {
				// Go back up the tree
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

	// Render the tree
	w.WriteString("<div class=\"directory-tree\">")

	var renderEntry func(entry *DirEntry)
	renderEntry = func(entry *DirEntry) {
		for _, child := range entry.Children {
			if child.IsFolder {
				w.WriteString(fmt.Sprintf("<div class=\"dir-entry dir-folder\"><span class=\"dir-icon\">📁</span> %s</div>", child.Name))
			} else {
				fileIcon := getFileIcon(child.Name)
				w.WriteString(fmt.Sprintf("<div class=\"dir-entry dir-file\"><span class=\"dir-icon\">%s</span> %s</div>", fileIcon, child.Name))
			}

			if len(child.Children) > 0 {
				w.WriteString("<div class=\"dir-children\">")
				renderEntry(child)
				w.WriteString("</div>")
			}
		}
	}

	renderEntry(root)

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

		// Handle directory structure special case
		if lang == "directory-structure" {
			renderDirectoryStructure(w, code.String())
			return ast.WalkSkipChildren, nil
		}

		w.WriteString("<div class=\"code-block\">")

		if hasFilename {
			w.WriteString("<div class=\"code-filename\">")
			w.WriteString(string(filenameAttr.([]byte)))
			w.WriteString("</div>")
		}

		lexer := lexers.Get(lang)
		if lexer == nil {
			lexer = lexers.Fallback
		}
		lexer = chroma.Coalesce(lexer)

		// We don't actually use monokai colors
		// Specifying a theme just helps us specify how to format: should classnames have the same color as strings ? etc ...
		// Colors are handled dynamcally using css + js
		theme := "monokai"
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
			w.WriteString("<pre><code>")
			w.WriteString(code.String())
			w.WriteString("</code></pre>")
		} else {
			w.WriteString("<div class=\"highlight\">")
			err = formatter.Format(w, style, iterator)
			w.WriteString("</div>")
		}

		w.WriteString("</div>")

		return ast.WalkSkipChildren, nil
	}

	return ast.WalkContinue, nil
}

func ordinal(day int) string {
	if day%10 == 1 && day != 11 {
		return "st"
	} else if day%10 == 2 && day != 12 {
		return "nd"
	} else if day%10 == 3 && day != 13 {
		return "rd"
	}
	return "th"
}

func formatDate(date time.Time) string {
	day := date.Day()
	month := date.Format("January")
	year := date.Year()
	return fmt.Sprintf("%s %d%s %d", month, day, ordinal(day), year)
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

func parseArticleMarkdown(filename string) (string, error) {
	var buf bytes.Buffer
	input, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

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
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithASTTransformers(
				util.Prioritized(&filenameTitleTransformer{}, 100),
			),
		),
	)

	err = p.Convert(input, &buf)
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
			stringifieldHTML, err := parseArticleMarkdown(markdownFullPath)
			if err != nil {
				return nil, err
			}
			article := Article{
				ManifestFilename: manifestFilename,
				HTMLFilename:     htmlFilename,
				Date:             date,
				FormatedDate:     formatDate(date),
				Manifest:         manifest,
				StringifiedHTML:  stringifieldHTML,
			}
			articles = append(articles, article)
		}
	}
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].Date.After(articles[j].Date)
	})
	return articles, nil
}
