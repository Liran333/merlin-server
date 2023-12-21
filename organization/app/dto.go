package app

import (
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
	Approves    []ApproveDTO `json:"-"`
}

type ApproveDTO struct {
	OrgName   string `json:"org_name"`
	UserName  string `json:"user_name"`
	Role      string `json:"role"`
	ExpiresAt int64  `json:"expires_at"`
	Fullname  string `json:"fullname"`
	Inviter   string `json:"inviter"`
}

func ToApproveDTO(m domain.Approve) ApproveDTO {
	return ApproveDTO{
		OrgName:   m.OrgName,
		UserName:  m.Username,
		Role:      string(m.Role),
		ExpiresAt: m.ExpireAt, // will expire in 14 days
		Inviter:   m.Inviter,
	}
}

type MemberDTO struct {
	OrgName     string `json:"org_name"`
	OrgFullName string `json:"org_full_name"`
	UserName    string `json:"user_name"`
	Role        string `json:"role"`
}

type OrgListOptions struct {
	Page     int
	PageSize int
	Owner    primitive.Account
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
