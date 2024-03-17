/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"github.com/openmerlin/merlin-server/search/domain"
	"github.com/openmerlin/merlin-server/search/domain/resourceadapter"
	"github.com/sirupsen/logrus"
)

type SearchAppService interface {
	Search(cmd *CmdToSearch, user domain.Account) (SearchDTO, error)
}

func NewSearchAppService(resourceAdapter resourceadapter.ResourceAdapter) SearchAppService {
	return &searchAppService{
		resourceAdapter: resourceAdapter,
	}
}

type searchAppService struct {
	resourceAdapter resourceadapter.ResourceAdapter
}

func (s *searchAppService) Search(cmd *CmdToSearch, user domain.Account) (SearchDTO, error) {

	var dto SearchDTO

	searchOption := &domain.SearchOption{
		SearchKey:  cmd.SearchKey.SearchKey(),
		SearchType: cmd.SearchType.SearchType(),
		Account:    user,
		Size:       cmd.Size.Size(),
	}

	searchResult, err := s.resourceAdapter.Search(searchOption)
	if err != nil {
		logrus.Error("Failed to search, error: ", err)
		return dto, err
	}

	dto.ResultSet = searchResult

	return dto, nil
}
