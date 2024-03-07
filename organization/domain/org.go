/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package domain

import (
	"fmt"

	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/user/domain"
	user "github.com/openmerlin/merlin-server/user/domain"
	"github.com/openmerlin/merlin-server/utils"
)

// InviteType represents the type of invitation.
type InviteType = string

// Organization represents a user's organization.
type Organization = user.User

const (
	// InviteTypeInvite represents the invite type for invitations.
	InviteTypeInvite InviteType = "invite"

	// InviteTypeRequest represents the invite type for requests.
	InviteTypeRequest InviteType = "request"
)

// OrgCreatedCmd represents the command for creating an organization.
type OrgCreatedCmd struct {
	Name        primitive.Account     `json:"name"`
	FullName    primitive.MSDFullname `json:"fullname"`
	Description primitive.MSDDesc     `json:"description"`
	Website     string                `json:"website"`
	AvatarId    primitive.AvatarId    `json:"avatar_id"`
	Owner       primitive.Account     `json:"owner"`
}

// OrgDeletedCmd represents the command for deleting an organization.
type OrgDeletedCmd struct {
	Actor primitive.Account
	Name  primitive.Account
}

// Validate validates the OrgDeletedCmd fields.
func (cmd OrgDeletedCmd) Validate() error {
	if cmd.Name == nil {
		return allerror.NewInvalidParam("invalid org name")
	}

	if cmd.Actor == nil {
		return allerror.NewInvalidParam("invalid actor name")
	}

	return nil
}

// OrgUpdatedBasicInfoCmd represents the command for updating basic information of an organization.
type OrgUpdatedBasicInfoCmd struct {
	Actor        primitive.Account
	OrgName      primitive.Account
	AllowRequest *bool
	DefaultRole  primitive.Role
	FullName     string
	Description  string
	Website      string
	AvatarId     string
}

// Validate validates the OrgUpdatedBasicInfoCmd fields.
func (cmd OrgUpdatedBasicInfoCmd) Validate() error {
	if cmd.Website != "" && !utils.IsUrl(cmd.Website) {
		return allerror.NewInvalidParam("invalid website")
	}

	if cmd.Actor == nil {
		return allerror.NewInvalidParam("account is nil")
	}

	if cmd.OrgName == nil {
		return allerror.NewInvalidParam("org name is nil")
	}

	return nil
}

// ToOrg updates the Organization object with the values from the OrgUpdatedBasicInfoCmd.
func (cmd OrgUpdatedBasicInfoCmd) ToOrg(o *Organization) (change bool, err error) {
	if cmd.AvatarId != o.AvatarId.AvatarId() {
		o.AvatarId, err = primitive.NewAvatarId(cmd.AvatarId)
		if err != nil {
			err = allerror.NewInvalidParam(fmt.Sprintf("failed to create avatar id, %s", err))
			return
		}
		change = true
	}

	if cmd.Website != o.Website && cmd.Website != "" {
		o.Website = cmd.Website
		change = true
	}

	if cmd.Description != o.Desc.MSDDesc() && cmd.Description != "" {
		o.Desc = primitive.CreateMSDDesc(cmd.Description)
		change = true
	}

	if cmd.FullName != o.Fullname.MSDFullname() && cmd.FullName != "" {
		o.Fullname = primitive.CreateMSDFullname(cmd.FullName)
		change = true
	}

	if cmd.AllowRequest != nil && *cmd.AllowRequest != o.AllowRequest {
		o.AllowRequest = *cmd.AllowRequest
		change = true
	}

	if cmd.DefaultRole != nil && cmd.DefaultRole != o.DefaultRole {
		o.DefaultRole = primitive.Role(cmd.DefaultRole)
		change = true
	}

	return
}

