// This is the core of Gowikibot, and has several helper methods that can be used
// in scripts. Scripts should rely on, or extend where possible, the methods
// provided here instead of constructing their own.
//
// If there is no method for what you are trying/wanting to do, then open a PR and
// help implement it.
//
// This is free software licensed under the MIT License
// See LICENSE
// (c) OAuthority 2024

package gowikibot

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

// default UserAgent in case the user did not provide one
const UserAgent = "gowikibot (https://github.com/oauthority/go-utils)"

// struct to represent the Client
type Client struct {
	httpClient *http.Client
	apiUrl     *url.URL
	UserAgent  string
	Tokens     map[string]string
	debug      io.Writer
}

// construct a new API client for use throughout, with a cookie jar to ensure
// that we retain the edit tokens/cookies et al.
func newApiClient(apiUrl string, userAgent string) (*Client, error) {
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
		apiUrl:     parsedUrl,
		UserAgent:  userAgent,
		Tokens:     make(map[string]string),
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

// Query the MediaWiki API for a specific token; can be any type of token
// we always get a fresh login token, do not check in the Client struct for
// the login token because it might not be fresh
func (c *Client) GetToken(tokenName string) (string, error) {

	// Always obtain a fresh edit token
	if tokenName != "login" {
		if tok, ok := c.Tokens[tokenName]; ok {
			return tok, nil
		}
	}

	p := Values{
		"action": "query",
		"meta":   "tokens",
		"type":   tokenName,
	}

	// Make the GET request
	response, err := c.Get(p)
	if err != nil {
		return "", err
	}

	// Navigate the response map to extract the token
	queryData, ok := response["query"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected API response format: missing 'query' field")
	}

	tokensData, ok := queryData["tokens"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected API response format: missing 'tokens' field")
	}

	tokenKey := tokenName + "token"
	token, ok := tokensData[tokenKey].(string)
	if !ok {
		return "", fmt.Errorf("unexpected API response format: missing or invalid token field '%s'", tokenKey)
	}

	// Only store tokens other than the login token
	if tokenName != "login" {
		c.Tokens[tokenName] = token
	}

	return token, nil
}

// Make a GET request to the API and return the JSON.
// MediaWiki's default return method is JSON. Whilst MediaWiki
// supports returning a result as XML or PHP, Gowikibot does.
// not support this and it is rarely—if ever—used in MediaWiki
// we just omit the format type and we will get JSON back
func (c *Client) Get(v Values) (map[string]interface{}, error) {

	v.Set("format", "json")
	v.Set("formatversion", "2")

	req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", c.apiUrl.String(), v.Encode()), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create HTTP request (params: %v): %v", v, err)
	}

	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error occurred during HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Read and parse the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var jsonResponse map[string]interface{}
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %v", err)
	}

	return jsonResponse, nil
}

// Make a POST request to the API and return the JSON response as a map.
func (c *Client) Post(v Values) (map[string]interface{}, error) {

	v.Set("format", "json")
	v.Set("formatversion", "2")

	reqBody := strings.NewReader(v.Encode())
	req, err := http.NewRequest("POST", c.apiUrl.String(), reqBody)
	if err != nil {
		return nil, fmt.Errorf("unable to create HTTP request (params: %v): %v", v, err)
	}

	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error occurred during HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Read and parse the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var jsonResponse map[string]interface{}
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %v", err)
	}

	return jsonResponse, nil
}

// Read a config file and return the family that we are operating on
// potentially could be refined to return anything rather than just the family?
func LoadCredentialsFromFile(filename string) (map[string]Family, error) {

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var credentials map[string]Family
	err = json.Unmarshal(data, &credentials)
	if err != nil {
		return nil, err
	}

	return credentials, nil
}
