package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"

	htmlquery "github.com/antchfx/xquery/html"
)

func panicOnError(err error, message string) {
	if err != nil {
		panic(message + err.Error())
	}
}

func writeFile(pageNum int, content []byte) {
	ioutil.WriteFile("pages/"+strconv.Itoa(pageNum)+".html", content, 0644)
}

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

func login(email, password string) *http.Client {
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
	resp, err = client.PostForm("https://www.tumblr.com/login", v)
	panicOnError(err, "Couldn't login (post login info)")
	return client
}

func main() {
	var email = flag.String("email", "", "tumblr username")
	var password = flag.String("password", "", "tumblr password")
	var token = flag.String("token", "", "OAuth Toekn. Looks like \"Bearer aaaaa...\"")
	var stopAfter = flag.Int("n", 2147483647, "Stop after the first n pages.")

	flag.Parse()

	if *email == "" || *password == "" {
		fmt.Println("Username, password or page count not specified!")
	}

	client := login(*email, *password)
	req, err := http.NewRequest("GET", "https://www.tumblr.com/api/v2/user/likes?fields%5Bblogs%5D=name%2Cavatar%2Ctitle%2Curl%2Cis_adult%2C%3Fis_member%2Cdescription_npf%2Cuuid%2Ccan_be_followed%2C%3Ffollowed%2C%3Fadvertiser_name%2Cis_paywall_on%2Ctheme%2Csubscription_plan%2C%3Fprimary&limit=21&reblog_info=true&before=1608550500", nil)
	for i := 0; i < *stopAfter; i++ {
		if err != nil {
			println("End of pages.")
			break
		}
		req.Header.Add("Authorization", *token)

		resp, err := client.Do(req)
		r, err := FromHTTPResponse(resp)
		panicOnError(err, "Error making request ")

		links := r.GetImages()
		for _, link := range links {
			fmt.Println(link)
		}

		req, err = r.NextRequest()
	}
}
