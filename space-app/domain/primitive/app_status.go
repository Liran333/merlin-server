package primitive

import (
	"errors"
	"strings"
)

const (
	Init = "init"

	Building          = "building"
	BuildFailed       = "build_failed"
	BuildSuccessfully = "build_successfully"

	Serving     = "serving"
	StartFailed = "start_failed"
)

var (
	AppStatusInit              = appStatus(Init)
	AppStatusServing           = appStatus(Serving)
	AppStatusBuilding          = appStatus(Building)
	AppStatusBuildFailed       = appStatus(BuildFailed)
	AppStatusStartFailed       = appStatus(StartFailed)
	AppStatusBuildSuccessfully = appStatus(BuildSuccessfully)
)

// AppStatus
type AppStatus interface {
	IsInit() bool
	AppStatus() string
}

// NewAppStatus
func NewAppStatus(v string) (AppStatus, error) {
	v = strings.ToLower(v)

	switch v {
	case Init:
	case Serving:
	case Building:
	case BuildFailed:
	case StartFailed:
	case BuildSuccessfully:

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
	return string(r) == Init
}
