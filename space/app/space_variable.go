/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package app

import (
	"fmt"
	spacerepo "github.com/openmerlin/merlin-server/space/domain/repository"
	"github.com/openmerlin/merlin-server/space/domain/securestorage"

	"github.com/sirupsen/logrus"

	commonapp "github.com/openmerlin/merlin-server/common/app"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/space/domain"
	"github.com/openmerlin/merlin-server/utils"
)

// SpaceVariableService is an interface for the space variable service.
type SpaceVariableService interface {
	CreateVariable(primitive.Account, primitive.Identity, *CmdToCreateSpaceVariable) (string, string, error)
	DeleteVariable(primitive.Account, primitive.Identity, primitive.Identity) (string, error)
	UpdateVariable(primitive.Account, primitive.Identity, primitive.Identity, *CmdToUpdateSpaceVariable) (string, error)
	ListVariableSecret(string) (SpaceVariableSecretDTO, error)
}

// NewSpaceVariableService creates a new instance of the space secret variable.
func NewSpaceVariableService(
	permission commonapp.ResourcePermissionAppService,
	repoAdapter spacerepo.SpaceRepositoryAdapter,
	variableAdapter spacerepo.SpaceVariableRepositoryAdapter,
	secureStorageAdapter securestorage.SpaceSecureManager,
) SpaceVariableService {
	return &spaceVariableService{
		permission:           permission,
		repoAdapter:          repoAdapter,
		variableAdapter:      variableAdapter,
		secureStorageAdapter: secureStorageAdapter,
	}
}

type spaceVariableService struct {
	permission           commonapp.ResourcePermissionAppService
	repoAdapter          spacerepo.SpaceRepositoryAdapter
	variableAdapter      spacerepo.SpaceVariableRepositoryAdapter
	secureStorageAdapter securestorage.SpaceSecureManager
}

// Create creates a new space with the given command and returns the ID of the created space.
func (s *spaceVariableService) CreateVariable(
	user primitive.Account,
	spaceId primitive.Identity,
	cmd *CmdToCreateSpaceVariable) (res string, action string, err error) {
	space, err := s.repoAdapter.FindById(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(err)
		}

		return
	}

	action = fmt.Sprintf(
		"owner:%s create spaceId:%s space variable name:%s value:%s",
		spaceId.Identity(), space.Owner.Account(), cmd.Name, cmd.Value,
	)

	err = s.permission.CanCreate(user, space.Owner, primitive.ObjTypeSpace)
	if err != nil {
		return "", action, err
	}

	now := utils.Now()

	variable := &domain.SpaceVariable{
		SpaceId:   space.Id,
		Desc:      cmd.Desc,
		Name:      cmd.Name,
		Value:     cmd.Value,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.spaceVariableCountCheck(space.Id); err != nil {
		return "", action, err
	}

	es := domain.NewSpaceVariableVault(variable)
	err = s.secureStorageAdapter.SaveSpaceEnvSecret(es)
	if err != nil {
		logrus.Errorf("failed to create space variable, space variable name:%s", variable.Name.MSDName())
		return "", action, err
	}

	err = s.variableAdapter.AddVariable(variable)

	return "successful", action, err
}

func (s *spaceVariableService) spaceVariableCountCheck(spaceId primitive.Identity) error {
	total, err := s.variableAdapter.CountVariable(spaceId)
	if err != nil {
		return err
	}

	if total >= config.MaxCountSpaceVariable {
		err = fmt.Errorf("space varibale count(now:%d max:%d) exceed", total, config.MaxCountSpaceVariable)
		return newSpaceVariableCountExceeded(err)
	}

	return nil
}

// DeleteVariable deletes the space variable with the given space ID and variable ID and returns the action performed.
func (s *spaceVariableService) DeleteVariable(
	user primitive.Account,
	spaceId primitive.Identity,
	variableId primitive.Identity,
) (action string, err error) {
	space, err := s.repoAdapter.FindById(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(err)
		}

		return
	}

	variable, err := s.variableAdapter.FindVariableById(variableId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceVariableNotFound(err)
		}

		return
	}

	action = fmt.Sprintf(
		"%s:%s delete space variable name:%s value:%s",
		spaceId.Identity(), space.Owner.Account(), variable.Name.MSDName(), variable.Value.MSDName(),
	)

	notFound, err := commonapp.CanDeleteOrNotFound(user, &space, s.permission)
	if err != nil {
		return
	}
	if notFound {
		err = newSpaceNotFound(fmt.Errorf("%s not found", spaceId.Identity()))

		return
	}

	err = s.secureStorageAdapter.DeleteSpaceEnvSecret(domain.VariablePath+space.Id.Identity(), variable.Name.MSDName())
	if err != nil {
		logrus.Errorf("failed to delete variable, variable id:%s", variable.Id.Identity())
		return
	}

	if err = s.variableAdapter.DeleteVariable(variable.Id); err != nil {
		logrus.Errorf("failed to delete space variable, space variable id:%s", variable.Id.Identity())
		return
	}

	return
}

// Update updates the space with the given space ID using the provided command and returns the action performed.
func (s *spaceVariableService) UpdateVariable(
	user primitive.Account, spaceId primitive.Identity,
	variableId primitive.Identity, cmd *CmdToUpdateSpaceVariable,
) (action string, err error) {
	space, err := s.repoAdapter.FindById(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(err)
		}

		return
	}

	variable, err := s.variableAdapter.FindVariableById(variableId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceVariableNotFound(err)
		}

		return
	}

	action = fmt.Sprintf(
		"%s:%s update space variable name:%s",
		spaceId.Identity(), space.Owner.Account(), variable.Name.MSDName(),
	)

	notFound, err := commonapp.CanUpdateOrNotFound(user, &space, s.permission)
	if err != nil {
		return
	}
	if notFound {
		err = newSpaceNotFound(fmt.Errorf("%s not found", spaceId.Identity()))

		return
	}

	b := cmd.toSpaceVariable(&variable)
	if !b {
		logrus.Errorf("failed to change variable, variable id:%s", variable.Id.Identity())
		return
	}

	es := domain.NewSpaceVariableVault(&variable)
	err = s.secureStorageAdapter.SaveSpaceEnvSecret(es)
	if err != nil {
		logrus.Errorf("failed to update variable, variable id:%s", variable.Id.Identity())
		return
	}

	err = s.variableAdapter.SaveVariable(&variable)
	return
}

// List retrieves a list of spaces based on the provided command parameters and returns the corresponding SpacesDTO.
func (s *spaceVariableService) ListVariableSecret(spaceId string) (
	SpaceVariableSecretDTO, error,
) {

	variableSecretList, err := s.variableAdapter.ListVariableSecret(spaceId)

	return SpaceVariableSecretDTO{
		SpaceVariableSecret: variableSecretList,
	}, err
}
