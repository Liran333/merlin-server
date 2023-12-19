package repository

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/user/domain"
)

type Token interface {
	Save(*domain.PlatformToken) (domain.PlatformToken, error)
	Delete(primitive.Account, string) error
	GetByAccount(primitive.Account) ([]domain.PlatformToken, error)
	GetByLastEight(string) ([]domain.PlatformToken, error)

	//Search(*UserSearchOption) (UserSearchResult, error)
}
