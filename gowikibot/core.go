package gowikibot

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

// default UserAgent in case the user did not provide one
const UserAgent = "gowikibot (https://github.com/oauthority/go-utils)"

// struct to represent the Client
type Client struct {
	httpClient *http.Client
	apiUrl *url.URL
	UserAgent string
	Tokens map[string]string
	debug io.Writer
}

// construct a new API client for use throughout, with a cookie jar to ensure
// that we retain the edit tokens/cookies et al. 
func newApiClient(apiUrl string, userAgent string ) (*Client, error) {
	parsedUrl, err := url.Parse(apiUrl)

	if err != nil {
		return nil, fmt.Errorf("invalid url passsed to the client: %v", err)
	}

	// if we didnt' provide a userAgent, then that's fine, use the defaul
	// one we have supplied earlier, to ensure that we are adhering to  
	// WP:User_Agent_Policy et al.
	if userAgent == "" {
		userAgent = UserAgent
	}

	jar, _ := cookiejar.New(nil)

	client := &Client{
		httpClient: &http.Client{Jar: jar},
		apiUrl: parsedUrl,
		UserAgent: userAgent,
		Tokens: make(map[string]string),
	}

	return client, nil
}

// override the default http.Client with ours
func (c *Client) SetHTTPClient(httpClient *http.Client) {
    if httpClient.Jar == nil {
        jar, _ := cookiejar.New(nil)
        httpClient.Jar = jar
    }
    c.httpClient = httpClient
}
