package app

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/space-app/domain"
)

type CmdToCreateApp = domain.SpaceAppIndex

type CmdToNotifyBuildIsStarted struct {
	domain.SpaceAppIndex

	LogURL primitive.URL
}

type CmdToNotifyServiceIsStarted struct {
	CmdToNotifyBuildIsStarted

	AppURL primitive.URL
}
