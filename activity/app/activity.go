/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package app provides functionality for the application.
package app

import (
	"github.com/sirupsen/logrus"

	"github.com/openmerlin/merlin-server/activity/domain"
	"github.com/openmerlin/merlin-server/activity/domain/repository"
	coderepoapp "github.com/openmerlin/merlin-server/coderepo/app"
	commonapp "github.com/openmerlin/merlin-server/common/app"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

// NewActivityAppService is an interface for the model application service.
type ActivityAppService interface {
	List(primitive.Account, []primitive.Account, *CmdToListActivities) (ActivityDTO, error)
	Create(*CmdToAddActivity) error
	Delete(*CmdToAddActivity) error
	HasLike(primitive.Account, primitive.Identity) (bool, error)
}

type activityAppService struct {
	permission  commonapp.ResourcePermissionAppService
	codeRepoApp coderepoapp.CodeRepoAppService
	repoAdapter repository.ActivitiesRepositoryAdapter
}

// NewActivityAppService creates a new instance of the Activity application service.
func NewActivityAppService(
	permission commonapp.ResourcePermissionAppService,
	codeRepoApp coderepoapp.CodeRepoAppService,
	repoAdapter repository.ActivitiesRepositoryAdapter,
) ActivityAppService {
	return &activityAppService{
		permission:  permission,
		codeRepoApp: codeRepoApp,
		repoAdapter: repoAdapter,
	}
}

// List retrieves a list of activity.
func (s *activityAppService) List(user primitive.Account, names []primitive.Account, cmd *CmdToListActivities) (
	ActivityDTO, error,
) {
	activities, _, err := s.repoAdapter.List(names, cmd)
	var filteredActivities []domain.Activity
	for _, activity := range activities {
		codeRepo, _ := s.codeRepoApp.GetById(activity.Resource.Index)
		activity.Name = codeRepo.Name
		activity.Resource.Owner = codeRepo.Owner
		if err := s.permission.CanRead(user, &codeRepo); err != nil {
			if allerror.IsNoPermission(err) {
				continue
			} else {
				logrus.Errorf("failed to read permission: %s", err)
			}
		}
		filteredActivities = append(filteredActivities, activity)
	}

	return ActivityDTO{
		Total:      len(filteredActivities),
		Activities: filteredActivities,
	}, err
}

// Create function to check if a "like" already exists before saving.
func (s *activityAppService) Create(cmd *CmdToAddActivity) error {
	// Check if there's already a like for the given resource by the owner.
	alreadyLiked, err := s.repoAdapter.HasLike(cmd.Owner, cmd.Resource.Index)
	if err != nil {
		return err
	}
	if alreadyLiked {
		return nil
	}

	// Retrieve the code repository information.
	codeRepo, err := s.codeRepoApp.GetById(cmd.Resource.Index)
	if err != nil {
		return err
	}

	// Only proceed if the repository is public.
	isPublic := codeRepo.IsPublic()
	if !isPublic {
		return nil
	}

	// Save the new activity.
	return s.repoAdapter.Save(cmd)
}

// Delete function to check if a "like" already exists before delete.
func (s *activityAppService) Delete(cmd *CmdToAddActivity) error {
	has, err := s.repoAdapter.HasLike(cmd.Owner, cmd.Resource.Index)
	if err != nil {
		return err
	}
	if !has {
		return nil
	}
	return s.repoAdapter.Delete(cmd)
}

// HasLike check if a user like a model or space
func (s *activityAppService) HasLike(acc primitive.Account, id primitive.Identity) (bool, error) {
	has, _ := s.repoAdapter.HasLike(acc, id)

	return has, nil
}
