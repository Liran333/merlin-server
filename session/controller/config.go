/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides the controller logic for managing configuration and settings.
package controller

import "time"

var config Config

// Init initializes the configuration with the provided values.
func Init(cfg *Config) {
	config = *cfg
}

// Config is a struct that holds the configuration for CSRF token cookie expiry.
type Config struct {
	SessionCookieExpiry   int64  `json:"session_cookie_expiry"`
	CSRFTokenCookieExpiry int64  `json:"csrf_token_cookie_expiry"`
	SessionDomain         string `json:"session_domain"`
}

// SetDefault sets default values for the Config struct.
func (cfg *Config) SetDefault() {
	if cfg.CSRFTokenCookieExpiry <= 0 {
		cfg.CSRFTokenCookieExpiry = 3600 // second
	}

	if cfg.SessionCookieExpiry <= 0 {
		cfg.SessionCookieExpiry = 3600 // second
	}
}

func (cfg *Config) csrfTokenCookieExpiry() time.Time {
	return time.Now().Add(time.Duration(cfg.CSRFTokenCookieExpiry) * time.Second)
}

func (cfg *Config) sessionCookieExpiry() time.Time {
	return time.Now().Add(time.Duration(cfg.SessionCookieExpiry) * time.Second)
}
