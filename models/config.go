package models

import (
	"github.com/openmerlin/merlin-server/models/app"
	"github.com/openmerlin/merlin-server/models/controller"
	"github.com/openmerlin/merlin-server/models/infrastructure/modelrepositoryadapter"
)

// Config
type Config struct {
	App        app.Config                    `json:"app"`
	Tables     modelrepositoryadapter.Tables `json:"tables"`
	Controller controller.Config             `json:"controller"`
}

func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.App,
		&cfg.Tables,
		&cfg.Controller,
	}
}

func (cfg *Config) Init() {
	app.Init(&cfg.App)
	controller.Init(&cfg.Controller)
}
