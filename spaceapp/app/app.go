/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package app provides the application layer for the space app service.
package app

import (
	"fmt"

	"github.com/sirupsen/logrus"

	commonapp "github.com/openmerlin/merlin-server/common/app"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	computilityapp "github.com/openmerlin/merlin-server/computility/app"
	computilitydomain "github.com/openmerlin/merlin-server/computility/domain"
	modelrepo "github.com/openmerlin/merlin-server/models/domain/repository"
	spacedomain "github.com/openmerlin/merlin-server/space/domain"
	spacerepo "github.com/openmerlin/merlin-server/space/domain/repository"
	"github.com/openmerlin/merlin-server/spaceapp/domain"
	"github.com/openmerlin/merlin-server/spaceapp/domain/message"
	appprimitive "github.com/openmerlin/merlin-server/spaceapp/domain/primitive"
	"github.com/openmerlin/merlin-server/spaceapp/domain/repository"
)

// SpaceappAppService is the interface for the space app service.
type SpaceappAppService interface {
	GetByName(primitive.Account, *spacedomain.SpaceIndex) (SpaceAppDTO, error)
	GetRequestDataStream(*domain.SeverSentStream) error
	RestartSpaceApp(primitive.Account, *spacedomain.SpaceIndex) error
	PauseSpaceApp(primitive.Account, *spacedomain.SpaceIndex) error
	ResumeSpaceApp(primitive.Account, *spacedomain.SpaceIndex) error
	CheckPermissionRead(primitive.Account, *spacedomain.SpaceIndex) error
	GetSpaceIdByName(index *spacedomain.SpaceIndex) (spacedomain.Space, error)
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
	computility computilityapp.ComputilityInternalAppService,
	repoAdapterModelSpace spacerepo.ModelSpaceRepositoryAdapter,
	modelRepoAdapter modelrepo.ModelRepositoryAdapter,
) *spaceappAppService {
	return &spaceappAppService{
		msg:                   msg,
		repo:                  repo,
		spaceRepo:             spaceRepo,
		permission:            permission,
		sse:                   sse,
		computility:           computility,
		repoAdapterModelSpace: repoAdapterModelSpace,
		modelRepoAdapter:      modelRepoAdapter,
	}
}

// spaceappAppService
type spaceappAppService struct {
	msg                   message.SpaceAppMessage
	repo                  repository.Repository
	spaceRepo             spaceRepository
	permission            commonapp.ResourcePermissionAppService
	sse                   domain.SeverSentEvent
	computility           computilityapp.ComputilityInternalAppService
	repoAdapterModelSpace spacerepo.ModelSpaceRepositoryAdapter
	modelRepoAdapter      modelrepo.ModelRepositoryAdapter
}

func (s *spaceappAppService) canHandleNotDisable(space *spacedomain.Space) error {
	if space.IsDisable() {
		errInfo := fmt.Sprintf("space %v was disable", space.Name.MSDName())
		logrus.Errorf("%s, do not allow to restart or resume", errInfo)
		return allerror.NewResourceDisabled(allerror.ErrorCodeResourceDisabled, errInfo, fmt.Errorf("resource disabled"))
	}

	modelIds, err := s.repoAdapterModelSpace.GetModelsBySpaceId(space.Id)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			return newSpaceNotFound(err)
		}

		return err
	}

	for _, id := range modelIds {
		model, err := s.modelRepoAdapter.FindById(id)
		if err != nil {
			if commonrepo.IsErrorResourceNotExists(err) {
				continue
			}
			logrus.Errorf("find model by id failed, id:%v, err:%v", id, err)
			return err
		}

		if model.IsDisable() {
			errInfo := fmt.Sprintf("related model %v was disable", model.Name.MSDName())
			logrus.Errorf("%s, do not allow to handle", errInfo)
			return allerror.NewResourceDisabled(allerror.ErrorCodeResourceDisabled, errInfo, fmt.Errorf("resource disabled"))
		}
	}

	return err
}

