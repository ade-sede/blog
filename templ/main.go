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
	outputDir := os.Getenv("OUTPUT_DIR")
	indexFileName := fmt.Sprintf("%s/index.html", outputDir)
	aboutFileName := fmt.Sprintf("%s/about.html", outputDir)
	indexFile, err := os.OpenFile(indexFileName, FLAGS, MODE)
	if err != nil {
		log.Fatal("Could not open file: %w", err)
	}

	aboutFile, err := os.OpenFile(aboutFileName, FLAGS, MODE)
	if err != nil {
		log.Fatal("Could not open file: %w", err)
	}

	home := home(make(map[string]string))
	home.Render(context.Background(), indexFile)

	about := about(make(map[string]string))
	about.Render(context.Background(), aboutFile)
}
