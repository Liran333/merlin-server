package app

import (
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/session/domain"
)

func (s *sessionAppService) CheckAndRefresh(cmd *CmdToCheck) (primitive.Account, string, error) {
	user, err := s.check(cmd)
	if err != nil {
		return nil, "", err
	}

	csrfToken := domain.NewCSRFToken(cmd.LoginId)

	if err := s.csrfTokenRepo.Add(&csrfToken); err != nil {
		return nil, "", err
	}

	return user, csrfToken.Id.String(), nil
}

func (s *sessionAppService) check(cmd *CmdToCheck) (primitive.Account, error) {
	// check csrf token
	csrfToken, err := s.csrfTokenRepo.Find(cmd.CSRFToken)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.New(allerror.ErrorCodeCSRFTokenNotFound, "no csrf token")
		}

		return nil, err
	}

	if err := csrfToken.Validate(cmd.LoginId); err != nil {
		return nil, err
	}

	if csrfToken.Reset() {
		if err := s.csrfTokenRepo.Save(&csrfToken); err != nil {
			return nil, err
		}
	}

	// check login
	login, err := s.loginRepo.Find(cmd.LoginId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.New(allerror.ErrorCodeLoginIdNotFound, "no login")
		}

		return nil, err
	}

	if err := login.Validate(cmd.IP, cmd.UserAgent); err != nil {
		return nil, err
	}

	return login.User, nil
}