// ToOrg creates a new Organization object based on the values in the OrgCreatedCmd.
func (cmd *OrgCreatedCmd) ToOrg() (o *Organization, err error) {
	if cmd.FullName == nil {
		err = allerror.NewInvalidParam("org fullname is empty")
		return
	}

	o = &Organization{
		Account:  cmd.Name,
		Fullname: cmd.FullName,
		Desc:     cmd.Description,
		Website:  cmd.Website,
		Owner:    cmd.Owner,
		AvatarId: cmd.AvatarId,
		Type:     domain.UserTypeOrganization,
	}

	return
}

// OrgListOptions represents the options for listing organizations.
type OrgListOptions struct {
	Username string // filter by member username
	Owner    string // filter by owner name
}

// ToApprove converts an OrgMember, expiry time, and inviter account to an Approval struct.
func ToApprove(member OrgMember, expiry int64, inviter primitive.Account) Approve {
	return Approve{
		OrgName:  member.OrgName,
		Username: member.Username,
		Role:     member.Role,
		ExpireAt: utils.Expiry(expiry),
		Inviter:  inviter,
	}
}

// ApproveStatus represents the status of an approval for a member to join an organization.
type ApproveStatus = string

const (
	// ApproveStatusPending represents the pending status for approval.
	ApproveStatusPending ApproveStatus = "pending"

	// ApproveStatusApproved represents the approved status for approval.
	ApproveStatusApproved ApproveStatus = "approved"

	// ApproveStatusRejected represents the rejected status for approval.
	ApproveStatusRejected ApproveStatus = "rejected"
)

// OrgMember represents an organization member with its details.
type OrgMember struct {
	Id        primitive.Identity `json:"id"`
	Username  primitive.Account  `json:"user_name"`
	UserId    primitive.Identity `json:"user_id"`
	OrgName   primitive.Account  `json:"org_name"`
	OrgId     primitive.Identity `json:"org_id"`
	Role      primitive.Role     `json:"role"`
	Type      InviteType         `json:"type"`
	CreatedAt int64              `json:"created_at"`
	UpdatedAt int64              `json:"updated_at"`
	Version   int
}

// MemberRequest represents a request to manage organization membership.
type MemberRequest struct {
	Id primitive.Identity `json:"id"`

	Username  primitive.Account  `json:"user_name"`
	UserId    primitive.Identity `json:"user_id"`
	OrgName   primitive.Account  `json:"org_name"`
	OrgId     primitive.Identity `json:"org_id"`
	Role      primitive.Role     `json:"role"`
	Status    ApproveStatus      `json:"status"`
	By        string             `json:"by"`
	Msg       string             `json:"msg"`
	CreatedAt int64              `json:"created_at"`
	UpdatedAt int64              `json:"updated_at"`
	Version   int
}

// Approve represents an approval for a member to join an organization.
type Approve struct {
	Id primitive.Identity `json:"id"`

	Username  primitive.Account  `json:"user_name"`
	UserId    primitive.Identity `json:"user_id"`
	OrgName   primitive.Account  `json:"org_name"`
	OrgId     primitive.Identity `json:"org_id"`
	Role      primitive.Role     `json:"role"`
	ExpireAt  int64              `json:"expire_at"`
	Inviter   primitive.Account  `json:"Inviter"`
	InviterId primitive.Identity `json:"InviterId"`
	Status    ApproveStatus      `json:"status"`
	By        string             `json:"by"`
	Msg       string             `json:"msg"`
	CreatedAt int64              `json:"created_at"`
	UpdatedAt int64              `json:"updated_at"`
	Version   int
}

// ToMember converts an Approve struct to an OrgMember struct.
func (a Approve) ToMember() OrgMember {
	return OrgMember{
		Username: a.Username,
		UserId:   a.UserId,
		OrgName:  a.OrgName,
		OrgId:    a.OrgId,
		Role:     a.Role,
	}
}

