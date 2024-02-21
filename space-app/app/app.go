package app

import (
	commonapp "github.com/openmerlin/merlin-server/common/app"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/space-app/domain/repository"
	spacedomain "github.com/openmerlin/merlin-server/space/domain"
)

// SpaceappAppService
type SpaceappAppService interface {
	GetByName(primitive.Account, *spacedomain.SpaceIndex) (SpaceAppDTO, error)
}

// spaceRepository
type spaceRepository interface {
	FindByName(*spacedomain.SpaceIndex) (spacedomain.Space, error)
}

// permissionValidator
type permissionValidator interface {
	Check(primitive.Account, primitive.Account, primitive.ObjType, primitive.Action) error
}

// NewSpaceappAppService
func NewSpaceappAppService(
	repo repository.Repository,
	spaceRepo spaceRepository,
	permission commonapp.ResourcePermissionAppService,
) *spaceappAppService {
	return &spaceappAppService{
		repo:       repo,
		spaceRepo:  spaceRepo,
		permission: permission,
	}
}

// spaceappAppService
type spaceappAppService struct {
	repo       repository.Repository
	spaceRepo  spaceRepository
	permission commonapp.ResourcePermissionAppService
}

// GetByName
func (s *spaceappAppService) GetByName(
	user primitive.Account, index *spacedomain.SpaceIndex,
) (SpaceAppDTO, error) {
	var dto SpaceAppDTO

	space, err := s.spaceRepo.FindByName(index)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errorSpaceAppNotFound
		}

		return dto, err
	}

	if err = s.permission.CanRead(user, &space); err != nil {
		if allerror.IsNoPermission(err) {
			err = errorSpaceAppNotFound
		}

		return dto, err
	}

	// TODO it should find by newest commit

	app, err := s.repo.FindBySpaceId(space.Id)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errorSpaceAppNotFound
		}

		return dto, err
	}

	return toSpaceAppDTO(&app), nil
}
