/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package coderepoadapter provides an adapter for interacting with a code repository service.
package coderepoadapter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/openmerlin/go-sdk/gitea"

	"github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	giteaclient "github.com/openmerlin/merlin-server/common/infrastructure/gitea"
)

const (
	repoAlreadyExistsErr = "The repository with the same name already exists."
)

// UserInfoAdapter is an interface that defines the methods for retrieving platform-specific user information.
type UserInfoAdapter interface {
	GetPlatformUserInfo(context.Context, primitive.Account) (string, error)
}

// codeRepoAdapter is an implementation of the CodeRepoAdapter interface.
type codeRepoAdapter struct {
	client      *gitea.Client
	config      Config
	userAdapter UserInfoAdapter
}

// NewRepoAdapter creates a new instance of the codeRepoAdapter.
func NewRepoAdapter(c *gitea.Client, userAdapter UserInfoAdapter, cfg *Config) *codeRepoAdapter {
	return &codeRepoAdapter{client: c, userAdapter: userAdapter, config: *cfg}
}

// Add adds a new code repository to the code repository service.
func (adapter *codeRepoAdapter) Add(ctx context.Context, repo *domain.CodeRepo, initReadme bool) error {
	defaultRef := primitive.InitCodeFileRef().FileRef()

	opt := gitea.CreateRepoOption{
		Name:          repo.Name.MSDName(),
		License:       repo.License.License()[0],
		Private:       repo.Visibility.IsPrivate() || adapter.config.ForceToBePrivate,
		DefaultBranch: defaultRef,
	}

	if initReadme {
		opt.Readme = "Default"
		opt.AutoInit = true
	}

	p, err := adapter.userAdapter.GetPlatformUserInfo(ctx, repo.CreatedBy)
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
	} else if err != nil && err.Error() == repoAlreadyExistsErr {
		err = commonrepo.NewErrorDuplicateCreating(err)
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

	private := repo.IsPrivate() || adapter.config.ForceToBePrivate
	opt.Private = &private

	_, _, err := adapter.client.EditRepo(
		index.Owner.Account(), index.Name.MSDName(), opt,
	)

	return err
}

// FindByIndex finds a codeRepo by its index.
func (adapter *codeRepoAdapter) FindByIndex(index primitive.Identity) (domain.CodeRepo, error) {
	repoID := index.Integer()
	repo, _, err := adapter.client.GetRepoByID(repoID)
	if err != nil {
		return domain.CodeRepo{}, err
	}
	visibility := "public"
	if repo.Private {
		visibility = "private"
	}
	return domain.CodeRepo{
		Id:    primitive.CreateIdentity(repo.ID),
		Name:  primitive.CreateMSDName(repo.Name),
		Owner: primitive.CreateAccount(repo.Owner.UserName),
		// todo fix license
		License:    primitive.CreateLicense([]string{"unkown"}),
		Visibility: primitive.CreateVisibility(visibility),
		CreatedBy:  primitive.CreateAccount(repo.Owner.UserName),
	}, err
}

// IsNotFound check whether a code repository is not found
func (adapter *codeRepoAdapter) IsNotFound(index primitive.Identity) bool {
	_, resp, _ := adapter.client.GetRepoByID(index.Integer())
	return resp.StatusCode == http.StatusNotFound
}
