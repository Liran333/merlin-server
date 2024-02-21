package domain

import "github.com/openmerlin/merlin-server/common/domain/primitive"

// Resource
type Resource interface {
	IsPublic() bool
	IsCreatedBy(user primitive.Account) bool
	ResourceType() primitive.ObjType
	ResourceOwner() primitive.Account
	OwnedByPerson() bool
}

type CodeRepoIndex struct {
	Name  primitive.MSDName
	Owner primitive.Account
}
