package platform

import "github.com/openmerlin/merlin-server/user/domain"

type User interface {
	Create(*domain.UserCreateCmd) (domain.User, error)
	Delete(user *domain.User) error
}
