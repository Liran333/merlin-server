package controller

import (
	"fmt"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/user/app"
	"github.com/openmerlin/merlin-server/user/domain"
)

type userBasicInfoUpdateRequest struct {
	AvatarId *string `json:"avatar_id"`
	Desc     *string `json:"description"`
	Email    *string `json:"email"`
	Fullname *string `json:"fullname"`
}

func (req *userBasicInfoUpdateRequest) toCmd() (
	cmd app.UpdateUserBasicInfoCmd,
	err error,
) {
	if req.Desc != nil {
		if cmd.Desc, err = primitive.NewMSDDesc(*req.Desc); err != nil {
			return
		}
	}

	if req.AvatarId != nil {
		if cmd.AvatarId, err = primitive.NewAvatarId(*req.AvatarId); err != nil {
			return
		}
	}

	if req.Email != nil {
		if cmd.Email, err = primitive.NewEmail(*req.Email); err != nil {
			return
		}
	}

	if req.Fullname != nil {
		if cmd.Fullname, err = primitive.NewMSDFullname(*req.Fullname); err != nil {
			return
		}
	}

	if req.AvatarId == nil && req.Desc == nil && req.Email == nil && req.Fullname == nil {
		err = fmt.Errorf("all param are empty")
		return
	}

	return
}

type userCreateRequest struct {
	Account  string `json:"account" binding:"required"`
	Fullname string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Bio      string `json:"bio"`
	AvatarId string `json:"avatar_id"`
}

func (req *userCreateRequest) toCmd() (cmd domain.UserCreateCmd, err error) {
	if cmd.Account, err = primitive.NewAccount(req.Account); err != nil {
		return
	}

	if cmd.Email, err = primitive.NewEmail(req.Email); err != nil {
		return
	}

	if cmd.Desc, err = primitive.NewMSDDesc(req.Bio); err != nil {
		return
	}

	if cmd.AvatarId, err = primitive.NewAvatarId(req.AvatarId); err != nil {
		return
	}

	if cmd.Fullname, err = primitive.NewMSDFullname(req.Fullname); err != nil {
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

func (req *tokenCreateRequest) toCmd(user domain.Account) (cmd domain.TokenCreatedCmd, err error) {
	if cmd.Permission, err = primitive.NewTokenPerm(req.Perm); err != nil {
		return
	}

	if cmd.Name, err = primitive.NewTokenName(req.Name); err != nil {
		return
	}

	cmd.Account = user

	return
}

type userToken struct {
	app.TokenDTO
}

// reqToGetUserInfo
type reqToGetUserInfo struct {
	Account string `form:"account"`
}

func (req *reqToGetUserInfo) toAccount() (primitive.Account, error) {
	if req.Account == "" {
		return nil, nil
	}

	return primitive.NewAccount(req.Account)
}
