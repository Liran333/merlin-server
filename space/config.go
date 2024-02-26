/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package space provides configuration and initialization functionality for the space application.
package space

import (
	"github.com/openmerlin/merlin-server/space/app"
	"github.com/openmerlin/merlin-server/space/controller"
	"github.com/openmerlin/merlin-server/space/domain/primitive"
	"github.com/openmerlin/merlin-server/space/infrastructure/messageadapter"
	"github.com/openmerlin/merlin-server/space/infrastructure/spacerepositoryadapter"
)

// Config is a struct that represents the overall configuration for the application.
type Config struct {
	App        app.Config                    `json:"app"`
	Tables     spacerepositoryadapter.Tables `json:"tables"`
	Topics     messageadapter.Topics         `json:"topics"`
	Primitive  primitive.Config              `json:"primitive"`
	Controller controller.Config             `json:"controller"`
}

// ConfigItems returns a slice of interface{} containing pointers to the configuration items in the Config struct.
func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.App,
		&cfg.Tables,
		&cfg.Topics,
		&cfg.Primitive,
		&cfg.Controller,
	}
}

// Init initializes the application using the configuration settings provided in the Config struct.
func (cfg *Config) Init() {
	app.Init(&cfg.App)
	primitive.Init(&cfg.Primitive)
	controller.Init(&cfg.Controller)
}
