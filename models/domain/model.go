/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package domain

import (
	"k8s.io/apimachinery/pkg/util/sets"

	coderepo "github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

// Model represents a model with its associated metadata and labels.
type Model struct {
	coderepo.CodeRepo

	Desc     primitive.MSDDesc
	Labels   ModelLabels
	Fullname primitive.MSDFullname

	Version       int
	CreatedAt     int64
	UpdatedAt     int64
	LikeCount     int
	DownloadCount int
}

// ResourceType returns the type of the model resource.
func (m *Model) ResourceType() primitive.ObjType {
	return primitive.ObjTypeModel
}

// ModelLabels represents the labels associated with a model, including task labels, other labels, and framework labels.
type ModelLabels struct {
	Task       string           // task label
	Others     sets.Set[string] // other labels
	Frameworks sets.Set[string] // framework labels
}

// ModelIndex represents the index for models in the code repository.
type ModelIndex = coderepo.CodeRepoIndex
