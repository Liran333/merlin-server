/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package modelrepositoryadapter provides an adapter for the model repository
package modelrepositoryadapter

import (
	"github.com/lib/pq"
	"k8s.io/apimachinery/pkg/util/sets"

	coderepo "github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/models/domain"
	"github.com/openmerlin/merlin-server/models/domain/repository"
)

const (
	fieldId            = "id"
	fieldName          = "name"
	fieldTask          = "task"
	fieldOwner         = "owner"
	fieldOthers        = "others"
	fieldLicense       = "license"
	fieldVersion       = "version"
	fieldFullName      = "fullname"
	fieldUpdatedAt     = "updated_at"
	fieldCreatedAt     = "created_at"
	fieldVisibility    = "visibility"
	fieldFrameworks    = "frameworks"
	fieldHardwares     = "hardwares"
	fieldLanguages     = "languages"
	fieldLikeCount     = "like_count"
	filedLibraryName   = "library_name"
	fieldDownloadCount = "download_count"
	fieldUseInOpenmind = "use_in_openmind"
)

var (
	modelTableName = ""
)

func toModelDO(m *domain.Model) modelDO {
	do := modelDO{
		Id:                   m.Id.Integer(),
		Desc:                 m.Desc.MSDDesc(),
		Name:                 m.Name.MSDName(),
		Owner:                m.Owner.Account(),
		Licenses:             m.License.License(),
		Fullname:             m.Fullname.MSDFullname(),
		CreatedBy:            m.CreatedBy.Account(),
		Visibility:           m.Visibility.Visibility(),
		Disable:              m.Disable,
		CreatedAt:            m.CreatedAt,
		LikeCount:            m.LikeCount,
		UpdatedAt:            m.UpdatedAt,
		Version:              m.Version,
		DownloadCount:        m.DownloadCount,
		UseInOpenmind:        m.UseInOpenmind,
		IsDiscussionDisabled: m.IsDiscussionDisabled,
	}

	if m.DisableReason != nil {
		do.DisableReason = m.DisableReason.DisableReason()
	}

	return do
}

func toModelStatisticDO(m *domain.Model) modelDO {
	do := modelDO{
		Id:            m.Id.Integer(),
		DownloadCount: m.DownloadCount,
	}

	return do
}

func toModelUseInOpenmindDO(m *domain.Model) modelDO {
	return modelDO{
		Id:            m.Id.Integer(),
		UseInOpenmind: m.UseInOpenmind,
	}
}

func toLabelsDO(labels *domain.ModelLabels) modelDO {
	do := modelDO{
		Task:        labels.Task,
		LibraryName: labels.LibraryName,
	}

	if labels.Others != nil {
		do.Others = labels.Others.UnsortedList()
	}

	if labels.Licenses != nil {
		do.Licenses = labels.Licenses.UnsortedList()
	}
	if labels.Frameworks != nil {
		do.Frameworks = labels.Frameworks.UnsortedList()
	}

	if labels.Hardwares != nil {
		do.Hardwares = labels.Hardwares.UnsortedList()
	}

	if labels.Languages != nil {
		do.Languages = labels.Languages.UnsortedList()
	}
	return do
}

type modelDO struct {
	Id                   int64          `gorm:"column:id;"`
	Desc                 string         `gorm:"column:desc"`
	Name                 string         `gorm:"column:name;index:model_index,unique,priority:2"`
	Owner                string         `gorm:"column:owner;index:model_index,unique,priority:1"`
	Licenses             pq.StringArray `gorm:"column:license;type:text[];default:'{}';index:licenses,type:gin"`
	Fullname             string         `gorm:"column:fullname"`
	CreatedBy            string         `gorm:"column:created_by"`
	Visibility           string         `gorm:"column:visibility"`
	Disable              bool           `gorm:"column:disable"`
	DisableReason        string         `gorm:"column:disable_reason"`
	CreatedAt            int64          `gorm:"column:created_at"`
	UpdatedAt            int64          `gorm:"column:updated_at"`
	Version              int            `gorm:"column:version"`
	LikeCount            int            `gorm:"column:like_count;not null;default:0"`
	DownloadCount        int            `gorm:"column:download_count;not null;default:0"`
	IsDiscussionDisabled bool           `gorm:"column:is_discussion_disabled"`

	// labels
	Task        string         `gorm:"column:task;index:task"`
	LibraryName string         `gorm:"column:library_name"`
	Others      pq.StringArray `gorm:"column:others;type:text[];default:'{}';index:others,type:gin"`
	Frameworks  pq.StringArray `gorm:"column:frameworks;type:text[];default:'{}';index:frameworks,type:gin"`
	Hardwares   pq.StringArray `gorm:"column:hardwares;type:text[];default:'{}';index:hardwares,type:gin"`
	Languages   pq.StringArray `gorm:"column:languages;type:text[];default:'{}';index:languages,type:gin"`

	// for openmind
	UseInOpenmind string `gorm:"column:use_in_openmind"`
}

// TableName returns the table name of the model.
func (do *modelDO) TableName() string {
	return modelTableName
}

func (do *modelDO) toModel() domain.Model {
	return domain.Model{
		CodeRepo: coderepo.CodeRepo{
			Id:         primitive.CreateIdentity(do.Id),
			Name:       primitive.CreateMSDName(do.Name),
			Owner:      primitive.CreateAccount(do.Owner),
			License:    primitive.CreateLicense(do.Licenses),
			CreatedBy:  primitive.CreateAccount(do.CreatedBy),
			Visibility: primitive.CreateVisibility(do.Visibility),
		},
		Disable:              do.Disable,
		DisableReason:        primitive.CreateDisableReason(do.DisableReason),
		Desc:                 primitive.CreateMSDDesc(do.Desc),
		Fullname:             primitive.CreateMSDFullname(do.Fullname),
		CreatedAt:            do.CreatedAt,
		UpdatedAt:            do.UpdatedAt,
		Version:              do.Version,
		LikeCount:            do.LikeCount,
		DownloadCount:        do.DownloadCount,
		UseInOpenmind:        do.UseInOpenmind,
		IsDiscussionDisabled: do.IsDiscussionDisabled,

		Labels: domain.ModelLabels{
			Task:        do.Task,
			LibraryName: do.LibraryName,
			Others:      sets.New[string](do.Others...),
			Licenses:    sets.New[string](do.Licenses...),
			Frameworks:  sets.New[string](do.Frameworks...),
			Languages:   sets.New[string](do.Languages...),
			Hardwares:   sets.New[string](do.Hardwares...),
		},
	}
}

func (do *modelDO) toModelSummary() repository.ModelSummary {
	return repository.ModelSummary{
		Id:            primitive.CreateIdentity(do.Id).Identity(),
		Name:          do.Name,
		Desc:          do.Desc,
		Task:          do.Task,
		Owner:         do.Owner,
		Licenses:      do.Licenses,
		Fullname:      do.Fullname,
		UpdatedAt:     do.UpdatedAt,
		Frameworks:    do.Frameworks,
		LikeCount:     do.LikeCount,
		DownloadCount: do.DownloadCount,
		Disable:       do.Disable,
		DisableReason: do.DisableReason,
	}
}
