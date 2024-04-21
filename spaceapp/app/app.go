/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package app provides the application layer for the space app service.
package app

import (
	"fmt"

	commonapp "github.com/openmerlin/merlin-server/common/app"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	computilityapp "github.com/openmerlin/merlin-server/computility/app"
	computilitydomain "github.com/openmerlin/merlin-server/computility/domain"
	spacedomain "github.com/openmerlin/merlin-server/space/domain"
	"github.com/openmerlin/merlin-server/spaceapp/domain"
	"github.com/openmerlin/merlin-server/spaceapp/domain/message"
	"github.com/openmerlin/merlin-server/spaceapp/domain/repository"
)

const compUtilityTypeNpu = "npu"

// SpaceappAppService is the interface for the space app service.
type SpaceappAppService interface {
	GetByName(primitive.Account, *spacedomain.SpaceIndex) (SpaceAppDTO, error)
	GetRequestDataStream(*domain.SeverSentStream) error
	RestartSpaceApp(primitive.Account, *spacedomain.SpaceIndex) error
	PauseSpaceApp(primitive.Account, *spacedomain.SpaceIndex) error
	ResumeSpaceApp(primitive.Account, *spacedomain.SpaceIndex) error
	CheckPermissionRead(primitive.Account, *spacedomain.SpaceIndex) error
}

// spaceRepository
type spaceRepository interface {
	FindByName(*spacedomain.SpaceIndex) (spacedomain.Space, error)
	FindById(primitive.Identity) (spacedomain.Space, error)
	Save(*spacedomain.Space) error
}

// NewSpaceappAppService creates a new instance of the space app service.
func NewSpaceappAppService(
	msg message.SpaceAppMessage,
	repo repository.Repository,
	spaceRepo spaceRepository,
	permission commonapp.ResourcePermissionAppService,
	sse domain.SeverSentEvent,
	compUtility computilityapp.ComputilityAppService,
) *spaceappAppService {
	return &spaceappAppService{
		msg:         msg,
		repo:        repo,
		spaceRepo:   spaceRepo,
		permission:  permission,
		sse:         sse,
		compUtility: compUtility,
	}
}

// spaceappAppService
type spaceappAppService struct {
	msg         message.SpaceAppMessage
	repo        repository.Repository
	spaceRepo   spaceRepository
	permission  commonapp.ResourcePermissionAppService
	sse         domain.SeverSentEvent
	compUtility computilityapp.ComputilityAppService
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

	if err := app.RestartService(); err != nil {
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

// PauseSpaceApp pause a SpaceApp in the spaceappAppService.
func (s *spaceappAppService) PauseSpaceApp(
	user primitive.Account, index *spacedomain.SpaceIndex,
) error {
	space, err := s.spaceRepo.FindByName(index)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeSpaceNotFound, "not found", err)
		}

		return err
	}

	if err = s.permission.CanUpdate(user, &space); err != nil {
		if allerror.IsNoPermission(err) {
			e := fmt.Errorf("no permission to exec spaceId:%s,err:%s", space.Id.Identity(), err)
			err = allerror.NewNotFound(allerror.ErrorCodeSpaceNotFound, "not found", e)
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

	if app.Status.IsPaused() {
		return nil
	}

	releaseSpaceCompQuota := func(compPowerAllocated bool, compUtility computilityapp.ComputilityAppService) error {
		if !compPowerAllocated {
			return nil
		}
		npu, err := primitive.NewComputilityType(compUtilityTypeNpu)
		if err != nil {
			return err
		}
		cmd := computilityapp.CmdToUserQuotaUpdate{
			ComputilityAccountIndex: computilitydomain.ComputilityAccountIndex{
				UserName:    space.CreatedBy,
				ComputeType: npu,
			},
			QuotaCount: 1,
		}
		return compUtility.UserQuotaRelease(cmd)
	}

	if err := app.PauseService(false, space.CompPowerAllocated,
		s.compUtility, releaseSpaceCompQuota); err != nil {
		return err
	}

	if space.Hardware.IsNpu() {
		space.CompPowerAllocated = false
		if err := s.spaceRepo.Save(&space); err != nil {
			return err
		}
	}

	if err := s.repo.Save(&app); err != nil {
		e := fmt.Errorf("failed to update pause spaceId:%s, err:%s", space.Id.Identity(), err)
		err = allerror.New(allerror.ErrorCodeSpaceAppPauseFailed, e.Error(), e)
		return err
	}

	e := domain.NewSpaceAppPauseEvent(&domain.SpaceAppIndex{
		SpaceId:  app.SpaceId,
		CommitId: app.CommitId,
	})
	return s.msg.SendSpaceAppPauseEvent(&e)
}

// ResumeSpaceApp a SpaceApp in the spaceappAppService.
func (s *spaceappAppService) ResumeSpaceApp(
	user primitive.Account, index *spacedomain.SpaceIndex,
) error {
	space, err := s.spaceRepo.FindByName(index)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeSpaceNotFound, "not found", err)
		}

		return err
	}

	if err = s.permission.CanUpdate(user, &space); err != nil {
		if allerror.IsNoPermission(err) {
			e := fmt.Errorf("no permission to exec spaceId:%s,err:%s", space.Id.Identity(), err)
			err = allerror.NewNotFound(allerror.ErrorCodeSpaceNotFound, "not found", e)
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

	applySpaceCompQuota := func(isNpu bool, compUtility computilityapp.ComputilityAppService) error {
		if !isNpu {
			return nil
		}
		npu, err := primitive.NewComputilityType(compUtilityTypeNpu)
		if err != nil {
			return err
		}
		cmd := computilityapp.CmdToUserQuotaUpdate{
			ComputilityAccountIndex: computilitydomain.ComputilityAccountIndex{
				UserName:    space.CreatedBy,
				ComputeType: npu,
			},
			QuotaCount: 1,
		}
		return compUtility.UserQuotaConsume(cmd)
	}

	if err := app.ResumeService(space.Hardware.IsNpu(), s.compUtility, applySpaceCompQuota); err != nil {
		e := fmt.Errorf("resume spaceId:%s failed, err:%s", space.Id.Identity(), err)
		err = allerror.New(allerror.ErrorCodeSpaceAppResumeFailed, "resume space failed", e)
		return err
	}

	if space.Hardware.IsNpu() {
		space.CompPowerAllocated = true
		if err := s.spaceRepo.Save(&space); err != nil {
			return err
		}
	}

	if err := s.repo.Save(&app); err != nil {
		e := fmt.Errorf("update resuming spaceId:%s failed, err:%s", space.Id.Identity(), err)
		err = allerror.New(allerror.ErrorCodeSpaceAppResumeFailed, "update space failed", e)
		return err
	}

	e := domain.NewSpaceAppResumeEvent(&domain.SpaceAppIndex{
		SpaceId:  app.SpaceId,
		CommitId: app.CommitId,
	})
	return s.msg.SendSpaceAppResumeEvent(&e)
}
