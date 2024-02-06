package repository

import "github.com/openmerlin/merlin-server/space-app/domain"

type Repository interface {
	Add(*domain.SpaceApp) error
}
