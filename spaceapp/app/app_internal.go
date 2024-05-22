/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"

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
	space, err := s.spaceRepo.FindById(cmd.SpaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(err)
			logrus.Errorf("spaceId：%s get space failed, err:%s", cmd.SpaceId.Identity(), err)
		}

		return err
	}

	if space.IsDisable() {
		e := xerrors.Errorf("spaceId:%s failed to create space failed, space is disable", space.Id.Identity())
		err = allerror.New(allerror.ErrorCodeSpaceAppCreateFailed, e.Error(), e)
		logrus.Errorf("spaceId：%s create space failed, err:%s", cmd.SpaceId.Identity(), err)
		return err
	}

	if space.Hardware.IsNpu() && !space.CompPowerAllocated {
		e := xerrors.Errorf("failed to create space failed, "+
			"spaceId:%s is npu but not allocate computility", space.Id.Identity())
		err = allerror.New(allerror.ErrorCodeSpaceAppCreateFailed, e.Error(), e)
		logrus.Errorf("spaceId：%s create space failed, err:%s", cmd.SpaceId.Identity(), err)
		return err
	}

	app, err := s.repo.FindBySpaceId(space.Id)
	if err == nil {
		if app.IsAppNotAllowToInit() {
			e := fmt.Errorf("spaceId:%s, not allow to init", space.Id.Identity())
			logrus.Errorf("create space app failed, err:%s", e)
			return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, e.Error(), e)
		}

		if err := s.repo.Remove(space.Id); err != nil {
			logrus.Errorf("spaceId:%s remove space app db failed, err:%s", space.Id.Identity(), err)
			return err
		}

		forceStopEvent := spacedomain.NewSpaceForceEvent(space.Id.Identity(), spacedomain.ForceTypeStop)
		if err := s.msg.SendSpaceAppForcePauseEvent(&forceStopEvent); err != nil {
			logrus.Errorf("spaceId:%s send force stop topic failed:%s", space.Id.Identity(), err)
			return err
		}
	}

	if space.Exception.Exception() != "" {
		e := xerrors.Errorf("spaceId:%s failed to create space failed, space has exception reason :%s",
			space.Id.Identity(), primitive.ExceptionMap[space.Exception.Exception()])
		err = allerror.New(allerror.ErrorCodeSpaceAppCreateFailed, e.Error(), e)
		logrus.Errorf("spaceId：%s create space failed, err:%s", cmd.SpaceId.Identity(), err)
		return err
	}

	v := domain.SpaceApp{
		Status:        appprimitive.AppStatusInit,
		SpaceAppIndex: *cmd,
	}
	if err := s.repo.Add(&v); err != nil {
		logrus.Errorf("spaceId:%s create space app db failed, err:%s", space.Id.Identity(), err)
		return err
	}
	e := domain.NewSpaceAppCreatedEvent(&v)
	if err := s.msg.SendSpaceAppCreatedEvent(&e); err != nil {
		logrus.Errorf("spaceId:%s send create topic failed, err:%v", space.Id.Identity(), err)
		return err
	}
	logrus.Infof("spaceId:%s create app successful", space.Id.Identity())
	return nil
}

func (s *spaceappInternalAppService) getSpaceApp(cmd CmdToCreateApp) (domain.SpaceApp, error) {
	space, err := s.spaceRepo.FindById(cmd.SpaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(xerrors.Errorf("space not found, err:%w", err))
		} else {
			err = xerrors.Errorf("failed to get space, err:%w", err)
		}
		logrus.Errorf("spaceId:%s get space failed, err:%s", cmd.SpaceId.Identity(), err)
		return domain.SpaceApp{}, err
	}

	if space.CommitId != cmd.CommitId {
		err = allerror.New(allerror.ErrorCodeSpaceCommitConflict, "commit conflict",
			xerrors.Errorf("spaceId:%s commit conflict", space.Id.Identity()))
		logrus.Errorf("spaceId:%s latest commitId:%s, old commitId:%s, err:%s",
			cmd.SpaceId.Identity(), space.CommitId, cmd.CommitId, err)
		return domain.SpaceApp{}, err
	}

	v, err := s.repo.FindBySpaceId(space.Id)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceAppNotFound(err)
		}
		logrus.Errorf("spaceId:%s get space app failed, err:%s", space.Id.Identity(), err)
		return domain.SpaceApp{}, err
	}
	return v, nil
}

