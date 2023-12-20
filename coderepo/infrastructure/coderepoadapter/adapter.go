package coderepoadapter

import (
	"code.gitea.io/sdk/gitea"
	"github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

var _codeRepoAdapter codeRepoAdapter

type codeRepoAdapter struct {
	client *gitea.Client
}

func Init(c *gitea.Client) {
	_codeRepoAdapter.client = c
	return
}

func (adapter *codeRepoAdapter) Get(repoIdentity primitive.Identity) (*gitea.Repository, error) {
	repoId := repoIdentity.Integer()
	repo, _, err := _codeRepoAdapter.client.GetRepoByID(repoId)
	return repo, err
}

func (adapter *codeRepoAdapter) Add(repo *domain.CodeRepo) error {
	opt := gitea.CreateRepoOption{}
	opt.Name = repo.Name.MSDName()
	owner := repo.Owner.Account()
	opt.License = repo.License.License()
	opt.Private = repo.Visibility.IsPrivate()
	obj, _, err := _codeRepoAdapter.client.AdminCreateRepo(owner, opt)
	if err != nil {
		return err
	}
	id := obj.ID
	repo.Id = primitive.CreateIdentity(id)
	return nil
}

func (adapter *codeRepoAdapter) Delete(repoIdentity primitive.Identity) error {
	repo, err := adapter.Get(repoIdentity)
	if err != nil {
		return err
	}
	ownerName := repo.Owner.UserName
	repoName := repo.Name
	_, err = _codeRepoAdapter.client.DeleteRepo(ownerName, repoName)
	if err != nil {
		return err
	}
	return nil
}

func (adapter *codeRepoAdapter) Save(repoIdentity primitive.Identity, repo *domain.CodeRepo) error {
	rep, err := adapter.Get(repoIdentity)
	if err != nil {
		return err
	}
	ownerName := rep.Owner.UserName
	repoName := rep.Name
	opt := new(gitea.EditRepoOption)
	MSDName := repo.Name.MSDName()
	opt.Name = &MSDName
	private := repo.IsPrivate()
	opt.Private = &private
	_, _, err = _codeRepoAdapter.client.EditRepo(ownerName, repoName, *opt)
	if err != nil {
		return err
	}
	return nil
}
