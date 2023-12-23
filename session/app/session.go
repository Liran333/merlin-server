package session

import (
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	session "github.com/openmerlin/merlin-server/session/domain"
	"github.com/openmerlin/merlin-server/user/domain"
)

type SessionCreateCmd struct {
	Account domain.Account
	Info    string
	Email   domain.Email
	UserId  string
}

func (cmd *SessionCreateCmd) Validate() error {
	if b := cmd.Account != nil && cmd.Info != "" && cmd.Email != nil && cmd.UserId != ""; !b {
		return allerror.NewInvalidParam("invalid cmd of creating login")
	}

	return nil
}

func (cmd *SessionCreateCmd) toSession() session.Session {
	return session.Session{
		Account: cmd.Account,
		Info:    cmd.Info,
		Email:   cmd.Email,
		UserId:  cmd.UserId,
	}
}

type SessionDTO struct {
	Info   string `json:"info"`
	Email  string `json:"email"`
	UserId string `json:"user_id"`
}

type SessionService interface {
	Create(*SessionCreateCmd) error
	Get(domain.Account) (SessionDTO, error)
}

func NewSessionService(repo session.SessionRepo) SessionService {
	return loginService{
		repo: repo,
	}
}

type loginService struct {
	repo session.SessionRepo
}

func (s loginService) Create(cmd *SessionCreateCmd) error {
	v := cmd.toSession()

	// new login
	return s.repo.Save(&v)
}

func (s loginService) Get(account domain.Account) (dto SessionDTO, err error) {
	v, err := s.repo.Get(account)
	if err != nil {
		return
	}

	s.toSessionDTO(&v, &dto)

	return
}

func (s loginService) toSessionDTO(u *session.Session, dto *SessionDTO) {
	dto.Info = u.Info
	dto.Email = u.Email.Email()
	dto.UserId = u.UserId
}
