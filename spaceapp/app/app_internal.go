/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/spaceapp/domain"
	"github.com/openmerlin/merlin-server/spaceapp/domain/message"
	appprimitive "github.com/openmerlin/merlin-server/spaceapp/domain/primitive"
	"github.com/openmerlin/merlin-server/spaceapp/domain/repository"
)

var (
	errorSpaceAppNotFound = allerror.NewNotFound(allerror.ErrorCodeSpaceAppNotFound, "not found")
)

// SpaceappInternalAppService is an interface that defines the methods for creating and managing a SpaceApp.
type SpaceappInternalAppService interface {
	Create(cmd *CmdToCreateApp) error
	NotifyBuildIsStarted(cmd *CmdToNotifyBuildIsStarted) error
	NotifyBuildIsDone(cmd *CmdToNotifyBuildIsDone) error
	NotifyServiceIsStarted(cmd *CmdToNotifyServiceIsStarted) error
}

func NewSpaceappInternalAppService(
	msg message.SpaceAppMessage,
	repo repository.Repository,
) *spaceappInternalAppService {
	return &spaceappInternalAppService{
		msg:  msg,
		repo: repo,
	}
}

// spaceappInternalAppService
type spaceappInternalAppService struct {
	msg  message.SpaceAppMessage
	repo repository.Repository
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
			err = errorSpaceAppNotFound
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
			err = errorSpaceAppNotFound
		}

		return err
	}

	if err := v.SetBuildIsDone(cmd.Success); err != nil {
		return err
	}

	return s.repo.Save(&v)
}

// NotifyServiceIsStarted notifies that a service of a SpaceApp has started.
func (s *spaceappInternalAppService) NotifyServiceIsStarted(cmd *CmdToNotifyServiceIsStarted) error {
	v, err := s.repo.Find(&cmd.SpaceAppIndex)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errorSpaceAppNotFound
		}

		return err
	}

	if err := v.StartService(cmd.AppURL, cmd.LogURL); err != nil {
		return err
	}

	return s.repo.Save(&v)
}
