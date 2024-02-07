package app

import (
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/space-app/domain"
	"github.com/openmerlin/merlin-server/space-app/domain/message"
	appprimitive "github.com/openmerlin/merlin-server/space-app/domain/primitive"
	"github.com/openmerlin/merlin-server/space-app/domain/repository"
	spaceapp "github.com/openmerlin/merlin-server/space/app"
)

var (
	errSpaceAppNotFound = allerror.NewNotFound(allerror.ErrorCodeSpaceAppNotFound, "not found")
)

type SpaceAppInternalAppService interface {
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
	msg   message.SpaceAppMessage
	repo  repository.Repository
	space spaceapp.SpaceAppService
}

// Create
func (s *spaceappInternalAppService) Create(cmd *CmdToCreateApp) error {
	// TODO check if it is the space repo
	// TODO check if it is the newest commit

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

// NotifyBuildIsStarted
func (s *spaceappInternalAppService) NotifyBuildIsStarted(cmd *CmdToNotifyBuildIsStarted) error {
	v, err := s.repo.Find(&cmd.SpaceAppIndex)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errSpaceAppNotFound
		}

		return err
	}

	if err := v.StartBuilding(cmd.LogURL); err != nil {
		return err
	}

	return s.repo.Save(&v)
}

// NotifyBuildIsDone
func (s *spaceappInternalAppService) NotifyBuildIsDone(cmd *CmdToNotifyBuildIsDone) error {
	v, err := s.repo.Find(&cmd.SpaceAppIndex)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errSpaceAppNotFound
		}

		return err
	}

	if err := v.SetBuildIsDone(cmd.Success); err != nil {
		return err
	}

	return s.repo.Save(&v)
}

// NotifyServiceIsStarted
func (s *spaceappInternalAppService) NotifyServiceIsStarted(cmd *CmdToNotifyServiceIsStarted) error {
	v, err := s.repo.Find(&cmd.SpaceAppIndex)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errSpaceAppNotFound
		}

		return err
	}

	if err := v.StartService(cmd.AppURL, cmd.LogURL); err != nil {
		return err
	}

	return s.repo.Save(&v)
}
