/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package controller

import (
	"github.com/gin-gonic/gin"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/user/app"
)

// AddRouterForUserInternalController adds routes for user internal controller to the given router group.
func AddRouterForUserInternalController(
	rg *gin.RouterGroup,
	us app.UserService,
	m middleware.UserMiddleWare,
) {
	ctl := UserInernalController{
		s: us,
		m: m,
	}

	rg.POST("/v1/user/token/verify", m.Write, ctl.VerifyToken)
	rg.GET("/v1/user/:name/platform", m.Write, ctl.GetPlatformUser)

}

// UserInernalController is a struct that holds references to user service and user middleware.
type UserInernalController struct {
	s app.UserService
	m middleware.UserMiddleWare
}

// @Summary  Verify token
// @Description  verify a platform token of user
// @Tags     User
// @Accept   json
// @Param    body  body  tokenVerifyRequest  true  "body of token"
// @Security Internal
// @Success  200  {object}  commonctl.ResponseData
// @Failure  400  token not provided
// @Failure  401  token empty
// @Failure  403  token invalid
// @Failure  500  internal error
// @Router   /v1/user/token/verify [post]
func (ctl *UserInernalController) VerifyToken(ctx *gin.Context) {
	var req tokenVerifyRequest

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, allerror.NewInvalidParam(err.Error()))
		return
	}

	token, perm, err := req.ToCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, allerror.NewInvalidParam(err.Error()))

		return
	}

	if v, err := ctl.s.VerifyToken(token, perm); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, tokenVerifyResp{
			Account: v.Account,
		})
	}
}

// @Summary  GetPlatformUser info
// @Description  Get platform user info
// @Tags     User
// @Accept   json
// @Param    name  path  string  true  "name of the user"
// @Security Internal
// @Success  200  {object}  commonctl.ResponseData
// @Router   /v1/user/{name}/platform [get]
func (ctl *UserInernalController) GetPlatformUser(ctx *gin.Context) {
	username, err := primitive.NewAccount(ctx.Param("name"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, allerror.NewInvalidParam(err.Error()))
		return
	}

	if v, err := ctl.s.GetPlatformUserInfo(username); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, v)
	}
}
