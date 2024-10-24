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
	filename string
	gen      func() templ.Component
}

var FLAGS = os.O_RDWR | os.O_CREATE
var MODE fs.FileMode = 0644

func main() {
	articleDir := os.Getenv("ARTICLE_DIR")
	outputDir := os.Getenv("OUTPUT_DIR")

	if articleDir == "" || outputDir == "" {
		log.Fatal("ARTICLE_DIR and OUTPUT_DIR must be set")
	}

	experiences, err := loadExperiencesFromJSON("experiences.json")
	if err != nil {
		log.Fatalf("Error while generating experiences: %v", err)
	}

	workExperiences := experiences.WorkExperiences
	schoolExperiences := experiences.SchoolExperiences

	allArticles, err := parseArticles(articleDir)
	if err != nil {
		log.Fatalf("Error while parsing articles: %v", err)
	}

	articleGenerators := make([]PageGenerator, 0)

	for _, a := range allArticles {
		generator := func() templ.Component {
			return article(a.Manifest.Title, a.Date, a.StringifiedHTML)
		}

		articleGenerators = append(articleGenerators, PageGenerator{
			filename: a.HTMLFilename,
			gen:      generator,
		})
	}

	homeGenerator := func() templ.Component {
		return home(allArticles)
	}

	resumeGenerator := func() templ.Component {
		return resumePage(workExperiences, schoolExperiences)
	}

	fullScreenResumeGenerator := func() templ.Component {
		return resume(workExperiences, schoolExperiences)
	}

	articlePageGenerator := func() templ.Component {
		return articles(allArticles)
	}

	pages := []PageGenerator{
		{filename: "index.html", gen: homeGenerator},
		{filename: "resume.html", gen: resumeGenerator},
		{filename: "articles.html", gen: articlePageGenerator},
		// Used for PDF generation
		{filename: "resume-light.html", gen: fullScreenResumeGenerator},
	}

	pages = append(pages, articleGenerators...)

	for _, page := range pages {
		filename := fmt.Sprintf("%s/%s", outputDir, page.filename)
		file, err := os.OpenFile(filename, FLAGS, MODE)
		if err != nil {
			log.Fatalf("Could not open file: %v", err)
		}

		templComponent := page.gen()
		templComponent.Render(context.Background(), file)
	}
}
