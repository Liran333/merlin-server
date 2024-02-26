/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package controller

import (
	"github.com/gin-gonic/gin"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/session/app"
)

// reqToLogin
type reqToLogin struct {
	Code        string `json:"code"          required:"true"`
	RedirectURI string `json:"redirect_uri"  required:"true"`
}

func (req *reqToLogin) toCmd(ctx *gin.Context) (cmd app.CmdToLogin, err error) {
	cmd.Code = req.Code
	cmd.RedirectURI = req.RedirectURI

	if cmd.IP, err = commonctl.GetIp(ctx); err != nil {
		return
	}

	cmd.UserAgent, err = commonctl.GetUserAgent(ctx)

	return
}

type logoutInfo struct {
	IdToken string `json:"id_token"`
}
