/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package repositoryimpl

import (
	"golang.org/x/xerrors"
	"gorm.io/gorm/clause"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/user/domain"
	"github.com/openmerlin/merlin-server/user/domain/repository"
)

// NewTokenRepo creates a new token repository with the given database implementation.
func NewTokenRepo(db postgresql.Impl) repository.Token {
	tokenTableName = db.TableName()

	if err := postgresql.DB().AutoMigrate(&TokenDO{}); err != nil {
		return nil
	}

	return &tokenRepoImpl{Impl: db}
}

type tokenRepoImpl struct {
	postgresql.Impl
}

// Add adds a new platform token to the repository and returns it.
func (impl *tokenRepoImpl) Add(u *domain.PlatformToken) (new domain.PlatformToken, err error) {
	u.Id = primitive.CreateIdentity(primitive.GetId())
	do := toTokenDO(u)

	err = impl.DB().Clauses(clause.Returning{}).Create(&do).Error
	if err != nil {
		err = xerrors.Errorf("failed to add token: %w", err)
		return
	}

	new = do.toToken()

	return
}

// Delete deletes a token from the repository based on the account and name.
func (impl *tokenRepoImpl) Delete(acc primitive.Account, name primitive.TokenName) (err error) {
	err = impl.DB().Where(impl.EqualQuery(fieldName),
		name.TokenName()).Where(impl.EqualQuery(fieldOwner), acc.Account()).Delete(&TokenDO{}).Error
	if err != nil {
		err = xerrors.Errorf("failed to delete token: %w", err)
		return
	}

	return
}

// GetByAccount retrieves all tokens owned by the given account.
func (impl *tokenRepoImpl) GetByAccount(account domain.Account) (r []domain.PlatformToken, err error) {
	query := impl.DB().Where(
		impl.EqualQuery(fieldOwner), account.Account(),
	)

	var dos []TokenDO

	err = query.Find(&dos).Error
	if err != nil || len(dos) == 0 {
		return nil, err
	}

	r = make([]domain.PlatformToken, len(dos))
	for i := range dos {
		r[i] = dos[i].toToken()
	}

	return
}

// GetByLastEight retrieves all tokens that match the given last eight characters.
func (impl *tokenRepoImpl) GetByLastEight(LastEight string) (r []domain.PlatformToken, err error) {
	query := impl.DB().Where(
		impl.EqualQuery(fieldLastEight), LastEight,
	)

	var dos []TokenDO

	err = query.Find(&dos).Error
	if err != nil {
		return nil, xerrors.Errorf("failed to get token by last eight: %w", err)
	} else if len(dos) == 0 {
		return nil, xerrors.Errorf("not a valid token")
	}

	r = make([]domain.PlatformToken, len(dos))
	for i := range dos {
		r[i] = dos[i].toToken()
	}

	return
}

// GetByName retrieves a token by its name and owner.
func (impl *tokenRepoImpl) GetByName(acc primitive.Account,
	name primitive.TokenName) (r domain.PlatformToken, err error) {
	tmpDo := &TokenDO{}
	tmpDo.Name = name.TokenName()
	tmpDo.Owner = acc.Account()

	if err = impl.GetRecord(&tmpDo, &tmpDo); err != nil {
		err = xerrors.Errorf("failed to get token by name: %w", err)
		return
	}

	r = tmpDo.toToken()

	return
}
