package repositoryimpl

import (
	"gorm.io/gorm/clause"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/user/domain"
	"github.com/openmerlin/merlin-server/user/domain/repository"
)

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

func (impl *tokenRepoImpl) Add(u *domain.PlatformToken) (new domain.PlatformToken, err error) {
	u.Id = primitive.CreateIdentity(primitive.GetId())
	do := toTokenDO(u)

	err = impl.DB().Clauses(clause.Returning{}).Create(&do).Error
	if err != nil {
		return
	}

	new = do.toToken()

	return
}

func (impl *tokenRepoImpl) Delete(acc primitive.Account, name primitive.TokenName) (err error) {
	return impl.DB().Where(impl.EqualQuery(fieldName), name.TokenName()).Where(impl.EqualQuery(fieldOwner), acc.Account()).Delete(&TokenDO{}).Error
}

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

func (impl *tokenRepoImpl) GetByLastEight(LastEight string) (r []domain.PlatformToken, err error) {
	query := impl.DB().Where(
		impl.EqualQuery(fieldLastEight), LastEight,
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

func (impl *tokenRepoImpl) GetByName(acc primitive.Account, name primitive.TokenName) (r domain.PlatformToken, err error) {
	tmpDo := &TokenDO{}
	tmpDo.Name = name.TokenName()
	tmpDo.Owner = acc.Account()

	if err = impl.GetRecord(&tmpDo, &tmpDo); err != nil {
		return
	}

	r = tmpDo.toToken()

	return
}
