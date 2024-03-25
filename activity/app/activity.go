package app

import (
	"github.com/openmerlin/merlin-server/activity/domain/repository"
	coderepoapp "github.com/openmerlin/merlin-server/coderepo/app"
	commonapp "github.com/openmerlin/merlin-server/common/app"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

// NewActivityAppService is an interface for the model application service.
type ActivityAppService interface {
	List([]primitive.Account, *CmdToListActivities) (AcctivityDTO, error)
	Create(*CmdToAddActivity) error
	Delete(*CmdToAddActivity) error
	HasLike(primitive.Account, string) (bool, error)
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
func (s *activityAppService) List(names []primitive.Account, cmd *CmdToListActivities) (
	AcctivityDTO, error,
) {

	v, total, err := s.repoAdapter.List(names, cmd)

	return AcctivityDTO{
		Total:      total,
		Activities: v,
	}, err
}

func (s *activityAppService) Create(cmd *CmdToAddActivity) error {
	err := s.repoAdapter.Save(cmd)

	return err
}

// Delete delete activities.
func (s *activityAppService) Delete(cmd *CmdToAddActivity) error {
	err := s.repoAdapter.Delete(cmd)

	return err
}

// HasLike check if a user like a model or space
func (s *activityAppService) HasLike(acc primitive.Account, id string) (bool, error) {
	has, _ := s.repoAdapter.HasLike(acc, id)

	return has, nil
}
