package parser

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// HTMLParser handles HTML document parsing and data extraction
type HTMLParser struct{}

func NewHTMLParser() *HTMLParser {
	return &HTMLParser{}
}

// ExtractMetaTags extracts all meta tags from HTML
func (p *HTMLParser) ExtractMetaTags(doc *goquery.Document) map[string]string {
	metaTags := make(map[string]string)

	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if name, exists := s.Attr("name"); exists {
			content, _ := s.Attr("content")
			metaTags[name] = content
		}

		if property, exists := s.Attr("property"); exists {
			content, _ := s.Attr("content")
			metaTags[property] = content
		}
	})

	return metaTags
}

// ExtractLinks extracts all internal and external links
func (p *HTMLParser) ExtractLinks(doc *goquery.Document, baseURL string) ([]string, []string) {
	var internal, external []string

	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if href == "" {
			return
		}

		// Simple check for internal vs external
		if strings.HasPrefix(href, "/") || strings.Contains(href, baseURL) {
			internal = append(internal, href)
		} else if strings.HasPrefix(href, "http") {
			external = append(external, href)
		}
	})

	return internal, external
}

// ExtractHeadings extracts all heading tags (H1-H6)
func (p *HTMLParser) ExtractHeadings(doc *goquery.Document) map[string][]string {
	headings := make(map[string][]string)

	for i := 1; i <= 6; i++ {
		tag := "h" + string(rune('0'+i))
		var texts []string

		doc.Find(tag).Each(func(j int, s *goquery.Selection) {
			texts = append(texts, strings.TrimSpace(s.Text()))
		})

		if len(texts) > 0 {
			headings[tag] = texts
		}
	}

	return headings
}

// ExtractImages extracts all image URLs and alt texts
func (p *HTMLParser) ExtractImages(doc *goquery.Document) []map[string]string {
	var images []map[string]string

	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		alt, _ := s.Attr("alt")

		images = append(images, map[string]string{
			"src": src,
			"alt": alt,
		})
	})

	return images
}

// CheckRobotsMetaTag checks for robots meta tag
func (p *HTMLParser) CheckRobotsMetaTag(doc *goquery.Document) (bool, string) {
	var robotsContent string
	found := false

	doc.Find("meta[name='robots']").Each(func(i int, s *goquery.Selection) {
		content, exists := s.Attr("content")
		if exists {
			found = true
			robotsContent = content
		}
	})

	return found, robotsContent
}

// GetCanonicalURL extracts canonical URL if present
func (p *HTMLParser) GetCanonicalURL(doc *goquery.Document) string {
	canonical, _ := doc.Find("link[rel='canonical']").Attr("href")
	return canonical
}
