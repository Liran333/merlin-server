package app

import (
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
	permission permissionValidator,
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
	permission permissionValidator
}

// GetByName
func (s *spaceappAppService) GetByName(
	user primitive.Account, index *spacedomain.SpaceIndex,
) (SpaceAppDTO, error) {
	var dto SpaceAppDTO

	spaceId, err := s.checkPermission(user, index, primitive.ActionRead)
	if err != nil {
		return dto, err
	}

	// TODO it should find by newest commit
	app, err := s.repo.FindBySpaceId(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errorSpaceAppNotFound
		}

		return dto, err
	}

	return toSpaceAppDTO(&app), nil
}

func (s *spaceappAppService) checkPermission(
	user primitive.Account, index *spacedomain.SpaceIndex, action primitive.Action,
) (primitive.Identity, error) {
	space, err := s.spaceRepo.FindByName(index)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errorSpaceAppNotFound
		}

		return nil, err
	}

	if space.IsPublic() {
		return space.Id, nil
	}

	// can't access private app anonymously
	if user == nil {
		return nil, errorSpaceAppNotFound
	}

	err = s.hasPermission(user, &space, action)

	return space.Id, err
}

func (s *spaceappAppService) hasPermission(
	user primitive.Account, space *spacedomain.Space, action primitive.Action,
) error {
	if space.OwnedBy(user) {
		return nil
	}

	if space.OwnedByPerson() {
		return errorSpaceAppNotFound
	}

	if err := s.permission.Check(user, space.Owner, primitive.ObjTypeModel, action); err != nil {
		return errorSpaceAppNotFound
	}

	return nil
}
