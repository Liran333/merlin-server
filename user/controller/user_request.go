/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package controller

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/user/app"
	"github.com/openmerlin/merlin-server/user/domain"
	"github.com/openmerlin/merlin-server/utils"
	"golang.org/x/xerrors"
)

type userBasicInfoUpdateRequest struct {
	AvatarId     string  `json:"avatar_id"`
	Desc         *string `json:"description"`
	Fullname     *string `json:"fullname"`
	RevokeDelete *bool   `json:"revoke_delete"`
}

func (req *userBasicInfoUpdateRequest) toCmd() (
	cmd app.UpdateUserBasicInfoCmd,
	err error,
) {
	if req.Desc != nil {
		if cmd.Desc, err = primitive.NewAccountDesc(*req.Desc); err != nil {
			return
		}
	}

	if req.Fullname != nil {
		if cmd.Fullname, err = primitive.NewAccountFullname(*req.Fullname); err != nil {
			return
		}
	}

	if req.RevokeDelete != nil {
		cmd.RevokeDelete = true
	}

	if req.Desc == nil && req.Fullname == nil && req.RevokeDelete == nil {
		err = fmt.Errorf("all param are empty")
		return
	}

	if req.AvatarId != "" {
		if cmd.AvatarId, err = primitive.NewAvatar(req.AvatarId); err != nil {
			return
		}
	}

	return
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

type bindEmailRequest struct {
	Email    string `json:"email" binding:"required"`
	PassCode string `json:"code" binding:"required"`
}

func (req *bindEmailRequest) action() string {
	return fmt.Sprintf("bind email %s", utils.AnonymizeEmail(req.Email))
}

func (req *bindEmailRequest) toCmd(user domain.Account) (cmd app.CmdToVerifyBindEmail, err error) {
	if cmd.Email, err = primitive.NewUserEmail(req.Email); err != nil {
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
	return fmt.Sprintf("send email verify code to %s", utils.AnonymizeEmail(req.Email))
}

func (req *sendEmailRequest) toCmd(user domain.Account) (cmd app.CmdToSendBindEmail, err error) {
	if cmd.Email, err = primitive.NewUserEmail(req.Email); err != nil {
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
		return "", nil, fmt.Errorf("invalid request: %w", err)
	}

	return req.Token, perm, nil
}

type tokenVerifyResp struct {
	Account string `json:"account"`
}

type revokePrivacyInfo struct {
	IdToken string `json:"id_token"`
}

func toUploadAvatarCmd(ctx *gin.Context, user primitive.Account) (app.CmdToUploadAvatar, error) {
	f, err := ctx.FormFile("file")
	if err != nil {
		return app.CmdToUploadAvatar{}, xerrors.Errorf("failed to parse request param: %w", err)
	}

	if f.Size > config.MaxAvatarFileSize {
		err = errors.New("file too big")

		return app.CmdToUploadAvatar{}, err
	}

	name, err := primitive.NewFileName(f.Filename)
	if err != nil {
		return app.CmdToUploadAvatar{}, err
	}

	if !name.IsPictureName() {
		err = errors.New("file format error")

		return app.CmdToUploadAvatar{}, err
	}

	p, err := f.Open()
	if err != nil {
		return app.CmdToUploadAvatar{}, xerrors.Errorf("can not get file: %w", err)
	}

	defer p.Close()

	cmd := app.CmdToUploadAvatar{
		User:  user,
		Image: p,
		FileName: commonctl.GetSaltHash(fmt.Sprintf("%s%s%v",
			user.Account(), f.Filename, utils.Now())) + name.GetFormat(),
	}

	return cmd, nil
}
