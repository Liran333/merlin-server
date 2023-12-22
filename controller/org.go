package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	commonapp "github.com/openmerlin/merlin-server/common/app"
	"github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/organization/domain"

	orgapp "github.com/openmerlin/merlin-server/organization/app"
	userapp "github.com/openmerlin/merlin-server/user/app"
)

func AddRouterForOrgController(
	rg *gin.RouterGroup,
	org orgapp.OrgService,
) {
	ctl := OrgController{
		org: org,
	}

	rg.PUT("/v1/organization/:name", ctl.Update)
	rg.POST("/v1/organization", ctl.Create)
	rg.GET("/v1/organization/:name", ctl.Get)
	rg.GET("/v1/organization", ctl.List)
	rg.POST("/v1/organization/:name", ctl.Leave)
	rg.DELETE("/v1/organization/:name", ctl.Delete)
	rg.HEAD("/v1/name", ctl.Check)

	rg.POST("/v1/organization/:name/invite", ctl.InviteMember)
	rg.GET("/v1/organization/:name/invite", ctl.ListInvitation)
	rg.DELETE("/v1/organization/:name/invite", ctl.RemoveInvitation)

	rg.DELETE("/v1/organization/:name/member", ctl.RemoveMember)
	rg.GET("/v1/organization/:name/member", ctl.ListMember)
	rg.PUT("/v1/organization/:name/member", ctl.EditMember)
	rg.POST("/v1/organization/:name/member", ctl.AddMember)
	rg.GET("/v1/:name", ctl.GetUser)
}

type OrgController struct {
	baseController

	org  orgapp.OrgService
	user userapp.UserService
}

// @Summary		Update org basic info
// @Description	update org basic info
// @Tags			Organization
// @Param			name	path	string	true	"name"
// @Param			body	body	orgBasicInfoUpdateRequest	true	"body of new organization"
// @Accept			json
// @Success		202	{object}			orgapp.OrganizationDTO
// @Failure		400	bad_request_param	account	is		invalid
// @Failure		401	resource_not_exists	user	does	not	exist
// @Failure		500	system_error		system	error
// @Produce		json
// @Router			/v1/organization/{name} [put]
func (ctl *OrgController) Update(ctx *gin.Context) {
	m := orgBasicInfoUpdateRequest{}

	if err := ctx.ShouldBindJSON(&m); err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	orgName, err := primitive.NewAccount(ctx.Param("name"))
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, fmt.Errorf("organization name not valid"),
		))

		return
	}

	cmd, err := m.toCmd()
	if err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "update org basic info")

	cmd.Actor = pl.DomainAccount()
	cmd.OrgName = orgName

	o, err := ctl.org.UpdateBasicInfo(&cmd)
	if err != nil {
		controller.SendError(ctx, err)

		return
	}

	ctx.JSON(http.StatusAccepted, newResponseData(o))
}

// @Summary		Get one organization info
// @Description	get organization info
// @Tags			Organization
// @Param			name	path	string	true	"name"
// @Accept			json
// @Success		200	{object}			orgapp.OrganizationDTO
// @Failure		400	bad_request_param	account	is		invalid
// @Failure		401	resource_not_exists	user	does	not	exist
// @Failure		500	system_error		system	error
// @Router			/v1/organization/{name} [get]
func (ctl *OrgController) Get(ctx *gin.Context) {
	orgName, err := primitive.NewAccount(ctx.Param("name"))
	if err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, fmt.Errorf("organization name not valid"),
		))

		return
	}

	_, _, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		logrus.Errorln("failed to get user info")
		return
	}

	// get org info
	if o, err := ctl.org.GetByAccount(orgName); err != nil {
		logrus.Error(err)

		controller.SendError(ctx, err)
	} else {
		controller.SendRespOfGet(ctx, o)
	}
}

// @Summary		Get one organization or user info
// @Description	get organization or user info
// @Tags			Organization
// @Param			name	path	string	true	"name"
// @Accept			json
// @Success		200	{object}			controller.User
// @Failure		400	bad_request_param	account	is		invalid
// @Failure		401	resource_not_exists	user	does	not	exist
// @Failure		500	system_error		system	error
// @Router			/v1/{name} [get]
func (ctl *OrgController) GetUser(ctx *gin.Context) {
	name, err := primitive.NewAccount(ctx.Param("name"))
	if err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, fmt.Errorf("organization name not valid"),
		))

		return
	}

	_, _, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		logrus.Errorln("failed to get user info")
		return
	}

	// get org info
	if o, err := ctl.org.GetByAccount(name); err != nil {
		if err != nil {
			u, err := ctl.user.GetByAccount(name, false)
			if err != nil {
				logrus.Error(err)

				controller.SendError(ctx, err)
				return
			}

			controller.SendRespOfGet(ctx, commonapp.FromUserDTO(u))
		}
	} else {
		controller.SendRespOfGet(ctx, commonapp.FromOrgDTO(o))
	}
}

