package branchclientadapter

import (
	"github.com/openmerlin/go-sdk/gitea"

	"github.com/openmerlin/merlin-server/coderepo/domain"
)

const (
	errorBranchInactive     = "branch_inactiver"
	errorBranchAlreadyExist = "branch_already_exist"
	errorBaseBranchNotFound = "base_branch_not_found"

	statusCodeInactive           = 403
	statusCodeBranchAlreadyExist = 409
	statusCodeBaseBranchNotFound = 404
)

type branchClientAdapter struct {
	client *gitea.Client
}

func NewBranchClientAdapter(c *gitea.Client) *branchClientAdapter {
	return &branchClientAdapter{client: c}
}

func (adapter *branchClientAdapter) CreateBranch(branch *domain.Branch) (n string, code string, err error) {
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

	code = generateCreateCode(r.StatusCode)

	return
}

func (adapter *branchClientAdapter) DeleteBranch(branch *domain.BranchIndex) error {
	_, _, err := adapter.client.DeleteRepoBranch(
		branch.Owner.Account(), branch.Repo.MSDName(), branch.Branch.BranchName(),
	)

	return err
}

func generateCreateCode(c int) (code string) {
	switch c {
	case statusCodeBaseBranchNotFound:
		code = errorBaseBranchNotFound
	case statusCodeBranchAlreadyExist:
		code = errorBranchAlreadyExist
	case statusCodeInactive:
		code = errorBranchInactive
	}

	return
}
