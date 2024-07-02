/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package repository provides interfaces for interacting with models and model labels in the domain.
package repository

import (
	"context"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/models/domain"
	orgrepo "github.com/openmerlin/merlin-server/organization/domain/repository"
)

// ModelSummary represents a summary of a model.
type ModelSummary struct {
	Id            string   `json:"id"`
	Name          string   `json:"name"`
	Desc          string   `json:"desc"`
	Task          string   `json:"task"`
	Owner         string   `json:"owner"`
	Licenses      []string `json:"license"`
	Fullname      string   `json:"fullname"`
	UpdatedAt     int64    `json:"updated_at"`
	LikeCount     int      `json:"like_count"`
	Frameworks    []string `json:"frameworks"`
	DownloadCount int      `json:"download_count"`
	Disable       bool     `json:"disable"`
	DisableReason string   `json:"disable_reason"`
}

// ListOption represents options for listing models.
type ListOption struct {
	// can't define Name as domain.ResourceName
	// because the Name can be subpart of the real resource name
	Name string

	ExcludeFullname bool

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
	InternalSaveStatistic(*domain.Model) error
	InternalSaveUseInOpenmind(*domain.Model) error
	List(context.Context, *ListOption, primitive.Account, orgrepo.OrgMember) ([]ModelSummary, int, error)
	Count(*ListOption) (int, error)
	AddLike(domain.Model) error
	DeleteLike(domain.Model) error
	FindByModelName(string) (domain.Model, error)
}

// ModelLabelsRepoAdapter represents an interface for managing model labels.
type ModelLabelsRepoAdapter interface {
	Save(primitive.Identity, *domain.ModelLabels) error
}

type ModelDeployRepoAdapter interface {
	Create(domain.ModelIndex, []domain.Deploy) error
	DeleteByOwnerName(domain.ModelIndex) error
	FindByOwnerName(*domain.ModelIndex) ([]domain.Deploy, error)
}
