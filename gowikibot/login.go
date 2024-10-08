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
	"parameters"
)

// Call out to the API and attempt to do a login, returning any errors we may
// encounter along the way.
func (c *Client) Login(username, password string) error {

	token, err := c.GetToken(LoginToken)

	if err != nil {
		return err
	}

	v := Values{
		"action":     "login",
		"lgname":     username,
		"lgpassword": password,
		"lgtoken":    token,
	}

	fmt.Printf("Attempting to log you in as %s", username)

	response, err := c.Post(v)

	if err != nil {
		return err
	}

	loginResult, err := response.GetString("login", "result")

	if err != nil {
		return fmt.Errorf("the API response is not valid, the API did not return a string for the login result.")
	}

	if loginResult != "Success" {
		apiError := ApiError{Code: loginResult}

		if reason, err := response.GetString("login", "reason"); err == nil {
			apiError.Info = reason
		}

		return apiError
	}

	fmt.Printf("Successfully logged in as %s", username)
	return nil
}
