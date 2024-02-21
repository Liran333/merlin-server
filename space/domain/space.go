package domain

import (
	"k8s.io/apimachinery/pkg/util/sets"

	coderepo "github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	spaceprimitive "github.com/openmerlin/merlin-server/space/domain/primitive"
)

type Space struct {
	coderepo.CodeRepo

	SDK       spaceprimitive.SDK
	Desc      primitive.MSDDesc
	Labels    SpaceLabels
	Fullname  primitive.MSDFullname
	Hardware  spaceprimitive.Hardware
	CreatedBy primitive.Account

	Version       int
	CreatedAt     int64
	UpdatedAt     int64
	LikeCount     int
	DownloadCount int
}

func (s *Space) ResourceOwner() primitive.Account {
	return s.Owner
}

func (s *Space) ResourceType() primitive.ObjType {
	return primitive.ObjTypeSpace
}

func (s *Space) IsCreatedBy(user primitive.Account) bool {
	return s.CreatedBy == user
}

func (s *Space) OwnedByPerson() bool {
	return s.Owner == s.CreatedBy
}

type SpaceLabels struct {
	Task       string           // task label
	Others     sets.Set[string] // other labels
	Frameworks sets.Set[string] // framework labels
}

type SpaceIndex = coderepo.CodeRepoIndex
