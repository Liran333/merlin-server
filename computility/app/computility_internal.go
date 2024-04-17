/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package app

import (
	"github.com/sirupsen/logrus"

	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/computility/domain"
	"github.com/openmerlin/merlin-server/computility/domain/message"
	"github.com/openmerlin/merlin-server/computility/domain/repository"
	"github.com/openmerlin/merlin-server/utils"
)

// ComputilityInternalAppService is an interface for computility internal application service
type ComputilityInternalAppService interface {
	UserJoin(CmdToUserOrgOperate) error
	UserRemove(CmdToUserOrgOperate) error
	OrgDelete(CmdToOrgDelete) error
}

// NewComputilityInternalAppService creates a new instance of ComputilityInternalAppService
func NewComputilityInternalAppService(
	orgAdapter repository.ComputilityOrgRepositoryAdapter,
	detailAdapter repository.ComputilityDetailRepositoryAdapter,
	accountAdapter repository.ComputilityAccountRepositoryAdapter,
	messageAdapter message.ComputilityMessage,
) ComputilityInternalAppService {
	return &computilityInternalAppService{
		orgAdapter:     orgAdapter,
		detailAdapter:  detailAdapter,
		accountAdapter: accountAdapter,
		messageAdapter: messageAdapter,
	}
}

type computilityInternalAppService struct {
	orgAdapter     repository.ComputilityOrgRepositoryAdapter
	accountAdapter repository.ComputilityAccountRepositoryAdapter
	detailAdapter  repository.ComputilityDetailRepositoryAdapter
	messageAdapter message.ComputilityMessage
}

func (s *computilityInternalAppService) UserJoin(cmd CmdToUserOrgOperate) error {
	exist, err := s.orgAdapter.CheckOrgExist(cmd.OrgName)
	if err != nil {
		return err
	}
	if !exist {
		return nil
	}

	userExist, err := s.accountAdapter.CheckAccountExist(cmd.UserName)
	if err != nil {
		return err
	}

	org, err := s.orgAdapter.FindByOrgName(cmd.OrgName)
	if err != nil {
		return err
	}

	index := domain.ComputilityAccountIndex{
		UserName:    cmd.UserName,
		ComputeType: org.ComputeType,
	}

	if !userExist {
		err = s.accountAdapter.Add(&domain.ComputilityAccount{
			ComputilityAccountIndex: index,
			UsedQuota:               0,
			QuotaCount:              0,
			CreatedAt:               utils.Now(),
			Version:                 0,
		})
		if err != nil {
			return err
		}
	}

	account, err := s.accountAdapter.FindByAccountIndex(index)
	if err != nil {
		return err
	}
	err = s.accountAdapter.IncreaseAccountAssignedQuota(account, org.DefaultAssignQuota)
	if err != nil {
		return err
	}

	err = s.detailAdapter.Add(&domain.ComputilityDetail{
		ComputilityIndex: cmd.ComputilityIndex,
		QuotaCount:       org.DefaultAssignQuota,
		ComputeType:      org.ComputeType,
		Version:          0,
	})
	if err != nil {
		return err
	}

	err = s.orgAdapter.OrgAssignQuota(org, org.DefaultAssignQuota)
	if err != nil {
		return err
	}

	return err
}

func (s *computilityInternalAppService) UserRemove(cmd CmdToUserOrgOperate) error {
	org, err := s.orgAdapter.FindByOrgName(cmd.OrgName)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			return nil
		}

		return err
	}

	_, err = s.accountAdapter.FindByAccountIndex(
		domain.ComputilityAccountIndex{
			UserName:    cmd.UserName,
			ComputeType: org.ComputeType,
		})
	if err != nil {
		logrus.Errorf("user: %s have no computility account", cmd.UserName.Account())
		return err
	}

	err = s.userRemoveOperate(&cmd.ComputilityIndex)

	return err
}

func (s *computilityInternalAppService) OrgDelete(cmd CmdToOrgDelete) error {
	exist, err := s.orgAdapter.CheckOrgExist(cmd.OrgName)
	if err != nil {
		return err
	}
	if !exist {
		return nil
	}

	r, err := s.detailAdapter.GetMembers(cmd.OrgName)
	if err != nil {
		return err
	}

	for _, v := range r {
		err = s.userRemoveOperate(&domain.ComputilityIndex{
			OrgName:  v.OrgName,
			UserName: v.UserName,
		})
		if err != nil {
			return err
		}
	}

	org, err := s.orgAdapter.FindByOrgName(cmd.OrgName)
	if err != nil {
		return err
	}

	err = s.orgAdapter.Delete(org.Id)

	return err
}

func (s *computilityInternalAppService) userRemoveOperate(index *domain.ComputilityIndex) error {
	detail, err := s.detailAdapter.FindByIndex(index)
	if err != nil {
		return err
	}

	assigned := detail.QuotaCount

	accountIndex := domain.ComputilityAccountIndex{
		UserName:    index.UserName,
		ComputeType: detail.ComputeType,
	}

	account, err := s.accountAdapter.FindByAccountIndex(accountIndex)
	if err != nil {
		return err
	}
	balance := account.QuotaCount - account.UsedQuota

	if assigned > balance {
		infoList := domain.RecallInfoList{}
		infoList.InfoList = append(infoList.InfoList, domain.RecallInfo{
			UserName:    index.UserName,
			QuotaCount:  assigned - balance,
			ComputeType: detail.ComputeType,
		})

		logrus.Infof("user:%s quota recall, balance:%v, quota need to return:%v", index.UserName.Account(), balance, assigned)

		e := domain.NewcomputeRecallEvent(&infoList)
		err = s.messageAdapter.SendComputilityRecallEvent(&e)
		if err != nil {
			logrus.Errorf("publish topic computility_recalled failed, user:%s, quota:%v", index.UserName.Account(), assigned-balance)
		} else {
			logrus.Infof("publish topic computility_recalled success, user:%s, quota:%v", index.UserName.Account(), assigned-balance)
		}
	}

	err = s.accountAdapter.DecreaseAccountAssignedQuota(account, assigned)
	if err != nil {
		return err
	}

	org, err := s.orgAdapter.FindByOrgName(index.OrgName)
	if err != nil {
		return err
	}

	err = s.orgAdapter.OrgRecallQuota(org, assigned)
	if err != nil {
		return err
	}

	err = s.detailAdapter.Delete(detail.Id)
	if err != nil {
		return err
	}

	err = s.accountAdapter.CancelAccount(accountIndex)

	return err
}
