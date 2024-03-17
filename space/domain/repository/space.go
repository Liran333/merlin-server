/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package repository provides interfaces for interacting with space-related data.
package repository

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/space/domain"
)

// SpaceSummary represents a summary of a space.
type SpaceSummary struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Desc          string `json:"desc"`
	Task          string `json:"task"`
	Owner         string `json:"owner"`
	Fullname      string `json:"fullname"`
	UpdatedAt     int64  `json:"updated_at"`
	LikeCount     int    `json:"like_count"`
	DownloadCount int    `json:"download_count"`
}

// ListOption contains options for listing spaces.
type ListOption struct {
	// can't define Name as domain.ResourceName
	// because the Name can be subpart of the real resource name
	Name string

	// list the space of Owner
	Owner primitive.Account

	// list by visibility
	Visibility primitive.Visibility

	// list space which have one of licenses
	License primitive.License

	// list space which have at least one label for each kind of lables.
	Labels domain.SpaceLabels

	// sort
	SortType primitive.SortType

	// whether to calculate the total
	Count        bool
	PageNum      int
	CountPerPage int
}

// Pagination returns a boolean indicating whether pagination is enabled and the offset for pagination.
func (opt *ListOption) Pagination() (bool, int) {
	if opt.PageNum > 0 && opt.CountPerPage > 0 {
		return true, (opt.PageNum - 1) * opt.CountPerPage
	}

	return false, 0
}

// SpaceRepositoryAdapter is an interface for interacting with space repositories.
type SpaceRepositoryAdapter interface {
	Add(*domain.Space) error
	FindByName(*domain.SpaceIndex) (domain.Space, error)
	FindById(primitive.Identity) (domain.Space, error)
	Delete(primitive.Identity) error
	Save(*domain.Space) error
	List(*ListOption) ([]SpaceSummary, int, error)
	Count(*ListOption) (int, error)
	SearchSpace(*ListOption, primitive.Account) ([]SpaceSummary, int, error)
}

// SpaceLabelsRepoAdapter is an interface for interacting with space label repositories.
type SpaceLabelsRepoAdapter interface {
	Save(*domain.SpaceIndex, *domain.SpaceLabels) error
}