// Validate validates the fields of the OrgInviteMemberCmd struct.
func (cmd OrgInviteMemberCmd) Validate() error {
	if cmd.Account == nil {
		return allerror.NewInvalidParam("invalid account")
	}

	if cmd.Org == nil {
		return allerror.NewInvalidParam("invalid org")
	}

	if cmd.Actor == nil {
		return allerror.NewInvalidParam("invalid actor")
	}
	return nil
}

// OrgRemoveMemberCmd represents a command to remove a member from an organization.
type OrgRemoveMemberCmd struct {
	Actor   primitive.Account
	Account primitive.Account
	Org     primitive.Account
	Msg     string
}

// Validate checks if the command is valid.
func (cmd OrgRemoveMemberCmd) Validate() error {
	if cmd.Account == nil {
		return allerror.NewInvalidParam("invalid account")
	}

	if cmd.Org == nil {
		return allerror.NewInvalidParam("invalid org")
	}

	if cmd.Actor == nil {
		return allerror.NewInvalidParam("invalid actor")
	}

	return nil
}

// ToMember converts the command to an OrgMember.
func (cmd OrgRemoveMemberCmd) ToMember() OrgMember {
	return OrgMember{
		Username: cmd.Account,
		OrgName:  cmd.Org,
	}
}

// OrgEditMemberCmd represents a command to edit a member's role in an organization.
type OrgEditMemberCmd struct {
	Actor   primitive.Account
	Account primitive.Account
	Org     primitive.Account
	Role    primitive.Role
}

// OrgInviteMemberCmd represents a command to invite a member to an organization.
type OrgInviteMemberCmd struct {
	Actor   primitive.Account
	Account primitive.Account
	Org     primitive.Account
	Role    primitive.Role
	Msg     string
}

// ToApprove converts the command to an Approval.
func (cmd OrgInviteMemberCmd) ToApprove(expire int64) *Approve {
	return &Approve{
		OrgName:  cmd.Org,
		Username: cmd.Account,
		Role:     cmd.Role,
		Status:   ApproveStatusPending,
		Inviter:  cmd.Actor,
		ExpireAt: utils.Expiry(expire),
		Msg:      cmd.Msg,
	}
}

// OrgAddMemberCmd represents a command to add a member to an organization.
type OrgAddMemberCmd struct {
	User   primitive.Account
	UserId primitive.Identity
	Org    primitive.Account
	OrgId  primitive.Identity
	Type   InviteType
	Role   primitive.Role
	Msg    string
}

// Validate checks if the command is valid.
func (cmd OrgAddMemberCmd) Validate() error {
	if cmd.User == nil {
		return allerror.NewInvalidParam("invalid user")
	}

	if cmd.Org == nil {
		return allerror.NewInvalidParam("invalid org")
	}

	return nil
}

// ToMember converts the command to an OrgMember.
func (cmd OrgAddMemberCmd) ToMember() OrgMember {
	return OrgMember{
		Username: cmd.User,
		UserId:   cmd.UserId,
		OrgName:  cmd.Org,
		OrgId:    cmd.OrgId,
		Role:     primitive.Role(cmd.Role),
		Type:     cmd.Type,
	}
}

// OrgRemoveInviteCmd represents a command to remove an invite from an organization.
type OrgRemoveInviteCmd = OrgRemoveMemberCmd

// OrgRequestMemberCmd represents a command to request membership in an organization.
type OrgRequestMemberCmd struct {
	OrgNormalCmd
	Msg string
}

// ToMemberRequest converts the command to a MemberRequest.
func (o *OrgRequestMemberCmd) ToMemberRequest(role primitive.Role) *MemberRequest {
	return &MemberRequest{
		Username: o.Actor,
		OrgName:  o.Org,
		Role:     role,
		Status:   ApproveStatusPending,
		Msg:      o.Msg,
	}
}

// OrgCancelRequestMemberCmd represents a command to cancel a membership request in an organization.
type OrgCancelRequestMemberCmd struct {
	Actor     primitive.Account
	Requester primitive.Account
	Org       primitive.Account
	Msg       string
}

