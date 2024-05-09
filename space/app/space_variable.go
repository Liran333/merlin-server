/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package app

import (
	"fmt"

	"golang.org/x/xerrors"

	"github.com/openmerlin/merlin-server/common/app"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/space/domain"
	"github.com/openmerlin/merlin-server/space/domain/message"
	spacerepo "github.com/openmerlin/merlin-server/space/domain/repository"
	"github.com/openmerlin/merlin-server/space/domain/securestorage"
	appprimitive "github.com/openmerlin/merlin-server/spaceapp/domain/primitive"
	"github.com/openmerlin/merlin-server/spaceapp/domain/repository"
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
	permission app.ResourcePermissionAppService,
	repoAdapter spacerepo.SpaceRepositoryAdapter,
	repo repository.Repository,
	variableAdapter spacerepo.SpaceVariableRepositoryAdapter,
	secureStorageAdapter securestorage.SpaceSecureManager,
	msgAdapter message.SpaceMessage,
) SpaceVariableService {
	return &spaceVariableService{
		permission:           permission,
		repoAdapter:          repoAdapter,
		repo:                 repo,
		variableAdapter:      variableAdapter,
		secureStorageAdapter: secureStorageAdapter,
		msgAdapter:           msgAdapter,
	}
}

type spaceVariableService struct {
	permission           app.ResourcePermissionAppService
	repoAdapter          spacerepo.SpaceRepositoryAdapter
	repo                 repository.Repository
	variableAdapter      spacerepo.SpaceVariableRepositoryAdapter
	secureStorageAdapter securestorage.SpaceSecureManager
	msgAdapter           message.SpaceMessage
}

func (s *spaceVariableService) setAppRestarting(spaceId primitive.Identity) error {
	app, err := s.repo.FindBySpaceId(spaceId)
	if err != nil {
		return nil
	}
	if app.Status.IsPaused() || app.Status.IsResuming() || app.Status.IsResumeFailed() {
		return nil
	}
	app.Status = appprimitive.AppStatusRestarted
	return s.repo.Save(&app)
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
		"add space variable of %s:%s/%s:%s",
		spaceId.Identity(), space.Owner.Account(), cmd.Name.ENVName(), cmd.Value.ENVValue(),
	)

	err = s.permission.CanCreate(user, space.Owner, primitive.ObjTypeSpace)
	if err != nil {
		err = newSpaceNotFound(err)
		return "", action, err
	}

	if err := s.spaceVariableCountCheck(space.Id); err != nil {
		err = newSpaceVariableCountExceeded(err)
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
	es := domain.NewSpaceVariableVault(variable)
	err = s.secureStorageAdapter.SaveSpaceEnvSecret(es)
	if err != nil {
		err = allerror.NewCommonRespError("failed to create space variable",
			xerrors.Errorf("space variable name:%s, err: %w", variable.Name.ENVName(), err))
		return "", action, err
	}

	if err = s.variableAdapter.AddVariable(variable); err != nil {
		err = allerror.NewCommonRespError("failed to create space variable db",
			xerrors.Errorf("space variable name:%s, err: %w", variable.Name.ENVName(), err))
		return "", action, err
	}

	e := domain.NewSpaceEnvChangedEvent(user, &space)
	if err = s.msgAdapter.SendSpaceEnvChangedEvent(&e); err != nil {
		err = allerror.NewCommonRespError("failed to send add space variable event",
			xerrors.Errorf("space id:%s, err: %w", spaceId.Identity(), err))
		return "", action, err
	}
	if err = s.setAppRestarting(space.Id); err != nil {
		err = allerror.NewCommonRespError("failed to restart space app",
			xerrors.Errorf("space id:%s, err: %w", spaceId.Identity(), err))
		return "", action, err
	}
	return "successful", action, err
}

