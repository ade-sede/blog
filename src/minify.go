package main

import (
	"github.com/tdewolff/minify/v2"
	// Using `css` identifier causes collisions everywhere
	cssminifier "github.com/tdewolff/minify/v2/css"
	jsminifier "github.com/tdewolff/minify/v2/js"
	"os"
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
