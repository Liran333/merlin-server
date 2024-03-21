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
	"github.com/openmerlin/merlin-server/spaceapp/domain"
	"github.com/openmerlin/merlin-server/spaceapp/domain/message"
	"github.com/openmerlin/merlin-server/spaceapp/domain/repository"
)

// SpaceappAppService is the interface for the space app service.
type SpaceappAppService interface {
	GetByName(primitive.Account, *spacedomain.SpaceIndex) (SpaceAppDTO, error)
	GetRequestDataStream(*domain.SeverSentStream) error
	RestartSpaceApp(primitive.Account, *spacedomain.SpaceIndex) error
	CheckPermissionRead(primitive.Account, *spacedomain.SpaceIndex) error
}

// spaceRepository
type spaceRepository interface {
	FindByName(*spacedomain.SpaceIndex) (spacedomain.Space, error)
}

// NewSpaceappAppService creates a new instance of the space app service.
func NewSpaceappAppService(
	msg message.SpaceAppMessage,
	repo repository.Repository,
	spaceRepo spaceRepository,
	permission commonapp.ResourcePermissionAppService,
	sse domain.SeverSentEvent,
) *spaceappAppService {
	return &spaceappAppService{
		msg:        msg,
		repo:       repo,
		spaceRepo:  spaceRepo,
		permission: permission,
		sse:        sse,
	}
}

// spaceappAppService
type spaceappAppService struct {
	msg        message.SpaceAppMessage
	repo       repository.Repository
	spaceRepo  spaceRepository
	permission commonapp.ResourcePermissionAppService
	sse        domain.SeverSentEvent
}

// GetByName retrieves the space app by name.
func (s *spaceappAppService) GetByName(
	user primitive.Account, index *spacedomain.SpaceIndex,
) (SpaceAppDTO, error) {
	var dto SpaceAppDTO

	space, err := s.spaceRepo.FindByName(index)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceAppNotFound(err)
		}

		return dto, err
	}

	if err = s.permission.CanRead(user, &space); err != nil {
		if allerror.IsNoPermission(err) {
			err = newSpaceAppNotFound(err)
		}

		return dto, err
	}

	app, err := s.repo.FindBySpaceId(space.Id)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceAppNotFound(err)
		}

		return dto, err
	}

	return toSpaceAppDTO(&app), nil
}

// GetRequestDataStream
func (s *spaceappAppService) GetRequestDataStream(cmd *domain.SeverSentStream) error {
	return s.sse.Request(cmd)
}

// RestartSpaceApp a SpaceApp in the spaceappAppService.
func (s *spaceappAppService) RestartSpaceApp(
	user primitive.Account, index *spacedomain.SpaceIndex,
) error {
	space, err := s.spaceRepo.FindByName(index)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceAppNotFound(err)
		}

		return err
	}

	if err = s.permission.CanUpdate(user, &space); err != nil {
		if allerror.IsNoPermission(err) {
			err = newSpaceAppNotFound(err)
		}

		return err
	}

	app, err := s.repo.FindBySpaceId(space.Id)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceAppNotFound(err)
		}

		return err
	}

	if err := app.RestartService(app.RestartedAt); err != nil {
		return err
	}

	if err := s.repo.Save(&app); err != nil {
		return err
	}

	v := domain.SpaceAppIndex{
		SpaceId:  app.SpaceId,
		CommitId: app.CommitId,
	}
	e := domain.NewSpaceAppRestartEvent(&v)
	return s.msg.SendSpaceAppRestartedEvent(&e)
}

// CheckPermissionRead  check user permission for read space app.
func (s *spaceappAppService) CheckPermissionRead(user primitive.Account, index *spacedomain.SpaceIndex) error {
	space, err := s.spaceRepo.FindByName(index)
	if err != nil {
		err = newSpaceAppNotFound(err)
		return err
	}

	return s.permission.CanRead(user, &space)
}
