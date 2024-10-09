// This script interacts with the MediaWiki API to attempt to log a user in
// returning any error to the caller that they may encounter along the journey
// @TODO: this could do a bit more than just login people in potentially?
//
// This is free software licensed under the MIT License
// See LICENSE
// (c) OAuthority 2024

package gowikibot

import (
	"fmt"
)

// Define the struct based on the API response format
// the reason is optional and may be omitted if empty
// MediaWiki will not pass back a reason if the login was successful.
type LoginResponse struct {
	Login struct {
		Result string `json:"result"`
		Reason string `json:"reason,omitempty"`
	} `json:"login"`
}

type ApiError struct {
	Code string
	Info string
}

func (e ApiError) Error() string {
	if e.Info != "" {
		return fmt.Sprintf("API Error - Code: %s, Reason: %s", e.Code, e.Info)
	}
	return fmt.Sprintf("API Error - Code: %s", e.Code)
}

// Call out to the API and attempt to do a login, returning any errors we may
// encounter along the way.
func (c *Client) Login(username, password string) error {
	// Get login token
	token, err := c.GetToken("login")
	if err != nil {
		return err
	}

	v := Values{
		"action":     "login",
		"lgname":     username,
		"lgpassword": password,
		"lgtoken":    token,
	}

	fmt.Printf("Attempting to log you in as %s\n", username)

	// Make the POST request
	response, err := c.Post(v)
	if err != nil {
		return err
	}

	// Extract the login result from the response map
	loginData, ok := response["login"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid API response: missing 'login' field")
	}

	// Extract the "result" field
	result, ok := loginData["result"].(string)
	if !ok {
		return fmt.Errorf("invalid API response: missing 'result' field")
	}

	// Check if login was successful
	if result != "Success" {
		reason, _ := loginData["reason"].(string)
		return ApiError{Code: result, Info: reason}
	}

	fmt.Printf("Successfully logged in as %s\n", username)
	return nil
}