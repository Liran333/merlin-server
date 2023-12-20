package git

import (
	"github.com/openmerlin/merlin-server/infrastructure/giteauser"
	"github.com/openmerlin/merlin-server/user/domain"
)

type User interface {
	Create(*domain.UserCreateCmd) (domain.User, error)
	Delete(user *domain.User) error
}

func NewUserGit(c *giteauser.UserClient) User {
	return &userGitImpl{c}
}

type userGitImpl struct {
	client *giteauser.UserClient
}

func (u *userGitImpl) Create(cmd *domain.UserCreateCmd) (domain.User, error) {
	return u.client.CreateUser(&giteauser.UserCreateCmd{
		Username: cmd.Account.Account(),
		Email:    cmd.Email.Email(),
	})
}

func (u *userGitImpl) Delete(user *domain.User) error {
	return u.client.DeleteUser(user.Account.Account())
}
