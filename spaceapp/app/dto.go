package app

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/spaceapp/domain"
)

type CmdToCreateApp = domain.SpaceAppIndex

type CmdToNotifyBuildIsStarted struct {
	domain.SpaceAppIndex

	LogURL primitive.URL
}

type CmdToNotifyBuildIsDone struct {
	domain.SpaceAppIndex

	Success bool
}

type CmdToNotifyServiceIsStarted struct {
	CmdToNotifyBuildIsStarted

	AppURL primitive.URL
}

type SpaceAppDTO struct {
	Id          int64  `json:"id"`
	Status      string `json:"status"`
	AppURL      string `json:"app_url"`
	AppLogURL   string `json:"app_log_url"`
	BuildLogURL string `json:"build_log_url"`
}

func toSpaceAppDTO(app *domain.SpaceApp) SpaceAppDTO {
	dto := SpaceAppDTO{
		Id:     app.Id,
		Status: app.Status.AppStatus(),
	}

	if app.AppURL != nil {
		dto.AppURL = app.AppURL.URL()
	}

	if app.AppLogURL != nil {
		dto.AppLogURL = app.AppLogURL.URL()
	}

	if app.BuildLogURL != nil {
		dto.BuildLogURL = app.BuildLogURL.URL()
	}

	return dto
}
