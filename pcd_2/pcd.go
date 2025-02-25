package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func isValidHTTPUrl(urlStr string) bool {
	// Ensure the URL starts with "http://" or "https://"
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		return false
	}

	// Parse the URL
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	// Check if the URL's scheme is "http" or "https"
	return parsedUrl.Scheme == "http" || parsedUrl.Scheme == "https"
}

func crawlPage(link string) []string {
	res, err := http.Get(link)
	if err != nil {
		log.Printf("Error fetching URL %s: %v", link, err)
		return []string{}
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Printf("Error parsing document for URL %s: %v", link, err)
		return []string{}
	}

	var links []string

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		linkRef := s.AttrOr("href", "")
		if isValidHTTPUrl(linkRef) {
			links = append(links, linkRef)
		}
	})

	return links
}

func printMap(m map[string][]string) {
	// Print the map in a structured format
	fmt.Println("Map contents:")
	for key, value := range m {
		// Print each key and its associated value
		fmt.Printf("%s:\n", key)
		// Print each item in the slice
		for i, item := range value {
			fmt.Printf("  %d - %s\n", i, item)
		}
	}
}

func main() {
	linksToCrawl := []string{
		"https://example.com",
		"https://golang.org",
		"https://github.com",
		"https://stackoverflow.com",
		"https://reddit.com",
		"https://cin.ufpe.br",
		"https://sigaa.ufpe.br/",
		"https://en.wikipedia.org/wiki/Go_(programming_language)",
	}

	linksHashMap := make(map[string][]string)

	// Process each link sequentially
	for _, link := range linksToCrawl {
		//fmt.Printf("Crawling %s...\n", link)
		result := crawlPage(link)
		linksHashMap[link] = result
	}

	// Print the results
	//printMap(linksHashMap)
}
