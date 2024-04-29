/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides the controllers for handling HTTP requests and managing the application's business logic.
package controller

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/spaceapp/app"
	appprimitive "github.com/openmerlin/merlin-server/spaceapp/domain/primitive"
)

// reqToCreateSpaceApp
type reqToCreateSpaceApp struct {
	SpaceId  string `json:"space_id"`
	CommitId string `json:"commit_id"`
}

func (req *reqToCreateSpaceApp) toCmd() (cmd app.CmdToCreateApp, err error) {
	cmd.SpaceId, err = primitive.NewIdentity(req.SpaceId)
	if err != nil {
		return
	}

	cmd.CommitId = req.CommitId

	return
}

// reqToUpdateBuildInfo
type reqToUpdateBuildInfo struct {
	reqToCreateSpaceApp

	LogURL string `json:"log_url"`
}

func (req *reqToUpdateBuildInfo) toCmd() (cmd app.CmdToNotifyBuildIsStarted, err error) {
	if cmd.SpaceAppIndex, err = req.reqToCreateSpaceApp.toCmd(); err != nil {
		return
	}

	cmd.LogURL, err = primitive.NewURL(req.LogURL)

	return
}

// reqToUpdateServiceInfo
type reqToUpdateServiceInfo struct {
	reqToUpdateBuildInfo

	AppURL string `json:"app_url"`
}

func (req *reqToUpdateServiceInfo) toCmd() (cmd app.CmdToNotifyServiceIsStarted, err error) {
	if cmd.CmdToNotifyBuildIsStarted, err = req.reqToUpdateBuildInfo.toCmd(); err != nil {
		return
	}

	cmd.AppURL, err = appprimitive.NewAppURL(req.AppURL)

	return
}

// reqToPauseSpaceApp
type reqToPauseSpaceApp struct {
	SpaceId string `json:"space_id"`
	IsForce bool   `json:"is_force"`
}

func (req *reqToPauseSpaceApp) toCmd() (cmd app.CmdToPauseSpaceApp, err error) {
	cmd.SpaceId, err = primitive.NewIdentity(req.SpaceId)
	if err != nil {
		return
	}

	cmd.IsForce = req.IsForce
	return
}

// reqToFailedStatus
type reqToFailedStatus struct {
	reqToCreateSpaceApp

	Status string `json:"status"`
	Reason string `json:"reason"`
}

func (req *reqToFailedStatus) toCmd() (cmd app.CmdToNotifyFailedStatus, err error) {
	if cmd.SpaceAppIndex, err = req.reqToCreateSpaceApp.toCmd(); err != nil {
		return
	}

	cmd.Reason = req.Reason

	cmd.Status, err = appprimitive.NewAppStatus(req.Status)

	return
}
