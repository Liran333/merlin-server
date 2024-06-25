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
	Id primitive.Identity

	SpaceAppIndex

	Status appprimitive.AppStatus
	Reason string

	ResumedAt   int64
	RestartedAt int64

	AppURL      appprimitive.AppURL
	AppLogURL   primitive.URL
	BuildLogURL primitive.URL

	Version int
}

// StartBuilding starts the building process for the space app and sets the build log URL.
func (app *SpaceApp) StartBuilding(logURL primitive.URL) error {
	if !app.Status.IsInit() {
		e := fmt.Errorf("old status is %s, can not set", app.Status.AppStatus())
		return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, e.Error(), e)
	}

	app.Status = appprimitive.AppStatusBuilding
	app.BuildLogURL = logURL

	return nil
}

// SetBuildFailed set app status is build failed.
func (app *SpaceApp) SetBuildFailed(status appprimitive.AppStatus, reason string) error {
	if !app.Status.IsBuilding() {
		e := fmt.Errorf("old status is %s, can not set", app.Status.AppStatus())
		return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, e.Error(), e)
	}

	app.Status = status
	app.Reason = reason

	return nil
}

// SetStarting sets the starting status of the space app based on the success parameter.
func (app *SpaceApp) SetStarting() error {
	if !app.Status.IsBuilding() {
		e := fmt.Errorf("old status is %s, can not set", app.Status.AppStatus())
		return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, e.Error(), e)
	}

	app.Status = appprimitive.AppStatusServeStarting

	return nil
}

// SetStartFailed set app status is start failed.
func (app *SpaceApp) SetStartFailed(status appprimitive.AppStatus, reason string) error {
	if !app.Status.IsStarting() {
		e := fmt.Errorf("old status is %s, can not set", app.Status.AppStatus())
		return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, e.Error(), e)
	}

	app.Status = status
	app.Reason = reason

	return nil
}

// StartServing starts the service for the space app with the specified app URL and log URL.
func (app *SpaceApp) StartServing(appURL appprimitive.AppURL, logURL primitive.URL) error {
	if app.Status.IsStarting() || app.Status.IsRestarting() || app.Status.IsResuming() {

		app.Status = appprimitive.AppStatusServing
		app.AppURL = appURL
		app.AppLogURL = logURL

		return nil
	}

	e := fmt.Errorf("old status is %s, can not set", app.Status.AppStatus())
	return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, e.Error(), e)
}

// SetRestartFailed set app status is restart failed.
func (app *SpaceApp) SetRestartFailed(status appprimitive.AppStatus, reason string) error {
	if !app.Status.IsRestarting() {
		e := fmt.Errorf("old status is %s, can not set", app.Status.AppStatus())
		return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, e.Error(), e)
	}

	app.Status = status
	app.Reason = reason

	return nil
}

// SetResumeFailed set app status is restart failed.
func (app *SpaceApp) SetResumeFailed(status appprimitive.AppStatus, reason string) error {
	if !app.Status.IsResuming() {
		e := fmt.Errorf("old status is %s, can not set", app.Status.AppStatus())
		return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, e.Error(), e)
	}

	app.Status = status
	app.Reason = reason

	return nil
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
	if app.Status.IsServing() || app.Status.IsStartFailed() || app.Status.IsRestartFailed() {
		app.Status = appprimitive.AppStatusRestarted
		app.RestartedAt = now
		return nil
	}
	e := fmt.Errorf("old status not %s, can not set", app.Status.AppStatus())
	return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, e.Error(), e)
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
func (app *SpaceApp) PauseService() error {
	if app.Status.IsPaused() {
		e := fmt.Errorf("spaceId:%s, is paused", app.SpaceId.Identity())
		return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, e.Error(), e)
	}

	if !app.Status.IsServing() {
		e := fmt.Errorf("spaceId:%s, not serving", app.SpaceId.Identity())
		return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, e.Error(), e)
	}

	app.Status = appprimitive.AppStatusPaused

	return nil
}

// ResumeService resume the service for space app with oldResumeTime
func (app *SpaceApp) ResumeService() error {
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

	if app.Status.IsPaused() || app.Status.IsResumeFailed() {
		app.Status = appprimitive.AppStatusResuming
		app.ResumedAt = now
		return nil
	}

	e := fmt.Errorf("spaceId:%s, not ready to resume", app.SpaceId.Identity())
	return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, e.Error(), e)
}

// SleepService sleep the service for space app
func (app *SpaceApp) SleepService() error {
	if !app.Status.IsServing() {
		e := fmt.Errorf("spaceId:%s, not serving", app.SpaceId.Identity())
		return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, e.Error(), e)
	}

	app.Status = appprimitive.AppStatusSleeping

	return nil
}

// WakeupService wakeup the service for space app
func (app *SpaceApp) WakeupService() error {
	if !app.Status.IsSleeping() {
		e := fmt.Errorf("spaceId:%s, is not sleeping", app.SpaceId.Identity())
		return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, e.Error(), e)
	}

	app.Status = appprimitive.AppStatusServeStarting

	return nil
}

// IsAppNotAllowToInit app can be init if return false
func (app *SpaceApp) IsAppNotAllowToInit() bool {
	if app.Status.IsPaused() || app.Status.IsResuming() || app.Status.IsResumeFailed() {
		return true
	}

	return false
}

// GetFailedReason app only return failed reason
func (app *SpaceApp) GetFailedReason() string {
	if !app.Status.IsUpdateStatusAccept() {
		return ""
	}
	return app.Reason
}

// SpaceAppBuildLog is the value object of log
type SpaceAppBuildLog struct {
	AppId primitive.Identity
	Logs  string
}
