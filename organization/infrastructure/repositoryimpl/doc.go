/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package repositoryimpl

import (
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
)

const (
	fieldName     = "name"
	fieldOwner    = "owner"
	fieldCount    = "count"
	fieldAccount  = "account"
	fieldBio      = "bio"
	fieldAvatarId = "avatar_id"
	fieldVersion  = "version"
	fieldType     = "type"
	fieldUser     = "user_name"
	fieldOrg      = "org_name"
	fieldRole     = "role"
	fieldInvitee  = "user_name"
	fieldInviter  = "inviter"
	fieldStatus   = "status"
)

// Member represents a member in the database.
type Member struct {
	postgresql.CommonModel
	Username string `gorm:"column:user_name;index:username_index"`
	UserId   int64  `gorm:"column:user_id;index:userid_index"`
	Orgname  string `gorm:"column:org_name;index:orgname_index"`
	OrgId    int64  `gorm:"column:org_id;index:orgid_index"`
	Role     string `gorm:"column:role"`
	Type     string `gorm:"column:type"`
	Version  int    `gorm:"column:version"`
}

// Approve both request and approve use the same DO
type Approve struct {
	postgresql.CommonModel

	Username string `gorm:"column:user_name;index:username_index"`
	UserId   int64  `gorm:"column:user_id;index:userid_index"`
	Orgname  string `gorm:"column:org_name;index:orgname_index"`
	OrgId    int64  `gorm:"column:org_id;index:orgid_index"`
	Role     string `gorm:"column:role"`
	Expire   int64  `gorm:"column:expire"`  // approve only attr
	Inviter  string `gorm:"column:inviter"` // aprove only attr
	Status   string `gorm:"column:status;index:status_index"`
	Type     string `gorm:"column:type"`
	By       string `gorm:"column:by"`
	Msg      string `gorm:"column:msg"`
	Version  int    `gorm:"column:version"`
}
