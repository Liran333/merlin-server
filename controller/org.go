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
	user userapp.UserService,
) {
	ctl := OrgController{
		org:  org,
		user: user,
	}

	rg.PUT("/v1/organization/:name", ctl.Update)
	rg.POST("/v1/organization", ctl.Create)
	rg.GET("/v1/organization/:name", ctl.Get)
	rg.GET("/v1/organization", ctl.List)
	rg.POST("/v1/organization/:name", ctl.Leave)
	rg.DELETE("/v1/organization/:name", ctl.Delete)
	rg.HEAD("/v1/name", ctl.Check)

	rg.POST("/v1/invite", ctl.InviteMember)
	rg.PUT("/v1/invite", ctl.AcceptInvite)
	rg.GET("/v1/invite", ctl.ListInvitation)
	rg.DELETE("/v1/invite", ctl.RemoveInvitation)

	rg.POST("/v1/request", ctl.RequestMember)
	rg.PUT("/v1/request", ctl.ApproveRequest)
	rg.GET("/v1/request", ctl.ListRequests)
	rg.DELETE("/v1/request", ctl.RemoveRequest)

	rg.DELETE("/v1/organization/:name/member", ctl.RemoveMember)
	rg.GET("/v1/organization/:name/member", ctl.ListMember)
	rg.PUT("/v1/organization/:name/member", ctl.EditMember)

	rg.GET("/v1/account/:name", ctl.GetUser)
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

	//cmd.Actor = pl.DomainAccount()
	cmd.Actor = primitive.CreateAccount(pl.Account)
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
// @Success		200	{object}			commonapp.UserDTO
// @Failure		400	bad_request_param	account	is		invalid
// @Failure		401	resource_not_exists	user	does	not	exist
// @Failure		500	system_error		system	error
// @Router			/v1/account/{name} [get]
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
// @Param			owner		query	string	false	"filter by owner "
// @Param			username	query	string	false	"filter by username"
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

	var req orgListRequest
	if err := ctx.BindQuery(&req); err != nil {
		controller.SendBadRequestParam(ctx, err)

		return
	}

	var os []orgapp.OrganizationDTO
	var user primitive.Account
	var err error
	if req.Owner != "" {
		user, err = primitive.NewAccount(req.Owner)
		if err != nil {
			controller.SendBadRequestParam(ctx, err)

			return
		}
		os, err = ctl.org.GetByOwner(pl.DomainAccount(), user)
	} else if req.Username != "" {
		user, err = primitive.NewAccount(req.Username)
		if err != nil {
			controller.SendBadRequestParam(ctx, err)

			return
		}
		os, err = ctl.org.GetByUser(pl.DomainAccount(), user)
	} else {
		os, err = ctl.org.GetByUser(pl.DomainAccount(), pl.DomainAccount())
	}

	// get org info
	if err != nil {
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

		controller.SendError(ctx, err)

		return
	}

	controller.SendRespOfGet(ctx, members)
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
// @Router			/v1/organization/{name}/member [put]
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

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "edit a member to organization")

	_, err = ctl.org.EditMember(&domain.OrgEditMemberCmd{
		Actor:   pl.DomainAccount(),
		Org:     orgName,
		Role:    req.Role,
		Account: acc,
	})
	if err != nil {
		logrus.Error(err)
		controller.SendError(ctx, err)

		return
	}

	controller.SendRespOfPut(ctx, nil)
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
		controller.SendError(ctx, err)

		return
	}

	controller.SendRespOfDelete(ctx)
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
// @Router			/v1/invite [post]
func (ctl *OrgController) InviteMember(ctx *gin.Context) {

	req := OrgInviteMemberRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	orgName, err := primitive.NewAccount(req.OrgName)
	if err != nil {
		controller.SendBadRequestBody(ctx, err)
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
		Msg:     req.Msg,
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
// @Param			org		query	string	false	"organization name"
// @Param			invitee	query	string	false	"invitee name"
// @Param			inviter	query	string	false	"inviter name"
// @Param			status	query	string	false	"invitation status, can be: pending/approved/rejected"
// @Param			page_size	query int	false	"page size"
// @Param			page		query int 	false	"page index"
// @Accept			json
// @Success		200	{object}	[]orgapp.ApproveDTO
// @Failure		400	bad_request_param	account	is	invalid
// @Failure		401	not_allowed			can't	get	info	of	other	user
// @Router			/v1/invite [get]
func (ctl *OrgController) ListInvitation(ctx *gin.Context) {
	var req OrgListInviteRequest
	if err := ctx.BindQuery(&req); err != nil {
		controller.SendBadRequestParam(ctx, err)

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	cmd := req.toCmd()
	cmd.Actor = pl.DomainAccount()

	dtos, err := ctl.org.ListInvitation(&cmd)
	if err != nil {
		logrus.Error(err)

		controller.SendError(ctx, err)

		return
	}

	controller.SendRespOfGet(ctx, dtos)
}

// @Summary			Revoke member request of the organization
// @Title			Revoke member request of the organization
// @Description Revoke member request of the organization
// @Tags			Organization
// @Param			body	body	OrgRevokeMemberReqRequest	true	"body of the member request"
// @Param			name	path	string	true	"organization name"
// @Accept			json
// @Success		200 {object} []orgapp.ApproveDTO
// @Failure		400	bad_request_param	account	is	invalid
// @Failure		401	not_allowed			can't	get	info	of	other	user
// @Router			/v1//request [delete]
func (ctl *OrgController) RemoveRequest(ctx *gin.Context) {
	req := OrgRevokeMemberReqRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	orgName, err := primitive.NewAccount(req.OrgName)
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

	acc, err := primitive.NewAccount(req.User)
	if err != nil {
		acc = pl.DomainAccount()
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "revoke a request of organization")

	cmd := domain.OrgCancelRequestMemberCmd{
		Actor:     pl.DomainAccount(),
		Org:       orgName,
		Requester: acc,
		Msg:       req.Msg,
	}

	_, err = ctl.org.CancelReqMember(&cmd)
	if err != nil {
		logrus.Error(err)

		controller.SendError(ctx, err)

		return
	}

	controller.SendRespOfDelete(ctx)
}

// @Summary			Request to be a member of the organization
// @Title			Request to be a member of the organization
// @Description Request to be a member of the organization
// @Tags			Organization
// @Param			body	body	OrgReqMemberRequest	true	"body of the member request"
// @Accept			json
// @Success		201 {object} orgapp.OrganizationDTO
// @Failure		400	bad_request_param	account	is	invalid
// @Failure		401	not_allowed			can't	get	info	of	other	user
// @Router			/v1/request [post]
func (ctl *OrgController) RequestMember(ctx *gin.Context) {
	req := OrgReqMemberRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	orgName, err := primitive.NewAccount(req.OrgName)
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

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "request a member to organization")

	cmd := &domain.OrgRequestMemberCmd{}
	cmd.Actor = pl.DomainAccount()
	cmd.Org = orgName
	cmd.Msg = req.Msg
	dto, err := ctl.org.RequestMember(cmd)
	if err != nil {
		logrus.Error(err)

		controller.SendError(ctx, err)

		return
	}

	controller.SendRespOfPost(ctx, dto)
}

// @Summary			Approve a user's member request of the organization
// @Title			Approve a user's member request of the organization
// @Description Approve a user's member request of the organization
// @Tags			Organization
// @Param			body	body	OrgApproveMemberRequest	true	"body of the accept"
// @Accept			json
// @Success		201 {object} orgapp.MemberRequestDTO
// @Failure		400	bad_request_param	account	is	invalid
// @Failure		401	not_allowed			can't	get	info	of	other	user
// @Router			/v1/request [put]
func (ctl *OrgController) ApproveRequest(ctx *gin.Context) {
	req := OrgApproveMemberRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	orgName, err := primitive.NewAccount(req.OrgName)
	if err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, fmt.Errorf("organization name not valid"),
		))

		return
	}

	acc, err := primitive.NewAccount(req.User)
	if err != nil {
		logrus.Error(err)

		controller.SendError(ctx, err)

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "approve a member request to organization")

	_, err = ctl.org.ApproveRequest(&domain.OrgApproveRequestMemberCmd{
		Actor:     pl.DomainAccount(),
		Org:       orgName,
		Msg:       req.Msg,
		Requester: acc,
	})
	if err != nil {
		logrus.Error(err)

		controller.SendError(ctx, err)

		return
	}

	controller.SendRespOfPut(ctx, nil)
}

