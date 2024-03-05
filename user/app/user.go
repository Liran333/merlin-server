/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	session "github.com/openmerlin/merlin-server/session/domain/repository"
	"github.com/openmerlin/merlin-server/user/domain"
	"github.com/openmerlin/merlin-server/user/domain/platform"
	"github.com/openmerlin/merlin-server/user/domain/repository"
	"github.com/openmerlin/merlin-server/user/infrastructure/git"
	"github.com/openmerlin/merlin-server/utils"
)

const tokenLen = 8

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

type SessionClearAppService interface {
	ClearAllSession(user primitive.Account) error
}

// UserService is an interface for user-related operations.
type UserService interface {
	Create(*domain.UserCreateCmd) (UserDTO, error)
	Delete(domain.Account) error
	RequestDelete(domain.Account) error
	UpdateBasicInfo(domain.Account, UpdateUserBasicInfoCmd) (UserDTO, error)
	UserInfo(domain.Account, domain.Account) (UserInfoDTO, error)
	GetByAccount(domain.Account, domain.Account) (UserDTO, error)
	GetUserAvatarId(domain.Account) (AvatarDTO, error)
	GetUserFullname(domain.Account) (string, error)
	GetUsersAvatarId([]domain.Account) ([]AvatarDTO, error)
	HasUser(primitive.Account) bool

	ListUsers(primitive.Account) ([]UserDTO, error)

	GetPlatformUser(domain.Account) (platform.BaseAuthClient, error)
	GetPlatformUserInfo(domain.Account) (string, error)

	CreateToken(*domain.TokenCreatedCmd, platform.BaseAuthClient) (TokenDTO, error)
	DeleteToken(*domain.TokenDeletedCmd, platform.BaseAuthClient) error
	ListTokens(domain.Account) ([]TokenDTO, error)
	GetToken(domain.Account, primitive.TokenName) (TokenDTO, error)
	VerifyToken(string, primitive.TokenPerm) (TokenDTO, error)

	SendBindEmail(*CmdToSendBindEmail) error
	VerifyBindEmail(*CmdToVerifyBindEmail) error

	PrivacyRevoke(user primitive.Account) error
}

// NewUserService creates a new UserService instance with the provided dependencies.
func NewUserService(
	repo repository.User,
	git git.User,
	token repository.Token,
	session session.SessionRepositoryAdapter,
	oidc session.OIDCAdapter,
	sc SessionClearAppService,
) UserService {
	return userService{
		repo:         repo,
		git:          git,
		oidc:         oidc,
		token:        token,
		session:      session,
		sessionClear: sc,
	}
}

type userService struct {
	repo         repository.User
	git          git.User
	oidc         session.OIDCAdapter
	token        repository.Token
	session      session.SessionRepositoryAdapter
	sessionClear SessionClearAppService
}

// Create creates a new user in the system.
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

	if cmd.Email != nil && cmd.Email.Email() != "" {
		// create git user when email is valid
		var repoUser domain.User
		if repoUser, err = s.git.Create(cmd); err != nil {
			err = allerror.NewInvalidParam(fmt.Sprintf("failed to create platform user: %s", err))
			return
		}

		v.PlatformId = repoUser.PlatformId
		v.PlatformPwd = repoUser.PlatformPwd
	}
	// create user
	u, err := s.repo.AddUser(&v)
	if err != nil {
		err = allerror.NewInvalidParam(fmt.Sprintf("failed to save user in db: %s", err))
		return
	}

	u.PlatformPwd = ""
	dto = newUserDTO(&u, cmd.Account)

	return
}

// UpdateBasicInfo updates the basic information of a user in the system.
func (s userService) UpdateBasicInfo(account domain.Account, cmd UpdateUserBasicInfoCmd) (dto UserDTO, err error) {
	user, err := s.repo.GetByAccount(account)
	if err != nil {
		return
	}

	if b := cmd.toUser(&user); !b {
		logrus.Warn("nothing changed")
		dto = newUserDTO(&user, account)
		return
	}

	if user, err = s.repo.SaveUser(&user); err != nil {
		err = allerror.NewInvalidParam("failed to update user info")
		return
	}

	// update git user when email is valid
	if user.Email != nil && user.Email.Email() != "" {
		if err = s.git.Update(&domain.UserCreateCmd{
			Account:  user.Account,
			Fullname: user.Fullname,
			Email:    user.Email,
			AvatarId: user.AvatarId,
			Desc:     user.Desc,
		}); err != nil {
			logrus.Error(err)
			err = allerror.NewInvalidParam("failed to update git user info")
			return
		}
	}

	dto = newUserDTO(&user, account)

	return
}

