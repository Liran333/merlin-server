package server

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/config"
	"github.com/openmerlin/merlin-server/other/controller"
)

func setRouterOfOther(rg *gin.RouterGroup, cfg *config.Config) {
	controller.AddRouterForOtherController(rg, cfg)
}
