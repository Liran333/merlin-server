/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package repository provides interfaces for managing approvals in an organization.
package repository

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/organization/domain"
)

// OrgMember interface defines the methods for managing organization members.
type OrgMember interface {
	Add(*domain.OrgMember) (domain.OrgMember, error)
	Save(*domain.OrgMember) (domain.OrgMember, error)
	Delete(*domain.OrgMember) error
	DeleteByOrg(primitive.Account) error
	GetByOrg(*domain.OrgListMemberCmd) ([]domain.OrgMember, error)
	GetByOrgAndRole(string, primitive.Role) ([]domain.OrgMember, error)
	GetByOrgAndUser(org, user string) (domain.OrgMember, error)
	GetByUser(string) ([]domain.OrgMember, error)
	GetOrgNamesByUserName(string) ([]string, error)
	GetByUserAndRoles(primitive.Account, []primitive.Role) ([]domain.OrgMember, error)
}
