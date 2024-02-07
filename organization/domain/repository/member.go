package repository

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/organization/domain"
)

type OrgMember interface {
	Add(*domain.OrgMember) (domain.OrgMember, error)
	Save(*domain.OrgMember) (domain.OrgMember, error)
	Delete(*domain.OrgMember) error
	DeleteByOrg(primitive.Account) error
	GetByOrg(string) ([]domain.OrgMember, error)
	GetByOrgAndRole(string, domain.OrgRole) ([]domain.OrgMember, error)
	GetByOrgAndUser(org, user string) (domain.OrgMember, error)
	GetByUser(string) ([]domain.OrgMember, error)
}
