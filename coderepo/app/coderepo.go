package app

import (
	"github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

type CodeRepoAppService interface {
	Create(cmd *CmdToCreateRepo) (domain.CodeRepo, error)
	Delete(primitive.Identity) error
	Update(*domain.CodeRepo, *CmdToUpdateRepo) (bool, error)
}

func NewCodeRepoAppService() *codeRepoAppService {
	return &codeRepoAppService{}
}

type codeRepoAppService struct{}

func (impl *codeRepoAppService) Create(cmd *CmdToCreateRepo) (domain.CodeRepo, error) {
	repo := cmd.toCodeRepo()

	// TODO should check if duplicate creating

	return repo, nil
}

func (impl *codeRepoAppService) Delete(primitive.Identity) error {
	// TODO if the repo does not exist, ignore it

	return nil
}

func (impl *codeRepoAppService) Update(*domain.CodeRepo, *CmdToUpdateRepo) (bool, error) {
	return false, nil
}
