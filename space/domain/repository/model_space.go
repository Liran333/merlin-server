/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package repository provides interfaces for interacting with models and model labels in the domain.
package repository

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

// ModelSpaceRepositoryAdapter provides an adapter for the model repository
type ModelSpaceRepositoryAdapter interface {
	GetModelsBySpaceId(spaceId primitive.Identity) ([]primitive.Identity, error)
	GetSpacesByModelId(modelId primitive.Identity) ([]primitive.Identity, error)
	UpdateRelation(spaceId primitive.Identity, modelIds []primitive.Identity) error
	DeleteByModelId(modelId primitive.Identity) error
	DeleteBySpaceId(spaceId primitive.Identity) error
}
