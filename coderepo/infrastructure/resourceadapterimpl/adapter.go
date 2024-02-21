package resourceadapterimpl

import (
	"errors"

	"github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/coderepo/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	modelrepo "github.com/openmerlin/merlin-server/models/domain/repository"
	spacerepo "github.com/openmerlin/merlin-server/space/domain/repository"
)

func NewResourceAdapterImpl(
	model modelrepo.ModelRepositoryAdapter,
	space spacerepo.SpaceRepositoryAdapter,
) *resourceAdapterImpl {
	return &resourceAdapterImpl{
		model: model,
		space: space,
	}
}

// resourceAdapterImpl
type resourceAdapterImpl struct {
	model modelrepo.ModelRepositoryAdapter
	space spacerepo.SpaceRepositoryAdapter
}

func (adapter *resourceAdapterImpl) GetByName(index *domain.CodeRepoIndex) (domain.Resource, error) {
	r, err := adapter.model.FindByName(index)
	if err == nil {
		return &r, nil
	}
	if !commonrepo.IsErrorResourceNotExists(err) {
		return nil, err
	}

	space, err := adapter.space.FindByName(index)

	return &space, err
}

func (adapter *resourceAdapterImpl) GetByType(t primitive.RepoType, index *domain.CodeRepoIndex) (domain.Resource, error) {
	if t.IsModel() {
		r, err := adapter.model.FindByName(index)

		return &r, err
	}

	if t.IsSpace() {
		r, err := adapter.space.FindByName(index)

		return &r, err

	}

	return nil, commonrepo.NewErrorResourceNotExists(errors.New("unknown repo type"))
}
