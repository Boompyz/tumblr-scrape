package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"

	"github.com/antchfx/xquery/html"
)

func panicOnError(err error, message string) {
	if err != nil {
		panic(message)
	}
}

func writeFile(pageNum int, content []byte) {
	ioutil.WriteFile("pages/"+strconv.Itoa(pageNum)+".html", content, 0644)
}

// Extracts links from a page, whose innerText has the given string
func getFormID(content io.Reader) (string, error) {
	doc, err := htmlquery.Parse(content)

	if err != nil {
		return "", err
	}

	candidate := htmlquery.FindOne(doc, "//input[@type=\"hidden\" and @name=\"form_key\"]")
	formID := htmlquery.SelectAttr(candidate, "value")

	return formID, nil
}

func getPage(pageNum int, client *http.Client) {
	resp, err := client.Get("https://www.tumblr.com/likes/page/" + strconv.Itoa(pageNum))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting page: %v\n", pageNum)
	}
	defer resp.Body.Close()

	content, _ := ioutil.ReadAll(resp.Body)
	writeFile(pageNum, content)
}

func getPageWorker(work <-chan int, ready chan<- bool, client *http.Client) {
	for {
		pageNum, gotWork := <-work
		if !gotWork {
			break
		}

		getPage(pageNum, client)
	}
	ready <- true
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

	jar, _ := cookiejar.New(&cookiejar.Options{})
	client := &http.Client{
		Jar: jar,
	}

	v := url.Values{}
	v.Set("determine_email", *email)
	v.Set("user[email]", *email)
	v.Set("user[password]", *password)
	v.Set("action", "signup_determine")

	// Get to the login page to aquire an ID
	resp, err := client.Get("https://www.tumblr.com/login")
	panicOnError(err, "Couldn't login (get login page)")

	// Extract the ID from the page and put it in the values
	formID, err := getFormID(resp.Body)
	panicOnError(err, "Couldn't find form id")
	v.Set("form_key", formID)

	// Post the login info
	_, err = client.PostForm("https://www.tumblr.com/login", v)
	panicOnError(err, "Couldn't login (post login info)")

	giveWork := make(chan int)
	ready := make(chan bool)

	// Spawn the workers
	for i := 0; i < *workerCount; i++ {
		go getPageWorker(giveWork, ready, client)
	}

	// Give the work
	for i := 1; i <= *pageCount; i++ {
		giveWork <- i
	}
	close(giveWork)

	// Wait for the workers to finish
	for i := 0; i < *workerCount; i++ {
		<-ready
	}
}
