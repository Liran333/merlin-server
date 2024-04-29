/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package app provides functionality for the application.
package app

import (
	"github.com/openmerlin/merlin-server/activity/domain"
	"github.com/openmerlin/merlin-server/activity/domain/repository"
)

// ActivitysDTO is a struct that represents a data transfer object for a list of activity.
type ActivitysDTO struct {
	Total      int                  `json:"total"`
	AvatarId   string               `json:"avatar_id"`
	Activities []ActivitySummaryDTO `json:"activity"`
}

// ActivitySummary struct represents the user activity entity with statistic.
type ActivitySummaryDTO struct {
	ActivityDTO
	StatDTO
}

// ActivityDTO struct represents the user activity entity.
type ActivityDTO struct {
	Type     string      `json:"Type"`
	Time     int64       `json:"Time"`
	Name     string      `json:"Name"`
	Owner    string      `json:"Owner"`
	Resource ResourceDTO `json:"Resource"`
}

// StatDTO struct represents the statistic of an activity.
type StatDTO struct {
	LikeCount     int `json:"like_count"`
	DownloadCount int `json:"download_count"`
}

func toStatDTO(stat *domain.Stat) StatDTO {
	dto := StatDTO{
		LikeCount:     stat.LikeCount,
		DownloadCount: stat.DownloadCount,
	}

	return dto
}

// ResourceDTO struct represents the resource object targeted by user activities.
type ResourceDTO struct {
	Type    string `json:"Type"`
	Index   int64  `json:"Index"`
	Owner   string `json:"Owner"`
	Disable bool   `json:"disable"`
}

func toResourceDTO(resource *domain.Resource) ResourceDTO {
	dto := ResourceDTO{
		Type:    string(resource.Type),
		Index:   resource.Index.Integer(),
		Owner:   resource.Owner.Account(),
		Disable: resource.Disable,
	}

	return dto
}

func toActivityDTO(activity *domain.Activity) ActivityDTO {
	dto := ActivityDTO{
		Type:     string(activity.Type),
		Time:     activity.Time,
		Name:     activity.Name.MSDName(),
		Owner:    activity.Owner.Account(),
		Resource: toResourceDTO(&activity.Resource),
	}

	return dto
}

func toActivitySummaryDTO(activity *domain.Activity, stat *domain.Stat) ActivitySummaryDTO {
	dto := ActivitySummaryDTO{
		ActivityDTO: toActivityDTO(activity),
		StatDTO:     toStatDTO(stat),
	}

	return dto
}

type CmdToAddActivity = domain.Activity

// CmdToListActivities is a type alias for repository.ListOption, representing a command to list models.
type CmdToListActivities = repository.ListOption
