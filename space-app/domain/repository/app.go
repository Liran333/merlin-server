package repository

import "github.com/openmerlin/merlin-server/space-app/domain"

type Repository interface {
	Add(*domain.SpaceApp) error
	Find(*domain.SpaceAppIndex) (domain.SpaceApp, error)
	Save(*domain.SpaceApp) error
}
