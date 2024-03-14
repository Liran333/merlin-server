/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	userapp "github.com/openmerlin/merlin-server/user/app"
)

// UserDTO represents the user data transfer object for the userapp package.
type UserDTO = userapp.UserDTO

// CmdToCheck represents a command to check the session and user agent information.
type CmdToCheck struct {
	SessionDTO

	IP          string
	UserAgent   primitive.UserAgent
	AutoRefresh bool
}

// CmdToLogin represents a command to login with the provided code, redirect URI, and user agent.
type CmdToLogin struct {
	IP          string
	Code        string
	RedirectURI string
	UserAgent   primitive.UserAgent
}

// SessionDTO represents the session data transfer object containing login ID and CSRF token.
type SessionDTO struct {
	SessionId primitive.RandomId
	CSRFToken primitive.RandomId
}
