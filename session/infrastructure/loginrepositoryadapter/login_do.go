/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package loginrepositoryadapter

import (
	"github.com/openmerlin/merlin-server/common/domain/crypto"
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

func toLoginDO(m *domain.Session, e crypto.Encrypter) (do loginDO, err error) {
	idToken := m.IdToken
	if e != nil {
		idToken, err = e.Encrypt(m.IdToken)
		if err != nil {
			return
		}
	}

	do = loginDO{
		Id:        m.Id.RandomId(),
		IP:        m.IP,
		User:      m.User.Account(),
		IdToken:   idToken,
		UserAgent: m.UserAgent.UserAgent(),
		CreatedAt: m.CreatedAt,
		UserId:    m.UserId,
	}

	return
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

func (do *loginDO) toLogin(e crypto.Encrypter) (s domain.Session, err error) {
	idToken := do.IdToken
	if e != nil {
		idToken, err = e.Decrypt(do.IdToken)
		if err != nil {
			return
		}
	}

	s = domain.Session{
		Id:        primitive.CreateRandomId(do.Id),
		IP:        do.IP,
		User:      primitive.CreateAccount(do.User),
		IdToken:   idToken,
		UserAgent: primitive.CreateUserAgent(do.UserAgent),
		CreatedAt: do.CreatedAt,
		UserId:    do.UserId,
	}

	return
}
