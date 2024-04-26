/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package app provides functionality for the application.
package app

import (
	"fmt"

	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/models/domain"
	"github.com/openmerlin/merlin-server/models/domain/repository"
	"github.com/openmerlin/merlin-server/utils"
)

// ModelInternalAppService is an interface for the internal model application service.
type ModelInternalAppService interface {
	ResetLabels(primitive.Identity, *CmdToResetLabels) error
	UpdateUseInOpenmind(primitive.Identity, string) error
	GetById(modelId primitive.Identity) (ModelDTO, error)
	GetByNames([]*domain.ModelIndex) []primitive.Identity
	UpdateStatistics(primitive.Identity, *CmdToUpdateStatistics) error
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
		err = allerror.NewNotFound(allerror.ErrorCodeModelNotFound, "not found", fmt.Errorf("%s not found, %w", modelId, err))
	}

	return err
}

// GetById retrieves a model by id.
func (s *modelInternalAppService) GetById(modelId primitive.Identity) (ModelDTO, error) {
	model, err := s.modelAdapter.FindById(modelId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeModelNotFound, "not found", fmt.Errorf("%s not found, %w", modelId, err))
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

// UpdateStatistics updates the statistics of a model.
func (s *modelInternalAppService) UpdateStatistics(modelId primitive.Identity, cmd *CmdToUpdateStatistics) error {
	model, err := s.modelAdapter.FindById(modelId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeModelNotFound, "not found", fmt.Errorf("%s not found, err: %w", modelId.Identity(), err))
		}
		return err
	}

	model.DownloadCount = cmd.DownloadCount
	model.UpdatedAt = utils.Now()

	return s.modelAdapter.Save(&model)
}

// UpdateUseInOpenmind set the use in openmind tag of a model.
func (s *modelInternalAppService) UpdateUseInOpenmind(modelId primitive.Identity, cmd string) error {
	model, err := s.modelAdapter.FindById(modelId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeModelNotFound, "not found", fmt.Errorf("%s not found, %w", modelId, err))
		}
		return err
	}

	model.UseInOpenmind = cmd
	model.UpdatedAt = utils.Now()

	if err := s.modelAdapter.Save(&model); err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeModelNotFound, "not found", fmt.Errorf("%s not found, %w", modelId, err))
		}

		if commonrepo.IsErrorConcurrentUpdating(err) {
			err = allerror.New(allerror.ErrorCodeConcurrentUpdating, "concurrent updating", fmt.Errorf("failed to update use_in_openmind, %w", err))
		}

		return err
	}

	return nil
}