// NotifyIsInvalid notifies change SpaceApp status.
func (s *spaceappInternalAppService) NotifyIsInvalid(cmd *CmdToNotifyFailedStatus) error {
	v, err := s.getSpaceApp(cmd.SpaceAppIndex)
	if err != nil {
		return err
	}
	if err := v.SetInvalid(cmd.Status, cmd.Reason); err != nil {
		logrus.Errorf("spaceId:%s set space app %s failed, err:%s",
			cmd.SpaceId.Identity(), cmd.Status.AppStatus(), err)
		return err
	}
	if err := s.repo.Save(&v); err != nil {
		logrus.Errorf("spaceId:%s save db failed", cmd.SpaceId.Identity())
		return err
	}
	logrus.Infof("spaceId:%s notify invalid successful", cmd.SpaceId.Identity())
	return nil
}

// NotifyIsBuilding notifies that the build process of a SpaceApp has started.
func (s *spaceappInternalAppService) NotifyIsBuilding(cmd *CmdToNotifyBuildIsStarted) error {
	v, err := s.getSpaceApp(cmd.SpaceAppIndex)
	if err != nil {
		return err
	}

	if err := v.StartBuilding(cmd.LogURL); err != nil {
		logrus.Errorf("spaceId:%s set space app building failed, err:%s", cmd.SpaceId.Identity(), err)
		return err
	}
	if err := s.repo.Save(&v); err != nil {
		logrus.Errorf("spaceId:%s save db failed", cmd.SpaceId.Identity())
		return err
	}
	logrus.Infof("spaceId:%s notify building successful", cmd.SpaceId.Identity())
	return nil
}

// NotifyIsBuildFailed notifies change SpaceApp status.
func (s *spaceappInternalAppService) NotifyIsBuildFailed(cmd *CmdToNotifyFailedStatus) error {
	v, err := s.getSpaceApp(cmd.SpaceAppIndex)
	if err != nil {
		return err
	}
	if err := v.SetBuildFailed(cmd.Status, cmd.Reason); err != nil {
		logrus.Errorf("spaceId:%s set space app %s failed, err:%s",
			cmd.SpaceId.Identity(), cmd.Status.AppStatus(), err)
		return err
	}
	if err := s.repo.Save(&v); err != nil {
		logrus.Errorf("spaceId:%s save db failed", cmd.SpaceId.Identity())
		return err
	}
	logrus.Infof("spaceId:%s notify build failed successful", cmd.SpaceId.Identity())
	return nil
}

// NotifyIsStarting notifies that the build process of a SpaceApp has finished.
func (s *spaceappInternalAppService) NotifyIsStarting(cmd *CmdToCreateApp) error {
	v, err := s.getSpaceApp(*cmd)
	if err != nil {
		return err
	}

	if err := v.SetStarting(); err != nil {
		logrus.Errorf("spaceId:%s set space app starting failed, err:%s", cmd.SpaceId.Identity(), err)
		return err
	}

	if err := s.repo.Save(&v); err != nil {
		logrus.Errorf("spaceId:%s save db failed", cmd.SpaceId.Identity())
		return err
	}
	logrus.Infof("spaceId:%s notify starting successful", cmd.SpaceId.Identity())
	return nil
}

// NotifyIsBuildFailed notifies change SpaceApp status.
func (s *spaceappInternalAppService) NotifyIsStartFailed(cmd *CmdToNotifyFailedStatus) error {
	v, err := s.getSpaceApp(cmd.SpaceAppIndex)
	if err != nil {
		return err
	}
	if err := v.SetStartFailed(cmd.Status, cmd.Reason); err != nil {
		logrus.Errorf("spaceId:%s set space app %s failed, err:%s",
			cmd.SpaceId.Identity(), cmd.Status.AppStatus(), err)
		return err
	}

	if err := s.repo.Save(&v); err != nil {
		logrus.Errorf("spaceId:%s save db failed", cmd.SpaceId.Identity())
		return err
	}
	logrus.Infof("spaceId:%s notify start failed successful", cmd.SpaceId.Identity())
	return nil
}

