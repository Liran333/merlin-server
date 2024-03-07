/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package controller

import (
	"fmt"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/user/app"
	"github.com/openmerlin/merlin-server/user/domain"
)

type userBasicInfoUpdateRequest struct {
	AvatarId     *string `json:"avatar_id"`
	Desc         *string `json:"description"`
	Fullname     *string `json:"fullname"`
	RevokeDelete *bool   `json:"revoke_delete"`
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

	if req.Fullname != nil {
		if cmd.Fullname, err = primitive.NewMSDFullname(*req.Fullname); err != nil {
			return
		}
	}

	if req.RevokeDelete != nil {
		cmd.RevokeDelete = true
	}

	if req.AvatarId == nil && req.Desc == nil && req.Fullname == nil && req.RevokeDelete == nil {
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

func (req *tokenCreateRequest) action() string {
	return fmt.Sprintf("create a new platform token named %s", req.Name)
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

type bindEmailRequest struct {
	Email    string `json:"email" binding:"required"`
	PassCode string `json:"code" binding:"required"`
}

func (req *bindEmailRequest) action() string {
	return fmt.Sprintf("bind email %s", req.Email)
}

func (req *bindEmailRequest) toCmd(user domain.Account) (cmd app.CmdToVerifyBindEmail, err error) {
	if cmd.Email, err = primitive.NewEmail(req.Email); err != nil {
		return
	}

	cmd.PassCode = req.PassCode

	cmd.User = user

	return
}

type sendEmailRequest struct {
	Email string `json:"email" binding:"required"`
	Capt  string `json:"capt" binding:"required"`
}

func (req *sendEmailRequest) action() string {
	return fmt.Sprintf("send email verify code to %s", req.Email)
}

func (req *sendEmailRequest) toCmd(user domain.Account) (cmd app.CmdToSendBindEmail, err error) {
	if cmd.Email, err = primitive.NewEmail(req.Email); err != nil {
		return
	}

	cmd.Capt = req.Capt

	cmd.User = user

	return
}

type tokenVerifyRequest struct {
	Token  string `json:"token" binding:"required"`
	Action string `json:"action" binding:"required"`
}

func (req *tokenVerifyRequest) ToCmd() (string, primitive.TokenPerm, error) {
	perm, err := primitive.NewTokenPerm(req.Action)
	if err != nil {
		return "", nil, fmt.Errorf("invalid action: %w", err)
	}

	return req.Token, perm, nil
}

type tokenVerifyResp struct {
	Account string `json:"account"`
}

type revokePrivacyInfo struct {
	IdToken string `json:"id_token"`
}
