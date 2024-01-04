package postgresql

import (
	"errors"
	"fmt"

	"github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/openmerlin/merlin-server/common/domain/repository"
)

func DAO(table string) *daoImpl {
	return &daoImpl{
		table: table,
	}
}

type daoImpl struct {
	table string
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
