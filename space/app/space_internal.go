/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/space/domain/repository"
)

// SpaceInternalAppService is an interface for space internal application service
type SpaceInternalAppService interface {
	GetById(primitive.Identity) (SpaceMetaDTO, error)
}

// NewSpaceInternalAppService creates a new instance of SpaceInternalAppService
func NewSpaceInternalAppService(
	repoAdapter repository.SpaceRepositoryAdapter,
) SpaceInternalAppService {
	return &spaceInternalAppService{
		repoAdapter: repoAdapter,
	}
}

type spaceInternalAppService struct {
	repoAdapter repository.SpaceRepositoryAdapter
}

// GetById retrieves a space by its ID and returns the corresponding SpaceMetaDTO
func (s *spaceInternalAppService) GetById(spaceId primitive.Identity) (SpaceMetaDTO, error) {
	space, err := s.repoAdapter.FindById(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errorSpaceNotFound
		}

		return SpaceMetaDTO{}, err
	}

	return toSpaceMetaDTO(&space), nil
}
