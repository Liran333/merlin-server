package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	login "github.com/openmerlin/merlin-server/login/domain"
	session "github.com/openmerlin/merlin-server/session/app"

	userapp "github.com/openmerlin/merlin-server/user/app"
	"github.com/openmerlin/merlin-server/user/domain"
	"github.com/openmerlin/merlin-server/utils"
)

type oldUserTokenPayload struct {
	Account string `json:"account"`
	Email   string `json:"email"`
}

func (pl *oldUserTokenPayload) DomainAccount() domain.Account {
	a, _ := domain.NewAccount(pl.Account)

	return a
}

func (pl *oldUserTokenPayload) isNotMe(a domain.Account) bool {
	return pl.Account != a.Account()
}

func (pl *oldUserTokenPayload) isMyself(a domain.Account) bool {
	return pl.Account == a.Account()
}

func (pl *oldUserTokenPayload) hasEmail() bool {
	return pl.Email != ""
}

func AddRouterForLoginController(
	rg *gin.RouterGroup,
	us userapp.UserService,
	auth login.User,
	session session.SessionService,
) {
	pc := LoginController{
		auth:    auth,
		us:      us,
		session: session,
	}

	rg.GET("/v1/login", pc.Login)
	rg.GET("/v1/login/:account", pc.Logout)
}

type LoginController struct {
	baseController

	auth    login.User
	us      userapp.UserService
	session session.SessionService
}

// @Title			Login
// @Description	callback of authentication by authing
// @Tags			Login
// @Param			code			query	string	true	"authing code"
// @Param			redirect_uri	query	string	true	"redirect uri"
// @Accept			json
// @Success		200	{object}			app.UserDTO
// @Failure		500	system_error		system	error
// @Failure		501	duplicate_creating	create	user	repeatedly	which	should	not	happen
// @Router			/v1/login [get]
func (ctl *LoginController) Login(ctx *gin.Context) {
	info, err := ctl.auth.GetByCode(
		ctl.getQueryParameter(ctx, "code"),
		ctl.getQueryParameter(ctx, "redirect_uri"),
	)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseCodeError(
			errorSystemError, err,
		))

		return
	}
	defer utils.ClearStringMemory(info.AccessToken)

	user, err := ctl.us.GetByAccount(info.Name)
	if err != nil {
		if d := newResponseError(err); d.Code != errorResourceNotExists {
			ctl.sendRespWithInternalError(ctx, d)

			return
		}

		if user, err = ctl.newUser(ctx, info); err != nil {
			utils.DoLog(user.Id, user.Account, "logup", "", "failed")

			return
		}

		utils.DoLog(user.Id, user.Account, "logup", "", "success")
	}

	prepareOperateLog(ctx, user.Account, OPERATE_TYPE_USER, "user login")

	if err := ctl.newLogin(ctx, info); err != nil {
		return
	}

	payload := oldUserTokenPayload{
		Account: user.Account,
		Email:   user.Email,
	}

	token, csrftoken, err := ctl.newApiToken(ctx, payload)
	if err != nil {
		ctl.sendRespWithInternalError(
			ctx, newResponseCodeError(errorSystemError, err),
		)

		return
	}

	if err = ctl.setRespToken(ctx, token, csrftoken, user.Account); err != nil {
		ctl.sendRespWithInternalError(
			ctx, newResponseCodeError(errorSystemError, err),
		)

		return
	}

	utils.DoLog(user.Id, user.Account, "login", "", "success")

	ctx.JSON(http.StatusOK, newResponseData(user))
}

func (ctl *LoginController) newLogin(ctx *gin.Context, info login.Login) (err error) {
	idToken, err := ctl.encryptData(info.IDToken)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseCodeError(
			errorSystemError, err,
		))

		return
	}

	err = ctl.session.Create(&session.SessionCreateCmd{
		Account: info.Name,
		Email:   info.Email,
		Info:    idToken,
		UserId:  info.UserId,
	})
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	}

	return
}

func (ctl *LoginController) newUser(ctx *gin.Context, info login.Login) (user userapp.UserDTO, err error) {
	cmd := domain.UserCreateCmd{
		Email:    info.Email,
		Account:  info.Name,
		Bio:      info.Bio,
		AvatarId: info.AvatarId,
	}

	if user, err = ctl.us.Create(&cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	return
}

// @Title			Logout
// @Description	get info of login
// @Tags			Login
// @Param			account	path	string	true	"account"
// @Accept			json
// @Success		200	{object}			session.SessionDTO
// @Failure		400	bad_request_param	account	is	invalid
// @Failure		401	not_allowed			can't	get	login	info	of	other	user
// @Failure		500	system_error		system	error
// @Router			/v1/logout [get]
func (ctl *LoginController) Logout(ctx *gin.Context) {
	account, err := domain.NewAccount(ctx.Param("account"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	pl, _, ok := ctl.checkUserApiTokenNoRefresh(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "user logout")

	if pl.isNotMe(account) {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorNotAllowed,
			"can't get login info of other user",
		))

		return
	}
	// TODO: need to delete session info
	info, err := ctl.session.Get(account)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		utils.DoLog(info.UserId, "", "logout", "", "failed")

		return
	}

	v, err := ctl.decryptData(info.Info)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseCodeError(
			errorSystemError, err,
		))

		return
	}

	utils.DoLog(info.UserId, "", "logout", "", "success")

	ctl.cleanCookie(ctx)

	info.Info = string(v)
	ctx.JSON(http.StatusOK, newResponseData(info))
}
