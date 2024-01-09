package app

import (
	"github.com/sirupsen/logrus"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/organization/domain"
	userapp "github.com/openmerlin/merlin-server/user/app"
)

type OrganizationDTO struct {
	Id           string `json:"id"`
	Name         string `json:"account"`
	FullName     string `json:"fullname"`
	AvatarId     string `json:"avatar_id"`
	Description  string `json:"description"`
	CreatedAt    int64  `json:"created_at"`
	Website      string `json:"website"`
	Owner        string `json:"owner"`
	DefaultRole  string `json:"default_role"`
	AllowRequest bool   `json:"allow_request"`
}

type ApproveDTO struct {
	OrgName   string `json:"org_name"`
	UserName  string `json:"user_name"`
	Role      string `json:"role"`
	ExpiresAt int64  `json:"expires_at"`
	Fullname  string `json:"fullname"`
	Inviter   string `json:"inviter"`
	Status    string `json:"status"`
	By        string `json:"by"`
	Msg       string `json:"msg"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func ToApproveDTO(m *domain.Approve, user userapp.UserService) ApproveDTO {
	var fullname string
	u, err := user.GetByAccount(primitive.CreateAccount(m.Username), false)
	if err != nil {
		logrus.Warnf("failed to get fullname for %s, err:%s", m.Username, err)
		fullname = ""
	}
	fullname = u.Fullname

	return ApproveDTO{
		OrgName:   m.OrgName,
		UserName:  m.Username,
		Role:      string(m.Role),
		ExpiresAt: m.ExpireAt, // will expire in 14 days
		Inviter:   m.Inviter,
		Status:    string(m.Status),
		Fullname:  fullname,
		Msg:       m.Msg,
		By:        m.By,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func ToApprove(dto *ApproveDTO) domain.Approve {
	return domain.Approve{
		OrgName:   dto.OrgName,
		Username:  dto.UserName,
		Role:      domain.OrgRole(dto.Role),
		ExpireAt:  dto.ExpiresAt,
		Inviter:   dto.Inviter,
		Status:    domain.ApproveStatus(dto.Status),
		By:        dto.By,
		Msg:       dto.Msg,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}

type MemberRequestDTO struct {
	Username  string `json:"username"`
	Fullname  string `json:"fullname"`
	Role      string `json:"role"`
	OrgName   string `json:"org_name"`
	Status    string `json:"status"`
	By        string `json:"by"`
	Msg       string `json:"msg"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func ToMemberRequestDTO(m *domain.MemberRequest, user userapp.UserService) MemberRequestDTO {
	var fullname string
	u, err := user.GetByAccount(primitive.CreateAccount(m.Username), false)
	if err != nil {
		logrus.Warnf("failed to get fullname for %s, err:%s", m.Username, err)
		fullname = ""
	}
	fullname = u.Fullname

	return MemberRequestDTO{
		Username:  m.Username,
		Role:      string(m.Role),
		OrgName:   m.OrgName,
		Status:    string(m.Status),
		Fullname:  fullname,
		By:        m.By,
		Msg:       m.Msg,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

type MemberDTO struct {
	OrgName     string `json:"org_name"`
	OrgFullName string `json:"org_fullname"`
	UserName    string `json:"user_name"`
	Role        string `json:"role"`
}

type OrgListOptions struct {
	Page     int
	PageSize int
	Owner    primitive.Account
	Member   primitive.Account
}

func ToDTO(org *domain.Organization, role domain.OrgRole) OrganizationDTO {
	if org.DefaultRole == "" {
		org.DefaultRole = role
	}
	return OrganizationDTO{
		Id:           org.PlatformId,
		Name:         org.Name.Account(),
		FullName:     org.FullName,
		Description:  org.Description,
		Website:      org.Website,
		CreatedAt:    org.CreatedAt,
		Owner:        org.Owner.Account(),
		AvatarId:     org.AvatarId.AvatarId(),
		DefaultRole:  string(org.DefaultRole),
		AllowRequest: org.AllowRequest,
	}
}

func ToMemberDTO(member *domain.OrgMember) MemberDTO {
	return MemberDTO{
		UserName: member.Username,
		Role:     string(member.Role),
	}
}
