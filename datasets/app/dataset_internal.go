/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package app provides functionality for the application.
package app

import (
	"golang.org/x/xerrors"

	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/datasets/domain"
	"github.com/openmerlin/merlin-server/datasets/domain/repository"
	"github.com/openmerlin/merlin-server/utils"
)

// DatasetInternalAppService is an interface for the internal dataset application service.
type DatasetInternalAppService interface {
	ResetLabels(primitive.Identity, *CmdToResetLabels) error
	GetById(datasetId primitive.Identity) (DatasetDTO, error)
	GetByNames([]*domain.DatasetIndex) []primitive.Identity
	UpdateStatistics(primitive.Identity, *CmdToUpdateStatistics) error
}

// NewDatasetInternalAppService creates a new instance of the internal dataset application service.
func NewDatasetInternalAppService(
	repoAdapter repository.DatasetLabelsRepoAdapter,
	mAdapter repository.DatasetRepositoryAdapter,
) DatasetInternalAppService {
	return &datasetInternalAppService{
		repoAdapter:    repoAdapter,
		datasetAdapter: mAdapter,
	}
}

type datasetInternalAppService struct {
	repoAdapter    repository.DatasetLabelsRepoAdapter
	datasetAdapter repository.DatasetRepositoryAdapter
}

// ResetLabels resets the labels of a dataset.
func (s *datasetInternalAppService) ResetLabels(datasetId primitive.Identity, cmd *CmdToResetLabels) error {
	err := s.repoAdapter.Save(datasetId, cmd)

	if err != nil && commonrepo.IsErrorResourceNotExists(err) {
		err = allerror.NewNotFound(allerror.ErrorCodeDatasetNotFound, "not found",
			xerrors.Errorf("%s not found, %w", datasetId, err))
	}

	return err
}

// GetById retrieves a dataset by id.
func (s *datasetInternalAppService) GetById(datasetId primitive.Identity) (DatasetDTO, error) {
	dataset, err := s.datasetAdapter.FindById(datasetId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeDatasetNotFound, "not found",
				xerrors.Errorf("%s not found, %w", datasetId, err))
		}
		return DatasetDTO{}, err
	}

	return toDatasetDTO(&dataset), nil
}

// GetByNames retrieves ids of datasets by names.
func (s *datasetInternalAppService) GetByNames(datasetsIndex []*domain.DatasetIndex) []primitive.Identity {
	var dtos []primitive.Identity

	for _, index := range datasetsIndex {
		dataset, err := s.datasetAdapter.FindByName(index)
		if err != nil {
			continue
		}

		dtos = append(dtos, dataset.Id)
	}

	return dtos
}

// UpdateStatistics updates the statistics of a dataset.
func (s *datasetInternalAppService) UpdateStatistics(datasetId primitive.Identity, cmd *CmdToUpdateStatistics) error {
	dataset, err := s.datasetAdapter.FindById(datasetId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeDatasetNotFound, "not found",
				xerrors.Errorf("%s not found, err: %w", datasetId.Identity(), err))
		}
		return xerrors.Errorf("failed to update statistics, %w", err)
	}

	dataset.DownloadCount = cmd.DownloadCount
	dataset.UpdatedAt = utils.Now()

	return s.datasetAdapter.Save(&dataset)
}
