/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
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
	MaxCountPerOwner      int `json:"max_count_per_owner"`
	MaxCountSpaceVariable int `json:"max_count_space_variable"`
	MaxCountSpaceSecret   int `json:"max_count_space_secret"`
}

// SetDefault sets the default values for the Config struct.
func (cfg *Config) SetDefault() {
	if cfg.MaxCountPerOwner <= 0 {
		cfg.MaxCountPerOwner = 1000
	}
	if cfg.MaxCountSpaceVariable <= 0 {
		cfg.MaxCountSpaceVariable = 100
	}
	if cfg.MaxCountSpaceSecret <= 0 {
		cfg.MaxCountSpaceSecret = 100
	}
}
