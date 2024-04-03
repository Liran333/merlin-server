/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package spacerepositoryadapter

import (
	"errors"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/space/domain"
)

type spaceSecretAdapter struct {
	daoImpl
}

// Add adds a new space secret to the database and returns an error if any occurs.
func (adapter *spaceSecretAdapter) AddSecret(secret *domain.SpaceSecret) error {
	do := toSpaceSecretDO(secret)

	v := adapter.db().Create(&do)

	return v.Error
}

// FindById finds a space secret by its ID and returns it along with an error if any occurs.
func (adapter *spaceSecretAdapter) FindSecretById(spaceId primitive.Identity) (domain.SpaceSecret, error) {
	do := spaceEnvSecretDO{Id: spaceId.Integer()}

	if err := adapter.GetByPrimaryKey(&do); err != nil {
		return domain.SpaceSecret{}, err
	}

	return do.toSpaceSecret(), nil
}

// Delete deletes a space secret from the database by its ID and returns an error if any occurs.
func (adapter *spaceSecretAdapter) DeleteSecret(secretId primitive.Identity) error {
	return adapter.DeleteByPrimaryKey(
		&spaceEnvSecretDO{Id: secretId.Integer()},
	)
}

// Save updates a space in the database and returns an error if any occurs.
func (adapter *spaceSecretAdapter) SaveSecret(secret *domain.SpaceSecret) error {
	do := toSpaceSecretDO(secret)
	do.Id = secret.Id.Integer()

	v := adapter.db().Model(
		&spaceEnvSecretDO{Id: secret.Id.Integer()},
	).Select(`*`).Updates(&do)

	if v.Error != nil {
		return v.Error
	}

	if v.RowsAffected == 0 {
		return repository.NewErrorConcurrentUpdating(
			errors.New("concurrent updating"),
		)
	}

	return nil
}

// Count is a method of spaceAdapter that takes a ListOption pointer as input
// and returns the total count of spaces and an error if any occurs.
func (adapter *spaceSecretAdapter) CountSecret(spaceId primitive.Identity) (int, error) {
	var total int64
	err := adapter.db().
		Where(equalQuery(fieldType), secretTypeName).
		Where(equalQuery(filedSpaceId), spaceId).Count(&total).Error

	return int(total), err
}
