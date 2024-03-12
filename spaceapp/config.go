/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package spaceapp

import (
	"github.com/openmerlin/merlin-server/spaceapp/controller"
	"github.com/openmerlin/merlin-server/spaceapp/domain"
	"github.com/openmerlin/merlin-server/spaceapp/infrastructure/messageadapter"
	"github.com/openmerlin/merlin-server/spaceapp/infrastructure/repositoryadapter"
)

// Config is a struct that holds the configuration for tables and topics.
type Config struct {
	Controller controller.Config        `json:"controller"`
	Domain     domain.Config            `json:"domain"`
	Tables     repositoryadapter.Tables `json:"tables"`
	Topics     messageadapter.Topics    `json:"topics"`
}

// ConfigItems returns a slice of interfaces containing references to the Tables and Topics fields of the Config struct.
func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Controller,
		&cfg.Tables,
		&cfg.Topics,
	}
}

func (cfg *Config) Init() {
	controller.Init(&cfg.Controller)
	domain.Init(&cfg.Domain)
}
