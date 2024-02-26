/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package repository

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/organization/domain"
)

// Organization interface defines the methods for managing organizations.
type Organization interface {
	AddOrg(*domain.Organization) (domain.Organization, error)
	SaveOrg(*domain.Organization) (domain.Organization, error)
	DeleteOrg(*domain.Organization) error
	CheckName(primitive.Account) bool
	GetOrgByName(primitive.Account) (domain.Organization, error)
	GetOrgByOwner(primitive.Account) ([]domain.Organization, error)
}
