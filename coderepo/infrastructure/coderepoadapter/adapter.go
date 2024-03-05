/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package coderepoadapter provides an adapter for interacting with a code repository service.
package coderepoadapter

import (
	"github.com/openmerlin/go-sdk/gitea"

	"github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
)

const msgRepoNotExist = "The target couldn't be found."

type codeRepoAdapter struct {
	client *gitea.Client
}

// NewRepoAdapter creates a new instance of the codeRepoAdapter.
func NewRepoAdapter(c *gitea.Client) *codeRepoAdapter {
	return &codeRepoAdapter{client: c}
}

// Add adds a new code repository to the code repository service.
func (adapter *codeRepoAdapter) Add(repo *domain.CodeRepo, initReadme bool) error {
	defaultRef := primitive.InitCodeFileRef().FileRef()

	opt := gitea.CreateRepoOption{
		Name:          repo.Name.MSDName(),
		License:       repo.License.License(),
		Private:       repo.Visibility.IsPrivate(),
		DefaultBranch: defaultRef,
	}

	if initReadme {
		opt.Readme = "Default"
		opt.AutoInit = true
	}

	obj, _, err := adapter.client.AdminCreateRepo(repo.Owner.Account(), opt)
	if err == nil {
		repo.Id = primitive.CreateIdentity(obj.ID)
	}

	return err
}

// Delete deletes a code repository from the code repository service.
func (adapter *codeRepoAdapter) Delete(index *domain.CodeRepoIndex) error {
	_, err := adapter.client.DeleteRepo(index.Owner.Account(), index.Name.MSDName())

	return err
}

// Save saves updates to a code repository in the code repository service.
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

func (adapter *codeRepoAdapter) Get(index *domain.CodeRepoIndex) (*domain.Repository, error) {
	repo, _, err := adapter.client.GetRepo(index.Owner.Account(), index.Name.MSDName())
	if err != nil && err.Error() == msgRepoNotExist {
		err = commonrepo.NewErrorResourceNotExists(err)
	}

	return repo, err
}
