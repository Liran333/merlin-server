/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package app provides functionality for the application.
package app

var config Config

// Init initializes the application with the provided configuration.
func Init(cfg *Config) {
	config = *cfg
}

// Config is a struct that holds the configuration for max count per owner.
type Config struct {
	MaxCountPerOrg  int    `json:"max_count_per_org"`
	MaxCountPerUser int    `json:"max_count_per_user"`
	GegexpRule      string `json:"regexp_rule"`
}

// SetDefault sets the default values for the Config struct.
func (cfg *Config) SetDefault() {
	if cfg.MaxCountPerUser <= 0 {
		cfg.MaxCountPerUser = 50
	}

	if cfg.MaxCountPerOrg <= 0 {
		cfg.MaxCountPerOrg = 200
	}
}
