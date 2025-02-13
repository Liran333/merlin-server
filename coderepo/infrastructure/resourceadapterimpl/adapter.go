/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package resourceadapterimpl provides an implementation of the resource adapter interface.
package resourceadapterimpl

import (
	"errors"

	"github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/coderepo/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	datasetsdomain "github.com/openmerlin/merlin-server/datasets/domain"
	datasetrepo "github.com/openmerlin/merlin-server/datasets/domain/repository"
	modeldomain "github.com/openmerlin/merlin-server/models/domain"
	modelrepo "github.com/openmerlin/merlin-server/models/domain/repository"
	spacedomain "github.com/openmerlin/merlin-server/space/domain"
	spacerepo "github.com/openmerlin/merlin-server/space/domain/repository"
)

// NewResourceAdapterImpl creates a new instance of the resourceAdapterImpl.
func NewResourceAdapterImpl(
	model modelrepo.ModelRepositoryAdapter,
	dataset datasetrepo.DatasetRepositoryAdapter,
	space spacerepo.SpaceRepositoryAdapter,
) *resourceAdapterImpl {
	return &resourceAdapterImpl{
		model:   model,
		dataset: dataset,
		space:   space,
	}
}

// resourceAdapterImpl
type resourceAdapterImpl struct {
	model   modelrepo.ModelRepositoryAdapter
	dataset datasetrepo.DatasetRepositoryAdapter
	space   spacerepo.SpaceRepositoryAdapter
}

// GetByName retrieves a resource by name.
func (adapter *resourceAdapterImpl) GetByName(index *domain.CodeRepoIndex) (domain.Resource, error) {
	r, err := adapter.model.FindByName(index)
	if err == nil {
		return &r, nil
	}
	if !commonrepo.IsErrorResourceNotExists(err) {
		return nil, err
	}

	dr, err := adapter.dataset.FindByName(index)
	if err == nil {
		return &dr, nil
	}
	if !commonrepo.IsErrorResourceNotExists(err) {
		return nil, err
	}

	space, err := adapter.space.FindByName(index)

	return &space, err
}

// GetByType retrieves a resource by type and name.
func (adapter *resourceAdapterImpl) GetByType(t primitive.RepoType,
	index *domain.CodeRepoIndex) (domain.Resource, error) {
	if t.IsModel() {
		r, err := adapter.model.FindByName(index)

		return &r, err
	}

	if t.IsSpace() {
		r, err := adapter.space.FindByName(index)

		return &r, err

	}

	if t.IsDataset() {
		r, err := adapter.dataset.FindByName(index)

		return &r, err
	}

	return nil, commonrepo.NewErrorResourceNotExists(errors.New("unknown repo type"))
}

// GetByIndex retrieves a resource by index.
func (adapter *resourceAdapterImpl) GetByIndex(index primitive.Identity) (domain.Resource, error) {
	r, err := adapter.model.FindById(index)
	if err == nil {
		return &r, nil
	}
	if !commonrepo.IsErrorResourceNotExists(err) {
		return nil, err
	}

	dr, err := adapter.dataset.FindById(index)
	if err == nil {
		return &dr, nil
	}
	if !commonrepo.IsErrorResourceNotExists(err) {
		return nil, err
	}

	space, err := adapter.space.FindById(index)

	return &space, err
}

func (adapter *resourceAdapterImpl) Save(r domain.Resource) error {
	switch t := r.(type) {
	case *modeldomain.Model:
		return adapter.model.Save(t)
	case *spacedomain.Space:
		return adapter.space.Save(t)
	case *datasetsdomain.Dataset:
		return adapter.dataset.Save(t)
	default:
		return commonrepo.NewErrorResourceNotExists(errors.New("unknown resource type"))
	}
}
