/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	computilityapp "github.com/openmerlin/merlin-server/computility/app"
	spacedomain "github.com/openmerlin/merlin-server/space/domain"
	"github.com/openmerlin/merlin-server/spaceapp/domain"
	"github.com/openmerlin/merlin-server/spaceapp/domain/message"
	appprimitive "github.com/openmerlin/merlin-server/spaceapp/domain/primitive"
	"github.com/openmerlin/merlin-server/spaceapp/domain/repository"
)

func newSpaceNotFound(err error) error {
	return allerror.NewNotFound(allerror.ErrorCodeSpaceNotFound, "space not found", err)
}

func newSpaceAppNotFound(err error) error {
	return allerror.NewNotFound(allerror.ErrorCodeSpaceAppNotFound, "space app not found", err)
}

// SpaceappInternalAppService is an interface that defines the methods for creating and managing a SpaceApp.
type SpaceappInternalAppService interface {
	Create(cmd *CmdToCreateApp) error

	NotifyIsInvalid(cmd *CmdToNotifyFailedStatus) error
	NotifyIsBuilding(cmd *CmdToNotifyBuildIsStarted) error
	NotifyIsBuildFailed(cmd *CmdToNotifyFailedStatus) error
	NotifyIsStarting(cmd *CmdToCreateApp) error
	NotifyIsStartFailed(cmd *CmdToNotifyFailedStatus) error
	NotifyIsServing(cmd *CmdToNotifyServiceIsStarted) error
	NotifyIsRestartFailed(cmd *CmdToNotifyFailedStatus) error
	NotifyIsResumeFailed(cmd *CmdToNotifyFailedStatus) error

	ForcePauseSpaceApp(primitive.Identity) error
	PauseSpaceApp(primitive.Identity) error
}

// NewSpaceappInternalAppService creates a new instance of spaceappInternalAppService
// with the provided message and repository.
func NewSpaceappInternalAppService(
	msg message.SpaceAppMessage,
	repo repository.Repository,
	buildLogAdapter repository.SpaceAppBuildLogAdapter,
	spaceRepo spaceRepository,
	computility computilityapp.ComputilityInternalAppService,
) *spaceappInternalAppService {
	return &spaceappInternalAppService{
		msg:             msg,
		repo:            repo,
		buildLogAdapter: buildLogAdapter,
		spaceRepo:       spaceRepo,
		computility:     computility,
	}
}

// spaceappInternalAppService
type spaceappInternalAppService struct {
	msg             message.SpaceAppMessage
	repo            repository.Repository
	buildLogAdapter repository.SpaceAppBuildLogAdapter
	spaceRepo       spaceRepository
	computility     computilityapp.ComputilityInternalAppService
}

// Create creates a new SpaceApp in the spaceappInternalAppService.
func (s *spaceappInternalAppService) Create(cmd *CmdToCreateApp) error {
	v := domain.SpaceApp{
		Status:        appprimitive.AppStatusInit,
		SpaceAppIndex: *cmd,
	}

	space, err := s.spaceRepo.FindById(cmd.SpaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(err)
			logrus.Infof("get space failed, err:%s", err)
		}

		return err
	}

	space.CommitId = cmd.CommitId
	if err := s.spaceRepo.Save(&space); err != nil {
		e := fmt.Errorf("failed to save latest commit id failed, spaceId:%s", space.Id.Identity())
		err = allerror.New(allerror.ErrorCodeSpaceAppCreateFailed, e.Error(), e)
		logrus.Infof("save space failed, err:%s", err)
		return err
	}

	if space.Hardware.IsNpu() && !space.CompPowerAllocated {
		e := fmt.Errorf("failed to create space failed, "+
			"spaceId:%s is npu but not allocate computility", space.Id.Identity())
		err = allerror.New(allerror.ErrorCodeSpaceAppCreateFailed, e.Error(), e)
		logrus.Infof("create space app failed, err:%s", err)
		return err
	}

	app, err := s.repo.FindBySpaceId(space.Id)
	if err == nil && app.IsAppInitAllow() {
		e := fmt.Errorf("spaceId:%s, not allow to init", app.SpaceId.Identity())
		logrus.Infof("create space app failed, err:%s", e)
		return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, e.Error(), e)
	}

	if err := s.repo.Add(&v); err != nil {
		logrus.Info("create space app db failed")
		return err
	}
	e := domain.NewSpaceAppCreatedEvent(&v)

	return s.msg.SendSpaceAppCreatedEvent(&e)
}

