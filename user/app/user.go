/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"

	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	orgrepository "github.com/openmerlin/merlin-server/organization/domain/repository"
	session "github.com/openmerlin/merlin-server/session/domain/repository"
	"github.com/openmerlin/merlin-server/user/domain"
	"github.com/openmerlin/merlin-server/user/domain/platform"
	"github.com/openmerlin/merlin-server/user/domain/repository"
	"github.com/openmerlin/merlin-server/user/infrastructure/git"
	"github.com/openmerlin/merlin-server/utils"
)

const tokenLen = 8

// SessionClearAppService defines the application service interface for clearing sessions.
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
	GetOrgOrUser(primitive.Account, primitive.Account) (UserDTO, error)
	GetUserAvatarId(domain.Account) (AvatarDTO, error)
	GetUserFullname(domain.Account) (string, error)
	GetUsersAvatarId([]domain.Account) ([]AvatarDTO, error)
	HasUser(primitive.Account) bool

	IsOrganization(domain.Account) bool

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

	PrivacyRevoke(user primitive.Account) (string, error)
	AgreePrivacy(user primitive.Account) error
	IsAgreePrivacy(user primitive.Account) (bool, error)
}

// NewUserService creates a new UserService instance with the provided dependencies.
func NewUserService(
	repo repository.User,
	mem orgrepository.OrgMember,
	git git.User,
	token repository.Token,
	session session.SessionRepositoryAdapter,
	oidc session.OIDCAdapter,
	sc SessionClearAppService,
	config *domain.Config,
) UserService {
	return userService{
		repo:         repo,
		member:       mem,
		git:          git,
		oidc:         oidc,
		token:        token,
		session:      session,
		sessionClear: sc,
		config:       config,
	}
}

type userService struct {
	repo         repository.User
	member       orgrepository.OrgMember
	git          git.User
	oidc         session.OIDCAdapter
	token        repository.Token
	session      session.SessionRepositoryAdapter
	sessionClear SessionClearAppService
	config       *domain.Config
}

// Create creates a new user in the system.
func (s userService) Create(cmd *domain.UserCreateCmd) (dto UserDTO, err error) {
	if cmd == nil {
		e := xerrors.Errorf("input param is empty")
		err = allerror.NewCommonRespError(e.Error(), e)
		return
	}

	if err = cmd.Validate(); err != nil {
		err = allerror.NewInvalidParam(err.Error(), xerrors.Errorf("create user cmd validate error: %w", err))
		return
	}

	if !s.repo.CheckName(cmd.Account) {
		e := xerrors.Errorf("user name %s is already taken", cmd.Account)
		err = allerror.New(allerror.ErrorUsernameIsAlreadyTaken, e.Error(), e)
		return
	}

	v := cmd.ToUser()

	if cmd.Email != nil && cmd.Email.Email() != "" {
		// create git user when email is valid
		var repoUser domain.User
		if repoUser, err = s.git.Create(cmd); err != nil {
			e := xerrors.Errorf("failed to create platform user: %w", err)
			err = allerror.New(allerror.ErrorFailedToCreatePlatformUser, "", e)
			return
		}

		v.PlatformId = repoUser.PlatformId
		v.PlatformPwd = repoUser.PlatformPwd
	}
	// create user
	u, err := s.repo.AddUser(&v)
	if err != nil {
		e := xerrors.Errorf("failed to save user in db: %w", err)
		err = allerror.New(allerror.ErrorFailToSaveUserInDb, "", e)
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
		if commonrepo.IsErrorResourceNotExists(err) {
			e := xerrors.Errorf("user %s not found: %w", account.Account(), err)
			err = allerror.NewNotFound(allerror.ErrorCodeUserNotFound, "", e)
		} else {
			e := xerrors.Errorf("failed to get user: %w", err)
			err = allerror.NewCommonRespError("", e)
		}
		return
	}

	if b := cmd.toUser(&user); !b {
		logrus.Warn("nothing changed")
		dto = newUserDTO(&user, account)
		return
	}

	if user, err = s.repo.SaveUser(&user); err != nil {
		e := xerrors.Errorf("failed to update user info: %w", err)
		err = allerror.New(allerror.ErrorFailedToUpdateUserInfo, "", e)
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
			e := xerrors.Errorf("failed to update git user info: %w", err)
			err = allerror.New(allerror.ErrorFailedToUPdateGitUserInfo, "", e)
			return
		}
	}

	dto = newUserDTO(&user, account)

	return
}

