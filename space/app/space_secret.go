/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package app

import (
	"fmt"

	"github.com/sirupsen/logrus"

	commonapp "github.com/openmerlin/merlin-server/common/app"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/space/domain"
	spacerepo "github.com/openmerlin/merlin-server/space/domain/repository"
	"github.com/openmerlin/merlin-server/space/domain/securestorage"
	"github.com/openmerlin/merlin-server/utils"
)

// SpaceSecretService is an interface for the space secret service.
type SpaceSecretService interface {
	CreateSecret(primitive.Account, primitive.Identity, *CmdToCreateSpaceSecret) (string, string, error)
	DeleteSecret(primitive.Account, primitive.Identity, primitive.Identity) (string, error)
	UpdateSecret(primitive.Account, primitive.Identity, primitive.Identity, *CmdToUpdateSpaceSecret) (string, error)
}

// NewSpaceSecretService creates a new instance of the space secret service.
func NewSpaceSecretService(
	permission commonapp.ResourcePermissionAppService,
	repoAdapter spacerepo.SpaceRepositoryAdapter,
	secretAdapter spacerepo.SpaceSecretRepositoryAdapter,
	secureStorageAdapter securestorage.SpaceSecureManager,
) SpaceSecretService {
	return &spaceSecretService{
		permission:           permission,
		repoAdapter:          repoAdapter,
		secretAdapter:        secretAdapter,
		secureStorageAdapter: secureStorageAdapter,
	}
}

type spaceSecretService struct {
	permission           commonapp.ResourcePermissionAppService
	repoAdapter          spacerepo.SpaceRepositoryAdapter
	secretAdapter        spacerepo.SpaceSecretRepositoryAdapter
	secureStorageAdapter securestorage.SpaceSecureManager
}

// Create creates a new space with the given command and returns the ID of the created space.
func (s *spaceSecretService) CreateSecret(
	user primitive.Account,
	spaceId primitive.Identity,
	cmd *CmdToCreateSpaceSecret) (res string, action string, err error) {
	space, err := s.repoAdapter.FindById(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(err)
		}

		return
	}

	action = fmt.Sprintf(
		"owner:%s create spaceId:%s space secret name:%s",
		spaceId.Identity(), space.Owner.Account(), cmd.Name,
	)

	err = s.permission.CanCreate(user, space.Owner, primitive.ObjTypeSpace)
	if err != nil {
		return "", action, err
	}

	if err := s.spaceSecretCountCheck(space.Id); err != nil {
		return "", action, err
	}

	now := utils.Now()
	secret := &domain.SpaceSecret{
		SpaceId:   space.Id,
		Desc:      cmd.Desc,
		Name:      cmd.Name,
		Value:     cmd.Value,
		CreatedAt: now,
		UpdatedAt: now,
	}
	es := domain.NewSpaceSecretVault(secret)
	err = s.secureStorageAdapter.SaveSpaceEnvSecret(es)
	if err != nil {
		logrus.Errorf("failed to create space secret, space secret id:%s", secret.Name.MSDName())
		return "", action, err
	}

	err = s.secretAdapter.AddSecret(secret)

	return "successful", action, err
}

func (s *spaceSecretService) spaceSecretCountCheck(spaceId primitive.Identity) error {
	total, err := s.secretAdapter.CountSecret(spaceId)
	if err != nil {
		return err
	}

	if total >= config.MaxCountSpaceSecret {
		err = fmt.Errorf("space secret count(now:%d max:%d) exceed", total, config.MaxCountSpaceSecret)
		return newSpaceSecretCountExceeded(err)
	}

	return nil
}

// DeleteVariable deletes the space variable with the given space ID and variable ID and returns the action performed.
func (s *spaceSecretService) DeleteSecret(
	user primitive.Account,
	spaceId primitive.Identity,
	secretId primitive.Identity,
) (action string, err error) {
	space, err := s.repoAdapter.FindById(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(err)
		}

		return
	}

	secret, err := s.secretAdapter.FindSecretById(secretId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceSecretNotFound(err)
		}

		return
	}

	action = fmt.Sprintf(
		"%s:%s delete space secret name:%s",
		spaceId.Identity(), space.Owner.Account(), secret.Name.MSDName(),
	)

	notFound, err := commonapp.CanDeleteOrNotFound(user, &space, s.permission)
	if err != nil {
		return
	}
	if notFound {
		err = newSpaceNotFound(fmt.Errorf("%s not found", spaceId.Identity()))

		return
	}

	err = s.secureStorageAdapter.DeleteSpaceEnvSecret(space.Id.Identity(), domain.SecretePath+secret.Name.MSDName())
	if err != nil {
		return
	}

	if err = s.secretAdapter.DeleteSecret(secret.Id); err != nil {
		return
	}

	return
}

// Update updates the space with the given space ID using the provided command and returns the action performed.
func (s *spaceSecretService) UpdateSecret(
	user primitive.Account, spaceId primitive.Identity,
	secretId primitive.Identity, cmd *CmdToUpdateSpaceSecret,
) (action string, err error) {
	space, err := s.repoAdapter.FindById(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(err)
		}

		return
	}

	secret, err := s.secretAdapter.FindSecretById(secretId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceSecretNotFound(err)
		}

		return
	}

	action = fmt.Sprintf(
		"%s:%s update space secret name:%s",
		spaceId.Identity(), space.Owner.Account(), secret.Name.MSDName(),
	)

	notFound, err := commonapp.CanUpdateOrNotFound(user, &space, s.permission)
	if err != nil {
		return
	}
	if notFound {
		err = newSpaceNotFound(fmt.Errorf("%s not found", spaceId.Identity()))

		return
	}

	b := cmd.toSpaceSecret(&secret)
	if !b {
		return
	}

	es := domain.NewSpaceSecretVault(&secret)
	err = s.secureStorageAdapter.SaveSpaceEnvSecret(es)
	if err != nil {
		return
	}

	err = s.secretAdapter.SaveSecret(&secret)
	return
}
