/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides the controllers for handling HTTP requests and managing the application's business logic.
package controller

import (
	"errors"
	"fmt"

	"github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/organization/domain"
)

type orgBasicInfoUpdateRequest struct {
	FullName     string  `json:"fullname"`
	Website      *string `json:"website"`
	AvatarId     string  `json:"avatar_id"`
	Description  string  `json:"description"`
	AllowRequest *bool   `json:"allow_request,omitempty"`
	DefaultRole  string  `json:"default_role"`
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
		if cmd.FullName, err = primitive.NewOrgFullname(req.FullName); err != nil {
			return
		}
		empty = false
	}

	if req.Website != nil {
		if cmd.Website, err = primitive.NewOrgWebsite(*req.Website); err != nil {
			return
		}
		empty = false
	}

	if req.Description != "" {
		if cmd.Description, err = primitive.NewAccountDesc(req.Description); err != nil {
			return
		}
		empty = false
	}

	if req.DefaultRole != "" {
		if cmd.DefaultRole, err = primitive.NewRole(req.DefaultRole); err != nil {
			return
		}

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
	Owner    string   `form:"owner"`
	Username string   `form:"username"`
	Roles    []string `form:"roles"`
}

func (req *orgListRequest) toCmd() (owner, user primitive.Account, roles []primitive.Role, err error) {
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

	if len(req.Roles) > 0 {
		roles = make([]primitive.Role, 0, len(req.Roles))
		var r primitive.Role

		for _, val := range req.Roles {
			if r, err = primitive.NewRole(val); err != nil {
				return
			} else {
				roles = append(roles, r)
			}
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

func (req *orgCreateRequest) action() string {
	return fmt.Sprintf("create organization %s", req.Name)
}

func (req *orgCreateRequest) toCmd() (cmd domain.OrgCreatedCmd, err error) {
	if cmd.Name, err = primitive.NewAccount(req.Name); err != nil {
		return
	}

	if req.FullName == "" {
		e := fmt.Errorf("fullname can't be empty")
		err = allerror.New(allerror.ErrorFullnameCanNotBeEmpty, "", e)
		return
	}

	if cmd.FullName, err = primitive.NewOrgFullname(req.FullName); err != nil {
		return
	}

	if cmd.AvatarId, err = primitive.NewAvatarId(req.AvatarId); err != nil {
		return
	}

	if cmd.Description, err = primitive.NewAccountDesc(req.Description); err != nil {
		return
	}

	if cmd.Website, err = primitive.NewOrgWebsite(req.Website); err != nil {
		return
	}
	return
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

// OrgListInviteRequest is a struct for handling organization invite requests.
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

// OrgListMemberReqRequest is a struct for handling organization member request list requests.
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

// OrgMemberEditRequest is a struct for handling organization member editing requests.
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

	if cmd.Role, err = primitive.NewRole(req.Role); err != nil {
		return
	}

	cmd.Actor = user

	return
}

// OrgInviteMemberRequest is a struct for handling organization member invite requests.
type OrgInviteMemberRequest struct {
	Role    string `json:"role" binding:"required"`
	User    string `json:"user" binding:"required"`
	Msg     string `json:"msg"`
	OrgName string `json:"org_name" binding:"required"`
}

func (req *OrgInviteMemberRequest) action() string {
	return fmt.Sprintf("invite %s as a %s in %s", req.User, req.Role, req.OrgName)
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

	if cmd.Role, err = primitive.NewRole(req.Role); err != nil {
		return
	}

	cmd.Actor = user
	cmd.Msg = req.Msg

	return
}

// OrgAcceptMemberRequest is a struct for handling organization member acceptance requests.
type OrgAcceptMemberRequest struct {
	Msg     string `json:"msg"`
	OrgName string `json:"org_name" binding:"required"`
}

func (req *OrgAcceptMemberRequest) action() string {
	return fmt.Sprintf("accept invite from organization %s", req.OrgName)
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

// OrgApproveMemberRequest is a struct for handling organization member approval requests.
type OrgApproveMemberRequest struct {
	User    string `json:"user"`
	Msg     string `json:"msg"`
	OrgName string `json:"org_name" binding:"required"`
}

func (req *OrgApproveMemberRequest) action() string {
	return fmt.Sprintf("approve %s to be a member of %s", req.User, req.OrgName)
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

// OrgReqMemberRequest is a struct for handling organization member request creation requests.
type OrgReqMemberRequest struct {
	Msg     string `json:"msg"`
	OrgName string `json:"org_name" binding:"required"`
}

func (req *OrgReqMemberRequest) action() string {
	return fmt.Sprintf("request to be a member of %s", req.OrgName)
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

// OrgRevokeInviteRequest is a struct for handling organization invite revocation requests.
type OrgRevokeInviteRequest struct {
	User    string `json:"user"`
	Msg     string `json:"msg"`
	OrgName string `json:"org_name" binding:"required"`
}

func (req *OrgRevokeInviteRequest) action() string {
	return fmt.Sprintf("revoke invite of %s from %s", req.User, req.OrgName)
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

// OrgRevokeMemberReqRequest is a struct for handling organization member request revocation requests.
type OrgRevokeMemberReqRequest struct {
	User    string `json:"user"`
	Msg     string `json:"msg"`
	OrgName string `json:"org_name" binding:"required"`
}

func (req *OrgRevokeMemberReqRequest) action() string {
	return fmt.Sprintf("revoke member request of %s to join %s", req.User, req.OrgName)
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

type orgListMemberRequest struct {
	Username string `form:"username"`
	Role     string `form:"role"`
}

func (req *orgListMemberRequest) toCmd(org primitive.Account) (cmd domain.OrgListMemberCmd, err error) {
	if org == nil {
		err = errors.New("organization name can not be empty")
		return
	}

	if req.Username != "" {
		if cmd.User, err = primitive.NewAccount(req.Username); err != nil {
			return
		}
	}

	if req.Role != "" {
		if cmd.Role, err = primitive.NewRole(req.Role); err != nil {
			return
		}
	}

	cmd.Org = org

	return
}
