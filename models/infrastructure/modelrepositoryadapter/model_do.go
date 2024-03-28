/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

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
	fieldName       = "name"
	fieldTask       = "task"
	fieldOwner      = "owner"
	fieldOthers     = "others"
	fieldLicense    = "license"
	fieldVersion    = "version"
	fieldFullName   = "fullname"
	fieldUpdatedAt  = "updated_at"
	fieldCreatedAt  = "created_at"
	fieldVisibility = "visibility"
	fieldFrameworks = "frameworks"
)

var (
	modelTableName = ""
)

func toModelDO(m *domain.Model) modelDO {
	return modelDO{
		Id:         m.Id.Integer(),
		Desc:       m.Desc.MSDDesc(),
		Name:       m.Name.MSDName(),
		Owner:      m.Owner.Account(),
		License:    m.License.License(),
		Fullname:   m.Fullname.MSDFullname(),
		CreatedBy:  m.CreatedBy.Account(),
		Visibility: m.Visibility.Visibility(),
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
		Version:    m.Version,
	}
}

func toLabelsDO(labels *domain.ModelLabels) modelDO {
	do := modelDO{
		Task:    labels.Task,
		License: labels.License,
	}

	if labels.Others != nil {
		do.Others = labels.Others.UnsortedList()
	}

	if labels.Frameworks != nil {
		do.Frameworks = labels.Frameworks.UnsortedList()
	}

	return do
}

type modelDO struct {
	Id         int64  `gorm:"column:id;"`
	Desc       string `gorm:"column:desc"`
	Name       string `gorm:"column:name;index:model_index,unique,priority:2"`
	Owner      string `gorm:"column:owner;index:model_index,unique,priority:1"`
	License    string `gorm:"column:license"`
	Fullname   string `gorm:"column:fullname"`
	CreatedBy  string `gorm:"column:created_by"`
	Visibility string `gorm:"column:visibility"`
	CreatedAt  int64  `gorm:"column:created_at"`
	UpdatedAt  int64  `gorm:"column:updated_at"`
	Version    int    `gorm:"column:version"`

	// labels
	Task       string         `gorm:"column:task;index:task"`
	Others     pq.StringArray `gorm:"column:others;type:text[];default:'{}';index:others,type:gin"`
	Frameworks pq.StringArray `gorm:"column:frameworks;type:text[];default:'{}';index:frameworks,type:gin"`
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
			License:    primitive.CreateLicense(do.License),
			CreatedBy:  primitive.CreateAccount(do.CreatedBy),
			Visibility: primitive.CreateVisibility(do.Visibility),
		},
		Desc:      primitive.CreateMSDDesc(do.Desc),
		Fullname:  primitive.CreateMSDFullname(do.Fullname),
		CreatedAt: do.CreatedAt,
		UpdatedAt: do.UpdatedAt,
		Version:   do.Version,

		Labels: domain.ModelLabels{
			Task:       do.Task,
			Others:     sets.New[string](do.Others...),
			Frameworks: sets.New[string](do.Frameworks...),
		},
	}
}

func (do *modelDO) toModelSummary() repository.ModelSummary {
	return repository.ModelSummary{
		Id:         primitive.CreateIdentity(do.Id).Identity(),
		Name:       do.Name,
		Desc:       do.Desc,
		Task:       do.Task,
		Owner:      do.Owner,
		License:    do.License,
		Fullname:   do.Fullname,
		UpdatedAt:  do.UpdatedAt,
		Frameworks: do.Frameworks,
	}
}
