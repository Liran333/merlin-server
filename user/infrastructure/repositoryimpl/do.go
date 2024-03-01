/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package repositoryimpl provides the implementation of repository interfaces for user and token entities.
package repositoryimpl

import (
	"github.com/openmerlin/merlin-server/common/domain/crypto"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	org "github.com/openmerlin/merlin-server/organization/domain"
	"github.com/openmerlin/merlin-server/user/domain"
)

const (
	fieldID         = "id"
	fieldName       = "name"
	fieldCount      = "count"
	fieldPlatformId = "platform_id"
	fieldAccount    = "account"
	fieldType       = "type"
	fieldAvatarId   = "avatar_id"
	fieldFullname   = "fullname"
	fieldVersion    = "version"
	fieldOwner      = "owner"
	fieldLastEight  = "last_eight"
	fieldCreatedAt  = "created_at"
	fieldUpdatedAt  = "updated_at"
)

var (
	tokenTableName string
	userTableName  string
)

// UserDO represents the database model for the User entity.
type UserDO struct {
	postgresql.CommonModel
	Desc              string `gorm:"column:desc"`
	Name              string `gorm:"column:name;uniqueIndex:name_index"`
	Fullname          string `gorm:"column:fullname"`
	CreatedBy         string `gorm:"column:created_by"`
	Email             string `gorm:"column:email"`
	Phone             string `gorm:"column:phone"`
	AvatarId          string `gorm:"column:avatar_id"`
	Type              int    `gorm:"column:type;index:type_index"`
	Website           string `gorm:"column:website"`
	OwnerId           int64  `gorm:"column:owner_id;index:owner_id_index"`
	Owner             string `gorm:"column:owner;index:owner_index"`
	DefaultRole       string `gorm:"column:default_role"`
	AllowRequest      bool   `gorm:"column:allow_request"`
	OwnerTeamId       int64  `gorm:"column:owner_team_id"`
	ReadTeamId        int64  `gorm:"column:read_team_id"`
	WriteTeamId       int64  `gorm:"column:write_team_id"`
	ContributorTeamId int64  `gorm:"column:contributor_team_id"`
	PlatformId        int64  `gorm:"column:platform_id"`
	PlatformPwd       string `gorm:"column:platform_pwd"`
	Version           int    `gorm:"column:version"`
}

func toUserDONoEnc(u *domain.User) (do UserDO) {
	do = UserDO{
		Name:        u.Account.Account(),
		Fullname:    u.Fullname.MSDFullname(),
		Email:       u.Email.Email(),
		AvatarId:    u.AvatarId.AvatarId(),
		Type:        u.Type,
		PlatformId:  u.PlatformId,
		PlatformPwd: u.PlatformPwd,
		Version:     u.Version,
		Phone:       u.Phone.PhoneNumber(),
	}

	if u.Desc != nil {
		do.Desc = u.Desc.MSDDesc()
	}

	do.ID = u.Id.Integer()

	return
}

func toUserDO(u *domain.User, e crypto.Encrypter) (do UserDO, err error) {
	email, err := e.Encrypt(u.Email.Email())
	if err != nil {
		return
	}

	phone, err := e.Encrypt(u.Phone.PhoneNumber())
	if err != nil {
		return
	}

	pwd, err := e.Encrypt(u.PlatformPwd)
	if err != nil {
		return
	}

	do = UserDO{
		Name:        u.Account.Account(),
		Fullname:    u.Fullname.MSDFullname(),
		Email:       email,
		AvatarId:    u.AvatarId.AvatarId(),
		Type:        u.Type,
		PlatformId:  u.PlatformId,
		PlatformPwd: pwd,
		Version:     u.Version,
		Phone:       phone,
	}

	if u.Desc != nil {
		do.Desc = u.Desc.MSDDesc()
	}

	do.ID = u.Id.Integer()

	return
}

func toOrgDO(u *org.Organization) (do UserDO) {
	do = UserDO{
		Desc:              u.Desc.MSDDesc(),
		Name:              u.Account.Account(),
		Fullname:          u.Fullname.MSDFullname(),
		AvatarId:          u.AvatarId.AvatarId(),
		Type:              u.Type,
		PlatformId:        u.PlatformId,
		Version:           u.Version,
		AllowRequest:      u.AllowRequest,
		DefaultRole:       u.DefaultRole,
		OwnerTeamId:       u.OwnerTeamId,
		ReadTeamId:        u.ReadTeamId,
		WriteTeamId:       u.WriteTeamId,
		ContributorTeamId: u.ContributorTeamId,
		Owner:             u.Owner.Account(),
		OwnerId:           u.OwnerId.Integer(),
		Website:           u.Website,
	}

	do.ID = u.Id.Integer()

	return
}

func (u *UserDO) toUser(e crypto.Encrypter) (domain.User, error) {
	return u.toOrg(e)
}

func (u *UserDO) toUserNoEnc() domain.User {
	return u.toOrgNoEnc()
}

