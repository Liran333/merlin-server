package repository

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/user/domain"
)

const (
	ErrorCodeError          = "email_code_error"
	ErrorCodeInvalid        = "email_code_invalid"
	ErrorUserDuplicateBind  = "email_user_duplicate_bind"
	ErrorEmailDuplicateBind = "email_email_duplicate_bind"
	ErrorEmailDuplicateSend = "email_email_duplicate_send"
)

type UserInfo struct {
	Desc     primitive.MSDDesc
	Name     domain.Account
	Email    primitive.Email
	AvatarId primitive.AvatarId
	Fullname primitive.MSDFullname
	Phone    primitive.Phone
	UserId   string
}

type Login struct {
	UserInfo

	IDToken     string
	AccessToken string
}

type OIDCAdapter interface {
	GetByCode(code, redirectURI string) (Login, error)

	SendBindEmail(email, capt string) (err error)
	VerifyBindEmail(email, passCode, userid string) (err error)

	PrivacyRevoke(userid string) error
}
