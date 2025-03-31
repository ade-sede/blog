package main

import (
	"context"
	"fmt"
	"github.com/a-h/templ"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type PageGenerator struct {
	filename       string
	HTMLgenerator  func(args ...interface{}) templ.Component
	cssFilename    []string
	scriptFilename []string
	arguments      []interface{}
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

func findAndMinifyFile(baseDirs []string, subDir, filename string, minifyFunc func(string) (string, error)) (string, error) {
	// This function is a bit of a monster.
	// I have hijacked existing functionalities to inject CSS files to inject javascript and css for articles
	// Because I need a lot of flexibility to accomodate many formats / locations I just kind of brute force it
	// Try every place it could be

	for _, baseDir := range baseDirs[:1] {
		path := filepath.Join(baseDir, subDir, filename)
		result, err := minifyFunc(path)
		if err == nil {
			return result, nil
		}
		if !os.IsNotExist(err) {
			return "", err
		}
	}

	for _, baseDir := range baseDirs[1:] {
		path := filepath.Join(baseDir, filename)
		result, err := minifyFunc(path)
		if err == nil {
			return result, nil
		}
		if !os.IsNotExist(err) {
			return "", err
		}
	}

	return "", fmt.Errorf("file %s not found in any location", filename)
}

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
		formattedDate, _ := args[3].(string)
		styleTags, _ := args[4].([]string)
		scriptTags, _ := args[5].([]string)
		return article(title, description, stringifiedHTML, formattedDate, styleTags, scriptTags)
	}

	quickNoteHTMLGenerator := func(args ...interface{}) templ.Component {
		title, _ := args[0].(string)
		description, _ := args[1].(string)
		stringifiedHTML, _ := args[2].(string)
		formattedDate, _ := args[3].(string)
		styleTags, _ := args[4].([]string)
		scriptTags, _ := args[5].([]string)
		return quickNote(title, description, stringifiedHTML, formattedDate, styleTags, scriptTags)
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
		generator := PageGenerator{
			filename:       a.HTMLFilename,
			HTMLgenerator:  articleHTMLGenerator,
			cssFilename:    []string{"article.css", "syntax-highlighting.css"},
			scriptFilename: []string{},
			arguments:      []interface{}{a.Manifest.Title, a.Manifest.Description, a.StringifiedHTML, a.FormatedDate},
		}

		if a.Manifest.CssFile != "" {
			generator.cssFilename = append(generator.cssFilename, a.Manifest.CssFile)
		}

		if a.Manifest.ScriptFile != "" {
			generator.scriptFilename = append(generator.scriptFilename, a.Manifest.ScriptFile)
		}
		allArticlesHTMLGenerators = append(allArticlesHTMLGenerators, generator)
	}

	for _, a := range allQuickNotes {
		generator := PageGenerator{
			filename:       a.HTMLFilename,
			HTMLgenerator:  quickNoteHTMLGenerator,
			cssFilename:    []string{"article.css", "syntax-highlighting.css"},
			scriptFilename: []string{},
			arguments:      []interface{}{a.Manifest.Title, a.Manifest.Description, a.StringifiedHTML, a.FormatedDate},
		}

		if a.Manifest.CssFile != "" {
			generator.cssFilename = append(generator.cssFilename, a.Manifest.CssFile)
		}

		if a.Manifest.ScriptFile != "" {
			generator.scriptFilename = append(generator.scriptFilename, a.Manifest.ScriptFile)
		}
		allQuickNotesHTMLGenerators = append(allQuickNotesHTMLGenerators, generator)
	}

	pages := []PageGenerator{
		{
			filename:       "index.html",
			HTMLgenerator:  homeHTMLGenerator,
			cssFilename:    []string{"home.css", "articles.css"},
			scriptFilename: []string{},
			arguments:      []interface{}{allArticles, allQuickNotes},
		},
		{
			filename:       "resume.html",
			HTMLgenerator:  resumePageHTMLGenerator,
			cssFilename:    []string{"resume.css"},
			scriptFilename: []string{},
			arguments:      []interface{}{workExperiences, schoolExperiences},
		},
		{
			filename:       "articles.html",
			HTMLgenerator:  articlesHTMLGenerator,
			cssFilename:    []string{"articles.css"},
			scriptFilename: []string{},
			arguments:      []interface{}{allArticles},
		},
		{
			filename:       "quick-notes.html",
			HTMLgenerator:  quickNotesHTMLGenerator,
			cssFilename:    []string{"articles.css"},
			scriptFilename: []string{},
			arguments:      []interface{}{allQuickNotes},
		},
		{
			filename:       "resume-printable.html",
			HTMLgenerator:  resumePrintReadyHTMLGenerator,
			cssFilename:    []string{"resume.css"},
			scriptFilename: []string{},
			arguments:      []interface{}{workExperiences, schoolExperiences},
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
		allScriptTags := make([]string, 0)

		for _, jsFile := range page.scriptFilename {
			minifiedJS, err := findAndMinifyFile([]string{srcDir, quickNoteDir, articleDir}, "scripts", jsFile, MinifyJS)
			if err != nil {
				log.Fatalf("Failed to find or minify JS file %s: %v", jsFile, err)
			}
			scriptTag := fmt.Sprintf("<script>%s</script>", minifiedJS)
			allScriptTags = append(allScriptTags, scriptTag)
		}

		for _, cssFile := range page.cssFilename {
			minifiedCSS, err := findAndMinifyFile([]string{srcDir, quickNoteDir, articleDir}, "css", cssFile, MinifyCSS)
			if err != nil {
				log.Fatalf("Failed to find or minify CSS file %s: %v", cssFile, err)
			}
			styleTag := fmt.Sprintf("<style type='text/css'>%s</style>", minifiedCSS)
			allStyleTags = append(allStyleTags, styleTag)
		}

		page.arguments = append(page.arguments, allStyleTags)
		page.arguments = append(page.arguments, allScriptTags)

		templComponent := page.HTMLgenerator(page.arguments...)
		templComponent.Render(context.Background(), file)
	}
}
