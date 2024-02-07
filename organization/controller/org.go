package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	orgapp "github.com/openmerlin/merlin-server/organization/app"
	"github.com/openmerlin/merlin-server/organization/domain"
	userapp "github.com/openmerlin/merlin-server/user/app"
	userctl "github.com/openmerlin/merlin-server/user/controller"
)

func AddRouterForOrgController(
	rg *gin.RouterGroup,
	org orgapp.OrgService,
	user userapp.UserService,
	l middleware.OperationLog,
	m middleware.UserMiddleWare,
) {
	ctl := OrgController{
		m:    m,
		org:  org,
		user: user,
	}

	rg.PUT("/v1/organization/:name", m.Write, userctl.CheckMail(ctl.m, ctl.user), l.Write, ctl.Update)
	rg.POST("/v1/organization", m.Write, userctl.CheckMail(ctl.m, ctl.user), l.Write, ctl.Create)
	rg.GET("/v1/organization/:name", m.Optional, ctl.Get)
	rg.GET("/v1/organization", m.Read, ctl.List)
	rg.POST("/v1/organization/:name", m.Write, userctl.CheckMail(ctl.m, ctl.user), l.Write, ctl.Leave)
	rg.DELETE("/v1/organization/:name", m.Write, userctl.CheckMail(ctl.m, ctl.user), l.Write, ctl.Delete)
	rg.HEAD("/v1/name", m.Read, ctl.Check)

	rg.POST("/v1/invite", m.Write, userctl.CheckMail(ctl.m, ctl.user), l.Write, ctl.InviteMember)
	rg.PUT("/v1/invite", m.Write, userctl.CheckMail(ctl.m, ctl.user), l.Write, ctl.AcceptInvite)
	rg.GET("/v1/invite", m.Read, ctl.ListInvitation)
	rg.DELETE("/v1/invite", m.Write, userctl.CheckMail(ctl.m, ctl.user), l.Write, ctl.RemoveInvitation)

	rg.POST("/v1/request", m.Write, userctl.CheckMail(ctl.m, ctl.user), l.Write, ctl.RequestMember)
	rg.PUT("/v1/request", m.Write, userctl.CheckMail(ctl.m, ctl.user), l.Write, ctl.ApproveRequest)
	rg.GET("/v1/request", m.Read, ctl.ListRequests)
	rg.DELETE("/v1/request", m.Write, userctl.CheckMail(ctl.m, ctl.user), l.Write, ctl.RemoveRequest)

	rg.DELETE("/v1/organization/:name/member", m.Write, userctl.CheckMail(ctl.m, ctl.user), l.Write, ctl.RemoveMember)
	rg.GET("/v1/organization/:name/member", m.Read, ctl.ListMember)
	rg.PUT("/v1/organization/:name/member", m.Write, userctl.CheckMail(ctl.m, ctl.user), l.Write, ctl.EditMember)

	rg.GET("/v1/account/:name", m.Optional, ctl.GetUser)
}

type OrgController struct {
	m    middleware.UserMiddleWare
	org  orgapp.OrgService
	user userapp.UserService
}

// @Summary  Update
// @Description  update org basic info
// @Tags     Organization
// @Param    name  path  string                     true  "name"
// @Param    body  body  orgBasicInfoUpdateRequest  true  "body of new organization"
// @Accept   json
// @Security Bearer
// @Success  202  {object}  commonctl.ResponseData
// @Router   /v1/organization/{name} [put]
func (ctl *OrgController) Update(ctx *gin.Context) {
	var req orgBasicInfoUpdateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	middleware.SetAction(ctx, fmt.Sprintf("update basic info of %s", ctx.Param("name")))

	cmd, err := req.toCmd(user, ctx.Param("name"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if o, err := ctl.org.UpdateBasicInfo(&cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, o)
	}
}

// @Summary  Get
// @Description  get organization info
// @Tags     Organization
// @Param    name  path  string  true  "name"
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData
// @Router   /v1/organization/{name} [get]
func (ctl *OrgController) Get(ctx *gin.Context) {
	orgName, err := primitive.NewAccount(ctx.Param("name"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if o, err := ctl.org.GetByAccount(orgName); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, o)
	}
}

// @Summary   User or organization info
// @Description  get organization or user info
// @Tags     Organization
// @Param    name  path  string  true  "name of the user of organization"
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData
// @Failure  404  "user not found"
// @Failure  400  {object}  commonctl.ResponseData
// @Router   /v1/account/{name} [get]
func (ctl *OrgController) GetUser(ctx *gin.Context) {
	name, err := primitive.NewAccount(ctx.Param("name"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if o, err := ctl.org.GetOrgOrUser(name); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, o)
	}
}

// @Summary  Check
// @Description  Check the name is available
// @Tags     Name
// @Param    name  query  string  true  "the name to be check whether it's usable"
// @Accept   json
// @Security Bearer
// @Success  200  "name is valid"
// @Failure  409  "name is invalid"
// @Router   /v1/name [head]
func (ctl *OrgController) Check(ctx *gin.Context) {
	// TODO why head method

	var req reqToCheckName

	if err := ctx.BindQuery(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	name, err := req.toAccount()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if ctl.org.CheckName(name) {
		ctx.JSON(http.StatusOK, nil)
	} else {
		ctx.JSON(http.StatusConflict, nil)
	}
}

// @Summary  List
// @Description  get organization info
// @Tags     Organization
// @Param    owner     query  string  false  "filter by owner"
// @Param    username  query  string  false  "filter by username"
// @Accept   json
// @Security Bearer
// @Success  200  {object}  commonctl.ResponseData
// @Failure  400  {object}  commonctl.ResponseData
// @Router   /v1/organization [get]
func (ctl *OrgController) List(ctx *gin.Context) {
	var req orgListRequest

	if err := ctx.BindQuery(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	me := ctl.m.GetUserAndExitIfFailed(ctx)
	if me == nil {
		return
	}

	owner, user, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if owner != nil {
		if os, err := ctl.org.GetByOwner(me, owner); err != nil {
			commonctl.SendError(ctx, err)
		} else {
			commonctl.SendRespOfGet(ctx, os)
		}

		return
	}

	if user == nil {
		user = me
	}

	if os, err := ctl.org.GetByUser(me, user); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, os)
	}
}

// @Summary  Create
// @Description  create a new organization
// @Tags     Organization
// @Param    body  body  orgCreateRequest  true  "body of new organization"
// @Accept   json
// @Security Bearer
// @Success  201 {object}  commonctl.ResponseData
// @Router   /v1/organization [post]
func (ctl *OrgController) Create(ctx *gin.Context) {
	var req orgCreateRequest

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	middleware.SetAction(ctx, req.action())

	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	cmd.Owner = user

	o, err := ctl.org.Create(&cmd)
	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, o)
	}
}

