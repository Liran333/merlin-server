package repositoryadapter

import (
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/space-app/domain"
)

const fieldSpaceId = "space_id"

var (
	spaceappTableName = ""
)

func toSpaceAppDO(m *domain.SpaceApp) spaceappDO {
	return spaceappDO{
		SpaceId:  m.SpaceId.Identity(),
		CommitId: m.CommitId,
		Status:   m.Status.AppStatus(),
	}
}

// spaceappDO
type spaceappDO struct {
	postgresql.CommonModel

	SpaceId  string `gorm:"column:space_id;index:,unique"`
	CommitId string `gorm:"column:commit_id"`

	Status string `gorm:"column:status"`

	AppURL    string `gorm:"column:app_url"`
	AppLogURL string `gorm:"column:app_log_url"`

	BuildLog    string `gorm:"column:build_log"`
	BuildLogURL string `gorm:"column:build_log_url"`

	Version int `gorm:"column:version"`
}

func (do *spaceappDO) TableName() string {
	return spaceappTableName
}
