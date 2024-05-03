/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package postgresql provides functionality for interacting with PostgreSQL databases.
package postgresql

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/utils"
)

// Impl is an interface for database operations.
type Impl interface {
	GetRecord(filter, result interface{}) error
	GetByPrimaryKey(row interface{}) error
	DeleteByPrimaryKey(row interface{}) error
	LikeFilter(field, value string) (query, arg string)
	IntersectionFilter(field string, value []string) (query string, arg pq.StringArray)
	EqualQuery(field string) string
	NotEqualQuery(field string) string
	OrderByDesc(field string) string
	InFilter(field string) string
	DB() *gorm.DB
	TableName() string
}

// DAO creates a new daoImpl instance with the specified table name.
func DAO(table string) *daoImpl {
	return &daoImpl{
		table: table,
	}
}

type daoImpl struct {
	table string
}

// CommonModel is a struct that represents a common model with ID, CreatedAt, and UpdatedAt fields.
type CommonModel struct {
	ID        int64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// DB Each operation must generate a new gorm.DB instance.
// If using the same gorm.DB instance by different operations, they will share the same error.
func (dao *daoImpl) DB() *gorm.DB {
	return db.Table(dao.table)
}

// GetRecord retrieves a single record that matches the given filter criteria.
func (dao *daoImpl) GetRecord(filter, result interface{}) error {
	err := dao.DB().Where(filter).First(result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repository.NewErrorResourceNotExists(errors.New("not found"))
	}

	return err
}

// GetByPrimaryKey retrieves a single record by its primary key.
func (dao *daoImpl) GetByPrimaryKey(row interface{}) error {
	err := dao.DB().First(row).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repository.NewErrorResourceNotExists(errors.New("not found"))
	}

	return err
}

// DeleteByPrimaryKey deletes a record by its primary key.
func (dao *daoImpl) DeleteByPrimaryKey(row interface{}) error {
	err := dao.DB().Delete(row).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repository.NewErrorResourceNotExists(errors.New("not found"))
	}

	return err
}

// LikeFilter generates a query string and argument for a "like" filter condition.
func (dao *daoImpl) LikeFilter(field, value string) (query, arg string) {
	query = fmt.Sprintf(`%s ilike ?`, field)

	arg = `%` + utils.EscapePgsqlValue(value) + `%`

	return
}

// IntersectionFilter generates a query string and argument for an "intersection" filter condition.
func (dao *daoImpl) IntersectionFilter(field string, value []string) (query string, arg pq.StringArray) {
	query = fmt.Sprintf(`%s @> ?`, field)

	arg = pq.StringArray(value)

	return
}

// EqualQuery generates a query string for an "equal" filter condition.
func (dao *daoImpl) EqualQuery(field string) string {
	return fmt.Sprintf(`%s = ?`, field)
}

// MultiEqualQuery generates a query string for multiple "equal" filter conditions.
func (dao *daoImpl) MultiEqualQuery(fields ...string) string {
	v := make([]string, len(fields))

	for i, field := range fields {
		v[i] = dao.EqualQuery(field)
	}

	return strings.Join(v, " AND ")
}

// NotEqualQuery generates a query string for a "not equal" filter condition.
func (dao *daoImpl) NotEqualQuery(field string) string {
	return fmt.Sprintf(`%s <> ?`, field)
}

// OrderByDesc generates a query string for ordering results in descending order by the specified field.
func (dao *daoImpl) OrderByDesc(field string) string {
	return field + " desc"
}

// InFilter generates a query string and argument for an "in" filter condition.
func (dao *daoImpl) InFilter(field string) string {
	return fmt.Sprintf(`%s IN ?`, field)
}

// TableName returns the name of the table associated with this daoImpl instance.
func (dao *daoImpl) TableName() string {
	return dao.table
}

// IsRecordExists checks if the given error indicates that a unique constraint violation occurred.
func (dao *daoImpl) IsRecordExists(err error) bool {
	var pgError *pgconn.PgError
	ok := errors.As(err, &pgError)
	if !ok {
		return false
	}

	return pgError != nil && pgError.Code == errorCodes.UniqueConstraint
}
