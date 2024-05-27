/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package domain provides domain models and types for the code repository branch.
package domain

import (
	"github.com/openmerlin/go-sdk/gitea"

	commondomain "github.com/openmerlin/merlin-server/common/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

// CodeRepo represents a code repository with its properties.
type CodeRepo struct {
	Id         primitive.Identity
	Name       primitive.MSDName
	Owner      primitive.Account
	License    primitive.License
	CreatedBy  primitive.Account
	Visibility primitive.Visibility
}

// ResourceType returns the resource type of the code repository.
func (r *CodeRepo) ResourceType() primitive.ObjType {
	return primitive.ObjTypeCodeRepo
}

// IsPrivate checks if the code repository is private.
func (r *CodeRepo) IsPrivate() bool {
	return r.Visibility.IsPrivate()
}

// IsPublic checks if the code repository is public.
func (r *CodeRepo) IsPublic() bool {
	return r.Visibility.IsPublic()
}

// RepoIndex returns the index of the code repository.
func (r *CodeRepo) RepoIndex() CodeRepoIndex {
	return CodeRepoIndex{
		Name:  r.Name,
		Owner: r.Owner,
		Id:    r.Id,
	}
}

// ResourceOwner returns the owner of the model resource.
func (m *CodeRepo) ResourceOwner() primitive.Account {
	return m.Owner
}

// IsCreatedBy checks if the model is created by the given user.
func (m *CodeRepo) IsCreatedBy(user primitive.Account) bool {
	return m.CreatedBy == user
}

// OwnedByPerson checks if the model is owned by the same person who created it.
func (m *CodeRepo) OwnedByPerson() bool {
	return m.Owner == m.CreatedBy
}

// ResourceVisibility returns the visibility of the code repository.
func (m *CodeRepo) ResourceVisibility() primitive.Visibility {
	return m.Visibility
}

// ResourceLicense returns the license of the code repository.
func (m *CodeRepo) ResourceLicense() primitive.License {
	return m.License
}

// CodeRepoIndex represents the index of a code repository.
type CodeRepoIndex = commondomain.CodeRepoIndex

// Resource represents a common resource.
type Resource = commondomain.Resource

// Repository is an alias type for the gitea.Repository.
type Repository = gitea.Repository
