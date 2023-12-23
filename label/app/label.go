package app

import (
	"strings"

	"code.gitea.io/gitea/modules/structs"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openmerlin/merlin-server/label/domain/gitea"
	"github.com/openmerlin/merlin-server/label/utils"
	"github.com/openmerlin/merlin-server/models/domain"
)

const (
	readMeFileName         = "README.md"
	readMeFileResideBranch = "main"
)

type LabelService interface {
	GetLabels(*structs.PushPayload) (domain.ModelLabels, error)
}

type labelService struct {
	gitea gitea.Gitea
}

func NewLabelHandler(g gitea.Gitea) *labelService {
	return &labelService{
		gitea: g,
	}
}

func (l *labelService) GetLabels(p *structs.PushPayload) (labels domain.ModelLabels, err error) {
	if strings.Split(p.Ref, "/")[2] != readMeFileResideBranch {
		return
	}

	modifiedFilesSet := sets.NewString(p.HeadCommit.Modified...)
	addedFilesSet := sets.NewString(p.HeadCommit.Added...)
	if !modifiedFilesSet.Has(readMeFileName) && !addedFilesSet.Has(readMeFileName) {
		return
	}

	org, repo := utils.GetOrgRepo(p.Repo)
	opt := gitea.Option{
		Org:  org,
		Repo: repo,
		Ref:  readMeFileResideBranch,
		Path: readMeFileName,
	}

	return l.gitea.GetLabels(&opt)
}
