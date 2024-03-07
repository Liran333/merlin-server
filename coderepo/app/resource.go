/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/coderepo/domain/resourceadapter"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
)

// ResourceAppService is an interface for code local repository application service.
type ResourceAppService interface {
	IsRepoExist(*domain.CodeRepoIndex) (bool, error)
}

func NewResourceAppService(r resourceadapter.ResourceAdapter) *resourceAppService {
	return &resourceAppService{resource: r}
}

type resourceAppService struct {
	resource resourceadapter.ResourceAdapter
}

// IsRepoExist check whether the repo is exists
func (s *resourceAppService) IsRepoExist(index *domain.CodeRepoIndex) (bool, error) {
	_, err := s.resource.GetByName(index)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		}

		return false, err
	}

	return true, nil
}