// GetPlatformUser retrieves the platform user info for the given account.
func (s userService) GetPlatformUserInfo(account domain.Account) (string, error) {
	if account == nil {
		e := xerrors.Errorf("username invalid")
		return "", allerror.New(allerror.ErrorUsernameInvalid, e.Error(), e)
	}
	// get user from db
	usernew, err := s.repo.GetByAccount(account)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeUserNotFound, "",
				xerrors.Errorf("failed to get platform user info: %w", err))
		} else {
			err = xerrors.Errorf("failed to get platform user info: %w", err)
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

	token, err = git.NewBaseAuthClient(
		account.Account(),
		p,
	)
	if err != nil {
		err = xerrors.Errorf("failed to generate platform user: %w", err)
	}

	return
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
			logrus.Warnf("user %s not found, no need to delete", account.Account())
			err = nil
		} else {
			err = xerrors.Errorf("failed to get user: %w", err)
		}

		return
	}

	// delete user
	err = s.repo.DeleteUser(&u)
	if err != nil {
		err = allerror.New(allerror.ErrorFailedToDeleteUser, "",
			xerrors.Errorf("failed to delete user in db, %w", err))
		return
	}

	if u.Email != nil && u.Email.Email() != "" {
		// delete git user when email is valid
		err = s.git.Delete(u.Account)
		if err != nil {
			err = allerror.New(allerror.ErrorFailedToDeleteUserInGitServer, "",
				xerrors.Errorf("failed to delete user in git server, %w", err))
		}
	}

	return
}

func (s userService) RequestDelete(user domain.Account) error {
	u, err := s.repo.GetByAccount(user)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.New(allerror.ErrorCodeUserNotFound, "",
				xerrors.Errorf("user %s not found: %w", user.Account(), err))
		}

		return err
	}

	if u.RequestDelete {
		e := xerrors.Errorf("user already requested to be delete")
		return allerror.New(allerror.ErrorUserAlreadyRequestedToBeDelete, e.Error(), e)
	}

	memList, err := s.member.GetByUserAndRoles(user, []primitive.Role{primitive.Admin})
	if err != nil {
		return xerrors.Errorf("failed to get member list: %w", err)
	}
	if len(memList) > 0 {
		e := xerrors.Errorf("user is an admin role of other organization, can't be deleted")
		return allerror.NewInvalidParam(e.Error(), e)
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
			e := xerrors.Errorf("user %s not found: %w", account.Account(), err)
			err = allerror.New(allerror.ErrorCodeUserNotFound, "", e)
		} else {
			err = xerrors.Errorf("failed to get user info: %w", err)
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
		if commonrepo.IsErrorResourceNotExists(err) {
			e := xerrors.Errorf("user %s not found: %w", account.Account(), err)
			err = allerror.NewNotFound(allerror.ErrorCodeUserNotFound, "", e)
		} else {
			err = xerrors.Errorf("failed to get user: %w", err)
		}

		return
	}

	u.PlatformPwd = ""

	dto = newUserDTO(&u, actor)

	return
}

// GetOrgOrUser retrieves either an organization or a user by their account and returns it as a UserDTO.
func (s userService) GetOrgOrUser(actor, acc primitive.Account) (dto UserDTO, err error) {
	u, err := s.repo.GetByAccount(acc)
	if err != nil && !commonrepo.IsErrorResourceNotExists(err) {
		return
	} else if err == nil {
		dto = NewUserDTO(&u, actor)
		return
	}

	o, err := s.repo.GetOrgByName(acc)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.New(allerror.ErrorCodeUserNotFound, fmt.Sprintf("org %s not found", acc.Account()),
				fmt.Errorf("org %s not found, %w", acc.Account(), err))
		}

		return
	}

	dto = ToDTO(&o)
	return
}

