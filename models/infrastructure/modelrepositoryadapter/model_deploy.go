package modelrepositoryadapter

import "github.com/openmerlin/merlin-server/models/domain"

type modelDeployAdapter struct {
	daoImpl
}

func (adapter *modelDeployAdapter) Create(index domain.ModelIndex, deploys []domain.Deploy) error {
	var dos []*modelDeployDO
	for _, v := range deploys {
		do := toModelDeployDO(&index, v)
		dos = append(dos, &do)
	}

	return adapter.db().Create(dos).Error
}

func (adapter *modelDeployAdapter) DeleteByOwnerName(index domain.ModelIndex) error {
	do := modelDeployDO{
		Owner: index.Owner.Account(),
		Name:  index.Name.MSDName(),
	}

	return adapter.db().Where(&do).Delete(&do).Error
}

func (adapter *modelDeployAdapter) FindByOwnerName(index *domain.ModelIndex) ([]domain.Deploy, error) {
	do := modelDeployDO{
		Owner: index.Owner.Account(),
		Name:  index.Name.MSDName(),
	}

	var dos []modelDeployDO
	err := adapter.db().Where(&do).Find(&dos).Error
	if err != nil {
		return nil, err
	}

	var deploys []domain.Deploy
	for _, v := range dos {
		deploys = append(deploys, v.toDeploy())
	}

	return deploys, err
}
