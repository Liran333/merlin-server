/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package repository provides interfaces for interacting with computility-related data.
package repository

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/computility/domain"
)

// ComputilityOrgRepositoryAdapter is an interface for interacting with computility org repositories.
type ComputilityOrgRepositoryAdapter interface {
	Delete(primitive.Identity) error
	FindByOrgName(primitive.Account) (domain.ComputilityOrg, error)

	OrgAssignQuota(domain.ComputilityOrg, int) error
	OrgRecallQuota(domain.ComputilityOrg, int) error
}

// ComputilityDetailRepositoryAdapter is an interface for interacting with computility detail repositories.
type ComputilityDetailRepositoryAdapter interface {
	Add(*domain.ComputilityDetail) error
	Delete(primitive.Identity) error
	FindByIndex(*domain.ComputilityIndex) (domain.ComputilityDetail, error)
	GetMembers(primitive.Account) ([]domain.ComputilityDetail, error)
}

// ComputilityAccountRepositoryAdapter is an interface for interacting with computility account repositories.
type ComputilityAccountRepositoryAdapter interface {
	Add(*domain.ComputilityAccount) error
	Delete(primitive.Identity) error
	FindByAccountIndex(domain.ComputilityAccountIndex) (domain.ComputilityAccount, error)
	CheckAccountExist(primitive.Account) (bool, error)

	DecreaseAccountAssignedQuota(domain.ComputilityAccount, int) error
	IncreaseAccountAssignedQuota(domain.ComputilityAccount, int) error

	ConsumeQuota(domain.ComputilityAccount, int) error
	ReleaseQuota(domain.ComputilityAccount, int) error

	CancelAccount(domain.ComputilityAccountIndex) error
}
