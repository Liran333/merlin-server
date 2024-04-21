/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive provides primitive types and functionality for managing application statuses.
package primitive

import (
	"errors"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	appInvalid = "invalid"

	appInit = "init"

	building    = "building"
	buildFailed = "build_failed"

	starting    = "starting"
	serving     = "serving"
	startFailed = "start_failed"

	restarting = "restarting"

	paused       = "paused"
	resuming     = "resuming"
	resumeFailed = "resume_failed"
)

var (
	// AppStatusInit represents the application status when it is in the initialization phase.
	AppStatusInit = appStatus(appInit)

	// AppStatusServing represents the application status when it is serving requests.
	AppStatusServing = appStatus(serving)

	// AppStatusBuilding represents the application status when it is being built.
	AppStatusBuilding = appStatus(building)

	// AppStatusBuildFailed represents the application status when the build process fails.
	AppStatusBuildFailed = appStatus(buildFailed)

	// AppStatusStartFailed represents the application status when the start process fails.
	AppStatusStartFailed = appStatus(startFailed)

	// AppStatusServeStarting represents the application status when the build process is successful.
	AppStatusServeStarting = appStatus(starting)

	// AppStatusRestarted represents the application status when the app is restarted.
	AppStatusRestarted = appStatus(restarting)

	// AppStatusPaused represents the application status when the app is pause.
	AppStatusPaused = appStatus(paused)

	// AppStatusResuming represents the application status when the app is resume.
	AppStatusResuming = appStatus(resuming)
)

var acceptAppStatusSets = sets.NewString(appInvalid, resumeFailed)

// AppStatus is an interface that defines methods for working with application statuses.
type AppStatus interface {
	IsInit() bool
	AppStatus() string
	IsBuilding() bool
	IsBuildSuccessful() bool
	IsRestarting() bool
	IsPaused() bool
	IsResuming() bool
	IsResumeFailed() bool
	IsUpdateStatusAccept() bool
	IsStartFailed() bool
	IsServing() bool
}

// NewAppStatus creates a new instance of AppStatus based on the provided value.
func NewAppStatus(v string) (AppStatus, error) {
	v = strings.ToLower(v)

	switch v {
	case appInvalid:
	case appInit:
	case serving:
	case building:
	case buildFailed:
	case startFailed:
	case starting:
	case resumeFailed:

	default:
		return nil, errors.New("unknown appStatus")
	}

	return appStatus(v), nil
}

// CreateAppStatus creates a new instance of AppStatus with the provided value.
func CreateAppStatus(v string) AppStatus {
	return appStatus(v)
}

// appStatus
type appStatus string

// AppStatus returns the string representation of the appStatus.
func (r appStatus) AppStatus() string {
	return string(r)
}

// IsInit checks if the appStatus is equal to appInit.
func (r appStatus) IsInit() bool {
	return string(r) == appInit
}

// IsBuilding checks if the appStatus is equal to building.
func (r appStatus) IsBuilding() bool {
	return string(r) == building
}

// IsBuildSuccessful checks if the appStatus is equal to buildSuccessfully.
func (r appStatus) IsBuildSuccessful() bool {
	return string(r) == starting
}

// IsRestarting checks if the appStatus is equal to restarting.
func (r appStatus) IsRestarting() bool {
	return string(r) == restarting
}

// IsPaused checks if the appStatus is equal to paused.
func (r appStatus) IsPaused() bool {
	return string(r) == paused
}

// IsResuming checks if the appStatus is equal to resuming.
func (r appStatus) IsResuming() bool {
	return string(r) == resuming
}

// IsResuming checks if the appStatus is equal to resuming.
func (r appStatus) IsResumeFailed() bool {
	return string(r) == resumeFailed
}

// IsStartFailed checks if the appStatus startFailed
func (r appStatus) IsStartFailed() bool {
	return string(r) == startFailed
}

// IsServing checks if the appStatus serving
func (r appStatus) IsServing() bool {
	return string(r) == serving
}

// checks if the appStatus can be update
func (r appStatus) IsUpdateStatusAccept() bool {
	return acceptAppStatusSets.Has(string(r))
}
