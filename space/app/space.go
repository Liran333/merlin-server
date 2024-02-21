package app

import (
	"fmt"

	"github.com/sirupsen/logrus"

	coderepoapp "github.com/openmerlin/merlin-server/coderepo/app"
	commonapp "github.com/openmerlin/merlin-server/common/app"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/space/domain"
	"github.com/openmerlin/merlin-server/space/domain/message"
	"github.com/openmerlin/merlin-server/space/domain/repository"
	"github.com/openmerlin/merlin-server/utils"
)

var (
	errorSpaceNotFound      = allerror.NewNotFound(allerror.ErrorCodeSpaceNotFound, "not found")
	errorSpaceCountExceeded = allerror.NewCountExceeded("space count exceed")
)

type Permission interface {
	Check(primitive.Account, primitive.Account, primitive.ObjType, primitive.Action) error
}

type SpaceAppService interface {
	Create(primitive.Account, *CmdToCreateSpace) (string, error)
	Delete(primitive.Account, primitive.Identity) (string, error)
	Update(primitive.Account, primitive.Identity, *CmdToUpdateSpace) (string, error)
	GetByName(primitive.Account, *domain.SpaceIndex) (SpaceDTO, error)
	List(primitive.Account, *CmdToListSpaces) (SpacesDTO, error)
}

func NewSpaceAppService(
	permission commonapp.ResourcePermissionAppService,
	msgAdapter message.SpaceMessage,
	codeRepoApp coderepoapp.CodeRepoAppService,
	repoAdapter repository.SpaceRepositoryAdapter,
) SpaceAppService {
	return &spaceAppService{
		permission:  permission,
		msgAdapter:  msgAdapter,
		codeRepoApp: codeRepoApp,
		repoAdapter: repoAdapter,
	}
}

type spaceAppService struct {
	permission  commonapp.ResourcePermissionAppService
	msgAdapter  message.SpaceMessage
	codeRepoApp coderepoapp.CodeRepoAppService
	repoAdapter repository.SpaceRepositoryAdapter
}

func (s *spaceAppService) Create(user primitive.Account, cmd *CmdToCreateSpace) (string, error) {
	err := s.permission.CanCreate(user, cmd.Owner, primitive.ObjTypeSpace)
	if err != nil {
		return "", err
	}

	if err := s.spaceCountCheck(cmd.Owner); err != nil {
		return "", err
	}

	coderepo, err := s.codeRepoApp.Create(&cmd.CmdToCreateRepo)
	if err != nil {
		return "", err
	}

	now := utils.Now()

	err = s.repoAdapter.Add(&domain.Space{
		SDK:       cmd.SDK,
		Desc:      cmd.Desc,
		Hardware:  cmd.Hardware,
		Fullname:  cmd.Fullname,
		CodeRepo:  coderepo,
		CreatedBy: user,
		CreatedAt: now,
		UpdatedAt: now,
	})

	return coderepo.Id.Identity(), err

	// TODO send space created event in order to add activity and operation log
}

func (s *spaceAppService) Delete(user primitive.Account, spaceId primitive.Identity) (action string, err error) {
	space, err := s.repoAdapter.FindById(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		}

		return
	}

	action = fmt.Sprintf(
		"delete space of %s:%s/%s",
		spaceId.Identity(), space.Owner.Account(), space.Name.MSDName(),
	)

	if err = s.permission.CanDelete(user, &space); err != nil {
		if allerror.IsNoPermission(err) {
			err = errorSpaceNotFound
		}

		return
	}

	if err = s.codeRepoApp.Delete(space.RepoIndex()); err != nil {
		return
	}

	if err = s.repoAdapter.Delete(space.Id); err != nil {
		return
	}

	e := domain.NewSpaceDeletedEvent(&space)
	if err1 := s.msgAdapter.SendSpaceDeletedEvent(&e); err1 != nil {
		logrus.Errorf("failed to send space deleted event, space id:%s", spaceId.Identity())
	}

	return
}

func (s *spaceAppService) Update(
	user primitive.Account, spaceId primitive.Identity, cmd *CmdToUpdateSpace,
) (action string, err error) {
	space, err := s.repoAdapter.FindById(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errorSpaceNotFound
		}

		return
	}

	action = fmt.Sprintf(
		"update space of %s:%s/%s",
		spaceId.Identity(), space.Owner.Account(), space.Name.MSDName(),
	)

	if err = s.permission.CanUpdate(user, &space); err != nil {
		if allerror.IsNoPermission(err) {
			err = errorSpaceNotFound
		}

		return
	}

	b, err := s.codeRepoApp.Update(&space.CodeRepo, &cmd.CmdToUpdateRepo)
	if err != nil {
		return
	}

	b1 := cmd.toSpace(&space)
	if !b && !b1 {
		return
	}

	err = s.repoAdapter.Save(&space)

	e := domain.NewSpaceUpdatedEvent(&space)
	if err1 := s.msgAdapter.SendSpaceUpdatedEvent(&e); err1 != nil {
		logrus.Errorf("failed to send space updated event, space id:%s", spaceId.Identity())
	}

	return
}

func (s *spaceAppService) GetByName(user primitive.Account, index *domain.SpaceIndex) (SpaceDTO, error) {
	var dto SpaceDTO

	space, err := s.repoAdapter.FindByName(index)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errorSpaceNotFound
		}

		return dto, err
	}

	if err := s.permission.CanRead(user, &space); err != nil {
		if allerror.IsNoPermission(err) {
			err = errorSpaceNotFound
		}

		return dto, err
	}

	return toSpaceDTO(&space), nil
}

func (s *spaceAppService) List(user primitive.Account, cmd *CmdToListSpaces) (
	SpacesDTO, error,
) {
	if user == nil {
		cmd.Visibility = primitive.VisibilityPublic
	} else {
		if cmd.Owner == nil {
			// It can list the private spaces of user,
			// but it maybe no need to do it.
			cmd.Visibility = primitive.VisibilityPublic
		} else {
			if user != cmd.Owner {
				err := s.permission.CanListOrgResource(
					user, cmd.Owner, primitive.ObjTypeSpace,
				)
				if err != nil {
					cmd.Visibility = primitive.VisibilityPublic
				}
			}
		}
	}

	v, total, err := s.repoAdapter.List(cmd)

	return SpacesDTO{
		Total:  total,
		Spaces: v,
	}, err
}

// DeleteById is an example for restful API.
func (s *spaceAppService) DeleteById(user primitive.Account, spaceId string) error {
	// get space by space id
	// check if user can delete it
	// delete it
	return nil
}

func (s *spaceAppService) spaceCountCheck(owner primitive.Account) error {
	cmdToList := CmdToListSpaces{
		Owner: owner,
	}

	total, err := s.repoAdapter.Count(&cmdToList)
	if err != nil {
		return err
	}

	if total >= config.MaxCountPerOwner {
		return errorSpaceCountExceeded
	}

	return nil
}
