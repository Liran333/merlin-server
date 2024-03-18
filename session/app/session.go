/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"fmt"

	"github.com/sirupsen/logrus"

	commondomain "github.com/openmerlin/merlin-server/common/domain"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/session/domain"
	"github.com/openmerlin/merlin-server/session/domain/repository"
	userapp "github.com/openmerlin/merlin-server/user/app"
	"github.com/openmerlin/merlin-server/utils"
)

// SessionAppService is an interface for session application service.
type SessionAppService interface {
	Login(*CmdToLogin) (dto SessionDTO, user UserDTO, err error)
	Logout(primitive.RandomId) (string, error)
	Clear(primitive.RandomId) error

	CheckAndRefresh(*CmdToCheck) (primitive.Account, string, error)
	CheckSession(*CmdToCheck) (primitive.Account, error)
}

// NewSessionAppService creates a new instance of sessionAppService.
func NewSessionAppService(
	oidc repository.OIDCAdapter,
	userApp userapp.UserService,
	maxLogin int,
	sessionRepo repository.SessionRepositoryAdapter,
	csrfTokenRepo repository.CSRFTokenRepositoryAdapter,
	sessionFastRepo repository.SessionFastRepositoryAdapter,
) SessionAppService {
	return &sessionAppService{
		oidc:            oidc,
		userApp:         userApp,
		maxLogin:        maxLogin,
		sessionRepo:     sessionRepo,
		csrfTokenRepo:   csrfTokenRepo,
		sessionFastRepo: sessionFastRepo,
	}
}

type sessionAppService struct {
	maxLogin int

	oidc            repository.OIDCAdapter
	userApp         userapp.UserService
	sessionRepo     repository.SessionRepositoryAdapter
	csrfTokenRepo   repository.CSRFTokenRepositoryAdapter
	sessionFastRepo repository.SessionFastRepositoryAdapter
}

// Login logs in a user and returns the session DTO, user DTO, and error.
func (s *sessionAppService) Login(cmd *CmdToLogin) (dto SessionDTO, user UserDTO, err error) {
	login, err := s.oidc.GetByCode(cmd.Code, cmd.RedirectURI)
	if err != nil {
		return
	}

	user, err = s.userApp.GetByAccount(login.Name, login.Name)
	if err != nil {
		if !allerror.IsNotFound(err) {
			return
		}

		if user, err = s.createUser(&login); err != nil {
			return
		}
	}

	if !user.IsAgreePrivacy {
		err = s.userApp.AgreePrivacy(primitive.CreateAccount(user.Name))
		if err != nil {
			return
		}
	}

	if err = s.clearLogin(login.Name, cmd.IP); err != nil {
		return
	}

	if dto.SessionId, err = primitive.NewRandomId(); err != nil {
		return
	}

	if dto.CSRFToken, err = primitive.NewRandomId(); err != nil {
		return
	}

	session := domain.Session{
		Id:        dto.SessionId,
		IP:        cmd.IP,
		User:      login.Name,
		UserId:    login.UserId,
		IdToken:   login.IDToken,
		UserAgent: cmd.UserAgent,
		CreatedAt: utils.Now(),
	}

	if err = s.sessionRepo.Add(&session); err != nil {
		return
	}

	if err = s.sessionFastRepo.Add(&session); err != nil {
		return
	}

	csrfToken := session.NewCSRFToken()
	err = s.csrfTokenRepo.Add(dto.CSRFToken, &csrfToken)

	return
}

// Logout logs out a user and returns the ID token and error.
func (s *sessionAppService) Logout(sessionId primitive.RandomId) (string, error) {
	session, err := s.sessionFastRepo.Find(sessionId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		}

		return "", err
	}

	if err := s.sessionRepo.Delete(session.Id); err != nil {
		return "", err
	}

	if err = s.sessionFastRepo.Delete(sessionId); err != nil {
		return "", err
	}

	return session.IdToken, nil
}

func (s *sessionAppService) Clear(sessionId primitive.RandomId) error {
	_, err := s.sessionRepo.Find(sessionId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		}

		return err
	}

	err = s.sessionRepo.Delete(sessionId)

	r := commondomain.OperationLogRecord{
		Time:    utils.Time(),
		User:    primitive.CreateAccount("service"),
		IP:      "local",
		Method:  "Auto",
		Action:  fmt.Sprintf("clear expired session automatically by id: %s", sessionId.RandomId()),
		Success: err == nil,
	}

	logrus.Info(r.String())

	return err
}

func (s *sessionAppService) createUser(login *repository.Login) (UserDTO, error) {
	return s.userApp.Create(&userapp.CmdToCreateUser{
		Desc:     login.Desc,
		Email:    login.Email,
		Account:  login.Name,
		AvatarId: login.AvatarId,
		Fullname: login.Fullname,
		Phone:    login.Phone,
	})
}

func (s *sessionAppService) clearLogin(name primitive.Account, ip string) error {
	logins, err := s.sessionRepo.FindByUser(name)
	if err != nil || len(logins) == 0 {
		return err
	}

	deleteSession := func(sessionId primitive.RandomId) error {
		if err := s.sessionFastRepo.Delete(sessionId); err != nil {
			return err
		}

		return s.sessionRepo.Delete(sessionId)
	}

	n := len(logins)
	for i := range logins {
		if item := &logins[i]; item.IsSameLogin(ip) {
			if err = deleteSession(item.Id); err != nil {
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
		if logins[i].IsInvalid() {
			continue
		}

		if err = deleteSession(logins[i].Id); err != nil {
			return err
		}

		if n--; n < s.maxLogin {
			break
		}
	}

	return nil
}
