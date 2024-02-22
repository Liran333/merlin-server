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

type Config struct {
	URL   string `json:"url"        required:"true"`
	Token string `json:"token"      required:"true"`
}

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

func Client() *gitea.Client {
	return cli
}

func NewClient(username, password string) (*gitea.Client, error) {
	return gitea.NewClient(endpoint, gitea.SetBasicAuth(username, password), gitea.SetHTTPClient(&http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // #nosec G402
		}},
		Timeout: time.Duration(timeout) * time.Second,
	}))
}
