package messageapp

import (
	"github.com/openmerlin/merlin-server/space/domain"
	"github.com/openmerlin/merlin-server/space/domain/repository"
)

type SpaceAppService interface {
	ResetLabels(*domain.SpaceIndex, *CmdToResetLabels) error
}

func NewSpaceAppService(repoAdapter repository.SpaceLabelsRepoAdapter) SpaceAppService {
	return &spaceAppService{
		repoAdapter: repoAdapter,
	}
}

type spaceAppService struct {
	repoAdapter repository.SpaceLabelsRepoAdapter
}

func (s *spaceAppService) ResetLabels(index *domain.SpaceIndex, cmd *CmdToResetLabels) error {
	return s.repoAdapter.Save(index, cmd)
}
