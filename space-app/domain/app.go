package domain

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	appprimitive "github.com/openmerlin/merlin-server/space-app/domain/primitive"
)

type SpaceApp struct {
	SpaceId  primitive.Identity
	CommitId string

	Status appprimitive.AppStatus

	AppURL    primitive.URL
	AppLogURL primitive.URL

	BuilgLog    string
	BuildLogURL primitive.URL
}
