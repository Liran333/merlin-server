/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"fmt"

	coderepoapp "github.com/openmerlin/merlin-server/coderepo/app"
	commonapp "github.com/openmerlin/merlin-server/common/app"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/models/domain"
	"github.com/openmerlin/merlin-server/models/domain/repository"
	"github.com/openmerlin/merlin-server/utils"
)

var (
	errorModelNotFound      = allerror.NewNotFound(allerror.ErrorCodeModelNotFound, "not found")
	errorModelCountExceeded = allerror.NewCountExceeded("model count exceed")
)

// ModelAppService is an interface for the model application service.
type ModelAppService interface {
	Create(primitive.Account, *CmdToCreateModel) (string, error)
	Delete(primitive.Account, primitive.Identity) (string, error)
	Update(primitive.Account, primitive.Identity, *CmdToUpdateModel) (string, error)
	GetByName(primitive.Account, *domain.ModelIndex) (ModelDTO, error)
	List(primitive.Account, *CmdToListModels) (ModelsDTO, error)
}

// NewModelAppService creates a new instance of the model application service.
func NewModelAppService(
	permission commonapp.ResourcePermissionAppService,
	codeRepoApp coderepoapp.CodeRepoAppService,
	repoAdapter repository.ModelRepositoryAdapter,
) ModelAppService {
	return &modelAppService{
		permission:  permission,
		codeRepoApp: codeRepoApp,
		repoAdapter: repoAdapter,
	}
}

type modelAppService struct {
	permission  commonapp.ResourcePermissionAppService
	codeRepoApp coderepoapp.CodeRepoAppService
	repoAdapter repository.ModelRepositoryAdapter
}

// Create creates a new model.
func (s *modelAppService) Create(user primitive.Account, cmd *CmdToCreateModel) (string, error) {
	if err := s.permission.CanCreate(user, cmd.Owner, primitive.ObjTypeModel); err != nil {
		return "", err
	}

	if err := s.modelCountCheck(cmd.Owner); err != nil {
		return "", err
	}

	coderepo, err := s.codeRepoApp.Create(&cmd.CmdToCreateRepo)
	if err != nil {
		return "", err
	}

	now := utils.Now()

	err = s.repoAdapter.Add(&domain.Model{
		Desc:      cmd.Desc,
		Fullname:  cmd.Fullname,
		CodeRepo:  coderepo,
		CreatedBy: user,
		CreatedAt: now,
		UpdatedAt: now,
	})

	return coderepo.Id.Identity(), err
}

// Delete deletes a model.
func (s *modelAppService) Delete(user primitive.Account, modelId primitive.Identity) (action string, err error) {
	model, err := s.repoAdapter.FindById(modelId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		}

		return
	}

	action = fmt.Sprintf(
		"delete model of %s:%s/%s",
		modelId.Identity(), model.Owner.Account(), model.Name.MSDName(),
	)

	if err = s.permission.CanDelete(user, &model); err != nil {
		if allerror.IsNoPermission(err) {
			err = errorModelNotFound
		}

		return
	}

	if err = s.codeRepoApp.Delete(model.RepoIndex()); err != nil {
		return
	}

	if err = s.repoAdapter.Delete(model.Id); err != nil {
		return
	}

	return
}

// Update updates a model.
func (s *modelAppService) Update(
	user primitive.Account, modelId primitive.Identity, cmd *CmdToUpdateModel,
) (action string, err error) {
	model, err := s.repoAdapter.FindById(modelId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errorModelNotFound
		}

		return
	}

	action = fmt.Sprintf(
		"update model of %s:%s/%s",
		modelId.Identity(), model.Owner.Account(), model.Name.MSDName(),
	)

	if err = s.permission.CanUpdate(user, &model); err != nil {
		if allerror.IsNoPermission(err) {
			err = errorModelNotFound
		}

		return
	}

	b, err := s.codeRepoApp.Update(&model.CodeRepo, &cmd.CmdToUpdateRepo)
	if err != nil {
		return
	}

	b1 := cmd.toModel(&model)
	if !b && !b1 {
		return
	}

	err = s.repoAdapter.Save(&model)

	return
}

// GetByName retrieves a model by its name.
func (s *modelAppService) GetByName(user primitive.Account, index *domain.ModelIndex) (ModelDTO, error) {
	var dto ModelDTO

	model, err := s.repoAdapter.FindByName(index)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errorModelNotFound
		}

		return dto, err
	}

	if err := s.permission.CanRead(user, &model); err != nil {
		if allerror.IsNoPermission(err) {
			err = errorModelNotFound
		}

		return dto, err
	}

	return toModelDTO(&model), nil
}

// List retrieves a list of models.
func (s *modelAppService) List(user primitive.Account, cmd *CmdToListModels) (
	ModelsDTO, error,
) {
	if user == nil {
		cmd.Visibility = primitive.VisibilityPublic
	} else {
		if cmd.Owner == nil {
			// It can list the private models of user,
			// but it maybe no need to do it.
			cmd.Visibility = primitive.VisibilityPublic
		} else {
			if user != cmd.Owner {
				err := s.permission.CanListOrgResource(
					user, cmd.Owner, primitive.ObjTypeModel,
				)
				if err != nil {
					cmd.Visibility = primitive.VisibilityPublic
				}
			}
		}
	}

	v, total, err := s.repoAdapter.List(cmd)

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

func (s *modelAppService) modelCountCheck(owner primitive.Account) error {
	cmdToList := CmdToListModels{
		Owner: owner,
	}

	total, err := s.repoAdapter.Count(&cmdToList)
	if err != nil {
		return err
	}

	if total >= config.MaxCountPerOwner {
		return errorModelCountExceeded
	}

	return nil
}