// @Summary  Delete
// @Description  delete a organization
// @Tags     Organization
// @Param    name  path  string  true  "name"
// @Accept   json
// @Security Bearer
// @Success  204
// @Router   /v1/organization/{name} [delete]
func (ctl *OrgController) Delete(ctx *gin.Context) {
	orgName, err := primitive.NewAccount(ctx.Param("name"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	middleware.SetAction(ctx, fmt.Sprintf("delete organization %s", orgName))

	err = ctl.org.Delete(&domain.OrgDeletedCmd{
		Actor: user,
		Name:  orgName,
	})
	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfDelete(ctx)
	}
}

// @Summary  ListMember
// @Description  list organization members
// @Tags     Organization
// @Param    name  path  string  true  "name"
// @Accept   json
// @Security Bearer
// @Success  200 {object}  commonctl.ResponseData
// @Router   /v1/organization/{name}/member [get]
func (ctl *OrgController) ListMember(ctx *gin.Context) {
	orgName, err := primitive.NewAccount(ctx.Param("name"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if members, err := ctl.org.ListMember(orgName); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, members)
	}
}

// @Summary  EditMember
// @Description Edit a member to the organization's role
// @Tags     Organization
// @Param    body  body  OrgMemberEditRequest  true  "body of new member"
// @Param    name  path  string  true  "name"
// @Accept   json
// @Security Bearer
// @Success  202 {object}  commonctl.ResponseData
// @Router   /v1/organization/{name}/member [put]
func (ctl *OrgController) EditMember(ctx *gin.Context) {
	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	var req OrgMemberEditRequest

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd, err := req.toCmd(ctx.Param("name"), user)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	middleware.SetAction(ctx,
		fmt.Sprintf("edit member %s to be %s of %s", req.User, req.Role, cmd.Org.Account()))

	if _, err = ctl.org.EditMember(&cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}

// @Summary  RemoveMember
// @Description Remove a member from a organization
// @Tags     Organization
// @Param    body  body  orgMemberRemoveRequest  true  "body of the removed member"
// @Param    name  path  string  true  "name"
// @Accept   json
// @Security Bearer
// @Success  204
// @Router   /v1/organization/{name}/member [delete]
func (ctl *OrgController) RemoveMember(ctx *gin.Context) {
	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	var req orgMemberRemoveRequest

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd, err := req.toCmd(ctx.Param("name"), user)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	middleware.SetAction(ctx,
		fmt.Sprintf("remove member %s from %s", req.User, cmd.Org.Account()))

	if err = ctl.org.RemoveMember(&cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfDelete(ctx)
	}
}

// @Summary  Leave
// @Description  leave the organization
// @Tags     Organization
// @Param    name  path  string  true  "name"
// @Accept   json
// @Security Bearer
// @Success  204
// @Router   /v1/organization/{name} [post]
func (ctl *OrgController) Leave(ctx *gin.Context) {
	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	orgName, err := primitive.NewAccount(ctx.Param("name"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	middleware.SetAction(ctx, fmt.Sprintf("leave organization %s", orgName))

	err = ctl.org.RemoveMember(&domain.OrgRemoveMemberCmd{
		Actor:   user,
		Org:     orgName,
		Account: user,
	})
	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfDelete(ctx)
	}
}

// @Summary  InviteMember
// @Description Send invitation to a user to join the organization
// @Tags     Organization
// @Param    body  body  OrgInviteMemberRequest  true  "body of the invitation"
// @Accept   json
// @Security Bearer
// @Success  201 {object}  commonctl.ResponseData
// @Router   /v1/invite [post]
func (ctl *OrgController) InviteMember(ctx *gin.Context) {
	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	var req OrgInviteMemberRequest

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	middleware.SetAction(ctx, req.action())

	cmd, err := req.toCmd(user)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if dto, err := ctl.org.InviteMember(&cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, dto)
	}
}

// @Summary  ListInvitation
// @Description List invitation of the organization
// @Tags     Organization
// @Param    org_name   query  string  false  "organization name"
// @Param    invitee    query  string  false  "invitee name"
// @Param    inviter    query  string  false  "inviter name"
// @Param    status     query  string  false  "invitation status, can be: pending/approved/rejected"
// @Param    page_size  query  int     false  "page size"
// @Param    page       query  int     false  "page index"
// @Accept   json
// @Security Bearer
// @Success  200  {object}  commonctl.ResponseData
// @Router   /v1/invite [get]
func (ctl *OrgController) ListInvitation(ctx *gin.Context) {
	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	var req OrgListInviteRequest

	if err := ctx.BindQuery(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	cmd := req.toCmd(user)

	if dtos, err := ctl.org.ListInvitation(&cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, dtos)
	}
}

// @Summary  RevokeMember
// @Description Revoke member request of the organization
// @Tags     Organization
// @Param    body  body  OrgRevokeMemberReqRequest  true  "body of the member request"
// @Accept   json
// @Security Bearer
// @Success  200 {object}  commonctl.ResponseData
// @Router   /v1/request [delete]
func (ctl *OrgController) RemoveRequest(ctx *gin.Context) {
	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	var req OrgRevokeMemberReqRequest

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	middleware.SetAction(ctx, req.action())

	cmd, err := req.toCmd(user)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if _, err = ctl.org.CancelReqMember(&cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfDelete(ctx)
	}
}

// @Summary  RequestMember
// @Description Request to be a member of the organization
// @Tags     Organization
// @Param    body  body  OrgReqMemberRequest  true  "body of the member request"
// @Accept   json
// @Security Bearer
// @Success  201 {object}  commonctl.ResponseData
// @Router   /v1/request [post]
func (ctl *OrgController) RequestMember(ctx *gin.Context) {
	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	var req OrgReqMemberRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	middleware.SetAction(ctx, req.action())

	cmd, err := req.toCmd(user)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if dto, err := ctl.org.RequestMember(&cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, dto)
	}
}

// @Summary  ApproveRequest
// @Description Approve a user's member request of the organization
// @Tags     Organization
// @Param    body  body  OrgApproveMemberRequest  true  "body of the accept"
// @Accept   json
// @Security Bearer
// @Success  201 {object}  commonctl.ResponseData
// @Router   /v1/request [put]
func (ctl *OrgController) ApproveRequest(ctx *gin.Context) {
	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	var req OrgApproveMemberRequest

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	middleware.SetAction(ctx, req.action())

	cmd, err := req.toCmd(user)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if _, err = ctl.org.ApproveRequest(&cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}

// @Summary  ListRequests
// @Description  List requests of the organization
// @Tags     Organization
// @Param    org_name   query  string  false  "organization name"
// @Param    requester  query  string  false  "invitee name"
// @Param    status     query  string  false  "invitation status, can be: pending/approved/rejected"
// @Param    page_size  query  int     false  "page size"
// @Param    page       query  int     false  "page index"
// @Accept   json
// @Security Bearer
// @Success  200 {object}  commonctl.ResponseData
// @Router   /v1/request [get]
func (ctl *OrgController) ListRequests(ctx *gin.Context) {
	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	var req OrgListMemberReqRequest

	if err := ctx.BindQuery(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	// TODO  pagination is not working
	cmd, err := req.toCmd(user)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if dtos, err := ctl.org.ListMemberReq(&cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, dtos)
	}
}

// @Summary  RevokeInvitation
// @Description Revoke invitation of the organization
// @Tags     Organization
// @Param    body  body  OrgRevokeInviteRequest  true  "body of the invitation"
// @Accept   json
// @Security Bearer
// @Success  204
// @Router   /v1/invite [delete]
func (ctl *OrgController) RemoveInvitation(ctx *gin.Context) {
	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	var req OrgRevokeInviteRequest

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	middleware.SetAction(ctx, req.action())

	cmd, err := req.toCmd(user)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if _, err = ctl.org.RevokeInvite(&cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfDelete(ctx)
	}
}

// @Summary  AcceptInvite
// @Description Accept invite of the organization
// @Tags     Organization
// @Param    body  body  OrgAcceptMemberRequest  true  "body of the invitation"
// @Accept   json
// @Security Bearer
// @Success  202 {object}  commonctl.ResponseData
// @Router   /v1/invite [put]
func (ctl *OrgController) AcceptInvite(ctx *gin.Context) {
	user := ctl.m.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	var req OrgAcceptMemberRequest

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	middleware.SetAction(ctx, req.action())

	cmd, err := req.toCmd(user)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if a, err := ctl.org.AcceptInvite(&cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, a)
	}
}
