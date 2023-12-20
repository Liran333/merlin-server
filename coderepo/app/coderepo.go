package app

import (
	"github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/coderepo/domain/repoadapter"
)

type CodeRepoAppService interface {
	Create(cmd *CmdToCreateRepo) (domain.CodeRepo, error)
	//  TODO should define Delete(primitive.Identity) error
	Delete(domain.CodeRepoIndex) error
	Update(*domain.CodeRepo, *CmdToUpdateRepo) (bool, error)
}

func NewCodeRepoAppService(repoAdapter repoadapter.RepoAdapter) *codeRepoAppService {
	return &codeRepoAppService{repoAdapter: repoAdapter}
}

type codeRepoAppService struct {
	repoAdapter repoadapter.RepoAdapter
}

func (s *codeRepoAppService) Create(cmd *CmdToCreateRepo) (domain.CodeRepo, error) {
	repo := cmd.toCodeRepo()

	if err := s.repoAdapter.Add(&repo); err != nil {
		// TODO should check if duplicate creating and send allerror

		return repo, err
	}

	return repo, nil
}

func (s *codeRepoAppService) Delete(index domain.CodeRepoIndex) error {
	// TODO if the repo does not exist, ignore it

	return s.repoAdapter.Delete(&index)
}

func (s *codeRepoAppService) Update(repo *domain.CodeRepo, cmd *CmdToUpdateRepo) (bool, error) {
	index := repo.RepoIndex()

	if !cmd.toRepo(repo) {
		return false, nil
	}

	return true, s.repoAdapter.Save(&index, repo)
}
