package domain

import (
	"errors"
	"fmt"
)

// user
type User struct {
	Id string

	Email   Email
	Account Account

	PlatformPwd    string                   //password for git user
	PlatformId     int64                    // id in gitea
	PlatformTokens map[string]PlatformToken // token for git user

	Bio      Bio
	AvatarId AvatarId

	Version int
}

type UserInfo struct {
	Account  Account
	AvatarId AvatarId
}

type TokenPerm string

const (
	TokenPermWrite TokenPerm = "write"
	TokenPermRead  TokenPerm = "read"
)

type PlatformToken struct {
	Token   string
	Name    string
	Account Account
	// TODO: expire not honor by gitea
	Expire     int64 // timeout in seconds
	CreatedAt  int64
	Permission TokenPerm
}

func ToPerms(t TokenPerm) []string {
	switch t {
	case TokenPermWrite:
		return []string{"write:organization", "write:repository"}
	case TokenPermRead:
		return []string{"read:organization", "read:repository"}
	default:
		return []string{}
	}
}

// user
type UserCreateCmd struct {
	Email    Email
	Account  Account
	Bio      Bio
	AvatarId AvatarId
}

type TokenCreatedCmd struct {
	Account    Account // user name
	Name       string  // name of the token
	Expire     int64   // timeout in seconds
	Permission TokenPerm
}

func (cmd TokenCreatedCmd) Validate() error {
	if cmd.Permission == "" {
		return fmt.Errorf("missing permission when creating token")
	}

	if cmd.Permission != TokenPermWrite &&
		cmd.Permission != TokenPermRead {
		return fmt.Errorf("invalid permission when creating token")
	}

	if cmd.Name == "" {
		return fmt.Errorf("missing name when creating token")
	}

	return nil
}

type TokenDeletedCmd struct {
	Account Account // user name
	Name    string  // name of the token
}

func (cmd TokenDeletedCmd) Validate() error {
	if cmd.Name == "" {
		return fmt.Errorf("missing name when creating token")
	}

	return nil
}

type FollowerInfo struct {
	User     Account
	Follower Account
}

type FollowerUserInfo struct {
	Account    Account
	AvatarId   AvatarId
	Bio        Bio
	IsFollower bool
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
