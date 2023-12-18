package platform

import (
	org "github.com/openmerlin/merlin-server/organization/domain"
	"github.com/openmerlin/merlin-server/user/domain"
)

type BaseAuthClient interface {
	CreateToken(*domain.TokenCreatedCmd) (domain.PlatformToken, error)
	DeleteToken(*domain.TokenDeletedCmd) error
	CreateOrg(*org.Organization) error
	DeleteOrg(string) error
	CanDelete(string) (bool, error)
	AddMember(*org.Organization, *org.OrgMember) error
	RemoveMember(*org.Organization, *org.OrgMember) error
	EditMemberRole(*org.Organization, org.OrgRole, *org.OrgMember) error
}
