package domain

import (
	"fmt"

	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/utils"
)

type InviteType = string

const (
	InviteTypeInvite  InviteType = "invite"
	InviteTypeRequest InviteType = "request"
)

type Organization struct {
	Id                primitive.Identity    `json:"id"`
	Name              primitive.Account     `json:"name"`
	Fullname          primitive.MSDFullname `json:"fullname"`
	AvatarId          primitive.AvatarId    `json:"avatar_id"`
	PlatformId        int64                 `json:"platform_id"`
	Description       primitive.MSDDesc     `json:"description"`
	CreatedAt         int64                 `json:"created_at"`
	UpdatedAt         int64                 `json:"updated_at"`
	Website           string                `json:"website"`
	Owner             primitive.Account     `json:"owner"`
	OwnerId           primitive.Identity    `json:"owner_id"`
	WriteTeamId       int64                 `json:"write_team_id"`
	ReadTeamId        int64                 `json:"read_team_id"`
	OwnerTeamId       int64                 `json:"owner_team_id"`
	ContributorTeamId int64                 `json:"contributor_team_id"`
	Type              int                   `json:"type"`
	DefaultRole       OrgRole               `json:"default_role"`
	AllowRequest      bool                  `json:"allow_request"`

	Version int
}

type OrgCreatedCmd struct {
	Name        primitive.Account     `json:"name"`
	FullName    primitive.MSDFullname `json:"fullname"`
	Description primitive.MSDDesc     `json:"description"`
	Website     string                `json:"website"`
	AvatarId    primitive.AvatarId    `json:"avatar_id"`
	Owner       primitive.Account     `json:"owner"`
}

type OrgDeletedCmd struct {
	Actor primitive.Account
	Name  primitive.Account
}

func (cmd OrgDeletedCmd) Validate() error {
	if cmd.Name == nil {
		return allerror.NewInvalidParam("invalid org name")
	}

	if cmd.Actor == nil {
		return allerror.NewInvalidParam("invalid actor name")
	}

	return nil
}

type OrgUpdatedBasicInfoCmd struct {
	Actor        primitive.Account
	OrgName      primitive.Account
	AllowRequest *bool
	DefaultRole  string
	FullName     string
	Description  string
	Website      string
	AvatarId     string
}

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

	if cmd.DefaultRole != "" && cmd.DefaultRole != string(OrgRoleAdmin) && cmd.DefaultRole != string(OrgRoleReader) && cmd.DefaultRole != string(OrgRoleWriter) && cmd.DefaultRole != string(OrgRoleContributor) {
		return allerror.NewInvalidParam("invalid default role")
	}

	return nil
}

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

	if cmd.Description != o.Description.MSDDesc() && cmd.Description != "" {
		o.Description = primitive.CreateMSDDesc(cmd.Description)
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

	if cmd.DefaultRole != "" && cmd.DefaultRole != string(o.DefaultRole) {
		o.DefaultRole = OrgRole(cmd.DefaultRole)
		change = true
	}

	return
}

func (cmd *OrgCreatedCmd) ToOrg() (o *Organization, err error) {
	o = &Organization{
		Name:        cmd.Name,
		Fullname:    cmd.FullName,
		Description: cmd.Description,
		Website:     cmd.Website,
		Owner:       cmd.Owner,
		AvatarId:    cmd.AvatarId,
	}

	return
}

type OrgListOptions struct {
	Username string // filter by member user name
	Owner    string // filter by owner name
}

func ToApprove(member OrgMember, expiry int64, inviter primitive.Account) Approve {
	return Approve{
		OrgName:  member.OrgName,
		Username: member.Username,
		Role:     member.Role,
		ExpireAt: utils.Expiry(expiry),
		Inviter:  inviter,
	}
}

type OrgRole = string

const (
	OrgRoleContributor OrgRole = "contributor" // in contributor team
	OrgRoleReader      OrgRole = "read"        // in read team
	OrgRoleWriter      OrgRole = "write"       // in write team
	OrgRoleAdmin       OrgRole = "admin"       // in owner team
)

type ApproveStatus = string

const (
	ApproveStatusPending  ApproveStatus = "pending"
	ApproveStatusApproved ApproveStatus = "approved"
	ApproveStatusRejected ApproveStatus = "rejected"
)

type OrgMember struct {
	Id        primitive.Identity `json:"id"`
	Username  primitive.Account  `json:"user_name"`
	UserId    primitive.Identity `json:"user_id"`
	OrgName   primitive.Account  `json:"org_name"`
	OrgId     primitive.Identity `json:"org_id"`
	Role      OrgRole            `json:"role"`
	Type      InviteType         `json:"type"`
	CreatedAt int64              `json:"created_at"`
	UpdatedAt int64              `json:"updated_at"`
	Version   int
}

type MemberRequest struct {
	Id primitive.Identity `json:"id"`

	Username  primitive.Account  `json:"user_name"`
	UserId    primitive.Identity `json:"user_id"`
	OrgName   primitive.Account  `json:"org_name"`
	OrgId     primitive.Identity `json:"org_id"`
	Role      OrgRole            `json:"role"`
	Status    ApproveStatus      `json:"status"`
	By        string             `json:"by"`
	Msg       string             `json:"msg"`
	CreatedAt int64              `json:"created_at"`
	UpdatedAt int64              `json:"updated_at"`
	Version   int
}

