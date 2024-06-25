/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package domain provides domain models and configuration for a specific functionality.
package domain

type tables struct {
	User  string `json:"user"  required:"true"`
	Token string `json:"token" required:"true"`
}

// Config is a struct that holds the configuration for the program.
type Config struct {
	Key             []byte `json:"key"    required:"true"`
	Tables          tables `json:"tables" required:"true"`
	ObsPath         string `json:"obs_path" required:"true"`
	ObsBucket       string `json:"obs_bucket" required:"true"`
	CdnEndpoint     string `json:"cdn_endpoint" required:"true"`
	MaxTokenPerUser int    `json:"max_token_per_user" required:"true"`
}

// SetDefault sets the default values for the Config struct if they are not already set.
func (cfg *Config) SetDefault() {
	if cfg.MaxTokenPerUser <= 0 {
		cfg.MaxTokenPerUser = 20
	}
}
