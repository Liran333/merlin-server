/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package spacerepositoryadapter

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

type modelSpaceRelationAdapter struct {
	daoImpl
}

// GetModelsBySpaceId find models related to a space
func (adapter *modelSpaceRelationAdapter) GetModelsBySpaceId(spaceId primitive.Identity) ([]primitive.Identity, error) {
	query := adapter.db().Where(equalQuery(fieldSpaceId), spaceId.Integer())

	var ids []modelSpaceRelationDO

	err := query.Find(&ids).Error

	dtos := make([]primitive.Identity, 0, len(ids))

	for _, id := range ids {
		dtos = append(dtos, primitive.CreateIdentity(id.ModelId))
	}

	return dtos, err
}

// GetSpacesByModelId find spaces related to a model
func (adapter *modelSpaceRelationAdapter) GetSpacesByModelId(modelId primitive.Identity) ([]primitive.Identity, error) {
	query := adapter.db().Where(equalQuery(fieldModelId), modelId.Integer())

	var ids []modelSpaceRelationDO

	err := query.Find(&ids).Error
	if err != nil || len(ids) == 0 {
		return []primitive.Identity{}, err
	}

	dtos := make([]primitive.Identity, 0, len(ids))

	for _, id := range ids {
		dtos = append(dtos, primitive.CreateIdentity(id.SpaceId))
	}

	return dtos, err
}

func (adapter *modelSpaceRelationAdapter) deleteBySpaceId(spaceId primitive.Identity) error {
	return adapter.db().Where(equalQuery(fieldSpaceId), spaceId.Identity()).Delete(&modelSpaceRelationDO{}).Error
}

func (adapter *modelSpaceRelationAdapter) create(spaceId primitive.Identity, modelId primitive.Identity) error {
	return adapter.db().Create(&modelSpaceRelationDO{
		SpaceId: spaceId.Integer(),
		ModelId: modelId.Integer(),
	}).Error
}

// UpdateRelation updates models related to a space
func (adapter *modelSpaceRelationAdapter) UpdateRelation(
	spaceId primitive.Identity, modelIds []primitive.Identity) error {
	return adapter.db().Transaction(func(tx *gorm.DB) error {
		if err := adapter.deleteBySpaceId(spaceId); err != nil {
			return fmt.Errorf("failed to delete model space relations by space id: %w", err)
		}

		for _, modelId := range modelIds {
			if err := adapter.create(spaceId, modelId); err != nil {
				return fmt.Errorf("failed to create model space relation: %w", err)
			}
		}

		return nil
	})
}

// DeleteBySpaceId deletes a model by its ID.
func (adapter *modelSpaceRelationAdapter) DeleteByModelId(modelId primitive.Identity) error {
	return adapter.db().Where(equalQuery(fieldModelId), modelId.Identity()).Delete(&modelSpaceRelationDO{}).Error
}

// DeleteByModelId deletes a space by its ID.
func (adapter *modelSpaceRelationAdapter) DeleteBySpaceId(spaceId primitive.Identity) error {
	return adapter.db().Where(equalQuery(fieldSpaceId), spaceId.Identity()).Delete(&modelSpaceRelationDO{}).Error
}
