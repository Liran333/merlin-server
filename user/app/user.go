package app

import (
	"github.com/openmerlin/merlin-server/user/domain"
	"github.com/openmerlin/merlin-server/user/domain/repository"
	"github.com/openmerlin/merlin-server/user/infrastructure/git"
)

type UserService interface {
	// user
	Create(*domain.UserCreateCmd) (UserDTO, error)
	UserInfo(domain.Account) (UserInfoDTO, error)
	UpdateBasicInfo(domain.Account, UpdateUserBasicInfoCmd) error
	GetByAccount(domain.Account) (UserDTO, error)
	GetByFollower(owner, follower domain.Account) (UserDTO, bool, error)
}

// ps: platform user service
func NewUserService(
	repo repository.User,
	git git.User,
) UserService {
	return userService{
		repo: repo,
		git:  git,
	}
}

type userService struct {
	repo repository.User
	git  git.User
}

func (s userService) Create(cmd *domain.UserCreateCmd) (dto UserDTO, err error) {
	v := cmd.ToUser()

	// create user
	u, err := s.repo.Save(&v)
	if err != nil {
		return
	}

	// create git user
	if err = s.git.Create(cmd); err != nil {
		return
	}
	dto = newUserDTO(&u)

	return
}

func (s userService) UserInfo(account domain.Account) (dto UserInfoDTO, err error) {
	if dto.UserDTO, err = s.GetByAccount(account); err != nil {
		return
	}

	return
}

func (s userService) GetByAccount(account domain.Account) (dto UserDTO, err error) {
	// update user
	u, err := s.repo.GetByAccount(account)
	if err != nil {
		return
	}

	dto = newUserDTO(&u)

	return
}

func (s userService) GetByFollower(owner, follower domain.Account) (
	dto UserDTO, isFollower bool, err error,
) {
	v, isFollower, err := s.repo.GetByFollower(owner, follower)
	if err != nil {
		return
	}

	dto = newUserDTO(&v)

	return
}
