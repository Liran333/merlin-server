package internalservice

var config Config

func Init(cfg *Config) {
	config = *cfg
}

type Config struct {
	TokenHash string `json:"token_hash" required:"true"`
	Salt      string `json:"salt" required:"true"`
}
