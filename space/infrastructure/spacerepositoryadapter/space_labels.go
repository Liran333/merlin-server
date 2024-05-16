/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package spacerepositoryadapter

import (
	"errors"

	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/space/domain"
)

type spaceLabelsAdapter struct {
	daoImpl
}

// Save saves the space labels to the database.
func (adapter *spaceLabelsAdapter) Save(index *domain.SpaceIndex, labels *domain.SpaceLabels) error {
	do := toLabelsDO(labels)

	v := adapter.db().Model(&spaceDO{}).Where(
		equalQuery(fieldOwner), index.Owner.Account(),
	).Where(
		equalQuery(fieldName), index.Name.MSDName(),
	).Select(
		fieldTask, fieldOthers, fieldBaseImage,
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
