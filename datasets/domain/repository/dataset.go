/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package repository provides interfaces for interacting with datasets and datasets labels in the domain.
package repository

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/datasets/domain"
	orgrepo "github.com/openmerlin/merlin-server/organization/domain/repository"
)

// DatasetSummary represents a summary of a dataset.
type DatasetSummary struct {
	Id            string   `json:"id"`
	Name          string   `json:"name"`
	Desc          string   `json:"desc"`
	Task          []string `json:"task"`
	Owner         string   `json:"owner"`
	License       string   `json:"license"`
	Fullname      string   `json:"fullname"`
	UpdatedAt     int64    `json:"updated_at"`
	LikeCount     int      `json:"like_count"`
	DownloadCount int      `json:"download_count"`
	Disable       bool     `json:"disable"`
	DisableReason string   `json:"disable_reason"`
	Size          string   `json:"size"`
	Language      []string `json:"language"`
	Domain        []string `json:"domain"`
}

// ListOption represents options for listing datasets.
type ListOption struct {
	// can't define Name as domain.ResourceName
	// because the Name can be subpart of the real resource name
	Name string

	// list the datasets of Owner
	Owner primitive.Account

	// list by visibility
	Visibility primitive.Visibility

	// list datasets which have one of licenses
	License primitive.License

	// list datasets which have at least one label for each kind of lables.
	Labels domain.DatasetLabels

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

// DatasetRepositoryAdapter represents an interface for managing datasets.
type DatasetRepositoryAdapter interface {
	Add(*domain.Dataset) error
	FindByName(*domain.DatasetIndex) (domain.Dataset, error)
	FindById(primitive.Identity) (domain.Dataset, error)
	Delete(primitive.Identity) error
	Save(*domain.Dataset) error
	List(*ListOption, primitive.Account, orgrepo.OrgMember) ([]DatasetSummary, int, error)
	Count(*ListOption) (int, error)
	SearchDataset(*ListOption, primitive.Account, orgrepo.OrgMember) ([]DatasetSummary, int, error)
	AddLike(domain.Dataset) error
	DeleteLike(domain.Dataset) error
}

// DatasetLabelsRepoAdapter represents an interface for managing dataset labels.
type DatasetLabelsRepoAdapter interface {
	Save(primitive.Identity, *domain.DatasetLabels) error
}
