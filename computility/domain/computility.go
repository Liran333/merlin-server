/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package domain

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

type ComputilityDetail struct {
	ComputilityIndex

	Id          primitive.Identity
	CreatedAt   int64
	QuotaCount  int
	ComputeType primitive.ComputilityType

	Version int
}

type ComputilityAccount struct {
	ComputilityAccountIndex

	Id         primitive.Identity
	UsedQuota  int
	QuotaCount int
	CreatedAt  int64

	Version int
}

type ComputilityOrg struct {
	Id                 primitive.Identity
	OrgId              primitive.Identity
	OrgName            primitive.Account
	UsedQuota          int
	QuotaCount         int
	ComputeType        primitive.ComputilityType
	DefaultAssignQuota int

	Version int
}

type ComputilityIndex struct {
	OrgName  primitive.Account
	UserName primitive.Account
}

type ComputilityAccountIndex struct {
	UserName    primitive.Account
	ComputeType primitive.ComputilityType
}

type RecallInfoList struct {
	InfoList []RecallInfo
}

type RecallInfo struct {
	UserName    primitive.Account
	QuotaCount  int
	ComputeType primitive.ComputilityType
}

type ComputilityAccountRecordIndex struct {
	UserName    primitive.Account
	SpaceId     primitive.Identity
	ComputeType primitive.ComputilityType
}

type ComputilityAccountRecord struct {
	ComputilityAccountRecordIndex

	Id         primitive.Identity
	CreatedAt  int64
	QuotaCount int

	Version int
}
