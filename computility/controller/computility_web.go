/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package controller

import (
	"github.com/gin-gonic/gin"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/computility/app"
	"github.com/openmerlin/merlin-server/computility/domain"
)

// AddRouterForComputilityWebController adds routes to the given router group for the AddRouterForComputilityWebController.
func AddRouterForComputilityWebController(
	r *gin.RouterGroup,
	s app.ComputilityAppService,
	m middleware.UserMiddleWare,
	l middleware.OperationLog,
) {
	ctl := ComputilityWebController{
		appService:     s,
		userMiddleWare: m,
	}

	r.GET("/v1/computility/account/:type", l.Write, m.Read, ctl.GetComputilityAccountDetail)
}

// ComputilityWebAppService is a struct that holds the necessary dependencies for handling computility-related operations.
type ComputilityWebController struct {
	appService     app.ComputilityAppService
	userMiddleWare middleware.UserMiddleWare
}

// @Summary  GetComputilityAccountDetail
// @Description  get user computility account detail
// @Tags     ComputilityWeb
// @Param    type   path  string  true  "computility type"
// @Accept   json
// @Success  200  {object} commonctl.ResponseData{data=app.AccountQuotaDetailDTO,msg=string,code=string}
// @Failure  400  {object} commonctl.ResponseData{data=error,msg=string,code=string}
// @Router   /v1/computility/account/{type} [get]
func (ctl *ComputilityWebController) GetComputilityAccountDetail(ctx *gin.Context) {
	user := ctl.userMiddleWare.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	compute, err := primitive.NewComputilityType(ctx.Param("type"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	r, err := ctl.appService.GetAccountDetail(domain.ComputilityAccountIndex{
		UserName:    user,
		ComputeType: compute,
	})
	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, r)
	}
}
