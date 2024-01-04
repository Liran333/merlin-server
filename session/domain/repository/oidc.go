package repository

import "github.com/openmerlin/merlin-server/user/domain"

type UserInfo struct {
	Bio      domain.Bio
	Name     domain.Account
	Email    domain.Email
	AvatarId domain.AvatarId
	Fullname string
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
