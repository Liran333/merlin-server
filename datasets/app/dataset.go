/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package app provides functionality for the application.
package app

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"

	coderepoapp "github.com/openmerlin/merlin-server/coderepo/app"
	commonapp "github.com/openmerlin/merlin-server/common/app"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/datasets/domain"
	"github.com/openmerlin/merlin-server/datasets/domain/message"
	"github.com/openmerlin/merlin-server/datasets/domain/repository"
	orgapp "github.com/openmerlin/merlin-server/organization/app"
	orgrepo "github.com/openmerlin/merlin-server/organization/domain/repository"
	userapp "github.com/openmerlin/merlin-server/user/app"
	"github.com/openmerlin/merlin-server/utils"
)

// DatasetAppService is an interface for the dataset application service.
type DatasetAppService interface {
	Create(primitive.Account, *CmdToCreateDataset) (string, error)
	Delete(primitive.Account, primitive.Identity) (string, error)
	Update(primitive.Account, primitive.Identity, *CmdToUpdateDataset) (string, error)
	Disable(primitive.Account, primitive.Identity, *CmdToDisableDataset) (string, error)
	GetByName(primitive.Account, *domain.DatasetIndex) (DatasetDTO, error)
	List(primitive.Account, *CmdToListDatasets) (DatasetsDTO, error)
	AddLike(primitive.Identity) error
	DeleteLike(primitive.Identity) error
}

// NewDatasetAppService creates a new instance of the dataset application service.
func NewDatasetAppService(
	permission commonapp.ResourcePermissionAppService,
	msgAdapter message.DatasetMessage,
	codeRepoApp coderepoapp.CodeRepoAppService,
	repoAdapter repository.DatasetRepositoryAdapter,
	member orgrepo.OrgMember,
	disableOrg orgapp.PrivilegeOrg,
	user userapp.UserService,
) DatasetAppService {
	return &datasetAppService{
		permission:  permission,
		msgAdapter:  msgAdapter,
		codeRepoApp: codeRepoApp,
		repoAdapter: repoAdapter,
		member:      member,
		disableOrg:  disableOrg,
		user:        user,
	}
}

type datasetAppService struct {
	permission  commonapp.ResourcePermissionAppService
	msgAdapter  message.DatasetMessage
	codeRepoApp coderepoapp.CodeRepoAppService
	repoAdapter repository.DatasetRepositoryAdapter
	member      orgrepo.OrgMember
	disableOrg  orgapp.PrivilegeOrg
	user        userapp.UserService
}

// Create creates a new dataset.
func (s *datasetAppService) Create(user primitive.Account, cmd *CmdToCreateDataset) (string, error) {
	if err := s.permission.CanCreate(user, cmd.Owner, primitive.ObjTypeDataset); err != nil {

		return "", xerrors.Errorf("permission check failed, err:%w", err)
	}

	if err := s.datasetsCountCheck(cmd.Owner); err != nil {
		return "", xerrors.Errorf("failed to check dataset count, err:%w", err)
	}

	coderepo, err := s.codeRepoApp.Create(user, &cmd.CmdToCreateRepo)
	if err != nil {
		return "", xerrors.Errorf("failed to create dataset code repo, err:%w", err)
	}

	now := utils.Now()
	dataset := domain.Dataset{
		Desc:      cmd.Desc,
		Fullname:  cmd.Fullname,
		CodeRepo:  coderepo,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err = s.repoAdapter.Add(&dataset); err != nil {
		return "", xerrors.Errorf("failed to add dataset info, err:%w", err)
	}

	e := domain.NewDatasetCreatedEvent(&dataset)
	if err1 := s.msgAdapter.SendDatasetCreatedEvent(&e); err1 != nil {
		logrus.Errorf("failed to send dataset created event, dataset id:%s", dataset.Id.Identity())

	}

	return dataset.Id.Identity(), nil
}

// Delete deletes a dataset.
func (s *datasetAppService) Delete(user primitive.Account, datasetId primitive.Identity) (action string, err error) {
	dataset, err := s.repoAdapter.FindById(datasetId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		} else {
			err = xerrors.Errorf("find dataset by id failed, err:%w", err)
		}

		return
	}

	action = fmt.Sprintf(
		"delete dataset of %s:%s/%s",
		datasetId.Identity(), dataset.Owner.Account(), dataset.Name.MSDName(),
	)

	notFound, err := commonapp.CanDeleteOrNotFound(user, &dataset, s.permission)
	if err != nil {
		err = xerrors.Errorf("can not delete dataset, err:%w", err)
		return
	}
	if notFound {
		err = allerror.NewNotFound(allerror.ErrorCodeDatasetNotFound, "not found",
			xerrors.Errorf("%s not found", datasetId.Identity()))

		return
	}

	if err = s.codeRepoApp.Delete(dataset.RepoIndex()); err != nil {
		err = xerrors.Errorf("failed to delete dataset code repo, err:%w", err)
		return
	}

	if err = s.repoAdapter.Delete(dataset.Id); err != nil {
		err = xerrors.Errorf("failed to delete dataset info, err:%w", err)
		return
	}

	e := domain.NewDatasetDeletedEvent(user, dataset)
	if err := s.msgAdapter.SendDatasetDeletedEvent(&e); err != nil {
		logrus.Errorf("failed to send dataset deleted event, dataset id:%s, error:%s", datasetId.Identity(), err)
	}

	return
}