// GetPlatformUser retrieves the platform user info for the given account.
func (s userService) GetPlatformUserInfo(account domain.Account) (string, error) {
	if account == nil {
		return "", allerror.NewInvalidParam("account is nil")
	}
	// get user from db
	usernew, err := s.repo.GetByAccount(account)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errUserNotFound(fmt.Sprintf("user %s not found", account.Account()))
		}

		return "", err
	}

	return usernew.PlatformPwd, nil
}

// GetPlatformUser retrieves the platform user for the given account.
func (s userService) GetPlatformUser(account domain.Account) (token platform.BaseAuthClient, err error) {
	p, err := s.GetPlatformUserInfo(account)
	if err != nil {
		return
	}

	return git.NewBaseAuthClient(
		account.Account(),
		p,
	)
}

// HasUser checks if a user with the given account exists in the system.
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

// Delete deletes a user from the system.
func (s userService) Delete(account domain.Account) (err error) {
	u, err := s.repo.GetByAccount(account)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			logrus.Warnf("user %s not found", account.Account())
			err = nil
		}

		return
	}

	// delete user
	err = s.repo.DeleteUser(&u)
	if err != nil {
		err = allerror.NewInvalidParam(fmt.Sprintf("failed to delete user in db: %s", err))
		return
	}

	if u.Email != nil && u.Email.Email() != "" {
		// delete git user when email is valid
		err = s.git.Delete(u.Account)
		if err != nil {
			err = allerror.NewInvalidParam(fmt.Sprintf("failed to delete user in git server: %s", err))
		}
	}

	return
}

func (s userService) RequestDelete(user domain.Account) error {
	u, err := s.repo.GetByAccount(user)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errUserNotFound(fmt.Sprintf("user %s not found", user.Account()))
		}

		return err
	}

	if u.RequestDelete {
		return allerror.NewInvalidParam("user already requested to be delete")
	}

	u.RequestDelete = true
	u.RequestDeleteAt = utils.Now()

	_, err = s.repo.SaveUser(&u)

	return err
}

// UserInfo returns the user information for the given actor and account.
func (s userService) UserInfo(actor, account domain.Account) (dto UserInfoDTO, err error) {
	if dto.UserDTO, err = s.GetByAccount(actor, account); err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errUserNotFound(fmt.Sprintf("user %s not found", account.Account()))
		}

		return
	}

	return
}

// GetByAccount retrieves the user information by the given account.
func (s userService) GetByAccount(actor, account domain.Account) (dto UserDTO, err error) {
	// get user
	u, err := s.repo.GetByAccount(account)
	if err != nil {
		logrus.Error(err)
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errUserNotFound(fmt.Sprintf("user %s not found", account.Account()))
		}

		return
	}

	u.PlatformPwd = ""

	dto = newUserDTO(&u, actor)

	return
}

// ListUsers returns a list of users.
func (s userService) ListUsers(actor primitive.Account) (dtos []UserDTO, err error) {
	// get user
	t := domain.UserTypeUser
	u, _, err := s.repo.ListAccount(&repository.ListOption{Type: &t})
	if err != nil {
		logrus.Error(err)
		return
	}

	dtos = make([]UserDTO, len(u))
	for i := range u {
		dtos[i] = newUserDTO(&u[i], actor)
	}

	return
}

// GetUserAvatarId returns the avatar ID for the given user.
func (s userService) GetUserAvatarId(user domain.Account) (
	AvatarDTO, error,
) {
	var ava AvatarDTO
	a, err := s.repo.GetUserAvatarId(user)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errUserNotFound(fmt.Sprintf("user %s not found", user.Account()))
		}

		return ava, err
	}

	return AvatarDTO{
		Name:     user.Account(),
		AvatarId: a.AvatarId(),
	}, nil
}

// GetUsersAvatarId returns the avatar IDs for the given users.
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
		}

		return
	}

	dtos = make([]AvatarDTO, len(us))
	for i := range us {
		dtos[i] = ToAvatarDTO(&us[i])
	}

	return
}

// GetUserFullname returns the full name of the given user.
func (s userService) GetUserFullname(user domain.Account) (
	string, error,
) {
	name, err := s.repo.GetUserFullname(user)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errUserNotFound(fmt.Sprintf("user %s not found", user.Account()))
		}

		return "", err
	}

	return name, nil
}

