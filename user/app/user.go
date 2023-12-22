package app

import (
	"fmt"

	"github.com/openmerlin/merlin-server/user/domain"
	"github.com/openmerlin/merlin-server/user/domain/platform"
	"github.com/openmerlin/merlin-server/user/domain/repository"
	"github.com/openmerlin/merlin-server/user/infrastructure/git"
	"github.com/openmerlin/merlin-server/utils"
	"github.com/sirupsen/logrus"
)

type UserService interface {
	// user
	Create(*domain.UserCreateCmd) (UserDTO, error)
	Delete(domain.Account) error
	UpdateBasicInfo(domain.Account, UpdateUserBasicInfoCmd) error

	UserInfo(domain.Account) (UserInfoDTO, error)
	GetByAccount(domain.Account, bool) (UserDTO, error)
	GetByFollower(owner, follower domain.Account) (UserDTO, bool, error)
	GetUserAvatarId(domain.Account) (AvatarDTO, error)
	GetUserFullname(domain.Account) (string, error)
	GetUsersAvatarId([]domain.Account) ([]AvatarDTO, error)

	GetPlatformUser(domain.Account) (platform.BaseAuthClient, error)

	CreateToken(*domain.TokenCreatedCmd, platform.BaseAuthClient) (TokenDTO, error)
	DeleteToken(*domain.TokenDeletedCmd, platform.BaseAuthClient) error
	ListTokens(domain.Account) ([]TokenDTO, error)
	VerifyToken(string) (TokenDTO, bool)
}

// ps: platform user service
func NewUserService(
	repo repository.User,
	git git.User,
	token repository.Token,
) UserService {
	return userService{
		repo:  repo,
		git:   git,
		token: token,
	}
}

type userService struct {
	repo  repository.User
	git   git.User
	token repository.Token
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
	v.CreatedAt = utils.Now()
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
	usernew, err := s.GetByAccount(account, true)
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
	if dto.UserDTO, err = s.GetByAccount(account, false); err != nil {
		return
	}

	return
}

func (s userService) GetByAccount(account domain.Account, pwd bool) (dto UserDTO, err error) {
	// update user
	u, err := s.repo.GetByAccount(account)
	if err != nil {
		return
	}

	if !pwd {
		u.PlatformPwd = ""
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

	v.PlatformPwd = ""

	dto = newUserDTO(&v)

	return
}

func (s userService) GetUserAvatarId(user domain.Account) (
	AvatarDTO, error,
) {
	var ava AvatarDTO
	a, err := s.repo.GetUserAvatarId(user)
	if err != nil {
		return ava, err
	}

	return AvatarDTO{
		Name:     user.Account(),
		AvatarId: a.AvatarId(),
	}, nil
}

func (s userService) GetUsersAvatarId(users []domain.Account) (
	[]AvatarDTO, error,
) {
	names := make([]string, len(users))
	for i := range users {
		names[i] = users[i].Account()
	}
	us, err := s.repo.GetUsersAvatarId(names)
	if err != nil {
		return nil, err
	}

	dtos := make([]AvatarDTO, len(us))
	for i := range us {
		dtos[i] = ToAvatarDTO(&us[i])
	}

	return dtos, nil
}

func (s userService) GetUserFullname(user domain.Account) (
	string, error,
) {
	return s.repo.GetUserFullname(user)

}

func (s userService) CreateToken(cmd *domain.TokenCreatedCmd, client platform.BaseAuthClient) (token TokenDTO, err error) {
	if err = cmd.Validate(); err != nil {
		return
	}

	t, err := client.CreateToken(cmd)
	if err != nil {
		return
	}

	enc, salt, err := domain.EncryptToken(t.Token)
	if err != nil {
		return
	}

	token = newTokenDTO(&t)

	t.Token = enc
	t.Salt = salt

	_, err = s.token.Save(&t)

	return
}

func (s userService) DeleteToken(cmd *domain.TokenDeletedCmd, client platform.BaseAuthClient) (err error) {
	if err = cmd.Validate(); err != nil {
		return
	}

	err = client.DeleteToken(cmd)
	if err != nil {
		return
	}

	err = s.token.Delete(cmd.Account, cmd.Name)

	return
}

func (s userService) ListTokens(u domain.Account) (tokens []TokenDTO, err error) {
	ts, err := s.token.GetByAccount(u)
	if err != nil {
		return
	}

	tokens = make([]TokenDTO, len(ts))
	for t := range ts {
		tokens[t] = newTokenDTO(&ts[t])
		tokens[t].Token = ""
	}

	return
}

func (s userService) VerifyToken(token string) (dto TokenDTO, b bool) {
	tokens, err := s.token.GetByLastEight(token[len(token)-8:])
	if err != nil {
		logrus.Errorf("get token by last eight failed: %s", err)
		return
	}

	for t := range tokens {
		if tokens[t].Compare(token) {
			b = true
			dto = newTokenDTO(&tokens[t])
			return
		}
	}

	return
}
