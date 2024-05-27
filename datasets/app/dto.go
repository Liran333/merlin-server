/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package app provides functionality for the application.
package app

import (
	coderepoapp "github.com/openmerlin/merlin-server/coderepo/app"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/datasets/domain"
	"github.com/openmerlin/merlin-server/datasets/domain/repository"
	"github.com/openmerlin/merlin-server/utils"
)

// CmdToCreateDataset is a struct that represents a command to create a dataset.
type CmdToCreateDataset struct {
	coderepoapp.CmdToCreateRepo

	Desc     primitive.MSDDesc
	Fullname primitive.MSDFullname
}

// CmdToUpdateDataset is a struct that represents a command to update a dataset.
type CmdToUpdateDataset struct {
	coderepoapp.CmdToUpdateRepo

	Desc     primitive.MSDDesc
	Fullname primitive.MSDFullname
}

func (cmd *CmdToUpdateDataset) toDataset(dataset *domain.Dataset) (b bool) {
	if v := cmd.Desc; v != nil && v != dataset.Desc {
		dataset.Desc = v
		b = true
	}

	if v := cmd.Fullname; v != nil && v != dataset.Fullname {
		dataset.Fullname = v
		b = true
	}

	if b {
		dataset.UpdatedAt = utils.Now()
	}

	return
}

// CmdToDisableDataset is a struct that represents a command to disable a dataset.
type CmdToDisableDataset struct {
	Disable       bool
	DisableReason primitive.DisableReason
}

func (cmd *CmdToDisableDataset) toDataset(dataset *domain.Dataset) {
	dataset.Disable = cmd.Disable
	dataset.DisableReason = cmd.DisableReason
	dataset.UpdatedAt = utils.Now()
}

// DatasetDTO is a struct that represents a data transfer object for a dataset.
type DatasetDTO struct {
	Id            string           `json:"id"`
	Name          string           `json:"name"`
	Desc          string           `json:"desc"`
	Owner         string           `json:"owner"`
	Labels        DatasetLabelsDTO `json:"labels"`
	Fullname      string           `json:"fullname"`
	CreatedAt     int64            `json:"created_at"`
	UpdatedAt     int64            `json:"updated_at"`
	LikeCount     int              `json:"like_count"`
	Visibility    string           `json:"visibility"`
	DownloadCount int              `json:"download_count"`
	Disable       bool             `json:"disable"`
	DisableReason string           `json:"disable_reason"`
}

// DatasetLabelsDTO is a struct that represents a data transfer object for dataset labels.
type DatasetLabelsDTO struct {
	Task     []string `json:"task"`
	License  string   `json:"license"`
	Size     string   `json:"size"`
	Language []string `json:"language"`
	Domain   []string `json:"domain"`
}

func toDatasetLabelsDTO(dataset *domain.Dataset) DatasetLabelsDTO {
	labels := &dataset.Labels

	return DatasetLabelsDTO{
		Task:     labels.Task.UnsortedList(),
		License:  dataset.License.License(),
		Size:     labels.Size,
		Language: labels.Language.UnsortedList(),
		Domain:   labels.Domain.UnsortedList(),
	}
}

func toDatasetDTO(dataset *domain.Dataset) DatasetDTO {
	dto := DatasetDTO{
		Id:            dataset.Id.Identity(),
		Name:          dataset.Name.MSDName(),
		Owner:         dataset.Owner.Account(),
		Labels:        toDatasetLabelsDTO(dataset),
		CreatedAt:     dataset.CreatedAt,
		UpdatedAt:     dataset.UpdatedAt,
		LikeCount:     dataset.LikeCount,
		Visibility:    dataset.Visibility.Visibility(),
		DownloadCount: dataset.DownloadCount,

		Disable:       dataset.Disable,
		DisableReason: dataset.DisableReason.DisableReason()}

	if dataset.Desc != nil {
		dto.Desc = dataset.Desc.MSDDesc()
	}

	if dataset.Fullname != nil {
		dto.Fullname = dataset.Fullname.MSDFullname()
	}

	return dto
}

// DatasetsDTO is a struct that represents a data transfer object for a list of datasets.
type DatasetsDTO struct {
	Total    int                         `json:"total"`
	Datasets []repository.DatasetSummary `json:"datasets"`
}

// CmdToListDatasets is a type alias for repository.ListOption,
// representing a command to list datasets.
type CmdToListDatasets = repository.ListOption

// CmdToResetLabels is a type alias for domain.DatasetsLabels,
// representing a command to reset dataset labels.
type CmdToResetLabels = domain.DatasetLabels

// CmdToUpdateStatistics is a type alias for domain.DatasetsLabels,
// representing a command to update datasets statistics.
type CmdToUpdateStatistics struct {
	DownloadCount int `json:"download_count"`
}