type Approve struct {
	Id primitive.Identity `json:"id"`

	Username  primitive.Account  `json:"user_name"`
	UserId    primitive.Identity `json:"user_id"`
	OrgName   primitive.Account  `json:"org_name"`
	OrgId     primitive.Identity `json:"org_id"`
	Role      OrgRole            `json:"role"`
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

func (a Approve) ToMember() OrgMember {
	return OrgMember{
		Username: a.Username,
		UserId:   a.UserId,
		OrgName:  a.OrgName,
		OrgId:    a.OrgId,
		Role:     a.Role,
	}
}

func (cmd OrgInviteMemberCmd) Validate() error {
	if cmd.Role != string(OrgRoleContributor) &&
		cmd.Role != string(OrgRoleReader) &&
		cmd.Role != string(OrgRoleWriter) &&
		cmd.Role != string(OrgRoleAdmin) {
		return allerror.NewInvalidParam(fmt.Sprintf("invalid role: %s", cmd.Role))
	}

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

type OrgRemoveMemberCmd struct {
	Actor   primitive.Account
	Account primitive.Account
	Org     primitive.Account
	Msg     string
}

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

func (cmd OrgRemoveMemberCmd) ToMember() OrgMember {
	return OrgMember{
		Username: cmd.Account,
		OrgName:  cmd.Org,
	}
}

type OrgEditMemberCmd struct {
	Actor   primitive.Account
	Account primitive.Account
	Org     primitive.Account
	Role    string
}
type OrgInviteMemberCmd struct {
	Actor   primitive.Account
	Account primitive.Account
	Org     primitive.Account
	Role    string
	Msg     string
}

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

type OrgAddMemberCmd struct {
	User   primitive.Account
	UserId primitive.Identity
	Org    primitive.Account
	OrgId  primitive.Identity
	Type   InviteType
	Role   OrgRole
	Msg    string
}

func (cmd OrgAddMemberCmd) Validate() error {
	if cmd.User == nil {
		return allerror.NewInvalidParam("invalid user")
	}

	if cmd.Org == nil {
		return allerror.NewInvalidParam("invalid org")
	}

	return nil
}

func (cmd OrgAddMemberCmd) ToMember() OrgMember {
	return OrgMember{
		Username: cmd.User,
		UserId:   cmd.UserId,
		OrgName:  cmd.Org,
		OrgId:    cmd.OrgId,
		Role:     OrgRole(cmd.Role),
		Type:     cmd.Type,
	}
}

type OrgRemoveInviteCmd = OrgRemoveMemberCmd

type OrgRequestMemberCmd struct {
	OrgNormalCmd
	Msg string
}

func (o *OrgRequestMemberCmd) ToMemberRequest(role OrgRole) *MemberRequest {
	return &MemberRequest{
		Username: o.Actor,
		OrgName:  o.Org,
		Role:     role,
		Status:   ApproveStatusPending,
		Msg:      o.Msg,
	}
}

type OrgCancelRequestMemberCmd struct {
	Actor     primitive.Account
	Requester primitive.Account
	Org       primitive.Account
	Msg       string
}

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

type OrgApproveRequestMemberCmd = OrgCancelRequestMemberCmd
type OrgAcceptInviteCmd = OrgRemoveInviteCmd

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

func (cmd OrgNormalCmd) Validate() error {
	if cmd.Actor == nil {
		return allerror.NewInvalidParam("invalid actor")
	}

	if cmd.Org == nil {
		return allerror.NewInvalidParam("invalid org")
	}

	return nil
}

type OrgNormalCmd struct {
	Actor primitive.Account
	Org   primitive.Account
}

type OrgInvitationListCmd struct {
	//TODO add sort and paginate
	OrgNormalCmd
	Inviter primitive.Account
	Invitee primitive.Account
	Status  ApproveStatus
}

type OrgMemberReqListCmd struct {
	//TODO add sort and paginate
	OrgNormalCmd
	Requester primitive.Account
	Status    ApproveStatus
}

func (cmd OrgMemberReqListCmd) Validate() error {
	if cmd.Actor == nil {
		return allerror.NewInvalidParam("invalid actor")
	}

	if cmd.Org == nil && cmd.Requester == nil {
		return allerror.NewInvalidParam("When list member requests, org_name/requester can't be all empty")
	}

	if cmd.Status != "" && cmd.Status != ApproveStatusPending && cmd.Status != ApproveStatusApproved && cmd.Status != ApproveStatusRejected {
		return allerror.NewInvalidParam("invalid status")
	}

	return nil
}

func (cmd OrgInvitationListCmd) Validate() error {
	if cmd.Actor == nil {
		return allerror.NewInvalidParam("invalid actor")
	}

	if cmd.Org == nil && cmd.Invitee == nil && cmd.Inviter == nil {
		return allerror.NewInvalidParam("When list invitation, org_name/invitee/inviter can't be all empty")
	}

	if cmd.Status != "" && cmd.Status != ApproveStatusPending && cmd.Status != ApproveStatusApproved && cmd.Status != ApproveStatusRejected {
		return allerror.NewInvalidParam("invalid status")
	}

	return nil
}

func (cmd OrgInviteMemberCmd) ToMember() OrgMember {
	return OrgMember{
		Username: cmd.Account,
		Role:     OrgRole(cmd.Role),
		OrgName:  cmd.Org,
	}
}
