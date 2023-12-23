package gitea

import (
	"github.com/openmerlin/merlin-server/models/domain"
)

type Option struct {
	Org  string
	Repo string
	Ref  string
	Path string
}

type Gitea interface {
	GetLabels(option *Option) (domain.ModelLabels, error)
}
