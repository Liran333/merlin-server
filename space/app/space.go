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
	"github.com/openmerlin/merlin-server/space/domain"
	"github.com/openmerlin/merlin-server/space/domain/message"
	"github.com/openmerlin/merlin-server/space/domain/repository"
	"github.com/openmerlin/merlin-server/space/domain/securestorage"
	spaceappRepository "github.com/openmerlin/merlin-server/spaceapp/domain/repository"
	"github.com/openmerlin/merlin-server/utils"
)

const (
	variableTypeName = "variable"
	secretTypeName   = "secret"
)

func newSpaceNotFound(err error) error {
	return allerror.NewNotFound(allerror.ErrorCodeSpaceNotFound, "not found", err)
}

func newSpaceCountExceeded(err error) error {
	return allerror.NewCountExceeded("space count exceed", err)
}

func newModelNotFound(err error) error {
	return allerror.NewNotFound(allerror.ErrorCodeModelNotFound, "not found", err)
}

func newSpaceSecretNotFound(err error) error {
	return allerror.NewNotFound(allerror.ErrorCodeSpaceSecretNotFound, "not found", err)
}

func newSpaceSecretCountExceeded(err error) error {
	return allerror.NewCountExceeded("space secret count exceed", err)
}

func newSpaceVariableNotFound(err error) error {
	return allerror.NewNotFound(allerror.ErrorCodeSpaceVariableNotFound, "not found", err)
}

func newSpaceVariableCountExceeded(err error) error {
	return allerror.NewCountExceeded("space variable count exceed", err)
}

// Permission is an interface for checking permissions.
type Permission interface {
	Check(primitive.Account, primitive.Account, primitive.ObjType, primitive.Action) error
}

// SpaceAppService is an interface for space application services.
type SpaceAppService interface {
	Create(primitive.Account, *CmdToCreateSpace) (string, error)
	Delete(primitive.Account, primitive.Identity) (string, error)
	Update(primitive.Account, primitive.Identity, *CmdToUpdateSpace) (string, error)
	GetByName(primitive.Account, *domain.SpaceIndex) (SpaceDTO, error)
	List(primitive.Account, *CmdToListSpaces) (SpacesDTO, error)
	AddLike(primitive.Identity) error
	DeleteLike(primitive.Identity) error
}

// NewSpaceAppService creates a new instance of SpaceAppService.
func NewSpaceAppService(
	permission commonapp.ResourcePermissionAppService,
	msgAdapter message.SpaceMessage,
	codeRepoApp coderepoapp.CodeRepoAppService,
	spaceappRepository spaceappRepository.Repository,
	variableAdapter repository.SpaceVariableRepositoryAdapter,
	secretAdapter repository.SpaceSecretRepositoryAdapter,
	secureStorageAdapter securestorage.SpaceSecureManager,
	repoAdapter repository.SpaceRepositoryAdapter,
) SpaceAppService {
	return &spaceAppService{
		permission:           permission,
		msgAdapter:           msgAdapter,
		codeRepoApp:          codeRepoApp,
		spaceappRepository:   spaceappRepository,
		variableAdapter:      variableAdapter,
		secretAdapter:        secretAdapter,
		secureStorageAdapter: secureStorageAdapter,
		repoAdapter:          repoAdapter,
	}
}

type spaceAppService struct {
	permission           commonapp.ResourcePermissionAppService
	msgAdapter           message.SpaceMessage
	codeRepoApp          coderepoapp.CodeRepoAppService
	spaceappRepository   spaceappRepository.Repository
	variableAdapter      repository.SpaceVariableRepositoryAdapter
	secretAdapter        repository.SpaceSecretRepositoryAdapter
	secureStorageAdapter securestorage.SpaceSecureManager
	repoAdapter          repository.SpaceRepositoryAdapter
}

