package modelrepositoryadapter

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/openmerlin/merlin-server/common/domain/repository"
)

var dbInstance *gorm.DB

type daoImpl struct {
	table string
}

// Each operation must generate a new gorm.DB instance.
// If using the same gorm.DB instance by different operations, they will share the same error.
func (dao *daoImpl) db() *gorm.DB {
	return dbInstance.Table(dao.table)
}

func (dao *daoImpl) GetRecord(filter, result interface{}) error {
	err := dao.db().Where(filter).First(result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repository.NewErrorResourceNotExists(errors.New("not found"))
	}

	return err
}

func (dao *daoImpl) GetByPrimaryKey(row interface{}) error {
	err := dao.db().First(row).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repository.NewErrorResourceNotExists(errors.New("not found"))
	}

	return err
}

func (dao *daoImpl) DeleteByPrimaryKey(row interface{}) error {
	err := dao.db().Delete(row).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repository.NewErrorResourceNotExists(errors.New("not found"))
	}

	return err
}

func likeFilter(field, value string) (query, arg string) {
	query = fmt.Sprintf(`%s ilike ?`, field)

	arg = `%` + value + `%`

	return
}

func equalQuery(field string) string {
	return fmt.Sprintf(`%s = ?`, field)
}

func orderByDesc(field string) string {
	return field + " desc"
}
