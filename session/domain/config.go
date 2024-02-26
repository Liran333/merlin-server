/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package domain provides domain models and configuration for the session management functionality.
package domain

const (
	hours   = 8
	seconds = 3600
)

var config Config

// Init initializes the configuration with the given Config struct.
func Init(cfg *Config) {
	config = *cfg
}

// Config is a struct that holds the configuration for max session num,
// csrf token timeout and csrf token timeout to reset.
type Config struct {
	MaxSessionNum           int   `json:"max_session_num"`
	CSRFTokenTimeout        int64 `json:"csrf_token_timeout"`
	CSRFTokenTimeoutToReset int64 `json:"csrf_token_timeout_to_reset"`
}

// SetDefault sets the default values for the Config struct.
func (cfg *Config) SetDefault() {
	if cfg.MaxSessionNum <= 0 {
		cfg.MaxSessionNum = 3
	}

	if cfg.CSRFTokenTimeout <= 0 {
		cfg.CSRFTokenTimeout = hours * seconds
	}

	if cfg.CSRFTokenTimeoutToReset <= 0 {
		cfg.CSRFTokenTimeoutToReset = 3
	}
}