// Update updates a dataset.
func (s *datasetAppService) Update(
	user primitive.Account, datasetId primitive.Identity, cmd *CmdToUpdateDataset,
) (action string, err error) {
	dataset, err := s.repoAdapter.FindById(datasetId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeDatasetNotFound, "not found", xerrors.Errorf("failed to find dataset by id, %w", err))
		} else {
			err = xerrors.Errorf("failed to find dataset by id, %w", err)
		}

		return
	}

	action = fmt.Sprintf(
		"update dataset of %s:%s/%s",
		datasetId.Identity(), dataset.Owner.Account(), dataset.Name.MSDName(),
	)

	notFound, err := commonapp.CanUpdateOrNotFound(user, &dataset, s.permission)
	if err != nil {
		err = xerrors.Errorf("failed to find dataset by id, %w", err)
		return
	}
	if notFound {
		err = allerror.NewNotFound(allerror.ErrorCodeDatasetNotFound, "not found",
			xerrors.Errorf("%s not found", datasetId.Identity()))

		return
	}

	if dataset.IsDisable() {
		err = allerror.NewResourceDisabled(allerror.ErrorCodeResourceDisabled, "resource was disabled, cant be modified.",
			xerrors.Errorf("cant change resource to public"))
		return
	}

	isPrivateToPublic := dataset.IsPrivate() && cmd.Visibility.IsPublic()

	b, err := s.codeRepoApp.Update(&dataset.CodeRepo, &cmd.CmdToUpdateRepo)
	if err != nil {
		err = xerrors.Errorf("failed to update code repo, %w", err)
		return
	}

	b1 := cmd.toDataset(&dataset)
	if !b && !b1 {
		return
	}

	if err = s.repoAdapter.Save(&dataset); err != nil {
		err = xerrors.Errorf("failed to save dataset info, %w", err)
		return
	}

	e := domain.NewDatasetUpdatedEvent(&dataset, user, isPrivateToPublic)
	if err1 := s.msgAdapter.SendDatasetUpdatedEvent(&e); err1 != nil {
		logrus.Errorf("failed to send dataset updated event, dataset id:%s", datasetId.Identity())
	}

	return
}

// Disable disable a dataset.
func (s *datasetAppService) Disable(
	user primitive.Account, datasetId primitive.Identity, cmd *CmdToDisableDataset,
) (action string, err error) {
	dataset, err := s.repoAdapter.FindById(datasetId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeDatasetNotFound, "not found",
				xerrors.Errorf("failed to find dataset by id %d, %w", datasetId, err))
		} else {
			err = xerrors.Errorf("failed to find dataset by id %d, %w", datasetId, err)
		}

		return
	}

	action = fmt.Sprintf(
		"disable dataset of %s:%s/%s",
		datasetId.Identity(), dataset.Owner.Account(), dataset.Name.MSDName(),
	)

	err = s.canDisable(user)
	if err != nil {
		err = xerrors.Errorf("cant disable dataset:%d, %w", datasetId, err)
		return
	}

	if dataset.IsDisable() {
		err = allerror.NewResourceDisabled(allerror.ErrorCodeResourceAlreadyDisabled, "already been disabled",
			xerrors.Errorf("dataset %s already been disabled", dataset.Name.MSDName()))
		return
	}

	cmdRepo := coderepoapp.CmdToUpdateRepo{
		Visibility: primitive.VisibilityPrivate,
	}
	_, err = s.codeRepoApp.Update(&dataset.CodeRepo, &cmdRepo)
	if err != nil {
		err = xerrors.Errorf("failed to update dataset code repo:%d, %w", datasetId, err)
		return
	}

	cmd.toDataset(&dataset)

	if err = s.repoAdapter.Save(&dataset); err != nil {
		err = xerrors.Errorf("failed to save dataset:%d, %w", datasetId, err)
		return
	}

	return
}

