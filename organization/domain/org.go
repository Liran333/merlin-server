package domain

import (
	"fmt"

	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	user "github.com/openmerlin/merlin-server/user/domain"
	"github.com/openmerlin/merlin-server/utils"
)

type InviteType string

const (
	InviteTypeInvite  InviteType = "invite"
	InviteTypeRequest InviteType = "request"
)

type Organization struct {
	Id                string            `json:"id"`
	Name              primitive.Account `json:"name"`
	FullName          string            `json:"fullname"`
	AvatarId          user.AvatarId     `json:"avatar_id"`
	PlatformId        string            `json:"platform_id"`
	Description       string            `json:"description"`
	CreatedAt         int64             `json:"created_at"`
	Website           string            `json:"website"`
	Owner             primitive.Account `json:"owner"`
	WriteTeamId       int64             `json:"write_team_id"`
	ReadTeamId        int64             `json:"read_team_id"`
	OwnerTeamId       int64             `json:"owner_team_id"`
	ContributorTeamId int64             `json:"contributor_team_id"`
	Type              int               `json:"type"`
	DefaultRole       OrgRole           `json:"default_role"`
	AllowRequest      bool              `json:"allow_request"`

	Version int
}

type OrgCreatedCmd struct {
	Name        string `json:"name"`
	FullName    string `json:"fullname"`
	Description string `json:"description"`
	Website     string `json:"website"`
	AvatarId    string `json:"avatar_id"`
	Owner       string `json:"owner"`
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
		o.AvatarId, err = user.NewAvatarId(cmd.AvatarId)
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

	if cmd.Description != o.Description && cmd.Description != "" {
		o.Description = cmd.Description
		change = true
	}

	if cmd.FullName != o.FullName && cmd.FullName != "" {
		o.FullName = cmd.FullName
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

func (cmd OrgCreatedCmd) Validate() error {
	if _, err := primitive.NewAccount(cmd.Name); err != nil {
		return allerror.NewInvalidParam("org name is invalid")
	}

	if _, err := user.NewAvatarId(cmd.AvatarId); err != nil {
		return allerror.NewInvalidParam(err.Error())
	}

	if _, err := primitive.NewAccount(cmd.Owner); err != nil {
		return allerror.NewInvalidParam("owner name is invalid")
	}

	if cmd.Website != "" && !utils.IsUrl(cmd.Website) {
		return allerror.NewInvalidParam("invalid website")
	}

	return nil
}

func (cmd *OrgCreatedCmd) ToOrg() *Organization {
	return &Organization{
		Name:        primitive.CreateAccount(cmd.Name),
		FullName:    cmd.FullName,
		Description: cmd.Description,
		Website:     cmd.Website,
		CreatedAt:   utils.Now(),
		Owner:       primitive.CreateAccount(cmd.Owner),
		AvatarId:    user.CreateAvatarId(cmd.AvatarId),
	}
}

type OrgListOptions struct {
	Username string // filter by member user name
	Owner    string // filter by owner name
}

func ToApprove(member OrgMember, expiry int64, inviter string) Approve {
	return Approve{
		OrgName:  member.OrgName,
		Username: member.Username,
		Role:     member.Role,
		ExpireAt: utils.Expiry(expiry),
		Inviter:  inviter,
	}
}

type OrgRole string

const (
	OrgRoleContributor OrgRole = "contributor" // in contributor team
	OrgRoleReader      OrgRole = "read"        // in read team
	OrgRoleWriter      OrgRole = "write"       // in write team
	OrgRoleAdmin       OrgRole = "admin"       // in owner team
)

type ApproveStatus string

const (
	ApproveStatusPending  ApproveStatus = "pending"
	ApproveStatusApproved ApproveStatus = "approved"
	ApproveStatusRejected ApproveStatus = "rejected"
)

type OrgMember struct {
	Id       string     `json:"id"`
	Username string     `json:"user_name"`
	OrgName  string     `json:"org_name"`
	Role     OrgRole    `json:"role"`
	Type     InviteType `json:"type"`

	Version int
}

type MemberRequest struct {
	Id string `json:"id"`

	Username  string        `json:"user_name"`
	OrgName   string        `json:"org_name"`
	Role      OrgRole       `json:"role"`
	Status    ApproveStatus `json:"status"`
	By        string        `json:"by"`
	Msg       string        `json:"msg"`
	CreatedAt int64         `json:"created_at"`
	UpdatedAt int64         `json:"updated_at"`
	Version   int
}

type Approve struct {
	Id string `json:"id"`

	Username  string        `json:"user_name"`
	OrgName   string        `json:"org_name"`
	Role      OrgRole       `json:"role"`
	ExpireAt  int64         `json:"expire_at"`
	Inviter   string        `json:"Inviter"`
	Status    ApproveStatus `json:"status"`
	By        string        `json:"by"`
	Msg       string        `json:"msg"`
	CreatedAt int64         `json:"created_at"`
	UpdatedAt int64         `json:"updated_at"`
	Version   int
}

func (a Approve) ToMember() OrgMember {
	return OrgMember{
		Username: a.Username,
		OrgName:  a.OrgName,
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
		Username: cmd.Account.Account(),
		OrgName:  cmd.Org.Account(),
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

func (cmd OrgInviteMemberCmd) ToApprove(expire int64) Approve {
	return Approve{
		OrgName:   cmd.Org.Account(),
		Username:  cmd.Account.Account(),
		Role:      OrgRole(cmd.Role),
		Status:    ApproveStatusPending,
		Inviter:   cmd.Actor.Account(),
		CreatedAt: utils.Now(),
		ExpireAt:  utils.Expiry(expire),
		Msg:       cmd.Msg,
	}
}

type OrgAddMemberCmd struct {
	User primitive.Account
	Org  primitive.Account
	Id   string
	Type InviteType
	Role OrgRole
	Msg  string
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
		Username: cmd.User.Account(),
		OrgName:  cmd.Org.Account(),
		Role:     OrgRole(cmd.Role),
		Type:     cmd.Type,
	}
}

type OrgRemoveInviteCmd = OrgRemoveMemberCmd

type OrgRequestMemberCmd struct {
	OrgNormalCmd
	Msg string
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
		Username: cmd.Account.Account(),
		Role:     OrgRole(cmd.Role),
		OrgName:  cmd.Org.Account(),
	}
}
