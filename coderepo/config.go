/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package coderepo provides configuration and initialization functionality for the code repository application.
package coderepo

import (
	"github.com/openmerlin/merlin-server/coderepo/domain/primitive"
	"github.com/openmerlin/merlin-server/coderepo/infrastructure/branchrepositoryadapter"
)

// Config is a struct that represents the overall configuration for the application.
type Config struct {
	Tables    branchrepositoryadapter.Tables `json:"tables"`
	Primitive primitive.Config               `json:"primitive"`
}

// ConfigItems returns a slice of interface{} containing pointers to the configuration items in the Config struct.
func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Tables,
		&cfg.Primitive,
	}
}

// Init initializes the application using the configuration settings provided in the Config struct.
func (cfg *Config) Init() {
	primitive.Init(&cfg.Primitive)
}
