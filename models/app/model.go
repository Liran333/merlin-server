package app

import (
	coderepoapp "github.com/openmerlin/merlin-server/coderepo/app"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/models/domain"
	"github.com/openmerlin/merlin-server/models/domain/repository"
	"github.com/openmerlin/merlin-server/utils"
)

var (
	errorModelNotFound = allerror.NewNotFound(allerror.ErrorCodeModelNotFound, "not found")
)

type ModelAppService interface {
	Create(primitive.Account, *CmdToCreateModel) (string, error)
	Delete(primitive.Account, primitive.Identity) error
	Update(primitive.Account, primitive.Identity, *CmdToUpdateModel) error
	GetByName(primitive.Account, *ModelIndex) (ModelDTO, error)
	List(primitive.Account, *CmdToListModels) (ModelsDTO, error)
}

func NewModelAppService(
	codeRepoApp coderepoapp.CodeRepoAppService,
	repoAdapter repository.ModelRepositoryAdapter,
) ModelAppService {
	return &modelAppService{
		codeRepoApp: codeRepoApp,
		repoAdapter: repoAdapter,
	}
}

type modelAppService struct {
	codeRepoApp coderepoapp.CodeRepoAppService
	repoAdapter repository.ModelRepositoryAdapter
}

func (s *modelAppService) Create(user primitive.Account, cmd *CmdToCreateModel) (string, error) {
	// TODO check if can create

	coderepo, err := s.codeRepoApp.Create(&cmd.CmdToCreateRepo)
	if err != nil {
		return "", err
	}

	now := utils.Now()

	err = s.repoAdapter.Add(&domain.Model{
		Desc:      cmd.Desc,
		Fullname:  cmd.Fullname,
		CodeRepo:  coderepo,
		CreatedAt: now,
		UpdatedAt: now,
	})

	return coderepo.Id.Identity(), err

	// TODO send model created event in order to add activity and operation log
}

func (s *modelAppService) Delete(user primitive.Account, modelId primitive.Identity) error {
	model, err := s.repoAdapter.FindById(modelId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		}

		return err
	}

	// TODO check if can delete

	if err = s.codeRepoApp.Delete(model.RepoIndex()); err != nil {
		return err
	}

	if err = s.repoAdapter.Delete(model.Id); err != nil {
		return err
	}

	// TODO send model deleted event

	return nil
}

func (s *modelAppService) Update(
	user primitive.Account, modelId primitive.Identity, cmd *CmdToUpdateModel,
) error {
	model, err := s.repoAdapter.FindById(modelId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errorModelNotFound
		}

		return err
	}

	// TODO check if can update

	b, err := s.codeRepoApp.Update(&model.CodeRepo, &cmd.CmdToUpdateRepo)
	if err != nil {
		return err
	}

	b1 := cmd.toModel(&model)
	if !b && !b1 {
		return nil
	}

	return s.repoAdapter.Save(&model)

	// send model updated event to add activity
}

func (s *modelAppService) GetByName(user primitive.Account, index *ModelIndex) (ModelDTO, error) {
	var dto ModelDTO

	model, err := s.repoAdapter.FindByName(index.Owner, index.Name)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errorModelNotFound
		}

		return dto, err
	}

	if model.IsPublic() {
		return toModelDTO(&model), nil
	}

	// can't access private model anonymously
	if user == nil {
		return dto, errorModelNotFound
	}

	if user != index.Owner {
		//TODO check if can get
	}

	return toModelDTO(&model), nil
}

func (s *modelAppService) List(user primitive.Account, cmd *CmdToListModels) (
	ModelsDTO, error,
) {
	option := cmd.toOption()

	if user == nil {
		option.Visibility = primitive.VisibilityPublic
	} else {
		if cmd.Owner == nil {
			// It can list the private models of user,
			// but it maybe no need to do it.
			option.Visibility = primitive.VisibilityPublic
		} else {
			if user != cmd.Owner {
				// TODO if user can't get, then
				// option.Visibility = primitive.VisibilityPublic
			}
		}
	}

	v, total, err := s.repoAdapter.List(&option)

	return ModelsDTO{
		Total:  total,
		Models: v,
	}, err
}

// DeleteById is an example for restful API.
func (s *modelAppService) DeleteById(user primitive.Account, modelId string) error {
	// get model by model id
	// check if user can delete it
	// delete it
	return nil
}
