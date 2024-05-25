/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package domain provides an domain for the repository
package domain

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

// ActivityType enum
type ActivityType string

// Activity const
const (
	Create ActivityType = "create"
	Update ActivityType = "update"
	Like   ActivityType = "like"
)

// Activity struct represents the user activity entity.
type Activity struct {
	Type     ActivityType
	Time     int64
	Name     primitive.MSDName
	Owner    primitive.Account
	Resource Resource
}

// ActivitySummary struct represents the user activity entity with statistic.
type ActivitySummary struct {
	Activity
	Stat
}

// Stat struct represents the statistic of an activity.
type Stat struct {
	LikeCount     int `json:"like_count"`
	DownloadCount int `json:"download_count"`
}

// Resource struct represents the resource object targeted by user activities.
type Resource struct {
	Type    primitive.ObjType  // Resource type
	Index   primitive.Identity // Resource index
	Owner   primitive.Account
	Disable bool
}
