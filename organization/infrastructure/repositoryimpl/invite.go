/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package repositoryimpl

import (
	"errors"
	"gorm.io/gorm"

	"gorm.io/gorm/clause"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/organization/domain"
	"github.com/openmerlin/merlin-server/organization/domain/repository"
)

// NewInviteRepo creates a new instance of inviteRepoImpl.
func NewInviteRepo(db postgresql.Impl) repository.Approve {
	if err := postgresql.DB().Table(db.TableName()).AutoMigrate(&Approve{}); err != nil {
		return nil
	}

	return &inviteRepoImpl{Impl: db}
}

type inviteRepoImpl struct {
	postgresql.Impl
}

// ListInvitation lists the invitations based on the provided command.
func (impl *inviteRepoImpl) ListInvitation(cmd *domain.OrgInvitationListCmd) (approves []domain.Approve, err error) {
	var v []Approve

	query := impl.DB().Where(impl.EqualQuery(fieldType), domain.InviteTypeInvite)

	if cmd.Org != nil {
		query = query.Where(impl.EqualQuery(fieldOrg), cmd.Org.Account())
	}

	if cmd.Invitee != nil {
		query = query.Where(impl.EqualQuery(fieldInvitee), cmd.Invitee.Account())
	}

	if cmd.Inviter != nil {
		query = query.Where(impl.EqualQuery(fieldInviter), cmd.Inviter.Account())
	}

	if cmd.Status != "" {
		query = query.Where(impl.EqualQuery(fieldStatus), cmd.Status)
	}

	err = query.Find(&v).Error
	if err != nil || len(v) == 0 {
		return nil, err
	}

	approves = make([]domain.Approve, len(v))
	for i := range v {
		approves[i] = toApprove(&v[i])
	}

	return

}

// AddInvite adds a new invite to the database.
func (impl *inviteRepoImpl) AddInvite(o *domain.Approve) (new domain.Approve, err error) {
	// Build query to check for an existing invitation
	query := impl.DB().Model(&Approve{}).
		Where(impl.EqualQuery(fieldUser), o.Username).
		Where(impl.EqualQuery(fieldOrg), o.OrgName).
		Where(impl.EqualQuery(fieldInviter), o.Inviter).
		Where(impl.EqualQuery(fieldStatus), o.Status)

	// Attempt to find an existing record
	var existingInvite Approve
	err = query.First(&existingInvite).Error
	if err == nil {
		// Found an existing record, delete it
		if deleteErr := impl.DB().Delete(&existingInvite).Error; deleteErr != nil {
			return domain.Approve{}, deleteErr
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// An error occurred other than "record not found", return the error
		return domain.Approve{}, err
	}

	// Record not found or deleted, proceed to add new one
	o.Id = primitive.CreateIdentity(primitive.GetId())
	do := toApproveDoc(o)

	err = impl.DB().Clauses(clause.Returning{}).Create(&do).Error
	if err != nil {
		return domain.Approve{}, err
	}

	new = toApprove(&do)
	return
}

// SaveInvite saves an existing invite in the database.
func (impl *inviteRepoImpl) SaveInvite(o *domain.Approve) (new domain.Approve, err error) {
	do := toApproveDoc(o)
	do.Version += 1

	tmpDo := &Approve{}
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
	new = toApprove(tmpDo)

	return
}

// DeleteInviteAndReqByOrg deletes invite and request records associated with the given organization account.
func (impl *inviteRepoImpl) DeleteInviteAndReqByOrg(acc primitive.Account) error {
	return impl.DB().Where(impl.EqualQuery(fieldOrg), acc.Account()).Delete(&Approve{}).Error
}

// AddRequest adds a new member request and returns the created request.
func (impl *inviteRepoImpl) AddRequest(r *domain.MemberRequest) (new domain.MemberRequest, err error) {
	r.Id = primitive.CreateIdentity(primitive.GetId())
	do := toRequestDoc(r)

	err = impl.DB().Clauses(clause.Returning{}).Create(&do).Error
	if err != nil {
		return
	}

	new = toMemberRequest(&do)

	return
}

// SaveRequest updates an existing member request and returns the updated request.
func (impl *inviteRepoImpl) SaveRequest(r *domain.MemberRequest) (new domain.MemberRequest, err error) {
	do := toRequestDoc(r)
	do.Version += 1

	tmpDo := &Approve{}
	tmpDo.ID = r.Id.Integer()
	v := impl.DB().Model(
		tmpDo,
	).Clauses(clause.Returning{}).Where(
		impl.EqualQuery(fieldVersion), r.Version,
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

	return toMemberRequest(tmpDo), nil

}

// ListRequests lists member requests based on the provided command criteria.
func (impl *inviteRepoImpl) ListRequests(cmd *domain.OrgMemberReqListCmd) (rs []domain.MemberRequest, err error) {
	var v []Approve

	query := impl.DB().Where(impl.EqualQuery(fieldType), domain.InviteTypeRequest)

	if cmd.Org != nil {
		query = query.Where(impl.EqualQuery(fieldOrg), cmd.Org.Account())
	}

	if cmd.Requester != nil {
		query = query.Where(impl.EqualQuery(fieldInvitee), cmd.Requester.Account())
	}

	if cmd.Status != "" {
		query = query.Where(impl.EqualQuery(fieldStatus), cmd.Status)
	}

	err = query.Find(&v).Error
	if err != nil || len(v) == 0 {
		return nil, err
	}

	rs = make([]domain.MemberRequest, len(v))
	for i := range v {
		rs[i] = toMemberRequest(&v[i])
	}

	return
}
