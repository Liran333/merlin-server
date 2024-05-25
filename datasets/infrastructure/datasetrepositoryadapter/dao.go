/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package datasetrepositoryadapter provides an adapter for the dataset repository
package datasetrepositoryadapter

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"golang.org/x/xerrors"
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
		return repository.NewErrorResourceNotExists(xerrors.Errorf("%w", errors.New("not found")))
	}

	if err != nil {
		return xerrors.Errorf("failed to get dataset record, %w", err)
	}

	return nil
}

// GetLowerDatasetName retrieves a single lower dataset name from the database based on the provided filter
// and stores it in the result parameter.
func (dao *daoImpl) GetLowerDatasetName(filter *datasetDO, result interface{}) error {
	err := dao.db().Where("LOWER(name) = ? AND LOWER(owner) = ?", strings.ToLower(filter.Name), strings.ToLower(filter.Owner)).First(result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repository.NewErrorResourceNotExists(xerrors.Errorf("%w", errors.New("not found")))
	}

	if err != nil {
		return xerrors.Errorf("failed to get lower dataset name, %w", err)
	}

	return nil
}

// GetByPrimaryKey retrieves a single record from the database based on the primary key of the row parameter.
func (dao *daoImpl) GetByPrimaryKey(row interface{}) error {
	err := dao.db().First(row).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repository.NewErrorResourceNotExists(xerrors.Errorf("%w", errors.New("not found")))
	}

	if err != nil {
		return xerrors.Errorf("failed to get dataset by primary key, %w", err)
	}

	return nil
}

// DeleteByPrimaryKey deletes a single record from the database based on the primary key of the row parameter.
func (dao *daoImpl) DeleteByPrimaryKey(row interface{}) error {
	err := dao.db().Delete(row).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repository.NewErrorResourceNotExists(xerrors.Errorf("%w", errors.New("not found")))
	}

	if err != nil {
		return xerrors.Errorf("failed to delete dataset by primary key, %w", err)
	}

	return nil
}

func likeFilter(field, value string) (query, arg string) {
	query = fmt.Sprintf(`%s ilike ?`, field)

	arg = `%` + utils.EscapePgsqlValue(value) + `%`

	return
}

func intersectionFilter(field string, value []string) (query string, arg pq.StringArray) {
	query = fmt.Sprintf(`%s @> ?`, field)

	arg = pq.StringArray(value)

	return
}

func equalQuery(field string) string {
	return fmt.Sprintf(`%s = ?`, field)
}

func notEqualQuery(field string) string {
	return fmt.Sprintf(`%s <> ?`, field)
}

func orderByDesc(field string) string {
	return field + " desc"
}

// IncrementLikeCount increments the LikeCount field by 1 for a record with the specified primary key.
func (dao *daoImpl) IncrementLikeCount(id int64, version int) error {
	// Update version to handle high concurrency scenario
	result := dao.db().Model(&datasetDO{Id: id}).Where(fieldVersion+" = ?", version).
		Updates(map[string]interface{}{
			fieldLikeCount: gorm.Expr("COALESCE("+fieldLikeCount+", 0) + ?", 1),
			fieldVersion:   gorm.Expr(fieldVersion+" + ?", 1),
		})

	// Check if any rows were updated
	if result.RowsAffected == 0 {
		return commonrepo.NewErrorConcurrentUpdating(
			xerrors.Errorf("%w", errors.New("concurrent updating")),
		)
	}
	return nil
}

// DescendLikeCount Descend the LikeCount field by 1 for a record with the specified primary key.
func (dao *daoImpl) DescendLikeCount(id int64, version int) error {
	// Update version to handle high concurrency scenario
	result := dao.db().Model(&datasetDO{Id: id}).Where(fieldVersion+" = ?", version).
		Updates(map[string]interface{}{
			fieldLikeCount: gorm.Expr("COALESCE("+fieldLikeCount+", 1) - ?", 1),
			fieldVersion:   gorm.Expr(fieldVersion+" + ?", 1),
		})

	// Check if any rows were updated
	if result.RowsAffected == 0 {
		return repository.NewErrorResourceNotExists(xerrors.Errorf("%w", errors.New("resource not found")))
	}
	return nil
}
