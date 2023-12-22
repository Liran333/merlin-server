package modelrepositoryadapter

// "gorm.io/plugin/optimisticlock"

import (
	"github.com/lib/pq"

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
	fieldVersion    = "version"
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
		Visibility: m.Visibility.Visibility(),
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
		Version:    m.Version,
	}
}

func toLabelsDO(labels *domain.ModelLabels) modelDO {
	return modelDO{
		Task:       labels.Task,
		Others:     labels.Others.UnsortedList(),
		Frameworks: labels.Frameworks.UnsortedList(),
	}
}

type modelDO struct {
	Id         int64  `gorm:"column:id;"`
	Desc       string `gorm:"column:desc"`
	Name       string `gorm:"column:name;index:model_index,unique,priority:2"`
	Owner      string `gorm:"column:owner;index:model_index,unique,priority:1"`
	License    string `gorm:"column:license"`
	Fullname   string `gorm:"column:fullname"`
	Visibility string `gorm:"column:visibility"`
	CreatedAt  int64  `gorm:"column:created_at"`
	UpdatedAt  int64  `gorm:"column:updated_at"`
	Version    int    `gorm:"column:version"`

	//labels
	Task       string         `gorm:"column:task"`
	Others     pq.StringArray `gorm:"column:others;type:text[];default:'{}'"`
	Frameworks pq.StringArray `gorm:"column:frameworks;type:text[];default:'{}'"`
}

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
			Visibility: primitive.CreateVisibility(do.Visibility),
		},
		Desc:      primitive.CreateMSDDesc(do.Desc),
		Fullname:  primitive.CreateMSDFullname(do.Fullname),
		CreatedAt: do.CreatedAt,
		UpdatedAt: do.UpdatedAt,
		Version:   do.Version,
	}
}

func (do *modelDO) toModelSummary() repository.ModelSummary {
	return repository.ModelSummary{
		Id:        primitive.CreateIdentity(do.Id).Identity(),
		Name:      do.Name,
		Desc:      do.Desc,
		Task:      do.Task,
		Owner:     do.Owner,
		Fullname:  do.Fullname,
		UpdatedAt: do.UpdatedAt,
	}
}
