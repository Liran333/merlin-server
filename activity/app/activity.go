/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package app provides functionality for the application.
package app

import (
	"github.com/openmerlin/merlin-server/common/domain/allerror"

	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"

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
	List(primitive.Account, []primitive.Account, *CmdToListActivities) (ActivitysDTO, error)
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
func (s *activityAppService) List(user primitive.Account, names []primitive.Account, cmd *CmdToListActivities) (ActivitysDTO, error) {
	activities, _, err := s.repoAdapter.List(names, cmd)
	if err != nil {
		e := xerrors.Errorf("failed to list activities: %w", err)
		err = allerror.New(allerror.ErrorFailToRetrieveActivityData, "", e)
		return ActivitysDTO{}, err
	}

	var filteredActivities []ActivitySummaryDTO
	for _, activity := range activities {
		activitySummary, err := s.getActivity(user, activity)
		if err != nil {
			logrus.Infof("User %v %v repo id:%v", user, err, activity.Resource.Index)
			continue
		}
		filteredActivities = append(filteredActivities, activitySummary)
	}

	return ActivitysDTO{
		Total:      len(filteredActivities),
		Activities: filteredActivities,
	}, nil
}

// getActivity retrieves statistics based on the activity type.
func (s *activityAppService) getActivity(user primitive.Account, activity domain.Activity) (ActivitySummaryDTO, error) {
	codeRepo, err := s.codeRepoApp.GetById(activity.Resource.Index)
	if err != nil {
		return ActivitySummaryDTO{}, xerrors.Errorf("failed to get code repository by ID: %w", err)
	}
	activity.Name = codeRepo.Name
	activity.Resource.Owner = codeRepo.Owner

	stat := domain.Stat{}
	err = s.getActivityInfo(user, codeRepo, &activity, &stat)
	if err != nil {
		return ActivitySummaryDTO{}, xerrors.Errorf("failed to get activity info by ID: %w", err)
	}

	return toActivitySummaryDTO(&activity, &stat), nil
}

func (s *activityAppService) getActivityInfo(user primitive.Account, codeRepo coderepo.CodeRepo, activity *domain.Activity, statsMap *domain.Stat) error {
	if activity.Resource.Type == primitive.ObjTypeModel {
		model, err := s.modelApp.GetByName(user, &coderepo.CodeRepoIndex{Name: codeRepo.Name, Owner: codeRepo.Owner})
		if err != nil {
			return xerrors.Errorf("failed to get model, id:%v, err: %w", activity.Resource.Index, err)
		}
		activity.Resource.Disable = model.Disable
		statsMap.LikeCount = model.LikeCount
		statsMap.DownloadCount = model.DownloadCount

	} else {
		space, err := s.spaceApp.GetByName(user, &coderepo.CodeRepoIndex{Name: codeRepo.Name, Owner: codeRepo.Owner})
		if err != nil {
			logrus.Errorf("failed to get space, id:%v, err: %v", activity.Resource.Index, err)
			return err
		}
		activity.Resource.Disable = space.Disable
		statsMap.LikeCount = space.LikeCount
		statsMap.DownloadCount = space.DownloadCount
	}
	return nil
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
