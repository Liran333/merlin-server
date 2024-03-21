/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package securestorage provides an adapter for working with vault-related functionality.
package securestorage

const (
	initBasePath = "space_env_secret"
)

// Config is a struct that represents the vault config
type Config struct {
	Address  string `json:"address" required:"true"`
	UserName string `json:"user_name" required:"true"`
	PassWord string `json:"pass_word" required:"true"`
	BasePath string `json:"base_path" required:"true"`
}

// SetDefault sets the default values for the configuration.
func (cfg *Config) SetDefault() {
	if cfg.BasePath == "" {
		cfg.BasePath = initBasePath
	}
}
