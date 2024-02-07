package app

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/models/domain/repository"
)

type ModelInternalAppService interface {
	ResetLabels(primitive.Identity, *CmdToResetLabels) error
}

func NewModelInternalAppService(repoAdapter repository.ModelLabelsRepoAdapter) ModelInternalAppService {
	return &modelInternalAppService{
		repoAdapter: repoAdapter,
	}
}

type modelInternalAppService struct {
	repoAdapter repository.ModelLabelsRepoAdapter
}

func (s *modelInternalAppService) ResetLabels(modelId primitive.Identity, cmd *CmdToResetLabels) error {
	err := s.repoAdapter.Save(modelId, cmd)

	if err != nil && commonrepo.IsErrorResourceNotExists(err) {
		err = errorModelNotFound
	}

	return err
}
