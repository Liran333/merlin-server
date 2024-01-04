package repositoryimpl

import (
	"context"
	"fmt"

	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	mongo "github.com/openmerlin/merlin-server/common/infrastructure/mongo"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/user/domain"
	"github.com/openmerlin/merlin-server/user/domain/repository"
)

func NewTokenRepo(m mongodbClient) repository.Token {
	return &tokenRepoImpl{m}
}

type tokenRepoImpl struct {
	cli mongodbClient
}

func (impl *tokenRepoImpl) Save(u *domain.PlatformToken) (r domain.PlatformToken, err error) {
	if u.Id != "" {
		if err = impl.update(u); err != nil {
			err = fmt.Errorf("failed to update token: %w", err)
		} else {
			r = *u
			r.Version += 1
		}

		return
	}

	v, err := impl.insert(u)
	if err != nil {
		err = fmt.Errorf("failed to add token info %w", err)
	} else {
		r = *u
		r.Id = v
	}

	return
}

func tokenFilterByAccountAndName(account, name string) bson.M {
	return bson.M{
		fieldName:    name,
		fieldAccount: account,
	}
}

func tokenFilterByAccount(account string) bson.M {
	return bson.M{
		fieldAccount: account,
	}
}

func (impl *tokenRepoImpl) Delete(acc domain.Account, name string) (err error) {
	f := func(ctx context.Context) error {
		return impl.cli.DeleteOne(
			ctx, tokenFilterByAccountAndName(acc.Account(), name),
		)
	}

	if err = primitive.WithContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = commonrepo.NewErrorResourceNotExists(err)
		}
	}

	return
}

func (impl *tokenRepoImpl) GetByAccount(account domain.Account) (r []domain.PlatformToken, err error) {
	var v []DToken
	f := func(ctx context.Context) error {
		return impl.cli.GetDocs(
			ctx, tokenFilterByAccount(account.Account()),
			bson.M{}, &v,
		)
	}

	if err = primitive.WithContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = commonrepo.NewErrorResourceNotExists(err)
		}

		return
	}

	r = make([]domain.PlatformToken, len(v))
	for i := range v {
		toToken(v[i], &r[i])
	}

	return
}

func tokenFilterByLastEight(LastEight string) bson.M {
	return bson.M{
		fieldLastEight: LastEight,
	}
}

func (impl *tokenRepoImpl) GetByLastEight(LastEight string) (r []domain.PlatformToken, err error) {
	var v []DToken
	f := func(ctx context.Context) error {
		return impl.cli.GetDocs(
			ctx, tokenFilterByLastEight(LastEight),
			bson.M{}, &v,
		)
	}

	if err = primitive.WithContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = nil
		}

		return
	}

	r = make([]domain.PlatformToken, len(v))
	for i := range v {
		toToken(v[i], &r[i])
	}

	return
}

func (impl *tokenRepoImpl) update(u *domain.PlatformToken) (err error) {
	var token DToken
	err = toTokenDoc(*u, &token)
	if err != nil {
		return
	}
	doc, err := mongo.GenDoc(token)
	if err != nil {
		return
	}

	filter, err := mongo.ObjectIdFilter(u.Id)
	if err != nil {
		return
	}

	f := func(ctx context.Context) error {
		return impl.cli.UpdateDoc(
			ctx, filter, doc, mongoCmdSet, u.Version,
		)
	}

	if err = primitive.WithContext(f); err != nil && impl.cli.IsDocNotExists(err) {
		err = fmt.Errorf("concurrent updating: %w", err)
	}

	return
}

func (impl *tokenRepoImpl) insert(u *domain.PlatformToken) (id string, err error) {
	var token DToken
	err = toTokenDoc(*u, &token)
	if err != nil {
		return
	}

	doc, err := mongo.GenDoc(token)
	if err != nil {
		return
	}

	doc[fieldVersion] = 0
	doc[fieldFollower] = bson.A{}
	doc[fieldFollowing] = bson.A{}

	f := func(ctx context.Context) error {
		v, err := impl.cli.NewDocIfNotExist(
			ctx, mongo.UserDocFilterByAccount(u.Account.Account()), doc,
		)

		id = v

		return err
	}

	if err = primitive.WithContext(f); err != nil && impl.cli.IsDocExists(err) {
		err = commonrepo.NewErrorDuplicateCreating(err)
	}

	return
}
