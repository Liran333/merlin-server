/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package user

import (
	"github.com/openmerlin/merlin-server/user/controller"
	"github.com/openmerlin/merlin-server/user/domain"
)

// Config is a struct that holds the configuration for domain and controller.
type Config struct {
	Domain     domain.Config     `json:"domain"`
	Controller controller.Config `json:"controller"`
}

// ConfigItems returns a slice of interfaces containing references to the Domain and Controller fields of the Config struct.
func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Domain,
		&cfg.Controller,
	}
}

// Init initializes the Config struct with default values.
func (cfg *Config) Init() {
	controller.Init(&cfg.Controller)
}
