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

func init() {
	m = minify.New()
	m.AddFunc("text/css", cssminifier.Minify)
	m.AddFunc("application/javascript", jsminifier.Minify)
}

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

func loadAndMinifyGlobalStyle(srcDir, filename string) (string, error) {
	path := filepath.Join(srcDir, "css", filename)
	return MinifyCSS(path)
}

func loadAndMinifyGlobalScript(srcDir, filename string) (string, error) {
	path := filepath.Join(srcDir, "scripts", filename)
	return MinifyJS(path)
}

func loadAndMinifyArticleStyle(articleDir, filename string) (string, error) {
	searchPaths := []string{
		filepath.Join(articleDir, filename),
	}
	return loadAndMinifyFileFromPaths(searchPaths, MinifyCSS)
}

func loadAndMinifyArticleScript(articleDir, filename string) (string, error) {
	searchPaths := []string{
		filepath.Join(articleDir, filename),
	}
	return loadAndMinifyFileFromPaths(searchPaths, MinifyJS)
}
