package branchclientadapter

import (
	"github.com/openmerlin/go-sdk/gitea"

	"github.com/openmerlin/merlin-server/coderepo/domain"
)

type branchClientAdapter struct {
	client *gitea.Client
}

func NewBranchClientAdapter(c *gitea.Client) *branchClientAdapter {
	return &branchClientAdapter{client: c}
}

func (adapter *branchClientAdapter) CreateBranch(branch *domain.Branch) (n string, err error) {
	opt := gitea.CreateBranchOption{}
	opt.BranchName = branch.Branch.BranchName()
	opt.OldBranchName = branch.BaseBranch.BranchName()

	b, _, err := adapter.client.CreateBranch(
		branch.Owner.Account(), branch.Repo.MSDName(), opt,
	)
	if err == nil {
		n = b.Name
	}

	return
}

func (adapter *branchClientAdapter) DeleteBranch(branch *domain.BranchIndex) error {
	_, _, err := adapter.client.DeleteRepoBranch(
		branch.Owner.Account(), branch.Repo.MSDName(), branch.Branch.BranchName(),
	)

	return err
}
