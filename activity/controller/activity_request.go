/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package controller provides functionality for managing the application's controllers.
package controller

import (
	"github.com/openmerlin/merlin-server/activity/app"
)

// nolint:golint,unused
const (
	firstPage = 1
	typeModel = "model"
)

// activityInfo
type activityInfo struct {
	AvatarId  string `json:"avatar_id"`
	OwnerType int    `json:"owner_type"`
	*app.ActivitySummaryDTO
}

// activitiesInfo
type activitiesInfo struct {
	Total      int            `json:"total"`
	Activities []activityInfo `json:"activities"`
}
