package login

import (
	"github.com/openmerlin/merlin-server/user/domain"
)

type UserInfo struct {
	Name     domain.Account
	Email    domain.Email
	Bio      domain.Bio
	AvatarId domain.AvatarId
	UserId   string
}

type Login struct {
	UserInfo

	IDToken     string
	AccessToken string
}

type User interface {
	GetByCode(code, redirectURI string) (Login, error)
	GetByAccessToken(accessToken string) (UserInfo, error)
}
