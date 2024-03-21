/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	commonapp "github.com/openmerlin/merlin-server/common/app"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	modelapp "github.com/openmerlin/merlin-server/models/app"
	"github.com/openmerlin/merlin-server/models/domain"
	modelrepo "github.com/openmerlin/merlin-server/models/domain/repository"
	spacerepo "github.com/openmerlin/merlin-server/space/domain/repository"
)

// ModelAppService is an interface for the model application service.
type ModelSpaceAppService interface {
	GetModelsBySpaceId(user primitive.Account, spaceId primitive.Identity) ([]SpaceModelDTO, error)
	GetSpacesByModelId(user primitive.Account, modelId primitive.Identity) ([]SpaceModelDTO, error)
	UpdateRelation(spaceId primitive.Identity, modelsIndex []*domain.ModelIndex) error
	DeleteBySpaceId(modelId primitive.Identity) error
	DeleteByModelId(spaceId primitive.Identity) error
}

// NewModelAppService creates a new instance of the model application service.
func NewModelSpaceAppService(
	permission commonapp.ResourcePermissionAppService,
	repoAdapter spacerepo.ModelSpaceRepositoryAdapter,
	modelRepoAdapter modelrepo.ModelRepositoryAdapter,
	spaceRepoAdapter spacerepo.SpaceRepositoryAdapter,
	modelInternalApp modelapp.ModelInternalAppService,
	spaceInternalApp SpaceInternalAppService,
) ModelSpaceAppService {
	return &modelSpaceAppService{
		permission:       permission,
		repoAdapter:      repoAdapter,
		modelRepoAdapter: modelRepoAdapter,
		spaceRepoAdapter: spaceRepoAdapter,
		modelInternalApp: modelInternalApp,
		spaceInternlApp:  spaceInternalApp,
	}
}

type modelSpaceAppService struct {
	permission       commonapp.ResourcePermissionAppService
	repoAdapter      spacerepo.ModelSpaceRepositoryAdapter
	modelRepoAdapter modelrepo.ModelRepositoryAdapter
	spaceRepoAdapter spacerepo.SpaceRepositoryAdapter
	modelInternalApp modelapp.ModelInternalAppService
	spaceInternlApp  SpaceInternalAppService
}

// GetModelsBySpaceId return models that exits and user can read
func (s *modelSpaceAppService) GetModelsBySpaceId(user primitive.Account, spaceId primitive.Identity) ([]SpaceModelDTO, error) {
	space, err := s.spaceRepoAdapter.FindById(spaceId)
	if err != nil && commonrepo.IsErrorResourceNotExists(err) {
		return []SpaceModelDTO{}, newSpaceNotFound(err)
	}

	if err = s.permission.CanRead(user, &space); err != nil {
		if allerror.IsNoPermission(err) {
			err = newSpaceNotFound(err)
		}

		return []SpaceModelDTO{}, err
	}

	modelIds, err := s.repoAdapter.GetModelsBySpaceId(spaceId)
	if err != nil && commonrepo.IsErrorResourceNotExists(err) {
		return []SpaceModelDTO{}, newSpaceNotFound(err)
	}

	var models []SpaceModelDTO
	for _, id := range modelIds {
		// check if model exists
		model, err := s.modelRepoAdapter.FindById(id)
		if err != nil && commonrepo.IsErrorResourceNotExists(err) {
			_ = s.DeleteByModelId(id)
			continue
		}

		// check if user can read the model
		if err := s.permission.CanRead(user, &model); err != nil {
			continue
		}

		spaceModel := SpaceModelDTO{
			Owner:         model.CodeRepo.Owner.Account(),
			Name:          model.CodeRepo.Name.MSDName(),
			UpdatedAt:     model.UpdatedAt,
			LikeCount:     model.LikeCount,
			DownloadCount: model.DownloadCount,
		}
		models = append(models, spaceModel)
	}

	return models, nil
}

func (s *modelSpaceAppService) GetSpacesByModelId(user primitive.Account, modelId primitive.Identity) ([]SpaceModelDTO, error) {
	model, err := s.modelRepoAdapter.FindById(modelId)
	if err != nil && commonrepo.IsErrorResourceNotExists(err) {
		return []SpaceModelDTO{}, newModelNotFound(err)
	}

	if err = s.permission.CanRead(user, &model); err != nil {
		if allerror.IsNoPermission(err) {
			err = newModelNotFound(err)
		}

		return []SpaceModelDTO{}, err
	}

	spaceIds, err := s.repoAdapter.GetSpacesByModelId(modelId)
	if err != nil && commonrepo.IsErrorResourceNotExists(err) {
		return []SpaceModelDTO{}, newModelNotFound(err)
	}

	var spaces []SpaceModelDTO
	for _, id := range spaceIds {
		// check if model exists
		space, err := s.spaceRepoAdapter.FindById(id)
		if err != nil && commonrepo.IsErrorResourceNotExists(err) {
			_ = s.DeleteBySpaceId(id)
			continue
		}

		// check if user can read the space
		if err := s.permission.CanRead(user, &space); err != nil {
			continue
		}

		spaceModel := SpaceModelDTO{
			Owner:         space.CodeRepo.Owner.Account(),
			Name:          space.CodeRepo.Name.MSDName(),
			UpdatedAt:     space.UpdatedAt,
			LikeCount:     space.LikeCount,
			DownloadCount: space.DownloadCount,
		}
		spaces = append(spaces, spaceModel)
	}

	return spaces, nil
}

func (s *modelSpaceAppService) UpdateRelation(spaceId primitive.Identity, modelsIndex []*domain.ModelIndex) error {
	modelsId := s.modelInternalApp.GetByNames(modelsIndex)

	return s.repoAdapter.UpdateRelation(spaceId, modelsId)
}

func (s *modelSpaceAppService) DeleteByModelId(modelId primitive.Identity) error {
	return s.repoAdapter.DeleteByModelId(modelId)
}

func (s *modelSpaceAppService) DeleteBySpaceId(spaceId primitive.Identity) error {
	return s.repoAdapter.DeleteBySpaceId(spaceId)
}
