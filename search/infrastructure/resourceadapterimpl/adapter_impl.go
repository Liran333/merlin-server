/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package resourceadapterimpl provides rimpl models and configuration for a specific functionality.
package resourceadapterimpl

import (
	"context"
	"fmt"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	datasetrepo "github.com/openmerlin/merlin-server/datasets/domain/repository"
	modelrepo "github.com/openmerlin/merlin-server/models/domain/repository"
	orgrepo "github.com/openmerlin/merlin-server/organization/domain/repository"
	"github.com/openmerlin/merlin-server/search/domain"
	spacerepo "github.com/openmerlin/merlin-server/space/domain/repository"
	"github.com/openmerlin/merlin-server/user/domain/repository"
	"github.com/openmerlin/merlin-server/utils"
)

// DefaultPage is the default page number
const DefaultPage = 1

// serchAdapter is an implementation of SearchRepositoryAdapter
type searchAdapter struct {
	model   modelrepo.ModelRepositoryAdapter
	dataset datasetrepo.DatasetRepositoryAdapter
	space   spacerepo.SpaceRepositoryAdapter
	user    repository.User
	member  orgrepo.OrgMember
}

// NewSearchRepositoryAdapter creates a new instance of SearchRepositoryAdapter
func NewSearchRepositoryAdapter(
	model modelrepo.ModelRepositoryAdapter,
	dataset datasetrepo.DatasetRepositoryAdapter,
	space spacerepo.SpaceRepositoryAdapter,
	user repository.User,
	member orgrepo.OrgMember,
) *searchAdapter {
	return &searchAdapter{
		model:   model,
		dataset: dataset,
		space:   space,
		user:    user,
		member:  member,
	}
}

// Search for models if the search type includes primitive.SearchTypeModel
func (adapter *searchAdapter) Search(ctx context.Context, opt *domain.SearchOption) (domain.SearchResult, error) {
	var result domain.SearchResult
	if utils.Contains(opt.SearchType, primitive.SearchTypeModel) {
		cmd := &modelrepo.ListOption{
			Name:            opt.SearchKey,
			PageNum:         DefaultPage,
			CountPerPage:    opt.Size,
			ExcludeFullname: true,
			Count:           true,
			Visibility:      primitive.VisibilityPublic,
		}
		models, err := adapter.SearchModel(ctx, cmd, opt.Account)
		if err != nil {
			return result, err
		}
		result.SearchResultModel = models
	}

	if utils.Contains(opt.SearchType, primitive.SearchTypeDataset) {
		cmd := &datasetrepo.ListOption{
			Name:            opt.SearchKey,
			PageNum:         DefaultPage,
			CountPerPage:    opt.Size,
			ExcludeFullname: true,
			Count:           true,
			Visibility:      primitive.VisibilityPublic,
		}
		datasets, err := adapter.SearchDataset(cmd, opt.Account)
		if err != nil {
			return result, err
		}
		result.SearchResultDataset = datasets
	}

	if utils.Contains(opt.SearchType, primitive.SearchTypeSpace) {
		cmd := &spacerepo.ListOption{
			Name:            opt.SearchKey,
			PageNum:         DefaultPage,
			CountPerPage:    opt.Size,
			ExcludeFullname: true,
			Count:           true,
			Visibility:      primitive.VisibilityPublic,
		}
		spaces, err := adapter.SearchSpace(cmd, opt.Account)
		if err != nil {
			return result, err
		}
		result.SearchResultSpace = spaces
	}

	if utils.Contains(opt.SearchType, primitive.SearchTypeUser) {
		cmd := &repository.ListOption{
			Name:         opt.SearchKey,
			PageNum:      DefaultPage,
			CountPerPage: opt.Size,
		}
		users, err := adapter.SearchUser(ctx, cmd)
		if err != nil {
			return result, err
		}
		result.SearchResultUser = users
	}

	if utils.Contains(opt.SearchType, primitive.SearchTypeOrg) {
		cmd := &repository.ListOption{
			Name:         opt.SearchKey,
			PageNum:      DefaultPage,
			CountPerPage: opt.Size,
		}
		orgs, err := adapter.SearchOrg(ctx, cmd)
		if err != nil {
			return result, err
		}
		result.SearchResultOrg = orgs
	}

	return result, nil
}

