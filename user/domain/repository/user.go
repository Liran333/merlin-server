package repository

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	org "github.com/openmerlin/merlin-server/organization/domain"
	"github.com/openmerlin/merlin-server/user/domain"
)

type FollowFindOption struct {
	Follower domain.Account

	CountPerPage int
	PageNum      int
}

type FollowerUserInfos struct {
	Users []domain.FollowerUserInfo
	Total int
}

type UserSearchOption struct {
	// can't define Name as domain.Account
	// because the Name can be subpart of the real account
	Name   string
	TopNum int
}

type UserSearchResult struct {
	Top []domain.Account

	Total int
}

type User interface {
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
}
