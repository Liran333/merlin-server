/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package app provides the application layer for the space app service.
package app

import (
	commonapp "github.com/openmerlin/merlin-server/common/app"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	spacedomain "github.com/openmerlin/merlin-server/space/domain"
	"github.com/openmerlin/merlin-server/spaceapp/domain/repository"
)

// SpaceappAppService is the interface for the space app service.
type SpaceappAppService interface {
	GetByName(primitive.Account, *spacedomain.SpaceIndex) (SpaceAppDTO, error)
}

// spaceRepository
type spaceRepository interface {
	FindByName(*spacedomain.SpaceIndex) (spacedomain.Space, error)
}

// NewSpaceappAppService creates a new instance of the space app service.
func NewSpaceappAppService(
	repo repository.Repository,
	spaceRepo spaceRepository,
	permission commonapp.ResourcePermissionAppService,
) *spaceappAppService {
	return &spaceappAppService{
		repo:       repo,
		spaceRepo:  spaceRepo,
		permission: permission,
	}
}

// spaceappAppService
type spaceappAppService struct {
	repo       repository.Repository
	spaceRepo  spaceRepository
	permission commonapp.ResourcePermissionAppService
}

// GetByName retrieves the space app by name.
func (s *spaceappAppService) GetByName(
	user primitive.Account, index *spacedomain.SpaceIndex,
) (SpaceAppDTO, error) {
	var dto SpaceAppDTO

	space, err := s.spaceRepo.FindByName(index)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errorSpaceAppNotFound
		}

		return dto, err
	}

	if err = s.permission.CanRead(user, &space); err != nil {
		if allerror.IsNoPermission(err) {
			err = errorSpaceAppNotFound
		}

		return dto, err
	}

	app, err := s.repo.FindBySpaceId(space.Id)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errorSpaceAppNotFound
		}

		return dto, err
	}

	return toSpaceAppDTO(&app), nil
}
