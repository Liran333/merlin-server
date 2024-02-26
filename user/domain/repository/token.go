/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package repository provides interfaces for managing platform tokens in the user domain.
package repository

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/user/domain"
)

// Token represents an interface for managing platform tokens.
type Token interface {
	Add(*domain.PlatformToken) (domain.PlatformToken, error)
	Delete(primitive.Account, primitive.TokenName) error
	GetByAccount(primitive.Account) ([]domain.PlatformToken, error)
	GetByLastEight(string) ([]domain.PlatformToken, error)
	GetByName(primitive.Account, primitive.TokenName) (domain.PlatformToken, error)
}
