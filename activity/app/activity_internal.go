/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"github.com/openmerlin/merlin-server/activity/domain/repository"
)

// ActivityInternalAppService is an interface for the internal model application service.
type ActivityInternalAppService interface {
	Create(*CmdToAddActivity) error
	DeleteAll(*CmdToAddActivity) error
}

// NewActivityInternalAppService creates a new instance of the internal model application service.
func NewActivityInternalAppService(repoAdapter repository.ActivityInternalAdapter) ActivityInternalAppService {
	return &activityInternalAppService{
		repoAdapter: repoAdapter,
	}
}

type activityInternalAppService struct {
	repoAdapter repository.ActivityInternalAdapter
}

func (s *activityInternalAppService) DeleteAll(cmd *CmdToAddActivity) error {
	err := s.repoAdapter.DeleteAll(cmd)

	return err
}

// Create add activities.
func (s *activityInternalAppService) Create(cmd *CmdToAddActivity) error {
	err := s.repoAdapter.Save(cmd)

	return err
}
