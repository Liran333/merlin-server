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
	AppStatusInit              = appStatus(appInit)
	AppStatusServing           = appStatus(serving)
	AppStatusBuilding          = appStatus(building)
	AppStatusBuildFailed       = appStatus(buildFailed)
	AppStatusStartFailed       = appStatus(startFailed)
	AppStatusBuildSuccessfully = appStatus(buildSuccessfully)
)

// AppStatus
type AppStatus interface {
	IsInit() bool
	AppStatus() string
	IsBuilding() bool
	IsBuildSuccessful() bool
}

// NewAppStatus
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

// CreateAppStatus
func CreateAppStatus(v string) AppStatus {
	return appStatus(v)
}

// appStatus
type appStatus string

func (r appStatus) AppStatus() string {
	return string(r)
}

func (r appStatus) IsInit() bool {
	return string(r) == appInit
}

func (r appStatus) IsBuilding() bool {
	return string(r) == building
}

func (r appStatus) IsBuildSuccessful() bool {
	return string(r) == buildSuccessfully
}
