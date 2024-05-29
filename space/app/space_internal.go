/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"fmt"

	sdk "github.com/openmerlin/merlin-sdk/space"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"

	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	modelrepo "github.com/openmerlin/merlin-server/models/domain/repository"
	"github.com/openmerlin/merlin-server/space/domain"
	"github.com/openmerlin/merlin-server/space/domain/message"
	"github.com/openmerlin/merlin-server/space/domain/repository"
	spacerepo "github.com/openmerlin/merlin-server/space/domain/repository"
	spaceappRepository "github.com/openmerlin/merlin-server/spaceapp/domain/repository"
	"github.com/openmerlin/merlin-server/utils"
)

// SpaceInternalAppService is an interface for space internal application service
type SpaceInternalAppService interface {
	GetById(primitive.Identity) (sdk.SpaceMetaDTO, error)
	UpdateLocalCMD(spaceId primitive.Identity, cmd string) error
	UpdateEnvInfo(spaceId primitive.Identity, envInfo string) error
	UpdateStatistics(primitive.Identity, *CmdToUpdateStatistics) error
	Disable(primitive.Identity) error
	ResetLabels(primitive.Identity, *CmdToResetLabels) error
	NotifyUpdateCodes(primitive.Identity, *CmdToNotifyUpdateCode) error
}

// NewSpaceInternalAppService creates a new instance of SpaceInternalAppService
func NewSpaceInternalAppService(
	repoAdapter repository.SpaceRepositoryAdapter,
	msgAdapter message.SpaceMessage,
	spaceappRepository spaceappRepository.Repository,
	repoAdapterModelSpace spacerepo.ModelSpaceRepositoryAdapter,
	modelRepoAdapter modelrepo.ModelRepositoryAdapter,
) SpaceInternalAppService {
	return &spaceInternalAppService{
		repoAdapter:           repoAdapter,
		msgAdapter:            msgAdapter,
		spaceappRepository:    spaceappRepository,
		repoAdapterModelSpace: repoAdapterModelSpace,
		modelRepoAdapter:      modelRepoAdapter,
	}
}

type spaceInternalAppService struct {
	repoAdapter           repository.SpaceRepositoryAdapter
	msgAdapter            message.SpaceMessage
	spaceappRepository    spaceappRepository.Repository
	repoAdapterModelSpace spacerepo.ModelSpaceRepositoryAdapter
	modelRepoAdapter      modelrepo.ModelRepositoryAdapter
}

// GetById retrieves a space by its ID and returns the corresponding SpaceMetaDTO
func (s *spaceInternalAppService) GetById(spaceId primitive.Identity) (sdk.SpaceMetaDTO, error) {
	space, err := s.repoAdapter.FindById(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(err)
		}

		return sdk.SpaceMetaDTO{}, err
	}

	return toSpaceMetaDTO(&space), nil
}

func (s *spaceInternalAppService) UpdateLocalCMD(spaceId primitive.Identity, cmd string) error {
	space, err := s.repoAdapter.FindById(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(err)
		}

		return err
	}

	space.LocalCmd = cmd
	return s.repoAdapter.Save(&space)
}

func (s *spaceInternalAppService) UpdateEnvInfo(spaceId primitive.Identity, envInfo string) error {
	space, err := s.repoAdapter.FindById(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(err)
		}

		return err
	}

	space.LocalEnvInfo = envInfo
	return s.repoAdapter.Save(&space)
}

// UpdateStatistics updates the statistics of a space.
func (s *spaceInternalAppService) UpdateStatistics(spaceId primitive.Identity, cmd *CmdToUpdateStatistics) error {
	space, err := s.repoAdapter.FindById(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(
				allerror.ErrorCodeSpaceNotFound,
				"not found",
				fmt.Errorf("%s not found, err: %w", spaceId.Identity(), err))
		}

		return err
	}

	space.DownloadCount = cmd.DownloadCount
	space.UpdatedAt = utils.Now()

	return s.repoAdapter.Save(&space)
}

// Disable disable the space with the given space ID using the provided command and returns the action performed.
func (s *spaceInternalAppService) Disable(spaceId primitive.Identity) (err error) {
	space, err := s.repoAdapter.FindById(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(err)
		}

		return
	}

	// del space app
	_, err = s.spaceappRepository.FindBySpaceId(spaceId)
	if err != nil && !commonrepo.IsErrorResourceNotExists(err) {
		logrus.Errorf("get space app by id %v failed, err:%v", spaceId, err)
		return
	} else if err == nil {
		if err = s.spaceappRepository.DeleteBySpaceId(spaceId); err != nil {
			logrus.Errorf("delete space app by id %v failed, err:%v", spaceId, err)
			return
		}
	}

	space.Exception = primitive.CreateException(primitive.RelatedModelDisabled)

	if err = s.repoAdapter.Save(&space); err != nil {
		return
	}

	e := domain.NewSpaceForceEvent(space.Id.Identity(), domain.ForceTypeStop)
	if err1 := s.msgAdapter.SendSpaceForceEvent(&e); err1 != nil {
		logrus.Errorf("failed to send space force stop event, space id:%s", spaceId.Identity())
	}

	logrus.Infof("send space force stop event success, space id:%s", spaceId.Identity())

	return
}

type SpaceMetaDTO1 = sdk.SpaceMetaDTO

func (s *spaceInternalAppService) ResetLabels(spaceId primitive.Identity, cmd *CmdToResetLabels) error {
	space, err := s.repoAdapter.FindById(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(xerrors.Errorf("not found, err: %w", err))
		} else {
			err = xerrors.Errorf("find space by id failed, err: %w", err)
		}

		return err
	}

	if space.IsDisable() {
		errInfo := fmt.Errorf("space %v was disable", space.Name.MSDName())
		return allerror.NewResourceDisabled(allerror.ErrorCodeResourceDisabled, errInfo.Error(), errInfo)
	}

	if cmd.Task != nil {
		space.Labels.Task = cmd.Task
	}

	if cmd.License != nil {
		space.License = cmd.License
	}

	err = s.repoAdapter.Save(&space)
	if err != nil {
		err = xerrors.Errorf("save space failed, err: %w", err)
	}

	return err
}

// NotifyUpdateCodes notify the space no application file and commitId with the given space ID
func (s *spaceInternalAppService) NotifyUpdateCodes(spaceId primitive.Identity, cmd *CmdToNotifyUpdateCode) error {
	space, err := s.repoAdapter.FindById(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(xerrors.Errorf("not found, err: %w", err))
		} else {
			err = xerrors.Errorf("find space by id failed, err: %w", err)
		}

		return err
	}

	space.SetSpaceCommitId(cmd.CommitId)
	space.SetNoApplicationFile(cmd.HasHtml, cmd.HasApp)
	err = s.repoAdapter.Save(&space)

	if err != nil {
		err = xerrors.Errorf("save space failed, err: %w", err)
	}
	logrus.Infof("spaceId:%s set notify commitId:%s, req:%v success",
		spaceId, cmd.CommitId, cmd)
	return err
}
