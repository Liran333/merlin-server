/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package activity provides an adapter for the model repository
package activity

import (
	"github.com/openmerlin/merlin-server/activity/app"
	"github.com/openmerlin/merlin-server/activity/controller"
	"github.com/openmerlin/merlin-server/activity/insfrastructure/activityrepositoryadapter"
)

// Config is a struct that represents the overall configuration for the application.
type Config struct {
	App        app.Config                       `json:"app"`
	Tables     activityrepositoryadapter.Tables `json:"tables"`
	Usages     activityrepositoryadapter.Config `json:"usages"`
	Controller controller.Config                `json:"controller"`
}

// ConfigItems returns a slice of interface{} containing pointers to the configuration items in the Config struct.
func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.App,
		&cfg.Tables,
		&cfg.Controller,
	}
}

// Init initializes the application using the configuration settings provided in the Config struct.
func (cfg *Config) Init() {
	app.Init(&cfg.App)
	controller.Init(&cfg.Controller)
	activityrepositoryadapter.InitUsage(&cfg.Usages)
}
