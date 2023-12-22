package domain

import (
	"k8s.io/apimachinery/pkg/util/sets"

	coderepo "github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

type Model struct {
	coderepo.CodeRepo

	Desc     primitive.MSDDesc
	Labels   ModelLabels
	Fullname primitive.MSDFullname

	Version       int
	CreatedAt     int64
	UpdatedAt     int64
	LikeCount     int
	DownloadCount int
}

type ModelLabels struct {
	Task       string           // task label
	Others     sets.Set[string] // other labels
	Frameworks sets.Set[string] // framework labels
}
