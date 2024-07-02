/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package branchclientadapter provides an adapter implementation for interacting with the Gitea branch client.
package branchclientadapter

import (
	"errors"

	"github.com/openmerlin/go-sdk/gitea"

	"github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
)

const (
	statusCodeInactive           = 403
	statusCodeBranchAlreadyExist = 409
	statusCodeBaseBranchNotFound = 404
)

type branchClientAdapter struct {
	client *gitea.Client
}

// NewBranchClientAdapter creates a new instance of branchClientAdapter with the given gitea.Client.
func NewBranchClientAdapter(c *gitea.Client) *branchClientAdapter {
	return &branchClientAdapter{client: c}
}

// CreateBranch creates a new branch in the repository using the provided branch information.
func (adapter *branchClientAdapter) CreateBranch(branch *domain.Branch) (n string, err error) {
	opt := gitea.CreateBranchOption{}
	opt.BranchName = branch.Branch.BranchName()
	opt.OldBranchName = branch.BaseBranch.BranchName()

	b, r, err := adapter.client.CreateBranch(
		branch.Owner.Account(), branch.Repo.MSDName(), opt,
	)
	if err == nil {
		n = b.Name
		return
	}

	err = parseCreateError(r.StatusCode, err)
	return
}

// DeleteBranch deletes the specified branch from the repository.
func (adapter *branchClientAdapter) DeleteBranch(branch *domain.BranchIndex) error {
	b, _, err := adapter.client.DeleteRepoBranch(
		branch.Owner.Account(), branch.Repo.MSDName(), branch.Branch.BranchName(),
	)
	if !b {
		err = errors.New("delete failed, branch not exist")
	}

	return err
}

func parseCreateError(c int, err error) error {
	switch c {
	case statusCodeBaseBranchNotFound:
		return allerror.New(allerror.ErrorCodeBaseBranchNotFound, "base branch not found", err)
	case statusCodeBranchAlreadyExist:
		return allerror.New(allerror.ErrorCodeBranchExist, "branch already exist", err)
	case statusCodeInactive:
		return allerror.New(allerror.ErrorCodeBranchInavtive, "branch inactive", err)
	default:
		// default case modified to return 500
		return allerror.New(allerror.ErrorBaseCase, "unexpected error when creating branch", err)
	}
}
