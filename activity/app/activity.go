/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package app provides functionality for the application.
package app

import (
	"fmt"
	"github.com/openmerlin/merlin-server/common/domain/allerror"

	"github.com/sirupsen/logrus"

	"github.com/openmerlin/merlin-server/activity/domain"
	"github.com/openmerlin/merlin-server/activity/domain/repository"
	coderepoapp "github.com/openmerlin/merlin-server/coderepo/app"
	coderepo "github.com/openmerlin/merlin-server/coderepo/domain"
	commonapp "github.com/openmerlin/merlin-server/common/app"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	modelApp "github.com/openmerlin/merlin-server/models/app"
	spaceApp "github.com/openmerlin/merlin-server/space/app"
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
	modelApp    modelApp.ModelAppService
	spaceApp    spaceApp.SpaceAppService
}

// NewActivityAppService creates a new instance of the Activity application service.
func NewActivityAppService(
	permission commonapp.ResourcePermissionAppService,
	codeRepoApp coderepoapp.CodeRepoAppService,
	repoAdapter repository.ActivitiesRepositoryAdapter,
	modelApp modelApp.ModelAppService,
	spaceApp spaceApp.SpaceAppService,
) ActivityAppService {
	return &activityAppService{
		permission:  permission,
		codeRepoApp: codeRepoApp,
		repoAdapter: repoAdapter,
		modelApp:    modelApp,
		spaceApp:    spaceApp,
	}
}

// List retrieves a list of activities with statistics for models and spaces.
func (s *activityAppService) List(user primitive.Account, names []primitive.Account, cmd *CmdToListActivities) (ActivityDTO, error) {
	activities, _, err := s.repoAdapter.List(names, cmd)
	if err != nil {
		e := fmt.Errorf("failed to list activities: %s", err)
		err = allerror.New(allerror.ErrorFailToRetrieveActivityData, e.Error(), e)
		return ActivityDTO{}, err
	}

	var filteredActivities []domain.ActivitySummary
	for i := range activities {
		activity := activities[i]
		stat, err := s.getActivityStats(user, &activity)
		if err != nil {
			// User does not have access to this repo
			e := fmt.Errorf("failed to get repo info by index: %s", err)
			logrus.Infof("User %v does not have access to this repo id: %s. %t", user, activity.Resource.Index, e)
			continue
		}
		filteredActivities = append(filteredActivities, domain.ActivitySummary{
			Activity: activity,
			Stat:     stat,
		})
	}

	return ActivityDTO{
		Total:      len(filteredActivities),
		Activities: filteredActivities,
	}, nil
}

// getActivityStats retrieves statistics based on the activity type.
func (s *activityAppService) getActivityStats(user primitive.Account, activity *domain.Activity) (domain.Stat, error) {
	codeRepo, err := s.codeRepoApp.GetById(activity.Resource.Index)
	if err != nil {
		return domain.Stat{}, fmt.Errorf("failed to get code repository by ID: %v", err)
	}
	activity.Name = codeRepo.Name
	activity.Resource.Owner = codeRepo.Owner

	switch activity.Resource.Type {
	case primitive.ObjTypeModel:
		return s.getModelStats(user, codeRepo)
	case primitive.ObjTypeSpace:
		return s.getSpaceStats(user, codeRepo)
	default:
		return domain.Stat{}, fmt.Errorf("unknown activity type: %s", activity.Type)
	}
}

// getModelStats retrieves statistics for a model activity.
func (s *activityAppService) getModelStats(user primitive.Account, codeRepo coderepo.CodeRepo) (domain.Stat, error) {
	model, err := s.modelApp.GetByName(user, &coderepo.CodeRepoIndex{Name: codeRepo.Name, Owner: codeRepo.Owner})
	if err != nil {
		return domain.Stat{}, fmt.Errorf("failed to find model by name: %v", err)
	}
	statsMap := domain.Stat{
		LikeCount:     model.LikeCount,
		DownloadCount: model.DownloadCount,
	}

	return statsMap, nil
}

// getSpaceStats retrieves statistics for a space activity.
func (s *activityAppService) getSpaceStats(user primitive.Account, codeRepo coderepo.CodeRepo) (domain.Stat, error) {
	space, err := s.spaceApp.GetByName(user, &coderepo.CodeRepoIndex{Name: codeRepo.Name, Owner: codeRepo.Owner})
	if err != nil {
		return domain.Stat{}, fmt.Errorf("failed to find space by name: %v", err)
	}
	statsMap := domain.Stat{
		LikeCount:     space.LikeCount,
		DownloadCount: space.DownloadCount,
	}

	return statsMap, nil
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
