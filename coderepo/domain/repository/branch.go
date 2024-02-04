package repository

import (
	"github.com/openmerlin/merlin-server/coderepo/domain"
	coderepoprimtive "github.com/openmerlin/merlin-server/coderepo/domain/primitive"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

type BranchRepositoryAdapter interface {
	Add(*domain.Branch) error
	Delete(primitive.Identity) error
	FindByIndex(*domain.BranchIndex) (domain.Branch, error)
}

type BranchClientAdapter interface {
	CreateBranch(*domain.Branch) (string, error)
	DeleteBranch(*domain.BranchIndex) error
}

type CheckRepoAdapter interface {
	CheckRepo(coderepoprimtive.RepoType, primitive.Account, primitive.MSDName) error
}
