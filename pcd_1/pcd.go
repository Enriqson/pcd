package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

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
		log.Fatal(err)
		return []string{}
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
		return []string{}
	}

	var links []string

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		link_ref := s.AttrOr("href", "")
		if isValidHTTPUrl(link_ref) {
			links = append(links, link_ref)
		}
	})
	//printArr(links)
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

func crawlThreadSafe(link string, wg *sync.WaitGroup, mx *sync.Mutex, linksHashMap map[string][]string) {
	defer wg.Done()
	result := crawlPage(link)
	mx.Lock()
	//fmt.Printf("Saving crawl results for %s\n", link)
	linksHashMap[link] = result
	mx.Unlock()

}

func main() {
	var links_to_crawl = []string{
		"https://example.com",
		"https://golang.org",
		"https://github.com",
		"https://stackoverflow.com",
		"https://reddit.com",
		"https://cin.ufpe.br",
		"https://sigaa.ufpe.br/",
		"https://en.wikipedia.org/wiki/Go_(programming_language)",
	}

	wg := sync.WaitGroup{}
	mx := sync.Mutex{}
	linksHashMap := make(map[string][]string)

	for _, link := range links_to_crawl {
		wg.Add(1)
		go crawlThreadSafe(link, &wg, &mx, linksHashMap)
	}
	wg.Wait()
	//printMap(linksHashMap)
}
