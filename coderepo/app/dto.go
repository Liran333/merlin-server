package app

import (
	"github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

type CmdToCreateRepo struct {
	Name       primitive.MSDName
	Owner      primitive.Account
	License    primitive.License
	Visibility primitive.Visibility
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
