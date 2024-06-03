/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package app

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"

	"github.com/openmerlin/merlin-server/common/domain/allerror"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/computility/domain"
	"github.com/openmerlin/merlin-server/computility/domain/message"
	"github.com/openmerlin/merlin-server/computility/domain/repository"
	orgapp "github.com/openmerlin/merlin-server/organization/app"
	"github.com/openmerlin/merlin-server/utils"
)

// ComputilityInternalAppService is an interface for computility internal application service
type ComputilityInternalAppService interface {
	UserJoin(CmdToUserOrgOperate) error
	UserRemove(CmdToUserOrgOperate) (QuotaRecallDTO, error)
	OrgDelete(CmdToOrgDelete) ([]QuotaRecallDTO, error)

	UserQuotaRelease(CmdToUserQuotaUpdate) error
	UserQuotaConsume(CmdToUserQuotaUpdate) error
	SpaceCreateSupply(CmdToSupplyRecord) error
}

// NewComputilityInternalAppService creates a new instance of ComputilityInternalAppService
func NewComputilityInternalAppService(
	orgAdapter repository.ComputilityOrgRepositoryAdapter,
	detailAdapter repository.ComputilityDetailRepositoryAdapter,
	accountAdapter repository.ComputilityAccountRepositoryAdapter,
	accountRecordAtapter repository.ComputilityAccountRecordRepositoryAdapter,
	messageAdapter message.ComputilityMessage,
	privilege orgapp.PrivilegeOrg,
) ComputilityInternalAppService {
	return &computilityInternalAppService{
		orgAdapter:           orgAdapter,
		detailAdapter:        detailAdapter,
		accountAdapter:       accountAdapter,
		accountRecordAtapter: accountRecordAtapter,
		messageAdapter:       messageAdapter,
		privilege:            privilege,
	}
}

type computilityInternalAppService struct {
	orgAdapter           repository.ComputilityOrgRepositoryAdapter
	accountAdapter       repository.ComputilityAccountRepositoryAdapter
	detailAdapter        repository.ComputilityDetailRepositoryAdapter
	accountRecordAtapter repository.ComputilityAccountRecordRepositoryAdapter
	messageAdapter       message.ComputilityMessage
	privilege            orgapp.PrivilegeOrg
}

func (s *computilityInternalAppService) UserJoin(cmd CmdToUserOrgOperate) error {
	if s.privilege == nil {
		return nil
	}

	org, err := s.orgAdapter.FindByOrgName(cmd.OrgName)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			return nil
		}

		return err
	}

	if org.UsedQuota >= org.QuotaCount {
		logrus.Errorf("quota assign failed | organization:%s has no balance to assign to user:%s",
			org.OrgName.Account(), cmd.UserName.Account(),
		)

		e := xerrors.Errorf("organization:%s has no balance to assign to user:%s",
			org.OrgName.Account(), cmd.UserName.Account(),
		)

		return allerror.New(
			allerror.ErrorCodeInsufficientQuota,
			"organization insufficient computing quota balance", e)
	}

	userExist, err := s.accountAdapter.CheckAccountExist(cmd.UserName)
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
		CreatedAt:        utils.Now(),
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

func (s *computilityInternalAppService) UserRemove(cmd CmdToUserOrgOperate) (
	QuotaRecallDTO, error,
) {
	if s.privilege == nil {
		return QuotaRecallDTO{}, nil
	}

	org, err := s.orgAdapter.FindByOrgName(cmd.OrgName)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			return QuotaRecallDTO{}, nil
		}

		return QuotaRecallDTO{}, err
	}

	_, err = s.accountAdapter.FindByAccountIndex(
		domain.ComputilityAccountIndex{
			UserName:    cmd.UserName,
			ComputeType: org.ComputeType,
		})
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			logrus.Errorf("user: %s have no computility account", cmd.UserName.Account())

			return QuotaRecallDTO{}, nil
		}

		return QuotaRecallDTO{}, err
	}

	recall, err := s.userRemoveOperate(&cmd.ComputilityIndex)

	return recall, err
}