// @Summary		Check the name is available
// @Description	 Check the name is available
// @Tags			Name
// @Param			name	query	string	true	"name"
// @Accept			json
// @Success		200	name is valid
// @Failure		409	name is been used
// @Router			/v1/name [head]
func (ctl *OrgController) Check(ctx *gin.Context) {
	name, err := primitive.NewAccount(ctl.getQueryParameter(ctx, "name"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, fmt.Errorf("name invalid"),
		))

		return
	}

	// get org info
	if bool := ctl.org.CheckName(name); bool {
		ctx.JSON(http.StatusOK, newResponseData(nil))
	} else {
		ctx.JSON(http.StatusConflict, newResponseData(nil))
	}
}

// @Summary		Get all organization of the user
// @Description	get organization info
// @Tags			Organization
// @Param			name	path	string	true	"name"
// @Accept			json
// @Success		200	{object}			[]orgapp.OrganizationDTO
// @Failure		400	bad_request_param	account	is		invalid
// @Failure		401	resource_not_exists	user	does	not	exist
// @Failure		500	system_error		system	error
// @Router			/v1/organization [get]
func (ctl *OrgController) List(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		logrus.Errorln("failed to get user info")
		return
	}

	// get org info
	if os, err := ctl.org.GetByUser(pl.DomainAccount()); err != nil {
		logrus.Error(err)

		controller.SendError(ctx, err)
	} else {
		controller.SendRespOfGet(ctx, os)
	}
}

// @Summary			Create organization
// @Title			Create organization
// @Description	create a new organization
// @Tags			Organization
// @Param			body	body	orgCreateRequest	true	"body of new organization"
// @Accept			json
// @Success		201 {object} orgapp.OrganizationDTO
// @Failure		400	bad_request_param	token	is	invalid
// @Failure		401	not_allowed			can't	get	info	of	other	user
// @Failure		404	not_found			no such token
// @Router			/v1/organization [post]
func (ctl *OrgController) Create(ctx *gin.Context) {
	req := orgCreateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
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

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "create org")

	o, err := ctl.org.Create(&domain.OrgCreatedCmd{
		Name:        req.Name,
		Owner:       pl.Account,
		Website:     req.Website,
		AvatarId:    req.AvatarId,
		Description: req.Description,
		FullName:    req.FullName,
	})
	if err != nil {
		logrus.Errorf("create org failed: %s", err)
		controller.SendError(ctx, err)

		return
	} else {
		controller.SendRespOfPost(ctx, o)
	}
}

// @Summary			Delete organization
// @Title			Delete organization
// @Description	delete a organization
// @Tags			Organization
// @Param			name	path	string	true	"name"
// @Accept			json
// @Success		204
// @Failure		400	bad_request_param	token	is	invalid
// @Failure		401	not_allowed			can't	get	info	of	other	user
// @Failure		403	permission denied user not allowed to delete organization
// @Router			/v1/organization/{name} [delete]
func (ctl *OrgController) Delete(ctx *gin.Context) {
	orgName, err := primitive.NewAccount(ctx.Param("name"))
	if err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, fmt.Errorf("organization name not valid"),
		))

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "delete organization")

	err = ctl.org.Delete(&domain.OrgDeletedCmd{
		Actor: pl.DomainAccount(),
		Name:  orgName,
	})
	if err != nil {
		logrus.Error(err)

		controller.SendError(ctx, err)

		return
	} else {
		controller.SendRespOfDelete(ctx)
	}
}

// @Summary			List organization members
// @Title			List organization members
// @Description	list organization members
// @Tags			Organization
// @Param			name	path	string	true	"name"
// @Accept			json
// @Success		200 {object} []orgapp.MemberDTO
// @Failure		401	not_allowed		not a valid session
// @Failure		403	permission denied user not allowed to get organization info
// @Router			/v1/organization/{name}/member [get]
func (ctl *OrgController) ListMember(ctx *gin.Context) {
	_, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	orgName, err := primitive.NewAccount(ctx.Param("name"))
	if err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, fmt.Errorf("organization name not valid"),
		))

		return
	}

	members, err := ctl.org.ListMember(orgName)
	if err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusInternalServerError, newResponseCodeError(
			errorSystemError, err,
		))

		return
	}
	ctx.JSON(http.StatusOK, newResponseData(members))
}

