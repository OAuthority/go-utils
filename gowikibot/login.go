package gowikibot

import (
	"fmt"
)

func (c *Client) Login(username, password string) error {

	token, err := c.GetToken(LoginToken)

	if err != nil {
		return err
	}

	v := parameters.Values{
		"action": "login",
		"lgname": username,
		"lgpassword": password,
		"lgtoken": token,
	}

	
}
