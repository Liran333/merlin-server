package app

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/user/domain"
	"github.com/openmerlin/merlin-server/user/domain/platform"
	"github.com/openmerlin/merlin-server/user/domain/repository"
	"github.com/openmerlin/merlin-server/user/infrastructure/git"
)

func errUserNotFound(msg string) error {
	if msg == "" {
		msg = "user not found"
	}

	return allerror.NewNotFound(allerror.ErrorCodeUserNotFound, msg)
}

func errTokenNotFound(msg string) error {
	if msg == "" {
		msg = "token not found"
	}

	return allerror.NewNotFound(allerror.ErrorCodeTokenNotFound, msg)
}

type UserService interface {
	// user
	Create(*domain.UserCreateCmd) (UserDTO, error)
	Delete(domain.Account) error
	UpdateBasicInfo(domain.Account, UpdateUserBasicInfoCmd) (UserDTO, error)

	UserInfo(domain.Account) (UserInfoDTO, error)
	GetByAccount(domain.Account, bool) (UserDTO, error)
	GetUserAvatarId(domain.Account) (AvatarDTO, error)
	GetUserFullname(domain.Account) (string, error)
	GetUsersAvatarId([]domain.Account) ([]AvatarDTO, error)
	HasUser(primitive.Account) bool

	GetPlatformUser(domain.Account) (platform.BaseAuthClient, error)

	CreateToken(*domain.TokenCreatedCmd, platform.BaseAuthClient) (TokenDTO, error)
	DeleteToken(*domain.TokenDeletedCmd, platform.BaseAuthClient) error
	ListTokens(domain.Account) ([]TokenDTO, error)
	GetToken(domain.Account, domain.Account) (TokenDTO, error)
	VerifyToken(string, primitive.TokenPerm) (TokenDTO, error)
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
	// TODO: need perm orgapp.Permission
}

func (s userService) Create(cmd *domain.UserCreateCmd) (dto UserDTO, err error) {
	if cmd == nil {
		err = errUserNotFound("input param is empty")
		return
	}

	if err = cmd.Validate(); err != nil {
		err = allerror.NewInvalidParam(err.Error())
		return
	}

	if !s.repo.CheckName(cmd.Account) {
		err = allerror.NewInvalidParam(fmt.Sprintf("user name %s is already taken", cmd.Account))
		return
	}

	v := cmd.ToUser()

	// create git user
	var repoUser domain.User
	if repoUser, err = s.git.Create(cmd); err != nil {
		err = allerror.NewInvalidParam(fmt.Sprintf("failed to create platform user: %s", err))
		return
	}

	v.PlatformId = repoUser.PlatformId
	v.PlatformPwd = repoUser.PlatformPwd
	// create user
	u, err := s.repo.AddUser(&v)
	if err != nil {
		err = allerror.NewInvalidParam(fmt.Sprintf("failed to save user in db: %s", err))
		s.git.Delete(&repoUser) // #nosec G104
		return
	}

	u.PlatformPwd = ""
	dto = newUserDTO(&u)

	return
}

func (s userService) UpdateBasicInfo(account domain.Account, cmd UpdateUserBasicInfoCmd) (dto UserDTO, err error) {
	user, err := s.repo.GetByAccount(account)
	if err != nil {
		return
	}

	if b := cmd.toUser(&user); !b {
		dto = newUserDTO(&user)
		return
	}

	if user, err = s.repo.SaveUser(&user); err != nil {
		err = allerror.NewInvalidParam("failed to update user info")
		return
	}

	dto = newUserDTO(&user)

	return
}

func (s userService) GetPlatformUser(account domain.Account) (token platform.BaseAuthClient, err error) {
	if account == nil {
		err = allerror.NewInvalidParam("account is nil")
		return
	}
	usernew, err := s.GetByAccount(account, true)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errUserNotFound(fmt.Sprintf("user %s not found", account.Account()))
		} else {
			logrus.Error(err)
			err = allerror.NewInvalidParam("failed to get user info")
		}
		return
	}

	return git.NewBaseAuthClient(
		usernew.Account,
		usernew.Password,
	)
}

func (s userService) HasUser(acc primitive.Account) bool {
	if acc == nil {
		logrus.Errorf("username invalid")
		return false
	}

	_, err := s.repo.GetByAccount(acc)
	if err != nil {
		logrus.Errorf("user %s not found", acc.Account())
		return false
	}

	return true
}

func (s userService) Delete(account domain.Account) (err error) {
	u, err := s.repo.GetByAccount(account)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			logrus.Warnf("user %s not found", account.Account())
			err = nil
		} else {
			logrus.Error(err)
			err = allerror.NewInvalidParam("failed to delete user")
		}
		return
	}

	// delete user
	err = s.repo.DeleteUser(&u)
	if err != nil {
		err = allerror.NewInvalidParam(fmt.Sprintf("failed to delete user in db: %s", err))
		return
	}

	// delete git user
	err = s.git.Delete(&u)
	if err != nil {
		err = allerror.NewInvalidParam(fmt.Sprintf("failed to delete user in git server: %s", err))
	}

	return
}

func (s userService) UserInfo(account domain.Account) (dto UserInfoDTO, err error) {
	if dto.UserDTO, err = s.GetByAccount(account, false); err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errUserNotFound(fmt.Sprintf("user %s not found", account.Account()))
		} else {
			logrus.Error(err)
			err = allerror.NewInvalidParam("failed to get user info")
		}
		return
	}

	return
}

