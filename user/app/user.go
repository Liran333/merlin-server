package app

import (
	"github.com/openmerlin/merlin-server/user/domain"
	"github.com/openmerlin/merlin-server/user/domain/repository"
)

type UserService interface {
	// user
	Create(*UserCreateCmd) (UserDTO, error)
	UserInfo(domain.Account) (UserInfoDTO, error)
	UpdateBasicInfo(domain.Account, UpdateUserBasicInfoCmd) error
	GetByAccount(domain.Account) (UserDTO, error)
	GetByFollower(owner, follower domain.Account) (UserDTO, bool, error)
}

// ps: platform user service
func NewUserService(
	repo repository.User,
) UserService {
	return userService{
		repo: repo,
	}
}

type userService struct {
	repo repository.User
}

func (s userService) Create(cmd *UserCreateCmd) (dto UserDTO, err error) {
	v := cmd.toUser()

	// update user
	u, err := s.repo.Save(&v)
	if err != nil {
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
