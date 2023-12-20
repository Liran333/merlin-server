package domain

import (
	"fmt"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	user "github.com/openmerlin/merlin-server/user/domain"
	"github.com/openmerlin/merlin-server/utils"
)

type Organization struct {
	Id          string            `json:"id"`
	Name        primitive.Account `json:"name"`
	FullName    string            `json:"full_name"`
	AvatarId    user.AvatarId     `json:"avatar_id"`
	PlatformId  string            `json:"platform_id"`
	Description string            `json:"description"`
	CreatedAt   int64             `json:"created_at"`
	Website     string            `json:"website"`
	Owner       primitive.Account `json:"owner"`
	WriteTeamId int64             `json:"write_team_id"`
	ReadTeamId  int64             `json:"read_team_id"`
	OwnerTeamId int64             `json:"owner_team_id"`
	AdminTeamId int64             `json:"admin_team_id"`
	Approves    []Approve         `json:"approves"`

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
	Name  string `json:"name"`
	Owner string `json:"owner"`
}

func (cmd OrgCreatedCmd) Validate() error {
	if _, err := primitive.NewAccount(cmd.Name); err != nil {
		return err
	}

	if _, err := user.NewAvatarId(cmd.AvatarId); err != nil {
		return err
	}

	if _, err := primitive.NewAccount(cmd.Owner); err != nil {
		return err
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

func ToApprove(member OrgMember, expiry int64) Approve {
	return Approve{
		OrgName:  member.OrgName,
		Username: member.Username,
		Role:     member.Role,
		ExpireAt: utils.Expiry(expiry),
	}
}

type OrgRole string

const (
	OrgRoleOwner  OrgRole = "owner"
	OrgRoleReader OrgRole = "read"
	OrgRoleWriter OrgRole = "write"
	OrgRoleAdmin  OrgRole = "admin"
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
}

func (a Approve) ToMember() OrgMember {
	return OrgMember{
		Username: a.Username,
		OrgName:  a.OrgName,
		Role:     a.Role,
	}
}

func (org *Organization) AddInvite(member OrgMember, expiry int64) error {
	if org.Approves == nil {
		org.Approves = make([]Approve, 1)
	} else {
		for _, m := range org.Approves {
			if member.Username == m.Username {
				return fmt.Errorf("member already exists")
			}
		}
	}

	org.Approves = append(org.Approves, ToApprove(member, expiry))
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
