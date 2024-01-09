package app

import (
	orgapp "github.com/openmerlin/merlin-server/organization/app"
	userapp "github.com/openmerlin/merlin-server/user/app"
)

type userType int

const (
	userTypeUser userType = iota
	userTypeOrganization
)

type UserDTO struct {
	Id           string `json:"id"`
	Name         string `json:"account"`
	FullName     string `json:"fullname"`
	AvatarId     string `json:"avatar_id"`
	Email        string `json:"email,omitempty"`
	Description  string `json:"description"`
	CreatedAt    int64  `json:"created_at"`
	Website      string `json:"website,omitempty"`
	Owner        string `json:"owner,omitempty"`
	Type         int    `json:"type"`
	AllowRequest bool   `json:"allow_request,omitempty"`
	DefaultRole  string `json:"default_role,omitempty"`
}

func FromOrgDTO(o orgapp.OrganizationDTO) UserDTO {
	return UserDTO{
		Id:           o.Id,
		Name:         o.Name,
		FullName:     o.FullName,
		AvatarId:     o.AvatarId,
		Email:        "",
		Description:  o.Description,
		CreatedAt:    o.CreatedAt,
		Website:      o.Website,
		Owner:        o.Owner,
		Type:         int(userTypeOrganization),
		AllowRequest: o.AllowRequest,
		DefaultRole:  o.DefaultRole,
	}
}

func FromUserDTO(u userapp.UserDTO) UserDTO {
	return UserDTO{
		Id:          u.Id,
		Name:        u.Account,
		FullName:    u.Fullname,
		AvatarId:    u.AvatarId,
		Email:       u.Email,
		Description: u.Bio,
		CreatedAt:   u.CreatedAt,
		Website:     "",
		Owner:       "",
		Type:        int(userTypeUser),
	}

}