func (s *spaceVariableService) spaceVariableCountCheck(spaceId primitive.Identity) error {
	total, err := s.variableAdapter.CountVariable(spaceId)
	if err != nil {
		return err
	}

	if total >= config.MaxCountSpaceVariable {
		err = fmt.Errorf("space varibale count(now:%d max:%d) exceed", total, config.MaxCountSpaceVariable)
		return err
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
		"delete space variable of %s:%s/%s",
		spaceId.Identity(), space.Owner.Account(), variable.Name.ENVName(),
	)

	notFound, err := app.CanDeleteOrNotFound(user, &space, s.permission)
	if err != nil {
		return
	}
	if notFound {
		err = newSpaceNotFound(xerrors.Errorf("%s not found", spaceId.Identity()))
		return
	}

	err = s.secureStorageAdapter.DeleteSpaceEnvSecret(variable.GetVariablePath(), variable.Name.ENVName())
	if err != nil {
		err = allerror.NewCommonRespError("failed to delete space variable",
			xerrors.Errorf("space variable name:%s, err:%w", variable.Name.ENVName(), err))
		return
	}

	if err = s.variableAdapter.DeleteVariable(variable.Id); err != nil {
		err = allerror.NewCommonRespError("failed to delete space variable db",
			xerrors.Errorf("space variable name:%s, err:%w", variable.Name.ENVName(), err))
		return
	}

	e := domain.NewSpaceEnvChangedEvent(user, &space)
	if err = s.msgAdapter.SendSpaceEnvChangedEvent(&e); err != nil {
		err = allerror.NewCommonRespError("failed to send delete space variable event",
			xerrors.Errorf("space id:%s, err:%w", spaceId.Identity(), err))
		return
	}
	if err = s.setAppRestarting(space.Id); err != nil {
		err = allerror.NewCommonRespError("failed to restart space app",
			xerrors.Errorf("space id:%s, err:%w", spaceId.Identity(), err))
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
		"update space variable of %s:%s/%s",
		spaceId.Identity(), space.Owner.Account(), variable.Name.ENVName(),
	)

	notFound, err := app.CanUpdateOrNotFound(user, &space, s.permission)
	if err != nil {
		return
	}
	if notFound {
		err = newSpaceNotFound(xerrors.Errorf("%s not found", spaceId.Identity()))
		return
	}

	b := cmd.toSpaceVariable(&variable)
	if !b {
		err = newSpaceVariableNotFound(xerrors.Errorf("%s not found", variableId.Identity()))
		return
	}

	es := domain.NewSpaceVariableVault(&variable)
	err = s.secureStorageAdapter.SaveSpaceEnvSecret(es)
	if err != nil {
		err = allerror.NewCommonRespError("failed to update space variable",
			xerrors.Errorf("space variable name:%s, err:%w", variable.Name.ENVName(), err))
		return
	}

	if err = s.variableAdapter.SaveVariable(&variable); err != nil {
		err = allerror.NewCommonRespError("failed to update space variable db",
			xerrors.Errorf("space variable name:%s, err:%w", variable.Name.ENVName(), err))
		return
	}

	e := domain.NewSpaceEnvChangedEvent(user, &space)
	if err = s.msgAdapter.SendSpaceEnvChangedEvent(&e); err != nil {
		err = allerror.NewCommonRespError("failed to send update space variable event",
			xerrors.Errorf("space id:%s, err:%w", spaceId.Identity(), err))
		return
	}
	if err = s.setAppRestarting(space.Id); err != nil {
		err = allerror.NewCommonRespError("failed to restart space app",
			xerrors.Errorf("space id:%s, err:%w", spaceId.Identity(), err))
	}
	return
}

// List retrieves a list of spaces based on the provided command parameters and returns the corresponding SpacesDTO.
func (s *spaceVariableService) ListVariableSecret(spaceId string) (
	SpaceVariableSecretDTO, error,
) {
	variableSecretList, err := s.variableAdapter.ListVariableSecret(spaceId)
	if err != nil {
		err = newSpaceNotFound(xerrors.Errorf("%s not found", spaceId))
		return SpaceVariableSecretDTO{}, err
	}

	return SpaceVariableSecretDTO{
		SpaceVariableSecret: variableSecretList,
	}, err
}
