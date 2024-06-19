/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package app provides application services for resource permission management.
package app

import (
	"context"
	"fmt"

	"github.com/openmerlin/merlin-server/common/domain"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	orgapp "github.com/openmerlin/merlin-server/organization/app"
)

// NewResourcePermissionAppService creates a new instance of the resourcePermissionAppService.
func NewResourcePermissionAppService(
	org orgResourcePermissionValidator,
	disableOrg orgapp.PrivilegeOrg,
) *resourcePermissionAppService {
	return &resourcePermissionAppService{
		org:        org,
		disableOrg: disableOrg,
	}
}

// ResourcePermissionAppService defines methods for checking resource permissions.
type ResourcePermissionAppService interface {
	CanRead(ctx context.Context, user primitive.Account, r domain.Resource) error
	CanUpdate(ctx context.Context, user primitive.Account, r domain.Resource) error
	CanDelete(ctx context.Context, user primitive.Account, r domain.Resource) error
	CanCreate(ctx context.Context, user, owner primitive.Account, t primitive.ObjType) error
	CanReadPrivate(ctx context.Context, user primitive.Account, r domain.Resource) error
	CanListOrgResource(context.Context, primitive.Account, primitive.Account, primitive.ObjType) error
}

// orgResourcePermissionValidator
type orgResourcePermissionValidator interface {
	Check(context.Context, primitive.Account, primitive.Account, primitive.ObjType, primitive.Action) error
}

// resourcePermissionAppService
type resourcePermissionAppService struct {
	org        orgResourcePermissionValidator
	disableOrg orgapp.PrivilegeOrg
}

// CanListOrgResource checks if the user has permission to list organization resources of a specific type.
func (impl *resourcePermissionAppService) CanListOrgResource(
	ctx context.Context, user, owner primitive.Account, t primitive.ObjType,
) error {
	return impl.org.Check(ctx, user, owner, t, primitive.ActionRead)
}

// CanRead checks if the user has permission to read the specified resource.
func (impl *resourcePermissionAppService) CanRead(
	ctx context.Context, user primitive.Account, r domain.Resource) error {
	if r.IsPublic() {
		return nil
	}
	return impl.CanReadPrivate(ctx, user, r)
}

// disable administrator can read model, space and repocode obj.
func (impl *resourcePermissionAppService) disableAdminCanRead(
	ctx context.Context, user primitive.Account, r domain.Resource) (err error) {
	if impl.disableOrg != nil {
		action, _ := orgapp.NewAction(string(orgapp.Disable))
		err = impl.disableOrg.Contains(ctx, user)
		if err == nil && impl.disableOrg.IsCanReadObj(action, r.ResourceType()) {
			return nil
		}
	}
	return allerror.NewNoPermission("no permission", fmt.Errorf("not config disable admin"))
}

// CanReadPrivate checks if the user has permission to read private the specified resource.
func (impl *resourcePermissionAppService) CanReadPrivate(
	ctx context.Context, user primitive.Account, r domain.Resource) error {
	err := impl.canReadPrivate(ctx, user, r)
	if err != nil {
		if err1 := impl.disableAdminCanRead(ctx, user, r); err1 == nil {
			return nil
		}
	}
	return err
}

// canReadPrivate checks if the user has permission to read private the specified resource.
func (impl *resourcePermissionAppService) canReadPrivate(
	ctx context.Context, user primitive.Account, r domain.Resource) error {
	// can't access private resource anonymously
	if user == nil {
		return allerror.NewNoPermission("no permission", fmt.Errorf("anno can not access private resource"))
	}

	// my own resource
	if user == r.ResourceOwner() {
		return nil
	}

	// can't access other individual's private resource
	if r.OwnedByPerson() {
		return allerror.NewNoPermission("no permission", fmt.Errorf("can't access other individual's private resource"))
	}

	return impl.org.Check(ctx, user, r.ResourceOwner(), r.ResourceType(), primitive.ActionRead)
}

// CanUpdate checks if the user has permission to update the specified resource.
func (impl *resourcePermissionAppService) CanUpdate(
	ctx context.Context, user primitive.Account, r domain.Resource) error {
	return impl.canModify(ctx, user, r, primitive.ActionWrite)
}

// CanDelete checks if the user has permission to delete the specified resource.
func (impl *resourcePermissionAppService) CanDelete(
	ctx context.Context, user primitive.Account, r domain.Resource) error {
	return impl.canModify(ctx, user, r, primitive.ActionDelete)
}

// CanCreate checks if the user has permission to create a resource of the specified type, owned by the specified owner.
func (impl *resourcePermissionAppService) CanCreate(
	ctx context.Context, user, owner primitive.Account, t primitive.ObjType) error {
	if user == owner {
		return nil
	}

	return impl.org.Check(ctx, user, owner, t, primitive.ActionCreate)
}

func (impl *resourcePermissionAppService) canModify(
	ctx context.Context, user primitive.Account, r domain.Resource, action primitive.Action,
) error {
	// can't modify resource anonymously
	if user == nil {
		return allerror.NewNoPermission("no permission", fmt.Errorf("can't modify resource anonymously"))
	}

	// my own resource
	if user == r.ResourceOwner() {
		return nil
	}

	// can't modify other individual's resource
	if r.OwnedByPerson() {
		return allerror.NewNoPermission("no permission", fmt.Errorf("can't modify other individual's resource"))
	}

	return impl.org.Check(ctx, user, r.ResourceOwner(), r.ResourceType(), action)
}
