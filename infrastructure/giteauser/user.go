/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package giteauser provides functionality for interacting with user accounts in a Gitea instance.
package giteauser

import (
	"fmt"
	"net/http"

	"github.com/openmerlin/go-sdk/gitea"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/user/domain"
)

// UserCreateCmd represents the command to create a user.
type UserCreateCmd struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// UserUpdateCmd represents the command to update a user.
type UserUpdateCmd = domain.UserCreateCmd

func toUser(u *gitea.User) (user domain.User, err error) {
	user.Account, err = primitive.NewAccount(u.UserName)
	if err != nil {
		return
	}

	user.Email, err = primitive.NewEmail(u.Email)
	if err != nil {
		return
	}

	user.PlatformId = u.ID

	return
}

// GetClient returns a UserClient with the provided gitea client.
func GetClient(c *gitea.Client) *UserClient {
	return &UserClient{
		client: c,
	}
}

// UserClient represents the admin client for user management.
type UserClient struct {
	client *gitea.Client
}

// CreateUser creates a new user with the provided command.
func (c *UserClient) CreateUser(cmd *UserCreateCmd) (user domain.User, err error) {
	changePwd := false

	pwd, err := primitive.NewPassword()
	if err != nil {
		return
	}
	defer pwd.Clear()

	o := gitea.CreateUserOption{
		Username:           cmd.Username,
		Email:              cmd.Email,
		Password:           pwd.Password(),
		MustChangePassword: &changePwd,
	}

	u, _, err := c.client.AdminCreateUser(o)
	if err != nil {
		return
	}

	user, err = toUser(u)
	user.PlatformPwd = pwd.Password()

	return
}

// DeleteUser deletes the user with the specified name.
func (c *UserClient) DeleteUser(name string) error {
	resp, err := c.client.AdminDeleteUser(name)
	if resp.StatusCode == http.StatusNotFound {
		return nil
	}

	return err
}

// DeleteOrg delete the org with the specified name.
func (c *UserClient) DeleteOrg(name string) error {
	resp, err := c.client.AdminDeleteOrg(name)
	if resp.StatusCode == http.StatusNotFound {
		return nil
	}

	return err
}

// UpdateUser updates the user with the provided command.
func (c *UserClient) UpdateUser(cmd *UserUpdateCmd) (err error) {
	if cmd == nil {
		return fmt.Errorf("cmd is nil")
	}

	if cmd.Account == nil {
		return fmt.Errorf("account is nil")
	}

	d := gitea.EditUserOption{
		LoginName: cmd.Account.Account(),
	}

	if cmd.Email != nil {
		email := cmd.Email.Email()
		d.Email = &email
	}

	if cmd.Desc != nil {
		desc := cmd.Desc.AccountDesc()
		d.Description = &desc
	}

	if cmd.Fullname != nil {
		fullname := cmd.Fullname.AccountFullname()
		d.FullName = &fullname
	}

	_, err = c.client.AdminEditUser(cmd.Account.Account(), d)

	return
}
