/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package repository

import (
	"context"

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

// Pagination calculates the offset for pagination.
func (opt *ListOption) Pagination() (bool, int) {
	if opt.PageNum > 0 && opt.CountPerPage > 0 {
		return true, (opt.PageNum - 1) * opt.CountPerPage
	}

	return false, 0
}

// ListOrgOption is a struct for defining options when listing organization resources.
type ListOrgOption struct {
	OrgIDs []int64
	Owner  primitive.Account
}

type ListPageOrgOption struct {
	OrgIDs   []int64
	Owner    primitive.Account
	PageNum  int
	Count    bool
	PageSize int
}

// User is an interface for user-related operations.
type User interface {
	AddUser(context.Context, *domain.User) (domain.User, error)
	SaveUser(context.Context, *domain.User) (domain.User, error)
	DeleteUser(context.Context, *domain.User) error
	GetByAccount(context.Context, domain.Account) (domain.User, error)
	GetUserAvatarId(context.Context, domain.Account) (primitive.AvatarId, error)
	GetUsersAvatarId(context.Context, []string) ([]domain.User, error)
	GetUserFullname(context.Context, domain.Account) (string, error)

	AddOrg(context.Context, *org.Organization) (org.Organization, error)
	SaveOrg(context.Context, *org.Organization) (org.Organization, error)
	DeleteOrg(context.Context, *org.Organization) error
	CheckName(context.Context, primitive.Account) bool
	GetOrgByName(context.Context, primitive.Account) (org.Organization, error)
	GetOrgByOwner(context.Context, primitive.Account) ([]org.Organization, error)
	GetOrgList(context.Context, *ListOrgOption) ([]org.Organization, error)
	GetOrgPageList(context.Context, *ListPageOrgOption) ([]org.Organization, int, error)
	GetOrgCountByOwner(context.Context, primitive.Account) (int64, error)

	ListAccount(context.Context, *ListOption) ([]domain.User, int, error)

	SearchUser(context.Context, *ListOption) ([]domain.User, int, error)
	SearchOrg(context.Context, *ListOption) ([]org.Organization, int, error)
}
