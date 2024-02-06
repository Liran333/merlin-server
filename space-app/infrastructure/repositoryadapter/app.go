package repositoryadapter

import (
	"errors"

	"gorm.io/gorm"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/space-app/domain"
)

type dao interface {
	DB() *gorm.DB
	EqualQuery(field string) string
	IsRecordExists(err error) bool
}

type appRepositoryAdapter struct {
	dao dao
}

func (adapter *appRepositoryAdapter) Add(m *domain.SpaceApp) error {
	if err := adapter.remove(m.SpaceId); err != nil {
		return err
	}

	do := toSpaceAppDO(m)

	err := adapter.dao.DB().Create(&do).Error

	if err != nil && adapter.dao.IsRecordExists(err) {
		return repository.NewErrorDuplicateCreating(
			errors.New("space app exists"),
		)
	}

	return err
}

func (adapter *appRepositoryAdapter) remove(spaceId primitive.Identity) error {
	return adapter.dao.DB().Where(
		adapter.dao.EqualQuery(fieldSpaceId), spaceId.Identity(),
	).Delete(
		spaceappDO{},
	).Error
}
