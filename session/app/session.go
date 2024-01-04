package app

import (
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/session/domain"
	"github.com/openmerlin/merlin-server/session/domain/repository"
	userapp "github.com/openmerlin/merlin-server/user/app"
	"github.com/openmerlin/merlin-server/utils"
)

type SessionAppService interface {
	Login(*CmdToLogin) (dto SessionDTO, user UserDTO, err error)
	Logout(primitive.UUID) (string, error)

	CheckAndRefresh(*CmdToCheck) (primitive.Account, string, error)
}

func NewSessionAppService(
	oidc repository.OIDCAdapter,
	userApp userapp.UserService,
	maxLogin int,
	loginRepo repository.LoginRepositoryAdapter,
	csrfTokenRepo repository.CSRFTokenRepositoryAdapter,
) SessionAppService {
	return &sessionAppService{
		oidc:          oidc,
		userApp:       userApp,
		maxLogin:      maxLogin,
		loginRepo:     loginRepo,
		csrfTokenRepo: csrfTokenRepo,
	}
}

type sessionAppService struct {
	maxLogin int

	oidc          repository.OIDCAdapter
	userApp       userapp.UserService
	loginRepo     repository.LoginRepositoryAdapter
	csrfTokenRepo repository.CSRFTokenRepositoryAdapter
}

func (s *sessionAppService) Login(cmd *CmdToLogin) (dto SessionDTO, user UserDTO, err error) {
	login, err := s.oidc.GetByCode(cmd.Code, cmd.RedirectURI)
	if err != nil {
		return
	}

	user, err = s.userApp.GetByAccount(login.Name, false)
	if err != nil {
		if !allerror.IsNotFound(err) {
			return
		}

		if user, err = s.createUser(&login); err != nil {
			return
		}
	}

	if err = s.clearLogin(login.Name, cmd.IP); err != nil {
		return
	}

	v := domain.Login{
		Id:        primitive.CreateUUID(),
		IP:        cmd.IP,
		User:      login.Name,
		IdToken:   login.IDToken,
		CreatedAt: utils.Now(),
		UserAgent: cmd.UserAgent,
	}

	if err = s.loginRepo.Add(&v); err != nil {
		return
	}

	csrfToken := domain.NewCSRFToken(v.Id)

	if err = s.csrfTokenRepo.Add(&csrfToken); err == nil {
		dto.LoginId = csrfToken.LoginId
		dto.CSRFToken = csrfToken.Id
	}

	return
}

func (s *sessionAppService) Logout(loginId primitive.UUID) (string, error) {
	login, err := s.loginRepo.Find(loginId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		}

		return "", err
	}

	if err := s.loginRepo.Delete(loginId); err != nil {
		return "", err
	}

	return login.IdToken, nil
}

func (s *sessionAppService) createUser(login *repository.Login) (UserDTO, error) {
	return s.userApp.Create(&userapp.CmdToCreateUser{
		Bio:      login.Bio,
		Email:    login.Email,
		Account:  login.Name,
		AvatarId: login.AvatarId,
		Fullname: login.Fullname,
	})
}

func (s *sessionAppService) clearLogin(name primitive.Account, ip string) error {
	logins, err := s.loginRepo.FindByUser(name)
	if err != nil || len(logins) == 0 {
		return err
	}

	n := len(logins)
	for i := range logins {
		if item := &logins[i]; item.IsSameLogin(ip) {
			if err = s.loginRepo.Delete(item.Id); err != nil {
				return err
			}

			item.Invalidate()
			n--
		}
	}

	if n < s.maxLogin {
		return nil
	}

	for i := range logins {
		if logins[i].Invalid() {
			continue
		}

		if err = s.loginRepo.Delete(logins[i].Id); err != nil {
			return err
		}

		if n--; n < s.maxLogin {
			break
		}
	}

	return nil
}
