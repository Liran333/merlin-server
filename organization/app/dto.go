package app

import (
	"fmt"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/organization/domain"
)

type OrganizationDTO struct {
	Id          string       `json:"id"`
	Name        string       `json:"name"`
	FullName    string       `json:"full_name"`
	AvatarId    string       `json:"avatar_id"`
	PlatformId  string       `json:"platform_id"`
	Description string       `json:"description"`
	CreatedAt   int64        `json:"created_at"`
	Website     string       `json:"website"`
	Owner       string       `json:"owner"`
	Approves    []ApproveDTO `json:"approves"`
}

type ApproveDTO struct {
	OrgName   string `json:"org_name"`
	UserName  string `json:"user_name"`
	Role      string `json:"role"`
	ExpiresAt int64  `json:"expires_at"`
}

func ToApproveDTO(m domain.Approve) ApproveDTO {
	return ApproveDTO{
		OrgName:   m.OrgName,
		UserName:  m.Username,
		Role:      string(m.Role),
		ExpiresAt: m.ExpireAt, // will expire in 14 days
	}
}

type MemberDTO struct {
	OrgName     string `json:"org_name"`
	OrgFullName string `json:"org_full_name"`
	UserName    string `json:"user_name"`
	Role        string `json:"role"`
}

type UpdateOrgBasicInfoCmd struct {
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Website     string `json:"website"`
	AvatarId    string `json:"avatar_id"`
}

type OrgListOptions struct {
	Page     int
	PageSize int
	Owner    primitive.Account
}

func (cmd OrgInviteMemberCmd) Validate() error {
	if cmd.Role != string(domain.OrgRoleOwner) &&
		cmd.Role != string(domain.OrgRoleReader) &&
		cmd.Role != string(domain.OrgRoleWriter) {
		return fmt.Errorf("invalid role: %s", cmd.Role)
	}

	if cmd.Account == nil {
		return fmt.Errorf("invalid account")
	}

	if cmd.Org == nil {
		return fmt.Errorf("invalid org")
	}

	return nil
}

func (cmd OrgInviteMemberCmd) ToMember() domain.OrgMember {
	return domain.OrgMember{
		Username: cmd.Account.Account(),
		Role:     domain.OrgRole(cmd.Role),
		OrgName:  cmd.Org.Account(),
	}
}

type OrgRemoveMemberCmd struct {
	Account primitive.Account
	Org     primitive.Account
}

func (cmd OrgRemoveMemberCmd) Validate() error {
	if cmd.Account == nil {
		return fmt.Errorf("invalid account")
	}

	if cmd.Org == nil {
		return fmt.Errorf("invalid org")
	}

	return nil
}

func (cmd OrgRemoveMemberCmd) ToMember() domain.OrgMember {
	return domain.OrgMember{
		Username: cmd.Account.Account(),
		OrgName:  cmd.Org.Account(),
	}
}

type OrgEditMemberCmd struct {
	Account primitive.Account
	Org     primitive.Account
	Role    string
}

func ToDTO(org *domain.Organization) OrganizationDTO {
	approveDTOs := make([]ApproveDTO, len(org.Approves))
	for i := range org.Approves {
		approveDTOs = append(approveDTOs, ToApproveDTO(org.Approves[i]))
	}
	return OrganizationDTO{
		Id:          org.PlatformId,
		Name:        org.Name.Account(),
		FullName:    org.FullName,
		Description: org.Description,
		Website:     org.Website,
		CreatedAt:   org.CreatedAt,
		Owner:       org.Owner.Account(),
		AvatarId:    org.AvatarId.AvatarId(),
		Approves:    approveDTOs,
	}
}

func ToMemberDTO(member *domain.OrgMember) MemberDTO {
	return MemberDTO{
		UserName: member.Username,
		Role:     string(member.Role),
	}
}

type OrgInviteMemberCmd struct {
	Account primitive.Account
	Org     primitive.Account
	Role    string
}
type OrgAddMemberCmd = OrgRemoveMemberCmd
type OrgRemoveInviteCmd = OrgRemoveMemberCmd
