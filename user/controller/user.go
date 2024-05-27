/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	orgapp "github.com/openmerlin/merlin-server/organization/app"
	"github.com/openmerlin/merlin-server/user/app"
	"github.com/openmerlin/merlin-server/user/domain"
	userrepo "github.com/openmerlin/merlin-server/user/domain/repository"
)

// PrivacyClear is an interface for the privacy clearing operation.
type PrivacyClear interface {
	ClearCookieAfterRevokePrivacy(*gin.Context)
}

// AddRouterForUserController adds routes for user-related operations to the given router group.
func AddRouterForUserController(
	rg *gin.RouterGroup,
	us app.UserService,
	repo userrepo.User,
	l middleware.OperationLog,
	sl middleware.SecurityLog,
	m middleware.UserMiddleWare,
	rl middleware.RateLimiter,
	p middleware.PrivacyCheck,
	d orgapp.PrivilegeOrg,
	c PrivacyClear,
) {
	ctl := UserController{
		repo:    repo,
		s:       us,
		m:       m,
		disable: d,
		clear:   c,
	}

	rg.PUT("/v1/user", m.Write, l.Write, rl.CheckLimit, p.Check, ctl.Update)
	rg.GET("/v1/user", m.Read, rl.CheckLimit, p.Check, ctl.Get)
	rg.DELETE("/v1/user", m.Write, l.Write, rl.CheckLimit, ctl.RequestDelete)

	rg.POST("/v1/user/token", m.Write,
		CheckMail(ctl.m, ctl.s, sl), l.Write, rl.CheckLimit, p.Check, ctl.CreatePlatformToken)
	rg.DELETE("/v1/user/token/:name", m.Write,
		CheckMail(ctl.m, ctl.s, sl), l.Write, rl.CheckLimit, p.Check, ctl.DeletePlatformToken)
	rg.GET("/v1/user/token", m.Read, rl.CheckLimit, p.Check, ctl.GetTokenInfo)

	rg.POST("/v1/user/email/bind", m.Write, l.Write, rl.CheckLimit, p.Check, ctl.BindEmail)
	rg.POST("/v1/user/email/send", m.Write, l.Write, rl.CheckLimit, p.Check, ctl.SendEmail)

	rg.PUT("/v1/user/privacy", m.Write, l.Write, rl.CheckLimit, p.Check, ctl.PrivacyRevoke)
}

// UserController is a struct that holds references to user repository, user service, and user middleware.
type UserController struct {
	repo    userrepo.User
	s       app.UserService
	m       middleware.UserMiddleWare
	disable orgapp.PrivilegeOrg
	clear   PrivacyClear
}

// @Summary  Update
// @Description  update user basic info
// @Tags     User
// @Param    body  body  userBasicInfoUpdateRequest  true  "body of updating user"
// @Accept   json
// @Security Bearer
// @Success  202   {object}  commonctl.ResponseData{data=app.UserDTO,msg=string,code=string}
// @Router   /v1/user [put]
func (ctl *UserController) Update(ctx *gin.Context) {
	middleware.SetAction(ctx, "update user basic info")

	var req userBasicInfoUpdateRequest

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, xerrors.Errorf("failed to parse request body: %w", err))

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to parse request param: %w", err))

		return
	}

	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	middleware.SetAction(ctx, fmt.Sprintf("update user:%s's basic info", user.Account()))

	if req.RevokeDelete != nil {
		middleware.SetAction(ctx, fmt.Sprintf("revoke user:%s's delete request", user.Account()))
	}

	if u, err := ctl.s.UpdateBasicInfo(user, cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, u)
	}
}

