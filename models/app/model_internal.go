/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/models/domain/repository"
)

// ModelInternalAppService is an interface for the internal model application service.
type ModelInternalAppService interface {
	ResetLabels(primitive.Identity, *CmdToResetLabels) error
}

// NewModelInternalAppService creates a new instance of the internal model application service.
func NewModelInternalAppService(repoAdapter repository.ModelLabelsRepoAdapter) ModelInternalAppService {
	return &modelInternalAppService{
		repoAdapter: repoAdapter,
	}
}

type modelInternalAppService struct {
	repoAdapter repository.ModelLabelsRepoAdapter
}

// ResetLabels resets the labels of a model.
func (s *modelInternalAppService) ResetLabels(modelId primitive.Identity, cmd *CmdToResetLabels) error {
	err := s.repoAdapter.Save(modelId, cmd)

	if err != nil && commonrepo.IsErrorResourceNotExists(err) {
		err = errorModelNotFound
	}

	return err
}
