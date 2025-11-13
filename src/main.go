package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/a-h/templ"
	"io/fs"
	"log"
	"os"
)

type PageType string

const (
	PageTypeGlobal  PageType = "global"
	PageTypeArticle PageType = "article"
)

type PageGenerator struct {
	filename                   string
	pageType                   PageType
	HTMLgenerator              func(args ...interface{}) templ.Component
	globalCssFiles             []string
	globalScriptFiles          []string
	articleSpecificCssFiles    []string
	articleSpecificScriptFiles []string
	arguments                  []interface{}
}

// List of 'global' css files that need to be minified and shipped.
// If not in the list, it will not be included.
// For page specific CSS prefer inline CSS.
var CSSFiles []string = []string{
	"global.css",
	"footer.css",
	"icons.css",
	"navbar.css",
	"syntax-highlighting.css",
}

var FLAGS = os.O_RDWR | os.O_CREATE
var MODE fs.FileMode = 0644

func createContentPageGenerator(article Article, pageType PageType, htmlGenerator func(args ...interface{}) templ.Component) PageGenerator {
	generator := PageGenerator{
		filename:                   article.HTMLFilename,
		pageType:                   pageType,
		HTMLgenerator:              htmlGenerator,
		globalCssFiles:             []string{"article.css", "syntax-highlighting.css"},
		globalScriptFiles:          []string{"toc.js", "anchors.js", "footnotes.js"},
		articleSpecificCssFiles:    []string{},
		articleSpecificScriptFiles: []string{},
		arguments:                  []interface{}{article.Manifest.Title, article.Manifest.Description, article.StringifiedHTML, article.FormatedDate, article.TOC},
	}

	if article.Manifest.CssFile != "" {
		generator.articleSpecificCssFiles = append(generator.articleSpecificCssFiles, article.Manifest.CssFile)
	}

	if article.Manifest.ScriptFile != "" {
		generator.articleSpecificScriptFiles = append(generator.articleSpecificScriptFiles, article.Manifest.ScriptFile)
	}

	return generator
}

func createInlineStyleAndScriptTags(page PageGenerator, srcDir, articleDir string) ([]string, []string, error) {
	allStyleTags := make([]string, 0)
	allScriptTags := make([]string, 0)

	for _, jsFile := range page.globalScriptFiles {
		minifiedJS, err := loadAndMinifyGlobalScript(srcDir, jsFile)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load global JS file %s: %v", jsFile, err)
		}
		scriptTag := fmt.Sprintf("<script>%s</script>", minifiedJS)
		allScriptTags = append(allScriptTags, scriptTag)
	}

	for _, jsFile := range page.articleSpecificScriptFiles {
		minifiedJS, err := loadAndMinifyArticleScript(articleDir, jsFile)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load article-specific JS file %s: %v", jsFile, err)
		}
		scriptTag := fmt.Sprintf("<script>%s</script>", minifiedJS)
		allScriptTags = append(allScriptTags, scriptTag)
	}

	for _, cssFile := range page.globalCssFiles {
		minifiedCSS, err := loadAndMinifyGlobalStyle(srcDir, cssFile)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load global CSS file %s: %v", cssFile, err)
		}
		styleTag := fmt.Sprintf("<style type='text/css'>%s</style>", minifiedCSS)
		allStyleTags = append(allStyleTags, styleTag)
	}

	for _, cssFile := range page.articleSpecificCssFiles {
		minifiedCSS, err := loadAndMinifyArticleStyle(articleDir, cssFile)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load article-specific CSS file %s: %v", cssFile, err)
		}
		styleTag := fmt.Sprintf("<style type='text/css'>%s</style>", minifiedCSS)
		allStyleTags = append(allStyleTags, styleTag)
	}

	return allStyleTags, allScriptTags, nil
}

