package session

import (
	"fmt"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	sessiondomain "github.com/openmerlin/merlin-server/session/domain"
	"github.com/openmerlin/merlin-server/user/domain"
)

type SessionMapper interface {
	Insert(SessionDO) error
	Get(string) (SessionDO, error)
}

// TODO: mapper can be mysql
func NewSessionRepository(mapper SessionMapper) session {
	return session{mapper}
}

type session struct {
	mapper SessionMapper
}

func (impl session) Get(account domain.Account) (r sessiondomain.Session, err error) {
	do, err := impl.mapper.Get(account.Account())
	if err != nil {
		err = fmt.Errorf("failed to get session: %w", err)

		return
	}

	return r, do.toSession(&r)
}

func (impl session) Save(u *sessiondomain.Session) (err error) {
	if err = impl.mapper.Insert(impl.toSessionDO(u)); err != nil {
		err = fmt.Errorf("failed to save session: %w", err)
	}

	return
}

func (impl session) toSessionDO(u *sessiondomain.Session) SessionDO {
	return SessionDO{
		Account: u.Account.Account(),
		Info:    u.Info,
		Email:   u.Email.Email(),
		UserId:  u.UserId,
	}
}

type SessionDO struct {
	Account string
	Info    string
	Email   string
	UserId  string
}

func (do *SessionDO) toSession(r *sessiondomain.Session) (err error) {
	if r.Account, err = primitive.NewAccount(do.Account); err != nil {
		return
	}

	if r.Email, err = domain.NewEmail(do.Email); err != nil {
		return
	}

	r.Info = do.Info
	r.UserId = do.UserId

	return
}
