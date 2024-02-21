package app

import (
	"github.com/openmerlin/merlin-server/common/domain"
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
	CanRead(user primitive.Account, r domain.Resource) error
	CanUpdate(user primitive.Account, r domain.Resource) error
	CanDelete(user primitive.Account, r domain.Resource) error
	CanCreate(user, owner primitive.Account, t primitive.ObjType) error
	CanListOrgResource(primitive.Account, primitive.Account, primitive.ObjType) error
}

// orgResourcePermissionValidator
type orgResourcePermissionValidator interface {
	Check(primitive.Account, primitive.Account, primitive.ObjType, primitive.Action, bool) error
}

// resourcePermissionAppService
type resourcePermissionAppService struct {
	org orgResourcePermissionValidator
}

func (impl *resourcePermissionAppService) CanListOrgResource(
	user, owner primitive.Account, t primitive.ObjType,
) error {
	return impl.org.Check(user, owner, t, primitive.ActionRead, true)
}

func (impl *resourcePermissionAppService) CanRead(user primitive.Account, r domain.Resource) error {
	if r.IsPublic() {
		return nil
	}

	// can't access private resource anonymously
	if user == nil {
		return errorNoPermission
	}

	// my own resource
	if user == r.ResourceOwner() {
		return nil
	}

	// can't access other individual's private resource
	if r.OwnedByPerson() {
		return errorNoPermission
	}

	return impl.org.Check(user, r.ResourceOwner(), r.ResourceType(), primitive.ActionRead, true)
}

func (impl *resourcePermissionAppService) CanUpdate(user primitive.Account, r domain.Resource) error {
	return impl.canModify(user, r, primitive.ActionWrite)
}

func (impl *resourcePermissionAppService) CanDelete(user primitive.Account, r domain.Resource) error {
	return impl.canModify(user, r, primitive.ActionDelete)
}

func (impl *resourcePermissionAppService) CanCreate(user, owner primitive.Account, t primitive.ObjType) error {
	if user == owner {
		return nil
	}

	return impl.org.Check(user, owner, t, primitive.ActionCreate, true)
}

func (impl *resourcePermissionAppService) canModify(
	user primitive.Account, r domain.Resource, action primitive.Action,
) error {
	// can't modify resource anonymously
	if user == nil {
		return errorNoPermission
	}

	// my own resource
	if user == r.ResourceOwner() {
		return nil
	}

	// can't modify other individual's resource
	if r.OwnedByPerson() {
		return errorNoPermission
	}

	return impl.org.Check(user, r.ResourceOwner(), r.ResourceType(), action, r.IsCreatedBy(user))
}
