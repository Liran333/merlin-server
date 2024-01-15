package platform

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	org "github.com/openmerlin/merlin-server/organization/domain"
	"github.com/openmerlin/merlin-server/user/domain"
)

type BaseAuthClient interface {
	CreateToken(*domain.TokenCreatedCmd) (domain.PlatformToken, error)
	DeleteToken(*domain.TokenDeletedCmd) error
	CreateOrg(*org.Organization) error
	DeleteOrg(primitive.Account) error
	CanDelete(primitive.Account) (bool, error)
	AddMember(*org.Organization, *org.OrgMember) error
	RemoveMember(*org.Organization, *org.OrgMember) error
	EditMemberRole(*org.Organization, org.OrgRole, *org.OrgMember) error
}
