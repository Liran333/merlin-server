package domain

import (
	"k8s.io/apimachinery/pkg/util/sets"

	coderepo "github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

type Model struct {
	coderepo.CodeRepo

	Desc      primitive.MSDDesc
	Labels    ModelLabels
	Fullname  primitive.MSDFullname
	CreatedBy primitive.Account

	Version       int
	CreatedAt     int64
	UpdatedAt     int64
	LikeCount     int
	DownloadCount int
}

func (m *Model) ResourceOwner() primitive.Account {
	return m.Owner
}

func (m *Model) ResourceType() primitive.ObjType {
	return primitive.ObjTypeModel
}

func (m *Model) IsCreatedBy(user primitive.Account) bool {
	return m.CreatedBy == user
}

func (m *Model) OwnedByPerson() bool {
	return m.Owner == m.CreatedBy
}

type ModelLabels struct {
	Task       string           // task label
	Others     sets.Set[string] // other labels
	Frameworks sets.Set[string] // framework labels
}

type ModelIndex = coderepo.CodeRepoIndex
