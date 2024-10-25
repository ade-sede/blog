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
}

var FLAGS = os.O_RDWR | os.O_CREATE
var MODE fs.FileMode = 0644

func main() {
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

		return home(allArticles, styleTags)
	}

	articleHTMLGenerator := func(args ...interface{}) templ.Component {
		title, _ := args[0].(string)
		description, _ := args[1].(string)
		stringifiedHTML, _ := args[2].(string)
		styleTags, _ := args[3].([]string)
		return article(title, description, stringifiedHTML, styleTags)
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

	allArticlesHTMLGenerators := make([]PageGenerator, 0)

	for _, a := range allArticles {
		allArticlesHTMLGenerators = append(allArticlesHTMLGenerators, PageGenerator{
			filename:      a.HTMLFilename,
			HTMLgenerator: articleHTMLGenerator,
			cssFilename:   []string{"article.css"},
			arguments:     []interface{}{a.Manifest.Title, a.Manifest.Description, a.StringifiedHTML},
		})
	}

	pages := []PageGenerator{
		{
			filename:      "index.html",
			HTMLgenerator: homeHTMLGenerator,
			cssFilename:   []string{"home.css", "articles.css"},
			arguments:     []interface{}{allArticles},
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
		// Used for PDF generation
		{
			filename:      "resume-printable.html",
			HTMLgenerator: resumePrintReadyHTMLGenerator,
			cssFilename:   []string{"resume.css"},
			arguments:     []interface{}{workExperiences, schoolExperiences},
		},
	}

	pages = append(pages, allArticlesHTMLGenerators...)

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
