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
	Index      domain.ComputilityAccountRecordIndex
	QuotaCount int
}

type CmdToSupplyRecord struct {
	Index      domain.ComputilityAccountRecordIndex
	QuotaCount int
	NewSpaceId primitive.Identity
}

type AccountQuotaDetailDTO struct {
	UserName     string `json:"user_name"`
	UsedQuota    int    `json:"used_quota"`
	TotalQuota   int    `json:"total_quota"`
	ComputeType  string `json:"compute_type"`
	QuotaBalance int    `json:"quota_balance"`
}

func toAccountQuotaDetailDTO(a *domain.ComputilityAccount) AccountQuotaDetailDTO {
	return AccountQuotaDetailDTO{
		UserName:     a.UserName.Account(),
		UsedQuota:    a.UsedQuota,
		TotalQuota:   a.QuotaCount,
		QuotaBalance: a.QuotaCount - a.UsedQuota,
		ComputeType:  a.ComputeType.ComputilityType(),
	}
}

type AccountRecordlDTO struct {
	UserName    string `json:"user_name"`
	SpaceId     string `json:"space_id"`
	QuotaCount  int    `json:"quota_count"`
	ComputeType string `json:"compute_type"`
}

func toAccountRecordlDTO(a *domain.ComputilityAccountRecord) AccountRecordlDTO {
	return AccountRecordlDTO{
		UserName:    a.UserName.Account(),
		SpaceId:     a.SpaceId.Identity(),
		QuotaCount:  a.QuotaCount,
		ComputeType: a.ComputeType.ComputilityType(),
	}
}

type QuotaRecallDTO struct {
	UserName  string              `json:"user_name"`
	Records   []AccountRecordlDTO `json:"records"`
	QuotaDebt int                 `json:"quota_debt"`
}

func toQuotaRecallDTO(user primitive.Account, a []domain.ComputilityAccountRecord, debt int) QuotaRecallDTO {
	var records []AccountRecordlDTO

	for i := range a {
		info := toAccountRecordlDTO(&a[i])

		records = append(records, info)
	}

	return QuotaRecallDTO{
		UserName:  user.Account(),
		Records:   records,
		QuotaDebt: debt,
	}
}