// CheckMail check mail middleware
func CheckMail(m middleware.UserMiddleWare, us app.UserService, securityLog middleware.SecurityLog) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := m.GetUserAndExitIfFailed(ctx)
		if user == nil {
			return
		}

		u, err := us.GetByAccount(user, user)
		if err != nil {
			securityLog.Warn(ctx, err.Error())
			_ = ctx.AbortWithError(http.StatusBadRequest, xerrors.Errorf("failed to get user info: %w", err))
		} else {
			if u.Email != nil && *u.Email != "" {
				ctx.Next()
			} else {
				// will call ctx.Abort() internally
				e := xerrors.Errorf("need bind user email firstly")
				err := allerror.New(allerror.ErrorCodeNeedBindEmail, e.Error(), e)
				securityLog.Warn(ctx, err.Error())
				commonctl.SendError(ctx, err)
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
// @Success  201   {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Router   /v1/user/email/bind [post]
func (ctl *UserController) BindEmail(ctx *gin.Context) {
	middleware.SetAction(ctx, "bind email")

	var req bindEmailRequest
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, xerrors.Errorf("failed to parse request body: %w", err))

		return
	}

	middleware.SetAction(ctx, req.action())

	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	cmd, err := req.toCmd(user)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to parse request param: %w", err))

		return
	}

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
// @Success  201   {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Router   /v1/user/email/send [post]
func (ctl *UserController) SendEmail(ctx *gin.Context) {
	middleware.SetAction(ctx, "send email verify code")

	var req sendEmailRequest
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, xerrors.Errorf("failed to parse request body: %w", err))

		return
	}

	middleware.SetAction(ctx, req.action())

	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	cmd, err := req.toCmd(user)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to parse request param: %w", err))

		return
	}

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
// @Success  200  {object}   commonctl.ResponseData{data=app.UserInfoDTO,msg=string,code=string}
// @Security Bearer
// @Router   /v1/user [get]
func (ctl *UserController) Get(ctx *gin.Context) {
	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	// get user own info
	u, err := ctl.s.UserInfo(user, user)
	if err != nil {
		commonctl.SendError(ctx, err)

		return
	}

	if ctl.disable != nil {
		err = ctl.disable.Contains(user)
		if err != nil {
			u.IsDisableAdmin = false
		} else {
			u.IsDisableAdmin = true
		}
	}
	commonctl.SendRespOfGet(ctx, u)
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
	middleware.SetAction(ctx, "delete token")

	user := ctl.m.GetUser(ctx)

	platform, err := ctl.s.GetPlatformUser(user)
	if err != nil {
		commonctl.SendError(ctx, err)

		return
	}

	name, err := primitive.NewTokenName(ctx.Param("name"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("invalid token name: %w", err))
		return
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
// @Success  201  {object}  commonctl.ResponseData{data=app.TokenDTO,msg=string,code=string}
// @Router   /v1/user/token [post]
func (ctl *UserController) CreatePlatformToken(ctx *gin.Context) {
	middleware.SetAction(ctx, "create a new platform token")

	var req tokenCreateRequest

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, xerrors.Errorf("failed to parse request body: %w", err))

		return
	}

	middleware.SetAction(ctx, req.action())

	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	cmd, err := req.toCmd(user)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to parse request param: %w", err))

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
// @Success  200  {object}  commonctl.ResponseData{data=[]app.TokenDTO,msg=string,code=string}
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
// @Success  202   {object}  commonctl.ResponseData{data=revokePrivacyInfo,msg=string,code=string}
// @Router   /v1/user/privacy [put]
func (ctl *UserController) PrivacyRevoke(ctx *gin.Context) {
	middleware.SetAction(ctx, "privacy revoke")

	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	if idToken, err := ctl.s.PrivacyRevoke(user); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		ctl.clear.ClearCookieAfterRevokePrivacy(ctx)
		commonctl.SendRespOfPut(ctx, revokePrivacyInfo{IdToken: idToken})
	}
}

// @Summary  RequestDelete User info
// @Description  delete
// @Tags     User
// @Accept   json
// @Security Bearer
// @Success  204
// @Router   /v1/user [delete]
func (ctl *UserController) RequestDelete(ctx *gin.Context) {
	middleware.SetAction(ctx, "request delete user")

	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	middleware.SetAction(ctx, fmt.Sprintf("request delete user:%s's info", user.Account()))

	if err := ctl.s.RequestDelete(user); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfDelete(ctx)
	}
}
