package domain

import (
	coderepoprimitive "github.com/openmerlin/merlin-server/coderepo/domain/primitive"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

type Branch struct {
	BranchIndex

	Id         primitive.Identity
	RepoType   coderepoprimitive.RepoType
	CreatedAt  int64
	BaseBranch coderepoprimitive.BranchName
}

type BranchIndex struct {
	Repo   primitive.MSDName
	Owner  primitive.Account
	Branch coderepoprimitive.BranchName
}
