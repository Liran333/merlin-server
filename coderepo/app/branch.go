/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package app provides application services for creating and managing branches.
package app

import (
	"github.com/openmerlin/merlin-server/coderepo/domain"
	repoprimitive "github.com/openmerlin/merlin-server/coderepo/domain/primitive"
	"github.com/openmerlin/merlin-server/coderepo/domain/repository"
	"github.com/openmerlin/merlin-server/coderepo/domain/resourceadapter"
	commonapp "github.com/openmerlin/merlin-server/common/app"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
)

// BranchAppService defines the interface for a branch application service.
type BranchAppService interface {
	Create(primitive.Account, *CmdToCreateBranch) (BranchCreateDTO, error)
	Delete(primitive.Account, *CmdToDeleteBranch) error
}

// NewBranchAppService creates a new instance of the BranchAppService.
func NewBranchAppService(
	permission commonapp.ResourcePermissionAppService,
	branchAdapter repository.BranchRepositoryAdapter,
	resourceAdapter resourceadapter.ResourceAdapter,
	branchClientAdapter repository.BranchClientAdapter,
) BranchAppService {
	return &branchAppService{
		permission:          permission,
		branchAdapter:       branchAdapter,
		resourceAdapter:     resourceAdapter,
		branchClientAdapter: branchClientAdapter,
	}
}

type branchAppService struct {
	permission          commonapp.ResourcePermissionAppService
	branchAdapter       repository.BranchRepositoryAdapter
	resourceAdapter     resourceadapter.ResourceAdapter
	branchClientAdapter repository.BranchClientAdapter
}

// Create creates a new branch based on the provided command and returns the created branch DTO.
func (s *branchAppService) Create(user primitive.Account, cmd *CmdToCreateBranch) (
	dto BranchCreateDTO, err error,
) {
	index := cmd.RepoIndex()
	if err = s.canModify(user, cmd.RepoType, &index); err != nil {
		return
	}

	branch := cmd.toBranch()

	v, err := s.branchClientAdapter.CreateBranch(&branch)
	if err != nil {
		return
	}

	if err = s.branchAdapter.Add(&branch); err == nil {
		dto = toBranchCreateDTO(v)
	}

	return
}

// Delete deletes a branch based on the provided command.
func (s *branchAppService) Delete(user primitive.Account, cmd *CmdToDeleteBranch) error {
	index := cmd.RepoIndex()
	if err := s.canModify(user, cmd.RepoType, &index); err != nil {
		return err
	}

	br, err := s.branchAdapter.FindByIndex(&cmd.BranchIndex)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		}

		return err
	}

	if err = s.branchClientAdapter.DeleteBranch(&cmd.BranchIndex); err != nil {
		return err
	}

	return s.branchAdapter.Delete(br.Id)
}

func (s *branchAppService) canModify(
	user primitive.Account, t repoprimitive.RepoType, index *domain.CodeRepoIndex,
) error {
	repo, err := s.resourceAdapter.GetByType(t, index)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			return allerror.NewNotFound(allerror.ErrorCodeRepoNotFound, "no repo", err)
		}

		return err
	}

	return s.permission.CanUpdate(user, repo)
}
