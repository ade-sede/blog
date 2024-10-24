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

	workExperiences := buildWorkExperiences()
	schoolExperiences := buildSchoolExperience()

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

	aboutGenerator := func() templ.Component {
		return about(workExperiences, schoolExperiences)
	}

	resumeGenerator := func() templ.Component {
		return resume(workExperiences, schoolExperiences)
	}

	articlePageGenerator := func() templ.Component {
		return articles(allArticles)
	}

	pages := []PageGenerator{
		{filename: "index.html", gen: homeGenerator},
		{filename: "about.html", gen: aboutGenerator},
		{filename: "articles.html", gen: articlePageGenerator},
		{filename: "resume-light.html", gen: resumeGenerator},
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
