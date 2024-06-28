/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"context"
	"fmt"
	"time"

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
	Create(ctx context.Context, cmd *CmdToCreateApp) error

	NotifyIsBuilding(ctx context.Context, cmd *CmdToNotifyBuildIsStarted) error
	NotifyIsBuildFailed(ctx context.Context, cmd *CmdToNotifyFailedStatus) error
	NotifyStarting(ctx context.Context, cmd *CmdToNotifyStarting) error
	NotifyIsStartFailed(ctx context.Context, cmd *CmdToNotifyFailedStatus) error
	NotifyIsServing(ctx context.Context, cmd *CmdToNotifyServiceIsStarted) error
	NotifyIsRestartFailed(ctx context.Context, cmd *CmdToNotifyFailedStatus) error
	NotifyIsResumeFailed(ctx context.Context, cmd *CmdToNotifyFailedStatus) error

	ForcePauseSpaceApp(context.Context, primitive.Identity) error
	PauseSpaceApp(context.Context, primitive.Identity) error

	SleepSpaceApp(context.Context, *CmdToSleepSpaceApp) error
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
func (s *spaceappInternalAppService) Create(ctx context.Context, cmd *CmdToCreateApp) error {
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

	app, err := s.repo.FindBySpaceId(ctx, space.Id)
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
		time.Sleep(time.Second)
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

func (s *spaceappInternalAppService) getSpaceApp(ctx context.Context, cmd CmdToCreateApp) (domain.SpaceApp, error) {
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

	v, err := s.repo.FindBySpaceId(ctx, space.Id)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceAppNotFound(err)
		}
		logrus.Errorf("spaceId:%s get space app failed, err:%s", space.Id.Identity(), err)
		return domain.SpaceApp{}, err
	}
	return v, nil
}

