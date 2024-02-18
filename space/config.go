package space

import (
	"github.com/openmerlin/merlin-server/space/app"
	"github.com/openmerlin/merlin-server/space/controller"
	"github.com/openmerlin/merlin-server/space/domain/primitive"
	"github.com/openmerlin/merlin-server/space/infrastructure/messageadapter"
	"github.com/openmerlin/merlin-server/space/infrastructure/spacerepositoryadapter"
)

// Config
type Config struct {
	App        app.Config                    `json:"app"`
	Tables     spacerepositoryadapter.Tables `json:"tables"`
	Topics     messageadapter.Topics         `json:"topics"`
	Primitive  primitive.Config              `json:"primitive"`
	Controller controller.Config             `json:"controller"`
}

func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.App,
		&cfg.Tables,
		&cfg.Topics,
		&cfg.Primitive,
		&cfg.Controller,
	}
}

func (cfg *Config) Init() {
	app.Init(&cfg.App)
	primitive.Init(&cfg.Primitive)
	controller.Init(&cfg.Controller)
}