// ListUsers returns a list of users.
func (s userService) ListUsers(actor primitive.Account) (dtos []UserDTO, err error) {
	// get user
	t := domain.UserTypeUser
	u, _, err := s.repo.ListAccount(&repository.ListOption{Type: &t})
	if err != nil {
		err = xerrors.Errorf("failed to list users: %w", err)
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
			e := xerrors.Errorf("user %s not found: %w", user.Account(), err)
			err = allerror.New(allerror.ErrorCodeUserNotFound, "", e)
		} else {
			err = xerrors.Errorf("failed to get user avatarid: %w", err)
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
		} else {
			err = xerrors.Errorf("failed to get users avatarid: %w", err)
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
			e := xerrors.Errorf("user %s not found: %w", user.Account(), err)
			err = allerror.New(allerror.ErrorCodeUserNotFound, "", e)
		} else {
			err = xerrors.Errorf("failed to get user fullname: %w", err)
		}

		return "", err
	}

	return name, nil
}

// CreateToken creates a token for the given command and client.
func (s userService) CreateToken(cmd *domain.TokenCreatedCmd,
	client platform.BaseAuthClient) (token TokenDTO, err error) {
	if err = cmd.Validate(); err != nil {
		err = allerror.NewInvalidParam(err.Error(), xerrors.Errorf("create token cmd validate error: %w", err))
		return
	}

	if ok, err1 := s.CanCreateToken(cmd.Account); !ok {
		err = allerror.NewCountExceeded("token count exceed", xerrors.Errorf("create token failed :%w", err1))
		return
	}

	owner, err := s.repo.GetByAccount(cmd.Account)
	if err != nil {
		err = xerrors.Errorf("failed to get user: %w", err)
		err = allerror.New(allerror.ErrorFailedToCreateToken, "failed to create token", err)
		return
	}

	_, err = s.token.GetByName(cmd.Account, cmd.Name)
	if err != nil && !commonrepo.IsErrorResourceNotExists(err) {
		err = xerrors.Errorf("failed to get token by name: %w", err)
		err = allerror.New(allerror.ErrorFailedToCreateToken, "failed to create token", err)
		return
	}

	t, err := client.CreateToken(cmd)
	if err != nil {
		err = xerrors.Errorf("failed to create token: %w", err)
		err = allerror.New(allerror.ErrorFailedToCreateToken, "failed to create token", err)
		return
	}

	enc, salt, err := domain.EncryptToken(t.Token)
	if err != nil {
		err = xerrors.Errorf("failed to encrypt token: %w", err)
		err = allerror.New(allerror.ErrorFailedToEcryptToken, "failed to encrypt token", err)
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
		err = allerror.NewInvalidParam(err.Error(), xerrors.Errorf("delete token cmd validate error: %w", err))
		return
	}

	_, err = s.token.GetByName(cmd.Account, cmd.Name)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeTokenNotFound, "",
				xerrors.Errorf("token not found: %w", err))
		} else {
			err = xerrors.Errorf("failed to get token: %w", err)
		}

		return
	}

	err = client.DeleteToken(cmd)
	if err != nil {
		err = allerror.New(allerror.ErrorFailedToDeleteToken, "",
			xerrors.Errorf("failed to delete token: %w", err))
		return
	}

	err = s.token.Delete(cmd.Account, cmd.Name)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			logrus.Warnf("token %s not found, no need to delete", cmd.Name)
			err = nil
		}
	}

	return
}

