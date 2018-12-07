package main

import (
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

// Gets the form ID from login screen in tumblr, needed for future requests
func getFormID(content io.Reader) (string, error) {
	doc, err := htmlquery.Parse(content)

	if err != nil {
		return "", err
	}

	candidate := htmlquery.FindOne(doc, "//input[@type=\"hidden\" and @name=\"form_key\"]")
	formID := htmlquery.SelectAttr(candidate, "value")

	return formID, nil
}

// Gets a liked page given it's number
func getPage(pageNum int, client *http.Client) string {
	resp, err := client.Get("https://www.tumblr.com/likes/page/" + strconv.Itoa(pageNum))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting page: %v\n", pageNum)
	}
	defer resp.Body.Close()

	content, _ := ioutil.ReadAll(resp.Body)
	return string(content)
}

// a worker for easier goroutine
func getPageWorker(work <-chan int, ready chan<- bool, client *http.Client, collect chan string) {
	for {
		pageNum, gotWork := <-work
		if !gotWork {
			break
		}

		/*
			page := getPage(pageNum, client)
			lenLink := len(parsePage(page))
			if lenLink < 3 {
				writeFile(pageNum, []byte(page))
				fmt.Fprintf(os.Stderr, "Found fewer than 3: %d.\n", lenLink)
			}
		*/

		collect <- getPage(pageNum, client)
	}
	ready <- true
}

// gets all the pages in likes and passes the strings through 'collect'
func getPages(email, password string, pageCount, workerCount int, collect chan string) {
	jar, _ := cookiejar.New(&cookiejar.Options{})
	client := &http.Client{
		Jar: jar,
	}

	v := url.Values{}
	v.Set("determine_email", email)
	v.Set("user[email]", email)
	v.Set("user[password]", password)
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
	for i := 0; i < workerCount; i++ {
		go getPageWorker(giveWork, ready, client, collect)
	}

	// Give the work
	for i := 1; i <= pageCount; i++ {
		giveWork <- i
	}
	close(giveWork)

	for i := 0; i < workerCount; i++ {
		<-ready
	}

	close(ready)
	close(collect)
}
