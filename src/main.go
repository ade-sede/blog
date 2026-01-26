package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/a-h/templ"
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

func initializeAssets(config Config) error {
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

func loadArticles(config Config) ([]Article, error) {
	return parseArticles(config.ArticleDir, config.Env)
}

func loadExperiences(config Config) (ExperiencesData, error) {
	exp, err := loadExperiencesFromJSON(config.SrcDir + "/experiences.json")
	if err != nil {
		return ExperiencesData{}, err
	}
	return *exp, nil
}

func generateAllPages(config Config, allArticles []Article, experiences ExperiencesData) error {
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
			arguments:                  []interface{}{experiences.WorkExperiences, experiences.SchoolExperiences},
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
			arguments:                  []interface{}{experiences.WorkExperiences, experiences.SchoolExperiences},
		},
	}

	pages = append(pages, allArticleGenerators...)

	for _, page := range pages {
		filename := config.OutputDir + "/" + page.filename
		file, err := os.OpenFile(filename, config.FileFlags, config.FileMode)
		if err != nil {
			return fmt.Errorf("could not open file %s: %v", filename, err)
		}

		allStyleTags, allScriptTags, err := createInlineStyleAndScriptTags(page, config.SrcDir, config.ArticleDir)
		if err != nil {
			return fmt.Errorf("failed to create inline style and script tags for page %s: %v", page.filename, err)
		}

		page.arguments = append(page.arguments, allStyleTags)
		page.arguments = append(page.arguments, allScriptTags)

		templComponent := page.HTMLgenerator(page.arguments...)
		templComponent.Render(context.Background(), file)
	}

	return nil
}

func postProcessing(config Config, allArticles []Article) error {
	if err := generateSitemap(config.OutputDir, config.BaseURL, allArticles); err != nil {
		log.Printf("Warning: Failed to generate sitemap: %v", err)
	}

	if config.GeneratePDF {
		if err := generatePDF(config.OutputDir, config.SrcDir); err != nil {
			log.Printf("Warning: Failed to generate PDF: %v", err)
		}
	}

	return nil
}

func main() {
	generatePDFFlag := flag.Bool("pdf", false, "Generate PDF after building site")
	flag.Parse()

	config := LoadConfig(*generatePDFFlag)
	InitMinifier()

	if err := initializeAssets(config); err != nil {
		log.Fatalf("Error initializing assets: %v", err)
	}

	experiences, err := loadExperiences(config)
	if err != nil {
		log.Fatalf("Error loading experiences: %v", err)
	}

	allArticles, err := loadArticles(config)
	if err != nil {
		log.Fatalf("Error loading articles: %v", err)
	}

	if err := generateAllPages(config, allArticles, experiences); err != nil {
		log.Fatalf("Error generating pages: %v", err)
	}

	if err := postProcessing(config, allArticles); err != nil {
		log.Fatalf("Error in post-processing: %v", err)
	}
}
