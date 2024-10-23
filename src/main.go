package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
)

var FLAGS = os.O_RDWR | os.O_CREATE
var MODE fs.FileMode = 0644

func main() {
	articleDir := os.Getenv("ARTICLE_DIR")
	outputDir := os.Getenv("OUTPUT_DIR")

	if articleDir == "" {
		log.Fatal("ARTICLE_DIR is not set")
	}

	if outputDir == "" {
		log.Fatal("OUTPUT_DIR is not set")
	}

	articles, err := parseArticles(articleDir)
	if err != nil {
		log.Fatalf("Error while parsing articles: %v", err)
	}

	for _, a := range articles {
		filename := fmt.Sprintf("%s/%s", outputDir, a.HTMLFilename)
		file, err := os.OpenFile(filename, FLAGS, MODE)
		if err != nil {
			log.Fatalf("Could not open file: %v", err)
		}

		aa := article(a.Manifest.Title, a.Date, a.StringifiedHTML)
		aa.Render(context.Background(), file)
	}

	indexFilename := fmt.Sprintf("%s/index.html", outputDir)
	indexFile, err := os.OpenFile(indexFilename, FLAGS, MODE)
	if err != nil {
		log.Fatalf("Could not open file: %v", err)
	}

	aboutFilename := fmt.Sprintf("%s/about.html", outputDir)
	aboutFile, err := os.OpenFile(aboutFilename, FLAGS, MODE)
	if err != nil {
		log.Fatalf("Could not open file: %v", err)
	}

	resumeLightFilename := fmt.Sprintf("%s/resume-light.html", outputDir)
	resumeLightFile, err := os.OpenFile(resumeLightFilename, FLAGS, MODE)
	if err != nil {
		log.Fatalf("Could not open file: %v", err)
	}

	home := home()
	home.Render(context.Background(), indexFile)

	workExperiences := buildWorkExperiences()
	schoolExperiences := buildSchoolExperience()

	about := about(workExperiences, schoolExperiences)
	about.Render(context.Background(), aboutFile)

	// Print resume to a standalone HTMl file so that we can easily create a PDF from it
	// Will most likely be removed before serving the content
	resumeLight := resumeLight(workExperiences, schoolExperiences)
	resumeLight.Render(context.Background(), resumeLightFile)
}
