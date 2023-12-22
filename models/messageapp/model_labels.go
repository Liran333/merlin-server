package messageapp

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/models/domain/repository"
)

type ModelAppService interface {
	ResetLabels(modelId primitive.Identity, cmd *CmdToResetLabels) error
}

func NewModelAppService(repoAdapter repository.ModelLabelsRepoAdapter) ModelAppService {
	return &modelAppService{
		repoAdapter: repoAdapter,
	}
}

type modelAppService struct {
	repoAdapter repository.ModelLabelsRepoAdapter
}

func (s *modelAppService) ResetLabels(modelId primitive.Identity, cmd *CmdToResetLabels) error {
	return s.repoAdapter.Save(modelId, cmd)
}
