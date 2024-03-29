/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package coderepo provides configuration and initialization functionality for the code repository application.
package coderepo

import (
	"github.com/openmerlin/merlin-server/coderepo/domain/primitive"
	"github.com/openmerlin/merlin-server/coderepo/infrastructure/branchrepositoryadapter"
	"github.com/openmerlin/merlin-server/coderepo/infrastructure/coderepoadapter"
)

// Config is a struct that represents the overall configuration for the application.
type Config struct {
	Tables     branchrepositoryadapter.Tables `json:"tables"`
	Primitive  primitive.Config               `json:"primitive"`
	Repository coderepoadapter.Config         `json:"repository"`
}

// ConfigItems returns a slice of interface{} containing pointers to the configuration items in the Config struct.
func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Tables,
		&cfg.Primitive,
		&cfg.Repository,
	}
}

// Init initializes the application using the configuration settings provided in the Config struct.
func (cfg *Config) Init() {
	primitive.Init(&cfg.Primitive)
}
