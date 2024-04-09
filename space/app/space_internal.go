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
	UpdateLocalCMD(spaceId primitive.Identity, cmd string) error
	UpdateEnvInfo(spaceId primitive.Identity, envInfo string) error
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

func (s *spaceInternalAppService) UpdateLocalCMD(spaceId primitive.Identity, cmd string) error {
	space, err := s.repoAdapter.FindById(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(err)
		}

		return err
	}

	space.LocalCmd = cmd
	return s.repoAdapter.Save(&space)
}

func (s *spaceInternalAppService) UpdateEnvInfo(spaceId primitive.Identity, envInfo string) error {
	space, err := s.repoAdapter.FindById(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(err)
		}

		return err
	}

	space.LocalEnvInfo = envInfo
	return s.repoAdapter.Save(&space)
}
