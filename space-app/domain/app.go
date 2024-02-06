package domain

import (
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	appprimitive "github.com/openmerlin/merlin-server/space-app/domain/primitive"
)

type SpaceAppIndex struct {
	SpaceId  primitive.Identity
	CommitId string
}

type SpaceApp struct {
	Id int64
	SpaceAppIndex

	Status appprimitive.AppStatus

	AppURL    primitive.URL
	AppLogURL primitive.URL

	AllBuildLog string
	BuildLogURL primitive.URL

	Version int
}

func (app *SpaceApp) StartBuilding(logURL primitive.URL) error {
	if !app.Status.IsInit() {
		return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, "not init")
	}

	app.Status = appprimitive.AppStatusBuilding
	app.BuildLogURL = logURL

	return nil
}
