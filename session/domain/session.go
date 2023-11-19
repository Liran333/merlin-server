package session

import "github.com/openmerlin/merlin-server/user/domain"

type Session struct {
	Account domain.Account
	Info    string
	Email   domain.Email
	UserId  string
}

type SessionRepo interface {
	Save(*Session) error
	Get(domain.Account) (Session, error)
}
