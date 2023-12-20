package gitea

import "code.gitea.io/sdk/gitea"

var (
	cli      *gitea.Client
	endpoint string
)

type Config struct {
	URL   string `json:"url"        required:"true"`
	Token string `json:"token"      required:"true"`
}

func Init(cfg *Config) error {
	client, err := gitea.NewClient(cfg.URL, gitea.SetToken(cfg.Token))
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
	return gitea.NewClient(endpoint, gitea.SetBasicAuth(username, password))
}