// @Summary			Edit organization member
// @Title			Edit organization member role
// @Description Edit a member to the organization's role
// @Tags			Organization
// @Param			body	body	OrgMemberEditRequest	true	"body of new member"
// @Param			name	path	string	true	"name"
// @Accept			json
// @Success		202 {object} app.MemberDTO
// @Failure		400	bad_request_param	account	is	invalid
// @Failure		401	not_allowed			can't	get	info	of	other	user
// @Router			/v1/organization/{name}/member [post]
func (ctl *OrgController) EditMember(ctx *gin.Context) {
	orgName, err := primitive.NewAccount(ctx.Param("name"))
	if err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, fmt.Errorf("organization name not valid"),
		))

		return
	}

	req := OrgMemberEditRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	acc, err := primitive.NewAccount(req.User)
	if err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, fmt.Errorf("user name not valid"),
		))

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "add a member to organization")

	newMember, err := ctl.org.EditMember(&domain.OrgEditMemberCmd{
		Actor:   pl.DomainAccount(),
		Org:     orgName,
		Role:    req.Role,
		Account: acc,
	})
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, newResponseCodeMsg(
			errorNotAllowed,
			fmt.Sprintf("can't add member %s to organization %s ", acc, orgName),
		))

		return
	}

	ctx.JSON(http.StatusAccepted, newResponseData(newMember))
}

// @Summary			Add organization members
// @Title			Add organization members
// @Description Add a member to the organization, the user must be on invite list before adding
// @Tags			Organization
// @Param			body	body	orgMemberAddRequest	true	"body of new member"
// @Param			name	path	string	true	"name"
// @Accept			json
// @Success		201
// @Failure		400	bad_request_param	account	is	invalid
// @Failure		401	not_allowed			can't	get	info	of	other	user
// @Router			/v1/organization/{name}/member [post]
func (ctl *OrgController) AddMember(ctx *gin.Context) {
	orgName, err := primitive.NewAccount(ctx.Param("name"))
	if err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, fmt.Errorf("organization name not valid"),
		))

		return
	}

	req := orgMemberAddRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	acc, err := primitive.NewAccount(req.User)
	if err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, fmt.Errorf("user name not valid"),
		))

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "add a member to organization")

	err = ctl.org.AddMember(&domain.OrgAddMemberCmd{
		Actor:   pl.DomainAccount(),
		Org:     orgName,
		Account: acc,
	})
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, newResponseCodeMsg(
			errorNotAllowed,
			fmt.Sprintf("can't add member %s to organization %s ", acc, orgName),
		))

		return
	}

	ctx.JSON(http.StatusCreated, newResponseData(nil))
}

// @Summary			Remove organization members
// @Title			Remove organization members
// @Description Remove a member from a organization
// @Tags			Organization
// @Param			body	body	orgMemberRemoveRequest	true	"body of the removed member"
// @Param			name	path	string	true	"name"
// @Accept			json
// @Success		204
// @Failure		400	bad_request_param	account	is	invalid
// @Failure		401	not_allowed			can't	get	info	of	other	user
// @Router			/v1/organization/{name}/member [delete]
func (ctl *OrgController) RemoveMember(ctx *gin.Context) {
	orgName, err := primitive.NewAccount(ctx.Param("name"))
	if err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, fmt.Errorf("organization name not valid"),
		))

		return
	}

	req := orgMemberRemoveRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
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

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "remove a member to organization")

	acc, err := primitive.NewAccount(req.User)
	if err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, fmt.Errorf("user name not valid"),
		))

		return
	}

	err = ctl.org.RemoveMember(&domain.OrgRemoveMemberCmd{
		Actor:   pl.DomainAccount(),
		Org:     orgName,
		Account: acc,
	})
	if err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusInternalServerError, newResponseCodeMsg(
			errorNotAllowed,
			fmt.Sprintf("can't remove organization %s user %s ", orgName, acc),
		))

		return
	}

	ctx.JSON(http.StatusCreated, newResponseData(nil))
}

