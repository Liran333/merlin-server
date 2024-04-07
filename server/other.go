package server

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/other"
	"github.com/openmerlin/merlin-server/other/controller"
)

func setRouterOfOther(rg *gin.RouterGroup, cfg *other.Config) {
	controller.AddRouterForOtherController(rg, &cfg.Analyse)
}
