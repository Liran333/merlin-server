/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package spacerepositoryadapter

import (
	"errors"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/space/domain"
	"github.com/openmerlin/merlin-server/space/domain/repository"
)

type spaceVariableAdapter struct {
	daoImpl
}

// Add adds a new space variable to the database and returns an error if any occurs.
func (adapter *spaceVariableAdapter) AddVariable(variable *domain.SpaceVariable) error {
	do := toSpaceVariableDO(variable)

	v := adapter.db().Create(&do)

	return v.Error
}

// FindById finds a space variable by its ID and returns it along with an error if any occurs.
func (adapter *spaceVariableAdapter) FindVariableById(variableId primitive.Identity) (domain.SpaceVariable, error) {
	do := spaceEnvSecretDO{Id: variableId.Integer()}

	if err := adapter.GetByPrimaryKey(&do); err != nil {
		return domain.SpaceVariable{}, err
	}

	return do.toSpaceVariable(), nil
}

// Delete deletes a space variable from the database by its ID and returns an error if any occurs.
func (adapter *spaceVariableAdapter) DeleteVariable(variableId primitive.Identity) error {
	return adapter.DeleteByPrimaryKey(
		&spaceEnvSecretDO{Id: variableId.Integer()},
	)
}

// Save updates a space variable in the database and returns an error if any occurs.
func (adapter *spaceVariableAdapter) SaveVariable(variable *domain.SpaceVariable) error {
	do := toSpaceVariableDO(variable)
	do.Id = variable.Id.Integer()

	v := adapter.db().Model(
		&spaceEnvSecretDO{Id: variable.Id.Integer()},
	).Select(`*`).Updates(&do)

	if v.Error != nil {
		return v.Error
	}

	if v.RowsAffected == 0 {
		return commonrepo.NewErrorConcurrentUpdating(
			errors.New("concurrent updating"),
		)
	}

	return nil
}

// CountVariable count a space variable in the database and returns an error if any occurs.
func (adapter *spaceVariableAdapter) CountVariable(spaceId primitive.Identity) (int, error) {
	var total int64
	err := adapter.db().
		Where(equalQuery(fieldType), variableTypeName).
		Where(equalQuery(filedSpaceId), spaceId).Count(&total).Error

	return int(total), err
}

// Count is a method of spaceAdapter that takes a ListOption pointer as input
// and returns the total count of space variables and an error if any occurs.
func (adapter *spaceVariableAdapter) ListVariableSecret(spaceId primitive.Identity) (
	[]repository.SpaceVariableSecretSummary, error) {

	// list
	var dos []spaceEnvSecretDO

	err := adapter.db().
		Where(equalQuery(filedSpaceId), spaceId).
		Find(&dos).Error
	if err != nil || len(dos) == 0 {
		return nil, nil
	}

	r := make([]repository.SpaceVariableSecretSummary, len(dos))
	for i := range dos {
		r[i] = dos[i].toSpaceVariableSecretSummary()
	}

	return r, nil
}
