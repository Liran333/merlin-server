/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package domain provides domain models and types.
package domain

import "github.com/openmerlin/merlin-server/common/domain/primitive"

// Resource represents an interface for a resource with various methods.
type Resource interface {
	IsPublic() bool
	IsPrivate() bool
	IsCreatedBy(user primitive.Account) bool
	ResourceType() primitive.ObjType
	ResourceOwner() primitive.Account
	OwnedByPerson() bool
	RepoIndex() CodeRepoIndex
	IsDisable() bool
	ResourceVisibility() primitive.Visibility
	ResourceLicense() primitive.License
	DiscussionDisabled() bool
	CloseDiscussion() error
	ReopenDiscussion() error
}

// CodeRepoIndex represents a code repository index with a name and owner.
type CodeRepoIndex struct {
	Name  primitive.MSDName
	Owner primitive.Account
	Id    primitive.Identity
}