func (s *computilityInternalAppService) OrgDelete(cmd CmdToOrgDelete) (
	[]QuotaRecallDTO, error,
) {
	if s.privilege == nil {
		return nil, nil
	}

	org, err := s.orgAdapter.FindByOrgName(cmd.OrgName)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			return nil, nil
		}

		return nil, err
	}

	r, err := s.detailAdapter.GetMembers(cmd.OrgName)
	if err != nil {
		return nil, err
	}

	rList := make([]QuotaRecallDTO, 0)
	for _, v := range r {
		s, err := s.userRemoveOperate(&domain.ComputilityIndex{
			OrgName:  v.OrgName,
			UserName: v.UserName,
		})

		if err != nil {
			logrus.Errorf("org deleted | computility remove user:%s error: %s", v.UserName.Account(), err)

			continue
		}

		rList = append(rList, s)
	}

	err = s.orgAdapter.Delete(org.Id)

	return rList, err
}

func (s *computilityInternalAppService) userRemoveOperate(index *domain.ComputilityIndex) (
	QuotaRecallDTO, error,
) {
	if s.privilege == nil {
		return QuotaRecallDTO{}, nil
	}

	detail, err := s.detailAdapter.FindByIndex(index)
	if err != nil {
		return QuotaRecallDTO{}, err
	}

	assigned := detail.QuotaCount

	accountIndex := domain.ComputilityAccountIndex{
		UserName:    index.UserName,
		ComputeType: detail.ComputeType,
	}

	account, err := s.accountAdapter.FindByAccountIndex(accountIndex)
	if err != nil {
		return QuotaRecallDTO{}, err
	}

	balance := account.QuotaCount - account.UsedQuota

	var recall QuotaRecallDTO
	if assigned > balance {
		debt := assigned - balance

		spaces, n, err := s.accountRecordAtapter.ListByAccountIndex(accountIndex)
		if err != nil {
			logrus.Errorf("find user bind space failed, %s", err)

			return QuotaRecallDTO{}, err
		}

		if n == 0 {
			logrus.Errorf("cannot find user:%s bind space, user debt: %v", account.UserName.Account(), debt)
		} else {
			recall = toQuotaRecallDTO(index.UserName, spaces, debt)
		}
	}

	err = s.accountAdapter.DecreaseAccountAssignedQuota(account, assigned)
	if err != nil {
		return QuotaRecallDTO{}, err
	}

	org, err := s.orgAdapter.FindByOrgName(index.OrgName)
	if err != nil {
		return QuotaRecallDTO{}, err
	}

	err = s.orgAdapter.OrgRecallQuota(org, assigned)
	if err != nil {
		return QuotaRecallDTO{}, err
	}

	err = s.detailAdapter.Delete(detail.Id)
	if err != nil {
		return QuotaRecallDTO{}, err
	}

	err = s.accountAdapter.CancelAccount(accountIndex)

	return recall, err
}

func (s *computilityInternalAppService) UserQuotaRelease(cmd CmdToUserQuotaUpdate) error {
	if cmd.Index.ComputeType.IsCpu() {
		return nil
	}

	if s.privilege == nil {
		return nil
	}

	record, err := s.accountRecordAtapter.FindByRecordIndex(cmd.Index)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			logrus.Errorf("user:%s has not cosume record to release", cmd.Index.UserName.Account())

			return nil
		}

		return err
	}

	accountIndex := domain.ComputilityAccountIndex{
		UserName:    cmd.Index.UserName,
		ComputeType: cmd.Index.ComputeType,
	}

	account, err := s.accountAdapter.FindByAccountIndex(accountIndex)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			logrus.Errorf("user:%s is not a computility account, can not release quota", cmd.Index.UserName.Account())
			return nil
		}
		return err
	}

	if account.UsedQuota == 0 {
		logrus.Errorf("user:%s has no quota to release", cmd.Index.UserName.Account())
		return nil
	}

	err = s.accountAdapter.ReleaseQuota(account, cmd.QuotaCount)
	if err != nil {
		return err
	}

	err = s.accountRecordAtapter.Delete(record.Id)
	if err != nil {
		logrus.Errorf("delete user:%s account record failed, %s", cmd.Index.UserName.Account(), err)

		return err
	}

	err = s.accountAdapter.CancelAccount(accountIndex)
	if err != nil {
		logrus.Errorf("cancel user:%s account failed, %s", cmd.Index.UserName.Account(), err)

		return err
	}

	return nil
}

