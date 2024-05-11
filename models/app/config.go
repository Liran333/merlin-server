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
	MaxCountPerOwner int              `json:"max_count_per_owner"`
	RecommendModels  []RecommendIndex `json:"recommend_models"`
}

type RecommendIndex struct {
	Owner    string `json:"owner" required:"true"`
	Reponame string `json:"reponame" required:"true"`
}

// SetDefault sets the default values for the Config struct.
func (cfg *Config) SetDefault() {
	if cfg.MaxCountPerOwner <= 0 {
		cfg.MaxCountPerOwner = 1000
	}
}