func (s *datasetAppService) canDisable(user primitive.Account) error {
	if s.disableOrg != nil {
		if err := s.disableOrg.Contains(user); err != nil {
			logrus.Errorf("user:%s cant disable dataset err:%s", user.Account(), err)
			return allerror.NewNoPermission("no permission", xerrors.Errorf("cant disable"))
		}
	} else {
		logrus.Errorf("do not config disable org, no permit to disable")
		return allerror.NewNoPermission("no permission", xerrors.Errorf("cant disable"))
	}

	return nil
}

// GetByName retrieves a dataset by its name.
func (s *datasetAppService) GetByName(user primitive.Account, index *domain.DatasetIndex) (DatasetDTO, error) {
	var dto DatasetDTO

	dataset, err := s.repoAdapter.FindByName(index)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeDatasetNotFound, "not found",
				xerrors.Errorf("failed to find dataset by name:%s, %w", index.Name.MSDName(), err))
		} else {
			err = xerrors.Errorf("failed to find dataset by name:%s, %w", index.Name.MSDName(), err)
		}

		return dto, err
	}

	if err := s.permission.CanRead(user, &dataset); err != nil {
		if allerror.IsNoPermission(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeDatasetNotFound, "not found",
				xerrors.Errorf("not have permission to get dataset:%s, %w", index.Name.MSDName(), err))
		}

		return dto, err
	}

	return toDatasetDTO(&dataset), nil
}

// List retrieves a list of datasets.
func (s *datasetAppService) List(user primitive.Account, cmd *CmdToListDatasets) (
	DatasetsDTO, error,
) {
	if user == nil {
		cmd.Visibility = primitive.VisibilityPublic
	} else {
		if cmd.Owner == nil {
			// It can list the private datasets of user,
			// but it maybe no need to do it.
			cmd.Visibility = primitive.VisibilityPublic
		} else {
			if user != cmd.Owner {
				err := s.permission.CanListOrgResource(
					user, cmd.Owner, primitive.ObjTypeDataset,
				)
				if err != nil {
					cmd.Visibility = primitive.VisibilityPublic
				}
			}
		}
	}

	v, total, err := s.repoAdapter.List(cmd, user, s.member)

	return DatasetsDTO{
		Total:    total,
		Datasets: v,
	}, err
}

func (s *datasetAppService) datasetsCountCheck(owner primitive.Account) error {
	cmdToList := CmdToListDatasets{
		Owner: owner,
	}

	total, err := s.repoAdapter.Count(&cmdToList)
	if err != nil {
		return xerrors.Errorf("get datasets count error: %w", err)
	}

	if s.user.IsOrganization(owner) && total >= config.MaxCountPerOrg {
		return allerror.NewCountExceeded("dataset count exceed",
			xerrors.Errorf("dataset count(now:%d max:%d) exceed", total, config.MaxCountPerOrg))
	}

	if !s.user.IsOrganization(owner) && total >= config.MaxCountPerUser {
		return allerror.NewCountExceeded("dataset count exceed",
			xerrors.Errorf("dataset count(now:%d max:%d) exceed", total, config.MaxCountPerUser))
	}

	return nil
}

func (s *datasetAppService) AddLike(datasetId primitive.Identity) error {
	// Retrieve the code repository information.
	dataset, err := s.repoAdapter.FindById(datasetId)
	if err != nil {
		return xerrors.Errorf("failed to find dataset by id:%d, %w", datasetId, err)
	}

	// Only proceed if the repository is public.
	isPublic := dataset.IsPublic()

	if !isPublic {
		return nil
	}

	if err := s.repoAdapter.AddLike(dataset); err != nil {
		return xerrors.Errorf("failed to add dataset(%d) like:, %w", datasetId, err)
	}
	return nil
}

func (s *datasetAppService) DeleteLike(datasetId primitive.Identity) error {
	// Retrieve the code repository information.
	dataset, err := s.repoAdapter.FindById(datasetId)
	if err != nil {
		return xerrors.Errorf("failed to find dataset by id:%d, %w", datasetId, err)
	}

	// Only proceed if the repository is public.
	isPublic := dataset.IsPublic()
	if !isPublic {
		return nil
	}

	if err := s.repoAdapter.DeleteLike(dataset); err != nil {
		return xerrors.Errorf("failed to delete dataset(%d) like:, %w", datasetId, err)
	}
	return nil
}
