package controller

import (
	"fmt"

	"github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/organization/domain"
)

type orgBasicInfoUpdateRequest struct {
	FullName     string `json:"fullname"`
	Website      string `json:"website"`
	AvatarId     string `json:"avatar_id"`
	Description  string `json:"description"`
	AllowRequest *bool  `json:"allow_request,omitempty"`
	DefaultRole  string `json:"default_role"`
}

func (req *orgBasicInfoUpdateRequest) toCmd(user primitive.Account, orgName string) (
	cmd domain.OrgUpdatedBasicInfoCmd,
	err error,
) {
	cmd.Actor = user

	if cmd.OrgName, err = primitive.NewAccount(orgName); err != nil {
		return
	}

	empty := true
	if req.FullName != "" {
		cmd.FullName = req.FullName
		empty = false
	}

	if req.AvatarId != "" {
		cmd.AvatarId = req.AvatarId
		empty = false
	}

	if req.Website != "" {
		cmd.Website = req.Website
		empty = false
	}

	if req.Description != "" {
		cmd.Description = req.Description
		empty = false
	}

	if req.DefaultRole != "" {
		cmd.DefaultRole = req.DefaultRole
		empty = false
	}

	if req.AllowRequest != nil {
		cmd.AllowRequest = req.AllowRequest
		empty = false
	}

	if empty {
		err = fmt.Errorf("edit org param can't be all empty")
	}

	return
}

type orgListRequest struct {
	Owner    string `form:"owner"`
	Username string `form:"username"`
}

func (req *orgListRequest) toCmd() (owner, user primitive.Account, err error) {
	if req.Owner != "" {
		if owner, err = primitive.NewAccount(req.Owner); err != nil {
			return
		}
	}

	if req.Username != "" {
		if user, err = primitive.NewAccount(req.Username); err != nil {
			return
		}
	}

	return
}

type orgCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Website     string `json:"website"`
	FullName    string `json:"fullname"`
	AvatarId    string `json:"avatar_id"`
	Description string `json:"description"`
}

func (req *orgCreateRequest) toCmd() (cmd domain.OrgCreatedCmd, err error) {
	if cmd.Name, err = primitive.NewAccount(req.Name); err != nil {
		return
	}

	if req.FullName == "" {
		err = allerror.NewInvalidParam("org fullname can't be empty")
		return
	}

	if cmd.FullName, err = primitive.NewMSDFullname(req.FullName); err != nil {
		return
	}

	if cmd.AvatarId, err = primitive.NewAvatarId(req.AvatarId); err != nil {
		return
	}

	if cmd.Description, err = primitive.NewMSDDesc(req.Description); err != nil {
		return
	}

	cmd.Website = req.Website

	return
}

type orgMemberAddRequest struct {
	User string `json:"user" binding:"required"`
}

type orgMemberRemoveRequest struct {
	User string `json:"user" binding:"required"`
}

func (req *orgMemberRemoveRequest) toCmd(orgName string, user primitive.Account) (
	cmd domain.OrgRemoveMemberCmd, err error,
) {
	if cmd.Org, err = primitive.NewAccount(orgName); err != nil {
		return
	}

	if cmd.Account, err = primitive.NewAccount(req.User); err != nil {
		return
	}

	cmd.Actor = user

	return
}

type OrgListInviteRequest struct {
	controller.CommonListRequest
	Inviter string `form:"inviter"`
	Invitee string `form:"invitee"`
	OrgName string `form:"org_name"`
	Status  string `form:"status"`
}

func (req *OrgListInviteRequest) toCmd(user primitive.Account) (cmd domain.OrgInvitationListCmd) {
	cmd.Actor = user

	if inviter, err := primitive.NewAccount(req.Inviter); err == nil {
		cmd.Inviter = inviter
	}

	if invitee, err := primitive.NewAccount(req.Invitee); err == nil {
		cmd.Invitee = invitee
	}

	if org, err := primitive.NewAccount(req.OrgName); err == nil {
		cmd.Org = org
	}

	if req.Status != "" {
		cmd.Status = domain.ApproveStatus(req.Status)
	}

	return

}

type OrgListMemberReqRequest struct {
	controller.CommonListRequest
	Requester string `form:"requester"`
	OrgName   string `form:"org_name"`
	Status    string `form:"status"`
}

func (req *OrgListMemberReqRequest) toCmd(user primitive.Account) (cmd domain.OrgMemberReqListCmd, err error) {
	cmd.Actor = user

	if req.Requester != "" {
		if cmd.Requester, err = primitive.NewAccount(req.Requester); err != nil {
			return
		}
	}

	if req.OrgName != "" {
		if cmd.Org, err = primitive.NewAccount(req.OrgName); err != nil {
			return
		}
	}

	cmd.Status = domain.ApproveStatus(req.Status)

	return

}

