/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/models/domain"
	"github.com/openmerlin/merlin-server/models/domain/repository"
)

// ModelInternalAppService is an interface for the internal model application service.
type ModelInternalAppService interface {
	ResetLabels(primitive.Identity, *CmdToResetLabels) error
	GetById(modelId primitive.Identity) (ModelDTO, error)
	GetByNames([]*domain.ModelIndex) []primitive.Identity
}

// NewModelInternalAppService creates a new instance of the internal model application service.
func NewModelInternalAppService(
	repoAdapter repository.ModelLabelsRepoAdapter,
	mAdapter repository.ModelRepositoryAdapter,
) ModelInternalAppService {
	return &modelInternalAppService{
		repoAdapter:  repoAdapter,
		modelAdapter: mAdapter,
	}
}

type modelInternalAppService struct {
	repoAdapter  repository.ModelLabelsRepoAdapter
	modelAdapter repository.ModelRepositoryAdapter
}

// ResetLabels resets the labels of a model.
func (s *modelInternalAppService) ResetLabels(modelId primitive.Identity, cmd *CmdToResetLabels) error {
	err := s.repoAdapter.Save(modelId, cmd)

	if err != nil && commonrepo.IsErrorResourceNotExists(err) {
		err = errorModelNotFound
	}

	return err
}

// GetById retrieves a model by id.
func (s *modelInternalAppService) GetById(modelId primitive.Identity) (ModelDTO, error) {
	model, err := s.modelAdapter.FindById(modelId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errorModelNotFound
		}
		return ModelDTO{}, err
	}

	return toModelDTO(&model), nil
}

// GetByNames retrieves ids of models by names.
func (s *modelInternalAppService) GetByNames(modelsIndex []*domain.ModelIndex) []primitive.Identity {
	var dtos []primitive.Identity

	for _, index := range modelsIndex {
		model, err := s.modelAdapter.FindByName(index)
		if err != nil {
			continue
		}

		dtos = append(dtos, model.Id)
	}

	return dtos
}