func (s *spaceappInternalAppService) getSpaceApp(spaceId primitive.Identity) (domain.SpaceApp, error) {
	space, err := s.spaceRepo.FindById(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(err)
		}

		return domain.SpaceApp{}, err
	}

	v, err := s.repo.FindBySpaceId(space.Id)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceAppNotFound(err)
		}

		return domain.SpaceApp{}, err
	}
	return v, nil
}

// NotifyIsInvalid notifies change SpaceApp status.
func (s *spaceappInternalAppService) NotifyIsInvalid(cmd *CmdToNotifyFailedStatus) error {
	v, err := s.getSpaceApp(cmd.SpaceId)
	if err != nil {
		logrus.Infof("get space app failed, err:%s", err)
		return err
	}
	if err := v.SetInvalid(cmd.Status, cmd.Reason); err != nil {
		logrus.Infof("set space %s status failed, err:%s", cmd.Status.AppStatus(), err)
		return err
	}

	return s.repo.Save(&v)
}

// NotifyIsBuilding notifies that the build process of a SpaceApp has started.
func (s *spaceappInternalAppService) NotifyIsBuilding(cmd *CmdToNotifyBuildIsStarted) error {
	v, err := s.getSpaceApp(cmd.SpaceId)
	if err != nil {
		logrus.Infof("get space app failed, err:%s", err)
		return err
	}

	if err := v.StartBuilding(cmd.LogURL); err != nil {
		logrus.Infof("set  space app building failed, err:%s", err)
		return err
	}

	return s.repo.Save(&v)
}

// NotifyIsBuildFailed notifies change SpaceApp status.
func (s *spaceappInternalAppService) NotifyIsBuildFailed(cmd *CmdToNotifyFailedStatus) error {
	v, err := s.getSpaceApp(cmd.SpaceId)
	if err != nil {
		logrus.Infof("get space app failed, err:%s", err)
		return err
	}
	if err := v.SetBuildFailed(cmd.Status, cmd.Reason); err != nil {
		logrus.Infof("set space %s status failed, err:%s", cmd.Status.AppStatus(), err)
		return err
	}

	return s.repo.Save(&v)
}

// NotifyIsStarting notifies that the build process of a SpaceApp has finished.
func (s *spaceappInternalAppService) NotifyIsStarting(cmd *CmdToCreateApp) error {
	v, err := s.getSpaceApp(cmd.SpaceId)
	if err != nil {
		logrus.Infof("get space app failed, err:%s", err)
		return err
	}

	if err := v.SetStarting(); err != nil {
		logrus.Infof("set space app starting failed, err:%s", err)
		return err
	}

	return s.repo.Save(&v)
}

// NotifyIsBuildFailed notifies change SpaceApp status.
func (s *spaceappInternalAppService) NotifyIsStartFailed(cmd *CmdToNotifyFailedStatus) error {
	v, err := s.getSpaceApp(cmd.SpaceId)
	if err != nil {
		logrus.Infof("get space app failed, err:%s", err)
		return err
	}
	if err := v.SetStartFailed(cmd.Status, cmd.Reason); err != nil {
		logrus.Infof("set space %s status failed, err:%s", cmd.Status.AppStatus(), err)
		return err
	}

	return s.repo.Save(&v)
}

// NotifyIsServing notifies that a service of a SpaceApp has serving.
func (s *spaceappInternalAppService) NotifyIsServing(cmd *CmdToNotifyServiceIsStarted) error {
	v, err := s.getSpaceApp(cmd.SpaceId)
	if err != nil {
		logrus.Infof("get space app failed, err:%s", err)
		return err
	}

	if err := v.StartServing(cmd.AppURL, cmd.LogURL); err != nil {
		logrus.Infof("set space app serving failed, err:%s", err)
		return err
	}

	return s.repo.Save(&v)
}

