/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package app provides application services for resource permission management.
package app

import (
	"fmt"

	"github.com/openmerlin/merlin-server/common/domain"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

// NewResourcePermissionAppService creates a new instance of the resourcePermissionAppService.
func NewResourcePermissionAppService(
	org orgResourcePermissionValidator,
) *resourcePermissionAppService {
	return &resourcePermissionAppService{org: org}
}

// ResourcePermissionAppService defines methods for checking resource permissions.
type ResourcePermissionAppService interface {
	CanRead(user primitive.Account, r domain.Resource) error
	CanUpdate(user primitive.Account, r domain.Resource) error
	CanDelete(user primitive.Account, r domain.Resource) error
	CanCreate(user, owner primitive.Account, t primitive.ObjType) error
	CanListOrgResource(primitive.Account, primitive.Account, primitive.ObjType) error
}

// orgResourcePermissionValidator
type orgResourcePermissionValidator interface {
	Check(primitive.Account, primitive.Account, primitive.ObjType, primitive.Action) error
}

// resourcePermissionAppService
type resourcePermissionAppService struct {
	org orgResourcePermissionValidator
}

// CanListOrgResource checks if the user has permission to list organization resources of a specific type.
func (impl *resourcePermissionAppService) CanListOrgResource(
	user, owner primitive.Account, t primitive.ObjType,
) error {
	return impl.org.Check(user, owner, t, primitive.ActionRead)
}

// CanRead checks if the user has permission to read the specified resource.
func (impl *resourcePermissionAppService) CanRead(user primitive.Account, r domain.Resource) error {
	if r.IsPublic() {
		return nil
	}

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

	return impl.org.Check(user, r.ResourceOwner(), r.ResourceType(), primitive.ActionRead)
}

// CanUpdate checks if the user has permission to update the specified resource.
func (impl *resourcePermissionAppService) CanUpdate(user primitive.Account, r domain.Resource) error {
	return impl.canModify(user, r, primitive.ActionWrite)
}

// CanDelete checks if the user has permission to delete the specified resource.
func (impl *resourcePermissionAppService) CanDelete(user primitive.Account, r domain.Resource) error {
	return impl.canModify(user, r, primitive.ActionDelete)
}

// CanCreate checks if the user has permission to create a resource of the specified type, owned by the specified owner.
func (impl *resourcePermissionAppService) CanCreate(user, owner primitive.Account, t primitive.ObjType) error {
	if user == owner {
		return nil
	}

	return impl.org.Check(user, owner, t, primitive.ActionCreate)
}

func (impl *resourcePermissionAppService) canModify(
	user primitive.Account, r domain.Resource, action primitive.Action,
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

	return impl.org.Check(user, r.ResourceOwner(), r.ResourceType(), action)
}
