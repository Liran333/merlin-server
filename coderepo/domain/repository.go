/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package domain

import (
	commondomain "github.com/openmerlin/merlin-server/common/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

// CodeRepo represents a code repository with its properties.
type CodeRepo struct {
	Id         primitive.Identity
	Name       primitive.MSDName
	Owner      primitive.Account
	License    primitive.License
	Visibility primitive.Visibility
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
	}
}

// CodeRepoIndex represents the index of a code repository.
type CodeRepoIndex struct {
	Name  primitive.MSDName
	Owner primitive.Account
}

// Resource represents a common resource.
type Resource = commondomain.Resource
