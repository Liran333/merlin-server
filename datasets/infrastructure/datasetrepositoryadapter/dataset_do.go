/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package datasetsrepositoryadapter provides an adapter for the dataset repository
package datasetrepositoryadapter

import (
	"github.com/lib/pq"
	"k8s.io/apimachinery/pkg/util/sets"

	coderepo "github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/datasets/domain"
	"github.com/openmerlin/merlin-server/datasets/domain/repository"
)

const (
	fieldId            = "id"
	fieldName          = "name"
	fieldTask          = "task"
	fieldOwner         = "owner"
	fieldLicense       = "license"
	fieldVersion       = "version"
	fieldFullName      = "fullname"
	fieldUpdatedAt     = "updated_at"
	fieldCreatedAt     = "created_at"
	fieldVisibility    = "visibility"
	fieldLikeCount     = "like_count"
	fieldDownloadCount = "download_count"
	fieldSize          = "size"
	fieldLanguage      = "language"
	fieldDomain        = "domain"
)

var (
	datasetTableName = ""
)

func toDatasetDO(m *domain.Dataset) datasetDO {
	do := datasetDO{
		Id:            m.Id.Integer(),
		Desc:          m.Desc.MSDDesc(),
		Name:          m.Name.MSDName(),
		Owner:         m.Owner.Account(),
		License:       m.License.License(),
		Fullname:      m.Fullname.MSDFullname(),
		CreatedBy:     m.CreatedBy.Account(),
		Visibility:    m.Visibility.Visibility(),
		Disable:       m.Disable,
		CreatedAt:     m.CreatedAt,
		LikeCount:     m.LikeCount,
		UpdatedAt:     m.UpdatedAt,
		Version:       m.Version,
		DownloadCount: m.DownloadCount,
	}

	if m.DisableReason != nil {
		do.DisableReason = m.DisableReason.DisableReason()
	}

	return do
}

func toDatasetStatisticDO(m *domain.Dataset) datasetDO {
	do := datasetDO{
		Id:            m.Id.Integer(),
		DownloadCount: m.DownloadCount,
	}

	if m.DisableReason != nil {
		do.DisableReason = m.DisableReason.DisableReason()
	}

	return do
}

func toLabelsDO(labels *domain.DatasetLabels) datasetDO {
	do := datasetDO{
		License: labels.License,
		Size:    labels.Size,
	}

	if labels.Task != nil {
		do.Task = labels.Task.UnsortedList()
	}

	if labels.Language != nil {
		do.Language = labels.Language.UnsortedList()
	}

	if labels.Domain != nil {
		do.Domain = labels.Domain.UnsortedList()
	}

	return do
}

type datasetDO struct {
	Id            int64  `gorm:"column:id;"`
	Desc          string `gorm:"column:desc"`
	Name          string `gorm:"column:name;index:dataset_index,unique,priority:2"`
	Owner         string `gorm:"column:owner;index:dataset_index,unique,priority:1"`
	License       string `gorm:"column:license"`
	Fullname      string `gorm:"column:fullname"`
	CreatedBy     string `gorm:"column:created_by"`
	Visibility    string `gorm:"column:visibility"`
	Disable       bool   `gorm:"column:disable"`
	DisableReason string `gorm:"column:disable_reason"`
	CreatedAt     int64  `gorm:"column:created_at"`
	UpdatedAt     int64  `gorm:"column:updated_at"`
	Version       int    `gorm:"column:version"`
	LikeCount     int    `gorm:"column:like_count;not null;default:0"`
	DownloadCount int    `gorm:"column:download_count;not null;default:0"`

	// labels
	Task     pq.StringArray `gorm:"column:task;type:text[];default:'{}';index:task,type:gin"`
	Language pq.StringArray `gorm:"column:language;type:text[];default:'{}';index:language,type:gin"`
	Domain   pq.StringArray `gorm:"column:domain;type:text[];default:'{}';index:domain,type:gin"`
	Size     string         `gorm:"column:size;index:size"`
}

// TableName returns the table name of the dataset.
func (do *datasetDO) TableName() string {
	return datasetTableName
}

func (do *datasetDO) toDataset() domain.Dataset {
	return domain.Dataset{
		CodeRepo: coderepo.CodeRepo{
			Id:         primitive.CreateIdentity(do.Id),
			Name:       primitive.CreateMSDName(do.Name),
			Owner:      primitive.CreateAccount(do.Owner),
			License:    primitive.CreateLicense(do.License),
			CreatedBy:  primitive.CreateAccount(do.CreatedBy),
			Visibility: primitive.CreateVisibility(do.Visibility),
		},
		Disable:       do.Disable,
		DisableReason: primitive.CreateDisableReason(do.DisableReason),
		Desc:          primitive.CreateMSDDesc(do.Desc),
		Fullname:      primitive.CreateMSDFullname(do.Fullname),
		CreatedAt:     do.CreatedAt,
		UpdatedAt:     do.UpdatedAt,
		Version:       do.Version,
		LikeCount:     do.LikeCount,
		DownloadCount: do.DownloadCount,

		Labels: domain.DatasetLabels{
			Task:     sets.New[string](do.Task...),
			Size:     do.Size,
			Language: sets.New[string](do.Language...),
			Domain:   sets.New[string](do.Domain...),
		},
	}
}

func (do *datasetDO) toDatasetSummary() repository.DatasetSummary {
	return repository.DatasetSummary{
		Id:            primitive.CreateIdentity(do.Id).Identity(),
		Name:          do.Name,
		Desc:          do.Desc,
		Task:          do.Task,
		Owner:         do.Owner,
		License:       do.License,
		Fullname:      do.Fullname,
		UpdatedAt:     do.UpdatedAt,
		LikeCount:     do.LikeCount,
		DownloadCount: do.DownloadCount,
		Disable:       do.Disable,
		DisableReason: do.DisableReason,
		Size:          do.Size,
		Language:      do.Language,
		Domain:        do.Domain,
	}
}
