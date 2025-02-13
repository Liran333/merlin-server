/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package spacerepositoryadapter

import (
	"github.com/lib/pq"

	coderepo "github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/space/domain"
	spaceprimitive "github.com/openmerlin/merlin-server/space/domain/primitive"
	"github.com/openmerlin/merlin-server/space/domain/repository"
)

const (
	filedId                = "id"
	fieldName              = "name"
	fieldTask              = "task"
	fieldHardwareType      = "hardware_type"
	fieldOwner             = "owner"
	fieldOthers            = "others"
	fieldLicense           = "license"
	fieldVersion           = "version"
	fieldFullName          = "fullname"
	fieldHardware          = "hardware"
	fieldLocalCMD          = "local_cmd"
	fieldUpdatedAt         = "updated_at"
	fieldCreatedAt         = "created_at"
	fieldVisibility        = "visibility"
	fieldBaseImage         = "base_image"
	fieldLikeCount         = "like_count"
	fieldDownloadCount     = "download_count"
	filedVisitCount        = "visit_count"
	fieldNoApplicationFile = "no_application_file"
)

var (
	spaceTableName = ""
)

func toSpaceDO(m *domain.Space) spaceDO {
	do := spaceDO{
		Id:                   m.Id.Integer(),
		SDK:                  m.SDK.SDK(),
		Desc:                 m.Desc.MSDDesc(),
		Name:                 m.Name.MSDName(),
		Owner:                m.Owner.Account(),
		License:              m.License.License(),
		AvatarId:             m.AvatarId.Storage(),
		Hardware:             m.Hardware.Hardware(),
		Fullname:             m.Fullname.MSDFullname(),
		Framework:            m.Labels.Framework,
		HardwareType:         m.Labels.HardwareType,
		BaseImage:            m.BaseImage.BaseImage(),
		CreatedBy:            m.CreatedBy.Account(),
		Visibility:           m.Visibility.Visibility(),
		Disable:              m.Disable,
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
		LikeCount:            m.LikeCount,
		Version:              m.Version,
		LocalCmd:             m.LocalCmd,
		LocalEnvInfo:         m.LocalEnvInfo,
		DownloadCount:        m.DownloadCount,
		VisitCount:           m.VisitCount,
		CompPowerAllocated:   m.CompPowerAllocated,
		NoApplicationFile:    m.NoApplicationFile,
		CommitId:             m.CommitId,
		IsDiscussionDisabled: m.IsDiscussionDisabled,
	}

	if m.DisableReason != nil {
		do.DisableReason = m.DisableReason.DisableReason()
	}

	if m.Exception != nil {
		do.Exception = m.Exception.Exception()
	}

	if m.Labels.Task != nil {
		do.Task = m.Labels.Task.Task()
	}

	if m.Labels.Licenses != nil {
		do.License = m.Labels.Licenses.License()
	}

	return do
}

func toSpaceStatisticDO(m *domain.Space) spaceDO {
	do := spaceDO{
		Id:            m.Id.Integer(),
		DownloadCount: m.DownloadCount,
		VisitCount:    m.VisitCount,
	}

	return do
}

func toLabelsDO(labels *domain.SpaceLabels) spaceDO {
	return spaceDO{
		Task:    labels.Task.Task(),
		License: labels.Licenses.License(),
	}
}