// NotifyIsReStartFailed notifies change SpaceApp status.
func (s *spaceappInternalAppService) NotifyIsRestartFailed(cmd *CmdToNotifyFailedStatus) error {
	v, err := s.getSpaceApp(cmd.SpaceId)
	if err != nil {
		logrus.Infof("get space app failed, err:%s", err)
		return err
	}
	if err := v.SetRestartFailed(cmd.Status, cmd.Reason); err != nil {
		logrus.Infof("set space %s status failed, err:%s", cmd.Status.AppStatus(), err)
		return err
	}

	return s.repo.Save(&v)
}

// NotifyIsResumeFailed notifies change SpaceApp status.
func (s *spaceappInternalAppService) NotifyIsResumeFailed(cmd *CmdToNotifyFailedStatus) error {
	space, err := s.spaceRepo.FindById(cmd.SpaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(err)
			logrus.Infof("get space failed, err:%s", err)
		}

		return err
	}

	app, err := s.repo.FindBySpaceId(space.Id)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceAppNotFound(err)
			logrus.Infof("get space app failed, err:%s", err)
		}

		return err
	}

	if err := app.SetResumeFailed(cmd.Status, cmd.Reason); err != nil {
		logrus.Infof("set space %s status failed, err:%s", cmd.Status.AppStatus(), err)
		return err
	}

	spaceCompCmd := spaceUserComputilityService {
		userName: space.CreatedBy,
		space: space,
		spaceRepo: s.spaceRepo,
		computility: s.computility,
	}

	if err := spaceCompCmd.unbindSpaceCompQuota(); err != nil {
		err := fmt.Errorf("failed to release spaceId:%s comp quota, err:%s", space.Id.Identity(), err)
		logrus.Infof("set space %s status failed, err:%s", cmd.Status.AppStatus(), err)
		return err
	}

	if err := s.repo.Save(&app); err != nil {
		if err := spaceCompCmd.bindSpaceCompQuota(); err != nil {
			return err
		}
		err := fmt.Errorf("failed to save spaceId:%s db failed, err:%s", space.Id.Identity(), err)
		logrus.Infof("set space %s status failed, err:%s", cmd.Status.AppStatus(), err)
		return err
	}

	return nil
}

// PauseSpaceApp pause a SpaceApp in the spaceappAppService.
func (s *spaceappInternalAppService) ForcePauseSpaceApp(spaceId primitive.Identity) error {
	space, err := s.spaceRepo.FindById(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(err)
			logrus.Infof("get space failed, err:%s", err)
		}

		return err
	}
	spaceCompCmd := spaceUserComputilityService {
		userName: space.CreatedBy,
		space: space,
		spaceRepo: s.spaceRepo,
		computility: s.computility,
	}
	app, err := s.repo.FindBySpaceId(space.Id)
	if err != nil {
		if err := spaceCompCmd.unbindSpaceCompQuota(); err != nil {
			logrus.Infof("release space comp quota failed, err:%s", err)
			return err
		}
		err = newSpaceAppNotFound(err)
		logrus.Infof("get space app failed, err:%s", err)
		return err
	}

	if app.Status.IsPaused() {
		logrus.Info("app is already paused")
		return nil
	}
	app.Status = appprimitive.AppStatusPaused

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
	e := spacedomain.NewSpaceForceEvent(space.Id.Identity(), spacedomain.ForceTypePause)
	return s.msg.SendSpaceAppForcePauseEvent(&e)
}

func (s *spaceappInternalAppService) PauseSpaceApp(spaceId primitive.Identity) error {
	space, err := s.spaceRepo.FindById(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(err)
			logrus.Infof("get space failed, err:%s", err)
		}

		return err
	}
	app, err := s.repo.FindBySpaceId(space.Id)
	if err != nil {
		err = newSpaceAppNotFound(err)
		logrus.Infof("get space app failed, err:%s", err)
		return err
	}

	if err := app.PauseService(); err != nil {
		return err
	}

	spaceCompCmd := spaceUserComputilityService {
		userName: space.CreatedBy,
		space: space,
		spaceRepo: s.spaceRepo,
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
		SpaceId: app.SpaceId,
	})
	return s.msg.SendSpaceAppPauseEvent(&e)
}
