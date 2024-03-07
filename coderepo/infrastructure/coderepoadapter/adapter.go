/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package coderepoadapter provides an adapter for interacting with a code repository service.
package coderepoadapter

import (
	"fmt"

	"github.com/openmerlin/go-sdk/gitea"

	"github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	giteaclient "github.com/openmerlin/merlin-server/common/infrastructure/gitea"
)

type UserInfoAdapter interface {
	GetPlatformUserInfo(primitive.Account) (string, error)
}

type codeRepoAdapter struct {
	client      *gitea.Client
	userAdapter UserInfoAdapter
}

// NewRepoAdapter creates a new instance of the codeRepoAdapter.
func NewRepoAdapter(c *gitea.Client, userAdapter UserInfoAdapter) *codeRepoAdapter {
	return &codeRepoAdapter{client: c, userAdapter: userAdapter}
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

	p, err := adapter.userAdapter.GetPlatformUserInfo(repo.CreatedBy)
	if err != nil {
		return fmt.Errorf("failed to get platform user info: %w", err)
	}

	client, err := giteaclient.NewClient(repo.CreatedBy.Account(), p)
	if err != nil {
		return fmt.Errorf("failed to create git client: %w", err)
	}

	var obj *gitea.Repository
	if repo.OwnedByPerson() {
		obj, _, err = client.CreateRepo(opt)
	} else {
		obj, _, err = client.CreateOrgRepo(repo.Owner.Account(), opt)
	}
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
