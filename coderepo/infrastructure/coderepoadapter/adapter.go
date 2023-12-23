package coderepoadapter

import (
	"code.gitea.io/sdk/gitea"
	"github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

type codeRepoAdapter struct {
	client *gitea.Client
}

func NewRepoAdapter(c *gitea.Client) *codeRepoAdapter {
	return &codeRepoAdapter{client: c}
}

func (adapter *codeRepoAdapter) Add(repo *domain.CodeRepo, initReadme bool) error {
	readme := ""
	if initReadme {
		readme = repo.Name.MSDName()
	}

	defaultRef := primitive.InitCodeFileRef().FileRef()
	obj, _, err := adapter.client.AdminCreateRepo(
		repo.Owner.Account(),
		gitea.CreateRepoOption{
			Name:          repo.Name.MSDName(),
			Readme:        readme,
			License:       repo.License.License(),
			Private:       repo.Visibility.IsPrivate(),
			DefaultBranch: defaultRef,
		},
	)
	if err == nil {
		repo.Id = primitive.CreateIdentity(obj.ID)
	}

	// TODO check if duplicate create

	return err
}

func (adapter *codeRepoAdapter) Delete(index *domain.CodeRepoIndex) error {
	_, err := adapter.client.DeleteRepo(index.Owner.Account(), index.Name.MSDName())

	// TODO check if delete the unavailable repo

	return err
}

func (adapter *codeRepoAdapter) Save(index *domain.CodeRepoIndex, repo *domain.CodeRepo) error {
	opt := gitea.EditRepoOption{}

	name := repo.Name.MSDName()
	opt.Name = &name

	private := repo.IsPrivate()
	opt.Private = &private

	_, _, err := adapter.client.EditRepo(
		index.Owner.Account(), index.Name.MSDName(), opt,
	)

	return err
}
