package msgraph

import (
	"golang.org/x/oauth2"
	"os"
)

type NullCache struct{}

func (nc NullCache) Load(c *Client) (*oauth2.Token, error) {
	return nil, os.ErrNotExist
}

func (nc NullCache) Save(c *Client, token *oauth2.Token) error {
	return nil
}
