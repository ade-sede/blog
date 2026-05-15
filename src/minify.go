package main

import (
	"fmt"
	"github.com/tdewolff/minify/v2"
	// Using `css` identifier causes collisions everywhere
	cssminifier "github.com/tdewolff/minify/v2/css"
	jsminifier "github.com/tdewolff/minify/v2/js"
	"os"
	"path/filepath"
)

var m *minify.M

// InitMinifier initialises the package-level minifier with CSS and JS handlers.
// Must be called before any MinifyCSS or MinifyJS calls.
func InitMinifier() {
	m = minify.New()
	m.AddFunc("text/css", cssminifier.Minify)
	m.AddFunc("application/javascript", jsminifier.Minify)
}

// MinifyCSS reads the file at filepath and returns its minified CSS content.
func MinifyCSS(filepath string) (string, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	minified, err := m.String("text/css", string(content))
	if err != nil {
		return string(content), err
	}
	return minified, nil
}

// MinifyJS reads the file at filepath and returns its minified JavaScript content.
func MinifyJS(filepath string) (string, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	minified, err := m.String("application/javascript", string(content))
	if err != nil {
		return string(content), err
	}
	return minified, nil
}

// loadAndMinifyFileFromPaths tries each path in order, returning the minified
// content of the first file found. Returns an error if none exist.
func loadAndMinifyFileFromPaths(searchPaths []string, minifyFunc func(string) (string, error)) (string, error) {
	for _, path := range searchPaths {
		result, err := minifyFunc(path)
		if err == nil {
			return result, nil
		}
		if !os.IsNotExist(err) {
			return "", err
		}
	}
	return "", fmt.Errorf("file not found in any of the provided locations")
}

// loadAndMinifyGlobalStyle minifies a CSS file from the srcDir/css directory.
func loadAndMinifyGlobalStyle(srcDir, filename string) (string, error) {
	path := filepath.Join(srcDir, "css", filename)
	return MinifyCSS(path)
}

// loadAndMinifyGlobalScript minifies a JS file from the srcDir/scripts directory.
func loadAndMinifyGlobalScript(srcDir, filename string) (string, error) {
	path := filepath.Join(srcDir, "scripts", filename)
	return MinifyJS(path)
}

// loadAndMinifyArticleStyle minifies a CSS file from the article directory.
func loadAndMinifyArticleStyle(articleDir, filename string) (string, error) {
	searchPaths := []string{
		filepath.Join(articleDir, filename),
	}
	return loadAndMinifyFileFromPaths(searchPaths, MinifyCSS)
}

// loadAndMinifyArticleScript minifies a JS file from the article directory.
func loadAndMinifyArticleScript(articleDir, filename string) (string, error) {
	searchPaths := []string{
		filepath.Join(articleDir, filename),
	}
	return loadAndMinifyFileFromPaths(searchPaths, MinifyJS)
}
