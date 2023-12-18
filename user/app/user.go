package app

import (
	"fmt"

	"github.com/openmerlin/merlin-server/user/domain"
	"github.com/openmerlin/merlin-server/user/domain/platform"
	"github.com/openmerlin/merlin-server/user/domain/repository"
	"github.com/openmerlin/merlin-server/user/infrastructure/git"
)

type UserService interface {
	// user
	Create(*domain.UserCreateCmd) (UserDTO, error)
	Delete(domain.Account) error
	UpdateBasicInfo(domain.Account, UpdateUserBasicInfoCmd) error

	UserInfo(domain.Account) (UserInfoDTO, error)
	GetByAccount(domain.Account) (UserDTO, error)
	GetByFollower(owner, follower domain.Account) (UserDTO, bool, error)

	GetPlatformUser(domain.Account) (platform.BaseAuthClient, error)

	CreateToken(*domain.TokenCreatedCmd, platform.BaseAuthClient) (TokenDTO, error)
	DeleteToken(*domain.TokenDeletedCmd, platform.BaseAuthClient) error
	ListTokens(domain.Account) ([]TokenDTO, error)
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

	// create git user
	var repoUser domain.User
	if repoUser, err = s.git.Create(cmd); err != nil {
		err = fmt.Errorf("failed to create platform user: %w", err)
		return
	}

	v.PlatformId = repoUser.PlatformId
	v.PlatformPwd = repoUser.PlatformPwd
	// create user
	u, err := s.repo.Save(&v)
	if err != nil {
		err = fmt.Errorf("failed to save user in db: %w", err)
		s.git.Delete(&repoUser) // #nosec G104
		return
	}

	u.PlatformPwd = ""
	dto = newUserDTO(&u)

	return
}

func (s userService) GetPlatformUser(account domain.Account) (token platform.BaseAuthClient, err error) {
	if account == nil {
		err = fmt.Errorf("account is nil")
		return
	}
	usernew, err := s.GetByAccount(account)
	if err != nil {
		return
	}

	return git.NewBaseAuthClient(
		usernew.Account,
		usernew.Password,
	)
}

func (s userService) Delete(account domain.Account) (err error) {
	u, err := s.repo.GetByAccount(account)
	if err != nil {
		return
	}

	// delete user
	err = s.repo.Delete(&u)
	if err != nil {
		err = fmt.Errorf("failed to delete user in db: %w", err)
		return
	}

	// delete git user
	err = s.git.Delete(&u)
	if err != nil {
		err = fmt.Errorf("failed to delete user in git server: %w", err)
	}

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

func (s userService) CreateToken(cmd *domain.TokenCreatedCmd, client platform.BaseAuthClient) (token TokenDTO, err error) {
	if err = cmd.Validate(); err != nil {
		return
	}

	user, err := s.repo.GetByAccount(cmd.Account)
	if err != nil {
		return
	}

	t, err := client.CreateToken(cmd)
	if err != nil {
		return
	}

	user.PlatformTokens[cmd.Name] = domain.PlatformToken{
		CreatedAt:  t.CreatedAt,
		Permission: t.Permission,
		Expire:     t.Expire,
		Account:    t.Account,
		Name:       t.Name,
	}

	token = newTokenDTO(&t)

	_, err = s.repo.Save(&user)

	return
}

func (s userService) DeleteToken(cmd *domain.TokenDeletedCmd, client platform.BaseAuthClient) (err error) {
	if err = cmd.Validate(); err != nil {
		return
	}

	user, err := s.repo.GetByAccount(cmd.Account)
	if err != nil {
		return
	}

	err = client.DeleteToken(cmd)
	if err != nil {
		return
	}

	delete(user.PlatformTokens, cmd.Name)
	_, err = s.repo.Save(&user)

	return
}

func (s userService) ListTokens(u domain.Account) (tokens []TokenDTO, err error) {
	user, err := s.repo.GetByAccount(u)
	if err != nil {
		return
	}

	tokens = make([]TokenDTO, 0, len(user.PlatformTokens))
	for k := range user.PlatformTokens {
		t := user.PlatformTokens[k]
		tokens = append(tokens, newTokenDTO(&t))
	}

	return
}
