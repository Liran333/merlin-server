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

func NewMemberRepo(m mongodbClient) repository.OrgMember {
	return &memberRepoImpl{m}
}

type memberRepoImpl struct {
	cli mongodbClient
}

func (impl *memberRepoImpl) Save(o *domain.OrgMember) (r domain.OrgMember, err error) {
	if o.Id != "" {
		if err = impl.update(o); err != nil {
			err = fmt.Errorf("failed to update org member: %w", err)
		} else {
			r = *o
			r.Version += 1
		}

		return
	}

	v, err := impl.insert(o)
	if err != nil {
		err = fmt.Errorf("failed to add member info %w", err)
	} else {
		r = *o
		r.Id = v
	}

	return
}

func (impl *memberRepoImpl) Delete(o *domain.OrgMember) (err error) {
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

func (impl *memberRepoImpl) update(o *domain.OrgMember) (err error) {
	org := toMemberDoc(*o)
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

func memberDocInsertFilter(username, orgname string) bson.M {
	return bson.M{
		fieldUser: username,
		fieldOrg:  orgname,
	}
}

func (impl *memberRepoImpl) insert(o *domain.OrgMember) (id string, err error) {
	org := toMemberDoc(*o)
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
			ctx, memberDocInsertFilter(o.Username, o.OrgName), doc,
		)

		id = v

		return err
	}

	if err = primitive.WithContext(f); err != nil && impl.cli.IsDocExists(err) {
		err = commonrepo.NewErrorResourceNotExists(err)
	}

	return
}

func memberDocOrgNameFilter(name string) bson.M {
	return bson.M{
		fieldOrg: name,
	}
}

func memberDocUserNameFilter(name string) bson.M {
	return bson.M{
		fieldUser: name,
	}
}

func memberDocOrgAndRoleFilter(org, role string) bson.M {
	return bson.M{
		fieldOrg:  org,
		fieldRole: role,
	}
}

func (impl *memberRepoImpl) GetByOrg(name string) (
	members []domain.OrgMember, err error,
) {
	var v []Member
	filter := memberDocOrgNameFilter(name)

	f := func(ctx context.Context) error {
		return impl.cli.GetDocs(
			ctx, filter,
			bson.M{}, &v,
		)
	}

	if err = primitive.WithContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = commonrepo.NewErrorResourceNotExists(err)
		}

		return
	}

	if len(v) == 0 {
		err = commonrepo.NewErrorResourceNotExists(fmt.Errorf("no member found"))
		return
	}

	members = make([]domain.OrgMember, len(v))
	for i := range v {
		members[i] = toOrgMember(&v[i])
	}
	return
}

func (impl *memberRepoImpl) DeleteByOrg(name string) (
	err error,
) {
	filter := memberDocOrgNameFilter(name)

	f := func(ctx context.Context) error {
		return impl.cli.DeleteMany(
			ctx, filter,
		)
	}

	if err = primitive.WithContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = commonrepo.NewErrorResourceNotExists(err)
		}

		return
	}

	return
}

func (impl *memberRepoImpl) GetByOrgAndUser(org, user string) (
	member domain.OrgMember, err error,
) {
	var v Member
	filter := memberDocInsertFilter(user, org)

	f := func(ctx context.Context) error {
		return impl.cli.GetDoc(
			ctx, filter,
			bson.M{}, &v,
		)
	}

	if err = primitive.WithContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = commonrepo.NewErrorResourceNotExists(err)
		}

		return
	}

	member = toOrgMember(&v)

	return
}

func (impl *memberRepoImpl) GetByOrgAndRole(org string, role domain.OrgRole) (members []domain.OrgMember, err error) {
	var v []Member
	filter := memberDocOrgAndRoleFilter(org, string(role))

	f := func(ctx context.Context) error {
		return impl.cli.GetDocs(
			ctx, filter,
			bson.M{}, &v,
		)
	}

	if err = primitive.WithContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = commonrepo.NewErrorResourceNotExists(err)
		}

		return
	}

	if len(v) == 0 {
		err = commonrepo.NewErrorResourceNotExists(fmt.Errorf("no member found"))
		return
	}

	members = make([]domain.OrgMember, len(v))
	for i := range v {
		members[i] = toOrgMember(&v[i])
	}

	return
}

func (impl *memberRepoImpl) GetByUser(name string) (
	members []domain.OrgMember, err error,
) {
	var v []Member
	filter := memberDocUserNameFilter(name)

	f := func(ctx context.Context) error {
		return impl.cli.GetDocs(
			ctx, filter,
			bson.M{}, &v,
		)
	}

	if err = primitive.WithContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = commonrepo.NewErrorResourceNotExists(err)
		}

		return
	}

	if len(v) == 0 {
		err = commonrepo.NewErrorResourceNotExists(fmt.Errorf("no member found"))
		return
	}

	members = make([]domain.OrgMember, len(v))
	for i := range v {
		members[i] = toOrgMember(&v[i])
	}
	return
}
