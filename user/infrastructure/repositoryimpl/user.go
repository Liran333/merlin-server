/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package repositoryimpl

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/xerrors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/openmerlin/merlin-server/common/domain/crypto"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	org "github.com/openmerlin/merlin-server/organization/domain"
	"github.com/openmerlin/merlin-server/user/domain"
	"github.com/openmerlin/merlin-server/user/domain/repository"
	"github.com/sirupsen/logrus"
)

// NewUserRepo creates a new user repository with the given database implementation and encryption utility.
func NewUserRepo(db postgresql.Impl, enc crypto.Encrypter) repository.User {
	userTableName = db.TableName()
	logrus.Infof("user table name: %s", userTableName)

	if err := postgresql.DB().AutoMigrate(&UserDO{}); err != nil {
		logrus.Fatalf("failed to auto migrate user table: %v", err)
		return nil
	}

	return &userRepoImpl{Impl: db, e: enc}
}

type userRepoImpl struct {
	postgresql.Impl
	e crypto.Encrypter
}

// AddOrg adds a new organization to the repository and returns it.
func (impl *userRepoImpl) AddOrg(ctx context.Context, o *org.Organization) (new org.Organization, err error) {
	o.Id = primitive.CreateIdentity(primitive.GetId())
	do := toOrgDO(o)

	err = impl.WithContext(ctx).Clauses(clause.Returning{}).Create(&do).Error
	if err != nil {
		return
	}

	return do.toOrg(impl.e)
}

// SaveOrg saves an existing organization in the repository and returns it.
func (impl *userRepoImpl) SaveOrg(ctx context.Context, o *org.Organization) (new org.Organization, err error) {
	do := toOrgDO(o)
	do.Version += 1

	tmpDo := &UserDO{}
	tmpDo.ID = o.Id.Integer()
	v := impl.WithContext(ctx).Model(
		tmpDo,
	).Clauses(clause.Returning{}).Where(
		impl.EqualQuery(fieldVersion), o.Version,
	).Select(`*`).Omit("created_at").Updates(&do) // should not update created_at

	if v.Error != nil {
		err = v.Error
		return
	}

	if v.RowsAffected == 0 {
		err = commonrepo.NewErrorConcurrentUpdating(
			xerrors.Errorf("concurrent updating: %w", err),
		)
		return
	}

	return tmpDo.toOrg(impl.e)
}

// DeleteOrg deletes an organization from the repository by its primary key.
func (impl *userRepoImpl) DeleteOrg(ctx context.Context, o *org.Organization) error {
	tmpDo := &UserDO{}
	tmpDo.ID = o.Id.Integer()

	return impl.DeleteByPrimaryKey(
		ctx,
		tmpDo,
	)
}

// CheckName checks if the given name exists in the database.
// check if the name is available
// return true mean the name is available
// return false mean the name is not available
func (impl *userRepoImpl) CheckName(ctx context.Context, name primitive.Account) bool {
	var count int64
	err := impl.WithContext(ctx).Where(fmt.Sprintf(`LOWER(%s) = ?`, fieldName),
		strings.ToLower(name.Account())).Count(&count).Error

	if err == nil && count == 0 {
		return true
	}

	return false
}

// GetOrgByName retrieves an organization by its name.
func (impl *userRepoImpl) GetOrgByName(ctx context.Context,
	account primitive.Account) (r org.Organization, err error) {
	tmpDo := &UserDO{}
	tmpDo.Name = account.Account()
	tmpDo.Type = domain.UserTypeOrganization

	if err = impl.GetRecord(ctx, &tmpDo, &tmpDo); err != nil {
		return
	}

	return tmpDo.toOrg(impl.e)
}

// GetOrgByOwner retrieves organizations owned by the given account.
func (impl *userRepoImpl) GetOrgByOwner(
	ctx context.Context, account primitive.Account) (os []org.Organization, err error) {
	query := impl.WithContext(ctx).Where(
		impl.EqualQuery(fieldOwner), account.Account(),
	).Where(impl.EqualQuery(fieldType), domain.UserTypeOrganization)

	var dos []UserDO

	err = query.Find(&dos).Error
	if err != nil || len(dos) == 0 {
		return nil, err
	}

	os = make([]org.Organization, len(dos))
	for i := range dos {
		os[i], err = dos[i].toOrg(impl.e)
		if err != nil {
			return
		}
	}

	return
}

// GetOrgList retrieves organizations owned by organization id array or owner.
func (impl *userRepoImpl) GetOrgList(
	ctx context.Context, opt *repository.ListOrgOption) (os []org.Organization, err error) {
	query := impl.WithContext(ctx)
	if len(opt.OrgIDs) > 0 {
		query = query.Where(impl.InFilter(fieldID), opt.OrgIDs)
	}
	if opt.Owner != nil {
		query = query.Where(impl.EqualQuery(fieldOwner), opt.Owner.Account())
	}

	var dos []UserDO
	err = query.Where(impl.EqualQuery(fieldType), domain.UserTypeOrganization).Find(&dos).Error
	if err != nil || len(dos) == 0 {
		return nil, err
	}

	os = make([]org.Organization, len(dos))
	for i := range dos {
		os[i], err = dos[i].toOrg(impl.e)
		if err != nil {
			return
		}
	}

	return
}

