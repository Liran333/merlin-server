package controller

import (
	"github.com/gin-gonic/gin"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/space-app/app"
)

func AddRouteForSpaceAppInternalController(
	r *gin.RouterGroup,
	s app.SpaceAppInternalAppService,
	m middleware.UserMiddleWare,
) {

	ctl := SpaceAppInternalController{
		appService: s,
	}

	r.POST(`/v1/space-app`, m.Write, ctl.Create)
}

type SpaceAppInternalController struct {
	appService app.SpaceAppInternalAppService
}

// @Summary  Create
// @Description  create space app
// @Tags     SpaceApp
// @Param    body  body      reqToCreateSpaceApp  true  "body of creating space app"
// @Accept   json
// @Security Bearer
// @Success  201   {object}  commonctl.ResponseData
// @Router   /v1/space-app/ [post]
func (ctl *SpaceAppInternalController) Create(ctx *gin.Context) {
	req := reqToCreateSpaceApp{}
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if err := ctl.appService.Create(&cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, "successfully")
	}
}
