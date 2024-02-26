/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive provides primitive types and functionality for managing application statuses.
package primitive

import (
	"errors"
	"strings"
)

const (
	appInit = "init"

	building          = "building"
	buildFailed       = "build_failed"
	buildSuccessfully = "build_successfully"

	serving     = "serving"
	startFailed = "start_failed"
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

	// AppStatusBuildSuccessfully represents the application status when the build process is successful.
	AppStatusBuildSuccessfully = appStatus(buildSuccessfully)
)

// AppStatus is an interface that defines methods for working with application statuses.
type AppStatus interface {
	IsInit() bool
	AppStatus() string
	IsBuilding() bool
	IsBuildSuccessful() bool
}

// NewAppStatus creates a new instance of AppStatus based on the provided value.
func NewAppStatus(v string) (AppStatus, error) {
	v = strings.ToLower(v)

	switch v {
	case appInit:
	case serving:
	case building:
	case buildFailed:
	case startFailed:
	case buildSuccessfully:

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
	return string(r) == buildSuccessfully
}
