package repoadapter

import (
	"github.com/openmerlin/merlin-server/coderepo/domain"
)

type RepoAdapter interface {
	Add(*domain.CodeRepo, bool) error
	// TODO delete by id
	Delete(*domain.CodeRepoIndex) error
	// TODO change domain.CodeRepoIndex to id
	Save(*domain.CodeRepoIndex, *domain.CodeRepo) error
}
