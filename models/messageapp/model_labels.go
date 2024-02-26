/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package messageapp

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/models/domain/repository"
)

// ModelAppService is an interface that defines the ResetLabels method.
type ModelAppService interface {
	ResetLabels(primitive.Identity, *CmdToResetLabels) error
}

// NewModelAppService creates a new instance of modelAppService with the given repoAdapter.
func NewModelAppService(repoAdapter repository.ModelLabelsRepoAdapter) ModelAppService {
	return &modelAppService{
		repoAdapter: repoAdapter,
	}
}

type modelAppService struct {
	repoAdapter repository.ModelLabelsRepoAdapter
}

// ResetLabels resets the labels for the given modelId using the provided cmd.
func (s *modelAppService) ResetLabels(modelId primitive.Identity, cmd *CmdToResetLabels) error {
	return s.repoAdapter.Save(modelId, cmd)
}
