/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package modelrepositoryadapter provides an adapter for the model repository
package modelrepositoryadapter

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/openmerlin/merlin-server/common/domain/repository"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/utils"
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

// GetRecord retrieves a single record from the database based on the provided filter
// and stores it in the result parameter.
func (dao *daoImpl) GetRecord(filter, result interface{}) error {
	err := dao.db().Where(filter).First(result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repository.NewErrorResourceNotExists(errors.New("not found"))
	}

	return err
}

// GetLowerModelName retrieves a single lower model name from the database based on the provided filter
// and stores it in the result parameter.
func (dao *daoImpl) GetLowerModelName(filter *modelDO, result interface{}) error {
	err := dao.db().Where("LOWER(name) = ? AND LOWER(owner) = ?",
		strings.ToLower(filter.Name), strings.ToLower(filter.Owner)).First(result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repository.NewErrorResourceNotExists(errors.New("not found"))
	}

	return err
}

// GetByPrimaryKey retrieves a single record from the database based on the primary key of the row parameter.
func (dao *daoImpl) GetByPrimaryKey(row interface{}) error {
	err := dao.db().First(row).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repository.NewErrorResourceNotExists(errors.New("not found"))
	}

	return err
}

// DeleteByPrimaryKey deletes a single record from the database based on the primary key of the row parameter.
func (dao *daoImpl) DeleteByPrimaryKey(row interface{}) error {
	err := dao.db().Delete(row).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repository.NewErrorResourceNotExists(errors.New("not found"))
	}

	return err
}

// likeFilter generates a ILIKE filter query and argument for a given field and value.
func likeFilter(field, value string) (query, arg string) {
	query = fmt.Sprintf(`%s ilike ?`, field)

	arg = `%` + utils.EscapePgsqlValue(value) + `%`

	return
}

// intersectionFilter generates an intersection filter query and argument for a given field and value.
func intersectionFilter(field string, value []string) (query string, arg pq.StringArray) {
	query = fmt.Sprintf(`%s @> ?`, field)

	arg = pq.StringArray(value)

	return
}

// equalQuery generates an equality filter query for a given field.
func equalQuery(field string) string {
	return fmt.Sprintf(`%s = ?`, field)
}

// notEqualQuery generates a not equal filter query for a given field.
func notEqualQuery(field string) string {
	return fmt.Sprintf(`%s <> ?`, field)
}

// orderByDesc generates an ORDER BY clause in descending order for a given field.
func orderByDesc(field string) string {
	return field + " desc"
}

// IncrementLikeCount increments the LikeCount field by 1 for a record with the specified primary key.
func (dao *daoImpl) IncrementLikeCount(id int64, version int) error {
	// Update version to handle high concurrency scenario
	result := dao.db().Model(&modelDO{Id: id}).Where(fieldVersion+" = ?", version).
		Updates(map[string]interface{}{
			fieldLikeCount: gorm.Expr("COALESCE("+fieldLikeCount+", 0) + ?", 1),
			fieldVersion:   gorm.Expr(fieldVersion+" + ?", 1),
		})

	// Check if any rows were updated
	if result.RowsAffected == 0 {
		return commonrepo.NewErrorConcurrentUpdating(
			errors.New("concurrent updating"),
		)
	}
	return nil
}

// DescendLikeCount Descend the LikeCount field by 1 for a record with the specified primary key.
func (dao *daoImpl) DescendLikeCount(id int64, version int) error {
	// Update version to handle high concurrency scenario
	result := dao.db().Model(&modelDO{Id: id}).Where(fieldVersion+" = ?", version).
		Updates(map[string]interface{}{
			fieldLikeCount: gorm.Expr("COALESCE("+fieldLikeCount+", 1) - ?", 1),
			fieldVersion:   gorm.Expr(fieldVersion+" + ?", 1),
		})

	// Check if any rows were updated
	if result.RowsAffected == 0 {
		return repository.NewErrorResourceNotExists(errors.New("resource not found"))
	}
	return nil
}
