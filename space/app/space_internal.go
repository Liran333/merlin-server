package app

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/space/domain/repository"
)

type SpaceInternalAppService interface {
	GetById(primitive.Identity) (SpaceMetaDTO, error)
}

func NewSpaceInternalAppService(
	repoAdapter repository.SpaceRepositoryAdapter,
) SpaceInternalAppService {
	return &spaceInternalAppService{
		repoAdapter: repoAdapter,
	}
}

type spaceInternalAppService struct {
	repoAdapter repository.SpaceRepositoryAdapter
}

func (s *spaceInternalAppService) GetById(spaceId primitive.Identity) (SpaceMetaDTO, error) {
	space, err := s.repoAdapter.FindById(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errorSpaceNotFound
		}

		return SpaceMetaDTO{}, err
	}

	return toSpaceMetaDTO(&space), nil
}
