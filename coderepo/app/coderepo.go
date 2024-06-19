/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package app provides application services for creating and managing branches.
package app

import (
	"context"

	"github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/coderepo/domain/repoadapter"
	commondomain "github.com/openmerlin/merlin-server/common/domain"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
)

// CodeRepoAppService is an interface for code repository application service.
type CodeRepoAppService interface {
	Create(context.Context, primitive.Account, *CmdToCreateRepo) (domain.CodeRepo, error)
	Delete(commondomain.CodeRepoIndex) error
	Update(*domain.CodeRepo, *CmdToUpdateRepo) (bool, error)
	GetById(primitive.Identity) (domain.CodeRepo, error)
	IsNotFound(primitive.Identity) bool
}

// NewCodeRepoAppService creates a new instance of CodeRepoAppService.
func NewCodeRepoAppService(repoAdapter repoadapter.RepoAdapter) *codeRepoAppService {
	return &codeRepoAppService{repoAdapter: repoAdapter}
}

type codeRepoAppService struct {
	repoAdapter repoadapter.RepoAdapter
}

// Create creates a new code repository.
func (s *codeRepoAppService) Create(
	ctx context.Context, user primitive.Account, cmd *CmdToCreateRepo) (domain.CodeRepo, error) {
	repo := cmd.toCodeRepo(user)

	if err := s.repoAdapter.Add(ctx, &repo, cmd.InitReadme); err != nil {
		if commonrepo.IsErrorDuplicateCreating(err) {
			err = allerror.New(allerror.ErrorDuplicateCreating, "dulicate creating", err)
		}

		return repo, err
	}

	return repo, nil
}

// Get a coderepo object by id.
func (s *codeRepoAppService) GetById(index primitive.Identity) (domain.CodeRepo, error) {
	repo, err := s.repoAdapter.FindByIndex(index)

	return repo, err
}

// Delete deletes a code repository by its index.
func (s *codeRepoAppService) Delete(index commondomain.CodeRepoIndex) error {

	return s.repoAdapter.Delete(&index)
}

// Update updates a code repository with the given command.
func (s *codeRepoAppService) Update(repo *domain.CodeRepo, cmd *CmdToUpdateRepo) (bool, error) {
	index := repo.RepoIndex()

	if !cmd.toRepo(repo) {
		return false, nil
	}

	return true, s.repoAdapter.Save(&index, repo)
}

// IsNotFound check whether a code repository is not found
func (s *codeRepoAppService) IsNotFound(index primitive.Identity) bool {
	return s.repoAdapter.IsNotFound(index)
}