func (u *UserDO) toOrgNoEnc() (o org.Organization) {

	o = org.Organization{
		Id:                primitive.CreateIdentity(u.ID),
		Version:           u.Version,
		Desc:              primitive.CreateMSDDesc(u.Desc),
		Account:           primitive.CreateAccount(u.Name),
		Fullname:          primitive.CreateMSDFullname(u.Fullname),
		AvatarId:          primitive.CreateAvatarId(u.AvatarId),
		PlatformId:        u.PlatformId,
		CreatedAt:         u.CreatedAt.Unix(),
		UpdatedAt:         u.UpdatedAt.Unix(),
		PlatformPwd:       u.PlatformPwd,                        // user only
		Email:             primitive.CreateEmail(u.Email),       // user only
		Phone:             primitive.CreatePhoneNumber(u.Phone), // user only,
		AllowRequest:      u.AllowRequest,                       // org only
		DefaultRole:       u.DefaultRole,                        // org only
		OwnerTeamId:       u.OwnerTeamId,                        // org only
		ReadTeamId:        u.ReadTeamId,                         // org only
		WriteTeamId:       u.WriteTeamId,                        // org only
		ContributorTeamId: u.ContributorTeamId,                  // org only
		Owner:             primitive.CreateAccount(u.Owner),     // org only
		OwnerId:           primitive.CreateIdentity(u.OwnerId),  // org only
		Website:           u.Website,                            // org only
		Type:              u.Type,
	}

	return
}

func (u *UserDO) toOrg(e crypto.Encrypter) (o org.Organization, err error) {
	email, err := e.Decrypt(u.Email)
	if err != nil {
		return
	}

	phone, err := e.Decrypt(u.Phone)
	if err != nil {
		return
	}

	pwd, err := e.Decrypt(u.PlatformPwd)
	if err != nil {
		return
	}

	o = org.Organization{
		Id:                primitive.CreateIdentity(u.ID),
		Version:           u.Version,
		Desc:              primitive.CreateMSDDesc(u.Desc),
		Account:           primitive.CreateAccount(u.Name),
		Fullname:          primitive.CreateMSDFullname(u.Fullname),
		AvatarId:          primitive.CreateAvatarId(u.AvatarId),
		PlatformId:        u.PlatformId,
		CreatedAt:         u.CreatedAt.Unix(),
		UpdatedAt:         u.UpdatedAt.Unix(),
		PlatformPwd:       pwd,                                 // user only
		Email:             primitive.CreateEmail(email),        // user only
		Phone:             primitive.CreatePhoneNumber(phone),  // user only,
		AllowRequest:      u.AllowRequest,                      // org only
		DefaultRole:       u.DefaultRole,                       // org only
		OwnerTeamId:       u.OwnerTeamId,                       // org only
		ReadTeamId:        u.ReadTeamId,                        // org only
		WriteTeamId:       u.WriteTeamId,                       // org only
		ContributorTeamId: u.ContributorTeamId,                 // org only
		Owner:             primitive.CreateAccount(u.Owner),    // org only
		OwnerId:           primitive.CreateIdentity(u.OwnerId), // org only
		Website:           u.Website,                           // org only
		Type:              u.Type,
	}

	return
}

// TableName returns the table name for the UserDO struct in the database.
func (do *UserDO) TableName() string {
	return userTableName
}

// TokenDO represents the database model for the Token entity.
type TokenDO struct {
	postgresql.CommonModel

	Name       string `gorm:"column:name;index:name_index"`
	OwnerId    int64  `gorm:"column:owner_id;index:owner_index"`
	Owner      string `gorm:"column:owner;index:owner_index"`
	Expire     int64  `gorm:"column:expire"` // timeout in seconds
	Salt       string `gorm:"column:salt"`
	Token      string `gorm:"column:token"`
	LastEight  string `gorm:"column:last_eight;index:last_eight_index"`
	Version    int    `gorm:"column:version"`
	Permission string `gorm:"column:permission"`
}

// TableName returns the table name for the TokenDO struct in the database.
func (do *TokenDO) TableName() string {
	return tokenTableName
}

func toTokenDO(u *domain.PlatformToken) (do TokenDO) {
	do = TokenDO{
		Name:       u.Name.TokenName(),
		Owner:      u.Account.Account(),
		OwnerId:    u.OwnerId.Integer(),
		Expire:     u.Expire,
		LastEight:  u.LastEight,
		Salt:       u.Salt,
		Token:      u.Token,
		Version:    u.Version,
		Permission: u.Permission.TokenPerm(),
	}

	do.ID = u.Id.Integer()

	return
}

func (t *TokenDO) toToken() (token domain.PlatformToken) {
	token = domain.PlatformToken{
		Id:         primitive.CreateIdentity(t.ID),
		Version:    t.Version,
		Name:       primitive.CreateTokenName(t.Name),
		Account:    primitive.CreateAccount(t.Owner),
		OwnerId:    primitive.CreateIdentity(t.OwnerId),
		Expire:     t.Expire,
		LastEight:  t.LastEight,
		Salt:       t.Salt,
		Token:      t.Token,
		CreatedAt:  t.CreatedAt.Unix(),
		UpdatedAt:  t.UpdatedAt.Unix(),
		Permission: primitive.CreateTokenPerm(t.Permission),
	}

	return
}