type OrgMemberEditRequest struct {
	Role string `json:"role" binding:"required"`
	User string `json:"user" binding:"required"`
}

func (req *OrgMemberEditRequest) toCmd(orgName string, user primitive.Account) (
	cmd domain.OrgEditMemberCmd, err error,
) {
	if cmd.Org, err = primitive.NewAccount(orgName); err != nil {
		return
	}

	if cmd.Account, err = primitive.NewAccount(req.User); err != nil {
		return
	}

	cmd.Actor = user
	cmd.Role = req.Role

	return
}

type OrgInviteMemberRequest struct {
	Role    string `json:"role" binding:"required"`
	User    string `json:"user" binding:"required"`
	Msg     string `json:"msg"`
	OrgName string `json:"org_name" binding:"required"`
}

func (req *OrgInviteMemberRequest) toCmd(user primitive.Account) (
	cmd domain.OrgInviteMemberCmd, err error,
) {
	if cmd.Org, err = primitive.NewAccount(req.OrgName); err != nil {
		return
	}

	if cmd.Account, err = primitive.NewAccount(req.User); err != nil {
		return
	}

	cmd.Actor = user
	cmd.Role = req.Role
	cmd.Msg = req.Msg

	return
}

type OrgAcceptMemberRequest struct {
	Msg     string `json:"msg"`
	OrgName string `json:"org_name" binding:"required"`
}

func (req *OrgAcceptMemberRequest) toCmd(user primitive.Account) (
	cmd domain.OrgAcceptInviteCmd, err error,
) {
	if cmd.Org, err = primitive.NewAccount(req.OrgName); err != nil {
		return
	}

	cmd.Account = user
	cmd.Actor = user
	cmd.Msg = req.Msg

	return
}

type OrgApproveMemberRequest struct {
	User    string `json:"user"`
	Msg     string `json:"msg"`
	OrgName string `json:"org_name" binding:"required"`
}

func (req *OrgApproveMemberRequest) toCmd(user primitive.Account) (
	cmd domain.OrgApproveRequestMemberCmd, err error,
) {
	if cmd.Org, err = primitive.NewAccount(req.OrgName); err != nil {
		return
	}

	if cmd.Requester, err = primitive.NewAccount(req.User); err != nil {
		return
	}

	cmd.Actor = user
	cmd.Msg = req.Msg

	return
}

type OrgReqMemberRequest struct {
	Msg     string `json:"msg"`
	OrgName string `json:"org_name" binding:"required"`
}

func (req *OrgReqMemberRequest) toCmd(user primitive.Account) (
	cmd domain.OrgRequestMemberCmd, err error,
) {
	if cmd.Org, err = primitive.NewAccount(req.OrgName); err != nil {
		return
	}

	cmd.Actor = user
	cmd.Msg = req.Msg

	return
}

type OrgRevokeInviteRequest struct {
	User    string `json:"user"`
	Msg     string `json:"msg"`
	OrgName string `json:"org_name" binding:"required"`
}

func (req *OrgRevokeInviteRequest) toCmd(user primitive.Account) (
	cmd domain.OrgRemoveInviteCmd, err error,
) {
	if cmd.Org, err = primitive.NewAccount(req.OrgName); err != nil {
		return
	}

	if req.User == "" {
		cmd.Account = user
	} else {
		if cmd.Account, err = primitive.NewAccount(req.User); err != nil {
			return
		}
	}

	cmd.Actor = user
	cmd.Msg = req.Msg

	return
}

type OrgRevokeMemberReqRequest struct {
	User    string `json:"user"`
	Msg     string `json:"msg"`
	OrgName string `json:"org_name" binding:"required"`
}

func (req *OrgRevokeMemberReqRequest) toCmd(user primitive.Account) (
	cmd domain.OrgCancelRequestMemberCmd, err error,
) {
	if cmd.Org, err = primitive.NewAccount(req.OrgName); err != nil {
		return
	}

	if req.User == "" {
		cmd.Requester = user
	} else {
		if cmd.Requester, err = primitive.NewAccount(req.User); err != nil {
			return
		}
	}

	cmd.Actor = user
	cmd.Msg = req.Msg

	return
}

// reqToCheckName
type reqToCheckName struct {
	Name string `form:"name"`
}

func (req *reqToCheckName) toAccount() (primitive.Account, error) {
	return primitive.NewAccount(req.Name)
}
