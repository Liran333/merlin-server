/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package domain provides domain models and functionality for managing space apps.
package domain

import (
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	appprimitive "github.com/openmerlin/merlin-server/spaceapp/domain/primitive"
	"github.com/openmerlin/merlin-server/utils"
)

// SpaceAppIndex represents the index for a space app.
type SpaceAppIndex struct {
	SpaceId  primitive.Identity
	CommitId string
}

// SpaceApp represents a space app.
type SpaceApp struct {
	Id int64

	SpaceAppIndex

	Status      appprimitive.AppStatus
	RestartedAt int64

	AppURL    primitive.URL
	AppLogURL primitive.URL

	AllBuildLog string
	BuildLogURL primitive.URL

	Version int
}

// StartBuilding starts the building process for the space app and sets the build log URL.
func (app *SpaceApp) StartBuilding(logURL primitive.URL) error {
	if !app.Status.IsInit() {
		return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, "not init")
	}

	app.Status = appprimitive.AppStatusBuilding
	app.BuildLogURL = logURL

	return nil
}

// SetBuildIsDone sets the build status of the space app based on the success parameter.
func (app *SpaceApp) SetBuildIsDone(success bool) error {
	if !app.Status.IsBuilding() {
		return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, "not building")
	}

	if success {
		app.Status = appprimitive.AppStatusServeStarting
	} else {
		app.Status = appprimitive.AppStatusBuildFailed
	}

	return nil
}

// StartService starts the service for the space app with the specified app URL and log URL.
func (app *SpaceApp) StartService(appURL, logURL primitive.URL) error {
	if !app.Status.IsBuildSuccessful() && !app.Status.IsRestarting() {
		return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, "not build successful or restarting")
	}

	if appURL != nil {
		app.Status = appprimitive.AppStatusServing
		app.AppURL = appURL
		app.AppLogURL = logURL
	} else {
		app.Status = appprimitive.AppStatusStartFailed
	}

	return nil
}

type SpaceAppBuildLog struct {
	AppId int64
	Logs  string
}

// RestartService restart the service for space app with oldRestartTime
func (app *SpaceApp) RestartService(oldRestartTime int64) error {
	now := utils.Now()
	if app.Status.IsRestarting() {
		if now-oldRestartTime < config.RestartOverTime {
			return allerror.New(allerror.ErrorCodeSpaceAppRestartOverTime, "not over time to restart")
		}
		app.RestartedAt = now
		return nil
	}
	if !app.Status.IsReadyToRestart() {
		return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, "not ready to restart")
	}
	app.Status = appprimitive.AppStatusRestarted
	app.RestartedAt = now
	return nil
}
