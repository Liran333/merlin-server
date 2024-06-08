/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package domain provides domain models and functionality for managing space apps.
package repositoryadapter

import (
	"errors"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/spaceapp/domain"
)

type buildLogAdapterImpl struct {
	dao dao
}

func (adapter *buildLogAdapterImpl) Save(log *domain.SpaceAppBuildLog) error {
	v := adapter.dao.DB().Model(
		&spaceappDO{Id: log.AppId.Integer()},
	).Select(
		fieldAllBuildLog,
	).Updates(&spaceappDO{AllBuildLog: log.Logs})

	if v.Error != nil {
		return v.Error
	}

	if v.RowsAffected == 0 {
		return repository.NewErrorResourceNotExists(
			errors.New("not found"),
		)
	}

	return nil
}

func (adapter *buildLogAdapterImpl) Find(appId primitive.Identity) (
	log domain.SpaceAppBuildLog, err error,
) {
	do := spaceappDO{Id: appId.Integer()}

	if err = adapter.dao.GetByPrimaryKey(&do); err == nil {
		log.Logs = do.AllBuildLog
		log.AppId = appId
	}

	return
}
