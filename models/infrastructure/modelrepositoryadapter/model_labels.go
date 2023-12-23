package modelrepositoryadapter

import (
	"errors"

	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/models/domain"
)

type modelLabelsAdapter struct {
	daoImpl
}

func (adapter *modelLabelsAdapter) Save(index *domain.ModelIndex, labels *domain.ModelLabels) error {
	do := toLabelsDO(labels)

	v := adapter.db().Model(&modelDO{}).Where(
		equalQuery(fieldOwner), index.Owner.Account(),
	).Where(
		equalQuery(fieldName), index.Name.MSDName(),
	).Select(
		fieldTask, fieldOthers, fieldFrameworks,
	).Updates(&do)

	if v.Error != nil {
		return v.Error
	}

	if v.RowsAffected == 0 {
		return commonrepo.NewErrorResourceNotExists(
			errors.New("not found"),
		)
	}

	return nil
}
