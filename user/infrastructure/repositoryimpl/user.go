package repositoryimpl

import (
	"errors"
	"fmt"

	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"gorm.io/gorm/clause"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	org "github.com/openmerlin/merlin-server/organization/domain"
	"github.com/openmerlin/merlin-server/user/domain"
	"github.com/openmerlin/merlin-server/user/domain/repository"
)

func NewUserRepo(db postgresql.Impl) repository.User {
	userTableName = db.TableName()

	if err := postgresql.DB().AutoMigrate(&UserDO{}); err != nil {
		return nil
	}

	return &userRepoImpl{db}
}

type userRepoImpl struct {
	postgresql.Impl
}

func (impl *userRepoImpl) AddOrg(o *org.Organization) (new org.Organization, err error) {
	o.Id = primitive.CreateIdentity(primitive.GetId())
	do := toOrgDO(o)

	err = impl.DB().Clauses(clause.Returning{}).Create(&do).Error
	if err != nil {
		return
	}

	new = do.toOrg()

	return
}

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

	new = tmpDo.toOrg()

	return
}

func (impl *userRepoImpl) DeleteOrg(o *org.Organization) error {
	tmpDo := &UserDO{}
	tmpDo.ID = o.Id.Integer()

	return impl.DeleteByPrimaryKey(
		tmpDo,
	)
}

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

func (impl *userRepoImpl) GetOrgByName(account primitive.Account) (r org.Organization, err error) {
	tmpDo := &UserDO{}
	tmpDo.Name = account.Account()
	tmpDo.Type = domain.UserTypeOrganization

	if err = impl.GetRecord(&tmpDo, &tmpDo); err != nil {
		return
	}

	r = tmpDo.toOrg()

	return
}

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
		os[i] = dos[i].toOrg()
	}

	return
}

func (impl *userRepoImpl) AddUser(u *domain.User) (new domain.User, err error) {
	u.Id = primitive.CreateIdentity(primitive.GetId())
	do := toUserDO(u)

	err = impl.DB().Clauses(clause.Returning{}).Create(&do).Error

	if err != nil {
		return
	}

	new = do.toUser()

	return
}

func (impl *userRepoImpl) SaveUser(u *domain.User) (new domain.User, err error) {
	do := toUserDO(u)
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

	new = tmpDo.toUser()

	return
}

func (impl *userRepoImpl) DeleteUser(u *domain.User) error {
	tmpDo := &UserDO{}
	tmpDo.ID = u.Id.Integer()

	return impl.DeleteByPrimaryKey(
		tmpDo,
	)
}

func (impl *userRepoImpl) GetByAccount(account domain.Account) (r domain.User, err error) {
	tmpDo := &UserDO{}
	tmpDo.Name = account.Account()
	tmpDo.Type = domain.UserTypeUser

	if err = impl.GetRecord(&tmpDo, &tmpDo); err != nil {
		return
	}

	r = tmpDo.toUser()

	return
}

func (impl *userRepoImpl) GetUserFullname(account domain.Account) (fullname string, err error) {
	tmpDo := &UserDO{}
	tmpDo.Name = account.Account()

	if err = impl.GetRecord(&tmpDo, &tmpDo); err != nil {
		return
	}

	return tmpDo.Fullname, nil
}

func (impl *userRepoImpl) GetUserAvatarId(account domain.Account) (id primitive.AvatarId, err error) {
	tmpDo := &UserDO{}
	tmpDo.Name = account.Account()

	if err = impl.GetRecord(&tmpDo, &tmpDo); err != nil {
		return
	}

	return primitive.CreateAvatarId(tmpDo.AvatarId), nil
}

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
		users[i] = dos[i].toUser()
	}

	return
}

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
		users[i] = dos[i].toUser()
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

	// TODO other type

	default:
		return ""
	}
}
