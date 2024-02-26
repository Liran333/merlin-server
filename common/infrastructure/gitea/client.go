/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package gitea provides a client for interacting with the Gitea API.
package gitea

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/openmerlin/go-sdk/gitea"
)

const timeout = 10

var (
	cli      *gitea.Client
	endpoint string
)

// Config represents the configuration for the Gitea client.
type Config struct {
	URL   string `json:"url"        required:"true"`
	Token string `json:"token"      required:"true"`
}

// Init initializes the Gitea client with the given configuration.
func Init(cfg *Config) error {
	client, err := gitea.NewClient(cfg.URL, gitea.SetToken(cfg.Token), gitea.SetHTTPClient(&http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // #nosec G402
		}},
		Timeout: time.Duration(timeout) * time.Second,
	}))
	if err == nil {
		cli = client
		endpoint = cfg.URL
	}

	return err
}

// Client returns the Gitea client.
func Client() *gitea.Client {
	return cli
}

// NewClient creates a new Gitea client with the given username and password.
func NewClient(username, password string) (*gitea.Client, error) {
	return gitea.NewClient(endpoint, gitea.SetBasicAuth(username, password), gitea.SetHTTPClient(&http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // #nosec G402
		}},
		Timeout: time.Duration(timeout) * time.Second,
	}))
}