// Validate checks if the command is valid.
func (cmd OrgCancelRequestMemberCmd) Validate() error {
	if cmd.Requester == nil {
		return allerror.NewInvalidParam("invalid requester")
	}

	if cmd.Org == nil {
		return allerror.NewInvalidParam("invalid org")
	}

	if cmd.Actor == nil {
		return allerror.NewInvalidParam("invalid actor")
	}

	return nil
}

// OrgApproveRequestMemberCmd represents a command to approve a membership request in an organization.
type OrgApproveRequestMemberCmd = OrgCancelRequestMemberCmd

// OrgAcceptInviteCmd represents a command to accept an invitation to join an organization.
type OrgAcceptInviteCmd = OrgRemoveInviteCmd

// ToListReqCmd converts the command to a list of member requests.
func (cmd OrgApproveRequestMemberCmd) ToListReqCmd() *OrgMemberReqListCmd {
	return &OrgMemberReqListCmd{
		OrgNormalCmd: OrgNormalCmd{
			Actor: cmd.Actor,
			Org:   cmd.Org,
		},
		Requester: cmd.Requester,
		Status:    ApproveStatusPending,
	}
}

// Validate checks if the command is valid.
func (cmd OrgNormalCmd) Validate() error {
	if cmd.Actor == nil {
		return allerror.NewInvalidParam("invalid actor")
	}

	if cmd.Org == nil {
		return allerror.NewInvalidParam("invalid org")
	}

	return nil
}

// OrgNormalCmd represents a normal command with actor and org fields.
type OrgNormalCmd struct {
	Actor primitive.Account
	Org   primitive.Account
}

// OrgInvitationListCmd represents a command to list invitations in an organization.
type OrgInvitationListCmd struct {
	OrgNormalCmd
	Inviter primitive.Account
	Invitee primitive.Account
	Status  ApproveStatus
}

// OrgMemberReqListCmd represents a command to list member requests in an organization.
type OrgMemberReqListCmd struct {
	OrgNormalCmd
	Requester primitive.Account
	Status    ApproveStatus
}

// Validate checks if the command is valid.
func (cmd OrgMemberReqListCmd) Validate() error {
	if cmd.Actor == nil {
		return allerror.NewInvalidParam("invalid actor")
	}

	if cmd.Org == nil && cmd.Requester == nil {
		return allerror.NewInvalidParam("When list member requests, org_name/requester can't be all empty")
	}

	if cmd.Status != "" && cmd.Status != ApproveStatusPending && cmd.Status !=
		ApproveStatusApproved && cmd.Status != ApproveStatusRejected {
		return allerror.NewInvalidParam("invalid status")
	}

	return nil
}

// Validate checks if the command is valid.
func (cmd OrgInvitationListCmd) Validate() error {
	if cmd.Actor == nil {
		return allerror.NewInvalidParam("invalid actor")
	}

	count := 0
	if cmd.Org != nil {
		count++
	}

	if cmd.Invitee != nil {
		count++
	}

	if cmd.Inviter != nil {
		count++
	}

	if count > 1 {
		return allerror.NewInvalidParam("Only one of the org_name/invitee/inviter can be used")
	}

	if count == 0 {
		return allerror.NewInvalidParam("When list member invitation, org_name/invitee/inviter can't be all empty")
	}

	if cmd.Status != "" && cmd.Status != ApproveStatusPending && cmd.Status !=
		ApproveStatusApproved && cmd.Status != ApproveStatusRejected {
		return allerror.NewInvalidParam("invalid status")
	}

	return nil
}

// ToMember converts the command to an organization member.
func (cmd OrgInviteMemberCmd) ToMember() OrgMember {
	return OrgMember{
		Username: cmd.Account,
		Role:     primitive.Role(cmd.Role),
		OrgName:  cmd.Org,
	}
}
