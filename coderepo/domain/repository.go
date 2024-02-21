package domain

import (
	commondomain "github.com/openmerlin/merlin-server/common/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

type CodeRepo struct {
	Id         primitive.Identity
	Name       primitive.MSDName
	Owner      primitive.Account
	License    primitive.License
	Visibility primitive.Visibility
}

func (r *CodeRepo) IsPrivate() bool {
	return r.Visibility.IsPrivate()
}

func (r *CodeRepo) IsPublic() bool {
	return r.Visibility.IsPublic()
}

func (r *CodeRepo) RepoIndex() CodeRepoIndex {
	return CodeRepoIndex{
		Name:  r.Name,
		Owner: r.Owner,
	}
}

type CodeRepoIndex struct {
	Name  primitive.MSDName
	Owner primitive.Account
}

type Resource = commondomain.Resource
