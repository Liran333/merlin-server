/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package domain

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/pbkdf2"

	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/utils"
)

// OrgRole is a type alias for string.
type OrgRole = string

// UserType is a type alias for int.
type UserType = int

const (
	tokenPermDenied = "token permission denied"
	tokenInvalid    = "token invalid"
	tokenExpired    = "token expired"

	// UserTypeUser is const
	UserTypeUser UserType = 0
	// UserTypeOrganization is const
	UserTypeOrganization UserType = 1
	iter                 int      = 10000
	keyLen               int      = 32 // keyLen is const
)

// User is a struct representing a user.
type User struct {
	Id              primitive.Identity
	Email           primitive.Email
	Phone           primitive.Phone
	Account         Account
	Fullname        primitive.AccountFullname
	PlatformPwd     string // password for git user
	PlatformId      int64  // id in gitea
	Website         primitive.Website
	Owner           primitive.Account
	OwnerId         primitive.Identity
	WriteTeamId     int64
	ReadTeamId      int64
	OwnerTeamId     int64
	CreatedAt       int64
	UpdatedAt       int64
	Desc            primitive.AccountDesc
	AvatarId        primitive.AvatarId
	Type            UserType
	DefaultRole     primitive.Role
	AllowRequest    bool
	RequestDelete   bool
	RequestDeleteAt int64
	Version         int
	IsAgreePrivacy  bool
}

// IsOrganization checks if the user is an organization.
func (u *User) IsOrganization() bool {
	return u.Type == UserTypeOrganization
}

// AgreePrivacy sets the IsAgreePrivacy field of the User instance to true,
// indicating that the user has agreed to the privacy terms.
func (u *User) AgreePrivacy() {
	u.IsAgreePrivacy = true
}

// RevokePrivacy sets the IsAgreePrivacy field of the User instance to false,
// indicating that the user has revoked their agreement to the privacy terms.
func (u *User) RevokePrivacy() {
	u.IsAgreePrivacy = false
}

// UserInfo represents additional information about a user.
type UserInfo struct {
	Account  Account
	AvatarId primitive.AvatarId
}

// PlatformToken represents a token associated with a platform account.
type PlatformToken struct {
	Id         primitive.Identity
	Token      string
	Name       primitive.TokenName
	Account    Account            // owner name
	OwnerId    primitive.Identity // owner id
	Expire     int64              // timeout in seconds
	CreatedAt  int64
	UpdatedAt  int64
	Permission primitive.TokenPerm
	Salt       string
	LastEight  string
	Version    int
}

func (t PlatformToken) isExpired() bool {
	return t.Expire != 0 && t.Expire < utils.Now()
}

// Check checks if the given token is valid and has the required permission.
func (t PlatformToken) Check(token string, perm primitive.TokenPerm) error {
	if t.isExpired() {
		return allerror.NewNoPermission(tokenExpired, errors.New(tokenExpired))
	}

	if !t.Match(token) {
		return allerror.NewNoPermission(tokenInvalid, errors.New(tokenInvalid))
	}

	if !t.Permission.PermissionAllow(perm) {
		return allerror.NewNoPermission(tokenPermDenied, errors.New(tokenPermDenied))
	}

	return nil
}

// Match checks if the given token matches the stored token.
func (t PlatformToken) Match(token string) bool {
	saltBtye, err := base64.RawStdEncoding.DecodeString(t.Salt)
	if err != nil {
		return false
	}

	srcBtye, err := base64.RawStdEncoding.DecodeString(t.Token)
	if err != nil {
		return false
	}

	derivedKey := pbkdf2.Key([]byte(token), saltBtye, iter, keyLen, sha256.New)

	return bytes.Equal(srcBtye, derivedKey)
}

// EncryptToken encrypts the given token using a randomly generated salt.
func EncryptToken(token string) (enc, salt string, err error) {
	saltBtye := make([]byte, keyLen)
	_, err = rand.Read(saltBtye)
	if err != nil {
		return
	}

	encBytes := pbkdf2.Key([]byte(token), saltBtye, iter, keyLen, sha256.New)

	enc = base64.RawStdEncoding.EncodeToString(encBytes)
	salt = base64.RawStdEncoding.EncodeToString(saltBtye)
	return
}

// ToPerms converts the given TokenPerm to a slice of strings representing the permissions.
func ToPerms(t primitive.TokenPerm) []string {
	switch t.TokenPerm() {
	case primitive.TokenPermWrite:
		return []string{"write:organization", "write:repository", "write:user"}
	case primitive.TokenPermRead:
		return []string{"read:organization", "read:repository", "read:user"}
	default:
		return []string{}
	}
}

// UserCreateCmd is a struct for creating a user.
type UserCreateCmd struct {
	Email    primitive.Email
	Account  Account
	Desc     primitive.AccountDesc
	AvatarId primitive.AvatarId
	Fullname primitive.AccountFullname
	Phone    primitive.Phone
}

// TokenCreatedCmd is a struct for creating a token.
type TokenCreatedCmd struct {
	Account    Account             // user name
	Name       primitive.TokenName // name of the token
	Expire     int64               // timeout in seconds
	Permission primitive.TokenPerm
}

// Validate validates the TokenCreatedCmd.
func (cmd TokenCreatedCmd) Validate() error {
	if cmd.Name == nil {
		e := fmt.Errorf("missing name when creating token")
		return allerror.New(allerror.ErrorMissingName, "", e)
	}

	if cmd.Account == nil {
		e := fmt.Errorf("missing account when creating token")
		return allerror.New(allerror.ErrorMissingAccount, "", e)
	}

	return nil
}

// TokenDeletedCmd is a struct for deleting a token.
type TokenDeletedCmd struct {
	Account Account             // actor username
	Name    primitive.TokenName // name of the token
}

// Validate validates the TokenDeletedCmd.
func (cmd TokenDeletedCmd) Validate() error {
	if cmd.Account == nil {
		e := fmt.Errorf("missing account when delete token")
		return allerror.New(allerror.ErrorMissingAccount, "", e)
	}

	if cmd.Name == nil {
		e := fmt.Errorf("missing name when delete token")
		return allerror.New(allerror.ErrorMissingName, "", e)
	}

	return nil
}

// FollowerInfo is a struct for storing follower information.
type FollowerInfo struct {
	User     Account
	Follower Account
}

// FollowerUserInfo is a struct for storing follower user information.
type FollowerUserInfo struct {
	Account    Account
	AvatarId   primitive.AvatarId
	Desc       primitive.AccountDesc
	IsFollower bool
}

// Validate validates the UserCreateCmd.
func (cmd *UserCreateCmd) Validate() error {
	b := cmd.Email != nil &&
		cmd.Account != nil &&
		cmd.Fullname != nil &&
		cmd.Phone != nil

	if !b {
		return errors.New("invalid cmd of creating user")
	}

	return nil
}

// ToUser converts UserCreateCmd to User.
func (cmd *UserCreateCmd) ToUser() User {
	return User{
		Email:    cmd.Email,
		Account:  cmd.Account,
		Desc:     cmd.Desc,
		AvatarId: cmd.AvatarId,
		Fullname: cmd.Fullname,
		Type:     UserTypeUser,
		Phone:    cmd.Phone,
	}
}
