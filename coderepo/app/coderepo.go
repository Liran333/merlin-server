/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/coderepo/domain/repoadapter"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
)

// CodeRepoAppService is an interface for code repository application service.
type CodeRepoAppService interface {
	Create(primitive.Account, *CmdToCreateRepo) (domain.CodeRepo, error)
	Delete(domain.CodeRepoIndex) error
	Update(*domain.CodeRepo, *CmdToUpdateRepo) (bool, error)
	IsRepoExist(*domain.CodeRepoIndex) (bool, error)
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
		return repo, err
	}

	return repo, nil
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

func (s *codeRepoAppService) IsRepoExist(index *domain.CodeRepoIndex) (bool, error) {
	_, err := s.repoAdapter.Get(index)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		}

		return false, err
	}

	return true, nil
}
