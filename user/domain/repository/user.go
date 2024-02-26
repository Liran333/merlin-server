/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package repository

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	org "github.com/openmerlin/merlin-server/organization/domain"
	"github.com/openmerlin/merlin-server/user/domain"
)

// ListOption is a struct for defining options when listing resources.
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

// User is an interface for user-related operations.
type User interface {
	AddUser(*domain.User) (domain.User, error)
	SaveUser(*domain.User) (domain.User, error)
	DeleteUser(*domain.User) error
	GetByAccount(domain.Account) (domain.User, error)
	GetUserAvatarId(domain.Account) (primitive.AvatarId, error)
	GetUsersAvatarId([]string) ([]domain.User, error)
	GetUserFullname(domain.Account) (string, error)

	AddOrg(*org.Organization) (org.Organization, error)
	SaveOrg(*org.Organization) (org.Organization, error)
	DeleteOrg(*org.Organization) error
	CheckName(primitive.Account) bool
	GetOrgByName(primitive.Account) (org.Organization, error)
	GetOrgByOwner(primitive.Account) ([]org.Organization, error)
	GetOrgCountByOwner(primitive.Account) (int64, error)

	ListAccount(*ListOption) ([]domain.User, int, error)
}
