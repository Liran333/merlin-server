package controller

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/user/app"
	"github.com/openmerlin/merlin-server/user/domain"
)

type userBasicInfoUpdateRequest struct {
	AvatarId *string `json:"avatar_id"`
	Bio      *string `json:"bio"`
	Email    *string `json:"email"`
}

func (req *userBasicInfoUpdateRequest) toCmd() (
	cmd app.UpdateUserBasicInfoCmd,
	err error,
) {
	if req.Bio != nil {
		if cmd.Bio, err = domain.NewBio(*req.Bio); err != nil {
			return
		}
	}

	if req.AvatarId != nil {
		cmd.AvatarId, err = domain.NewAvatarId(*req.AvatarId)
	}

	if req.Email != nil {
		cmd.Email, err = domain.NewEmail(*req.Email)
	}

	return
}

type userCreateRequest struct {
	Account  string `json:"account" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Bio      string `json:"bio"`
	AvatarId string `json:"avatar_id"`
}

func (req *userCreateRequest) toCmd() (cmd domain.UserCreateCmd, err error) {
	if cmd.Account, err = primitive.NewAccount(req.Account); err != nil {
		return
	}

	if cmd.Email, err = domain.NewEmail(req.Email); err != nil {
		return
	}

	if cmd.Bio, err = domain.NewBio(req.Bio); err != nil {
		return
	}

	if cmd.AvatarId, err = domain.NewAvatarId(req.AvatarId); err != nil {
		return
	}

	err = cmd.Validate()

	return
}

type userDetail struct {
	*app.UserDTO
}

type tokenCreateRequest struct {
	Name string `json:"name" binding:"required"`
	Perm string `json:"perm" binding:"required"`
}

type userToken struct {
	app.TokenDTO
}