func (impl *userRepoImpl) GetOrgPageList(
	ctx context.Context, opt *repository.ListPageOrgOption) ([]org.Organization, int, error) {
	query := impl.WithContext(ctx)
	if len(opt.OrgIDs) > 0 {
		query = query.Where(impl.InFilter(fieldID), opt.OrgIDs)
	}
	if opt.Owner != nil {
		query = query.Where(impl.EqualQuery(fieldOwner), opt.Owner.Account())
	}
	query = query.Where(impl.EqualQuery(fieldType), 1)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := 0
	if opt.PageNum > 0 && opt.PageSize > 0 {
		offset = (opt.PageNum - 1) * opt.PageSize
	}
	if offset > 0 {
		query = query.Limit(opt.PageSize).Offset(offset)
	} else {
		query = query.Limit(opt.PageSize)
	}
	var dos []UserDO
	err := query.Find(&dos).Error
	if err != nil || len(dos) == 0 {
		return nil, 0, err
	}
	r := make([]org.Organization, len(dos))
	for i := range dos {
		r[i], err = dos[i].toOrg(impl.e)
		if err != nil {
			return nil, 0, err
		}
	}
	return r, int(total), nil
}

// GetOrgCountByOwner retrieves the count of organizations owned by the given account.
func (impl *userRepoImpl) GetOrgCountByOwner(ctx context.Context, account primitive.Account) (total int64, err error) {
	err = impl.WithContext(ctx).
		Where(impl.EqualQuery(fieldOwner), account.Account()).
		Where(impl.EqualQuery(fieldType), domain.UserTypeOrganization).
		Count(&total).Error

	return
}

// AddUser adds a new user to the database.
func (impl *userRepoImpl) AddUser(ctx context.Context, u *domain.User) (new domain.User, err error) {
	u.Id = primitive.CreateIdentity(primitive.GetId())
	do, err := toUserDO(u, impl.e)
	if err != nil {
		return
	}

	err = impl.DB().WithContext(ctx).Clauses(clause.Returning{}).Create(&do).Error

	if err != nil {
		return
	}

	return do.toUser(impl.e)
}

// SaveUser saves the given user to the database and returns the updated user.
func (impl *userRepoImpl) SaveUser(ctx context.Context, u *domain.User) (new domain.User, err error) {
	do, err := toUserDO(u, impl.e)
	if err != nil {
		return
	}

	do.Version += 1

	tmpDo := &UserDO{}
	tmpDo.ID = u.Id.Integer()
	v := impl.WithContext(ctx).Model(
		tmpDo,
	).Clauses(clause.Returning{}).Where(
		impl.EqualQuery(fieldVersion), u.Version,
	).Select(`*`).Omit("created_at").Updates(&do) // should not update created_at

	if v.Error != nil {
		err = v.Error
		return
	}

	if v.RowsAffected == 0 {
		err = commonrepo.NewErrorConcurrentUpdating(
			xerrors.Errorf("concurrent updating: %w", err),
		)
		return
	}

	return tmpDo.toUser(impl.e)
}

// DeleteUser deletes the given user from the database.
func (impl *userRepoImpl) DeleteUser(ctx context.Context, u *domain.User) error {
	tmpDo := &UserDO{}
	tmpDo.ID = u.Id.Integer()

	return impl.DeleteByPrimaryKey(
		ctx,
		tmpDo,
	)
}

// GetByAccount retrieves a user by their account information.
func (impl *userRepoImpl) GetByAccount(ctx context.Context, account domain.Account) (r domain.User, err error) {
	tmpDo := &UserDO{}
	// note: gorm struct query will ignore zero value field
	// so using Where instead of impl.GetRecord
	query := impl.WithContext(ctx).Where(
		impl.EqualQuery(fieldName), account.Account(),
	).Where(impl.EqualQuery(fieldType), domain.UserTypeUser)

	err = query.First(&tmpDo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = commonrepo.NewErrorResourceNotExists(xerrors.Errorf("user not found: %w", err))
			return
		}

		return
	}

	return tmpDo.toUser(impl.e)
}

// GetUserFullname retrieves the full name of a user by their account information.
func (impl *userRepoImpl) GetUserFullname(ctx context.Context, account domain.Account) (fullname string, err error) {
	tmpDo := &UserDO{}
	tmpDo.Name = account.Account()

	if err = impl.GetRecord(ctx, &tmpDo, &tmpDo); err != nil {
		return
	}

	return tmpDo.Fullname, nil
}

// GetUserAvatarId retrieves the avatar ID of a user by their account information.
func (impl *userRepoImpl) GetUserAvatarId(
	ctx context.Context, account domain.Account) (id primitive.Avatar, err error) {
	tmpDo := &UserDO{}
	tmpDo.Name = account.Account()

	if err = impl.GetRecord(ctx, &tmpDo, &tmpDo); err != nil {
		return
	}

	return primitive.CreateAvatar(tmpDo.AvatarId), nil
}

