/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package app

import (
	"context"
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

// SpaceSecretService is an interface for the space secret service.
type SpaceSecretService interface {
	CreateSecret(
		context.Context, primitive.Account,
		primitive.Identity, *CmdToCreateSpaceSecret) (string, string, error)
	DeleteSecret(context.Context, primitive.Account, primitive.Identity, primitive.Identity) (string, error)
	UpdateSecret(
		context.Context, primitive.Account, primitive.Identity,
		primitive.Identity, *CmdToUpdateSpaceSecret) (string, error)
}

// NewSpaceSecretService creates a new instance of the space secret service.
func NewSpaceSecretService(
	permission app.ResourcePermissionAppService,
	repoAdapter spacerepo.SpaceRepositoryAdapter,
	repo repository.Repository,
	secretAdapter spacerepo.SpaceSecretRepositoryAdapter,
	secureStorageAdapter securestorage.SpaceSecureManager,
	msgAdapter message.SpaceMessage,
) SpaceSecretService {
	return &spaceSecretService{
		permission:           permission,
		repoAdapter:          repoAdapter,
		repo:                 repo,
		secretAdapter:        secretAdapter,
		secureStorageAdapter: secureStorageAdapter,
		msgAdapter:           msgAdapter,
	}
}

type spaceSecretService struct {
	permission           app.ResourcePermissionAppService
	repoAdapter          spacerepo.SpaceRepositoryAdapter
	repo                 repository.Repository
	secretAdapter        spacerepo.SpaceSecretRepositoryAdapter
	secureStorageAdapter securestorage.SpaceSecureManager
	msgAdapter           message.SpaceMessage
}

func (s *spaceSecretService) setAppRestarting(ctx context.Context, spaceId primitive.Identity) error {
	app, err := s.repo.FindBySpaceId(ctx, spaceId)
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
func (s *spaceSecretService) CreateSecret(
	ctx context.Context,
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
		"add space secret of %s:%s/%s",
		spaceId.Identity(), space.Owner.Account(), cmd.Name.ENVName(),
	)

	err = s.permission.CanCreate(ctx, user, space.Owner, primitive.ObjTypeSpace)
	if err != nil {
		err = newSpaceNotFound(err)
		return "", action, err
	}

	if err = s.spaceSecretCountCheck(space.Id); err != nil {
		err = newSpaceSecretCountExceeded(err)
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
		err = allerror.NewCommonRespError("failed to create space secret",
			xerrors.Errorf("space secret name:%s, err: %w", secret.Name.ENVName(), err))
		return "", action, err
	}

	if err = s.secretAdapter.AddSecret(secret); err != nil {
		err = allerror.NewCommonRespError("failed to create space secret db",
			xerrors.Errorf("space secret name:%s, err: %w", secret.Name.ENVName(), err))
		return "", action, err
	}

	e := domain.NewSpaceEnvChangedEvent(user, &space)
	if err = s.msgAdapter.SendSpaceEnvChangedEvent(&e); err != nil {
		err = allerror.NewCommonRespError("failed to send create space secret event",
			xerrors.Errorf("space id:%s, err: %w", spaceId.Identity(), err))
		return "", action, err
	}
	if err = s.setAppRestarting(ctx, space.Id); err != nil {
		err = allerror.NewCommonRespError("failed to restart space app",
			xerrors.Errorf("space id:%s, err: %w", spaceId.Identity(), err))
		return "", action, err
	}
	return "successful", action, err
}

func (s *spaceSecretService) spaceSecretCountCheck(spaceId primitive.Identity) error {
	total, err := s.secretAdapter.CountSecret(spaceId)
	if err != nil {
		return err
	}

	if total >= config.MaxCountSpaceSecret {
		err = fmt.Errorf("space secret count(now:%d max:%d) exceed", total, config.MaxCountSpaceSecret)
		return err
	}

	return nil
}

// DeleteVariable deletes the space variable with the given space ID and variable ID and returns the action performed.
func (s *spaceSecretService) DeleteSecret(
	ctx context.Context,
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
		"delete space secret of %s:%s/%s",
		secretId.Identity(), space.Owner.Account(), secret.Name.ENVName(),
	)

	notFound, err := app.CanDeleteOrNotFound(ctx, user, &space, s.permission)
	if err != nil {
		return
	}
	if notFound {
		err = newSpaceNotFound(xerrors.Errorf("%s not found", spaceId.Identity()))
		return
	}

	err = s.secureStorageAdapter.DeleteSpaceEnvSecret(secret.GetSecretPath(), secret.Name.ENVName())
	if err != nil {
		err = allerror.NewCommonRespError("failed to delete secret",
			xerrors.Errorf("space secret name:%s, err: %w", secret.Name.ENVName(), err))
		return
	}

	if err = s.secretAdapter.DeleteSecret(secret.Id); err != nil {
		err = allerror.NewCommonRespError("failed to delete secret db",
			xerrors.Errorf("space secret name:%s, err: %w", secret.Name.ENVName(), err))
		return
	}

	e := domain.NewSpaceEnvChangedEvent(user, &space)
	if err = s.msgAdapter.SendSpaceEnvChangedEvent(&e); err != nil {
		err = allerror.NewCommonRespError("failed to send delete space secret event",
			xerrors.Errorf("space id:%s, err: %w", spaceId.Identity(), err))
		return
	}
	if err = s.setAppRestarting(ctx, space.Id); err != nil {
		err = allerror.NewCommonRespError("failed to restart space app",
			xerrors.Errorf("space id:%s, err: %w", spaceId.Identity(), err))
	}
	return
}

// Update updates the space with the given space ID using the provided command and returns the action performed.
func (s *spaceSecretService) UpdateSecret(
	ctx context.Context, user primitive.Account, spaceId primitive.Identity,
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
		"update space secret of %s:%s/%s",
		spaceId.Identity(), space.Owner.Account(), secret.Name.ENVName(),
	)

	notFound, err := app.CanUpdateOrNotFound(ctx, user, &space, s.permission)
	if err != nil {
		return
	}
	if notFound {
		err = newSpaceNotFound(xerrors.Errorf("%s not found", spaceId.Identity()))
		return
	}

	b := cmd.toSpaceSecret(&secret)
	if !b {
		return
	}

	es := domain.NewSpaceSecretVault(&secret)
	err = s.secureStorageAdapter.SaveSpaceEnvSecret(es)
	if err != nil {
		err = allerror.NewCommonRespError("failed to update secret",
			xerrors.Errorf("space secret name:%s, err: %w", secret.Name.ENVName(), err))
		return
	}

	err = s.secretAdapter.SaveSecret(&secret)
	if err != nil {
		err = allerror.NewCommonRespError("failed to update secret db",
			xerrors.Errorf("space secret name:%s, err: %w", secret.Name.ENVName(), err))
		return
	}

	e := domain.NewSpaceEnvChangedEvent(user, &space)
	if err = s.msgAdapter.SendSpaceEnvChangedEvent(&e); err != nil {
		err = allerror.NewCommonRespError("failed to send update space secret event",
			xerrors.Errorf("space id:%s, err: %w", spaceId.Identity(), err))
		return
	}
	if err = s.setAppRestarting(ctx, space.Id); err != nil {
		err = allerror.NewCommonRespError("failed to restart space app",
			xerrors.Errorf("space id:%s, err: %w", spaceId.Identity(), err))
	}
	return
}
