/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package domain provides domain models and functionality for managing space apps.
package domain

import (
	"fmt"

	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	computilityapp "github.com/openmerlin/merlin-server/computility/app"
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
	ResumedAt   int64
	RestartedAt int64

	AppURL    appprimitive.AppURL
	AppLogURL primitive.URL

	AllBuildLog string
	BuildLogURL primitive.URL

	Version int
}

type PauseSpaceAppHandle func(bool, computilityapp.ComputilityAppService) error

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
func (app *SpaceApp) StartService(appURL appprimitive.AppURL, logURL primitive.URL) error {
	if !app.Status.IsBuildSuccessful() &&
		!app.Status.IsRestarting() &&
		!app.Status.IsResuming() {
		e := fmt.Errorf("spaceId:%s, not build successful "+
			"or restarting "+
			"or resuming", app.SpaceId.Identity())
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
func (app *SpaceApp) RestartService() error {
	now := utils.Now()
	if app.Status.IsRestarting() {
		if now-app.RestartedAt < config.RestartOverTime {
			return allerror.New(allerror.ErrorCodeSpaceAppRestartOverTime, "not over time to restart",
				fmt.Errorf("restart cost(%d) not over time(%d) to restart", now-app.RestartedAt, config.RestartOverTime),
			)
		}
		app.RestartedAt = now
		return nil
	}
	if !app.Status.IsServing() && !app.Status.IsStartFailed() {
		e := fmt.Errorf("not ready to restart")
		return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, e.Error(), e)
	}
	app.Status = appprimitive.AppStatusRestarted
	app.RestartedAt = now
	return nil
}

// IsAppStatusAllow get status can be update.
func (app *SpaceApp) IsAppStatusAllow(status appprimitive.AppStatus) error {
	if !status.IsUpdateStatusAccept() {
		e := fmt.Errorf("app status not accept")
		return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, e.Error(), e)
	}

	app.Status = status

	return nil
}

// PauseService pause the service for space app
func (app *SpaceApp) PauseService(
	isForce bool,
	compPowerAllocated bool,
	compUtility computilityapp.ComputilityAppService,
	pauseHook PauseSpaceAppHandle,
	) error {
	if !app.Status.IsServing() && !isForce {
		e := fmt.Errorf("spaceId:%s, not serving", app.SpaceId.Identity())
		return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, e.Error(), e)
	}

	app.Status = appprimitive.AppStatusPaused

	return pauseHook(compPowerAllocated, compUtility)
}

// ResumeService resume the service for space app with oldResumeTime
func (app *SpaceApp) ResumeService(
	isNpu bool,
	compUtility computilityapp.ComputilityAppService,
	pauseHook PauseSpaceAppHandle,
	) error {
	now := utils.Now()
	if app.Status.IsResuming() {
		if now-app.ResumedAt < config.ResumeOverTime {
			return allerror.New(allerror.ErrorCodeSpaceAppResumeOverTime, "not over time to resume",
				fmt.Errorf("resume cost(%d) not over time(%d) to resume", now-app.ResumedAt, config.ResumeOverTime),
			)
		}
		app.ResumedAt = now
		return nil
	}
	if !app.Status.IsPaused() && !app.Status.IsResumeFailed() {
		e := fmt.Errorf("spaceId:%s, not ready to resume", app.SpaceId.Identity())
		return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, e.Error(), e)
	}
	app.Status = appprimitive.AppStatusResuming
	app.ResumedAt = now

	return pauseHook(isNpu, compUtility)
}
