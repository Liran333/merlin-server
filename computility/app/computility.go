/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package app

import (
	"fmt"

	"github.com/openmerlin/merlin-server/common/domain/allerror"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/computility/domain"
	"github.com/openmerlin/merlin-server/computility/domain/repository"
	"github.com/sirupsen/logrus"
)

// ComputilityAppService is an interface for computility internal application service
type ComputilityAppService interface {
	UserQuotaRelease(CmdToUserQuotaUpdate) error
	UserQuotaConsume(CmdToUserQuotaUpdate) error
}

// NewComputilityAppService creates a new instance of ComputilityAppService
func NewComputilityAppService(
	orgAdapter repository.ComputilityOrgRepositoryAdapter,
	detailAdapter repository.ComputilityDetailRepositoryAdapter,
	accountAdapter repository.ComputilityAccountRepositoryAdapter,
) ComputilityAppService {
	return &computilityAppService{
		orgAdapter:     orgAdapter,
		detailAdapter:  detailAdapter,
		accountAdapter: accountAdapter,
	}
}

type computilityAppService struct {
	orgAdapter     repository.ComputilityOrgRepositoryAdapter
	accountAdapter repository.ComputilityAccountRepositoryAdapter
	detailAdapter  repository.ComputilityDetailRepositoryAdapter
}

func (s *computilityAppService) UserQuotaRelease(cmd CmdToUserQuotaUpdate) error {
	user := cmd.UserName

	accountIndex := domain.ComputilityAccountIndex{
		UserName:    user,
		ComputeType: cmd.ComputeType,
	}

	account, err := s.accountAdapter.FindByAccountIndex(accountIndex)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			logrus.Errorf("user %s is not a computility account, can not release quota", user)
			return nil
		}
		return err
	}

	if account.UsedQuota == 0 {
		e := fmt.Errorf("user %s does not occupy any quota", user.Account())

		return allerror.New(
			allerror.ErrorCodeNoUsedQuota,
			"user does not occupy any quota", e)
	}

	err = s.accountAdapter.ReleaseQuota(account, cmd.QuotaCount)
	if err != nil {
		return err
	}

	err = s.accountAdapter.CancelAccount(accountIndex)

	return err
}

func (s *computilityAppService) UserQuotaConsume(cmd CmdToUserQuotaUpdate) error {
	user := cmd.UserName

	b, err := s.accountAdapter.CheckAccountExist(user)
	if err != nil {
		return err
	}
	if !b {
		e := fmt.Errorf("user %s no permission for npu space", user.Account())

		return allerror.New(
			allerror.ErrorCodeNoNpuPermission,
			"no permission for npu space", e)
	}

	index := domain.ComputilityAccountIndex{
		UserName:    user,
		ComputeType: cmd.ComputeType,
	}

	account, err := s.accountAdapter.FindByAccountIndex(index)
	if err != nil {
		return err
	}

	balance := account.QuotaCount - account.UsedQuota
	if balance < 1 {
		e := fmt.Errorf("user %s insufficient computing quota balance", user.Account())

		return allerror.New(
			allerror.ErrorCodeInsufficientQuota,
			"insufficient computing quota balance", e)
	}

	err = s.accountAdapter.ConsumeQuota(account, cmd.QuotaCount)
	if err != nil {
		return err
	}

	return nil
}
