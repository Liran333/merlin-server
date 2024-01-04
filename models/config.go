package models

import (
	"github.com/openmerlin/merlin-server/models/controller"
	"github.com/openmerlin/merlin-server/models/infrastructure/modelrepositoryadapter"
)

// Config
type Config struct {
	Tables     modelrepositoryadapter.Tables `json:"tables"`
	Controller controller.Config             `json:"controller"`
}

func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Tables,
		&cfg.Controller,
	}
}

func (cfg *Config) Init() {
	controller.Init(&cfg.Controller)
}
