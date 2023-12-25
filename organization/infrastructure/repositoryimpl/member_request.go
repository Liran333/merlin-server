package repositoryimpl

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	mongo "github.com/openmerlin/merlin-server/common/infrastructure/mongo"
	"github.com/openmerlin/merlin-server/organization/domain"
	"github.com/openmerlin/merlin-server/organization/domain/repository"
)

func NewMemberRequestRepo(m mongodbClient) repository.MemberRequest {
	return &requestRepoImpl{m}
}

type requestRepoImpl struct {
	cli mongodbClient
}

func (impl *requestRepoImpl) ListInvitation(cmd *domain.OrgMemberReqListCmd) (MemberRequests []domain.MemberRequest, err error) {
	var v []MemberRequest

	filter := bson.M{}
	if cmd.Org != nil {
		filter[fieldOrg] = cmd.Org.Account()
	}

	if cmd.Requester != nil {
		filter[fieldInvitee] = cmd.Requester.Account()
	}

	if cmd.Status != "" {
		filter[fieldStatus] = cmd.Status
	}

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

	MemberRequests = make([]domain.MemberRequest, len(v))
	for i := range v {
		MemberRequests[i] = toMemberRequest(&v[i])
	}

	return

}

func (impl *requestRepoImpl) Save(o *domain.MemberRequest) (r domain.MemberRequest, err error) {
	if o.Id != "" {
		if err = impl.update(o); err == nil {
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

func (impl *requestRepoImpl) update(o *domain.MemberRequest) (err error) {
	org := toMemberRequestDoc(*o)
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

	if err = primitive.WithContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = commonrepo.NewErrorConcurrentUpdating(err)
		}
	}

	return
}

func MemberRequestDocFilter(org, invitee, status string) bson.M {
	return bson.M{
		fieldInvitee: invitee,
		fieldOrg:     org,
		fieldStatus:  status,
	}
}

func (impl *requestRepoImpl) insert(o *domain.MemberRequest) (id string, err error) {
	org := toMemberRequestDoc(*o)
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
			ctx, MemberRequestDocFilter(org.Orgname, org.Username, org.Status), doc,
		)

		id = v

		return err
	}

	if err = primitive.WithContext(f); err != nil && impl.cli.IsDocExists(err) {
		err = commonrepo.NewErrorDuplicateCreating(err)
	}

	return
}

func (impl *requestRepoImpl) DeleteByOrg(acc primitive.Account) (err error) {
	if acc == nil {
		return fmt.Errorf("invalid org name when deleting member requests")
	}

	filter := bson.M{}
	filter[fieldOrg] = acc.Account()

	f := func(ctx context.Context) error {
		return impl.cli.DeleteMany(ctx, filter)
	}

	if err = primitive.WithContext(f); err != nil && impl.cli.IsDocExists(err) {
		err = nil
	}

	return
}
