package app

import (
	"github.com/openmerlin/merlin-server/space-app/domain"
	"github.com/openmerlin/merlin-server/space-app/domain/message"
	"github.com/openmerlin/merlin-server/space-app/domain/repository"
	spaceapp "github.com/openmerlin/merlin-server/space/app"
)

type SpaceAppInternalAppService interface {
	Create(cmd *CmdToCreateApp) error
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

func (s *spaceappInternalAppService) Create(cmd *CmdToCreateApp) error {
	// TODO check if it is the space repo
	// TODO check if it is the newest commit

	v := cmd.toApp()

	if err := s.repo.Add(&v); err != nil {
		return err
	}

	e := domain.NewSpaceAppCreatedEvent(&v)

	return s.msg.SendSpaceAppCreatedEvent(&e)
}