func (s *computilityInternalAppService) UserQuotaConsume(cmd CmdToUserQuotaUpdate) error {
	if cmd.Index.ComputeType.IsCpu() {
		return nil
	}

	if s.privilege == nil {
		return nil
	}

	user := cmd.Index.UserName
	_, err := s.accountRecordAtapter.FindByRecordIndex(cmd.Index)
	if err == nil {
		logrus.Errorf("user:%s already bind space:%s", user, cmd.Index.SpaceId.Identity())

		return nil
	}

	b, err := s.accountAdapter.CheckAccountExist(user)
	if err != nil {
		return err
	}
	if !b {
		e := xerrors.Errorf("user %s no quota banlance for %s",
			user.Account(), cmd.Index.ComputeType.ComputilityType())

		logrus.Errorf("consume quota error| %s", e)

		return allerror.New(
			allerror.ErrorCodeNoNpuPermission,
			"no quota balance", e)
	}

	index := domain.ComputilityAccountIndex{
		UserName:    user,
		ComputeType: cmd.Index.ComputeType,
	}

	account, err := s.accountAdapter.FindByAccountIndex(index)
	if err != nil {
		logrus.Errorf("find user:%s account failed, %s", cmd.Index.UserName.Account(), err)

		return err
	}

	balance := account.QuotaCount - account.UsedQuota
	if balance < 1 {
		e := xerrors.Errorf("user %s insufficient computing quota balance", user.Account())

		logrus.Errorf("consume quota error| %s", e)

		return allerror.New(
			allerror.ErrorCodeInsufficientQuota,
			"insufficient computing quota balance", e)
	}

	err = s.accountAdapter.ConsumeQuota(account, cmd.QuotaCount)
	if err != nil {
		return err
	}

	err = s.accountRecordAtapter.Add(&domain.ComputilityAccountRecord{
		ComputilityAccountRecordIndex: cmd.Index,
		CreatedAt:                     utils.Now(),
		QuotaCount:                    cmd.QuotaCount,
		Version:                       0,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *computilityInternalAppService) SpaceCreateSupply(cmd CmdToSupplyRecord) error {
	if cmd.Index.ComputeType.IsCpu() {
		return nil
	}

	if s.privilege == nil {
		return nil
	}

	b, err := s.accountAdapter.CheckAccountExist(cmd.Index.UserName)
	if err != nil {
		return err
	}
	if !b {
		e := xerrors.Errorf("user %s no permission for npu space", cmd.Index.UserName)

		return allerror.New(
			allerror.ErrorCodeNoNpuPermission,
			"no permission for npu space", e)
	}

	record, err := s.accountRecordAtapter.FindByRecordIndex(cmd.Index)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			logrus.Errorf("user %s has not cosume record to release", cmd.Index.UserName.Account())
			return nil
		}
		return err
	}

	record.SpaceId = cmd.NewSpaceId

	err = s.accountRecordAtapter.Save(&record)
	if err != nil {
		logrus.Errorf("user %s no permission for %s space", cmd.Index.UserName, cmd.Index.ComputeType.ComputilityType())
	}

	return nil
}