// GetUsersAvatarId retrieves the avatar IDs of multiple users by their names.
func (impl *userRepoImpl) GetUsersAvatarId(ctx context.Context, names []string) (users []domain.User, err error) {
	query := impl.WithContext(ctx).Where(
		impl.InFilter(fieldName), names,
	)

	var dos []UserDO

	err = query.Find(&dos).Error
	if err != nil || len(dos) == 0 {
		return nil, err
	}

	users = make([]domain.User, len(dos))
	for i := range dos {
		users[i], err = dos[i].toUser(impl.e)
		if err != nil {
			return
		}
	}

	return
}

// ListAccount lists users based on the provided ListOption.
func (impl *userRepoImpl) ListAccount(
	ctx context.Context, opt *repository.ListOption) (users []domain.User, count int, err error) {
	query := impl.WithContext(ctx)

	if opt == nil {
		err = fmt.Errorf("list option is nil")
		return
	}

	if opt.Name != "" {
		filter, arg := impl.LikeFilter(fieldName, opt.Name)

		query = query.Where(filter, arg)
	}

	if opt.Owner != nil {
		query = query.Where(impl.EqualQuery(fieldOwner), opt.Owner.Account())
	}

	if opt.Type != nil {
		query = query.Where(impl.EqualQuery(fieldType), *opt.Type)
	}

	if v := impl.order(opt.SortType); v != "" {
		query = query.Order(v)
	}

	var dos []UserDO

	err = query.Find(&dos).Error
	if err != nil || len(dos) == 0 {
		return nil, 0, nil
	}

	users = make([]domain.User, len(dos))
	for i := range dos {
		users[i], err = dos[i].toUser(impl.e)
		if err != nil {
			return
		}
	}

	return
}

func (impl *userRepoImpl) SearchUser(ctx context.Context, opt *repository.ListOption) ([]domain.User, int, error) {
	db := impl.WithContext(ctx)

	if opt == nil {
		err := fmt.Errorf("list option is nil")
		return nil, 0, err
	}

	if opt.Name == "" {
		err := fmt.Errorf("search key is empty")
		return nil, 0, err
	}

	queryName, argName := impl.LikeFilter(fieldName, opt.Name)

	if !opt.ExcludeFullname {
		_, argName2 := impl.LikeFilter(fieldFullname, opt.Name)
		db = db.Where(fieldName+" ilike ? OR "+fieldFullname+
			" ilike ?", argName, argName2).Session(&gorm.Session{})
	} else {
		db = db.Where(queryName, argName)
	}

	db = db.Where(impl.EqualQuery(fieldType), domain.UserTypeUser)

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if b, offset := opt.Pagination(); b {
		if offset > 0 {
			db = db.Limit(opt.CountPerPage).Offset(offset)
		} else {
			db = db.Limit(opt.CountPerPage)
		}
	}

	var dos []UserDO

	if err := db.Find(&dos).Error; err != nil {
		return nil, 0, err
	}

	users := make([]domain.User, len(dos))

	for i := range dos {
		u, err := dos[i].toUser(impl.e)
		if err != nil {
			return nil, 0, err
		}
		users[i] = u
	}

	return users, int(total), nil
}

func (impl *userRepoImpl) SearchOrg(ctx context.Context, opt *repository.ListOption) ([]org.Organization, int, error) {
	db := impl.WithContext(ctx)

	if opt == nil {
		err := fmt.Errorf("list option is nil")
		return nil, 0, err
	}

	if opt.Name == "" {
		err := fmt.Errorf("search key is empty")
		return nil, 0, err
	}

	queryName, argName := impl.LikeFilter(fieldName, opt.Name)

	if !opt.ExcludeFullname {
		_, argName2 := impl.LikeFilter(fieldFullname, opt.Name)
		db = db.Where(fieldName+" ilike ? OR "+fieldFullname+
			" ilike ?", argName, argName2).Session(&gorm.Session{})
	} else {
		db = db.Where(queryName, argName)
	}

	db = db.Where(impl.EqualQuery(fieldType), domain.UserTypeOrganization)

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if b, offset := opt.Pagination(); b {
		if offset > 0 {
			db = db.Limit(opt.CountPerPage).Offset(offset)
		} else {
			db = db.Limit(opt.CountPerPage)
		}
	}

	var dos []UserDO

	if err := db.Find(&dos).Error; err != nil {
		return nil, 0, err
	}

	orgs := make([]org.Organization, len(dos))
	for i := range dos {
		o, err := dos[i].toOrg(impl.e)
		if err != nil {
			return nil, 0, err
		}
		orgs[i] = o
	}

	return orgs, int(total), nil
}

func (impl *userRepoImpl) order(t primitive.SortType) string {
	if t == nil {
		return ""
	}

	switch t.SortType() {
	case primitive.SortByAlphabetical:
		return fieldName

	case primitive.SortByRecentlyUpdated:
		return impl.OrderByDesc(fieldUpdatedAt)

	case primitive.SortByRecentlyCreated:
		return impl.OrderByDesc(fieldCreatedAt)

	default:
		return ""
	}
}
