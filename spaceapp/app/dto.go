/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	spacedomain "github.com/openmerlin/merlin-server/space/domain"
	"github.com/openmerlin/merlin-server/spaceapp/domain"
	appprimitive "github.com/openmerlin/merlin-server/spaceapp/domain/primitive"
)

// CmdToCreateApp is a command to create an app.
type CmdToCreateApp = domain.SpaceAppIndex

// CmdToNotifyBuildIsStarted is a command to notify that the build has started.
type CmdToNotifyBuildIsStarted struct {
	domain.SpaceAppIndex

	LogURL primitive.URL
}

// CmdToNotifyBuildIsDone is a command to notify that the build has finished.
type CmdToNotifyBuildIsDone struct {
	domain.SpaceAppIndex

	Logs    string
	Success bool
}

// CmdToNotifyServiceIsStarted is a command to notify that the service has started.
type CmdToNotifyServiceIsStarted struct {
	CmdToNotifyBuildIsStarted

	AppURL appprimitive.AppURL
}

// SpaceAppDTO is a data transfer object for space app.
type SpaceAppDTO struct {
	Id          int64  `json:"id"`
	Status      string `json:"status"`
	Reason      string `json:"reason"`
	AppURL      string `json:"app_url"`
	AppLogURL   string `json:"-"`
	BuildLogURL string `json:"-"`
}

func toSpaceDTO(space *spacedomain.Space) SpaceAppDTO {
	dto := SpaceAppDTO{
		Id:     space.Id.Integer(),
		Status: space.Exception.Exception(),
		Reason: primitive.ExceptionMap[space.Exception.Exception()],
	}
	return dto
}

func toSpaceNoCompQuotaDTO(space *spacedomain.Space) SpaceAppDTO {
	dto := SpaceAppDTO{
		Id:     space.Id.Integer(),
		Status: primitive.NoCompQuotaException,
		Reason: primitive.ExceptionMap[primitive.NoCompQuotaException],
	}
	return dto
}

func toSpaceAppDTO(app *domain.SpaceApp) SpaceAppDTO {
	dto := SpaceAppDTO{
		Id:     app.Id,
		Status: app.Status.AppStatus(),
		Reason: app.GetFailedReason(),
	}

	if app.AppURL != nil {
		dto.AppURL = app.AppURL.AppURL()
	}

	if app.AppLogURL != nil {
		dto.AppLogURL = app.AppLogURL.URL()
	}

	if app.BuildLogURL != nil {
		dto.BuildLogURL = app.BuildLogURL.URL()
	}

	return dto
}

// CmdToNotifyUpdateStatus is a command to notify that status has update.
type CmdToNotifyUpdateStatus struct {
	domain.SpaceAppIndex

	Status appprimitive.AppStatus
}

// CmdToPauseSpaceApp is a command to pause space app
type CmdToPauseSpaceApp struct {
	SpaceId primitive.Identity
	IsForce bool
}

// CmdToNotifyFailedStatus is a command to notify that status has update.
type CmdToNotifyFailedStatus struct {
	domain.SpaceAppIndex

	Reason string
	Status appprimitive.AppStatus
}
