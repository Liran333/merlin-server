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
	UpdatedAt    int64  `json:"updated_at"`
	Website      string `json:"website"`
	Owner        string `json:"owner"`
	OwnerId      string `json:"owner_id"`
	DefaultRole  string `json:"default_role"`
	AllowRequest bool   `json:"allow_request"`
}

type ApproveDTO struct {
	Id        string `json:"id"`
	OrgName   string `json:"org_name"`
	OrgId     string `json:"org_id"`
	UserName  string `json:"user_name"`
	UserId    string `json:"user_id"`
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
	u, err := user.GetByAccount(m.Username, false)
	if err != nil {
		logrus.Warnf("failed to get fullname for %s, err:%s", m.Username, err)
		fullname = ""
	}
	fullname = u.Fullname

	return ApproveDTO{
		Id:        m.Id.Identity(),
		OrgName:   m.OrgName.Account(),
		OrgId:     m.OrgId.Identity(),
		UserName:  m.Username.Account(),
		UserId:    m.UserId.Identity(),
		Role:      m.Role,
		ExpiresAt: m.ExpireAt, // will expire in 14 days
		Inviter:   m.Inviter.Account(),
		Status:    string(m.Status),
		Fullname:  fullname,
		Msg:       m.Msg,
		By:        m.By,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

type MemberRequestDTO struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	UserId    string `json:"user_id"`
	Fullname  string `json:"fullname"`
	Role      string `json:"role"`
	OrgName   string `json:"org_name"`
	OrgId     string `json:"org_id"`
	Status    string `json:"status"`
	By        string `json:"by"`
	Msg       string `json:"msg"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func ToMemberRequestDTO(m *domain.MemberRequest, user userapp.UserService) MemberRequestDTO {
	var fullname string
	u, err := user.GetByAccount(m.Username, false)
	if err != nil {
		logrus.Warnf("failed to get fullname for %s, err:%s", m.Username, err)
		fullname = ""
	}
	fullname = u.Fullname

	return MemberRequestDTO{
		Id:        m.Id.Identity(),
		Username:  m.Username.Account(),
		UserId:    m.UserId.Identity(),
		Role:      m.Role,
		OrgName:   m.OrgName.Account(),
		OrgId:     m.OrgId.Identity(),
		Status:    string(m.Status),
		Fullname:  fullname,
		By:        m.By,
		Msg:       m.Msg,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

type MemberDTO struct {
	Id          string `json:"id"`
	OrgName     string `json:"org_name"`
	OrgId       string `json:"org_id"`
	OrgFullName string `json:"org_fullname"`
	UserName    string `json:"user_name"`
	UserId      string `json:"user_id"`
	Role        string `json:"role"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
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
		Id:           org.Id.Identity(),
		Name:         org.Name.Account(),
		FullName:     org.Fullname.MSDFullname(),
		Description:  org.Description.MSDDesc(),
		Website:      org.Website,
		CreatedAt:    org.CreatedAt,
		UpdatedAt:    org.UpdatedAt,
		Owner:        org.Owner.Account(),
		OwnerId:      org.OwnerId.Identity(),
		AvatarId:     org.AvatarId.AvatarId(),
		DefaultRole:  string(org.DefaultRole),
		AllowRequest: org.AllowRequest,
	}
}

func ToMemberDTO(member *domain.OrgMember) MemberDTO {
	return MemberDTO{
		Id:        member.Id.Identity(),
		UserName:  member.Username.Account(),
		UserId:    member.UserId.Identity(),
		Role:      string(member.Role),
		OrgName:   member.OrgName.Account(),
		OrgId:     member.OrgId.Identity(),
		CreatedAt: member.CreatedAt,
		UpdatedAt: member.UpdatedAt,
	}
}
