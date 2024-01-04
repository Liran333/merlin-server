package loginrepositoryadapter

import (
	"gorm.io/gorm"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/session/domain"
)

type dao interface {
	DB() *gorm.DB
	GetRecord(filter, result interface{}) error
	DeleteByPrimaryKey(row interface{}) error

	EqualQuery(field string) string
}

type loginAdapter struct {
	dao
}

func (adapter *loginAdapter) Add(login *domain.Login) error {
	do := toLoginDO(login)

	v := adapter.DB().Create(&do)

	return v.Error
}

func (adapter *loginAdapter) Delete(loginId primitive.UUID) error {
	return adapter.DeleteByPrimaryKey(
		&loginDO{Id: loginId},
	)
}

func (adapter *loginAdapter) Find(loginId primitive.UUID) (domain.Login, error) {
	do := loginDO{Id: loginId}

	if err := adapter.GetRecord(&do, &do); err != nil {
		return domain.Login{}, err
	}

	return do.toLogin(), nil
}

func (adapter *loginAdapter) FindByUser(user primitive.Account) ([]domain.Login, error) {
	query := adapter.DB().Where(
		adapter.EqualQuery(fieldUser), user.Account(),
	).Order(fieldCreatedAt)

	var dos []loginDO

	err := query.Find(&dos).Error
	if err != nil || len(dos) == 0 {
		return nil, nil
	}

	r := make([]domain.Login, len(dos))
	for i := range dos {
		r[i] = dos[i].toLogin()
	}

	return r, nil
}