// ListTokens lists all tokens for the given account.
func (s userService) ListTokens(u domain.Account) (tokens []TokenDTO, err error) {
	if u == nil {
		e := xerrors.Errorf("input param is empty")
		err = allerror.New(allerror.ErrorInputParamIsEmpty, e.Error(), e)
		return
	}

	ts, err := s.token.GetByAccount(u)
	if err != nil {
		err = allerror.New(allerror.ErrorFailedToGetUserInfo, "",
			xerrors.Errorf("failed to get user info: %w", err))
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
		e := xerrors.Errorf("input param is empty")
		err = allerror.New(allerror.ErrorCodeAccessTokenInvalid, e.Error(), e)
		return
	}

	if len(token) < tokenLen {
		e := xerrors.Errorf("token too short")
		err = allerror.New(allerror.ErrorCodeAccessTokenInvalid, e.Error(), e)
		return
	}

	tokens, err := s.token.GetByLastEight(token[len(token)-tokenLen:])
	if err != nil {
		err = xerrors.Errorf("failed to get token: %w", err)
		err = allerror.New(allerror.ErrorCodeAccessTokenInvalid, "", err)
		return
	}

	if len(tokens) == 0 {
		e := xerrors.Errorf("not a valid token")
		err = allerror.New(allerror.ErrorCodeAccessTokenInvalid, e.Error(), e)
		return
	}

	for t := range tokens {
		if err = tokens[t].Check(token, perm); err == nil {
			dto = newTokenDTO(&tokens[t])
			return
		}
	}

	err = allerror.NewNoPermission("", xerrors.Errorf("not a valid token"))
	return
}

// GetToken gets a token by account and name.
func (s userService) GetToken(acc domain.Account, name primitive.TokenName) (TokenDTO, error) {
	token, err := s.token.GetByName(acc, name)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.New(allerror.ErrorCodeTokenNotFound, "",
				xerrors.Errorf("token not found: %w", err))
		} else {
			err = xerrors.Errorf("failed to get token: %w", err)
		}

		return TokenDTO{}, err
	}

	if token.Account.Account() != acc.Account() {
		return TokenDTO{}, allerror.NewNoPermission("", xerrors.Errorf("can't get others token"))
	}

	newToken := newTokenDTO(&token)
	newToken.Token = ""

	return newToken, nil
}

// SendBindEmail sends an email to bind the account.
func (s userService) SendBindEmail(cmd *CmdToSendBindEmail) error {
	return s.oidc.SendBindEmail(cmd.Email.Email(), cmd.Capt)
}

// VerifyBindEmail verifies the email binding for a user.
func (s userService) VerifyBindEmail(cmd *CmdToVerifyBindEmail) error {
	userId, err := s.getUserIdOfLogin(cmd.User)
	if err != nil {
		return allerror.New(allerror.ErrorCodeUserNotFound, "",
			xerrors.Errorf("user %s not found: %w", cmd.User.Account(), err))
	}

	u, err := s.repo.GetByAccount(cmd.User)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.New(allerror.ErrorCodeUserNotFound, "",
				xerrors.Errorf("user %s not found: %w", cmd.User.Account(), err))
		} else {
			err = xerrors.Errorf("failed to get user: %w", err)
		}

		return err
	}

	if u.Email.Email() == cmd.Email.Email() {
		return nil
	}

	if u.Email.Email() != "" {
		e := xerrors.Errorf("user already bind another email address")
		return allerror.New(allerror.ErrorCodeEmailDuplicateBind, e.Error(), e)
	}

	err = s.oidc.VerifyBindEmail(cmd.Email.Email(), cmd.PassCode, userId)
	if err != nil && !allerror.IsUserDuplicateBind(err) {
		return allerror.New(allerror.ErrorCodeEmailVerifyFailed, "",
			xerrors.Errorf("failed to verify email: %w", err))
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
			return allerror.New(allerror.ErrorVerifyEmailFailed, "",
				xerrors.Errorf("failed to create platform user: %w", err))
		}

		u.PlatformId = user.PlatformId
		u.PlatformPwd = user.PlatformPwd
	} else {
		err = s.git.Update(userCmd)
		if err != nil {
			return allerror.New(allerror.ErrorVerifyEmailFailed, "",
				xerrors.Errorf("failed to update platform user: %w", err))
		}
	}
	// we must create git user before save
	// bcs we need save platform id&pwd
	_, err = s.repo.SaveUser(&u)

	return xerrors.Errorf("failed to save user: %w", err)
}

