/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package domain provides domain models and functionality for managing space apps.
package domain

import (
	"fmt"

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
		e := fmt.Errorf("not init")
		return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, e.Error(), e)
	}

	app.Status = appprimitive.AppStatusBuilding
	app.BuildLogURL = logURL

	return nil
}

// SetBuildIsDone sets the build status of the space app based on the success parameter.
func (app *SpaceApp) SetBuildIsDone(success bool) error {
	if !app.Status.IsBuilding() {
		e := fmt.Errorf("not building")
		return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, e.Error(), e)
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
		e := fmt.Errorf("not build successful or restarting")
		return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, e.Error(), e)
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
			return allerror.New(allerror.ErrorCodeSpaceAppRestartOverTime, "not over time to restart",
				fmt.Errorf("restart cost(%d) not over time(%d) to restart", now-oldRestartTime, config.RestartOverTime),
			)
		}
		app.RestartedAt = now
		return nil
	}
	if !app.Status.IsReadyToRestart() {
		e := fmt.Errorf("not ready to restart")
		return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, e.Error(), e)
	}
	app.Status = appprimitive.AppStatusRestarted
	app.RestartedAt = now
	return nil
}
