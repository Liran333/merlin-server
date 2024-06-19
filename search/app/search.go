/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package app provides application services for creating and managing branches.
package app

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/openmerlin/merlin-server/search/domain"
	"github.com/openmerlin/merlin-server/search/domain/resourceadapter"
)

// SearchAppService search app service
type SearchAppService interface {
	Search(ctx context.Context, cmd *CmdToSearch, user domain.Account) (SearchDTO, error)
}

// NewSearchAppService creates a new instance of the SearchAppService.
func NewSearchAppService(resourceAdapter resourceadapter.ResourceAdapter) SearchAppService {
	return &searchAppService{
		resourceAdapter: resourceAdapter,
	}
}

// searchAppService implements the SearchAppService interface.
type searchAppService struct {
	resourceAdapter resourceadapter.ResourceAdapter
}

// Search is a method of the SearchAppService interface that performs a search operation.
func (s *searchAppService) Search(ctx context.Context, cmd *CmdToSearch, user domain.Account) (SearchDTO, error) {

	var dto SearchDTO

	searchOption := &domain.SearchOption{
		SearchKey:  cmd.SearchKey.SearchKey(),
		SearchType: cmd.SearchType.SearchType(),
		Account:    user,
		Size:       cmd.Size.Size(),
	}

	searchResult, err := s.resourceAdapter.Search(ctx, searchOption)
	if err != nil {
		logrus.Error("Failed to search, error: ", err)
		return dto, err
	}

	dto.ResultSet = searchResult

	return dto, nil
}
