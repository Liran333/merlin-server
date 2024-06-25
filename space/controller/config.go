/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller contains the controllers for handling various operations and logic in the application.
package controller

var config Config

// Init initializes the controller package with the provided configuration.
func Init(cfg *Config) {
	config = *cfg
}

// Config represents the configuration settings for the controller package.
type Config struct {
	MaxCountPerPage  int   `json:"max_count_per_page"`
	MaxCoverFileSize int64 `json:"max_cover_file_size"`
}

// SetDefault sets the default values for the configuration.
func (cfg *Config) SetDefault() {
	if cfg.MaxCountPerPage <= 0 {
		cfg.MaxCountPerPage = 100
	}

	if cfg.MaxCoverFileSize <= 0 {
		cfg.MaxCoverFileSize = 5242880
	}
}
