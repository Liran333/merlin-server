/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// The server package provides functionality for setting up and running the server.
package server

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/config"
	"github.com/openmerlin/merlin-server/other/controller"
)

func setRouterOfOther(rg *gin.RouterGroup, cfg *config.Config) {
	controller.AddRouterForOtherController(rg, cfg)
}
