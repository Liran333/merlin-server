/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package repositoryadapter

import (
	"errors"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/computility/domain"
)

type computilityOrgAdapter struct {
	daoImpl
}

// FindByOrgName finds a computility org in the repository based on the org name and returns an error if any occurs.
func (adapter *computilityOrgAdapter) FindByOrgName(name primitive.Account) (
	domain.ComputilityOrg, error,
) {
	do := computilityOrgDO{OrgName: name.Account()}

	result := computilityOrgDO{}
	if err := adapter.daoImpl.GetRecord(&do, &result); err != nil {
		return domain.ComputilityOrg{}, err
	}

	return result.toComputilityOrg(), nil
}

// Delete deletes a computility org in the database and returns an error if any occurs.
func (adapter *computilityOrgAdapter) Delete(id primitive.Identity) error {
	return adapter.DeleteByPrimaryKey(
		&computilityOrgDO{Id: id.Integer()},
	)
}

// OrgAssignQuota increases used_quota field in computility org in the database and returns an error if any occurs.
func (adapter *computilityOrgAdapter) OrgAssignQuota(
	org domain.ComputilityOrg, quota int,
) error {
	do := toComputilityOrgDO(&org)

	do.Version += 1
	do.UsedQuota = do.UsedQuota + quota

	result := adapter.db().Model(
		&computilityOrgDO{Id: do.Id},
	).Where(
		equalQuery(filedVersion), org.Version,
	).Select(`*`).Omit(fieldQuotaCount).Updates(&do)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return repository.NewErrorResourceNotExists(errors.New("resource not found"))
	}

	return nil
}

// OrgRecallQuota decreases used_quota field in computility org in the database and returns an error if any occurs.
func (adapter *computilityOrgAdapter) OrgRecallQuota(
	org domain.ComputilityOrg, quota int,
) error {
	do := toComputilityOrgDO(&org)

	do.Version += 1
	do.UsedQuota = do.UsedQuota - quota

	result := adapter.db().Model(
		&computilityOrgDO{OrgName: do.OrgName},
	).Where(
		equalQuery(filedVersion), org.Version,
	).Select(`*`).Omit(fieldQuotaCount).Updates(&do)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return repository.NewErrorResourceNotExists(errors.New("resource not found"))
	}

	return nil
}
