/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package repository

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/session/domain"
)

// LoginRepositoryAdapter is an interface that defines the methods for interacting with the login repository.
type LoginRepositoryAdapter interface {
	Add(*domain.Login) error

	Delete(primitive.UUID) error

	// FindByUser sort by created_at aesc
	FindByUser(primitive.Account) ([]domain.Login, error)

	Find(primitive.UUID) (domain.Login, error)
}

// CSRFTokenRepositoryAdapter is an interface that defines the methods for interacting with the CSRF token repository.
type CSRFTokenRepositoryAdapter interface {
	Add(*domain.CSRFToken) error
	Save(*domain.CSRFToken) error
	Find(primitive.UUID) (domain.CSRFToken, error)
}
