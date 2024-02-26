/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package repositoryimpl

import (
	"errors"
	"fmt"

	"gorm.io/gorm/clause"

	"github.com/openmerlin/merlin-server/common/domain/crypto"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	org "github.com/openmerlin/merlin-server/organization/domain"
	"github.com/openmerlin/merlin-server/user/domain"
	"github.com/openmerlin/merlin-server/user/domain/repository"
)

// NewUserRepo creates a new user repository with the given database implementation and encryption utility.
func NewUserRepo(db postgresql.Impl, enc crypto.Encrypter) repository.User {
	userTableName = db.TableName()

	if err := postgresql.DB().AutoMigrate(&UserDO{}); err != nil {
		return nil
	}

	return &userRepoImpl{Impl: db, e: enc}
}

type userRepoImpl struct {
	postgresql.Impl
	e crypto.Encrypter
}

// AddOrg adds a new organization to the repository and returns it.
func (impl *userRepoImpl) AddOrg(o *org.Organization) (new org.Organization, err error) {
	o.Id = primitive.CreateIdentity(primitive.GetId())
	do := toOrgDO(o)

	err = impl.DB().Clauses(clause.Returning{}).Create(&do).Error
	if err != nil {
		return
	}

	return do.toOrg(impl.e)
}

// SaveOrg saves an existing organization in the repository and returns it.
func (impl *userRepoImpl) SaveOrg(o *org.Organization) (new org.Organization, err error) {
	do := toOrgDO(o)
	do.Version += 1

	tmpDo := &UserDO{}
	tmpDo.ID = o.Id.Integer()
	v := impl.DB().Model(
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
			errors.New("concurrent updating"),
		)
		return
	}

	return tmpDo.toOrg(impl.e)
}

// DeleteOrg deletes an organization from the repository by its primary key.
func (impl *userRepoImpl) DeleteOrg(o *org.Organization) error {
	tmpDo := &UserDO{}
	tmpDo.ID = o.Id.Integer()

	return impl.DeleteByPrimaryKey(
		tmpDo,
	)
}

// CheckName checks if the given name exists in the database.
// check if the name is available
// return true mean the name is available
// return false mean the name is not available
func (impl *userRepoImpl) CheckName(name primitive.Account) bool {
	var count int64
	err := impl.DB().Where(impl.EqualQuery(fieldName), name.Account()).Count(&count).Error

	if err == nil && count == 0 {
		return true
	}

	return false
}

// GetOrgByName retrieves an organization by its name.
func (impl *userRepoImpl) GetOrgByName(account primitive.Account) (r org.Organization, err error) {
	tmpDo := &UserDO{}
	tmpDo.Name = account.Account()
	tmpDo.Type = domain.UserTypeOrganization

	if err = impl.GetRecord(&tmpDo, &tmpDo); err != nil {
		return
	}

	return tmpDo.toOrg(impl.e)
}

// GetOrgByOwner retrieves organizations owned by the given account.
func (impl *userRepoImpl) GetOrgByOwner(account primitive.Account) (os []org.Organization, err error) {
	query := impl.DB().Where(
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

// GetOrgCountByOwner retrieves the count of organizations owned by the given account.
func (impl *userRepoImpl) GetOrgCountByOwner(account primitive.Account) (total int64, err error) {
	err = impl.DB().
		Where(impl.EqualQuery(fieldOwner), account.Account()).
		Where(impl.EqualQuery(fieldType), domain.UserTypeOrganization).
		Count(&total).Error

	return
}

// AddUser adds a new user to the database.
func (impl *userRepoImpl) AddUser(u *domain.User) (new domain.User, err error) {
	u.Id = primitive.CreateIdentity(primitive.GetId())
	do, err := toUserDO(u, impl.e)
	if err != nil {
		return
	}

	err = impl.DB().Clauses(clause.Returning{}).Create(&do).Error

	if err != nil {
		return
	}

	return do.toUser(impl.e)
}

// SaveUser saves the given user to the database and returns the updated user.
func (impl *userRepoImpl) SaveUser(u *domain.User) (new domain.User, err error) {
	do, err := toUserDO(u, impl.e)
	if err != nil {
		return
	}

	do.Version += 1

	tmpDo := &UserDO{}
	tmpDo.ID = u.Id.Integer()
	v := impl.DB().Model(
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
			errors.New("concurrent updating"),
		)
		return
	}

	return tmpDo.toUser(impl.e)
}

// DeleteUser deletes the given user from the database.
func (impl *userRepoImpl) DeleteUser(u *domain.User) error {
	tmpDo := &UserDO{}
	tmpDo.ID = u.Id.Integer()

	return impl.DeleteByPrimaryKey(
		tmpDo,
	)
}

// GetByAccount retrieves a user by their account information.
func (impl *userRepoImpl) GetByAccount(account domain.Account) (r domain.User, err error) {
	tmpDo := &UserDO{}
	tmpDo.Name = account.Account()
	tmpDo.Type = domain.UserTypeUser

	if err = impl.GetRecord(&tmpDo, &tmpDo); err != nil {
		return
	}

	return tmpDo.toUser(impl.e)
}

// GetUserFullname retrieves the full name of a user by their account information.
func (impl *userRepoImpl) GetUserFullname(account domain.Account) (fullname string, err error) {
	tmpDo := &UserDO{}
	tmpDo.Name = account.Account()

	if err = impl.GetRecord(&tmpDo, &tmpDo); err != nil {
		return
	}

	return tmpDo.Fullname, nil
}

// GetUserAvatarId retrieves the avatar ID of a user by their account information.
func (impl *userRepoImpl) GetUserAvatarId(account domain.Account) (id primitive.AvatarId, err error) {
	tmpDo := &UserDO{}
	tmpDo.Name = account.Account()

	if err = impl.GetRecord(&tmpDo, &tmpDo); err != nil {
		return
	}

	return primitive.CreateAvatarId(tmpDo.AvatarId), nil
}

// GetUsersAvatarId retrieves the avatar IDs of multiple users by their names.
func (impl *userRepoImpl) GetUsersAvatarId(names []string) (users []domain.User, err error) {
	query := impl.DB().Where(
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
func (impl *userRepoImpl) ListAccount(opt *repository.ListOption) (users []domain.User, count int, err error) {
	query := impl.DB()

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
