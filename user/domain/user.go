package domain

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/pbkdf2"

	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/utils"
)

type OrgRole = string
type UserType = int

const (
	tokenPermDenied               = "token permission denied"
	tokenInvalid                  = "token invalid"
	tokenExpired                  = "token expired"
	OrgRoleContributor   OrgRole  = "contributor" // in contributor team
	OrgRoleReader        OrgRole  = "read"        // in read team
	OrgRoleWriter        OrgRole  = "write"       // in write team
	OrgRoleAdmin         OrgRole  = "admin"       // in owner team
	UserTypeUser         UserType = 0
	UserTypeOrganization UserType = 1
)

// user
type User struct {
	Id                primitive.Identity
	Email             primitive.Email
	Phone             primitive.Phone
	Account           Account
	Fullname          primitive.MSDFullname
	PlatformPwd       string //password for git user
	PlatformId        int64  // id in gitea
	Website           string
	Owner             primitive.Account
	OwnerId           primitive.Identity
	WriteTeamId       int64
	ReadTeamId        int64
	OwnerTeamId       int64
	ContributorTeamId int64
	CreatedAt         int64
	UpdatedAt         int64
	Desc              primitive.MSDDesc
	AvatarId          primitive.AvatarId
	Type              UserType
	DefaultRole       OrgRole
	AllowRequest      bool
	Version           int
}

func (u User) IsOrganization() bool {
	return u.Type == UserTypeOrganization
}

func (u *User) ClearSenstiveData() {
	u.Email = nil
	u.Phone = nil
}

type UserInfo struct {
	Account  Account
	AvatarId primitive.AvatarId
}

type PlatformToken struct {
	Id      primitive.Identity
	Token   string
	Name    primitive.TokenName
	Account Account            // owner name
	OwnerId primitive.Identity // owner id
	// TODO: expire not honor by gitea
	Expire     int64 // timeout in seconds
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

func (t PlatformToken) Check(token string, perm primitive.TokenPerm) error {
	if t.isExpired() {
		return allerror.NewNoPermission(tokenExpired)
	}

	if !t.Match(token) {
		return allerror.NewNoPermission(tokenInvalid)
	}

	if !t.Permission.PermissionAllow(perm) {
		return allerror.NewNoPermission(tokenPermDenied)
	}

	return nil
}

func (t PlatformToken) Match(token string) bool {
	saltBtye, err := base64.RawStdEncoding.DecodeString(t.Salt)
	if err != nil {
		return false
	}

	srcBtye, err := base64.RawStdEncoding.DecodeString(t.Token)
	if err != nil {
		return false
	}

	derivedKey := pbkdf2.Key([]byte(token), saltBtye, 10000, 32, sha256.New)

	return bytes.Equal(srcBtye, derivedKey)
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

// user
type UserCreateCmd struct {
	Email    primitive.Email
	Account  Account
	Desc     primitive.MSDDesc
	AvatarId primitive.AvatarId
	Fullname primitive.MSDFullname
	Phone    primitive.Phone
}

type TokenCreatedCmd struct {
	Account    Account             // user name
	Name       primitive.TokenName // name of the token
	Expire     int64               // timeout in seconds
	Permission primitive.TokenPerm
}

func (cmd TokenCreatedCmd) Validate() error {
	if cmd.Name == nil {
		return allerror.NewInvalidParam("missing name when creating token")
	}

	if cmd.Account == nil {
		return allerror.NewInvalidParam("missing account when creating token")
	}

	return nil
}

type TokenDeletedCmd struct {
	Account Account             // actor user name
	Name    primitive.TokenName // name of the token
}

func (cmd TokenDeletedCmd) Validate() error {
	if cmd.Account == nil {
		return allerror.NewInvalidParam("missing account when delete token")
	}

	if cmd.Name == nil {
		return allerror.NewInvalidParam("missing name when delete token")
	}

	return nil
}

type FollowerInfo struct {
	User     Account
	Follower Account
}

type FollowerUserInfo struct {
	Account    Account
	AvatarId   primitive.AvatarId
	Desc       primitive.MSDDesc
	IsFollower bool
}

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
