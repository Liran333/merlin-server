package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	login "github.com/openmerlin/merlin-server/login/domain"
	userapp "github.com/openmerlin/merlin-server/user/app"
	"github.com/openmerlin/merlin-server/user/domain"

	userrepo "github.com/openmerlin/merlin-server/user/domain/repository"
	userinfra "github.com/openmerlin/merlin-server/user/infrastructure/git"
)

func AddRouterForUserController(
	rg *gin.RouterGroup,
	us userapp.UserService,
	repo userrepo.User,
	auth login.User,
) {
	ctl := UserController{
		auth: auth,
		repo: repo,
		s:    us,
	}

	rg.PUT("/v1/user", ctl.Update)
	rg.GET("/v1/user", ctl.Get)

	rg.DELETE("/v1/user/token/:name", checkUserEmailMiddleware(&ctl.baseController), ctl.DeletePlatformToken)
	rg.GET("/v1/user/token", checkUserEmailMiddleware(&ctl.baseController), ctl.GetTokenInfo)
	rg.POST("/v1/user/token", checkUserEmailMiddleware(&ctl.baseController), ctl.CreatePlatformToken)
}

type UserController struct {
	baseController

	repo userrepo.User
	auth login.User
	s    userapp.UserService
}

// @Summary		Update
// @Description	update user basic info
// @Tags			User
// @Param			body	body	userBasicInfoUpdateRequest	true	"body of updating user"
// @Accept			json
// @Produce		json
// @Router			/v1/user [put]
func (ctl *UserController) Update(ctx *gin.Context) {
	m := userBasicInfoUpdateRequest{}

	if err := ctx.ShouldBindJSON(&m); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	cmd, err := m.toCmd()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "update user basic info")

	if err := ctl.s.UpdateBasicInfo(pl.DomainAccount(), cmd); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusAccepted, newResponseData(nil))
}

// @Summary		Get
// @Description	get user
// @Tags			User
// @Param			account	query	string	false	"account"
// @Accept			json
// @Success		200	{object}			userDetail
// @Failure		400	bad_request_param	account	is		invalid
// @Failure		401	resource_not_exists	user	does	not	exist
// @Failure		500	system_error		system	error
// @Router			/v1/user [get]
func (ctl *UserController) Get(ctx *gin.Context) {
	var target domain.Account

	if account := ctl.getQueryParameter(ctx, "account"); account != "" {
		v, err := primitive.NewAccount(account)
		if err != nil {
			logrus.Error(err)

			ctx.JSON(http.StatusBadRequest, newResponseCodeError(
				errorBadRequestParam, err,
			))

			return
		}

		target = v
	}

	pl, visitor, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		logrus.Errorln("failed to get user info")
		return
	}

	resp := func(u *userapp.UserDTO, points int, isFollower bool) {
		ctx.JSON(http.StatusOK, newResponseData(
			userDetail{
				UserDTO: u,
			}),
		)
	}

	if visitor {
		if target == nil {
			logrus.Error("got invalid user info")
			// clear cookie if we got an invalid user info
			ctl.cleanCookie(ctx)

			ctx.JSON(http.StatusOK, newResponseData(nil))
			return
		}

		// get by visitor
		if u, err := ctl.s.GetByAccount(target); err != nil {
			logrus.Errorf("get by visitor err: %s", err)
			ctl.sendRespWithInternalError(ctx, newResponseError(fmt.Errorf("failed to get user(%s) info", target.Account())))
		} else {
			u.Email = ""
			//u.Password = ""
			resp(&u, 0, false)
		}

		return
	}

	if target != nil && pl.isNotMe(target) {
		// get by follower, and pl.Account is follower
		if u, isFollower, err := ctl.s.GetByFollower(target, pl.DomainAccount()); err != nil {
			logrus.Error(err)

			ctl.sendRespWithInternalError(ctx, newResponseError(err))
		} else {
			u.Email = ""
			//u.Password = ""
			resp(&u, 0, isFollower)
		}

		return
	}

	// get user own info
	if u, err := ctl.s.UserInfo(pl.DomainAccount()); err != nil {
		logrus.Error(err)

		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		resp(&u.UserDTO, u.Points, true)
	}
}

