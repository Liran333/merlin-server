/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"

	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	orgrepository "github.com/openmerlin/merlin-server/organization/domain/repository"
	session "github.com/openmerlin/merlin-server/session/domain/repository"
	"github.com/openmerlin/merlin-server/user/domain"
	"github.com/openmerlin/merlin-server/user/domain/obs"
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
	Create(context.Context, *domain.UserCreateCmd) (UserDTO, error)
	Delete(context.Context, domain.Account) error
	RequestDelete(context.Context, domain.Account) error
	UpdateBasicInfo(context.Context, domain.Account, UpdateUserBasicInfoCmd) (UserDTO, error)
	UserInfo(context.Context, domain.Account, domain.Account) (UserInfoDTO, error)
	GetByAccount(context.Context, domain.Account, domain.Account) (UserDTO, error)
	GetOrgOrUser(context.Context, primitive.Account, primitive.Account) (UserDTO, error)
	GetUserAvatarId(context.Context, domain.Account) (AvatarDTO, error)
	GetUserFullname(context.Context, domain.Account) (string, error)
	GetUsersAvatarId(context.Context, []domain.Account) ([]AvatarDTO, error)
	HasUser(context.Context, primitive.Account) bool

	IsOrganization(context.Context, domain.Account) bool

	ListUsers(context.Context, primitive.Account) ([]UserDTO, error)

	GetPlatformUser(context.Context, domain.Account) (platform.BaseAuthClient, error)
	GetPlatformUserInfo(context.Context, domain.Account) (string, error)

	CreateToken(context.Context, *domain.TokenCreatedCmd, platform.BaseAuthClient) (TokenDTO, error)
	DeleteToken(context.Context, *domain.TokenDeletedCmd, platform.BaseAuthClient) error
	ListTokens(domain.Account) ([]TokenDTO, error)
	GetToken(context.Context, domain.Account, primitive.TokenName) (TokenDTO, error)
	VerifyToken(string, primitive.TokenPerm) (TokenDTO, error)

	SendBindEmail(*CmdToSendBindEmail) error
	VerifyBindEmail(context.Context, *CmdToVerifyBindEmail) error

	PrivacyRevoke(context.Context, primitive.Account) (string, error)
	AgreePrivacy(context.Context, primitive.Account) error
	IsAgreePrivacy(context.Context, primitive.Account) (bool, error)

	UploadAvatar(*CmdToUploadAvatar) (AvatarUrlDTO, error)
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
	obs obs.ObsService,
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
		obs:          obs,
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
	obs          obs.ObsService
}

