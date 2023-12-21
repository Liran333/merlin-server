package domain

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

// user
type User struct {
	Id string

	Email    Email
	Account  Account
	Fullname string

	PlatformPwd string //password for git user
	PlatformId  int64  // id in gitea

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
	Id      string
	Token   string
	Name    string
	Account Account
	// TODO: expire not honor by gitea
	Expire     int64 // timeout in seconds
	CreatedAt  int64
	Permission TokenPerm
	Salt       string
	LastEight  string
	Version    int
}

func (t PlatformToken) Compare(token string) bool {
	saltBtye, err := base64.RawStdEncoding.DecodeString(t.Salt)
	if err != nil {
		return false
	}

	srcBtye, err := base64.RawStdEncoding.DecodeString(t.Token)
	if err != nil {
		return false
	}

	dstBytes := pbkdf2.Key([]byte(token), saltBtye, 10000, 32, sha256.New)

	return bytes.Equal(srcBtye, dstBytes)
}

func EncryptToken(token string) (enc, salt string, err error) {
	saltBtye := make([]byte, 32)
	_, err = rand.Read(saltBtye)
	if err != nil {
		return
	}

	encBytes := pbkdf2.Key([]byte(token), saltBtye, 10000, 32, sha256.New)

	enc = base64.RawStdEncoding.EncodeToString(encBytes)
	salt = base64.RawStdEncoding.EncodeToString(saltBtye)
	return
}

func ToPerms(t TokenPerm) []string {
	switch t {
	case TokenPermWrite:
		return []string{"write:organization", "write:repository", "write:user"}
	case TokenPermRead:
		return []string{"read:organization", "read:repository", "read:user"}
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
	Fullname string
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
		Fullname: cmd.Fullname,
	}
}
