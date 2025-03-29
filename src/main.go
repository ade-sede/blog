package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/a-h/templ"
)

type PageGenerator struct {
	filename      string
	HTMLgenerator func(args ...interface{}) templ.Component
	cssFilename   []string
	arguments     []interface{}
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

func main() {
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
		styleTags, _ := args[3].([]string)
		return article(title, description, stringifiedHTML, styleTags)
	}

	quickNoteHTMLGenerator := func(args ...interface{}) templ.Component {
		title, _ := args[0].(string)
		description, _ := args[1].(string)
		stringifiedHTML, _ := args[2].(string)
		styleTags, _ := args[3].([]string)
		return quickNote(title, description, stringifiedHTML, styleTags)
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
		allArticlesHTMLGenerators = append(allArticlesHTMLGenerators, PageGenerator{
			filename:      a.HTMLFilename,
			HTMLgenerator: articleHTMLGenerator,
			cssFilename:   []string{"article.css", "syntax-highlighting.css"},
			arguments:     []interface{}{a.Manifest.Title, a.Manifest.Description, a.StringifiedHTML},
		})
	}

	for _, a := range allQuickNotes {
		allQuickNotesHTMLGenerators = append(allQuickNotesHTMLGenerators, PageGenerator{
			filename:      a.HTMLFilename,
			HTMLgenerator: quickNoteHTMLGenerator,
			cssFilename:   []string{"article.css", "syntax-highlighting.css"},
			arguments:     []interface{}{a.Manifest.Title, a.Manifest.Description, a.StringifiedHTML},
		})
	}

	pages := []PageGenerator{
		{
			filename:      "index.html",
			HTMLgenerator: homeHTMLGenerator,
			cssFilename:   []string{"home.css", "articles.css"},
			arguments:     []interface{}{allArticles, allQuickNotes},
		},
		{
			filename:      "resume.html",
			HTMLgenerator: resumePageHTMLGenerator,
			cssFilename:   []string{"resume.css"},
			arguments:     []interface{}{workExperiences, schoolExperiences},
		},
		{
			filename:      "articles.html",
			HTMLgenerator: articlesHTMLGenerator,
			cssFilename:   []string{"articles.css"},
			arguments:     []interface{}{allArticles},
		},
		{
			filename:      "quick-notes.html",
			HTMLgenerator: quickNotesHTMLGenerator,
			cssFilename:   []string{"articles.css"},
			arguments:     []interface{}{allQuickNotes},
		},
		// Used for PDF generation
		{
			filename:      "resume-printable.html",
			HTMLgenerator: resumePrintReadyHTMLGenerator,
			cssFilename:   []string{"resume.css"},
			arguments:     []interface{}{workExperiences, schoolExperiences},
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

		allStyleTags := make([]string, 0)

		for _, filename := range page.cssFilename {
			filename := srcDir + "/css/" + filename
			minifiedCSS, err := MinifyCSS(filename)

			if err != nil {
				log.Fatalf("%v", err)
			}

			styleTag := fmt.Sprintf("<style type='text/css'>%s</style>", minifiedCSS)

			allStyleTags = append(allStyleTags, styleTag)
		}

		if len(allStyleTags) > 0 {
			page.arguments = append(page.arguments, allStyleTags)
		}

		templComponent := page.HTMLgenerator(page.arguments...)
		templComponent.Render(context.Background(), file)
	}
}

// For some reason I can't substitute within <style></style> directly so I'm trying to bypass the issue by wrapping in a function
func makeStyleTag(css string) templ.Component {
	return templ.Raw(fmt.Sprintf("<style type='text/css'>%s</style>", css))
}
