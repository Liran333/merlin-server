package app

import (
	"context"

	coderedomain "github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/coderepo/domain/resourceadapter"
	"github.com/openmerlin/merlin-server/common/app"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

type resourcePermission struct {
	resource   resourceadapter.ResourceAdapter
	permission app.ResourcePermissionAppService
}

func (rp *resourcePermission) CanRead(ctx context.Context, resourceId primitive.Identity, user primitive.Account,
) (r coderedomain.Resource, err error) {
	r, err = rp.resource.GetByIndex(resourceId)
	if err != nil {
		err = allerror.New(allerror.ErrorCodeRepoNotFound, "resource not found", err)

		return
	}

	if r.DiscussionDisabled() {
		err = discussionDisabledErr

		return
	}

	err = rp.permission.CanRead(ctx, user, r)

	return
}

func (rp *resourcePermission) CanUpdate(ctx context.Context, resourceId primitive.Identity, user primitive.Account,
) error {
	r, err := rp.resource.GetByIndex(resourceId)
	if err != nil {
		return allerror.New(allerror.ErrorCodeRepoNotFound, "resource not found", err)
	}

	if r.DiscussionDisabled() {
		return discussionDisabledErr
	}

	return rp.permission.CanUpdate(ctx, user, r)
}