// Create creates a new user in the system.
func (s userService) Create(ctx context.Context, cmd *domain.UserCreateCmd) (dto UserDTO, err error) {
	if cmd == nil {
		e := xerrors.Errorf("input param is empty")
		err = allerror.NewCommonRespError(e.Error(), e)
		return
	}

	if err = cmd.Validate(); err != nil {
		err = allerror.NewInvalidParam(err.Error(), xerrors.Errorf("create user cmd validate error: %w", err))
		return
	}

	if !s.repo.CheckName(ctx, cmd.Account) {
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
	u, err := s.repo.AddUser(ctx, &v)
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
func (s userService) UpdateBasicInfo(
	ctx context.Context, account domain.Account, cmd UpdateUserBasicInfoCmd) (dto UserDTO, err error) {
	user, err := s.repo.GetByAccount(ctx, account)
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

	if user, err = s.repo.SaveUser(ctx, &user); err != nil {
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
func (s userService) GetPlatformUserInfo(ctx context.Context, account domain.Account) (string, error) {
	if account == nil {
		e := xerrors.Errorf("username invalid")
		return "", allerror.New(allerror.ErrorUsernameInvalid, e.Error(), e)
	}
	// get user from db
	usernew, err := s.repo.GetByAccount(ctx, account)
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
func (s userService) GetPlatformUser(
	ctx context.Context, account domain.Account) (token platform.BaseAuthClient, err error) {
	p, err := s.GetPlatformUserInfo(ctx, account)
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
func (s userService) HasUser(ctx context.Context, acc primitive.Account) bool {
	if acc == nil {
		logrus.Errorf("username invalid")
		return false
	}

	_, err := s.repo.GetByAccount(ctx, acc)
	if err != nil {
		logrus.Errorf("user %s not found", acc.Account())
		return false
	}

	return true
}

// Delete deletes a user from the system.
func (s userService) Delete(ctx context.Context, account domain.Account) (err error) {
	u, err := s.repo.GetByAccount(ctx, account)
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
	err = s.repo.DeleteUser(ctx, &u)
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

func (s userService) RequestDelete(ctx context.Context, user domain.Account) error {
	u, err := s.repo.GetByAccount(ctx, user)
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

	_, err = s.repo.SaveUser(ctx, &u)

	return err
}

// UserInfo returns the user information for the given actor and account.
func (s userService) UserInfo(
	ctx context.Context, actor, account domain.Account) (dto UserInfoDTO, err error) {
	if dto.UserDTO, err = s.GetByAccount(ctx, actor, account); err != nil {
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
func (s userService) GetByAccount(
	ctx context.Context, actor, account domain.Account) (dto UserDTO, err error) {
	// get user
	u, err := s.repo.GetByAccount(ctx, account)
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
func (s userService) GetOrgOrUser(ctx context.Context, actor, acc primitive.Account) (dto UserDTO, err error) {
	u, err := s.repo.GetByAccount(ctx, acc)
	if err != nil && !commonrepo.IsErrorResourceNotExists(err) {
		return
	} else if err == nil {
		dto = NewUserDTO(&u, actor)
		return
	}

	o, err := s.repo.GetOrgByName(ctx, acc)
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
func (s userService) ListUsers(ctx context.Context, actor primitive.Account) (dtos []UserDTO, err error) {
	// get user
	t := domain.UserTypeUser
	u, _, err := s.repo.ListAccount(ctx, &repository.ListOption{Type: &t})
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
func (s userService) GetUserAvatarId(ctx context.Context, user domain.Account) (
	AvatarDTO, error,
) {
	var ava AvatarDTO
	a, err := s.repo.GetUserAvatarId(ctx, user)
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
		AvatarId: a.URL(),
	}, nil
}

// GetUsersAvatarId returns the avatar IDs for the given users.
func (s userService) GetUsersAvatarId(ctx context.Context, users []domain.Account) (
	dtos []AvatarDTO, err error,
) {
	names := make([]string, len(users))
	for i := range users {
		names[i] = users[i].Account()
	}

	us, err := s.repo.GetUsersAvatarId(ctx, names)
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
func (s userService) GetUserFullname(ctx context.Context, user domain.Account) (
	string, error,
) {
	name, err := s.repo.GetUserFullname(ctx, user)
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
func (s userService) CreateToken(ctx context.Context, cmd *domain.TokenCreatedCmd,
	client platform.BaseAuthClient) (token TokenDTO, err error) {
	if err = cmd.Validate(); err != nil {
		err = allerror.NewInvalidParam(err.Error(), xerrors.Errorf("create token cmd validate error: %w", err))
		return
	}

	if ok, err1 := s.CanCreateToken(cmd.Account); !ok {
		err = allerror.NewCountExceeded("token count exceed", xerrors.Errorf("create token failed :%w", err1))
		return
	}

	owner, err := s.repo.GetByAccount(ctx, cmd.Account)
	if err != nil {
		err = xerrors.Errorf("failed to get user: %w", err)
		err = allerror.New(allerror.ErrorFailedToCreateToken, "failed to create token", err)
		return
	}

	_, err = s.token.GetByName(ctx, cmd.Account, cmd.Name)
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
func (s userService) DeleteToken(
	ctx context.Context, cmd *domain.TokenDeletedCmd, client platform.BaseAuthClient) (err error) {
	if err = cmd.Validate(); err != nil {
		err = allerror.NewInvalidParam(err.Error(), xerrors.Errorf("delete token cmd validate error: %w", err))
		return
	}

	_, err = s.token.GetByName(ctx, cmd.Account, cmd.Name)
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
func (s userService) GetToken(ctx context.Context, acc domain.Account, name primitive.TokenName) (TokenDTO, error) {
	token, err := s.token.GetByName(ctx, acc, name)
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
func (s userService) VerifyBindEmail(ctx context.Context, cmd *CmdToVerifyBindEmail) error {
	userId, err := s.getUserIdOfLogin(cmd.User)
	if err != nil {
		return allerror.New(allerror.ErrorCodeUserNotFound, "",
			xerrors.Errorf("user %s not found: %w", cmd.User.Account(), err))
	}

	u, err := s.repo.GetByAccount(ctx, cmd.User)
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
		return allerror.New(allerror.ErrorCodeUserDuplicateBind, e.Error(), e)
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
	}

	if u.PlatformId == 0 {
		// create new user if user doesnot exist
		user, err := s.git.Create(userCmd)
		if err != nil {
			return allerror.New(allerror.ErrorVerifyEmailGitError, "",
				xerrors.Errorf("failed to create platform user: %w", err))
		}

		u.PlatformId = user.PlatformId
		u.PlatformPwd = user.PlatformPwd
	} else {
		err = s.git.Update(userCmd)
		if err != nil {
			return allerror.New(allerror.ErrorVerifyEmailGitError, "",
				xerrors.Errorf("failed to update platform user: %w", err))
		}
	}
	// we must create git user before save
	// bcs we need save platform id&pwd
	_, err = s.repo.SaveUser(ctx, &u)
	if err != nil {
		return xerrors.Errorf("failed to save user: %s", err)
	}

	return nil
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
func (s userService) PrivacyRevoke(ctx context.Context, user primitive.Account) (string, error) {
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

	userInfo, err := s.repo.GetByAccount(ctx, user)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			e := xerrors.Errorf("user %s not found: %w", user.Account(), err)
			return "", allerror.New(allerror.ErrorCodeUserNotFound, "", e)
		} else {
			return "", xerrors.Errorf("failed to get user: %w", err)
		}
	}

	userInfo.RevokePrivacy()
	if _, err = s.repo.SaveUser(ctx, &userInfo); err != nil {
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

func (s userService) AgreePrivacy(ctx context.Context, user primitive.Account) error {
	userInfo, err := s.repo.GetByAccount(ctx, user)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			e := xerrors.Errorf("user %s not found: %w", user.Account(), err)
			return allerror.New(allerror.ErrorCodeUserNotFound, "", e)
		} else {
			return xerrors.Errorf("failed to get user: %w", err)
		}
	}

	userInfo.AgreePrivacy()

	_, err = s.repo.SaveUser(ctx, &userInfo)
	if err != nil {
		err = allerror.New(allerror.ErrorFailedToAgreePrivacy, "",
			xerrors.Errorf("failed to save user: %w", err))
	}

	return err
}

func (s userService) IsAgreePrivacy(ctx context.Context, user primitive.Account) (bool, error) {
	userInfo, err := s.repo.GetByAccount(ctx, user)
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
func (s userService) IsOrganization(ctx context.Context, user domain.Account) bool {
	userInfo, err := s.repo.GetByAccount(ctx, user)
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

// UploadAvatar upload avatar.
func (s userService) UploadAvatar(cmd *CmdToUploadAvatar) (AvatarUrlDTO, error) {
	if cmd == nil {
		e := xerrors.Errorf("input param is empty")
		err := allerror.NewCommonRespError(e.Error(), e)

		return AvatarUrlDTO{}, err
	}

	avatar := domain.AvatarInfo{
		Path:        s.config.ObsPath,
		CdnEndpoint: s.config.CdnEndpoint,
		Account:     cmd.User,
		FileName:    cmd.FileName,
	}

	err := s.obs.CreateObject(cmd.Image, s.config.ObsBucket, avatar.GetObsPath())
	if err != nil {
		return AvatarUrlDTO{}, allerror.New(allerror.ErrorCodeFileUploadFailed, "",
			xerrors.Errorf("failed to upload avatar: %w", err))
	}

	dto := AvatarUrlDTO{
		URL: avatar.GetAvatarURL(),
	}

	return dto, nil
}
