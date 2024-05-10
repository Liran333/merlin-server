/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package activityrepositoryadapter provides an adapter for the model repository
package activityrepositoryadapter

import (
	"errors"

	"gorm.io/gorm"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

var dbInstance *gorm.DB

type daoImpl struct {
	table string
}

// Each operation must generate a new gorm.DB instance.
// If using the same gorm.DB instance by different operations, they will share the same error.
func (dao *daoImpl) db() *gorm.DB {
	if dbInstance == nil {
		return nil
	}
	return dbInstance.Table(dao.table)
}

func orderByDesc(field string) string {
	return field + " desc"
}

func (adapter *activityAdapter) deleteLikeByOwnerAndIndex(owner primitive.Account, index primitive.Identity) error {
	db := adapter.daoImpl.db() // Get the gorm.DB instance.
	if db == nil {
		return errors.New("database instance is not initialized")
	}

	// Define the condition string as a constant or directly use it in the Where clause.
	const condition = "type = ? AND owner = ? AND resource_id = ?"

	// Execute the deletion with the provided conditions.
	if err := db.Where(condition, fieldLike, owner.Account(), index).
		Delete(&activityDO{}).Error; err != nil {
		return err
	}

	return nil
}
