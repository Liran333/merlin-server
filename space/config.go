package space

import (
	"github.com/openmerlin/merlin-server/space/controller"
	"github.com/openmerlin/merlin-server/space/infrastructure/spacerepositoryadapter"
)

// Config
type Config struct {
	Tables     spacerepositoryadapter.Tables `json:"tables"`
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
