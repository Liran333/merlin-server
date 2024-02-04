package branchclientadapter

import (
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

func NewBranchClientAdapter(c *gitea.Client) *branchClientAdapter {
	return &branchClientAdapter{client: c}
}

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

func (adapter *branchClientAdapter) DeleteBranch(branch *domain.BranchIndex) error {
	_, _, err := adapter.client.DeleteRepoBranch(
		branch.Owner.Account(), branch.Repo.MSDName(), branch.Branch.BranchName(),
	)

	return err
}

func parseCreateError(c int, err error) error {
	switch c {
	case statusCodeBaseBranchNotFound:
		return allerror.New(allerror.ErrorCodeBaseBranchNotFound, "base branch not found")
	case statusCodeBranchAlreadyExist:
		return allerror.New(allerror.ErrorCodeBranchExist, "branch already exist")
	case statusCodeInactive:
		return allerror.New(allerror.ErrorCodeBranchInavtive, "branch inactive")
	}

	return err
}