// @Title			CheckEmail
// @Description	check user email
// @Tags			User
// @Accept			json
// @Success		200
// @Failure		400	no	email	this	api	need	email	of	user"
// @Router			/v1/user/check_email [get]
func (ctl *UserController) CheckEmail(ctx *gin.Context) {
	ctl.sendRespOfGet(ctx, "")
}

// @Title			DeletePlatformToken
// @Description	delete a new platform token of user
// @Tags			User
// @Param			name	query	string	false	"name"
// @Accept			json
// @Success		204
// @Failure		400	bad_request_param	token	is	invalid
// @Failure		401	not_allowed			can't	get	info	of	other	user
// @Failure		404	not_found			no such token
// @Router			/v1/user/token/{name} [delete]
func (ctl *UserController) DeletePlatformToken(ctx *gin.Context) {
	name := ctx.Param("name")

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	platform, err := ctl.s.GetPlatformUser(pl.DomainAccount())
	if err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusInternalServerError, newResponseCodeError(
			errorSystemError, fmt.Errorf("failed to init platform client"),
		))

		return
	}
	err = ctl.s.DeleteToken(&domain.TokenDeletedCmd{
		Account: pl.DomainAccount(),
		Name:    name,
	}, platform)

	if err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusInternalServerError, newResponseCodeError(
			errorSystemError, err,
		))

		return
	} else {
		ctx.JSON(http.StatusNoContent, newResponseData(nil))
	}
}

// @Title			CreatePlatformToken
// @Description	create a new platform token of user
// @Tags			User
// @Param			body	body	tokenCreateRequest	true	"body of create token"
// @Accept			json
// @Success		201 created userapp.TokenDTO
// @Failure		400	bad_request_param	account	is	invalid
// @Failure		401	not_allowed			can't	get	info	of	other	user
// @Router			/v1/user/token [post]
func (ctl *UserController) CreatePlatformToken(ctx *gin.Context) {
	r := tokenCreateRequest{}

	if err := ctx.ShouldBindJSON(&r); err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	cmd := domain.TokenCreatedCmd{
		Account:    pl.DomainAccount(),
		Name:       r.Name,
		Permission: domain.TokenPerm(r.Perm),
	}

	usernew, err := ctl.s.GetByAccount(pl.DomainAccount())
	if err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusInternalServerError, newResponseCodeError(
			errorSystemError, err,
		))

		return
	}

	platform, err := userinfra.NewBaseAuthClient(
		usernew.Account,
		usernew.Password,
	)
	if err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusInternalServerError, newResponseCodeError(
			errorSystemError, fmt.Errorf("failed to init platform client"),
		))

		return
	}
	createdToken, err := ctl.s.CreateToken(&cmd, platform)
	if err != nil {
		logrus.Errorf("failed to create token %s", err)
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorSystemError,
			"can't refresh token",
		))
		return
	}
	// create new token
	f := func() (token, csrftoken string) {

		if err != nil {
			logrus.Error(err)

			return
		}

		payload := oldUserTokenPayload{
			Account: usernew.Account,
			Email:   usernew.Email,
		}

		token, csrftoken, err = ctl.newApiToken(ctx, payload)
		if err != nil {
			logrus.Error(err)

			return
		}

		return
	}

	token, csrftoken := f()

	if token != "" {
		if err = ctl.setRespToken(ctx, token, csrftoken, usernew.Account); err != nil {
			logrus.Error(err)

			return
		}
	}

	ctl.sendRespOfPost(ctx, createdToken)
}

// @Title			ListUserTokens
// @Description	list all platform tokens of user
// @Tags			User
// @Accept			json
// @Success		200	{object}			[]userapp.TokenDTO
// @Failure		400	bad_request_param	account	is	invalid
// @Failure		401	not_allowed			can't	get	info	of	other	user
// @Router			/v1/user/token [get]
func (ctl *UserController) GetTokenInfo(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	usernew, err := ctl.s.GetByAccount(pl.DomainAccount())
	if err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusInternalServerError, newResponseCodeMsg(
			errorNotAllowed,
			fmt.Sprintf("can't get token of user %s ", pl.Account),
		))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(usernew.Tokens))
}
