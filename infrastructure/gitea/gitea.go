package gitea

import (
	"code.gitea.io/sdk/gitea"
)

var _client Client

// Client admin client
type Client struct {
	url    string
	client *gitea.Client
}

func Init(c *gitea.Client) {
	_client.client = c
	return
}

func GetClient() *Client {
	return &_client
}

func (c *Client) Url() string {
	return _client.url
}