// Create creates a new space with the given command and returns the ID of the created space.
func (s *spaceAppService) Create(user primitive.Account, cmd *CmdToCreateSpace) (string, error) {
	err := s.permission.CanCreate(user, cmd.Owner, primitive.ObjTypeSpace)
	if err != nil {
		return "", err
	}

	if err := s.spaceCountCheck(cmd.Owner); err != nil {
		return "", err
	}

	coderepo, err := s.codeRepoApp.Create(user, &cmd.CmdToCreateRepo)
	if err != nil {
		return "", err
	}

	now := utils.Now()
	space := domain.Space{
		SDK:       cmd.SDK,
		Desc:      cmd.Desc,
		Hardware:  cmd.Hardware,
		Fullname:  cmd.Fullname,
		CodeRepo:  coderepo,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err = s.repoAdapter.Add(&space); err != nil {
		return "", err
	}

	e := domain.NewSpaceCreatedEvent(&space)
	if err1 := s.msgAdapter.SendSpaceCreatedEvent(&e); err1 != nil {
		logrus.Errorf("failed to send space created event, space id:%s", space.Id.Identity())
	}

	return space.Id.Identity(), nil
}

// Delete deletes the space with the given space ID and returns the action performed.
func (s *spaceAppService) Delete(user primitive.Account, spaceId primitive.Identity) (action string, err error) {
	space, err := s.repoAdapter.FindById(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		}

		return
	}

	action = fmt.Sprintf(
		"delete space of %s:%s/%s",
		spaceId.Identity(), space.Owner.Account(), space.Name.MSDName(),
	)

	notFound, err := commonapp.CanDeleteOrNotFound(user, &space, s.permission)
	if err != nil {
		return
	}
	if notFound {
		err = newSpaceNotFound(fmt.Errorf("%s not found", spaceId.Identity()))

		return
	}

	if err = s.codeRepoApp.Delete(space.RepoIndex()); err != nil {
		return
	}

	// del space app
	if err = s.spaceappRepository.DeleteBySpaceId(space.Id); err != nil {
		return
	}

	// del space variable secret
	if err = s.delSpaceVariableSecret(space.Id); err != nil {
		return
	}

	if err = s.repoAdapter.Delete(space.Id); err != nil {
		return
	}

	e := domain.NewSpaceDeletedEvent(user, &space)
	if err1 := s.msgAdapter.SendSpaceDeletedEvent(&e); err1 != nil {
		logrus.Errorf("failed to send space deleted event, space id:%s", spaceId.Identity())
	}

	return
}

func (s *spaceAppService) delSpaceVariableSecret(spaceId primitive.Identity) error {
	spaceVariableSecretList, err := s.variableAdapter.ListVariableSecret(spaceId.Identity())
	if err != nil {
		return err
	}
	for _, envSecret := range spaceVariableSecretList {
		envSecretId, err := primitive.NewIdentity(envSecret.Id)
		if err != nil {
			logrus.Errorf("failed to get envSecretId, err:%s", err)
			continue
		}
		if envSecret.Type == variableTypeName {
			variable, err := s.variableAdapter.FindVariableById(envSecretId)
			if err != nil {
				logrus.Errorf("failed to get variable, err:%s", err)
				continue
			}
			if err = s.secureStorageAdapter.DeleteSpaceEnvSecret(
				variable.GetVariablePath(), variable.Name.MSDName()); err != nil {
				logrus.Errorf("failed to delete variable, err:%s", err)
				continue
			}
			if err = s.variableAdapter.DeleteVariable(envSecretId); err != nil {
				logrus.Errorf("failed to delete variable db, err:%s", err)
				continue
			}
		}
		if envSecret.Type == secretTypeName {
			secret, err := s.secretAdapter.FindSecretById(envSecretId)
			if err != nil {
				logrus.Errorf("failed to get secret, err:%s", err)
				continue
			}
			if err = s.secureStorageAdapter.DeleteSpaceEnvSecret(
				secret.GetSecretPath(), secret.Name.MSDName()); err != nil {
				logrus.Errorf("failed to delete secret, err:%s", err)
				continue
			}
			if err = s.secretAdapter.DeleteSecret(envSecretId); err != nil {
				logrus.Errorf("failed to delete secret db, err:%s", err)
				continue
			}
		}
	}
	return nil
}

