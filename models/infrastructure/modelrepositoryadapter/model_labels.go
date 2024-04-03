/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package modelrepositoryadapter provides an adapter for the model repository
package modelrepositoryadapter

import (
	"errors"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/models/domain"
)

type modelLabelsAdapter struct {
	daoImpl
}

// Save saves the model labels to the database.
func (adapter *modelLabelsAdapter) Save(modelId primitive.Identity, labels *domain.ModelLabels) error {
	do := toLabelsDO(labels)

	v := adapter.db().Model(
		&modelDO{Id: modelId.Integer()},
	).Select(
		fieldTask, fieldOthers, fieldFrameworks, fieldLicense,
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
