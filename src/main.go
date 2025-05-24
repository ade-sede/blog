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
	PageTypeGlobal    PageType = "global"
	PageTypeArticle   PageType = "article"
	PageTypeQuickNote PageType = "quicknote"
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
		globalScriptFiles:          []string{"toc.js"},
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

func createInlineStyleAndScriptTags(page PageGenerator, srcDir, articleDir, quickNoteDir string) ([]string, []string, error) {
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
		var minifiedJS string
		var err error

		if page.pageType == PageTypeQuickNote {
			minifiedJS, err = loadAndMinifyQuickNoteScript(quickNoteDir, jsFile)
		} else {
			minifiedJS, err = loadAndMinifyArticleScript(articleDir, jsFile)
		}

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
		var minifiedCSS string
		var err error

		if page.pageType == PageTypeQuickNote {
			minifiedCSS, err = loadAndMinifyQuickNoteStyle(quickNoteDir, cssFile)
		} else {
			minifiedCSS, err = loadAndMinifyArticleStyle(articleDir, cssFile)
		}

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

	quickNoteDir := os.Getenv("QUICK_NOTE_DIR")
	articleDir := os.Getenv("ARTICLE_DIR")
	outputDir := os.Getenv("OUTPUT_DIR")
	srcDir := os.Getenv("SRC_DIR")

	if articleDir == "" || outputDir == "" || srcDir == "" || quickNoteDir == "" {
		log.Fatal("ARTICLE_DIR, QUICK_NOTE_DIR, OUTPUT_DIR and SRC_DIR must be set")
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

	allQuickNotes, err := parseArticles(quickNoteDir)
	if err != nil {
		log.Fatalf("Error while parsing quick notes: %v", err)
	}

	homeHTMLGenerator := func(args ...interface{}) templ.Component {
		allArticles, _ := args[0].([]Article)
		allQuickNotes, _ := args[1].([]Article)
		styleTags, _ := args[2].([]string)

		if len(allQuickNotes) > 2 {
			allQuickNotes = allQuickNotes[:2]
		}

		if len(allArticles) > 2 {
			allArticles = allArticles[:2]
		}

		return home(allArticles, allQuickNotes, styleTags)
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

	quickNoteHTMLGenerator := func(args ...interface{}) templ.Component {
		title, _ := args[0].(string)
		description, _ := args[1].(string)
		stringifiedHTML, _ := args[2].(string)
		formattedDate, _ := args[3].(string)
		toc, _ := args[4].([]TOCEntry)
		styleTags, _ := args[5].([]string)
		scriptTags, _ := args[6].([]string)
		return quickNote(title, description, stringifiedHTML, formattedDate, scriptTags, styleTags, toc)
	}

	articlesHTMLGenerator := func(args ...interface{}) templ.Component {
		allArticles, _ := args[0].([]Article)
		styleTags, _ := args[1].([]string)
		return articles(allArticles, styleTags)
	}

	quickNotesHTMLGenerator := func(args ...interface{}) templ.Component {
		allQuickNotes, _ := args[0].([]Article)
		styleTags, _ := args[1].([]string)
		return quickNotes(allQuickNotes, styleTags)
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

	allArticlesHTMLGenerators := make([]PageGenerator, 0)
	allQuickNotesHTMLGenerators := make([]PageGenerator, 0)

	for _, a := range allArticles {
		generator := createContentPageGenerator(a, PageTypeArticle, articleHTMLGenerator)
		allArticlesHTMLGenerators = append(allArticlesHTMLGenerators, generator)
	}

	for _, a := range allQuickNotes {
		generator := createContentPageGenerator(a, PageTypeQuickNote, quickNoteHTMLGenerator)
		allQuickNotesHTMLGenerators = append(allQuickNotesHTMLGenerators, generator)
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
			arguments:                  []interface{}{allArticles, allQuickNotes},
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
			filename:                   "quick-notes.html",
			pageType:                   PageTypeGlobal,
			HTMLgenerator:              quickNotesHTMLGenerator,
			globalCssFiles:             []string{"articles.css"},
			globalScriptFiles:          []string{},
			articleSpecificCssFiles:    []string{},
			articleSpecificScriptFiles: []string{},
			arguments:                  []interface{}{allQuickNotes},
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

	pages = append(pages, allArticlesHTMLGenerators...)
	pages = append(pages, allQuickNotesHTMLGenerators...)

	for _, page := range pages {
		filename := outputDir + "/" + page.filename
		file, err := os.OpenFile(filename, FLAGS, MODE)
		if err != nil {
			log.Fatalf("Could not open file: %v", err)
		}

		allStyleTags, allScriptTags, err := createInlineStyleAndScriptTags(page, srcDir, articleDir, quickNoteDir)
		if err != nil {
			log.Fatalf("Failed to create inline style and script tags for page %s: %v", page.filename, err)
		}

		page.arguments = append(page.arguments, allStyleTags)
		page.arguments = append(page.arguments, allScriptTags)

		templComponent := page.HTMLgenerator(page.arguments...)
		templComponent.Render(context.Background(), file)
	}

	if *generatePDFFlag {
		if err := generatePDF(outputDir, srcDir); err != nil {
			log.Printf("Warning: Failed to generate PDF: %v", err)
		}
	}
}