// getUserIdOfLogin get user id of login
func (s userService) getUserIdOfLogin(user primitive.Account) (userId string, err error) {
	loginInfo, err := s.session.FindByUser(user)
	if err != nil {
		err = xerrors.Errorf("failed to get user session: %w", err)
		return
	}

	if len(loginInfo) == 0 {
		err = xerrors.Errorf("user session not found")
		return
	}

	if loginInfo[0].UserId == "" {
		err = xerrors.Errorf("user id is empty")
		return
	}

	userId = loginInfo[0].UserId

	return
}

// PrivacyRevoke revokes the privacy settings for a user.
func (s userService) PrivacyRevoke(user primitive.Account) (string, error) {
	sessions, err := s.session.FindByUser(user)
	if err != nil {
		return "", allerror.New(allerror.ErrorCodeUserNotFound, "",
			xerrors.Errorf("failed to get user session: %w", err))
	}

	if len(sessions) == 0 {
		e := xerrors.Errorf("user session not found")
		return "", allerror.NewNoPermission(e.Error(), e)
	}

	ss := sessions[0]
	if ss.UserId == "" || ss.IdToken == "" {
		return "", allerror.New(allerror.ErrorCodeSessionNotFound, "",
			xerrors.Errorf("session info not found"))
	}

	if err = s.oidc.PrivacyRevoke(ss.UserId); err != nil {
		return "", allerror.New(allerror.ErrorFailedToRevokePrivacy, "",
			xerrors.Errorf("failed to revoke privacy: %w", err))
	}

	userInfo, err := s.repo.GetByAccount(user)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			e := xerrors.Errorf("user %s not found: %w", user.Account(), err)
			return "", allerror.New(allerror.ErrorCodeUserNotFound, "", e)
		} else {
			return "", xerrors.Errorf("failed to get user: %w", err)
		}
	}

	userInfo.RevokePrivacy()
	if _, err = s.repo.SaveUser(&userInfo); err != nil {
		return "", allerror.New(allerror.ErrorFailedToRevokePrivacy, "",
			xerrors.Errorf("failed to save user: %w", err))
	}

	action, err := ss.IdToken, s.sessionClear.ClearAllSession(user)
	if err != nil {
		action = ""
		err = allerror.New(allerror.ErrorFailedToRevokePrivacy, "",
			xerrors.Errorf("failed to clear session: %w", err))
	}

	return action, err
}

func (s userService) AgreePrivacy(user primitive.Account) error {
	userInfo, err := s.repo.GetByAccount(user)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			e := xerrors.Errorf("user %s not found: %w", user.Account(), err)
			return allerror.New(allerror.ErrorCodeUserNotFound, "", e)
		} else {
			return xerrors.Errorf("failed to get user: %w", err)
		}
	}

	userInfo.AgreePrivacy()

	_, err = s.repo.SaveUser(&userInfo)
	if err != nil {
		err = allerror.New(allerror.ErrorFailedToAgreePrivacy, "",
			xerrors.Errorf("failed to save user: %w", err))
	}

	return err
}

func (s userService) IsAgreePrivacy(user primitive.Account) (bool, error) {
	userInfo, err := s.repo.GetByAccount(user)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			return true, nil
		}

		return false, xerrors.Errorf("failed to get user: %w", err)
	}

	return userInfo.IsAgreePrivacy, nil
}

// IsOrganization checks if the given user is an organization.
//
// Parameters:
// - user: The user account to check.
//
// Returns:
// - bool: True if the user is an organization, false otherwise.
func (s userService) IsOrganization(user domain.Account) bool {
	userInfo, err := s.repo.GetByAccount(user)
	if err != nil {
		return true
	}

	return userInfo.IsOrganization()
}

// CanCreateToken checks if the given user can create a token.
//
// Parameters:
// - user: The user account to check.
//
// Returns:
// - bool: True if the user can create a token, false otherwise.
func (s userService) CanCreateToken(user domain.Account) (bool, error) {
	c, err := s.token.Count(user)
	if err != nil {
		return false, xerrors.Errorf("failed to count token: %w", err)
	}

	if c >= int64(s.config.MaxTokenPerUser) {
		return false, xerrors.Errorf("token count(now:%d max:%d) exceed", c, s.config.MaxTokenPerUser)
	}

	return true, nil
}
