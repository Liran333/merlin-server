package controller

import "time"

var config Config

func Init(cfg *Config) {
	config = *cfg
}

type Config struct {
	CSRFTokenCookieExpiry int64 `json:"csrf_token_cookie_expiry"`
}

func (cfg *Config) SetDefault() {
	if cfg.CSRFTokenCookieExpiry <= 0 {
		cfg.CSRFTokenCookieExpiry = 5 * 60 // second
	}
}

func (cfg *Config) csrfTokenCookieExpiry() time.Time {
	return time.Now().Add(time.Duration(cfg.CSRFTokenCookieExpiry) * time.Second)
}
