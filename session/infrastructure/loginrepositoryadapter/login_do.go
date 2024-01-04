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

func toLoginDO(m *domain.Login) loginDO {
	return loginDO{
		Id:        m.Id,
		IP:        m.IP,
		User:      m.User.Account(),
		IdToken:   m.IdToken,
		UserAgent: m.UserAgent.UserAgent(),
		CreatedAt: m.CreatedAt,
	}
}

type loginDO struct {
	Id        primitive.UUID `gorm:"column:id;type:uuid;primaryKey"`
	IP        string         `gorm:"column:ip"`
	User      string         `gorm:"column:account;index:login_user"` // column'name can't be user, because it is a buildin name of pg.
	IdToken   string         `gorm:"column:id_token"`
	UserAgent string         `gorm:"column:user_agent"`
	CreatedAt int64          `gorm:"column:created_at"`
}

func (do *loginDO) TableName() string {
	return loginTableName
}

func (do *loginDO) toLogin() domain.Login {
	return domain.Login{
		Id:        do.Id,
		IP:        do.IP,
		User:      primitive.CreateAccount(do.User),
		IdToken:   do.IdToken,
		UserAgent: primitive.CreateUserAgent(do.UserAgent),
		CreatedAt: do.CreatedAt,
	}
}
