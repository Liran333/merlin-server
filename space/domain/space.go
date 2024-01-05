package domain

import (
	"k8s.io/apimachinery/pkg/util/sets"

	coderepo "github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

type Space struct {
	coderepo.CodeRepo

	Desc      primitive.MSDDesc
	Labels    SpaceLabels
	Fullname  primitive.MSDFullname
	CreatedBy primitive.Account

	Version       int
	CreatedAt     int64
	UpdatedAt     int64
	LikeCount     int
	DownloadCount int
}

func (s *Space) OwnedBy(user primitive.Account) bool {
	return s.Owner == user || s.CreatedBy == user
}

func (s *Space) OwnedByPerson() bool {
	return s.Owner == s.CreatedBy
}

type SpaceLabels struct {
	Task       string           // task label
	Others     sets.Set[string] // other labels
	Frameworks sets.Set[string] // framework labels
}

type SpaceIndex struct {
	Owner primitive.Account
	Name  primitive.MSDName
}
