package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
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
	FormattedDate    string
	Manifest         *ArticleManifest
	TOC              []TOCEntry
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
				FormattedDate:    formattedDate,
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
