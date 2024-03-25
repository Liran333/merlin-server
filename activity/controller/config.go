/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides functionality for managing the application's controllers.
package controller

// nolint:golint,unused
var config Config

// Init initializes the configuration.
func Init(cfg *Config) {
	config = *cfg
}

// Config represents the application configuration.
type Config struct {
	MaxCountPerPage int `json:"max_count_per_page"`
}

// SetDefault sets the default values for the configuration.
func (cfg *Config) SetDefault() {
	if cfg.MaxCountPerPage <= 0 {
		cfg.MaxCountPerPage = 100
	}
}
