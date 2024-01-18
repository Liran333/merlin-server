package repository

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	org "github.com/openmerlin/merlin-server/organization/domain"
	"github.com/openmerlin/merlin-server/user/domain"
)

type ListOption struct {
	// can't define Name as domain.ResourceName
	// because the Name can be subpart of the real resource name
	Name string

	// list the Owner only used when type is organization
	Owner primitive.Account

	// list by type
	Type *domain.UserType

	// sort
	SortType primitive.SortType

	// whether to calculate the total
	Count        bool
	PageNum      int
	CountPerPage int
}

type User interface {
	// user
	AddUser(*domain.User) (domain.User, error)
	SaveUser(*domain.User) (domain.User, error)
	DeleteUser(*domain.User) error
	GetByAccount(domain.Account) (domain.User, error)
	GetUserAvatarId(domain.Account) (primitive.AvatarId, error)
	GetUsersAvatarId([]string) ([]domain.User, error)
	GetUserFullname(domain.Account) (string, error)
	// org
	AddOrg(*org.Organization) (org.Organization, error)
	SaveOrg(*org.Organization) (org.Organization, error)
	DeleteOrg(*org.Organization) error
	CheckName(primitive.Account) bool
	GetOrgByName(primitive.Account) (org.Organization, error)
	GetOrgByOwner(primitive.Account) ([]org.Organization, error)
	// list
	ListAccount(*ListOption) ([]domain.User, int, error)
}
