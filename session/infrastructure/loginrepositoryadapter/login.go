/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package loginrepositoryadapter

import (
	"context"

	"gorm.io/gorm"

	"github.com/openmerlin/merlin-server/common/domain/crypto"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/session/domain"
)

type dao interface {
	DB() *gorm.DB
	GetRecord(ctx context.Context, filter, result interface{}) error
	DeleteByPrimaryKey(ctx context.Context, row interface{}) error

	EqualQuery(field string) string
}

type loginAdapter struct {
	dao
	e crypto.Encrypter
}

// Add adds a new login to the database.
func (adapter *loginAdapter) Add(login *domain.Session) error {
	do, err := toLoginDO(login, adapter.e)
	if err != nil {
		return err
	}

	v := adapter.DB().Create(&do)

	return v.Error
}

// Delete deletes a login from the database by its ID.
func (adapter *loginAdapter) Delete(ctx context.Context, loginId primitive.RandomId) error {
	return adapter.DeleteByPrimaryKey(
		ctx, &loginDO{Id: loginId.RandomId()},
	)
}

func (adapter *loginAdapter) DeleteByUser(user primitive.Account) error {
	return adapter.DB().
		Where(adapter.EqualQuery(fieldUser), user.Account()).
		Delete(&loginDO{}).Error
}

// Find finds a login in the database by its ID.
func (adapter *loginAdapter) Find(ctx context.Context, loginId primitive.RandomId) (domain.Session, error) {
	do := loginDO{Id: loginId.RandomId()}

	if err := adapter.GetRecord(ctx, &do, &do); err != nil {
		return domain.Session{}, err
	}

	return do.toLogin(adapter.e)
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
		if r[i], err = dos[i].toLogin(adapter.e); err != nil {
			return r, err
		}
	}

	return r, nil
}
