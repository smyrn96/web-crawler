package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type CrawlResult struct {
	Title          string
	HTMLVersion    string
	HasLoginForm   bool
	InternalLinks  int
	ExternalLinks  int
	BrokenLinks    string // JSON-encoded list of broken link objects
}

type BrokenLink struct {
	URL        string `json:"url"`
	StatusCode int    `json:"status_code"`
}

func CrawlURL(rawURL string) (*CrawlResult, error) {
	resp, err := http.Get(rawURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, errors.New("non-OK response code: " + resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	result := &CrawlResult{}
	result.Title = doc.Find("title").First().Text()
	result.HTMLVersion = detectHTMLVersion(doc)

	internal := 0
	external := 0
	broken := []BrokenLink{}

	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists || href == "" {
			return
		}

		link, err := parsedURL.Parse(href)
		if err != nil {
			return
		}

		isInternal := link.Host == parsedURL.Host || link.Host == ""
		if isInternal {
			internal++
		} else {
			external++
		}

		// Check if link is broken
		status := checkLink(link.String())
		if status >= 400 {
			broken = append(broken, BrokenLink{
				URL:        link.String(),
				StatusCode: status,
			})
		}
	})

	result.InternalLinks = internal
	result.ExternalLinks = external

	// Check for login form
	result.HasLoginForm = doc.Find("input[type='password']").Length() > 0

	// Encode broken links
	brokenJSON, _ := json.Marshal(broken)
	result.BrokenLinks = string(brokenJSON)

	return result, nil
}

func checkLink(link string) int {
	resp, err := http.Head(link)
	if err != nil {
		return 500
	}
	defer resp.Body.Close()
	return resp.StatusCode
}

func detectHTMLVersion(doc *goquery.Document) string {
	html := strings.ToLower(doc.Text())
	if strings.Contains(html, "<!doctype html>") {
		return "HTML5"
	}
	return "Unknown"
}