// NotifyIsServing notifies that a service of a SpaceApp has serving.
func (s *spaceappInternalAppService) NotifyIsServing(cmd *CmdToNotifyServiceIsStarted) error {
	v, err := s.getSpaceApp(cmd.SpaceAppIndex)
	if err != nil {
		return err
	}

	if err := v.StartServing(cmd.AppURL, cmd.LogURL); err != nil {
		logrus.Errorf("spaceId:%s set space app serving failed, err:%s", cmd.SpaceId.Identity(), err)
		return err
	}

	if err := s.repo.Save(&v); err != nil {
		logrus.Errorf("spaceId:%s save db failed", cmd.SpaceId.Identity())
		return err
	}
	logrus.Infof("spaceId:%s notify serving successful", cmd.SpaceId.Identity())

	return nil
}

// NotifyIsReStartFailed notifies change SpaceApp status.
func (s *spaceappInternalAppService) NotifyIsRestartFailed(cmd *CmdToNotifyFailedStatus) error {
	v, err := s.getSpaceApp(cmd.SpaceAppIndex)
	if err != nil {
		return err
	}
	if err := v.SetRestartFailed(cmd.Status, cmd.Reason); err != nil {
		logrus.Errorf("spaceId:%s set space app %s failed, err:%s",
			cmd.SpaceId.Identity(), cmd.Status.AppStatus(), err)
		return err
	}

	if err := s.repo.Save(&v); err != nil {
		logrus.Errorf("spaceId:%s save db failed", cmd.SpaceId.Identity())
		return err
	}
	logrus.Infof("spaceId:%s notify restart failed successful", cmd.SpaceId.Identity())
	return nil
}

// NotifyIsResumeFailed notifies change SpaceApp status.
func (s *spaceappInternalAppService) NotifyIsResumeFailed(cmd *CmdToNotifyFailedStatus) error {
	space, err := s.spaceRepo.FindById(cmd.SpaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(err)
			logrus.Errorf("spaceId:%s get space failed, err:%s", cmd.SpaceId.Identity(), err)
		}

		return err
	}

	app, err := s.repo.FindBySpaceId(space.Id)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceAppNotFound(err)
			logrus.Errorf("spaceId:%s get space app failed, err:%s", space.Id.Identity(), err)
		}

		return err
	}

	if err := app.SetResumeFailed(cmd.Status, cmd.Reason); err != nil {
		logrus.Errorf("spaceId:%s set space app %s failed, err:%s",
			space.Id.Identity(), cmd.Status.AppStatus(), err)
		return err
	}

	spaceCompCmd := spaceUserComputilityService{
		userName:    space.CreatedBy,
		space:       space,
		spaceRepo:   s.spaceRepo,
		computility: s.computility,
	}

	if err := spaceCompCmd.unbindSpaceCompQuota(); err != nil {
		err := fmt.Errorf("failed to release spaceId:%s comp quota, err: %w", space.Id.Identity(), err)
		logrus.Errorf("spaceId:%s set space %s status failed, err: %s",
			space.Id.Identity(), cmd.Status.AppStatus(), err)
		return err
	}

	if err := s.repo.Save(&app); err != nil {
		if err := spaceCompCmd.bindSpaceCompQuota(); err != nil {
			return err
		}
		err := fmt.Errorf("failed to save spaceId:%s db failed, err: %w", space.Id.Identity(), err)
		logrus.Errorf("spaceId:%s set space %s status failed, err:%s",
			space.Id.Identity(), cmd.Status.AppStatus(), err)
		return err
	}
	logrus.Infof("spaceId:%s notify resume failed successful", cmd.SpaceId.Identity())
	return nil
}

