/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package repositoryadapter provides an adapter implementation for working with the repository of space applications.
package repositoryadapter

import (
	"fmt"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/spaceapp/domain"
	appprimitive "github.com/openmerlin/merlin-server/spaceapp/domain/primitive"
)

const (
	fieldSpaceId     = "space_id"
	fieldCommitId    = "commit_id"
	fieldVersion     = "version"
	fieldAllBuildLog = "all_build_log"
)

var (
	spaceappTableName = ""
)

func toSpaceAppDO(m *domain.SpaceApp) spaceappDO {
	do := spaceappDO{
		Status:      m.Status.AppStatus(),
		SpaceId:     m.SpaceId.Integer(),
		Version:     m.Version,
		CommitId:    m.CommitId,
		Reason:      m.Reason,
		RestartedAt: m.RestartedAt,
		ResumedAt:   m.ResumedAt,
	}

	if m.Id != nil {
		do.Id = m.Id.Integer()
	}

	if m.AppURL != nil {
		do.AppURL = m.AppURL.AppURL()
	}

	if m.AppLogURL != nil {
		do.AppLogURL = m.AppLogURL.URL()
	}

	if m.BuildLogURL != nil {
		do.BuildLogURL = m.BuildLogURL.URL()
	}

	return do
}

// spaceappDO
type spaceappDO struct {
	Id       int64  `gorm:"primarykey"`
	SpaceId  int64  `gorm:"column:space_id;index:,unique"`
	CommitId string `gorm:"column:commit_id"`

	Status string `gorm:"column:status"`
	Reason string `gorm:"column:reason"`

	RestartedAt int64 `gorm:"column:restarted_at"`
	ResumedAt   int64 `gorm:"column:resumed_at"`

	AppURL    string `gorm:"column:app_url"`
	AppLogURL string `gorm:"column:app_log_url"`

	AllBuildLog string `gorm:"column:all_build_log;type:text;"`
	BuildLogURL string `gorm:"column:build_log_url"`

	Version int `gorm:"column:version"`
}

// TableName returns the name of the table for the spaceappDO struct.
func (do *spaceappDO) TableName() string {
	return spaceappTableName
}

func (do *spaceappDO) toSpaceApp() domain.SpaceApp {
	v := domain.SpaceApp{
		Id: primitive.CreateIdentity(do.Id),
		SpaceAppIndex: domain.SpaceAppIndex{
			SpaceId:  primitive.CreateIdentity(do.SpaceId),
			CommitId: do.CommitId,
		},
		Status:      appprimitive.CreateAppStatus(do.Status),
		Reason:      do.Reason,
		RestartedAt: do.RestartedAt,
		ResumedAt:   do.ResumedAt,
		Version:     do.Version,
	}

	if do.AppURL != "" {
		v.AppURL = appprimitive.CreateAppURL(do.AppURL)
	}

	if do.AppLogURL != "" {
		v.AppLogURL = primitive.CreateURL(do.AppLogURL)
	}

	if do.BuildLogURL != "" {
		v.BuildLogURL = primitive.CreateURL(do.BuildLogURL)
	}

	return v
}

func equalQuery(field string) string {
	return fmt.Sprintf(`%s = ?`, field)
}
