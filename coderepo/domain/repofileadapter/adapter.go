package repofileadapter

import (
	"github.com/openmerlin/merlin-server/coderepo/domain"
)

type CodeRepoFileAdapter interface {
	List(*domain.CodeRepoFile) (*domain.ListFileInfo, error)

	Get(*domain.CodeRepoFile) (*domain.DetailFileInfo, error)

	Download(*domain.CodeRepoFile) (*domain.DownLoadFileInfo, error)
}
