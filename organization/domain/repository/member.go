package repository

import (
	"github.com/openmerlin/merlin-server/organization/domain"
)

type OrgMember interface {
	Save(*domain.OrgMember) (domain.OrgMember, error)
	Delete(*domain.OrgMember) error
	DeleteByOrg(string) error
	GetByOrg(string) ([]domain.OrgMember, error)
	GetByOrgAndRole(string, domain.OrgRole) ([]domain.OrgMember, error)
	GetByOrgAndUser(org, user string) (domain.OrgMember, error)
	GetByUser(string) ([]domain.OrgMember, error)
	//Search(*UserSearchOption) (UserSearchResult, error)
}
