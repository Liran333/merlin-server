package controller

import (
	"fmt"

	"github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/organization/domain"
)

type orgBasicInfoUpdateRequest struct {
	FullName     string `json:"full_name"`
	Website      string `json:"website"`
	AvatarId     string `json:"avatar_id"`
	Description  string `json:"description"`
	AllowRequest *bool  `json:"allow_request,omitempty"`
	DefaultRole  string `json:"default_role"`
}

func (req *orgBasicInfoUpdateRequest) toCmd() (
	cmd domain.OrgUpdatedBasicInfoCmd,
	err error,
) {
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
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Owner    string `form:"owner"`
	Username string `form:"username"`
}

type orgCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Website     string `json:"website"`
	FullName    string `json:"full_name"`
	AvatarId    string `json:"avatar_id"`
	Description string `json:"description"`
}

func (req *orgCreateRequest) toCmd() (cmd domain.OrgCreatedCmd, err error) {
	if cmd.Name == "" {
		err = fmt.Errorf("org name can't be empty")
		return
	}

	if cmd.FullName == "" {
		err = fmt.Errorf("org full name can't be empty")
		return
	}

	err = cmd.Validate()

	return
}

type orgMemberAddRequest struct {
	User string `json:"user" binding:"required"`
}

type orgMemberRemoveRequest struct {
	User string `json:"user" binding:"required"`
}

type OrgListInviteRequest struct {
	controller.CommonListRequest
	Inviter string `form:"inviter"`
	Invitee string `form:"invitee"`
	OrgName string `form:"org_name"`
	Status  string `form:"status"`
}

type OrgListMemberReqRequest struct {
	controller.CommonListRequest
	Requester string `form:"requester"`
	OrgName   string `form:"org_name"`
	Status    string `form:"status"`
}

func (req *OrgListInviteRequest) toCmd() (cmd domain.OrgInvitationListCmd) {
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

func (req *OrgListMemberReqRequest) toCmd() (cmd domain.OrgMemberReqListCmd) {
	if requester, err := primitive.NewAccount(req.Requester); err == nil {
		cmd.Requester = requester
	}

	if org, err := primitive.NewAccount(req.OrgName); err == nil {
		cmd.Org = org
	}

	cmd.Status = domain.ApproveStatus(req.Status)

	return

}

type OrgMemberEditRequest struct {
	Role string `json:"role" binding:"required"`
	User string `json:"user" binding:"required"`
}

type OrgInviteMemberRequest struct {
	Role    string `json:"role" binding:"required"`
	User    string `json:"user" binding:"required"`
	Msg     string `json:"msg"`
	OrgName string `json:"org_name" binding:"required"`
}

type OrgAcceptMemberRequest struct {
	Msg     string `json:"msg"`
	OrgName string `json:"org_name" binding:"required"`
}

type OrgApproveMemberRequest struct {
	User    string `json:"user"`
	Msg     string `json:"msg"`
	OrgName string `json:"org_name" binding:"required"`
}

type OrgReqMemberRequest struct {
	Msg     string `json:"msg"`
	OrgName string `json:"org_name" binding:"required"`
}

type OrgRevokeInviteRequest struct {
	User    string `json:"user"`
	Msg     string `json:"msg"`
	OrgName string `json:"org_name" binding:"required"`
}

type OrgRevokeMemberReqRequest struct {
	User    string `json:"user"`
	Msg     string `json:"msg"`
	OrgName string `json:"org_name" binding:"required"`
}
