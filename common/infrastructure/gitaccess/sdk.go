package gitaccess

import (
	"github.com/openmerlin/git-access-sdk/httpclient"
)

func Init(cfg *Config) {
	httpclient.Init(cfg)
}

// Config is for http client config
type Config = httpclient.Config
