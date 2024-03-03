/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"github.com/sirupsen/logrus"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/session/domain/repository"
)

// NewSessionClearAppService creates a new instance of sessionClearAppService
func NewSessionClearAppService(
	s repository.SessionRepositoryAdapter,
	sa repository.SessionFastRepositoryAdapter,
) sessionClearAppService {
	return sessionClearAppService{
		sessionRepo:     s,
		sessionFastRepo: sa,
	}
}

type sessionClearAppService struct {
	sessionRepo     repository.SessionRepositoryAdapter
	sessionFastRepo repository.SessionFastRepositoryAdapter
}

// ClearAllSession clears all session
func (s sessionClearAppService) ClearAllSession(user primitive.Account) error {
	sessions, err := s.sessionRepo.FindByUser(user)
	if err != nil || len(sessions) == 0 {
		return err
	}

	for i := range sessions {
		if err := s.sessionFastRepo.Delete(sessions[i].Id); err != nil {
			logrus.Errorf(
				"clear session stored at fast repo failed, session id:%s, err:%s",
				sessions[i].Id.RandomId(), err.Error(),
			)
		}
	}

	if err = s.sessionRepo.DeleteByUser(user); err != nil {
		logrus.Errorf("clear user's all sessions failed, err: %s", err.Error())
	}

	return nil
}
