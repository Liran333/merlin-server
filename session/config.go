package session

import (
	"github.com/openmerlin/merlin-server/session/controller"
	"github.com/openmerlin/merlin-server/session/domain"
	"github.com/openmerlin/merlin-server/session/infrastructure/loginrepositoryadapter"
	"github.com/openmerlin/merlin-server/session/infrastructure/oidcimpl"
)

// Config
type Config struct {
	OIDC       oidcimpl.Config               `json:"oidc"`
	Login      loginrepositoryadapter.Tables `json:"login"`
	Domain     domain.Config                 `json:"domain"`
	Controller controller.Config             `json:"controller"`
}

func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.OIDC,
		&cfg.Login,
		&cfg.Domain,
		&cfg.Controller,
	}
}

func (cfg *Config) Init() error {
	domain.Init(&cfg.Domain)
	oidcimpl.Init(&cfg.OIDC)
	controller.Init(&cfg.Controller)

	if err := loginrepositoryadapter.Init(&cfg.Login); err != nil {
		return err
	}

	return nil
}
