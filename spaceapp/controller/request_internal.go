/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

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

// reqToSetBuildIsDone
type reqToSetBuildIsDone struct {
	reqToCreateSpaceApp

	Logs    string `json:"logs"`
	Success bool   `json:"success"`
}

func (req *reqToSetBuildIsDone) toCmd() (cmd app.CmdToNotifyBuildIsDone, err error) {
	if cmd.SpaceAppIndex, err = req.reqToCreateSpaceApp.toCmd(); err != nil {
		return
	}

	cmd.Logs = req.Logs
	cmd.Success = req.Success

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

	cmd.AppURL, err = primitive.NewURL(req.AppURL)

	return
}

// reqToSetStatus
type reqToSetStatus struct {
	reqToCreateSpaceApp

	Status string `json:"status"`
}

func (req *reqToSetStatus) toCmd() (cmd app.CmdToNotifyUpdateStatus, err error) {
	if cmd.SpaceAppIndex, err = req.reqToCreateSpaceApp.toCmd(); err != nil {
		return
	}

	cmd.Status, err = appprimitive.NewAppStatus(req.Status)
	return
}
