package repository

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/session/domain"
)

type LoginRepositoryAdapter interface {
	Add(*domain.Login) error

	Delete(primitive.UUID) error

	// sort by created_at aesc
	FindByUser(primitive.Account) ([]domain.Login, error)

	Find(primitive.UUID) (domain.Login, error)
}

type CSRFTokenRepositoryAdapter interface {
	Add(*domain.CSRFToken) error
	Save(*domain.CSRFToken) error
	Find(primitive.UUID) (domain.CSRFToken, error)
}
