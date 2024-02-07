package controller

import (
	"github.com/gin-gonic/gin"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/space/app"
)

func AddRouterForSpaceInternalController(
	r *gin.RouterGroup,
	s app.SpaceInternalAppService,
	m middleware.UserMiddleWare,
) {
	ctl := SpaceInternalController{
		appService: s,
	}

	r.GET("/v1/space/:id", m.Write, ctl.Get)
}

type SpaceInternalController struct {
	appService app.SpaceInternalAppService
}

// @Summary  Get
// @Description  get space
// @Tags     SpaceInternal
// @Param    id  path  string  true  "id of space"
// @Accept   json
// @Success  200  {object}  app.SpaceMetaDTO
// @Router   /v2/space/{id} [get]
func (ctl *SpaceInternalController) Get(ctx *gin.Context) {
	spaceId, err := primitive.NewIdentity(ctx.Param("id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if dto, err := ctl.appService.GetById(spaceId); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, &dto)
	}
}
