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
	CSRFTokenCookieExpiry int64 `json:"csrf_token_cookie_expiry"`
}

// SetDefault sets default values for the Config struct.
func (cfg *Config) SetDefault() {
	if cfg.CSRFTokenCookieExpiry <= 0 {
		cfg.CSRFTokenCookieExpiry = 5 * 60 // second
	}
}

func (cfg *Config) csrfTokenCookieExpiry() time.Time {
	return time.Now().Add(time.Duration(cfg.CSRFTokenCookieExpiry) * time.Second)
}
