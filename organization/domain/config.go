/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package domain provides domain models and configuration for a specific functionality.
package domain

import (
	"github.com/openmerlin/merlin-server/organization/domain/primitive"
	"github.com/openmerlin/merlin-server/organization/infrastructure/messageadapter"
)

// Config is a structure that holds the configuration settings for the application.
type Config struct {
	MaxCountPerOwner int64                 `json:"max_count_per_owner"`
	InviteExpiry     int64                 `json:"invite_expiry"`
	DefaultRole      string                `json:"default_role"`
	Tables           tables                `json:"tables"`
	Topics           messageadapter.Topics `json:"topics"`
	Primitive        primitive.Config      `json:"primitive"`
	CertificateEmail []string              `json:"certificate_email"`
}

type tables struct {
	Member      string `json:"member"      required:"true"`
	Invite      string `json:"invite"      required:"true"`
	Certificate string `json:"certificate" required:"true"`
}

// SetDefault sets the default values for the Config struct if they are not already set.
func (cfg *Config) SetDefault() {
	if cfg.MaxCountPerOwner <= 0 {
		cfg.MaxCountPerOwner = 10
	}
}

func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Primitive,
	}
}

func (cfg *Config) Init() error {
	return primitive.Init(cfg.Primitive)
}
