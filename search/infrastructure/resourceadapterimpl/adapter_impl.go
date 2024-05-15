/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package resourceadapterimpl provides rimpl models and configuration for a specific functionality.
package resourceadapterimpl

import (
	"fmt"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	modelrepo "github.com/openmerlin/merlin-server/models/domain/repository"
	orgrepo "github.com/openmerlin/merlin-server/organization/domain/repository"
	"github.com/openmerlin/merlin-server/search/domain"
	spacerepo "github.com/openmerlin/merlin-server/space/domain/repository"
	"github.com/openmerlin/merlin-server/user/domain/repository"
	"github.com/openmerlin/merlin-server/utils"
)

const DefaultPage = 1

type searchAdapter struct {
	model  modelrepo.ModelRepositoryAdapter
	space  spacerepo.SpaceRepositoryAdapter
	user   repository.User
	member orgrepo.OrgMember
}

func NewSearchRepositoryAdapter(
	model modelrepo.ModelRepositoryAdapter,
	space spacerepo.SpaceRepositoryAdapter,
	user repository.User,
	member orgrepo.OrgMember,
) *searchAdapter {
	return &searchAdapter{
		model:  model,
		space:  space,
		user:   user,
		member: member,
	}
}

func (adapter *searchAdapter) Search(opt *domain.SearchOption) (domain.SearchResult, error) {
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
		models, err := adapter.SearchModel(cmd, opt.Account)
		if err != nil {
			return result, err
		}
		result.SearchResultModel = models
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
		users, err := adapter.SearchUser(cmd)
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
		orgs, err := adapter.SearchOrg(cmd)
		if err != nil {
			return result, err
		}
		result.SearchResultOrg = orgs
	}

	return result, nil
}

func (adapter *searchAdapter) SearchModel(cmd *modelrepo.ListOption, account domain.Account) (
	domain.SearchResultModel, error) {
	var result domain.SearchResultModel
	v, count, err := adapter.model.List(cmd, account, adapter.member)
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

func (adapter *searchAdapter) SearchUser(cmd *repository.ListOption) (domain.SearchResultUser, error) {
	var result domain.SearchResultUser

	v, count, err := adapter.user.SearchUser(cmd)
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

func (adapter *searchAdapter) SearchOrg(cmd *repository.ListOption) (domain.SearchResultOrg, error) {
	var result domain.SearchResultOrg

	v, count, err := adapter.user.SearchOrg(cmd)
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
