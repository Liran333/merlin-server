package modelrepositoryadapter

import (
	"errors"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/models/domain"
)

type modelLabelsAdapter struct {
	daoImpl
}

func (adapter *modelLabelsAdapter) Save(modelId primitive.Identity, labels *domain.ModelLabels) error {
	do := toLabelsDO(labels)

	v := adapter.db().Model(
		&modelDO{Id: modelId.Integer()},
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