// NotifyIsBuilding notifies that the build process of a SpaceApp has started.
func (s *spaceappInternalAppService) NotifyIsBuilding(ctx context.Context, cmd *CmdToNotifyBuildIsStarted) error {
	v, err := s.getSpaceApp(ctx, cmd.SpaceAppIndex)
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
func (s *spaceappInternalAppService) NotifyIsBuildFailed(ctx context.Context, cmd *CmdToNotifyFailedStatus) error {
	v, err := s.getSpaceApp(ctx, cmd.SpaceAppIndex)
	if err != nil {
		return err
	}

	if err := v.SetBuildFailed(cmd.Status, cmd.Reason); err != nil {
		logrus.Errorf("spaceId:%s set space app %s failed, err:%s",
			cmd.SpaceId.Identity(), cmd.Status.AppStatus(), err)
		return err
	}

	if err := s.repo.SaveWithBuildLog(&v, &domain.SpaceAppBuildLog{
		Logs: cmd.Logs,
	}); err != nil {
		logrus.Errorf("spaceId:%s save with build log db failed, err:%s", cmd.SpaceId.Identity(), err)
		return err
	}

	logrus.Infof("spaceId:%s notify build failed successful, save build logs:%d",
		cmd.SpaceId.Identity(), len(cmd.Logs))
	return nil
}

// NotifyStarting notifies that the build process of a SpaceApp has finished.
func (s *spaceappInternalAppService) NotifyStarting(ctx context.Context, cmd *CmdToNotifyStarting) error {
	v, err := s.getSpaceApp(ctx, cmd.SpaceAppIndex)
	if err != nil {
		return err
	}

	if err := v.SetStarting(); err != nil {
		logrus.Errorf("spaceId:%s set space app starting failed, err:%s", cmd.SpaceId.Identity(), err)
		return err
	}

	if err := s.repo.SaveWithBuildLog(&v, &domain.SpaceAppBuildLog{
		Logs: cmd.Logs,
	}); err != nil {
		logrus.Errorf("spaceId:%s save with build log db failed, err:%s", cmd.SpaceId.Identity(), err)
		return err
	}

	logrus.Infof("spaceId:%s notify starting successful, save build logs:%d",
		cmd.SpaceId.Identity(), len(cmd.Logs))
	return nil
}

// NotifyIsBuildFailed notifies change SpaceApp status.
func (s *spaceappInternalAppService) NotifyIsStartFailed(ctx context.Context, cmd *CmdToNotifyFailedStatus) error {
	v, err := s.getSpaceApp(ctx, cmd.SpaceAppIndex)
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
func (s *spaceappInternalAppService) NotifyIsServing(ctx context.Context, cmd *CmdToNotifyServiceIsStarted) error {
	v, err := s.getSpaceApp(ctx, cmd.SpaceAppIndex)
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
func (s *spaceappInternalAppService) NotifyIsRestartFailed(ctx context.Context, cmd *CmdToNotifyFailedStatus) error {
	v, err := s.getSpaceApp(ctx, cmd.SpaceAppIndex)
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
func (s *spaceappInternalAppService) NotifyIsResumeFailed(ctx context.Context, cmd *CmdToNotifyFailedStatus) error {
	space, err := s.spaceRepo.FindById(cmd.SpaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(err)
			logrus.Errorf("spaceId:%s get space failed, err:%s", cmd.SpaceId.Identity(), err)
		}

		return err
	}

	app, err := s.repo.FindBySpaceId(ctx, space.Id)
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
func (s *spaceappInternalAppService) ForcePauseSpaceApp(ctx context.Context, spaceId primitive.Identity) error {
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
	app, err := s.repo.FindBySpaceId(ctx, space.Id)
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

func (s *spaceappInternalAppService) PauseSpaceApp(ctx context.Context, spaceId primitive.Identity) error {
	space, err := s.spaceRepo.FindById(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(err)
			logrus.Errorf("spaceId:%s get space failed, err:%s", spaceId.Identity(), err)
		}

		return err
	}
	app, err := s.repo.FindBySpaceId(ctx, space.Id)
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

func (s *spaceappInternalAppService) SleepSpaceApp(ctx context.Context, cmd *CmdToSleepSpaceApp) error {
	space, err := s.spaceRepo.FindById(cmd.SpaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(xerrors.Errorf("not found, err: %w", err))
		} else {
			err = xerrors.Errorf("find space by id failed, err: %w", err)
		}
		logrus.Errorf("spaceId:%s get space failed, err:%s", cmd.SpaceId.Identity(), err)
		return err
	}

	app, err := s.repo.FindBySpaceId(ctx, space.Id)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceAppNotFound(xerrors.Errorf("space app not found, err:%w", err))
		} else {
			err = xerrors.Errorf("failed to get space app, err:%w", err)
		}
		logrus.Errorf("spaceId:%s get space app failed, err:%s", space.Id.Identity(), err)
		return err
	}

	if app.CommitId != cmd.CommitId {
		err = xerrors.Errorf("sleep commit id:%s is not equal app:%s, failed to sleep space app",
			cmd.CommitId, app.CommitId)
		logrus.Errorf("spaceId:%s sleep failed, err:%s", space.Id.Identity(), err)
		return err
	}

	if app.Status.IsSleeping() {
		logrus.Infof("spaceId:%s is sleeping", space.Id.Identity())
		return s.sendSpaceAppSleepMsg(app)
	}

	if err := app.SleepService(); err != nil {
		logrus.Errorf("spaceId:%s sleep service failed", space.Id.Identity())
		return err
	}

	if err := s.repo.Save(&app); err != nil {
		e := fmt.Errorf("failed to save spaceId:%s db failed, err: %w", space.Id.Identity(), err)
		err = allerror.New(allerror.ErrorCodeSpaceAppSleepFailed, e.Error(), e)
		logrus.Errorf("spaceId:%s space sleep app failed:%s", space.Id.Identity(), err)
		return err
	}
	return s.sendSpaceAppSleepMsg(app)
}

func (s *spaceappInternalAppService) sendSpaceAppSleepMsg(app domain.SpaceApp) error {
	e := domain.NewSpaceAppSleepEvent(&domain.SpaceAppIndex{
		SpaceId: app.SpaceId,
	})
	if err := s.msg.SendSpaceAppSleepEvent(&e); err != nil {
		logrus.Errorf("spaceId:%s send sleep topic failed:%s", app.SpaceId.Identity(), err)
		return err
	}
	logrus.Infof("spaceId:%s sleep app successful", app.SpaceId.Identity())
	return nil
}
