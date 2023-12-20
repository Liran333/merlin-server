package gitea

import (
	"code.gitea.io/sdk/gitea"
	"fmt"
)

type Config struct {
	Url   string `json:"url"        required:"true"`
	Token string `json:"token"      required:"true"`
}

func Init(cfg *Config) (client *gitea.Client, err error) {
	if cfg == nil {
		return nil, fmt.Errorf("cfg is nil")
	}
	client, err = gitea.NewClient(cfg.Url, gitea.SetToken(cfg.Token))
	return client, err
}
