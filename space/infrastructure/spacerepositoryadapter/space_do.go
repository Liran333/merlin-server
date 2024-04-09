/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package spacerepositoryadapter

import (
	"github.com/lib/pq"
	"k8s.io/apimachinery/pkg/util/sets"

	coderepo "github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/space/domain"
	spaceprimitive "github.com/openmerlin/merlin-server/space/domain/primitive"
	"github.com/openmerlin/merlin-server/space/domain/repository"
)

const (
	filedId            = "id"
	fieldName          = "name"
	fieldTask          = "task"
	fieldOwner         = "owner"
	fieldOthers        = "others"
	fieldLicense       = "license"
	fieldVersion       = "version"
	fieldFullName      = "fullname"
	fieldLocalCMD      = "local_cmd"
	fieldUpdatedAt     = "updated_at"
	fieldCreatedAt     = "created_at"
	fieldVisibility    = "visibility"
	fieldFrameworks    = "frameworks"
	fieldLikeCount     = "like_count"
	fieldDownloadCount = "download_count"
)

var (
	spaceTableName = ""
)

func toSpaceDO(m *domain.Space) spaceDO {
	return spaceDO{
		Id:         m.Id.Integer(),
		SDK:        m.SDK.SDK(),
		Desc:       m.Desc.MSDDesc(),
		Name:       m.Name.MSDName(),
		Owner:      m.Owner.Account(),
		License:    m.License.License(),
		Hardware:   m.Hardware.Hardware(),
		Fullname:   m.Fullname.MSDFullname(),
		CreatedBy:  m.CreatedBy.Account(),
		Visibility: m.Visibility.Visibility(),
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
		LikeCount:  m.LikeCount,
		Version:    m.Version,
		LocalCmd:   m.LocalCmd,
		LocalEnvInfo:   m.LocalEnvInfo,
	}
}

func toLabelsDO(labels *domain.SpaceLabels) spaceDO {
	return spaceDO{
		Task:       labels.Task,
		Others:     labels.Others.UnsortedList(),
		Frameworks: labels.Frameworks.UnsortedList(),
	}
}

type spaceDO struct {
	Id            int64  `gorm:"column:id;"`
	SDK           string `gorm:"column:sdk"`
	Desc          string `gorm:"column:desc"`
	Name          string `gorm:"column:name;index:space_index,unique,priority:2"`
	Owner         string `gorm:"column:owner;index:space_index,unique,priority:1"`
	License       string `gorm:"column:license"`
	Hardware      string `gorm:"column:hardware"`
	Fullname      string `gorm:"column:fullname"`
	CreatedBy     string `gorm:"column:created_by"`
	Visibility    string `gorm:"column:visibility"`
	CreatedAt     int64  `gorm:"column:created_at"`
	UpdatedAt     int64  `gorm:"column:updated_at"`
	Version       int    `gorm:"column:version"`
	LikeCount     int    `gorm:"column:like_count;not null;default:0"`
	DownloadCount int    `gorm:"column:download_count;not null;default:0"`

	// local cmd
	LocalCmd string `gorm:"column:local_cmd;type:text;default:'{}'"`
	// local EnvInfo
	LocalEnvInfo string `gorm:"column:local_envInfo;type:text;default:'{}'"`

	// labels
	Task       string         `gorm:"column:task;index:task"`
	Others     pq.StringArray `gorm:"column:others;type:text[];default:'{}';index:others,type:gin"`
	Frameworks pq.StringArray `gorm:"column:frameworks;type:text[];default:'{}';index:frameworks,type:gin"`
}

// TableName returns the table name of spaceDO.
func (do *spaceDO) TableName() string {
	return spaceTableName
}

func (do *spaceDO) toSpace() domain.Space {
	return domain.Space{
		CodeRepo: coderepo.CodeRepo{
			Id:         primitive.CreateIdentity(do.Id),
			Name:       primitive.CreateMSDName(do.Name),
			Owner:      primitive.CreateAccount(do.Owner),
			License:    primitive.CreateLicense(do.License),
			Visibility: primitive.CreateVisibility(do.Visibility),
			CreatedBy:  primitive.CreateAccount(do.CreatedBy),
		},
		SDK:       spaceprimitive.CreateSDK(do.SDK),
		Desc:      primitive.CreateMSDDesc(do.Desc),
		Fullname:  primitive.CreateMSDFullname(do.Fullname),
		Hardware:  spaceprimitive.CreateHardware(do.Hardware),
		CreatedAt: do.CreatedAt,
		UpdatedAt: do.UpdatedAt,
		Version:   do.Version,
		LocalCmd:  do.LocalCmd,
		LikeCount: do.LikeCount,
		LocalEnvInfo:  do.LocalEnvInfo,
		Labels: domain.SpaceLabels{
			Task:       do.Task,
			Others:     sets.New[string](do.Others...),
			Frameworks: sets.New[string](do.Frameworks...),
		},
	}
}

func (do *spaceDO) toSpaceSummary() repository.SpaceSummary {
	return repository.SpaceSummary{
		Id:            primitive.CreateIdentity(do.Id).Identity(),
		Name:          do.Name,
		Desc:          do.Desc,
		Task:          do.Task,
		Owner:         do.Owner,
		Fullname:      do.Fullname,
		UpdatedAt:     do.UpdatedAt,
		LikeCount:     do.LikeCount,
		DownloadCount: do.DownloadCount,
	}
}
