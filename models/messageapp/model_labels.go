package messageapp

import (
	"github.com/openmerlin/merlin-server/models/domain"
	"github.com/openmerlin/merlin-server/models/domain/repository"
)

type ModelAppService interface {
	ResetLabels(*domain.ModelIndex, *CmdToResetLabels) error
}

func NewModelAppService(repoAdapter repository.ModelLabelsRepoAdapter) ModelAppService {
	return &modelAppService{
		repoAdapter: repoAdapter,
	}
}

type modelAppService struct {
	repoAdapter repository.ModelLabelsRepoAdapter
}

func (s *modelAppService) ResetLabels(index *domain.ModelIndex, cmd *CmdToResetLabels) error {
	return s.repoAdapter.Save(index, cmd)
}
