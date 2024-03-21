/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	sdk "github.com/openmerlin/merlin-sdk/space"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/space/domain/repository"
)

// SpaceInternalAppService is an interface for space internal application service
type SpaceInternalAppService interface {
	GetById(primitive.Identity) (sdk.SpaceMetaDTO, error)
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
func (s *spaceInternalAppService) GetById(spaceId primitive.Identity) (sdk.SpaceMetaDTO, error) {
	space, err := s.repoAdapter.FindById(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(err)
		}

		return sdk.SpaceMetaDTO{}, err
	}

	return toSpaceMetaDTO(&space), nil
}
