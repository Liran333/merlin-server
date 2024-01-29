package app

import (
	"github.com/openmerlin/merlin-server/coderepo/domain"
	repoprimitive "github.com/openmerlin/merlin-server/coderepo/domain/primitive"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	spaceprimitive "github.com/openmerlin/merlin-server/space/domain/primitive"
	"github.com/openmerlin/merlin-server/utils"
)

type CmdToCreateRepo struct {
	SDK        spaceprimitive.SDK
	Name       primitive.MSDName
	Owner      primitive.Account
	License    primitive.License
	Hardware   spaceprimitive.Hardware
	Visibility primitive.Visibility
	InitReadme bool
}

func (cmd *CmdToCreateRepo) toCodeRepo() domain.CodeRepo {
	return domain.CodeRepo{
		Name:       cmd.Name,
		Owner:      cmd.Owner,
		License:    cmd.License,
		Visibility: cmd.Visibility,
	}
}

type CmdToUpdateRepo struct {
	Name       primitive.MSDName
	Visibility primitive.Visibility
}

func (cmd *CmdToUpdateRepo) toRepo(repo *domain.CodeRepo) (b bool) {
	if v := cmd.Name; v != nil && v != repo.Name {
		repo.Name = v
		b = true
	}

	if v := cmd.Visibility; v != nil && v != repo.Visibility {
		repo.Visibility = v
		b = true
	}

	return
}

type CmdToCreateBranch struct {
	domain.BranchIndex

	RepoType   repoprimitive.RepoType
	BaseBranch repoprimitive.BranchName
}

func (cmd *CmdToCreateBranch) toBranch() domain.Branch {
	branch := domain.Branch{
		BranchIndex: cmd.BranchIndex,
		BaseBranch:  cmd.BaseBranch,
		RepoType:    cmd.RepoType,
		CreatedAt:   utils.Now(),
	}

	return branch
}

type BranchCreateDTO struct {
	Name string `json:"branch_name"`
}

func toBranchCreateDTO(v string) BranchCreateDTO {
	return BranchCreateDTO{
		Name: v,
	}
}

type CmdToDeleteBranch struct {
	domain.BranchIndex

	RepoType repoprimitive.RepoType
}
