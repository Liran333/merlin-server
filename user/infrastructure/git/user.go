/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package git

import (
	"fmt"

	common "github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/infrastructure/giteauser"
	"github.com/openmerlin/merlin-server/user/domain"
)

// User is an interface for user operations.
type User interface {
	Create(*domain.UserCreateCmd) (domain.User, error)
	Delete(common.Account) error
	DeleteOrg(common.Account) error
	Update(*domain.UserCreateCmd) error
}

// NewUserGit creates a new instance of userGitImpl with the given giteauser.UserClient.
func NewUserGit(c *giteauser.UserClient) User {
	return &userGitImpl{client: c}
}

type userGitImpl struct {
	client *giteauser.UserClient
}

// Create creates a new user using the provided command.
func (u *userGitImpl) Create(cmd *domain.UserCreateCmd) (domain.User, error) {
	return u.client.CreateUser(&giteauser.UserCreateCmd{
		Username: cmd.Account.Account(),
		Email:    cmd.Email.Email(),
	})
}

// Delete deletes the specified user.
func (u *userGitImpl) Delete(user common.Account) error {
	if user == nil {
		return fmt.Errorf("account is nil")
	}
	return u.client.DeleteUser(user.Account())
}

// Update updates the user with the provided command.
func (u *userGitImpl) Update(cmd *domain.UserCreateCmd) error {
	return u.client.UpdateUser(cmd)
}

func (u *userGitImpl) DeleteOrg(name common.Account) error {
	if name == nil {
		return fmt.Errorf("account is nil")
	}
	return u.client.DeleteOrg(name.Account())
}
