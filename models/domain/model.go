/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package domain provides domain for models.
package domain

import (
	"golang.org/x/xerrors"
	"k8s.io/apimachinery/pkg/util/sets"

	coderepo "github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

// Model represents a model with its associated metadata and labels.
type Model struct {
	coderepo.CodeRepo

	Desc          primitive.MSDDesc
	Labels        ModelLabels
	Fullname      primitive.MSDFullname
	UseInOpenmind string

	Version       int
	CreatedAt     int64
	UpdatedAt     int64
	LikeCount     int
	DownloadCount int

	Disable       bool
	DisableReason primitive.DisableReason

	IsDiscussionDisabled bool
}

// ResourceType returns the type of the model resource.
func (m *Model) ResourceType() primitive.ObjType {
	return primitive.ObjTypeModel
}

// IsDisable checks if the space is disable.
func (m *Model) IsDisable() bool {
	return m.Disable
}

func (m *Model) DiscussionDisabled() bool {
	return m.IsDiscussionDisabled
}

func (m *Model) CloseDiscussion() error {
	if m.IsDiscussionDisabled {
		return allerror.New(
			allerror.ErrorCodeDiscussionDisabled,
			"failed to close discussion",
			xerrors.Errorf("discussion is closed"),
		)
	}

	m.IsDiscussionDisabled = true

	return nil
}

func (m *Model) ReopenDiscussion() error {
	if !m.IsDiscussionDisabled {
		return allerror.New(
			allerror.ErrorCodeDiscussionEnabled,
			"failed to reopen discussion",
			xerrors.Errorf("discussion is open"),
		)
	}

	m.IsDiscussionDisabled = false

	return nil
}

// ModelLabels represents the labels associated with a model, including task labels, other labels, and framework labels.
type ModelLabels struct {
	Task        string           // task label
	LibraryName string           // library label
	Licenses    sets.Set[string] // license label
	Others      sets.Set[string] // other labels
	Frameworks  sets.Set[string] // framework labels
	Hardwares   sets.Set[string] // hardware label
	Languages   sets.Set[string] // language labels
}

// ModelIndex represents the index for models in the code repository.
type ModelIndex = coderepo.CodeRepoIndex
