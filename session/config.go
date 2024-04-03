/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package session provides an config for the session repository,
package session

import (
	"github.com/openmerlin/merlin-server/common/domain/crypto"
	"github.com/openmerlin/merlin-server/session/controller"
	"github.com/openmerlin/merlin-server/session/domain"
	"github.com/openmerlin/merlin-server/session/infrastructure/loginrepositoryadapter"
	"github.com/openmerlin/merlin-server/session/infrastructure/oidcimpl"
)

// Config is a struct that represents the overall configuration for the application.
type Config struct {
	OIDC       oidcimpl.Config               `json:"oidc"`
	Login      loginrepositoryadapter.Tables `json:"login"`
	Domain     domain.Config                 `json:"domain"`
	Controller controller.Config             `json:"controller"`
}

// ConfigItems returns a slice of interface{} containing pointers to the configuration items in the Config struct.
func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.OIDC,
		&cfg.Login,
		&cfg.Domain,
		&cfg.Controller,
	}
}

// Init initializes the application using the configuration settings provided in the Config struct.
func (cfg *Config) Init() error {
	domain.Init(&cfg.Domain)
	oidcimpl.Init(&cfg.OIDC)
	controller.Init(&cfg.Controller)

	if err := loginrepositoryadapter.Init(&cfg.Login, crypto.NewEncryption(cfg.Login.Key)); err != nil {
		return err
	}

	return nil
}
