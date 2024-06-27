package app

import (
	"context"

	coderepo "github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/coderepo/domain/resourceadapter"
	"github.com/openmerlin/merlin-server/common/app"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	datasetrepo "github.com/openmerlin/merlin-server/datasets/domain/repository"
	modelrepo "github.com/openmerlin/merlin-server/models/domain/repository"
	spacerepo "github.com/openmerlin/merlin-server/space/domain/repository"
)

type DiscussionService interface {
	CloseDiscussion(context.Context, primitive.Identity, primitive.Account) error
	ReopenDiscussion(context.Context, primitive.Identity, primitive.Account) error
}

func NewDiscussionService(
	r resourceadapter.ResourceAdapter,
	p app.ResourcePermissionAppService,
	m modelrepo.ModelRepositoryAdapter,
	s spacerepo.SpaceRepositoryAdapter,
	d datasetrepo.DatasetRepositoryAdapter,
) *discussionService {
	return &discussionService{
		resource:   r,
		permission: p,
		model:      m,
		space:      s,
		dataset:    d,
	}
}

type discussionService struct {
	resource   resourceadapter.ResourceAdapter
	permission app.ResourcePermissionAppService
	model      modelrepo.ModelRepositoryAdapter
	space      spacerepo.SpaceRepositoryAdapter
	dataset    datasetrepo.DatasetRepositoryAdapter
}

func (d *discussionService) CloseDiscussion(ctx context.Context, id primitive.Identity, user primitive.Account) error {
	op := func(r coderepo.Resource) error {
		return r.CloseDiscussion()
	}

	return d.updateResourceDiscussionSwitch(ctx, id, user, op)
}

func (d *discussionService) ReopenDiscussion(ctx context.Context, id primitive.Identity, user primitive.Account) error {
	op := func(r coderepo.Resource) error {
		return r.ReopenDiscussion()
	}

	return d.updateResourceDiscussionSwitch(ctx, id, user, op)
}

func (d *discussionService) updateResourceDiscussionSwitch(
	ctx context.Context, id primitive.Identity, user primitive.Account, op func(resource coderepo.Resource) error,
) error {
	r, err := d.resource.GetByIndex(id)
	if err != nil {
		return allerror.New(allerror.ErrorCodeRepoNotFound, "resource not found", err)
	}

	if err = d.permission.CanUpdate(ctx, user, r); err != nil {
		return err
	}

	if err = op(r); err != nil {
		return err
	}

	return d.resource.Save(r)
}
