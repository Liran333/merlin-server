package gitea

import (
	"fmt"

	"code.gitea.io/sdk/gitea"
)

var _client Client

type Config struct {
	Url   string `json:"url"        required:"true"`
	Token string `json:"token"      required:"true"`
}

// admin client
type Client struct {
	url    string
	client *gitea.Client
}

func GetClient() *Client {
	return &_client
}

func (c *Client) Url() string {
	return _client.url
}

func Init(cfg *Config) (err error) {
	if cfg == nil {
		return fmt.Errorf("cfg is nil")
	}

	_client.client, err = gitea.NewClient(cfg.Url, gitea.SetToken(cfg.Token))
	if err == nil {
		_client.url = cfg.Url
	}

	return
}
