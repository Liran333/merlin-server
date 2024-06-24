/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package domain provides domain for datasets.
package domain

import (
	"k8s.io/apimachinery/pkg/util/sets"

	coderepo "github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

// Dataset represents a dataset with its associated metadata and labels.
type Dataset struct {
	coderepo.CodeRepo

	Desc     primitive.MSDDesc
	Labels   DatasetLabels
	Fullname primitive.MSDFullname

	Version       int
	CreatedAt     int64
	UpdatedAt     int64
	LikeCount     int
	DownloadCount int

	Disable       bool
	DisableReason primitive.DisableReason
}

// ResourceType returns the type of the dataset resource.
func (m *Dataset) ResourceType() primitive.ObjType {
	return primitive.ObjTypeDataset
}

// IsDisable checks if the dataset is disable.
func (m *Dataset) IsDisable() bool {
	return m.Disable
}

// DatasetLabels represents the labels associated with a dataset, including task labels and other labels.
type DatasetLabels struct {
	Task     sets.Set[string] // task label
	License  sets.Set[string] // license label
	Size     string           // Size label
	Language sets.Set[string] // Language label
	Domain   sets.Set[string] // Domain label
}

// DatasetIndex represents the index for dataset in the code repository.
type DatasetIndex = coderepo.CodeRepoIndex