// SearchModel for search models
func (adapter *searchAdapter) SearchModel(ctx context.Context, cmd *modelrepo.ListOption, account domain.Account) (
	domain.SearchResultModel, error) {
	var result domain.SearchResultModel
	v, count, err := adapter.model.List(ctx, cmd, account, adapter.member)
	if err != nil {
		return result, err
	}
	models := make([]domain.ModelResult, 0)
	for _, m := range v {
		models = append(models, domain.ModelResult{
			Owner: m.Owner,
			Name:  m.Name,
			Path:  fmt.Sprintf("%s/%s", m.Owner, m.Name),
		})
	}
	result.ModelResult = models
	result.ModelResultCount = count
	return result, nil
}

// SearchDataset for search datasets
func (adapter *searchAdapter) SearchDataset(cmd *datasetrepo.ListOption, account domain.Account) (
	domain.SearchResultDataset, error) {
	var result domain.SearchResultDataset
	v, count, err := adapter.dataset.List(cmd, account, adapter.member)
	if err != nil {
		return result, err
	}
	datasets := make([]domain.DatasetResult, 0)
	for _, m := range v {
		datasets = append(datasets, domain.DatasetResult{
			Owner: m.Owner,
			Name:  m.Name,
			Path:  fmt.Sprintf("%s/%s", m.Owner, m.Name),
		})
	}
	result.DatasetResult = datasets
	result.DatasetResultCount = count
	return result, nil
}

// SearchSpace for search spaces
func (adapter *searchAdapter) SearchSpace(cmd *spacerepo.ListOption, account domain.Account) (
	domain.SearchResultSpace, error) {
	var result domain.SearchResultSpace

	v, count, err := adapter.space.List(cmd, account, adapter.member)
	if err != nil {
		return result, err
	}
	spaces := make([]domain.SpaceResult, 0)
	for _, s := range v {
		spaces = append(spaces, domain.SpaceResult{
			Owner: s.Owner,
			Name:  s.Name,
			Path:  fmt.Sprintf("%s/%s", s.Owner, s.Name),
		})
	}
	result.SpaceResult = spaces
	result.SpaceResultCount = count
	return result, nil
}

// SearchUser provides a method to search for users
func (adapter *searchAdapter) SearchUser(
	ctx context.Context, cmd *repository.ListOption) (domain.SearchResultUser, error) {
	var result domain.SearchResultUser

	v, count, err := adapter.user.SearchUser(ctx, cmd)
	if err != nil {
		return result, err
	}

	users := make([]domain.UserResult, 0)
	for _, u := range v {
		users = append(users, domain.UserResult{
			Account:  u.Account.Account(),
			FullName: u.Fullname.AccountFullname(),
			AvatarId: u.AvatarId.AvatarId(),
		})
	}
	result.UserResult = users
	result.UserResultCount = count
	return result, nil
}

// SearchOrg provides a method to search for orgs
func (adapter *searchAdapter) SearchOrg(
	ctx context.Context, cmd *repository.ListOption) (domain.SearchResultOrg, error) {
	var result domain.SearchResultOrg

	v, count, err := adapter.user.SearchOrg(ctx, cmd)
	if err != nil {
		return result, err
	}

	orgs := make([]domain.OrgResult, 0)
	for _, u := range v {
		orgs = append(orgs, domain.OrgResult{
			Id:       u.Id.Identity(),
			Name:     u.Account.Account(),
			FullName: u.Fullname.AccountFullname(),
		})
	}
	result.OrgResult = orgs
	result.OrgResultCount = count
	return result, nil
}
