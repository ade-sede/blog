package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark"
)

type ArticleManifest struct {
	Title        string `json:"name"`
	Date         string `json:"date"`
	MarkdownFile string `json:"markdownFile"`
}

type Article struct {
	ManifestFilename string
	HTMLFilename     string
	StringifiedHTML  string
	Date             time.Time
	Manifest         *ArticleManifest
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

	md := goldmark.New()

	err = md.Convert(input, &buf)
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
			manifestFullPath := fmt.Sprintf("%s/%s", articleDir, manifestFilename)

			manifest, err := readArticleManifest(manifestFullPath)
			if err != nil {
				return nil, err
			}

			htmlFilename := strings.TrimSuffix(manifest.MarkdownFile, ".md") + ".html"
			date, err := time.Parse(time.DateOnly, manifest.Date)
			if err != nil {
				return nil, err
			}

			markdownFullPath := fmt.Sprintf("%s/%s", articleDir, manifest.MarkdownFile)
			stringifieldHTML, err := parseArticleMarkdown(markdownFullPath)
			if err != nil {
				return nil, err
			}

			article := Article{
				ManifestFilename: manifestFilename,
				HTMLFilename:     htmlFilename,
				Date:             date,
				Manifest:         manifest,
				StringifiedHTML:  stringifieldHTML,
			}

			articles = append(articles, article)
		}
	}

	// Newest articles first
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].Date.After(articles[j].Date)
	})

	return articles, nil
}
