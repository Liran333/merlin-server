/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"fmt"

	"github.com/sirupsen/logrus"

	coderepoapp "github.com/openmerlin/merlin-server/coderepo/app"
	commonapp "github.com/openmerlin/merlin-server/common/app"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/models/domain"
	"github.com/openmerlin/merlin-server/models/domain/message"
	"github.com/openmerlin/merlin-server/models/domain/repository"
	"github.com/openmerlin/merlin-server/utils"
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
	msgAdapter message.ModelMessage,
	codeRepoApp coderepoapp.CodeRepoAppService,
	repoAdapter repository.ModelRepositoryAdapter,
) ModelAppService {
	return &modelAppService{
		permission:  permission,
		msgAdapter:  msgAdapter,
		codeRepoApp: codeRepoApp,
		repoAdapter: repoAdapter,
	}
}

type modelAppService struct {
	permission  commonapp.ResourcePermissionAppService
	msgAdapter  message.ModelMessage
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

	coderepo, err := s.codeRepoApp.Create(user, &cmd.CmdToCreateRepo)
	if err != nil {
		return "", err
	}

	now := utils.Now()
	model := domain.Model{
		Desc:      cmd.Desc,
		Fullname:  cmd.Fullname,
		CodeRepo:  coderepo,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err = s.repoAdapter.Add(&model); err != nil {
		return "", err
	}

	e := domain.NewModelCreatedEvent(&model)
	if err1 := s.msgAdapter.SendModelCreatedEvent(&e); err1 != nil {
		logrus.Errorf("failed to send model created event, model id:%s", model.Id.Identity())

	}

	return model.Id.Identity(), nil
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

	notFound, err := commonapp.CanDeleteOrNotFound(user, &model, s.permission)
	if err != nil {
		return
	}
	if notFound {
		err = allerror.NewNotFound(allerror.ErrorCodeModelNotFound, "not found", fmt.Errorf("%s not found", modelId.Identity()))

		return
	}

	if err = s.codeRepoApp.Delete(model.RepoIndex()); err != nil {
		return
	}

	if err = s.repoAdapter.Delete(model.Id); err != nil {
		return
	}

	e := domain.NewModelDeletedEvent(user, modelId)
	if err1 := s.msgAdapter.SendModelDeletedEvent(&e); err1 != nil {
		logrus.Errorf("failed to send model deleted event, model id:%s", modelId.Identity())
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
			err = allerror.NewNotFound(allerror.ErrorCodeModelNotFound, "not found", err)
		}

		return
	}

	action = fmt.Sprintf(
		"update model of %s:%s/%s",
		modelId.Identity(), model.Owner.Account(), model.Name.MSDName(),
	)

	notFound, err := commonapp.CanUpdateOrNotFound(user, &model, s.permission)
	if err != nil {
		return
	}
	if notFound {
		err = allerror.NewNotFound(allerror.ErrorCodeModelNotFound, "not found", fmt.Errorf("%s not found", modelId.Identity()))

		return
	}

	isPrivateToPublic := model.IsPrivate() && cmd.Visibility.IsPublic()

	b, err := s.codeRepoApp.Update(&model.CodeRepo, &cmd.CmdToUpdateRepo)
	if err != nil {
		return
	}

	b1 := cmd.toModel(&model)
	if !b && !b1 {
		return
	}

	if err = s.repoAdapter.Save(&model); err != nil {
		return
	}

	e := domain.NewModelUpdatedEvent(&model, user, isPrivateToPublic)
	if err1 := s.msgAdapter.SendModelUpdatedEvent(&e); err1 != nil {
		logrus.Errorf("failed to send model updated event, model id:%s", modelId.Identity())
	}

	return
}

// GetByName retrieves a model by its name.
func (s *modelAppService) GetByName(user primitive.Account, index *domain.ModelIndex) (ModelDTO, error) {
	var dto ModelDTO

	model, err := s.repoAdapter.FindByName(index)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeModelNotFound, "not found", err)
		}

		return dto, err
	}

	if err := s.permission.CanRead(user, &model); err != nil {
		if allerror.IsNoPermission(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeModelNotFound, "not found", err)
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
		return allerror.NewCountExceeded("model count exceed", fmt.Errorf("model count(now:%d max:%d) exceed", total, config.MaxCountPerOwner))
	}

	return nil
}
