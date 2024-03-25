package activityrepositoryadapter

import (
	"errors"

	"gorm.io/gorm"

	"github.com/openmerlin/merlin-server/activity/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/common/domain/repository"
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
	if err := db.Where(condition, fieldLike, owner, index).
		Delete(&domain.Activity{}).Error; err != nil {
		return err
	}

	return nil
}
