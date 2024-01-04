package app

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	userapp "github.com/openmerlin/merlin-server/user/app"
)

type UserDTO = userapp.UserDTO

type CmdToCheck struct {
	SessionDTO

	IP        string
	UserAgent primitive.UserAgent
}

type CmdToLogin struct {
	IP          string
	Code        string
	RedirectURI string
	UserAgent   primitive.UserAgent
}

type SessionDTO struct {
	LoginId   primitive.UUID
	CSRFToken primitive.UUID
}
