package resourceadapter

import (
	"github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/coderepo/domain/primitive"
)

type ResourceAdapter interface {
	GetByName(*domain.CodeRepoIndex) (domain.Resource, error)
	GetByType(primitive.RepoType, *domain.CodeRepoIndex) (domain.Resource, error)
}