// PauseSpaceApp pause a SpaceApp in the spaceappAppService.
func (s *spaceappInternalAppService) ForcePauseSpaceApp(spaceId primitive.Identity) error {
	space, err := s.spaceRepo.FindById(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(err)
			logrus.Errorf("spaceId:%s get space failed, err:%s", spaceId.Identity(), err)
		}

		return err
	}
	spaceCompCmd := spaceUserComputilityService{
		userName:    space.CreatedBy,
		space:       space,
		spaceRepo:   s.spaceRepo,
		computility: s.computility,
	}
	app, err := s.repo.FindBySpaceId(space.Id)
	if err != nil {
		if err := spaceCompCmd.unbindSpaceCompQuota(); err != nil {
			logrus.Errorf("spaceId:%s release space comp quota failed, err:%s", space.Id.Identity(), err)
			return err
		}
		err = newSpaceAppNotFound(err)
		logrus.Errorf("spaceId:%s get space app failed, err:%s", space.Id.Identity(), err)
		return err
	}

	if app.Status.IsPaused() {
		logrus.Infof("spaceId:%s app is already paused", space.Id.Identity())
		return nil
	}
	app.Status = appprimitive.AppStatusPaused

	if err := spaceCompCmd.unbindSpaceCompQuota(); err != nil {
		e := fmt.Errorf("failed to release spaceId:%s comp quota, err: %w", space.Id.Identity(), err)
		err = allerror.New(allerror.ErrorCodeSpaceAppPauseFailed, e.Error(), e)
		logrus.Errorf("spaceId:%s space unbind quota failed:%s", space.Id.Identity(), err)
		return err
	}

	if err := s.repo.Save(&app); err != nil {
		if err := spaceCompCmd.bindSpaceCompQuota(); err != nil {
			logrus.Errorf("spaceId:%s space bind quota failed:%s", space.Id.Identity(), err)
			return err
		}
		e := fmt.Errorf("failed to save spaceId:%s db failed, err: %w", space.Id.Identity(), err)
		err = allerror.New(allerror.ErrorCodeSpaceAppPauseFailed, e.Error(), e)
		logrus.Errorf("spaceId:%s save db failed:%s", space.Id.Identity(), err)
		return err
	}
	e := spacedomain.NewSpaceForceEvent(space.Id.Identity(), spacedomain.ForceTypePause)
	if err := s.msg.SendSpaceAppForcePauseEvent(&e); err != nil {
		logrus.Errorf("spaceId:%s send force pause topic failed:%s", space.Id.Identity(), err)
		return err
	}
	logrus.Infof("spaceId:%s force paused app successful", space.Id.Identity())
	return nil
}

func (s *spaceappInternalAppService) PauseSpaceApp(spaceId primitive.Identity) error {
	space, err := s.spaceRepo.FindById(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(err)
			logrus.Errorf("spaceId:%s get space failed, err:%s", spaceId.Identity(), err)
		}

		return err
	}
	app, err := s.repo.FindBySpaceId(space.Id)
	if err != nil {
		err = newSpaceAppNotFound(err)
		logrus.Errorf("spaceId:%s get space app failed, err:%s", space.Id.Identity(), err)
		return err
	}

	if err := app.PauseService(); err != nil {
		logrus.Errorf("spaceId:%s paused service failed", space.Id.Identity())
		return err
	}

	spaceCompCmd := spaceUserComputilityService{
		userName:    space.CreatedBy,
		space:       space,
		spaceRepo:   s.spaceRepo,
		computility: s.computility,
	}

	if err := spaceCompCmd.unbindSpaceCompQuota(); err != nil {
		e := fmt.Errorf("failed to release spaceId:%s comp quota, err: %w", space.Id.Identity(), err)
		err = allerror.New(allerror.ErrorCodeSpaceAppPauseFailed, e.Error(), e)
		logrus.Errorf("spaceId:%s space unbind quota failed:%s", space.Id.Identity(), err)
		return err
	}

	if err := s.repo.Save(&app); err != nil {
		if err := spaceCompCmd.bindSpaceCompQuota(); err != nil {
			return err
		}
		e := fmt.Errorf("failed to save spaceId:%s db failed, err: %w", space.Id.Identity(), err)
		err = allerror.New(allerror.ErrorCodeSpaceAppPauseFailed, e.Error(), e)
		logrus.Errorf("spaceId:%s space bind quota failed:%s", space.Id.Identity(), err)
		return err
	}
	e := domain.NewSpaceAppPauseEvent(&domain.SpaceAppIndex{
		SpaceId: app.SpaceId,
	})
	if err := s.msg.SendSpaceAppPauseEvent(&e); err != nil {
		logrus.Errorf("spaceId:%s send pause topic failed:%s", space.Id.Identity(), err)
		return err
	}
	logrus.Infof("spaceId:%s paused app successful", space.Id.Identity())
	return nil
}
