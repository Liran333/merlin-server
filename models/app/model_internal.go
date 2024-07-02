/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package app provides functionality for the application.
package app

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"

	coderepo "github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/models/domain"
	"github.com/openmerlin/merlin-server/models/domain/repository"
)

// ModelInternalAppService is an interface for the internal model application service.
type ModelInternalAppService interface {
	ResetLabels(primitive.Identity, *CmdToResetLabels) error
	UpdateUseInOpenmind(primitive.Identity, string) error
	GetById(modelId primitive.Identity) (ModelDTO, error)
	GetByNames([]*domain.ModelIndex) ([]primitive.Identity, error)
	UpdateStatistics(primitive.Identity, *CmdToUpdateStatistics) error
	SaveDeploy(domain.ModelIndex, CmdToDeploy) error
}

// NewModelInternalAppService creates a new instance of the internal model application service.
func NewModelInternalAppService(
	repoAdapter repository.ModelLabelsRepoAdapter,
	mAdapter repository.ModelRepositoryAdapter,
	deploy repository.ModelDeployRepoAdapter,
) ModelInternalAppService {
	return &modelInternalAppService{
		repoAdapter:   repoAdapter,
		modelAdapter:  mAdapter,
		deployAdapter: deploy,
	}
}

type modelInternalAppService struct {
	repoAdapter   repository.ModelLabelsRepoAdapter
	modelAdapter  repository.ModelRepositoryAdapter
	deployAdapter repository.ModelDeployRepoAdapter
}

// ResetLabels resets the labels of a model.
func (s *modelInternalAppService) ResetLabels(modelId primitive.Identity, cmd *CmdToResetLabels) error {
	err := s.repoAdapter.Save(modelId, cmd)

	if err != nil && commonrepo.IsErrorResourceNotExists(err) {
		err = allerror.NewNotFound(allerror.ErrorCodeModelNotFound, "not found",
			fmt.Errorf("%s not found, %w", modelId, err))
	}

	return err
}

// GetById retrieves a model by id.
func (s *modelInternalAppService) GetById(modelId primitive.Identity) (ModelDTO, error) {
	model, err := s.modelAdapter.FindById(modelId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeModelNotFound, "not found",
				fmt.Errorf("%s not found, %w", modelId, err))
		}
		return ModelDTO{}, err
	}

	return toModelDTO(&model), nil
}

// GetByNames retrieves ids of models by names.
func (s *modelInternalAppService) GetByNames(modelsIndex []*domain.ModelIndex) ([]primitive.Identity, error) {
	var dtos []primitive.Identity
	var resErr error

	for _, index := range modelsIndex {
		model, err := s.modelAdapter.FindByName(index)
		if err != nil {
			errInfo := fmt.Errorf("related model %v was not found", index.Name.MSDName())
			logrus.Errorf("%s, do not allow to remove exception", errInfo)
			resErr = allerror.NewNotFound(allerror.ErrorCodeModelNotFound, errInfo.Error(), errInfo)
			continue
		}
		if model.IsDisable() {
			errInfo := fmt.Errorf("related model %v was disable", model.Name.MSDName())
			logrus.Errorf("%s, do not allow to remove exception", errInfo)
			resErr = allerror.NewResourceDisabled(allerror.ErrorCodeResourceDisabled, errInfo.Error(), errInfo)
		}
		dtos = append(dtos, model.Id)
	}

	return dtos, resErr
}

// UpdateStatistics updates the statistics of a model.
func (s *modelInternalAppService) UpdateStatistics(modelId primitive.Identity, cmd *CmdToUpdateStatistics) error {
	_, err := s.modelAdapter.FindById(modelId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeModelNotFound, "not found",
				fmt.Errorf("%s not found, err: %w", modelId.Identity(), err))
		}
		return err
	}
	newModel := domain.Model{
		CodeRepo: coderepo.CodeRepo{
			Id: modelId,
		},
		DownloadCount: cmd.DownloadCount,
	}
	return s.modelAdapter.InternalSaveStatistic(&newModel)
}

// UpdateUseInOpenmind set the use in openmind tag of a model.
func (s *modelInternalAppService) UpdateUseInOpenmind(modelId primitive.Identity, cmd string) error {
	_, err := s.modelAdapter.FindById(modelId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeModelNotFound, "not found",
				fmt.Errorf("%s not found, %w", modelId, err))
		}
		return err
	}
	newModel := domain.Model{
		CodeRepo: coderepo.CodeRepo{
			Id: modelId,
		},
		UseInOpenmind: cmd,
	}

	if err := s.modelAdapter.InternalSaveUseInOpenmind(&newModel); err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeModelNotFound, "not found",
				fmt.Errorf("%s not found, %w", modelId, err))
		}

		if commonrepo.IsErrorConcurrentUpdating(err) {
			err = allerror.New(allerror.ErrorCodeConcurrentUpdating, "concurrent updating",
				fmt.Errorf("failed to update use_in_openmind, %w", err))
		}

		return err
	}

	return nil
}

func (s *modelInternalAppService) SaveDeploy(index domain.ModelIndex, deploy CmdToDeploy) error {
	if err := s.deployAdapter.DeleteByOwnerName(index); err != nil {
		return xerrors.Errorf("delete deploy of [%s/%s] error:%w",
			index.Owner.Account(), index.Name.MSDName(), err)
	}

	return s.deployAdapter.Create(index, deploy)
}
