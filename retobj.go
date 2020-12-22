package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// RetObj is used to parse the Likes response from Tumblr
type RetObj struct {
	Response struct {
		LikedPosts []struct {
			ObjectType   string `json:"object_type"`
			OriginalType string `json:"original_type"`
			BlogName     string `json:"blog_name"`
			IsNSFW       bool   `json:"is_nsfw"`
			Content      []struct {
				Type  string `json:"type"`
				Media []struct {
					Type                  string `json:"type"`
					Width                 int    `json:"width"`
					Height                int    `json:"height"`
					URL                   string `json:"url"`
					HasOriginalDimensions bool   `json:"has_original_dimensions"`
				} `json:"media"`
			} `json:"content"`
		} `json:"liked_posts"`
		Links struct {
			Next struct {
				Href string `json:"href"`
			} `json:"next"`
		} `json:"_links"`
	} `json:"response"`
}

// GetImages returns links to the the highest quality images from each post.
func (r *RetObj) GetImages() []string {
	links := make([]string, len(r.Response.LikedPosts))

	for idx, post := range r.Response.LikedPosts {
		size := 0
		link := ""
		for _, content := range post.Content {
			if content.Type != "image" {
				continue
			}
			for _, media := range content.Media {
				if media.Width*media.Height > size {
					size = media.Height * media.Width
					link = media.URL
				}
			}
		}
		links[idx] = link
	}

	return links
}

// NextRequest returns the next request we have to do for the following images.
// It does not contain authentication.
func (r *RetObj) NextRequest() (*http.Request, error) {
	if r.Response.Links.Next.Href == "" {
		return nil, errors.New("end of pages")
	}
	link := "https://www.tumblr.com/api" + r.Response.Links.Next.Href
	return http.NewRequest("GET", link, nil)
}

// FromHTTPResponse creates an RetObj from the response given by the api.
func FromHTTPResponse(resp *http.Response) (*RetObj, error) {
	var r *RetObj = &RetObj{}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	err := json.Unmarshal(body, r)
	return r, err
}
