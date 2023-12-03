package git

import (
	git "github.com/openmerlin/merlin-server/infrastructure/gitea"
	"github.com/openmerlin/merlin-server/user/domain"
)

type User interface {
	Create(*domain.UserCreateCmd) error
}

func NewUserGit(c *git.Client) User {
	return &userGitImpl{c}
}

type userGitImpl struct {
	client *git.Client
}

func (u *userGitImpl) Create(cmd *domain.UserCreateCmd) error {
	return u.client.CreateUser(&git.UserCreateCmd{
		Username: cmd.Account.Account(),
		Email:    cmd.Email.Email(),
	})
}