func main() {
	generatePDFFlag := flag.Bool("pdf", false, "Generate PDF after building site")
	flag.Parse()

	articleDir := os.Getenv("ARTICLE_DIR")
	outputDir := os.Getenv("OUTPUT_DIR")
	srcDir := os.Getenv("SRC_DIR")

	if articleDir == "" || outputDir == "" || srcDir == "" {
		log.Fatal("ARTICLE_DIR, OUTPUT_DIR and SRC_DIR must be set")
	}

	experiences, err := loadExperiencesFromJSON(srcDir + "/experiences.json")
	if err != nil {
		log.Fatalf("Error while generating experiences: %v", err)
	}

	workExperiences := experiences.WorkExperiences
	schoolExperiences := experiences.SchoolExperiences

	for _, filename := range CSSFiles {
		minifiedCSS, err := MinifyCSS(srcDir + "/css/" + filename)
		if err != nil {
			log.Fatalf("%v", err)
		}

		newFilename := outputDir + "/css/" + filename
		os.WriteFile(newFilename, []byte(minifiedCSS), MODE)
	}

	allArticles, err := parseArticles(articleDir)
	if err != nil {
		log.Fatalf("Error while parsing articles: %v", err)
	}

	homeHTMLGenerator := func(args ...interface{}) templ.Component {
		allArticles, _ := args[0].([]Article)
		styleTags, _ := args[1].([]string)

		var mostRecent []Article
		if len(allArticles) > 0 {
			mostRecent = allArticles[:1]
		}

		return home(mostRecent, styleTags)
	}

	articleHTMLGenerator := func(args ...interface{}) templ.Component {
		title, _ := args[0].(string)
		description, _ := args[1].(string)
		stringifiedHTML, _ := args[2].(string)
		formattedDate, _ := args[3].(string)
		toc, _ := args[4].([]TOCEntry)
		styleTags, _ := args[5].([]string)
		scriptTags, _ := args[6].([]string)
		return article(title, description, stringifiedHTML, formattedDate, scriptTags, styleTags, toc)
	}

	articlesHTMLGenerator := func(args ...interface{}) templ.Component {
		allArticles, _ := args[0].([]Article)
		styleTags, _ := args[1].([]string)
		return articles(allArticles, styleTags)
	}

	resumePageHTMLGenerator := func(args ...interface{}) templ.Component {
		workExperience, _ := args[0].([]ExperienceEntry)
		schoolExperience, _ := args[1].([]ExperienceEntry)
		styleTags, _ := args[2].([]string)
		return resumePage(workExperience, schoolExperience, styleTags)
	}

	resumePrintReadyHTMLGenerator := func(args ...interface{}) templ.Component {
		workExperience, _ := args[0].([]ExperienceEntry)
		schoolExperience, _ := args[1].([]ExperienceEntry)
		styleTags, _ := args[2].([]string)
		return resumePrintReady(workExperience, schoolExperience, styleTags)
	}

	allArticleGenerators := make([]PageGenerator, 0)

	for _, a := range allArticles {
		generator := createContentPageGenerator(a, PageTypeArticle, articleHTMLGenerator)
		allArticleGenerators = append(allArticleGenerators, generator)
	}

	pages := []PageGenerator{
		{
			filename:                   "index.html",
			pageType:                   PageTypeGlobal,
			HTMLgenerator:              homeHTMLGenerator,
			globalCssFiles:             []string{"home.css", "articles.css"},
			globalScriptFiles:          []string{},
			articleSpecificCssFiles:    []string{},
			articleSpecificScriptFiles: []string{},
			arguments:                  []interface{}{allArticles},
		},
		{
			filename:                   "resume.html",
			pageType:                   PageTypeGlobal,
			HTMLgenerator:              resumePageHTMLGenerator,
			globalCssFiles:             []string{"resume.css"},
			globalScriptFiles:          []string{},
			articleSpecificCssFiles:    []string{},
			articleSpecificScriptFiles: []string{},
			arguments:                  []interface{}{workExperiences, schoolExperiences},
		},
		{
			filename:                   "articles.html",
			pageType:                   PageTypeGlobal,
			HTMLgenerator:              articlesHTMLGenerator,
			globalCssFiles:             []string{"articles.css"},
			globalScriptFiles:          []string{},
			articleSpecificCssFiles:    []string{},
			articleSpecificScriptFiles: []string{},
			arguments:                  []interface{}{allArticles},
		},
		{
			filename:                   "resume-printable.html",
			pageType:                   PageTypeGlobal,
			HTMLgenerator:              resumePrintReadyHTMLGenerator,
			globalCssFiles:             []string{"resume.css"},
			globalScriptFiles:          []string{},
			articleSpecificCssFiles:    []string{},
			articleSpecificScriptFiles: []string{},
			arguments:                  []interface{}{workExperiences, schoolExperiences},
		},
	}

	pages = append(pages, allArticleGenerators...)

	for _, page := range pages {
		filename := outputDir + "/" + page.filename
		file, err := os.OpenFile(filename, FLAGS, MODE)
		if err != nil {
			log.Fatalf("Could not open file: %v", err)
		}

		allStyleTags, allScriptTags, err := createInlineStyleAndScriptTags(page, srcDir, articleDir)
		if err != nil {
			log.Fatalf("Failed to create inline style and script tags for page %s: %v", page.filename, err)
		}

		page.arguments = append(page.arguments, allStyleTags)
		page.arguments = append(page.arguments, allScriptTags)

		templComponent := page.HTMLgenerator(page.arguments...)
		templComponent.Render(context.Background(), file)
	}

	if err := generateSitemap(outputDir, allArticles); err != nil {
		log.Printf("Warning: Failed to generate sitemap: %v", err)
	}

	if *generatePDFFlag {
		if err := generatePDF(outputDir, srcDir); err != nil {
			log.Printf("Warning: Failed to generate PDF: %v", err)
		}
	}
}
