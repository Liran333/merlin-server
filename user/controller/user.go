package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/user/app"
	"github.com/openmerlin/merlin-server/user/domain"
	userrepo "github.com/openmerlin/merlin-server/user/domain/repository"
)

func AddRouterForUserController(
	rg *gin.RouterGroup,
	us app.UserService,
	repo userrepo.User,
	l middleware.OperationLog,
	m middleware.UserMiddleWare,
) {
	ctl := UserController{
		repo: repo,
		s:    us,
		m:    m,
	}

	rg.PUT("/v1/user", m.Write, l.Write, ctl.Update)
	rg.GET("/v1/user", m.Read, ctl.Get)

	rg.POST("/v1/user/token", m.Write, CheckMail(ctl.m, ctl.s), l.Write, ctl.CreatePlatformToken)
	rg.DELETE("/v1/user/token/:name", m.Write, CheckMail(ctl.m, ctl.s), l.Write, ctl.DeletePlatformToken)
	rg.GET("/v1/user/token", m.Read, ctl.GetTokenInfo)

	rg.POST("/v1/user/email/bind", m.Write, l.Write, ctl.BindEmail)
	rg.POST("/v1/user/email/send", m.Write, l.Write, ctl.SendEmail)

	rg.PUT("/v1/user/privacy", m.Write, l.Write, ctl.PrivacyRevoke)
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

	middleware.SetAction(ctx, "update user basic info")
	//prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "update user basic info")

	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	if u, err := ctl.s.UpdateBasicInfo(user, cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, u)
	}
}

// check mail middleware
func CheckMail(m middleware.UserMiddleWare, us app.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := m.GetUserAndExitIfFailed(ctx)
		if user == nil {
			return
		}

		u, err := us.GetByAccount(user, user)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("user not found"))
		} else {
			if u.Email != nil && *u.Email != "" {
				ctx.Next()
			} else {
				// will call ctx.Abort() internally
				commonctl.SendError(ctx, allerror.New(allerror.ErrorCodeNeedBindEmail, "need bind user email firstly"))
			}
		}
	}
}

// @Summary  Bind User Email
// @Description  bind user's email
// @Tags     User
// @Param    body  body  bindEmailRequest  true  "body of bind email info"
// @Accept   json
// @Security Bearer
// @Success  201   {object}  commonctl.ResponseData
// @Router   /v1/user/email/bind [post]
func (ctl *UserController) BindEmail(ctx *gin.Context) {
	var req bindEmailRequest
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	cmd, err := req.toCmd(user)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	middleware.SetAction(ctx, req.action())
	//prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "update user basic info")

	if err := ctl.s.VerifyBindEmail(&cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, nil)
	}
}

// @Summary  Send User Email Verify code
// @Description  send user's email verify code
// @Tags     User
// @Param    body  body  sendEmailRequest  true  "body of bind email info"
// @Accept   json
// @Security Bearer
// @Success  201   {object}  commonctl.ResponseData
// @Router   /v1/user/email/send [post]
func (ctl *UserController) SendEmail(ctx *gin.Context) {
	var req sendEmailRequest
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	cmd, err := req.toCmd(user)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	middleware.SetAction(ctx, req.action())
	//prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "update user basic info")

	if err := ctl.s.SendBindEmail(&cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, nil)
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
	if user == nil {
		return
	}

	// get user own info
	if u, err := ctl.s.UserInfo(user, user); err != nil {
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

	middleware.SetAction(ctx, fmt.Sprintf("delete token %s", name))

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

	middleware.SetAction(ctx, req.action())

	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

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
	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	if v, err := ctl.s.ListTokens(user); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, v)
	}
}

// @Summary  PrivacyRevoke
// @Description  revoke
// @Tags     User
// @Accept   json
// @Security Bearer
// @Success  202   {object}  commonctl.ResponseData
// @Router   /v1/user/privacy [put]
func (ctl *UserController) PrivacyRevoke(ctx *gin.Context) {
	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	middleware.SetAction(ctx, "privacy revoke")

	if err := ctl.s.PrivacyRevoke(user); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}
