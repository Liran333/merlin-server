/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package models provides configuration and initialization functionality for the application.
package models

import (
	"github.com/openmerlin/merlin-server/models/app"
	"github.com/openmerlin/merlin-server/models/controller"
	"github.com/openmerlin/merlin-server/models/infrastructure/messageadapter"
	"github.com/openmerlin/merlin-server/models/infrastructure/modelrepositoryadapter"
)

// Config is a struct that represents the overall configuration for the application.
type Config struct {
	App        app.Config                    `json:"app"`
	Tables     modelrepositoryadapter.Tables `json:"tables"`
	Topics     messageadapter.Topics         `json:"topics"`
	Controller controller.Config             `json:"controller"`
}

// ConfigItems returns a slice of interface{} containing pointers to the configuration items in the Config struct.
func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.App,
		&cfg.Tables,
		&cfg.Topics,
		&cfg.Controller,
	}
}

// Init initializes the application using the configuration settings provided in the Config struct.
func (cfg *Config) Init() {
	app.Init(&cfg.App)
	controller.Init(&cfg.Controller)
}
