package repository

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/organization/domain"
)

type Organization interface {
	Save(*domain.Organization) (domain.Organization, error)
	Delete(*domain.Organization) error
	CheckName(primitive.Account) bool
	GetByName(primitive.Account) (domain.Organization, error)
	GetByOwner(primitive.Account) ([]domain.Organization, error)
	//Search(*UserSearchOption) (UserSearchResult, error)
}
