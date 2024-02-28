/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

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

// Add adds a new login to the database.
func (adapter *loginAdapter) Add(login *domain.Session) error {
	do := toLoginDO(login)

	v := adapter.DB().Create(&do)

	return v.Error
}

// Delete deletes a login from the database by its ID.
func (adapter *loginAdapter) Delete(loginId primitive.RandomId) error {
	return adapter.DeleteByPrimaryKey(
		&loginDO{Id: loginId.RandomId()},
	)
}

// Find finds a login in the database by its ID.
func (adapter *loginAdapter) Find(loginId primitive.RandomId) (domain.Session, error) {
	do := loginDO{Id: loginId.RandomId()}

	if err := adapter.GetRecord(&do, &do); err != nil {
		return domain.Session{}, err
	}

	return do.toLogin(), nil
}

// FindByUser finds all logins in the database associated with a user.
func (adapter *loginAdapter) FindByUser(user primitive.Account) ([]domain.Session, error) {
	query := adapter.DB().Where(
		adapter.EqualQuery(fieldUser), user.Account(),
	).Order(fieldCreatedAt)

	var dos []loginDO

	err := query.Find(&dos).Error
	if err != nil || len(dos) == 0 {
		return nil, nil
	}

	r := make([]domain.Session, len(dos))
	for i := range dos {
		r[i] = dos[i].toLogin()
	}

	return r, nil
}
