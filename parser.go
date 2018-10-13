package main

import (
	"strings"

	"github.com/antchfx/xquery/html"
)

// Finds the data-big-photo links in the page
func parsePage(page string) []string {
	doc, err := htmlquery.Parse(strings.NewReader(page))

	panicOnError(err, "Couldn't parse page!")
	elements := htmlquery.Find(doc, "//a[@data-big-photo]")

	ret := make([]string, 0, len(elements))
	for _, element := range elements {
		ret = append(ret, htmlquery.SelectAttr(element, "data-big-photo"))
	}
	return ret
}

// A worker got handle work input and output
func parseAllLinks(getWork chan string, result chan string) {
	for {
		work, got := <-getWork
		if !got {
			break
		}
		links := parsePage(work)

		for _, link := range links {
			result <- link
		}
	}
	close(result)
}
