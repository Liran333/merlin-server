/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package app provides application-level functionality for the user domain.
package app

import (
	"io"

	sdk "github.com/openmerlin/merlin-sdk/user"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	orgdomain "github.com/openmerlin/merlin-server/organization/domain"
	"github.com/openmerlin/merlin-server/user/domain"
	"github.com/openmerlin/merlin-server/user/domain/repository"
)

// CmdToCreateUser represents the command to create a user.
type CmdToCreateUser = domain.UserCreateCmd

// CmdToListModels represents the command to list models.
type CmdToListModels = repository.ListOption

// UserInfoDTO represents the data transfer object for user information.
type UserInfoDTO struct {
	UserDTO
	IsDisableAdmin bool `json:"is_disable_admin"`
}

type UserPaginationDTO struct {
	Total  int       `json:"total"`
	Labels []UserDTO `json:"Labels"`
}

// UserDTO represents the data transfer object for a user.
type UserDTO struct {
	Id              string  `json:"id"`
	Name            string  `json:"account"`
	Fullname        string  `json:"fullname"`
	AvatarId        string  `json:"avatar_id"`
	Email           *string `json:"email,omitempty"`
	Phone           *string `json:"phone,omitempty"`
	Description     string  `json:"description"`
	CreatedAt       int64   `json:"created_at"`
	UpdatedAt       int64   `json:"updated_at"`
	Website         *string `json:"website,omitempty"`
	Owner           *string `json:"owner,omitempty"`
	OwnerId         *string `json:"owner_id,omitempty"`
	Type            int     `json:"type"`
	AllowRequest    *bool   `json:"allow_request,omitempty"`
	RequestDelete   bool    `json:"request_delete"`
	RequestDeleteAt int64   `json:"request_delete_at"`
	DefaultRole     string  `json:"default_role,omitempty"`
	IsAgreePrivacy  bool    `json:"-"`
}

// NewUserDTO creates a new UserDTO based on the given domain.User object.
func NewUserDTO(u *domain.User, actor primitive.Account) (dto UserDTO) {
	return newUserDTO(u, actor)
}

// ToDTO converts a domain.Organization object to a userapp.UserDTO object.
func ToDTO(org *orgdomain.Organization) UserDTO {
	return NewUserDTO(org, nil)
}

func newUserDTO(u *domain.User, actor primitive.Account) (dto UserDTO) {
	dto.Name = u.Account.Account()
	if u.AvatarId != nil {
		dto.AvatarId = u.AvatarId.URL()
	}

	if u.Desc != nil {
		dto.Description = u.Desc.AccountDesc()
	}

	if u.IsOrganization() {
		website := u.Website.Website()
		owner := u.Owner.Account()
		ownerId := u.OwnerId.Identity()
		allow := u.AllowRequest
		dto.Website = &website
		dto.Owner = &owner
		dto.OwnerId = &ownerId
		dto.AllowRequest = &allow
		dto.DefaultRole = u.DefaultRole.Role()
	} else {
		email := ""
		if u.Email != nil && actor == u.Account {
			email = u.Email.Email()
		}
		dto.Email = &email

		phone := ""
		if u.Phone != nil && actor == u.Account {
			phone = u.Phone.PhoneNumber()
		}
		dto.Phone = &phone
		dto.RequestDelete = u.RequestDelete
		dto.RequestDeleteAt = u.RequestDeleteAt
	}

	dto.Type = u.Type
	dto.Fullname = u.Fullname.AccountFullname()
	dto.CreatedAt = u.CreatedAt
	dto.UpdatedAt = u.UpdatedAt
	dto.Id = u.Id.Identity()
	dto.IsAgreePrivacy = u.IsAgreePrivacy

	return
}

// TokenDTO represents the data transfer object for a token.
type TokenDTO struct {
	Id         string `json:"id"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
	Expire     int64  `json:"expired"`
	Account    string `json:"account"`
	OwnerId    string `json:"owner_id"`
	Name       string `json:"name"`
	Permission string `json:"permission"`
	Token      string `json:"token"`
}

func newTokenDTO(t *domain.PlatformToken) (dto TokenDTO) {
	dto.CreatedAt = t.CreatedAt
	dto.UpdatedAt = t.UpdatedAt
	dto.Expire = t.Expire
	dto.Account = t.Account.Account()
	dto.Name = t.Name.TokenName()
	dto.Permission = t.Permission.TokenPerm()
	dto.Token = t.Token
	dto.Id = t.Id.Identity()
	dto.OwnerId = t.OwnerId.Identity()

	return

}

// AvatarDTO represents the data transfer object for an avatar.
type AvatarDTO = sdk.AvatarDTO

// ToAvatarDTO converts a domain.User object to an AvatarDTO.
func ToAvatarDTO(a *domain.User) (dto AvatarDTO) {
	dto.AvatarId = a.AvatarId.URL()
	dto.Name = a.Account.Account()
	return
}

// ToUserDTO is an interface for creating a new UserDTO.
type ToUserDTO interface {
	NewUserDTO() UserDTO
}

// UpdateUserBasicInfoCmd represents the command to update basic user information.
type UpdateUserBasicInfoCmd struct {
	Desc         primitive.AccountDesc
	AvatarId     primitive.Avatar
	Fullname     primitive.AccountFullname
	RevokeDelete bool
}

func (cmd *UpdateUserBasicInfoCmd) toUser(u *domain.User) (changed bool) {
	if cmd.AvatarId != nil && cmd.AvatarId.URL() != u.AvatarId.URL() {
		u.AvatarId = cmd.AvatarId
		changed = true
	}

	if cmd.Desc != nil && cmd.Desc.AccountDesc() != u.Desc.AccountDesc() {
		u.Desc = cmd.Desc
		changed = true
	}

	if cmd.Fullname != nil && u.Fullname.AccountFullname() != cmd.Fullname.AccountFullname() {
		u.Fullname = cmd.Fullname
		changed = true
	}

	if cmd.RevokeDelete {
		u.RequestDelete = false
		u.RequestDeleteAt = 0
		changed = true
	}

	return
}

// CmdToSendBindEmail represents the command to send a bind email.
type CmdToSendBindEmail struct {
	User  primitive.Account
	Email primitive.Email
	Capt  string
}

// CmdToVerifyBindEmail represents the command to verify a bind email.
type CmdToVerifyBindEmail struct {
	User     primitive.Account
	Email    primitive.Email
	PassCode string
}

// CmdToUploadAvatar represents the command to upload avatar.
type CmdToUploadAvatar struct {
	Image    io.Reader
	User     primitive.Account
	FileName string
}

type AvatarUrlDTO struct {
	URL string `json:"url"`
}