// CreateToken creates a token for the given command and client.
func (s userService) CreateToken(cmd *domain.TokenCreatedCmd,
	client platform.BaseAuthClient) (token TokenDTO, err error) {
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
	if err != nil && !commonrepo.IsErrorResourceNotExists(err) {
		logrus.Error(err)
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

// DeleteToken deletes a token for the given account and name.
func (s userService) DeleteToken(cmd *domain.TokenDeletedCmd, client platform.BaseAuthClient) (err error) {
	if err = cmd.Validate(); err != nil {
		return
	}

	_, err = s.token.GetByName(cmd.Account, cmd.Name)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeTokenNotFound, "token not found")
		}

		return
	}

	err = client.DeleteToken(cmd)
	if err != nil {
		logrus.Error(err)
		return
	}

	err = s.token.Delete(cmd.Account, cmd.Name)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		}
	}

	return
}

// ListTokens lists all tokens for the given account.
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

// VerifyToken verifies a token with the given permission.
func (s userService) VerifyToken(token string, perm primitive.TokenPerm) (dto TokenDTO, err error) {
	if token == "" {
		err = allerror.New(allerror.ErrorCodeAccessTokenInvalid, "empty token")
		return
	}

	if len(token) < tokenLen {
		err = allerror.New(allerror.ErrorCodeAccessTokenInvalid, "token too short")
		return
	}

	tokens, err := s.token.GetByLastEight(token[len(token)-tokenLen:])
	if err != nil {
		logrus.Errorf("failed to find token: %s", err)
		err = allerror.New(allerror.ErrorCodeAccessTokenInvalid, "failed to find token")
		return
	}

	if len(tokens) == 0 {
		err = allerror.New(allerror.ErrorCodeAccessTokenInvalid, "not a valid token")
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

// GetToken gets a token by account and name.
func (s userService) GetToken(acc domain.Account, name primitive.TokenName) (TokenDTO, error) {
	token, err := s.token.GetByName(acc, name)
	if err != nil {
		logrus.Error(err)
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errTokenNotFound("token not found")
		}

		return TokenDTO{}, err
	}

	if token.Account.Account() != acc.Account() {
		return TokenDTO{}, allerror.NewNoPermission("token not found")
	}

	return newTokenDTO(&token), nil
}

// SendBindEmail sends an email to bind the account.
func (s userService) SendBindEmail(cmd *CmdToSendBindEmail) error {
	return s.oidc.SendBindEmail(cmd.Email.Email(), cmd.Capt)
}

// VerifyBindEmail verifies the email binding for a user.
func (s userService) VerifyBindEmail(cmd *CmdToVerifyBindEmail) error {
	userId, err := s.getUserIdOfLogin(cmd.User)
	if err != nil {
		return err
	}

	u, err := s.repo.GetByAccount(cmd.User)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errUserNotFound(fmt.Sprintf("user %s not found", cmd.User.Account()))
		}

		return err
	}

	if u.Email.Email() == cmd.Email.Email() {
		return nil
	}

	if u.Email.Email() != "" {
		return allerror.New(allerror.ErrorCodeEmailDuplicateBind, "user already bind another email address")
	}

	err = s.oidc.VerifyBindEmail(cmd.Email.Email(), cmd.PassCode, userId)
	if err != nil && !allerror.IsUserDuplicateBind(err) {
		return err
	}

	u.Email = cmd.Email

	userCmd := &domain.UserCreateCmd{
		Email:    u.Email,
		Account:  u.Account,
		Fullname: u.Fullname,
		Desc:     u.Desc,
		AvatarId: u.AvatarId,
	}

	if u.PlatformId == 0 {
		// create new user if user doesnot exist
		user, err := s.git.Create(userCmd)
		if err != nil {
			return err
		}

		u.PlatformId = user.PlatformId
		u.PlatformPwd = user.PlatformPwd
	} else {
		err = s.git.Update(userCmd)
		if err != nil {
			return err
		}
	}
	// we must create git user before save
	// bcs we need save platform id&pwd
	_, err = s.repo.SaveUser(&u)

	return err
}

// getUserIdOfLogin get user id of login
func (s userService) getUserIdOfLogin(user primitive.Account) (userId string, err error) {
	loginInfo, err := s.session.FindByUser(user)
	if err != nil {
		return
	}

	if len(loginInfo) == 0 {
		err = fmt.Errorf("user session not found")

		return
	}

	if loginInfo[0].UserId == "" {
		err = fmt.Errorf("user id not found")

		return
	}

	userId = loginInfo[0].UserId

	return
}

// PrivacyRevoke revokes the privacy settings for a user.
func (s userService) PrivacyRevoke(user primitive.Account) error {
	sessions, err := s.session.FindByUser(user)
	if err != nil {
		return err
	}

	if len(sessions) == 0 || sessions[0].UserId == "" {
		return fmt.Errorf("user id not found")
	}

	if err = s.oidc.PrivacyRevoke(sessions[0].UserId); err != nil {
		return err
	}

	return s.sessionClear.ClearAllSession(user)
}
