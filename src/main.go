package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/a-h/templ"
)

// AssetKind identifies the type of an inlined asset.
type AssetKind string

const (
	AssetCSS AssetKind = "css"
	AssetJS  AssetKind = "js"
)

// AssetScope determines which directory to load the asset from.
type AssetScope string

const (
	AssetScopeGlobal  AssetScope = "global"  // loaded from srcDir
	AssetScopeArticle AssetScope = "article" // loaded from articleDir
)

// Asset describes a CSS or JS file that will be minified and inlined
// into the <head> of a page.
type Asset struct {
	Filename string
	Kind     AssetKind
	Scope    AssetScope
}

// Page describes a single HTML page to be generated: its output filename,
// the assets to inline, and a Render function that produces the templ component.
//
// Render receives the resolved inline style and script tags so that each
// page template can embed them directly in the <head>.
type Page struct {
	Filename string
	Assets   []Asset
	Render   func(styleTags, scriptTags []string) templ.Component
}

// inlineAssets minifies all assets declared on a page and returns
// the resulting <style> and <script> tag strings for inline embedding.
func inlineAssets(page Page, srcDir, articleDir string) ([]string, []string, error) {
	var styleTags []string
	var scriptTags []string

	for _, asset := range page.Assets {
		var content string
		var err error

		switch asset.Scope {
		case AssetScopeGlobal:
			switch asset.Kind {
			case AssetCSS:
				content, err = loadAndMinifyGlobalStyle(srcDir, asset.Filename)
			case AssetJS:
				content, err = loadAndMinifyGlobalScript(srcDir, asset.Filename)
			}
		case AssetScopeArticle:
			switch asset.Kind {
			case AssetCSS:
				content, err = loadAndMinifyArticleStyle(articleDir, asset.Filename)
			case AssetJS:
				content, err = loadAndMinifyArticleScript(articleDir, asset.Filename)
			}
		}

		if err != nil {
			return nil, nil, fmt.Errorf("failed to inline asset %s: %w", asset.Filename, err)
		}

		switch asset.Kind {
		case AssetCSS:
			styleTags = append(styleTags, fmt.Sprintf("<style type='text/css'>%s</style>", content))
		case AssetJS:
			scriptTags = append(scriptTags, fmt.Sprintf("<script>%s</script>", content))
		}
	}

	return styleTags, scriptTags, nil
}

// globalCSS returns an Asset for a CSS file loaded from srcDir.
func globalCSS(filename string) Asset {
	return Asset{Filename: filename, Kind: AssetCSS, Scope: AssetScopeGlobal}
}

// globalJS returns an Asset for a JS file loaded from srcDir.
func globalJS(filename string) Asset {
	return Asset{Filename: filename, Kind: AssetJS, Scope: AssetScopeGlobal}
}

// articleCSS returns an Asset for a CSS file loaded from articleDir.
func articleCSS(filename string) Asset {
	return Asset{Filename: filename, Kind: AssetCSS, Scope: AssetScopeArticle}
}

// articleJS returns an Asset for a JS file loaded from articleDir.
func articleJS(filename string) Asset {
	return Asset{Filename: filename, Kind: AssetJS, Scope: AssetScopeArticle}
}

// articlePage builds a Page for a single blog article.
func articlePage(a Article) Page {
	assets := []Asset{
		globalCSS("article.css"),
		globalCSS("syntax-highlighting.css"),
		globalJS("toc.js"),
		globalJS("anchors.js"),
		globalJS("footnotes.js"),
	}
	if a.Manifest.CssFile != "" {
		assets = append(assets, articleCSS(a.Manifest.CssFile))
	}
	if a.Manifest.ScriptFile != "" {
		assets = append(assets, articleJS(a.Manifest.ScriptFile))
	}

	return Page{
		Filename: a.HTMLFilename,
		Assets:   assets,
		Render: func(styleTags, scriptTags []string) templ.Component {
			return article(a.Manifest.Title, a.Manifest.Description, a.StringifiedHTML, a.FormattedDate, scriptTags, styleTags, a.TOC)
		},
	}
}

