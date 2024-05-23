/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package datasetrepositoryadapter provides an adapter for the dataset repository
package datasetrepositoryadapter

import (
	"errors"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/datasets/domain"
)

type datasetLabelsAdapter struct {
	daoImpl
}

// Save saves the dataset labels to the database.
func (adapter *datasetLabelsAdapter) Save(datasetId primitive.Identity, labels *domain.DatasetLabels) error {
	do := toLabelsDO(labels)

	v := adapter.db().Model(
		&datasetDO{Id: datasetId.Integer()},
	).Select(
		fieldTask, fieldLicense, fieldSize, fieldLanguage, fieldDomain,
	).Updates(&do)

	if v.Error != nil {
		return v.Error
	}

	if v.RowsAffected == 0 {
		return commonrepo.NewErrorResourceNotExists(
			errors.New("not found"),
		)
	}

	return nil
}
