package discussion

import (
	"github.com/openmerlin/merlin-server/discussion/domain/primitive"
	"github.com/openmerlin/merlin-server/discussion/infrastructure/emailimpl"
	"github.com/openmerlin/merlin-server/discussion/infrastructure/messageimpl"
	"github.com/openmerlin/merlin-server/discussion/infrastructure/repositoryimpl"
)

type Config struct {
	Tables    repositoryimpl.Tables `json:"tables"`
	Primitive primitive.Config      `json:"primitive"`
	Topics    messageimpl.Topics    `json:"topics"`
	Report    emailimpl.Config      `json:"report"`
}

func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Tables,
		&cfg.Primitive,
		&cfg.Topics,
		&cfg.Report,
	}
}

func (cfg *Config) SetDefault() {
	cfg.Primitive.SetDefault()
}

func (cfg *Config) Init() {
	primitive.InitConfig(&cfg.Primitive)
}
