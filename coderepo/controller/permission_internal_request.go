/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides the controllers for handling restful requests and converting them into commands
package controller

import (
	"github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

type reqToCheckPermission struct {
	User  string `json:"user"`
	Name  string `json:"name"`
	Owner string `json:"owner"`
}

func (req *reqToCheckPermission) toCmd() (
	user primitive.Account, index domain.CodeRepoIndex, err error,
) {
	if req.User != "" {
		if user, err = primitive.NewAccount(req.User); err != nil {
			return
		}
	}

	if index.Owner, err = primitive.NewAccount(req.Owner); err != nil {
		return
	}

	index.Name, err = primitive.NewMSDName(req.Name)

	return
}
