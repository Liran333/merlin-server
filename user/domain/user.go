package domain

import "errors"

// user
type User struct {
	Id      string
	Email   Email
	Account Account

	Bio      Bio
	AvatarId AvatarId

	Version int
}

type UserInfo struct {
	Account  Account
	AvatarId AvatarId
}

// user
type UserCreateCmd struct {
	Email    Email
	Account  Account
	Bio      Bio
	AvatarId AvatarId
}

func (cmd *UserCreateCmd) Validate() error {
	b := cmd.Email != nil &&
		cmd.Account != nil

	if !b {
		return errors.New("invalid cmd of creating user")
	}

	return nil
}

func (cmd *UserCreateCmd) ToUser() User {
	return User{
		Email:   cmd.Email,
		Account: cmd.Account,

		Bio:      cmd.Bio,
		AvatarId: cmd.AvatarId,
	}
}

func (cmd *UserCreateCmd) toUserDTO() User {
	return User{
		Email:   cmd.Email,
		Account: cmd.Account,

		Bio:      cmd.Bio,
		AvatarId: cmd.AvatarId,
	}
}
