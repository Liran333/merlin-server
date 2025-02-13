/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package repository provides an domain for the repository
package repository

import (
	"github.com/openmerlin/merlin-server/activity/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

// ActivityInternalAdapter represents an interface for managing activities.
type ActivityInternalAdapter interface {
	Save(activities *domain.Activity) error
	DeleteAll(activities *domain.Activity) error
}

// ActivitiesRepositoryAdapter represents an interface for managing activities.
type ActivitiesRepositoryAdapter interface {
	List([]primitive.Account, *ListOption) ([]domain.Activity, int, error)
	Save(activities *domain.Activity) error
	Delete(activities *domain.Activity) error
	HasLike(primitive.Account, primitive.Identity) (bool, error)
}

// ListOption represents options for listing activities.
type ListOption struct {
	Name    []string
	Space   string
	Model   string
	Dataset string
	Like    string
	// sort
	SortType primitive.SortType

	// whether to calculate the total
	Count        bool
	PageNum      int
	CountPerPage int
}

// Pagination calculates the offset for pagination.
func (opt *ListOption) Pagination() (bool, int) {
	if opt.PageNum > 0 && opt.CountPerPage > 0 {
		return true, (opt.PageNum - 1) * opt.CountPerPage
	}

	return false, 0
}
