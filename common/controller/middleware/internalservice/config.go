package internalservice

var config Config

func Init(cfg *Config) {
	config = *cfg
}

type Config struct {
	Token string `json:"token" required:"true"`
}
