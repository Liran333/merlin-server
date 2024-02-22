package spaceapp

import (
	"github.com/openmerlin/merlin-server/spaceapp/infrastructure/messageadapter"
	"github.com/openmerlin/merlin-server/spaceapp/infrastructure/repositoryadapter"
)

// Config
type Config struct {
	Tables repositoryadapter.Tables `json:"tables"`
	Topics messageadapter.Topics    `json:"topics"`
}

func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Tables,
		&cfg.Topics,
	}
}