// @Summary			List requests of the organization
// @Title			List requests of the organization
// @Description List invitation of the organization
// @Tags			Organization
// @Param			org_name	query	string	false	"organization name"
// @Param			requester	query	string	false	"invitee name"
// @Param			status	query	string	false	"invitation status, can be: pending/approved/rejected"
// @Param			page_size	query	int false	"page size"
// @Param			page	query	int	false	"page index"
// @Accept			json
// @Success		200 {object} []orgapp.MemberRequestDTO
// @Failure		400	bad_request_param	account	is	invalid
// @Failure		401	not_allowed			can't	get	info	of	other	user
// @Router			/v1/request [get]
func (ctl *OrgController) ListRequests(ctx *gin.Context) {
	var req OrgListMemberReqRequest
	if err := ctx.BindQuery(&req); err != nil {
		controller.SendBadRequestParam(ctx, err)

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	cmd := req.toCmd()
	cmd.Actor = pl.DomainAccount()

	dtos, err := ctl.org.ListMemberReq(&cmd)
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
// @Router			/v1/invite [delete]
func (ctl *OrgController) RemoveInvitation(ctx *gin.Context) {
	req := OrgRevokeInviteRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	orgName, err := primitive.NewAccount(req.OrgName)
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

	acc, err := primitive.NewAccount(req.User)
	if err != nil {
		acc = pl.DomainAccount()
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "revoke a invitation of organization")

	_, err = ctl.org.RevokeInvite(&domain.OrgRemoveInviteCmd{
		Actor:   pl.DomainAccount(),
		Org:     orgName,
		Account: acc,
		Msg:     req.Msg,
	})
	if err != nil {
		logrus.Error(err)

		controller.SendError(ctx, err)

		return
	}

	controller.SendRespOfDelete(ctx)
}

// @Summary			Accept invite of the organization
// @Title			Accept invite of the organization
// @Description Accept invite of the organization
// @Tags			Organization
// @Param			body	body	OrgAcceptMemberRequest	true	"body of the invitation"
// @Accept			json
// @Success		200 {object} []orgapp.ApproveDTO
// @Failure		400	bad_request_param	account	is	invalid
// @Failure		401	not_allowed			can't	get	info	of	other	user
// @Router			/v1/invite [put]
func (ctl *OrgController) AcceptInvite(ctx *gin.Context) {
	req := OrgAcceptMemberRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	orgName, err := primitive.NewAccount(req.OrgName)
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

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "accept a invitation of organization")

	_, err = ctl.org.AcceptInvite(&domain.OrgAcceptInviteCmd{
		Actor:   pl.DomainAccount(),
		Org:     orgName,
		Msg:     req.Msg,
		Account: pl.DomainAccount(),
	})
	if err != nil {
		logrus.Error(err)

		controller.SendError(ctx, err)

		return
	}

	controller.SendRespOfPut(ctx, nil)
}
