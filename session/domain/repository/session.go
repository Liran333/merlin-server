/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package repository

import (
	"context"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/session/domain"
)

// SessionRepositoryAdapter is an interface that defines the methods for interacting with the login repository.
type SessionRepositoryAdapter interface {
	Add(*domain.Session) error

	Delete(context.Context, primitive.RandomId) error

	DeleteByUser(primitive.Account) error

	// FindByUser sort by created_at aesc
	FindByUser(primitive.Account) ([]domain.Session, error)

	Find(context.Context, primitive.RandomId) (domain.Session, error)
}

// CSRFTokenRepositoryAdapter is an interface that defines the methods for interacting with the CSRF token repository.
type CSRFTokenRepositoryAdapter interface {
	Add(primitive.RandomId, *domain.CSRFToken) error
	Find(primitive.RandomId) (domain.CSRFToken, error)
}

// SessionFastRepositoryAdapter defines the interface for an adapter that interacts with a fast session repository.
type SessionFastRepositoryAdapter interface {
	Add(*domain.Session) error
	Save(*domain.Session) error
	Find(primitive.RandomId) (domain.Session, error)
	Delete(primitive.RandomId) error
}
