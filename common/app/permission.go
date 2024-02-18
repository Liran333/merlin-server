package app

import (
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

var errorNoPermission = allerror.NewNoPermission("no permission")

// NewResourcePermissionAppService
func NewResourcePermissionAppService(
	org orgResourcePermissionValidator,
) *resourcePermissionAppService {
	return &resourcePermissionAppService{org: org}
}

// ResourcePermissionAppService
type ResourcePermissionAppService interface {
	CanRead(user primitive.Account, r Resource) error
	CanUpdate(user primitive.Account, r Resource) error
	CanDelete(user primitive.Account, r Resource) error
	CanCreate(user, owner primitive.Account, t primitive.ObjType) error
}

// orgResourcePermissionValidator
type orgResourcePermissionValidator interface {
	Check(primitive.Account, primitive.Account, primitive.ObjType, primitive.Action) error
}

// Resource
type Resource interface {
	OwnedBy(user primitive.Account) bool
	IsPublic() bool
	ResourceType() primitive.ObjType
	ResourceOwner() primitive.Account
	OwnedByPerson() bool
}

// resourcePermissionAppService
type resourcePermissionAppService struct {
	org orgResourcePermissionValidator
}

func (impl *resourcePermissionAppService) CanRead(user primitive.Account, r Resource) error {
	if r.IsPublic() {
		return nil
	}

	return impl.hasPermission(user, r, primitive.ActionRead)
}

func (impl *resourcePermissionAppService) CanUpdate(user primitive.Account, r Resource) error {
	return impl.hasPermission(user, r, primitive.ActionWrite)
}

func (impl *resourcePermissionAppService) CanDelete(user primitive.Account, r Resource) error {
	return impl.hasPermission(user, r, primitive.ActionDelete)
}

func (impl *resourcePermissionAppService) CanCreate(user, owner primitive.Account, t primitive.ObjType) error {
	if user == owner {
		return nil
	}

	return impl.org.Check(user, owner, t, primitive.ActionCreate)
}

func (impl *resourcePermissionAppService) hasPermission(user primitive.Account, r Resource, action primitive.Action) error {
	if user == nil {
		return errorNoPermission
	}

	if r.OwnedBy(user) {
		return nil
	}

	if r.OwnedByPerson() {
		return errorNoPermission
	}

	return impl.org.Check(user, r.ResourceOwner(), r.ResourceType(), action)
}
