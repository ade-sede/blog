package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"
)

// URLSet is the root element of a sitemap XML document.
type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	URLs    []URL    `xml:"url"`
}

// URL represents a single entry in the sitemap.
type URL struct {
	Loc        string  `xml:"loc"`
	LastMod    string  `xml:"lastmod,omitempty"`
	ChangeFreq string  `xml:"changefreq,omitempty"`
	Priority   float64 `xml:"priority,omitempty"`
}

// generateSitemap writes a sitemap.xml to outputDir covering all static pages
// and published articles. Draft articles are excluded.
func generateSitemap(outputDir, baseURL string, allArticles []Article) error {

	urlset := URLSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  []URL{},
	}

	now := time.Now().Format("2006-01-02")

	urlset.URLs = append(urlset.URLs, URL{
		Loc:        baseURL + "/",
		LastMod:    now,
		ChangeFreq: "weekly",
		Priority:   1.0,
	})

	urlset.URLs = append(urlset.URLs, URL{
		Loc:        baseURL + "/articles.html",
		LastMod:    now,
		ChangeFreq: "weekly",
		Priority:   0.8,
	})

	urlset.URLs = append(urlset.URLs, URL{
		Loc:        baseURL + "/resume.html",
		LastMod:    now,
		ChangeFreq: "monthly",
		Priority:   0.6,
	})

	for _, article := range allArticles {
		if article.Manifest.Draft {
			continue
		}
		urlset.URLs = append(urlset.URLs, URL{
			Loc:        baseURL + "/" + article.HTMLFilename,
			LastMod:    article.Date.Format("2006-01-02"),
			ChangeFreq: "monthly",
			Priority:   0.9,
		})
	}

	xmlData, err := xml.MarshalIndent(urlset, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal sitemap XML: %v", err)
	}

	xmlHeader := []byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	fullXML := append(xmlHeader, xmlData...)

	sitemapPath := outputDir + "/sitemap.xml"
	err = os.WriteFile(sitemapPath, fullXML, 0644)
	if err != nil {
		return fmt.Errorf("failed to write sitemap file: %v", err)
	}

	return nil
}
