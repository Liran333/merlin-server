/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package app provides the application layer functionality for managing sessions and authentication.
package app

import (
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/session/domain"
)

// CheckAndRefresh checks and refreshes the session for a given command.
func (s *sessionAppService) CheckAndRefresh(cmd *CmdToCheck) (
	user primitive.Account, token string, err error,
) {
	session, refreshToken, err := s.check(cmd)
	if err != nil {
		return
	}

	user = session.User
	token = cmd.CSRFToken.RandomId()

	if !refreshToken {
		return
	}

	tokenId, err := primitive.NewRandomId()
	if err != nil {
		return
	}

	csrfToken := session.NewCSRFToken()
	if err = s.csrfTokenRepo.Add(tokenId, &csrfToken); err != nil {
		return
	}

	token = tokenId.RandomId()

	return
}

func (s *sessionAppService) check(cmd *CmdToCheck) (
	session domain.Session, refreshToken bool, err error,
) {
	session, err = s.sessionFastRepo.Find(cmd.SessionId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.New(allerror.ErrorCodeSessionNotFound, "no session")
		}

		return
	}

	csrfToken, err := s.csrfTokenRepo.Find(cmd.CSRFToken)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.New(allerror.ErrorCodeCSRFTokenNotFound, "no csrf token")
		}

		return
	}

	if err = csrfToken.Validate(cmd.SessionId); err != nil {
		return
	}

	if err = session.Validate(cmd.IP, cmd.UserAgent); err != nil {
		return
	}

	if err = s.sessionFastRepo.Save(&session); err != nil {
		return
	}

	refreshToken = csrfToken.IsExpired()

	return
}

// CheckSession checks the session for a given command.
func (s *sessionAppService) CheckSession(cmd *CmdToCheck) (user primitive.Account, err error) {
	session, err := s.sessionFastRepo.Find(cmd.SessionId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.New(allerror.ErrorCodeSessionNotFound, "no session")
		}

		return
	}

	err = s.sessionFastRepo.Save(&session)

	user = session.User
	return
}
