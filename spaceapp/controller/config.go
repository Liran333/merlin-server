package controller

var config Config

func Init(cfg *Config) {
	config = *cfg
}

type Config struct {
	SSEToken string `json:"sse_token"`
}