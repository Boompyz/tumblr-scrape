package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
)

func panicOnError(err error, message string) {
	if err != nil {
		panic(message)
	}
}

func writeFile(pageNum int, content []byte) {
	ioutil.WriteFile("pages/"+strconv.Itoa(pageNum)+".html", content, 0644)
}

func main() {
	var email = flag.String("email", "", "tumblr username")
	var password = flag.String("password", "", "tumblr password")
	var pageCount = flag.Int("pageCount", 0, "The number of pages to scrape")
	var workerCount = flag.Int("workerCount", 5, "The number of threads to get pages with")

	flag.Parse()

	if *email == "" || *password == "" || *pageCount == 0 {
		fmt.Println("Username, password or page count not specified!")
	}

	collectPages := make(chan string)
	go getPages(*email, *password, *pageCount, *workerCount, collectPages)

	collectLinks := make(chan string)
	go parseAllLinks(collectPages, collectLinks)

	for {
		link, got := <-collectLinks
		if !got {
			break
		}
		fmt.Println(link)
	}
}