type spaceDO struct {
	Id            int64          `gorm:"column:id;"`
	SDK           string         `gorm:"column:sdk"`
	Desc          string         `gorm:"column:desc"`
	Name          string         `gorm:"column:name;index:space_index,unique,priority:2"`
	Owner         string         `gorm:"column:owner;index:space_index,unique,priority:1"`
	License       pq.StringArray `gorm:"column:license;type:text[];default:'{}';index:licenses,type:gin"`
	Hardware      string         `gorm:"column:hardware"`
	HardwareType  string         `gorm:"column:hardware_type"`
	Framework     string         `gorm:"column:framework"`
	Fullname      string         `gorm:"column:fullname"`
	AvatarId      string         `gorm:"column:avatar_id"`
	CreatedBy     string         `gorm:"column:created_by"`
	Visibility    string         `gorm:"column:visibility"`
	Disable       bool           `gorm:"column:disable"`
	DisableReason string         `gorm:"column:disable_reason"`
	Exception     string         `gorm:"column:exception"`
	CreatedAt     int64          `gorm:"column:created_at"`
	UpdatedAt     int64          `gorm:"column:updated_at"`
	Version       int            `gorm:"column:version"`
	LikeCount     int            `gorm:"column:like_count;not null;default:0"`
	DownloadCount int            `gorm:"column:download_count;not null;default:0"`
	VisitCount    int            `gorm:"column:visit_count;not null;default:0"`
	BaseImage     string         `gorm:"column:base_image"`
	// local cmd
	LocalCmd string `gorm:"column:local_cmd;type:text;default:'{}'"`
	// local EnvInfo
	LocalEnvInfo string `gorm:"column:local_envInfo;type:text;default:'{}'"`

	// labels
	Task string `gorm:"column:task;index:task"`

	// comp power allocated
	CompPowerAllocated bool `gorm:"column:comp_power_allocated"`
	// no application file
	NoApplicationFile bool `gorm:"column:no_application_file;default:t"`
	// latest commit id
	CommitId string `gorm:"column:commit_id"`

	IsDiscussionDisabled bool `gorm:"column:is_discussion_disabled"`
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
		Disable:       do.Disable,
		DisableReason: primitive.CreateDisableReason(do.DisableReason),
		Exception:     primitive.CreateException(do.Exception),
		SDK:           spaceprimitive.CreateSDK(do.SDK),
		Desc:          primitive.CreateMSDDesc(do.Desc),
		Fullname:      primitive.CreateMSDFullname(do.Fullname),
		Hardware:      spaceprimitive.CreateHardware(do.Hardware),
		BaseImage:     spaceprimitive.CreateBaseImage(do.BaseImage),
		AvatarId:      primitive.CreateAvatar(do.AvatarId),
		CreatedAt:     do.CreatedAt,
		UpdatedAt:     do.UpdatedAt,
		Version:       do.Version,
		LocalCmd:      do.LocalCmd,
		LikeCount:     do.LikeCount,
		LocalEnvInfo:  do.LocalEnvInfo,
		DownloadCount: do.DownloadCount,
		VisitCount:    do.VisitCount,
		Labels: domain.SpaceLabels{
			Task:         spaceprimitive.CreateTask(do.Task),
			Licenses:     primitive.CreateLicense(do.License),
			Framework:    do.Framework,
			HardwareType: do.HardwareType,
		},
		CompPowerAllocated:   do.CompPowerAllocated,
		NoApplicationFile:    do.NoApplicationFile,
		CommitId:             do.CommitId,
		IsDiscussionDisabled: do.IsDiscussionDisabled,
	}
}

func (do *spaceDO) toSpaceSummary() repository.SpaceSummary {
	return repository.SpaceSummary{
		Id:            primitive.CreateIdentity(do.Id).Identity(),
		Name:          do.Name,
		Desc:          do.Desc,
		Owner:         do.Owner,
		Fullname:      do.Fullname,
		BaseImage:     do.BaseImage,
		AvatarId:      primitive.CreateAvatar(do.AvatarId).URL(),
		UpdatedAt:     do.UpdatedAt,
		LikeCount:     do.LikeCount,
		DownloadCount: do.DownloadCount,
		VisitCount:    do.VisitCount,
		Disable:       do.Disable,
		DisableReason: do.DisableReason,
		Labels: domain.SpaceLabels{
			Task:      spaceprimitive.CreateTask(do.Task),
			Licenses:  primitive.CreateLicense(do.License),
			Framework: do.Framework,
		},
		IsNpu:              spaceprimitive.CreateHardware(do.Hardware).IsNpu(),
		Exception:          do.Exception,
		CompPowerAllocated: do.CompPowerAllocated,
		NoApplicationFile:  do.NoApplicationFile,
	}
}
