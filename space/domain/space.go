/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package domain

import (
	"k8s.io/apimachinery/pkg/util/sets"

	coderepo "github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	spaceprimitive "github.com/openmerlin/merlin-server/space/domain/primitive"
)

// Space represents a space with its associated properties and methods.
type Space struct {
	coderepo.CodeRepo

	SDK      spaceprimitive.SDK
	Desc     primitive.MSDDesc
	Labels   SpaceLabels
	Fullname primitive.MSDFullname
	Hardware spaceprimitive.Hardware

	Version       int
	CreatedAt     int64
	UpdatedAt     int64
	LikeCount     int
	DownloadCount int
}

// ResourceType returns the type of the model resource.
func (m *Space) ResourceType() primitive.ObjType {
	return primitive.ObjTypeSpace
}

// SpaceLabels represents labels associated with a space.
type SpaceLabels struct {
	Task       string           // task label
	Others     sets.Set[string] // other labels
	Frameworks sets.Set[string] // framework labels
}

// SpaceIndex represents an index for spaces in the code repository.
type SpaceIndex = coderepo.CodeRepoIndex
