package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/user/app"
	"github.com/openmerlin/merlin-server/user/domain"
	userrepo "github.com/openmerlin/merlin-server/user/domain/repository"
)

func AddRouterForUserController(
	rg *gin.RouterGroup,
	us app.UserService,
	repo userrepo.User,
	m middleware.UserMiddleWare,
) {
	ctl := UserController{
		repo: repo,
		s:    us,
		m:    m,
	}

	rg.PUT("/v1/user", m.Write, ctl.Update)
	rg.GET("/v1/user", m.Read, ctl.Get)

	rg.POST("/v1/user/token", m.Write, ctl.CreatePlatformToken)
	rg.DELETE("/v1/user/token/:name", m.Write, ctl.DeletePlatformToken)
	rg.GET("/v1/user/token", m.Read, ctl.GetTokenInfo)
}

type UserController struct {
	repo userrepo.User
	s    app.UserService
	m    middleware.UserMiddleWare
}

// @Summary  Update
// @Description  update user basic info
// @Tags     User
// @Param    body  body  userBasicInfoUpdateRequest  true  "body of updating user"
// @Accept   json
// @Security Bearer
// @Success  202   {object}  commonctl.ResponseData
// @Router   /v1/user [put]
func (ctl *UserController) Update(ctx *gin.Context) {
	var req userBasicInfoUpdateRequest

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	//prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "update user basic info")

	user := ctl.m.GetUserAndExitIfFailed(ctx)

	if u, err := ctl.s.UpdateBasicInfo(user, cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, u)
	}
}

// @Summary  Get current user info
// @Description  get current sign-in user info
// @Tags     User
// @Accept   json
// @Success  200  {object}      commonctl.ResponseData
// @Security Bearer
// @Router   /v1/user [get]
func (ctl *UserController) Get(ctx *gin.Context) {
	user := ctl.m.GetUserAndExitIfFailed(ctx)

	// get user own info
	if u, err := ctl.s.UserInfo(user); err != nil {
		logrus.Error(err)

		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, u)
	}
}

// @Summary  DeletePlatformToken
// @Description  delete a new platform token of user
// @Tags     User
// @Param    name  path  string  true  "token name"
// @Accept   json
// @Success  204
// @Security Bearer
// @Router   /v1/user/token/{name} [delete]
func (ctl *UserController) DeletePlatformToken(ctx *gin.Context) {
	user := ctl.m.GetUser(ctx)

	platform, err := ctl.s.GetPlatformUser(user)
	if err != nil {
		commonctl.SendError(ctx, err)

		return
	}

	name, err := primitive.NewTokenName(ctx.Param("name"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, fmt.Errorf("invalid token name"))
	}

	err = ctl.s.DeleteToken(
		&domain.TokenDeletedCmd{
			Account: user,
			Name:    name,
		},
		platform,
	)

	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfDelete(ctx)
	}
}

// @Summary  CreatePlatformToken
// @Description  create a new platform token of user
// @Tags     User
// @Param    body  body  tokenCreateRequest  true  "body of create token"
// @Accept   json
// @Security Bearer
// @Success  201  {object}  commonctl.ResponseData
// @Router   /v1/user/token [post]
func (ctl *UserController) CreatePlatformToken(ctx *gin.Context) {
	var req tokenCreateRequest

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	user := ctl.m.GetUserAndExitIfFailed(ctx)

	cmd, err := req.toCmd(user)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	pl, err := ctl.s.GetPlatformUser(user)
	if err != nil {
		commonctl.SendError(ctx, err)

		return
	}

	if v, err := ctl.s.CreateToken(&cmd, pl); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, v)
	}
}

// @Summary  GetTokenInfo
// @Description  list all platform tokens of user
// @Tags     User
// @Accept   json
// @Security Bearer
// @Success  200  {object}  commonctl.ResponseData
// @Router   /v1/user/token [get]
func (ctl *UserController) GetTokenInfo(ctx *gin.Context) {
	if v, err := ctl.s.ListTokens(ctl.m.GetUserAndExitIfFailed(ctx)); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, v)
	}
}
