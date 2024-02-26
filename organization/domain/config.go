/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package domain provides domain models and configuration for a specific functionality.
package domain

// Config is a structure that holds the configuration settings for the application.
type Config struct {
	MaxCountPerOwner int64  `json:"max_count_per_owner"`
	InviteExpiry     int64  `json:"invite_expiry"`
	DefaultRole      string `json:"default_role"`
	Tables           tables `json:"tables"`
}

type tables struct {
	Member string `json:"member" required:"true"`
	Invite string `json:"invite" required:"true"`
}

// SetDefault sets the default values for the Config struct if they are not already set.
func (cfg *Config) SetDefault() {
	if cfg.MaxCountPerOwner <= 0 {
		cfg.MaxCountPerOwner = 10
	}
}
