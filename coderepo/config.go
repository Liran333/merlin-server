package coderepo

import (
	"github.com/openmerlin/merlin-server/coderepo/domain/primitive"
	"github.com/openmerlin/merlin-server/coderepo/infrastructure/branchrepositoryadapter"
)

// Config
type Config struct {
	Tables    branchrepositoryadapter.Tables `json:"tables"`
	Primitive primitive.Config               `json:"primitive"`
}

func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Tables,
		&cfg.Primitive,
	}
}

func (cfg *Config) Init() {
	primitive.Init(&cfg.Primitive)
}
