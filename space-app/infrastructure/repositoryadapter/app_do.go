package repositoryadapter

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/space-app/domain"
	appprimitive "github.com/openmerlin/merlin-server/space-app/domain/primitive"
)

const (
	fieldSpaceId = "space_id"
	fieldVersion = "version"
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
		AllBuildLog: m.AllBuildLog,
	}

	do.Id = m.Id

	if m.AppURL != nil {
		do.AppURL = m.AppURL.URL()
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

	AppURL    string `gorm:"column:app_url"`
	AppLogURL string `gorm:"column:app_log_url"`

	AllBuildLog string `gorm:"column:all_build_log"`
	BuildLogURL string `gorm:"column:build_log_url"`

	Version int `gorm:"column:version"`
}

func (do *spaceappDO) TableName() string {
	return spaceappTableName
}

func (do *spaceappDO) toSpaceApp() domain.SpaceApp {
	v := domain.SpaceApp{
		Id: do.Id,
		SpaceAppIndex: domain.SpaceAppIndex{
			SpaceId:  primitive.CreateIdentity(do.SpaceId),
			CommitId: do.CommitId,
		},
		Status:      appprimitive.CreateAppStatus(do.Status),
		Version:     do.Version,
		AllBuildLog: do.AllBuildLog,
	}

	if do.AppURL != "" {
		v.AppURL = primitive.CreateURL(do.AppURL)
	}

	if do.AppLogURL != "" {
		v.AppLogURL = primitive.CreateURL(do.AppLogURL)
	}

	if do.BuildLogURL != "" {
		v.BuildLogURL = primitive.CreateURL(do.BuildLogURL)
	}

	return v
}
