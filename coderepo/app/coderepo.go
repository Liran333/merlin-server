/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/coderepo/domain/repoadapter"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
)

// CodeRepoAppService is an interface for code repository application service.
type CodeRepoAppService interface {
	Create(primitive.Account, *CmdToCreateRepo) (domain.CodeRepo, error)
	Delete(domain.CodeRepoIndex) error
	Update(*domain.CodeRepo, *CmdToUpdateRepo) (bool, error)
	GetById(primitive.Identity) (domain.CodeRepo, error)
}

// NewCodeRepoAppService creates a new instance of CodeRepoAppService.
func NewCodeRepoAppService(repoAdapter repoadapter.RepoAdapter) *codeRepoAppService {
	return &codeRepoAppService{repoAdapter: repoAdapter}
}

type codeRepoAppService struct {
	repoAdapter repoadapter.RepoAdapter
}

// Create creates a new code repository.
func (s *codeRepoAppService) Create(user primitive.Account, cmd *CmdToCreateRepo) (domain.CodeRepo, error) {
	repo := cmd.toCodeRepo(user)

	if err := s.repoAdapter.Add(&repo, cmd.InitReadme); err != nil {
		if commonrepo.IsErrorDuplicateCreating(err) {
			err = allerror.NewInvalidParam("dulicate creating", err)
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
func (s *codeRepoAppService) Delete(index domain.CodeRepoIndex) error {

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
