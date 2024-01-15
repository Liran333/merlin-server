package repository

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/user/domain"
)

type UserInfo struct {
	Desc     primitive.MSDDesc
	Name     domain.Account
	Email    primitive.Email
	AvatarId primitive.AvatarId
	Fullname primitive.MSDFullname
	UserId   string
}

type Login struct {
	UserInfo

	IDToken     string
	AccessToken string
}

type OIDCAdapter interface {
	GetByCode(code, redirectURI string) (Login, error)
}
