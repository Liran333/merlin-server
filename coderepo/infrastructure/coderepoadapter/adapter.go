package coderepoadapter

import (
	"code.gitea.io/sdk/gitea"

	"github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

type codeRepoAdapter struct {
	client *gitea.Client
}

func (adapter *codeRepoAdapter) Add(repo *domain.CodeRepo) error {
	opt := gitea.CreateRepoOption{}
	opt.Name = repo.Name.MSDName()
	obj, _, err := adapter.client.CreateRepo(opt)
	if err != nil {
		return err
	}

	repo.Id = primitive.CreateIdentity(obj.ID)

	return nil
}

func (adapter *codeRepoAdapter) Delete(primitive.Identity) error {
	return nil
}

func (adapter *codeRepoAdapter) Save(repo *domain.CodeRepo) error {
	return nil
}
