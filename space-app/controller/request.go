package controller

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/space-app/app"
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
