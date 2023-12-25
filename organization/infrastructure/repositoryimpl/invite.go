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

func NewInviteRepo(m mongodbClient) repository.Approve {
	return &inviteRepoImpl{m}
}

type inviteRepoImpl struct {
	cli mongodbClient
}

func (impl *inviteRepoImpl) ListInvitation(cmd *domain.OrgInvitationListCmd) (approves []domain.Approve, err error) {
	var v []Approve

	filter := bson.M{}
	if cmd.Org != nil {
		filter[fieldOrg] = cmd.Org.Account()
	}

	if cmd.Invitee != nil {
		filter[fieldInvitee] = cmd.Invitee.Account()
	} else if cmd.Inviter != nil {
		filter[fieldInviter] = cmd.Inviter.Account()
	}

	if cmd.Status != "" {
		filter[fieldStatus] = string(cmd.Status)
	}

	f := func(ctx context.Context) error {
		return impl.cli.GetDocs(
			ctx, filter,
			bson.M{}, &v,
		)
	}

	if err = primitive.WithContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = nil
		}

		return
	}

	approves = make([]domain.Approve, len(v))
	for i := range v {
		approves[i] = toApprove(&v[i])
	}

	return

}

func (impl *inviteRepoImpl) Save(o *domain.Approve) (r domain.Approve, err error) {
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

func (impl *inviteRepoImpl) update(o *domain.Approve) (err error) {
	org := ToApproveDoc(*o)
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

func approveDocFilter(org, invitee, inviter, status string) bson.M {
	return bson.M{
		fieldInviter: inviter,
		fieldInvitee: invitee,
		fieldOrg:     org,
		fieldStatus:  status,
	}
}

func (impl *inviteRepoImpl) insert(o *domain.Approve) (id string, err error) {
	org := ToApproveDoc(*o)
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
			ctx, approveDocFilter(org.Orgname, org.Username, org.Inviter, org.Status), doc,
		)

		id = v

		return err
	}

	if err = primitive.WithContext(f); err != nil && impl.cli.IsDocExists(err) {
		err = commonrepo.NewErrorDuplicateCreating(err)
	}

	return
}

func (impl *inviteRepoImpl) DeleteByOrg(acc primitive.Account) (err error) {
	if acc == nil {
		return fmt.Errorf("invalid org name when deleting invitation")
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
