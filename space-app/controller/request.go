package controller

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/space-app/app"
)

type reqToCreateSpaceApp struct {
	SpaceId  string
	CommitId string
}

func (req *reqToCreateSpaceApp) toCmd() (cmd app.CmdToCreateApp, err error) {
	cmd.SpaceId, err = primitive.NewIdentity(req.SpaceId)
	if err != nil {
		return
	}

	cmd.CommitId = req.CommitId

	return
}
