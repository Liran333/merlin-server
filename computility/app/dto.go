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

type AccountQuotaDetailDTO struct {
	UserName     string `json:"user_name"`
	UsedQuota    int    `json:"used_quota"`
	TotalQuota   int    `json:"total_quota"`
	ComputeType  string `json:"compute_type"`
	QuotaBalance int    `json:"quota_balance"`
}

func toAccountDTO(a *domain.ComputilityAccount) AccountQuotaDetailDTO {
	return AccountQuotaDetailDTO{
		UserName:     a.UserName.Account(),
		UsedQuota:    a.UsedQuota,
		TotalQuota:   a.QuotaCount,
		QuotaBalance: a.QuotaCount - a.UsedQuota,
		ComputeType:  a.ComputeType.ComputilityType(),
	}
}
