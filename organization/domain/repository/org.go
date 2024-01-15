package repository

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/organization/domain"
)

type Organization interface {
	AddOrg(*domain.Organization) (domain.Organization, error)
	SaveOrg(*domain.Organization) (domain.Organization, error)
	DeleteOrg(*domain.Organization) error
	CheckName(primitive.Account) bool
	GetOrgByName(primitive.Account) (domain.Organization, error)
	GetOrgByOwner(primitive.Account) ([]domain.Organization, error)
	//Search(*UserSearchOption) (UserSearchResult, error)
}