// Update updates the space with the given space ID using the provided command and returns the action performed.
func (s *spaceAppService) Update(
	user primitive.Account, spaceId primitive.Identity, cmd *CmdToUpdateSpace,
) (action string, err error) {
	space, err := s.repoAdapter.FindById(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(err)
		}

		return
	}

	action = fmt.Sprintf(
		"update space of %s:%s/%s",
		spaceId.Identity(), space.Owner.Account(), space.Name.MSDName(),
	)

	notFound, err := commonapp.CanUpdateOrNotFound(user, &space, s.permission)
	if err != nil {
		return
	}
	if notFound {
		err = newSpaceNotFound(fmt.Errorf("%s not found", spaceId.Identity()))

		return
	}

	isPrivateToPublic := space.IsPrivate() && cmd.Visibility.IsPublic()

	b, err := s.codeRepoApp.Update(&space.CodeRepo, &cmd.CmdToUpdateRepo)
	if err != nil {
		return
	}

	b1 := cmd.toSpace(&space)
	if !b && !b1 {
		return
	}

	if err = s.repoAdapter.Save(&space); err != nil {
		return
	}

	e := domain.NewSpaceUpdatedEvent(user, &space, isPrivateToPublic)
	if err1 := s.msgAdapter.SendSpaceUpdatedEvent(&e); err1 != nil {
		logrus.Errorf("failed to send space updated event, space id:%s", spaceId.Identity())
	}

	return
}

// GetByName retrieves a space by its name and returns the corresponding SpaceDTO.
func (s *spaceAppService) GetByName(user primitive.Account, index *domain.SpaceIndex) (SpaceDTO, error) {
	var dto SpaceDTO

	space, err := s.repoAdapter.FindByName(index)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(err)
		}

		return dto, err
	}

	if err := s.permission.CanRead(user, &space); err != nil {
		if allerror.IsNoPermission(err) {
			err = newSpaceNotFound(err)
		}

		return dto, err
	}

	return toSpaceDTO(&space), nil
}

// List retrieves a list of spaces based on the provided command parameters and returns the corresponding SpacesDTO.
func (s *spaceAppService) List(user primitive.Account, cmd *CmdToListSpaces) (
	SpacesDTO, error,
) {
	if user == nil {
		cmd.Visibility = primitive.VisibilityPublic
	} else {
		if cmd.Owner == nil {
			// It can list the private spaces of user,
			// but it maybe no need to do it.
			cmd.Visibility = primitive.VisibilityPublic
		} else {
			if user != cmd.Owner {
				err := s.permission.CanListOrgResource(
					user, cmd.Owner, primitive.ObjTypeSpace,
				)
				if err != nil {
					cmd.Visibility = primitive.VisibilityPublic
				}
			}
		}
	}

	v, total, err := s.repoAdapter.List(cmd)

	return SpacesDTO{
		Total:  total,
		Spaces: v,
	}, err
}

// DeleteById is an example for restful API.
func (s *spaceAppService) DeleteById(user primitive.Account, spaceId string) error {
	// get space by space id
	// check if user can delete it
	// delete it
	return nil
}

func (s *spaceAppService) spaceCountCheck(owner primitive.Account) error {
	cmdToList := CmdToListSpaces{
		Owner: owner,
	}

	total, err := s.repoAdapter.Count(&cmdToList)
	if err != nil {
		return err
	}

	if total >= config.MaxCountPerOwner {
		return newSpaceCountExceeded(fmt.Errorf("space count(now:%d max:%d) exceed", total, config.MaxCountPerOwner))
	}

	return nil
}

func (s *spaceAppService) AddLike(spaceId primitive.Identity) error {
	// Retrieve the code repository information.
	codeRepo, err := s.codeRepoApp.GetById(spaceId)
	if err != nil {
		return err
	}

	// Only proceed if the repository is public.
	isPublic := codeRepo.IsPublic()
	if !isPublic {
		return nil
	}

	if err := s.repoAdapter.AddLike(spaceId); err != nil {
		return err
	}
	return nil
}

func (s *spaceAppService) DeleteLike(spaceId primitive.Identity) error {
	if err := s.repoAdapter.DeleteLike(spaceId); err != nil {
		return err
	}
	return nil
}