// publishGlobalCSS minifies each CSS file listed in config.CSSFiles and
// writes the result to the output CSS directory.
func publishGlobalCSS(config Config) error {
	for _, filename := range config.CSSFiles {
		minifiedCSS, err := MinifyCSS(config.SrcDir + "/css/" + filename)
		if err != nil {
			return fmt.Errorf("failed to minify %s: %v", filename, err)
		}

		newFilename := config.OutputDir + "/css/" + filename
		if err := os.WriteFile(newFilename, []byte(minifiedCSS), config.FileMode); err != nil {
			return fmt.Errorf("failed to write %s: %v", newFilename, err)
		}
	}
	return nil
}

// buildPages constructs all Page descriptors for the site: the four static
// pages (home, articles, resume, resume-printable) plus one page per article.
func buildPages(allArticles []Article, experiences ExperiencesData) []Page {
	pages := []Page{
		{
			Filename: "index.html",
			Assets: []Asset{
				globalCSS("home.css"),
				globalCSS("articles.css"),
			},
			Render: func(styleTags, scriptTags []string) templ.Component {
				var mostRecent []Article
				if len(allArticles) > 0 {
					mostRecent = allArticles[:1]
				}
				return home(mostRecent, styleTags)
			},
		},
		{
			Filename: "articles.html",
			Assets: []Asset{
				globalCSS("articles.css"),
			},
			Render: func(styleTags, scriptTags []string) templ.Component {
				return articles(allArticles, styleTags)
			},
		},
		{
			Filename: "resume.html",
			Assets: []Asset{
				globalCSS("resume.css"),
			},
			Render: func(styleTags, scriptTags []string) templ.Component {
				return resumePage(experiences.WorkExperiences, experiences.SchoolExperiences, styleTags)
			},
		},
		{
			Filename: "resume-printable.html",
			Assets: []Asset{
				globalCSS("resume.css"),
			},
			Render: func(styleTags, scriptTags []string) templ.Component {
				return resumePrintReady(experiences.WorkExperiences, experiences.SchoolExperiences, styleTags)
			},
		},
	}

	for _, a := range allArticles {
		pages = append(pages, articlePage(a))
	}

	return pages
}

// generateAllPages renders each Page to an HTML file in the output directory.
func generateAllPages(config Config, allArticles []Article, experiences ExperiencesData) error {
	pages := buildPages(allArticles, experiences)

	for _, page := range pages {
		filename := config.OutputDir + "/" + page.Filename
		file, err := os.OpenFile(filename, config.FileFlags, config.FileMode)
		if err != nil {
			return fmt.Errorf("could not open file %s: %v", filename, err)
		}

		styleTags, scriptTags, err := inlineAssets(page, config.SrcDir, config.ArticleDir)
		if err != nil {
			return fmt.Errorf("failed to inline assets for page %s: %v", page.Filename, err)
		}

		page.Render(styleTags, scriptTags).Render(context.Background(), file)
	}

	return nil
}

// postProcessing runs tasks that depend on all pages already being written:
// sitemap generation.
func postProcessing(config Config, allArticles []Article) error {
	if err := generateSitemap(config.OutputDir, config.BaseURL, allArticles); err != nil {
		log.Printf("Warning: Failed to generate sitemap: %v", err)
	}

	return nil
}

func main() {
	config := LoadConfig()
	InitMinifier()

	if err := publishGlobalCSS(config); err != nil {
		log.Fatalf("Error publishing global CSS: %v", err)
	}

	experiences, err := loadExperiencesFromJSON(config.SrcDir + "/experiences.json")
	if err != nil {
		log.Fatalf("Error loading experiences: %v", err)
	}

	allArticles, err := parseArticles(config.ArticleDir, config.Env)
	if err != nil {
		log.Fatalf("Error loading articles: %v", err)
	}

	if err := generateAllPages(config, allArticles, *experiences); err != nil {
		log.Fatalf("Error generating pages: %v", err)
	}

	if err := postProcessing(config, allArticles); err != nil {
		log.Fatalf("Error in post-processing: %v", err)
	}
}
