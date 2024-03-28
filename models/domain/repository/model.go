/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package repository provides interfaces for interacting with models and model labels in the domain.
package repository

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/models/domain"
)

// ModelSummary represents a summary of a model.
type ModelSummary struct {
	Id            string   `json:"id"`
	Name          string   `json:"name"`
	Desc          string   `json:"desc"`
	Task          string   `json:"task"`
	Owner         string   `json:"owner"`
	License       string   `json:"license"`
	Fullname      string   `json:"fullname"`
	UpdatedAt     int64    `json:"updated_at"`
	LikeCount     int      `json:"like_count"`
	Frameworks    []string `json:"frameworks"`
	DownloadCount int      `json:"download_count"`
}

// ListOption represents options for listing models.
type ListOption struct {
	// can't define Name as domain.ResourceName
	// because the Name can be subpart of the real resource name
	Name string

	// list the models of Owner
	Owner primitive.Account

	// list by visibility
	Visibility primitive.Visibility

	// list models which have one of licenses
	License primitive.License

	// list models which have at least one label for each kind of lables.
	Labels domain.ModelLabels

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

// ModelRepositoryAdapter represents an interface for managing models.
type ModelRepositoryAdapter interface {
	Add(*domain.Model) error
	FindByName(*domain.ModelIndex) (domain.Model, error)
	FindById(primitive.Identity) (domain.Model, error)
	Delete(primitive.Identity) error
	Save(*domain.Model) error
	List(*ListOption) ([]ModelSummary, int, error)
	Count(*ListOption) (int, error)
	SearchModel(*ListOption, primitive.Account) ([]ModelSummary, int, error)
}

// ModelLabelsRepoAdapter represents an interface for managing model labels.
type ModelLabelsRepoAdapter interface {
	Save(primitive.Identity, *domain.ModelLabels) error
}
