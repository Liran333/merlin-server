/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// The controller package provides functionality for handling requests and
// processing business logic related to user organizations.
package controller

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/computility/app"
)

type reqToUserOrgOperate struct {
	UserName string `json:"user_name"        required:"true"`
	OrgName  string `json:"org_name"         required:"true"`
}

func (req *reqToUserOrgOperate) toCmd() (cmd app.CmdToUserOrgOperate, err error) {
	if cmd.UserName, err = primitive.NewAccount(req.UserName); err != nil {
		return
	}

	if cmd.OrgName, err = primitive.NewAccount(req.OrgName); err != nil {
		return
	}

	return
}

type reqToOrgDelete struct {
	OrgName string `json:"org_name"           required:"true"`
}

func (req *reqToOrgDelete) toCmd() (cmd app.CmdToOrgDelete, err error) {
	if cmd.OrgName, err = primitive.NewAccount(req.OrgName); err != nil {
		return
	}

	return
}
