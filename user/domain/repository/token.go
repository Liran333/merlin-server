package repository

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/user/domain"
)

type Token interface {
	Add(*domain.PlatformToken) (domain.PlatformToken, error)
	Delete(primitive.Account, primitive.TokenName) error
	GetByAccount(primitive.Account) ([]domain.PlatformToken, error)
	GetByLastEight(string) ([]domain.PlatformToken, error)
	GetByName(primitive.Account, primitive.TokenName) (domain.PlatformToken, error)
}
