package main

import (
	"io/fs"
	"log"
	"os"
)

// Config holds all runtime configuration for the site generator,
// derived from environment variables and CLI flags.
type Config struct {
	ArticleDir  string
	OutputDir   string
	SrcDir      string
	Env         string
	BaseURL     string
	GeneratePDF bool
	CSSFiles    []string
	FileFlags   int
	FileMode    fs.FileMode
}

// LoadConfig reads environment variables and the generatePDF flag to
// produce a Config. Exits fatally if required variables are missing.
func LoadConfig(generatePDF bool) Config {
	articleDir := os.Getenv("ARTICLE_DIR")
	outputDir := os.Getenv("OUTPUT_DIR")
	srcDir := os.Getenv("SRC_DIR")

	if articleDir == "" || outputDir == "" || srcDir == "" {
		log.Fatal("ARTICLE_DIR, OUTPUT_DIR and SRC_DIR must be set")
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "production"
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "https://blog.ade-sede.dev"
	}

	cssFiles := []string{
		"global.css",
		"footer.css",
		"icons.css",
		"navbar.css",
		"syntax-highlighting.css",
	}

	return Config{
		ArticleDir:  articleDir,
		OutputDir:   outputDir,
		SrcDir:      srcDir,
		Env:         env,
		BaseURL:     baseURL,
		GeneratePDF: generatePDF,
		CSSFiles:    cssFiles,
		FileFlags:   os.O_RDWR | os.O_CREATE,
		FileMode:    0644,
	}
}
