package controller

import (
	"github.com/openmerlin/merlin-server/activity/domain"
)

// nolint:golint,unused
const (
	firstPage = 1
	typeModel = "model"
)

// activityInfo
type activityInfo struct {
	AvatarId string `json:"avatar_id"`
	*domain.Activity
}

// activitiesInfo
type activitiesInfo struct {
	Total      int            `json:"total"`
	Activities []activityInfo `json:"activities"`
}
