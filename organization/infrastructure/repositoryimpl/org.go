package repositoryimpl

import (
	"context"
	"fmt"

	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	mongo "github.com/openmerlin/merlin-server/common/infrastructure/mongo"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/organization/domain"
	"github.com/openmerlin/merlin-server/organization/domain/repository"
)

func NewOrgRepo(m mongodbClient) repository.Organization {
	return &orgRepoImpl{m}
}

type orgRepoImpl struct {
	cli mongodbClient
}

func (impl *orgRepoImpl) Save(o *domain.Organization) (r domain.Organization, err error) {
	if o.Id != "" {
		if err = impl.update(o); err != nil {
			err = fmt.Errorf("failed to update user: %w", err)
		} else {
			r = *o
			r.Version += 1
		}

		return
	}

	v, err := impl.insert(o)
	if err != nil {
		err = fmt.Errorf("failed to add org info %w", err)
	} else {
		r = *o
		r.Id = v
	}

	return
}

func (impl *orgRepoImpl) Delete(o *domain.Organization) (err error) {
	filter, err := mongo.ObjectIdFilter(o.Id)
	if err != nil {
		return
	}

	f := func(ctx context.Context) error {
		return impl.cli.DeleteOne(
			ctx, filter,
		)
	}

	if err = primitive.WithContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = commonrepo.NewErrorResourceNotExists(err)
		}
	}

	return
}

func (impl *orgRepoImpl) update(o *domain.Organization) (err error) {
	org := toOrgDoc(*o)
	if err != nil {
		return
	}
	doc, err := mongo.GenDoc(org)
	if err != nil {
		return
	}

	filter, err := mongo.ObjectIdFilter(o.Id)
	if err != nil {
		return
	}

	f := func(ctx context.Context) error {
		return impl.cli.UpdateDoc(
			ctx, filter, doc, mongoCmdSet, o.Version,
		)
	}

	if err = primitive.WithContext(f); err != nil && impl.cli.IsDocNotExists(err) {
		err = fmt.Errorf("concurrent updating: %w", err)
	}

	return
}

func (impl *orgRepoImpl) insert(o *domain.Organization) (id string, err error) {
	org := toOrgDoc(*o)
	if err != nil {
		return
	}

	doc, err := mongo.GenDoc(org)
	if err != nil {
		return
	}

	doc[fieldVersion] = 0

	f := func(ctx context.Context) error {
		v, err := impl.cli.NewDocIfNotExist(
			ctx, mongo.UserDocFilterByAccount(o.Name.Account()), doc,
		)

		id = v

		return err
	}

	if err = primitive.WithContext(f); err != nil && impl.cli.IsDocExists(err) {
		err = commonrepo.NewErrorResourceNotExists(err)
	}

	return
}

func (impl *orgRepoImpl) GetByName(orgName primitive.Account) (
	o domain.Organization, err error,
) {
	var v Organization
	f := func(ctx context.Context) error {
		return impl.cli.GetDoc(
			ctx, mongo.UserDocFilterByAccount(orgName.Account()),
			bson.M{}, &v,
		)
	}

	if err = primitive.WithContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = commonrepo.NewErrorResourceNotExists(err)
		}

		return
	}

	err = toOrganization(v, &o)

	return
}

func inviteDocFilterByUser(account string) bson.M {
	return bson.M{
		fieldOwner: account,
	}
}

func (impl *orgRepoImpl) GetInviteByUser(acc primitive.Account) (
	os []domain.Organization, err error,
) {

	return
}

func orgDocFilterByOwner(account string) bson.M {
	return bson.M{
		fieldOwner: account,
	}
}

func (impl *orgRepoImpl) GetByOwner(owner primitive.Account) (
	o []domain.Organization, err error,
) {
	var v []Organization
	f := func(ctx context.Context) error {
		return impl.cli.GetDocs(
			ctx, orgDocFilterByOwner(owner.Account()),
			bson.M{}, &v,
		)
	}

	if err = primitive.WithContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = commonrepo.NewErrorResourceNotExists(err)
		}

		return
	}

	o = make([]domain.Organization, len(v))
	for i := range v {
		item := &v[i]
		if err := toOrganization(*item, &o[i]); err != nil {
			return nil, err
		}
	}

	return
}
