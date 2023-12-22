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
	Id          string `json:"id"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	AvatarId    string `json:"avatar_id"`
	Email       string `json:"email"`
	Description string `json:"description"`
	CreatedAt   int64  `json:"created_at"`
	Website     string `json:"website"`
	Owner       string `json:"owner"`
	Type        int    `json:"type"`
}

func FromOrgDTO(o orgapp.OrganizationDTO) UserDTO {
	return UserDTO{
		Id:          o.Id,
		Name:        o.Name,
		FullName:    o.FullName,
		AvatarId:    o.AvatarId,
		Email:       "",
		Description: o.Description,
		CreatedAt:   o.CreatedAt,
		Website:     o.Website,
		Owner:       o.Owner,
		Type:        int(userTypeOrganization),
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
