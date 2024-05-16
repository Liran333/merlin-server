/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package repository provides interfaces for managing approvals in an organization.
package repository

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/organization/domain"
)

// Approve is an interface that defines the methods for handling approval-related operations.
type Approve interface {
	AddInvite(*domain.Approve) (domain.Approve, error)
	SaveInvite(*domain.Approve) (domain.Approve, error)
	AddRequest(*domain.MemberRequest) (domain.MemberRequest, error)
	SaveRequest(*domain.MemberRequest) (domain.MemberRequest, error)
	DeleteInviteAndReqByOrg(primitive.Account) error
	Count(primitive.Account) (int64, error)
	// DeleteRequestByOrg(primitive.Account) error
	ListInvitation(*domain.OrgInvitationListCmd) ([]domain.Approve, error)
	ListRequests(*domain.OrgMemberReqListCmd) ([]domain.MemberRequest, error)
	UpdateAllApproveStatus(primitive.Account, primitive.Account, domain.ApproveStatus) error
}
