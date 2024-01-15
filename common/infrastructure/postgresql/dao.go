package postgresql

import (
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/openmerlin/merlin-server/common/domain/repository"
)

type Impl interface {
	GetRecord(filter, result interface{}) error
	GetByPrimaryKey(row interface{}) error
	DeleteByPrimaryKey(row interface{}) error
	LikeFilter(field, value string) (query, arg string)
	IntersectionFilter(field string, value []string) (query string, arg pq.StringArray)
	EqualQuery(field string) string
	NotEqualQuery(field string) string
	InFilter(field string) string
	DB() *gorm.DB
	TableName() string
}

func DAO(table string) *daoImpl {
	return &daoImpl{
		table: table,
	}
}

type daoImpl struct {
	table string
}

type CommonModel struct {
	ID        int64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Each operation must generate a new gorm.DB instance.
// If using the same gorm.DB instance by different operations, they will share the same error.
func (dao *daoImpl) DB() *gorm.DB {
	return db.Table(dao.table)
}

func (dao *daoImpl) GetRecord(filter, result interface{}) error {
	err := dao.DB().Where(filter).First(result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repository.NewErrorResourceNotExists(errors.New("not found"))
	}

	return err
}

func (dao *daoImpl) GetByPrimaryKey(row interface{}) error {
	err := dao.DB().First(row).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repository.NewErrorResourceNotExists(errors.New("not found"))
	}

	return err
}

func (dao *daoImpl) DeleteByPrimaryKey(row interface{}) error {
	err := dao.DB().Delete(row).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repository.NewErrorResourceNotExists(errors.New("not found"))
	}

	return err
}

func (dao *daoImpl) LikeFilter(field, value string) (query, arg string) {
	query = fmt.Sprintf(`%s ilike ?`, field)

	arg = `%` + value + `%`

	return
}

func (dao *daoImpl) IntersectionFilter(field string, value []string) (query string, arg pq.StringArray) {
	query = fmt.Sprintf(`%s @> ?`, field)

	arg = pq.StringArray(value)

	return
}

func (dao *daoImpl) EqualQuery(field string) string {
	return fmt.Sprintf(`%s = ?`, field)
}

func (dao *daoImpl) NotEqualQuery(field string) string {
	return fmt.Sprintf(`%s <> ?`, field)
}

func (dao *daoImpl) OrderByDesc(field string) string {
	return field + " desc"
}

func (dao *daoImpl) InFilter(field string) string {
	return fmt.Sprintf(`%s IN ?`, field)
}

func (dao *daoImpl) TableName() string {
	return dao.table
}
