/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package repositoryimpl

import (
	"errors"

	"gorm.io/gorm/clause"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/organization/domain"
	"github.com/openmerlin/merlin-server/organization/domain/repository"
)

// NewMemberRepo creates a new MemberRepo instance with the given postgresql.Impl as the database implementation.
func NewMemberRepo(db postgresql.Impl) repository.OrgMember {
	if err := postgresql.DB().Table(db.TableName()).AutoMigrate(&Member{}); err != nil {
		return nil
	}

	return &memberRepoImpl{Impl: db}
}

type memberRepoImpl struct {
	postgresql.Impl
}

// Add adds a new org member to the database and returns the created org member.
func (impl *memberRepoImpl) Add(o *domain.OrgMember) (new domain.OrgMember, err error) {
	o.Id = primitive.CreateIdentity(primitive.GetId())
	do := toMemberDoc(o)

	err = impl.DB().Clauses(clause.Returning{}).Create(&do).Error
	if err != nil {
		return
	}

	new = toOrgMember(&do)

	return
}

// Save updates an existing org member in the database and returns the updated org member.
func (impl *memberRepoImpl) Save(o *domain.OrgMember) (new domain.OrgMember, err error) {
	do := toMemberDoc(o)
	do.Version += 1

	tmpDo := &Member{}
	tmpDo.ID = o.Id.Integer()
	v := impl.DB().Model(
		tmpDo,
	).Where(
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

	new = toOrgMember(tmpDo)

	return
}

// Delete deletes an org member from the database by its primary key.
func (impl *memberRepoImpl) Delete(o *domain.OrgMember) (err error) {
	tmpDo := &Member{}
	tmpDo.ID = o.Id.Integer()

	return impl.DeleteByPrimaryKey(
		tmpDo,
	)
}

// GetByOrg retrieves a list of members by organization name.
func (impl *memberRepoImpl) GetByOrg(name string) (
	members []domain.OrgMember, err error,
) {
	var v []Member

	query := impl.DB()

	query = query.Where(impl.EqualQuery(fieldOrg), name)

	err = query.Find(&v).Error
	if err != nil || len(v) == 0 {
		return nil, err
	}

	members = make([]domain.OrgMember, len(v))
	for i := range v {
		members[i] = toOrgMember(&v[i])
	}

	return
}

// DeleteByOrg deletes members by organization name.
func (impl *memberRepoImpl) DeleteByOrg(name primitive.Account) (
	err error,
) {
	return impl.DB().Where(impl.EqualQuery(fieldOrg), name.Account()).Delete(&Member{}).Error
}

// GetByOrgAndUser retrieves a member by organization and user names.
func (impl *memberRepoImpl) GetByOrgAndUser(org, user string) (
	member domain.OrgMember, err error,
) {
	var v Member

	v.Orgname = org
	v.Username = user

	if err = impl.GetRecord(&v, &v); err != nil {
		return
	}

	member = toOrgMember(&v)

	return
}

// GetByOrgAndRole retrieves members by organization and role.
func (impl *memberRepoImpl) GetByOrgAndRole(org string, role domain.OrgRole) (members []domain.OrgMember, err error) {
	var v []Member

	query := Member{}
	query.Orgname = org
	query.Role = role

	err = impl.DB().Where(&query).Find(&v).Error
	if err != nil || len(v) == 0 {
		return
	}

	members = make([]domain.OrgMember, len(v))
	for i := range v {
		members[i] = toOrgMember(&v[i])
	}

	return
}

// GetByUser retrieves members by user name.
func (impl *memberRepoImpl) GetByUser(name string) (
	members []domain.OrgMember, err error,
) {
	var v []Member

	query := Member{}
	query.Username = name

	err = impl.DB().Where(&query).Find(&v).Error
	if err != nil || len(v) == 0 {
		return
	}

	members = make([]domain.OrgMember, len(v))
	for i := range v {
		members[i] = toOrgMember(&v[i])
	}

	return
}
