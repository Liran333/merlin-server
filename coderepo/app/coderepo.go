/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"github.com/openmerlin/git-access-sdk/filescan"
	"github.com/openmerlin/git-access-sdk/filescan/api"
	"github.com/sirupsen/logrus"

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
			err = allerror.NewInvalidParam(err.Error())
		}

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

	// when the Visibility of repo change from private to public, must scan it
	whetherToScan := repo.IsPrivate() && cmd.Visibility.IsPublic()

	if !cmd.toRepo(repo) {
		return false, nil
	}

	err := s.repoAdapter.Save(&index, repo)

	if err == nil && whetherToScan {
		_, err1 := api.CreateFileInfo(&filescan.ReqToCreateFileInfo{
			Owner:  repo.Owner.Account(),
			Repo:   repo.Name.MSDName(),
			RepoId: repo.Id.Integer(),
		})
		if err1 != nil {
			logrus.Errorf("create file scan task of %s/%s error:%s",
				repo.Owner.Account(), repo.Name.MSDName(), err1.Error())
		}
	}

	return true, err
}