// spaceUserComputilityService
type spaceUserComputilityService struct {
	userName    primitive.Account
	space       spacedomain.Space
	spaceRepo   spaceRepository
	computility computilityapp.ComputilityInternalAppService
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

	if err = s.permission.CanReadPrivate(user, &space); err != nil {
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

	if err = s.canHandleNotDisable(&space); err != nil {
		logrus.Errorf("space %v cant restart, because was disabled", space.Name.MSDName())
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
		err = newSpaceNotFound(err)
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
			err = newSpaceNotFound(err)
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

	if err := app.PauseService(); err != nil {
		return err
	}

	spaceCompCmd := spaceUserComputilityService{
		userName:    user,
		space:       space,
		spaceRepo:   s.spaceRepo,
		computility: s.computility,
	}

	if err := spaceCompCmd.unbindSpaceCompQuota(); err != nil {
		e := fmt.Errorf("failed to release spaceId:%s comp quota, err:%s", space.Id.Identity(), err)
		err = allerror.New(allerror.ErrorCodeSpaceAppPauseFailed, e.Error(), e)
		return err
	}

	if err := s.repo.Save(&app); err != nil {
		if err := spaceCompCmd.bindSpaceCompQuota(); err != nil {
			return err
		}
		e := fmt.Errorf("failed to save spaceId:%s db failed, err:%s", space.Id.Identity(), err)
		err = allerror.New(allerror.ErrorCodeSpaceAppPauseFailed, e.Error(), e)
		return err
	}

	e := domain.NewSpaceAppPauseEvent(&domain.SpaceAppIndex{
		SpaceId: space.Id,
	})
	return s.msg.SendSpaceAppPauseEvent(&e)
}

func (sc *spaceUserComputilityService) unbindSpaceCompQuota() error {
	if !sc.space.CompPowerAllocated {
		logrus.Info("no allocated power, no need release")
		return nil
	}
	compCmd := computilityapp.CmdToUserQuotaUpdate{
		Index: computilitydomain.ComputilityAccountRecordIndex{
			UserName:    sc.userName,
			ComputeType: sc.space.GetComputeType(),
			SpaceId:     sc.space.Id,
		},
		QuotaCount: sc.space.GetQuotaCount(),
	}
	if err := sc.computility.UserQuotaRelease(compCmd); err != nil {
		return err
	}
	sc.space.CompPowerAllocated = false
	if err := sc.spaceRepo.Save(&sc.space); err != nil {
		return sc.computility.UserQuotaConsume(compCmd)
	}
	return nil
}

func (sc *spaceUserComputilityService) bindSpaceCompQuota() error {
	if !sc.space.Hardware.IsNpu() {
		logrus.Info("no allow consume type")
		return nil
	}
	if sc.space.CompPowerAllocated {
		return fmt.Errorf("already bind power, no consume")
	}
	compCmd := computilityapp.CmdToUserQuotaUpdate{
		Index: computilitydomain.ComputilityAccountRecordIndex{
			UserName:    sc.userName,
			ComputeType: sc.space.GetComputeType(),
			SpaceId:     sc.space.Id,
		},
		QuotaCount: sc.space.GetQuotaCount(),
	}
	if err := sc.computility.UserQuotaConsume(compCmd); err != nil {
		return err
	}
	sc.space.CompPowerAllocated = true
	if err := sc.spaceRepo.Save(&sc.space); err != nil {
		return sc.computility.UserQuotaRelease(compCmd)
	}
	return nil
}

// ResumeSpaceApp a SpaceApp in the spaceappAppService.
func (s *spaceappAppService) ResumeSpaceApp(
	user primitive.Account, index *spacedomain.SpaceIndex,
) error {
	space, err := s.spaceRepo.FindByName(index)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(err)
		}

		return err
	}

	if err = s.canHandleNotDisable(&space); err != nil {
		logrus.Errorf("space %v cant resume, because was disabled", space.Name.MSDName())
		return err
	}

	if err = s.permission.CanUpdate(user, &space); err != nil {
		if allerror.IsNoPermission(err) {
			e := fmt.Errorf("no permission to exec spaceId:%s,err:%s", space.Id.Identity(), err)
			err = allerror.NewNotFound(allerror.ErrorCodeSpaceNotFound, "not found", e)
		}

		return err
	}

	spaceCompCmd := spaceUserComputilityService{
		userName:    user,
		space:       space,
		spaceRepo:   s.spaceRepo,
		computility: s.computility,
	}

	if err := spaceCompCmd.bindSpaceCompQuota(); err != nil {
		err := allerror.New(allerror.ErrorCodeSpaceAppResumeFailed, "resume space failed", err)
		return err
	}

	app, err := s.repo.FindBySpaceId(space.Id)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceAppNotFound(err)
		}
		if !space.Hardware.IsNpu() {
			return err
		}
		app, err = s.reCreateApp(space)
		if err != nil {
			return err
		}
	}

	if err := app.ResumeService(); err != nil {
		e := fmt.Errorf("resume spaceId:%s failed, err:%s", space.Id.Identity(), err)
		err = allerror.New(allerror.ErrorCodeSpaceAppResumeFailed, "resume space failed", e)
		return err
	}

	if err := s.repo.Save(&app); err != nil {
		if err := spaceCompCmd.unbindSpaceCompQuota(); err != nil {
			e := fmt.Errorf("failed to release spaceId:%s comp quota, err:%s", space.Id.Identity(), err)
			err = allerror.New(allerror.ErrorCodeSpaceAppPauseFailed, e.Error(), e)
			return err
		}
		e := fmt.Errorf("update resuming spaceId:%s failed, err:%s", space.Id.Identity(), err)
		err = allerror.New(allerror.ErrorCodeSpaceAppResumeFailed, "update space failed", e)
		return err
	}

	e := domain.NewSpaceAppResumeEvent(&domain.SpaceAppIndex{
		SpaceId:  app.SpaceId,
		CommitId: space.CommitId,
	})
	return s.msg.SendSpaceAppResumeEvent(&e)
}

func (s *spaceappAppService) reCreateApp(space spacedomain.Space) (domain.SpaceApp, error) {
	if err := s.repo.Add(&domain.SpaceApp{
		Status: appprimitive.AppStatusInit,
		SpaceAppIndex: domain.SpaceAppIndex{
			SpaceId:  space.Id,
			CommitId: space.CommitId,
		},
	}); err != nil {
		return domain.SpaceApp{}, err
	}
	return s.repo.FindBySpaceId(space.Id)
}

// GetSpaceIdByName get space id by name.
func (s *spaceappAppService) GetSpaceIdByName(index *spacedomain.SpaceIndex) (spacedomain.Space, error) {
	space, err := s.spaceRepo.FindByName(index)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceAppNotFound(err)
		}

		return spacedomain.Space{}, err
	}

	return space, nil
}
