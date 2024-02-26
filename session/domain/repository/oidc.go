/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package repository provides interfaces and types for user-related functionality.
package repository

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/user/domain"
)

const (
	// ErrorCodeError represents the error code for email code errors.
	ErrorCodeError = "email_code_error"

	// ErrorCodeInvalid represents the error code for invalid email codes.
	ErrorCodeInvalid = "email_code_invalid"

	// ErrorUserDuplicateBind represents the error code for duplicate user bindings.
	ErrorUserDuplicateBind = "email_user_duplicate_bind"

	// ErrorEmailDuplicateBind represents the error code for duplicate email bindings.
	ErrorEmailDuplicateBind = "email_email_duplicate_bind"

	// ErrorEmailDuplicateSend represents the error code for duplicate email sending.
	ErrorEmailDuplicateSend = "email_email_duplicate_send"
)

// UserInfo represents the user information structure.
type UserInfo struct {
	Desc     primitive.MSDDesc
	Name     domain.Account
	Email    primitive.Email
	AvatarId primitive.AvatarId
	Fullname primitive.MSDFullname
	Phone    primitive.Phone
	UserId   string
}

// Login represents the login structure containing user information and tokens.
type Login struct {
	UserInfo

	IDToken     string
	AccessToken string
}

// OIDCAdapter is an interface for OpenID Connect adapters.
type OIDCAdapter interface {
	GetByCode(code, redirectURI string) (Login, error)

	SendBindEmail(email, capt string) (err error)
	VerifyBindEmail(email, passCode, userid string) (err error)

	PrivacyRevoke(userid string) error
}