// @Summary			Leave the organization
// @Title			Leave the organization
// @Description	delete a organization
// @Tags			Organization
// @Param			name	path	string	true	"name"
// @Accept			json
// @Success		204
// @Failure		400	bad_request_param	token	is	invalid
// @Failure		401	not_allowed			can't	get	info	of	other	user
// @Failure		403	permission denied user not allowed to delete organization
// @Router			/v1/organization/{name} [post]
func (ctl *OrgController) Leave(ctx *gin.Context) {
	orgName, err := primitive.NewAccount(ctx.Param("name"))
	if err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, fmt.Errorf("organization name not valid"),
		))

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "leave organization")

	err = ctl.org.RemoveMember(&domain.OrgRemoveMemberCmd{
		Actor:   pl.DomainAccount(),
		Org:     orgName,
		Account: pl.DomainAccount(),
	})
	if err != nil {
		logrus.Error(err)

		controller.SendError(ctx, err)

		return
	}

	controller.SendRespOfDelete(ctx)
}

// @Summary			Invite a user to be a member of the organization
// @Title			Send invitation to a user
// @Description Send invitation to a user to join the organization
// @Tags			Organization
// @Param			body	body	OrgInviteMemberRequest	true	"body of the invitation"
// @Param			name	path	string	true	"name"
// @Accept			json
// @Success		201 {object} orgapp.OrganizationDTO
// @Failure		400	bad_request_param	account	is	invalid
// @Failure		401	not_allowed			can't	get	info	of	other	user
// @Router			/v1/organization/{name}/invite [post]
func (ctl *OrgController) InviteMember(ctx *gin.Context) {
	orgName, err := primitive.NewAccount(ctx.Param("name"))
	if err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, fmt.Errorf("organization name not valid"),
		))

		return
	}

	req := OrgInviteMemberRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
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

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "invite a member to organization")

	acc, err := primitive.NewAccount(req.User)
	if err != nil {
		logrus.Error(err)

		controller.SendError(ctx, err)

		return
	}

	dto, err := ctl.org.InviteMember(&domain.OrgInviteMemberCmd{
		Actor:   pl.DomainAccount(),
		Org:     orgName,
		Account: acc,
		Role:    req.Role,
	})
	if err != nil {
		logrus.Error(err)

		controller.SendError(ctx, err)

		return
	}

	controller.SendRespOfPost(ctx, dto)
}

// @Summary			List invitation of the organization
// @Title			List invitation of the organization
// @Description List invitation of the organization
// @Tags			Organization
// @Param			name	path	string	true	"organization name"
// @Accept			json
// @Success		200 {object} []orgapp.ApproveDTO
// @Failure		400	bad_request_param	account	is	invalid
// @Failure		401	not_allowed			can't	get	info	of	other	user
// @Router			/v1/organization/{name}/invite [get]
func (ctl *OrgController) ListInvitation(ctx *gin.Context) {
	orgName, err := primitive.NewAccount(ctx.Param("name"))
	if err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, fmt.Errorf("organization name not valid"),
		))

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "invite a member to organization")

	dtos, err := ctl.org.ListInvitation(&domain.OrgNormalCmd{
		Actor: pl.DomainAccount(),
		Org:   orgName,
	})
	if err != nil {
		logrus.Error(err)

		controller.SendError(ctx, err)

		return
	}

	controller.SendRespOfGet(ctx, dtos)
}

// @Summary			Revoke invitation of the organization
// @Title			Revoke invitation of the organization
// @Description Revoke invitation of the organization
// @Tags			Organization
// @Param			body	body	OrgRevokeInviteRequest	true	"body of the invitation"
// @Param			name	path	string	true	"organization name"
// @Accept			json
// @Success		200 {object} []orgapp.ApproveDTO
// @Failure		400	bad_request_param	account	is	invalid
// @Failure		401	not_allowed			can't	get	info	of	other	user
// @Router			/v1/organization/{name}/invite [delete]
func (ctl *OrgController) RemoveInvitation(ctx *gin.Context) {
	orgName, err := primitive.NewAccount(ctx.Param("name"))
	if err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, fmt.Errorf("organization name not valid"),
		))

		return
	}

	req := OrgRevokeInviteRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	acc, err := primitive.NewAccount(req.User)
	if err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, fmt.Errorf("user name not valid"),
		))

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "revoke a invitation of organization")

	dto, err := ctl.org.RevokeInvite(&domain.OrgRemoveInviteCmd{
		Actor:   pl.DomainAccount(),
		Org:     orgName,
		Account: acc,
	})
	if err != nil {
		logrus.Error(err)

		controller.SendError(ctx, err)

		return
	}

	controller.SendRespOfPost(ctx, dto)
}
