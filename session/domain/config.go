/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package domain provides domain models and configuration for the session management functionality.
package domain

import "time"

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
	MaxSessionNum    int   `json:"max_session_num"`
	SessionTimeout   int64 `json:"session_timeout"`
	CSRFTokenTimeout int64 `json:"csrf_token_timeout"`
	sessionTimeout   time.Duration
	csrfTokenTimeout time.Duration
}

// SetDefault sets the default values for the Config struct.
func (cfg *Config) SetDefault() {
	if cfg.MaxSessionNum <= 0 {
		cfg.MaxSessionNum = 3
	}

	if cfg.CSRFTokenTimeout <= 0 {
		cfg.CSRFTokenTimeout = hours * seconds
	}

	if cfg.SessionTimeout <= 0 {
		cfg.SessionTimeout = 60 * 60
	}

	cfg.csrfTokenTimeout = time.Duration(cfg.CSRFTokenTimeout) * time.Second
	cfg.sessionTimeout = time.Duration(cfg.SessionTimeout) * time.Second
}
