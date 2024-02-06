package app

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/space-app/domain"
	appprimitive "github.com/openmerlin/merlin-server/space-app/domain/primitive"
)

type CmdToCreateApp struct {
	SpaceId  primitive.Identity
	CommitId string
}

func (cmd *CmdToCreateApp) toApp() domain.SpaceApp {
	return domain.SpaceApp{
		Status:   appprimitive.AppStatusInit,
		SpaceId:  cmd.SpaceId,
		CommitId: cmd.CommitId,
	}
}
