/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package app provides functionality for the application.
package app

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"

	"github.com/openmerlin/merlin-server/activity/domain"
	"github.com/openmerlin/merlin-server/activity/domain/message"
	"github.com/openmerlin/merlin-server/activity/domain/repository"
	coderepoapp "github.com/openmerlin/merlin-server/coderepo/app"
	coderepo "github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/coderepo/domain/resourceadapter"
	commonapp "github.com/openmerlin/merlin-server/common/app"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
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
	permission      commonapp.ResourcePermissionAppService
	codeRepoApp     coderepoapp.CodeRepoAppService
	repoAdapter     repository.ActivitiesRepositoryAdapter
	modelApp        modelApp.ModelAppService
	spaceApp        spaceApp.SpaceAppService
	msgAdapter      message.ActivityMessage
	resourceAdapter resourceadapter.ResourceAdapter
}

// NewActivityAppService creates a new instance of the Activity application service.
func NewActivityAppService(
	permission commonapp.ResourcePermissionAppService,
	codeRepoApp coderepoapp.CodeRepoAppService,
	repoAdapter repository.ActivitiesRepositoryAdapter,
	modelApp modelApp.ModelAppService,
	spaceApp spaceApp.SpaceAppService,
	msgAdapter message.ActivityMessage,
	resourceAdapter resourceadapter.ResourceAdapter,
) ActivityAppService {
	return &activityAppService{
		permission:      permission,
		codeRepoApp:     codeRepoApp,
		repoAdapter:     repoAdapter,
		modelApp:        modelApp,
		spaceApp:        spaceApp,
		msgAdapter:      msgAdapter,
		resourceAdapter: resourceAdapter,
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
		activitySummary, errProcess := s.processActivity(user, activity)
		if errProcess != nil {
			continue
		}
		filteredActivities = append(filteredActivities, activitySummary)
	}

	return ActivitysDTO{
		Total:      len(filteredActivities),
		Activities: filteredActivities,
	}, nil
}

func (s *activityAppService) processActivity(user primitive.Account, activity domain.Activity) (ActivitySummaryDTO, error) {
	codeRepo, err := s.codeRepoApp.GetById(activity.Resource.Index)
	if err != nil {
		return ActivitySummaryDTO{}, xerrors.Errorf("failed to get code repository by ID: %w", err)
	}

	activity.Name = codeRepo.Name
	activity.Resource.Owner = codeRepo.Owner

	switch activity.Resource.Type {
	case primitive.ObjTypeModel:
		return s.processModelActivity(user, codeRepo, activity)
	case primitive.ObjTypeSpace:
		return s.processSpaceActivity(user, codeRepo, activity)
	default:
		return ActivitySummaryDTO{}, xerrors.Errorf("unknown resource type")
	}
}

func (s *activityAppService) processModelActivity(user primitive.Account, codeRepo coderepo.CodeRepo, activity domain.Activity) (ActivitySummaryDTO, error) {
	model, err := s.modelApp.GetByName(user, &coderepo.CodeRepoIndex{Name: codeRepo.Name, Owner: codeRepo.Owner})
	if err != nil {
		return ActivitySummaryDTO{}, err
	}
	activity.Resource.Disable = model.Disable
	stat := domain.Stat{
		LikeCount:     model.LikeCount,
		DownloadCount: model.DownloadCount,
	}
	additionInfo := fromModelDTO(model, &activity, &stat)
	return additionInfo, nil
}

func (s *activityAppService) processSpaceActivity(user primitive.Account, codeRepo coderepo.CodeRepo, activity domain.Activity) (ActivitySummaryDTO, error) {
	space, err := s.spaceApp.GetByName(user, &coderepo.CodeRepoIndex{Name: codeRepo.Name, Owner: codeRepo.Owner})
	if err != nil {
		return ActivitySummaryDTO{}, err
	}
	activity.Resource.Disable = space.Disable
	stat := domain.Stat{
		LikeCount:     space.LikeCount,
		DownloadCount: space.DownloadCount,
	}
	additionInfo := fromSpaceDTO(space, &activity, &stat)
	return additionInfo, nil
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
		return xerrors.Errorf("failed to get a coderepo by id, error: %w", err)
	}

	// Only proceed if the repository is public.
	isPublic := codeRepo.IsPublic()
	if !isPublic {
		return nil
	}

	// Retrieve the code repository information.
	Repo, err := s.resourceAdapter.GetByIndex(cmd.Resource.Index)
	if err != nil {
		return xerrors.Errorf("failed to get a resouce by id, error: %w", err)
	}

	e := domain.NewLikeCreatedEvent(&codeRepo, string(Repo.ResourceType()))
	if err := s.msgAdapter.SendLikeCreatedEvent(&e); err != nil {
		logrus.Errorf("failed to send like created event, error:%s", err)

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

	// Retrieve the code repository information.
	codeRepo, err := s.codeRepoApp.GetById(cmd.Resource.Index)
	if err != nil {
		return xerrors.Errorf("failed to get a coderepo by id, error:%w", err)
	}

	// Retrieve the code repository information.
	Repo, err := s.resourceAdapter.GetByIndex(cmd.Resource.Index)
	if err != nil {
		return xerrors.Errorf("failed to get a resource by id, error: %w", err)
	}

	e := domain.NewLikeCreatedEvent(&codeRepo, string(Repo.ResourceType()))
	if err := s.msgAdapter.SendLikeDeletedEvent(&e); err != nil {
		logrus.Errorf("failed to send like deleted event, error:%s", err)
	}

	return s.repoAdapter.Delete(cmd)
}

// HasLike check if a user like a model or space
func (s *activityAppService) HasLike(acc primitive.Account, id primitive.Identity) (bool, error) {
	has, _ := s.repoAdapter.HasLike(acc, id)

	return has, nil
}