func (s userService) GetByAccount(account domain.Account, pwd bool) (dto UserDTO, err error) {
	// update user
	u, err := s.repo.GetByAccount(account)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errUserNotFound(fmt.Sprintf("user %s not found", account.Account()))
		} else {
			logrus.Error(err)
			err = allerror.NewInvalidParam("failed to get user")
		}
		return
	}

	if !pwd {
		u.PlatformPwd = ""
	}

	dto = newUserDTO(&u)

	return
}

func (s userService) GetUserAvatarId(user domain.Account) (
	AvatarDTO, error,
) {
	var ava AvatarDTO
	a, err := s.repo.GetUserAvatarId(user)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errUserNotFound(fmt.Sprintf("user %s not found", user.Account()))
		} else {
			logrus.Error(err)
			err = allerror.NewInvalidParam("failed to get user avatar id")
		}
		return ava, err
	}

	return AvatarDTO{
		Name:     user.Account(),
		AvatarId: a.AvatarId(),
	}, nil
}

func (s userService) GetUsersAvatarId(users []domain.Account) (
	dtos []AvatarDTO, err error,
) {
	names := make([]string, len(users))
	for i := range users {
		names[i] = users[i].Account()
	}

	us, err := s.repo.GetUsersAvatarId(names)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		} else {
			logrus.Error(err)
			err = allerror.NewInvalidParam("failed to get user avatar id")
		}
		return
	}

	dtos = make([]AvatarDTO, len(us))
	for i := range us {
		dtos[i] = ToAvatarDTO(&us[i])
	}

	return
}

func (s userService) GetUserFullname(user domain.Account) (
	string, error,
) {
	name, err := s.repo.GetUserFullname(user)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errUserNotFound(fmt.Sprintf("user %s not found", user.Account()))
		} else {
			logrus.Error(err)
			err = allerror.NewInvalidParam("failed to get user fullname")
		}
		return "", err
	}

	return name, nil
}

func (s userService) CreateToken(cmd *domain.TokenCreatedCmd, client platform.BaseAuthClient) (token TokenDTO, err error) {
	if err = cmd.Validate(); err != nil {
		return
	}

	owner, err := s.repo.GetByAccount(cmd.Account)
	if err != nil {
		logrus.Error(err)
		err = allerror.NewInvalidParam("failed to get owenr info")
		return
	}

	_, err = s.token.GetByName(cmd.Account, cmd.Name)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		} else {
			logrus.Error(err)

			err = allerror.NewInvalidParam("failed to get token info")
			return
		}
	} else {
		err = allerror.NewInvalidParam("token is already existed")
		return
	}

	t, err := client.CreateToken(cmd)
	if err != nil {
		logrus.Error(err)
		err = allerror.NewInvalidParam("failed to create platform token")
		return
	}

	enc, salt, err := domain.EncryptToken(t.Token)
	if err != nil {
		logrus.Error(err)
		err = allerror.NewInvalidParam("failed to encrypt token")
		return
	}

	// token without encrypted
	orgtoken := t.Token
	t.Token = enc
	t.Salt = salt
	t.OwnerId = owner.Id

	t, err = s.token.Add(&t)
	token = newTokenDTO(&t)
	token.Token = orgtoken

	return
}

func (s userService) DeleteToken(cmd *domain.TokenDeletedCmd, client platform.BaseAuthClient) (err error) {
	if err = cmd.Validate(); err != nil {
		return
	}

	_, err = s.token.GetByName(cmd.Account, cmd.Name)
	if err != nil {
		logrus.Error(err)
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errTokenNotFound("token not found")
		}
		return
	}

	err = client.DeleteToken(cmd)
	if err != nil {
		logrus.Error(err)
		err = allerror.NewInvalidParam("failed to delete platform token")
		return
	}

	err = s.token.Delete(cmd.Account, cmd.Name)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		} else {
			logrus.Error(err)
			err = allerror.NewInvalidParam("failed to delete token")
		}
	}
	return
}

func (s userService) ListTokens(u domain.Account) (tokens []TokenDTO, err error) {
	if u == nil {
		err = allerror.NewInvalidParam("username is empty")
		return
	}

	ts, err := s.token.GetByAccount(u)
	if err != nil {
		logrus.Error(err)
		err = allerror.NewInvalidParam("failed to get user info")
		return
	}

	tokens = make([]TokenDTO, len(ts))
	for t := range ts {
		tokens[t] = newTokenDTO(&ts[t])
		tokens[t].Token = ""
	}

	return
}

func (s userService) VerifyToken(token string, perm primitive.TokenPerm) (dto TokenDTO, err error) {
	if token == "" {
		err = allerror.New(allerror.ErrorCodeAccessTokenInvalid, "empty token")
		return
	}

	tokens, err := s.token.GetByLastEight(token[len(token)-8:])
	if err != nil {
		logrus.Errorf("failed to find token: %s", err)
		err = allerror.NewNoPermission("failed to find token")
		return
	}

	if len(tokens) == 0 {
		err = allerror.NewNoPermission("not a valid token")
		return
	}

	for t := range tokens {
		if err = tokens[t].Check(token, perm); err == nil {
			dto = newTokenDTO(&tokens[t])
			return
		}
	}

	return
}

func (s userService) GetToken(acc, name domain.Account) (TokenDTO, error) {
	token, err := s.token.GetByName(acc, name)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errTokenNotFound("token not found")
		}
		logrus.Error(err)
		return TokenDTO{}, err
	}

	if token.Account.Account() != acc.Account() {
		return TokenDTO{}, allerror.NewNoPermission("token not found")
	}

	return newTokenDTO(&token), nil

}
