/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package loginrepositoryadapter

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/session/domain"
)

var (
	loginTableName = ""
)

const (
	fieldUser      = "account"
	fieldCreatedAt = "created_at"
)

func toLoginDO(m *domain.Session) loginDO {
	return loginDO{
		Id:        m.Id.RandomId(),
		IP:        m.IP,
		User:      m.User.Account(),
		IdToken:   m.IdToken,
		UserAgent: m.UserAgent.UserAgent(),
		CreatedAt: m.CreatedAt,
		UserId:    m.UserId,
	}
}

type loginDO struct {
	Id string `gorm:"column:id;primaryKey"`
	IP string `gorm:"column:ip"`
	// column'name can't be user, because it is a buildin name of pg.
	User      string `gorm:"column:account;index:login_user"`
	IdToken   string `gorm:"column:id_token"`
	UserAgent string `gorm:"column:user_agent"`
	CreatedAt int64  `gorm:"column:created_at"`
	UserId    string `gorm:"column:user_id"` // user id in OIDC provider
}

// TableName returns the table name for the loginDO struct.
func (do *loginDO) TableName() string {
	return loginTableName
}

func (do *loginDO) toLogin() domain.Session {
	return domain.Session{
		Id:        primitive.CreateRandomId(do.Id),
		IP:        do.IP,
		User:      primitive.CreateAccount(do.User),
		IdToken:   do.IdToken,
		UserAgent: primitive.CreateUserAgent(do.UserAgent),
		CreatedAt: do.CreatedAt,
		UserId:    do.UserId,
	}
}
