/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package platform provides interfaces for interacting
// with the platform's authentication and organization functionality.
package platform

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	org "github.com/openmerlin/merlin-server/organization/domain"
	"github.com/openmerlin/merlin-server/user/domain"
)

// BaseAuthClient is an interface that defines the methods required for authentication and organization management.
type BaseAuthClient interface {
	CreateToken(*domain.TokenCreatedCmd) (domain.PlatformToken, error)
	DeleteToken(*domain.TokenDeletedCmd) error
	CreateOrg(*org.Organization) error
	DeleteOrg(primitive.Account) error
	CanDelete(primitive.Account) (bool, error)
	AddMember(*org.Organization, *org.OrgMember) error
	RemoveMember(*org.Organization, *org.OrgMember) error
	EditMemberRole(*org.Organization, primitive.Role, *org.OrgMember) error
}
