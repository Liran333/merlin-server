package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/organization/domain"

	orgapp "github.com/openmerlin/merlin-server/organization/app"
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
	rg.DELETE("/v1/organization/:name", ctl.Delete)

	rg.POST("/v1/organization/:name/invite", ctl.InviteMember)
	rg.GET("/v1/organization/:name/invite", ctl.ListInvitation)
	rg.DELETE("/v1/organization/:name/invite", ctl.RemoveInvitation)

	rg.DELETE("/v1/organization/:name/member", ctl.RemoveMember)
	rg.GET("/v1/organization/:name/member", ctl.ListMember)
	rg.POST("/v1/organization/:name/member", ctl.AddMember)
}

type OrgController struct {
	baseController

	org orgapp.OrgService
}

// @Summary		Update
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

	o, err := ctl.org.UpdateBasicInfo(orgName, &cmd)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusAccepted, newResponseData(o))
}

// @Summary		Get one organization info
// @Description	get organization info
// @Tags			Organization
// @Param			name	path	string	true	"name"
// @Accept			json
// @Success		200	{object}			userDetail
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

		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusOK, newResponseData(o))
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

		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusOK, newResponseData(os))
	}
}

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
	})
	if err != nil {
		logrus.Errorf("create org failed: %s", err)
		ctx.JSON(http.StatusInternalServerError, newResponseCodeError(
			errorSystemError, err,
		))

		return
	} else {
		ctx.JSON(http.StatusCreated, newResponseData(o))
	}
}

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

	err = ctl.org.Delete(orgName)
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

// @Title			Add organization members
// @Description Add a member to the organization
// @Tags			Organization
// @Param			body	body	orgCreateRequest	true	"body of new organization"
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

	err = ctl.org.AddMember(&orgapp.OrgAddMemberCmd{
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

// @Title			Remove organization members
// @Description Remove a member from a organization
// @Tags			Organization
// @Param			body	body	orgCreateRequest	true	"body of new organization"
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

	err = ctl.org.RemoveMember(&orgapp.OrgRemoveMemberCmd{
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

		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, fmt.Errorf("user name not valid"),
		))

		return
	}

	dto, err := ctl.org.InviteMember(&orgapp.OrgInviteMemberCmd{
		Org:     orgName,
		Account: acc,
		Role:    req.Role,
	})
	if err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusInternalServerError, newResponseCodeMsg(
			errorNotAllowed,
			fmt.Sprintf("can't invite user %s to org %s ", acc, orgName),
		))

		return
	}

	ctx.JSON(http.StatusCreated, newResponseData(dto))
}

// @Title			List invitation of the organization
// @Description List invitation of the organization
// @Tags			Organization
// @Param			name	path	string	true	"organization name"
// @Accept			json
// @Success		200 {object} []orgapp.ApproveDTO
// @Failure		400	bad_request_param	account	is	invalid
// @Failure		401	not_allowed			can't	get	info	of	other	user
// @Router			/v1/organization/{name}/invite [post]
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

	dtos, err := ctl.org.ListInvitation(orgName)
	if err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusInternalServerError, newResponseCodeMsg(
			errorNotAllowed,
			fmt.Sprintf("can't get invitation of org %s ", orgName),
		))

		return
	}

	ctx.JSON(http.StatusCreated, newResponseData(dtos))
}

// @Title			Revoke invitation of the organization
// @Description Revoke invitation of the organization
// @Tags			Organization
// @Param			body	body	OrgInviteMemberRequest	true	"body of the invitation"
// @Param			name	path	string	true	"organization name"
// @Accept			json
// @Success		200 {object} []orgapp.ApproveDTO
// @Failure		400	bad_request_param	account	is	invalid
// @Failure		401	not_allowed			can't	get	info	of	other	user
// @Router			/v1/organization/{name}/invite [post]
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

	dto, err := ctl.org.RevokeInvite(&orgapp.OrgRemoveInviteCmd{
		Org:     orgName,
		Account: acc,
	})
	if err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusInternalServerError, newResponseCodeMsg(
			errorNotAllowed,
			fmt.Sprintf("can't get invitation of org %s ", orgName),
		))

		return
	}

	ctx.JSON(http.StatusCreated, newResponseData(dto))
}
