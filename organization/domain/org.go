package domain

import (
	"fmt"

	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	user "github.com/openmerlin/merlin-server/user/domain"
	"github.com/openmerlin/merlin-server/utils"
)

type Organization struct {
	Id                string            `json:"id"`
	Name              primitive.Account `json:"name"`
	FullName          string            `json:"full_name"`
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
	Approves          []Approve         `json:"approves"`

	Version int
}

type OrgCreatedCmd struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
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
	Actor       primitive.Account
	OrgName     primitive.Account
	FullName    string
	Description string
	Website     string
	AvatarId    string
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

	return nil
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

type OrgMember struct {
	Id       string  `json:"id"`
	Username string  `json:"user_name"`
	OrgName  string  `json:"org_name"`
	Role     OrgRole `json:"role"`

	Version int
}

type Approve struct {
	Username string  `json:"user_name"`
	OrgName  string  `json:"org_name"`
	Role     OrgRole `json:"role"`
	ExpireAt int64   `json:"expire_at"`
	Inviter  string  `json:"Inviter"`
}

func (a Approve) ToMember() OrgMember {
	return OrgMember{
		Username: a.Username,
		OrgName:  a.OrgName,
		Role:     a.Role,
	}
}

func (org *Organization) AddInvite(member OrgMember, expiry int64, inviter string) error {
	if org.Approves == nil {
		org.Approves = make([]Approve, 1)
	} else {
		for _, m := range org.Approves {
			if member.Username == m.Username {
				return fmt.Errorf("member already exists")
			}
		}
	}

	org.Approves = append(org.Approves, ToApprove(member, expiry, inviter))
	return nil
}

func (org *Organization) RemoveInvite(member OrgMember) (approve Approve, err error) {
	if org.Approves == nil {
		err = fmt.Errorf("no approve to remove")
		return
	}

	for i, m := range org.Approves {
		if member.Username == m.Username {
			org.Approves = append(org.Approves[:i], org.Approves[i+1:]...)
			approve = m
			return
		}
	}

	err = fmt.Errorf("the target approve record not found")
	return
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
}
type OrgAddMemberCmd = OrgRemoveMemberCmd
type OrgRemoveInviteCmd = OrgRemoveMemberCmd

type OrgNormalCmd struct {
	Actor primitive.Account
	Org   primitive.Account
}

func (cmd OrgInviteMemberCmd) ToMember() OrgMember {
	return OrgMember{
		Username: cmd.Account.Account(),
		Role:     OrgRole(cmd.Role),
		OrgName:  cmd.Org.Account(),
	}
}
