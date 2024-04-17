/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package app

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/computility/domain"
)

// CmdToUserJoin is a struct used for user join computility.
type CmdToUserOrgOperate struct {
	domain.ComputilityIndex
}

type CmdToOrgDelete struct {
	OrgName primitive.Account
}

type CmdToUserQuotaUpdate struct {
	domain.ComputilityAccountIndex
	QuotaCount int
}
