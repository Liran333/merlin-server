/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"fmt"

	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	computilityapp "github.com/openmerlin/merlin-server/computility/app"
	computilitydomain "github.com/openmerlin/merlin-server/computility/domain"
	spacedomain "github.com/openmerlin/merlin-server/space/domain"
	"github.com/openmerlin/merlin-server/spaceapp/domain"
	"github.com/openmerlin/merlin-server/spaceapp/domain/message"
	appprimitive "github.com/openmerlin/merlin-server/spaceapp/domain/primitive"
	"github.com/openmerlin/merlin-server/spaceapp/domain/repository"
)

func newSpaceAppNotFound(err error) error {
	return allerror.NewNotFound(allerror.ErrorCodeSpaceAppNotFound, "not found", err)
}

// SpaceappInternalAppService is an interface that defines the methods for creating and managing a SpaceApp.
type SpaceappInternalAppService interface {
	Create(cmd *CmdToCreateApp) error
	NotifyBuildIsStarted(cmd *CmdToNotifyBuildIsStarted) error
	NotifyBuildIsDone(cmd *CmdToNotifyBuildIsDone) error
	NotifyServiceIsStarted(cmd *CmdToNotifyServiceIsStarted) error
	NotifyUpdateStatus(cmd *CmdToNotifyUpdateStatus) error
	PauseSpaceApp(*spacedomain.SpaceIndex, bool) error
}

// NewSpaceappInternalAppService creates a new instance of spaceappInternalAppService
// with the provided message and repository.
func NewSpaceappInternalAppService(
	msg message.SpaceAppMessage,
	repo repository.Repository,
	buildLogAdapter repository.SpaceAppBuildLogAdapter,
	spaceRepo spaceRepository,
	compUtility computilityapp.ComputilityAppService,
) *spaceappInternalAppService {
	return &spaceappInternalAppService{
		msg:             msg,
		repo:            repo,
		buildLogAdapter: buildLogAdapter,
		spaceRepo:       spaceRepo,
		compUtility:     compUtility,
	}
}

// spaceappInternalAppService
type spaceappInternalAppService struct {
	msg             message.SpaceAppMessage
	repo            repository.Repository
	buildLogAdapter repository.SpaceAppBuildLogAdapter
	spaceRepo       spaceRepository
	compUtility     computilityapp.ComputilityAppService
}

// Create creates a new SpaceApp in the spaceappInternalAppService.
func (s *spaceappInternalAppService) Create(cmd *CmdToCreateApp) error {

	v := domain.SpaceApp{
		Status:        appprimitive.AppStatusInit,
		SpaceAppIndex: *cmd,
	}

	if err := s.repo.Add(&v); err != nil {
		return err
	}

	e := domain.NewSpaceAppCreatedEvent(&v)

	return s.msg.SendSpaceAppCreatedEvent(&e)
}

// NotifyBuildIsStarted notifies that the build process of a SpaceApp has started.
func (s *spaceappInternalAppService) NotifyBuildIsStarted(cmd *CmdToNotifyBuildIsStarted) error {
	v, err := s.repo.Find(&cmd.SpaceAppIndex)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceAppNotFound(err)
		}

		return err
	}

	if err := v.StartBuilding(cmd.LogURL); err != nil {
		return err
	}

	return s.repo.Save(&v)
}

// NotifyBuildIsDone notifies that the build process of a SpaceApp has finished.
func (s *spaceappInternalAppService) NotifyBuildIsDone(cmd *CmdToNotifyBuildIsDone) error {
	v, err := s.repo.Find(&cmd.SpaceAppIndex)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceAppNotFound(err)
		}

		return err
	}

	if err := v.SetBuildIsDone(cmd.Success); err != nil {
		return err
	}

	log := domain.SpaceAppBuildLog{
		AppId: v.Id,
		Logs:  cmd.Logs,
	}

	if err := s.buildLogAdapter.Save(&log); err != nil {
		return err
	}

	return s.repo.Save(&v)
}

// NotifyServiceIsStarted notifies that a service of a SpaceApp has started.
func (s *spaceappInternalAppService) NotifyServiceIsStarted(cmd *CmdToNotifyServiceIsStarted) error {
	v, err := s.repo.Find(&cmd.SpaceAppIndex)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceAppNotFound(err)
		}

		return err
	}

	if err := v.StartService(cmd.AppURL, cmd.LogURL); err != nil {
		return err
	}

	return s.repo.Save(&v)
}

// NotifyUpdateStatus notifies change SpaceApp status.
func (s *spaceappInternalAppService) NotifyUpdateStatus(cmd *CmdToNotifyUpdateStatus) error {
	v, err := s.repo.Find(&cmd.SpaceAppIndex)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceAppNotFound(err)
		}

		return err
	}

	if err := v.IsAppStatusAllow(cmd.Status); err != nil {
		return err
	}

	return s.repo.Save(&v)
}

// PauseSpaceApp pause a SpaceApp in the spaceappAppService.
func (s *spaceappInternalAppService) PauseSpaceApp(
	index *spacedomain.SpaceIndex, isForce bool,
) error {
	space, err := s.spaceRepo.FindByName(index)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeSpaceNotFound, "not found", err)
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

	releaseSpaceCompQuota := func(compPowerAllocated bool,compUtility computilityapp.ComputilityAppService) error {
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

	if err := app.PauseService(isForce, space.CompPowerAllocated,
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

	v := domain.SpaceAppIndex{
		SpaceId:  app.SpaceId,
		CommitId: app.CommitId,
	}
	e := domain.NewSpaceAppPauseEvent(&v)
	return s.msg.SendSpaceAppPauseEvent(&e)
}
