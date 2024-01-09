package space

import (
	"github.com/openmerlin/merlin-server/space/controller"
	"github.com/openmerlin/merlin-server/space/domain/primitive"
	"github.com/openmerlin/merlin-server/space/infrastructure/spacerepositoryadapter"
)

// Config
type Config struct {
	Tables     spacerepositoryadapter.Tables `json:"tables"`
	Primitive  primitive.Config              `json:"primitive"`
	Controller controller.Config             `json:"controller"`
}

func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Tables,
		&cfg.Primitive,
		&cfg.Controller,
	}
}

func (cfg *Config) Init() {
	primitive.Init(&cfg.Primitive)
	controller.Init(&cfg.Controller)
}
