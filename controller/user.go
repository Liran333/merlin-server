package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	login "github.com/openmerlin/merlin-server/login/domain"
	userapp "github.com/openmerlin/merlin-server/user/app"
	"github.com/openmerlin/merlin-server/user/domain"
	userrepo "github.com/openmerlin/merlin-server/user/domain/repository"
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

	ctx.JSON(http.StatusAccepted, newResponseData(m))
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
		v, err := domain.NewAccount(account)
		if err != nil {
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
				UserDTO:    u,
				Points:     points,
				IsFollower: isFollower,
			}),
		)
	}

	if visitor {
		if target == nil {
			// clear cookie if we got an invalid user info
			ctl.cleanCookie(ctx)

			ctx.JSON(http.StatusOK, newResponseData(nil))
			return
		}

		// get by visitor
		if u, err := ctl.s.GetByAccount(target); err != nil {
			ctl.sendRespWithInternalError(ctx, newResponseError(err))
		} else {
			u.Email = ""
			resp(&u, 0, false)
		}

		return
	}

	if target != nil && pl.isNotMe(target) {
		// get by follower, and pl.Account is follower
		if u, isFollower, err := ctl.s.GetByFollower(target, pl.DomainAccount()); err != nil {
			ctl.sendRespWithInternalError(ctx, newResponseError(err))
		} else {
			u.Email = ""
			resp(&u, 0, isFollower)
		}

		return
	}

	// get user own info
	if u, err := ctl.s.UserInfo(pl.DomainAccount()); err != nil {
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
