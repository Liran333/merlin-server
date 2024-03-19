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
	"github.com/openmerlin/merlin-server/organization/domain"
	"github.com/openmerlin/merlin-server/organization/domain/repository"
	userapp "github.com/openmerlin/merlin-server/user/app"
	userrepo "github.com/openmerlin/merlin-server/user/domain/repository"
	"github.com/openmerlin/merlin-server/user/infrastructure/git"
	"github.com/openmerlin/merlin-server/utils"
)

func errOrgNotFound(msg string) error {
	if msg == "" {
		msg = "not found"
	}

	return allerror.NewNotFound(allerror.ErrorCodeOrganizationNotFound, msg)
}

// OrgService is an interface that defines the methods for organization-related operations.
type OrgService interface {
	Create(*domain.OrgCreatedCmd) (userapp.UserDTO, error)
	Delete(*domain.OrgDeletedCmd) error
	UpdateBasicInfo(*domain.OrgUpdatedBasicInfoCmd) (userapp.UserDTO, error)
	GetByAccount(primitive.Account) (userapp.UserDTO, error)
	GetOrgOrUser(primitive.Account, primitive.Account) (userapp.UserDTO, error)
	ListAccount(*userrepo.ListOption) ([]userapp.UserDTO, error)

	CheckName(primitive.Account) bool
	GetByOwner(primitive.Account, primitive.Account) ([]userapp.UserDTO, error)
	GetByUser(primitive.Account, primitive.Account) ([]userapp.UserDTO, error)
	List(*OrgListOptions) ([]userapp.UserDTO, error)
	InviteMember(*domain.OrgInviteMemberCmd) (ApproveDTO, error)
	RequestMember(*domain.OrgRequestMemberCmd) (MemberRequestDTO, error)
	CancelReqMember(*domain.OrgCancelRequestMemberCmd) (MemberRequestDTO, error)
	ApproveRequest(*domain.OrgApproveRequestMemberCmd) (MemberRequestDTO, error)
	AcceptInvite(*domain.OrgAcceptInviteCmd) (ApproveDTO, error)
	RevokeInvite(*domain.OrgRemoveInviteCmd) (ApproveDTO, error)
	ListMemberReq(*domain.OrgMemberReqListCmd) ([]MemberRequestDTO, error)
	ListInvitationByInvitee(primitive.Account, primitive.Account, domain.ApproveStatus) ([]ApproveDTO, error)
	ListInvitationByInviter(primitive.Account, primitive.Account, domain.ApproveStatus) ([]ApproveDTO, error)
	ListInvitationByOrg(primitive.Account, primitive.Account, domain.ApproveStatus) ([]ApproveDTO, error)
	AddMember(*domain.OrgAddMemberCmd) error
	RemoveMember(*domain.OrgRemoveMemberCmd) error
	EditMember(*domain.OrgEditMemberCmd) (MemberDTO, error)
	ListMember(primitive.Account) ([]MemberDTO, error)
	GetMemberByUserAndOrg(primitive.Account, primitive.Account) (MemberDTO, error)
}

// NewOrgService creates a new instance of the OrgService.
func NewOrgService(
	user userapp.UserService,
	repo userrepo.User,
	member repository.OrgMember,
	invite repository.Approve,
	perm *permService,
	cfg *domain.Config,
	git git.User,
) OrgService {
	return &orgService{
		user:             user,
		repo:             repo,
		member:           member,
		perm:             perm,
		invite:           invite,
		defaultRole:      primitive.CreateRole(cfg.DefaultRole),
		inviteExpiry:     cfg.InviteExpiry,
		MaxCountPerOwner: cfg.MaxCountPerOwner,
		git:              git,
	}
}

type orgService struct {
	MaxCountPerOwner int64
	inviteExpiry     int64
	defaultRole      primitive.Role
	user             userapp.UserService
	repo             userrepo.User
	member           repository.OrgMember
	invite           repository.Approve
	perm             *permService
	git              git.User
}

// Create creates a new organization with the given command and returns the created organization as a UserDTO.
func (org *orgService) Create(cmd *domain.OrgCreatedCmd) (o userapp.UserDTO, err error) {
	orgTemp, err := cmd.ToOrg()
	if err != nil {
		return
	}

	if !org.repo.CheckName(cmd.Name) {
		err = allerror.NewInvalidParam(fmt.Sprintf("name %s is already been taken", cmd.Name.Account()))
		return
	}

	if err = org.orgCountCheck(cmd.Owner); err != nil {
		return
	}

	owner, err := org.repo.GetByAccount(cmd.Owner)
	if err != nil {
		logrus.Error(err)
		err = allerror.NewInvalidParam("failed to get owner info")
		return
	}

	pl, err := org.user.GetPlatformUser(orgTemp.Owner)
	if err != nil {
		err = allerror.NewInvalidParam(fmt.Sprintf("failed to get platform user, %s", err))
		return
	}

	err = pl.CreateOrg(orgTemp)
	if err != nil {
		err = allerror.NewInvalidParam(fmt.Sprintf("failed to create org, %s", err))
		return
	}

	orgTemp.DefaultRole = org.defaultRole
	orgTemp.AllowRequest = false
	orgTemp.OwnerId = owner.Id

	*orgTemp, err = org.repo.AddOrg(orgTemp)
	if err != nil {
		err = allerror.NewInvalidParam(fmt.Sprintf("failed to create to org, %s", err))
		_ = pl.DeleteOrg(cmd.Name)
		return
	}

	_, err = org.member.Add(&domain.OrgMember{
		OrgName:  cmd.Name,
		OrgId:    orgTemp.Id,
		Username: cmd.Owner,
		UserId:   owner.Id,
		Role:     primitive.NewAdminRole(),
	})
	if err != nil {
		err = allerror.NewInvalidParam(fmt.Sprintf("failed to save org member, %s", err))
		_ = pl.DeleteOrg(cmd.Name)
		return
	}

	o = ToDTO(orgTemp)

	return
}

func (org *orgService) orgCountCheck(owner primitive.Account) error {
	total, err := org.repo.GetOrgCountByOwner(owner)
	if err != nil {
		return err
	}

	if total >= org.MaxCountPerOwner {
		return allerror.NewCountExceeded("org count exceed")
	}

	return nil
}

// GetByAccount retrieves an organization by its account and returns it as a UserDTO.
func (org *orgService) GetByAccount(acc primitive.Account) (dto userapp.UserDTO, err error) {
	o, err := org.repo.GetOrgByName(acc)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("org %s not found", acc.Account()))
		}

		return
	}

	dto = ToDTO(&o)
	return
}

// GetOrgOrUser retrieves either an organization or a user by their account and returns it as a UserDTO.
func (org *orgService) GetOrgOrUser(actor, acc primitive.Account) (dto userapp.UserDTO, err error) {
	u, err := org.repo.GetByAccount(acc)
	if err != nil && !commonrepo.IsErrorResourceNotExists(err) {
		return
	} else if err == nil {
		dto = userapp.NewUserDTO(&u, actor)
		return
	}

	o, err := org.repo.GetOrgByName(acc)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.New(allerror.ErrorCodeUserNotFound, fmt.Sprintf("org %s not found", acc.Account()))
		}

		return
	}

	dto = ToDTO(&o)
	return
}

// ListAccount lists organizations based on the provided options and returns them as a slice of UserDTOs.
func (org *orgService) ListAccount(opt *userrepo.ListOption) (dtos []userapp.UserDTO, err error) {
	return
}

// Delete deletes an organization based on the provided command and returns an error if any occurs.
func (org *orgService) Delete(cmd *domain.OrgDeletedCmd) error {
	err := org.perm.Check(cmd.Actor, cmd.Name, primitive.ObjTypeOrg, primitive.ActionDelete)
	if err != nil {
		return err
	}
	o, err := org.repo.GetOrgByName(cmd.Name)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		}

		return err
	}

	pl, err := org.user.GetPlatformUser(o.Owner)
	if err != nil {
		logrus.Error(err)

		return allerror.NewInvalidParam("failed to get platform user")
	}

	can, err := pl.CanDelete(cmd.Name)
	if err != nil {
		logrus.Error(err)

		return allerror.NewInvalidParam(fmt.Sprintf("%s can't delete the org", cmd.Name.Account()))
	}

	if !can {
		return allerror.New(allerror.ErrorCodeOrgExistResource,
			"can't delete the organization, while some repos still existed")
	}

	err = org.member.DeleteByOrg(o.Account)
	if err != nil {
		logrus.Error(err)

		return allerror.New(allerror.ErrorBaseCase, "failed to delete org member")
	}

	err = org.invite.DeleteInviteAndReqByOrg(o.Account)
	if err != nil {
		logrus.Error(err)

		return allerror.New(allerror.ErrorBaseCase, "failed to delete org invite")
	}

	err = org.git.DeleteOrg(o.Account)
	if err != nil {
		logrus.Error(err)

		return allerror.New(allerror.ErrorBaseCase, "failed to delete git org")
	}

	return org.repo.DeleteOrg(&o)
}

// UpdateBasicInfo updates the basic information of an organization based on the provided command
// and returns the updated organization as a UserDTO.
func (org *orgService) UpdateBasicInfo(cmd *domain.OrgUpdatedBasicInfoCmd) (dto userapp.UserDTO, err error) {
	if cmd == nil {
		err = allerror.NewInvalidParam("cmd is nil")
		return
	}

	if err = cmd.Validate(); err != nil {
		return
	}

	err = org.perm.Check(cmd.Actor, cmd.OrgName, primitive.ObjTypeOrg, primitive.ActionWrite)
	if err != nil {
		return
	}

	o, err := org.repo.GetOrgByName(cmd.OrgName)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("org %s not found", cmd.OrgName.Account()))
		}

		return
	}

	change := cmd.ToOrg(&o)
	if change {
		o, err = org.repo.SaveOrg(&o)
		if err != nil {
			err = fmt.Errorf("failed to save org, %w", err)
			return
		}
		dto = ToDTO(&o)
		return
	}
	err = allerror.NewInvalidParam("nothing changed")
	return
}

// GetByOwner retrieves organizations owned by the specified account and returns them as a slice of UserDTOs.
func (org *orgService) GetByOwner(actor, acc primitive.Account) (orgs []userapp.UserDTO, err error) {
	if acc == nil {
		err = fmt.Errorf("account is nil")
		return
	}

	orgs, err = org.List(&OrgListOptions{
		Owner: acc,
	})
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		}
	}

	return
}

// GetByUser retrieves organizations associated with a user.
func (org *orgService) GetByUser(actor, acc primitive.Account) (orgs []userapp.UserDTO, err error) {
	if acc == nil {
		err = fmt.Errorf("account is nil")
		return
	}

	members, err := org.member.GetByUser(acc.Account())
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		}

		return
	}

	orgs = make([]userapp.UserDTO, len(members))
	for i := range members {
		o, e := org.repo.GetOrgByName(members[i].OrgName)
		if e != nil {
			err = fmt.Errorf("failed to get org when get org by user, %w", e)
			return
		}
		orgs[i] = ToDTO(&o)
	}

	return
}

// List retrieves a list of organizations based on the provided options.
func (org *orgService) List(l *OrgListOptions) (orgs []userapp.UserDTO, err error) {
	if l == nil {
		return nil, fmt.Errorf("list options is nil")
	}
	orgs = []userapp.UserDTO{}

	var orgIDs []int64
	if l.Member != nil {
		orgIDs, err = org.getOrgIDsByUserAndRoles(l.Member, l.Roles)
		if err != nil || len(orgIDs) == 0 {
			return
		}
	}

	listOption := &userrepo.ListOrgOption{
		OrgIDs: orgIDs,
		Owner:  l.Owner,
	}
	os, err := org.repo.GetOrgList(listOption)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		}

		return
	}

	orgs = make([]userapp.UserDTO, len(os))
	for i := range os {
		orgs[i] = ToDTO(&os[i])
	}

	return
}

func (org *orgService) getOrgIDsByUserAndRoles(user primitive.Account,
	roles []primitive.Role) (orgIDs []int64, err error) {
	members, err := org.member.GetByUserAndRoles(user, roles)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		}
		return
	}

	for _, mem := range members {
		orgIDs = append(orgIDs, mem.OrgId.Integer())
	}

	return
}

// ListMember retrieves a list of members for a given organization.
func (org *orgService) ListMember(acc primitive.Account) (dtos []MemberDTO, err error) {
	if acc == nil {
		err = fmt.Errorf("account is nil")
		return
	}

	o, err := org.GetByAccount(acc)
	if err != nil {
		return
	}

	members, e := org.member.GetByOrg(o.Name)
	if e != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		}

		return
	}

	dtos = make([]MemberDTO, len(members))
	for i := range members {
		dtos[i] = ToMemberDTO(&members[i])
		dtos[i].OrgName = o.Name
		dtos[i].OrgFullName = o.Fullname
	}

	return
}

// AddMember adds a new member to an organization.
func (org *orgService) AddMember(cmd *domain.OrgAddMemberCmd) error {
	err := cmd.Validate()
	if err != nil {
		return allerror.NewInvalidParam(err.Error())
	}

	o, err := org.repo.GetOrgByName(cmd.Org)
	if err != nil {
		logrus.Error(err)
		return allerror.NewInvalidParam("failed to get org info")

	}

	m := cmd.ToMember()

	pl, err := org.user.GetPlatformUser(o.Owner)
	if err != nil {
		logrus.Error(err)
		return allerror.NewInvalidParam("failed to get platform user for adding member")
	}

	err = pl.AddMember(&o, &m)
	if err != nil {
		logrus.Error(err)
		return allerror.NewInvalidParam(fmt.Sprintf("failed to add member:%s to org:%s",
			m.Username.Account(), o.Account.Account()))
	}

	_, err = org.member.Add(&m)
	if err != nil {
		logrus.Error(err)
		return allerror.NewInvalidParam("failed to save member for adding member")
	}

	return nil
}

func (org *orgService) canRemoveMember(cmd *domain.OrgRemoveMemberCmd) (err error) {
	// check if this is the only owner
	members, err := org.member.GetByOrg(cmd.Org.Account())
	if err != nil {
		err = fmt.Errorf("failed to get members by org name: %s, %s", cmd.Org, err)
		return
	}

	member := cmd.ToMember()

	count := len(members)
	if count == 1 {
		err = fmt.Errorf("the org has only one member")
		return
	}

	if count == 0 {
		err = fmt.Errorf("the org has no member")
		return
	}

	ownerCount := 0
	removeOwner := false
	can := false
	for _, m := range members {
		if m.Role == primitive.Admin {
			ownerCount++
			if m.Username == member.Username {
				removeOwner = true
				can = true
			}
		}
		if m.Username == member.Username {
			can = true
		}
	}

	if ownerCount == 1 && removeOwner {
		err = allerror.NewNoPermission("the only owner can not be removed")
		return
	}

	if !can {
		err = allerror.NewNoPermission("the member is not in the org")
		return
	}

	return
}

// RemoveMember removes a member from an organization.
func (org *orgService) RemoveMember(cmd *domain.OrgRemoveMemberCmd) error {
	err := cmd.Validate()
	if err != nil {
		return allerror.NewInvalidParam(err.Error())
	}

	err = org.canRemoveMember(cmd)
	if err != nil {
		return allerror.NewInvalidParam(err.Error())
	}

	o, err := org.repo.GetOrgByName(cmd.Org)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("org %s not found", cmd.Org.Account()))
		}

		return err
	}

	_, err = org.repo.GetByAccount(cmd.Actor)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeUserNotFound, fmt.Sprintf("user %s not existed",
				cmd.Actor.Account()))
		}

		return err
	}

	if cmd.Actor.Account() != cmd.Account.Account() {
		_, err = org.repo.GetByAccount(cmd.Account)
		if err != nil {
			if commonrepo.IsErrorResourceNotExists(err) {
				err = allerror.NewNotFound(allerror.ErrorCodeUserNotFound, fmt.Sprintf("user %s not existed",
					cmd.Account.Account()))
			}

			return err
		}
		err = org.perm.Check(cmd.Actor, cmd.Org, primitive.ObjTypeMember, primitive.ActionDelete)
		if err != nil {
			return err
		}
	} else {
		err = org.perm.Check(cmd.Actor, cmd.Org, primitive.ObjTypeMember, primitive.ActionRead)
		if err != nil {
			return err
		}
	}

	pl, err := org.user.GetPlatformUser(o.Owner)
	if err != nil {
		return fmt.Errorf("failed to get platform user, %w", err)
	}

	m, err := org.member.GetByOrgAndUser(cmd.Org.Account(), cmd.Account.Account())
	if err != nil {
		return fmt.Errorf("failed to get member when remove member by org %s and user %s, %w", cmd.Org.Account(),
			cmd.Account.Account(), err)
	}

	err = pl.RemoveMember(&o, &m)
	if err != nil {
		return fmt.Errorf("failed to delete git member, %w", err)
	}

	err = org.member.Delete(&m)
	if err != nil {
		_ = pl.AddMember(&o, &m)
		return fmt.Errorf("failed to delete member, %w", err)
	}

	// when owner is removed, a new owner must be set
	if cmd.Account == o.Owner {
		o.Owner = cmd.Actor
		_, err = org.repo.SaveOrg(&o)
		if err != nil {
			return fmt.Errorf("failed to change owner of org, %w", err)
		}
	}

	return nil
}

// EditMember edits the role of a member in an organization.
func (org *orgService) EditMember(cmd *domain.OrgEditMemberCmd) (dto MemberDTO, err error) {
	err = org.perm.Check(cmd.Actor, cmd.Org, primitive.ObjTypeMember, primitive.ActionWrite)
	if err != nil {
		return
	}

	m, err := org.member.GetByOrgAndUser(cmd.Org.Account(), cmd.Account.Account())
	if err != nil {
		err = fmt.Errorf("failed to get member when edit member by org:%s and user:%s, %w",
			cmd.Org.Account(), cmd.Account.Account(), err)
		return
	}

	o, err := org.repo.GetOrgByName(cmd.Org)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("org %s not found", cmd.Org.Account()))
		}

		return
	}

	if o.Owner == cmd.Account {
		err = fmt.Errorf("can't edit owner's role")
		return
	}

	pl, err := org.user.GetPlatformUser(o.Owner)
	if err != nil {
		err = fmt.Errorf("failed to get platform user, %w", err)
		return
	}

	if m.Role != cmd.Role {
		origRole := m.Role
		m.Role = cmd.Role
		err = pl.EditMemberRole(&o, origRole, &m)
		if err != nil {
			err = fmt.Errorf("failed to edit git member, %w", err)
			return
		}

		m, err = org.member.Save(&m)
		if err != nil {
			err = fmt.Errorf("failed to save member, %w", err)
			return
		}
		dto = ToMemberDTO(&m)
	} else {
		logrus.Warn("role not changed")
	}

	return
}

// InviteMember invites a new member to an organization.
func (org *orgService) InviteMember(cmd *domain.OrgInviteMemberCmd) (dto ApproveDTO, err error) {
	if err = cmd.Validate(); err != nil {
		return
	}

	if org.hasMember(cmd.Org, cmd.Account) {
		err = allerror.NewInvalidParam("the user is already a member of the org")
		return
	}

	invitee, err := org.repo.GetByAccount(cmd.Account)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeUserNotFound, "invitee not found")
		}

		return
	}

	inviter, err := org.repo.GetByAccount(cmd.Actor)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeUserNotFound, "inviter not found")
		}

		return
	}

	o, err := org.repo.GetOrgByName(cmd.Org)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeUserNotFound, "organization not found")
		}

		return
	}

	err = org.perm.Check(cmd.Actor, cmd.Org, primitive.ObjTypeMember, primitive.ActionCreate)
	if err != nil {
		return
	}

	invite := cmd.ToApprove(org.inviteExpiry)
	invite.InviterId = inviter.Id
	invite.UserId = invitee.Id
	invite.OrgId = o.Id

	*invite, err = org.invite.AddInvite(invite)
	if err != nil {
		err = fmt.Errorf("failed to save member, %w", err)
		return
	}

	dto = ToApproveDTO(invite, org.user)

	return
}

func (org *orgService) hasMember(o, user primitive.Account) bool {
	_, err := org.member.GetByOrgAndUser(o.Account(), user.Account())
	if err != nil && !commonrepo.IsErrorResourceNotExists(err) {
		logrus.Errorf("failed to get member when check existence by org:%s and user:%s, %s", o.Account(), user.Account(), err)
		return true
	}

	if err == nil {
		logrus.Warnf("the user %s is already a member of the org %s", user.Account(), o.Account())
		return true
	}

	return false
}

// RequestMember sends a membership request to join an organization.
func (org *orgService) RequestMember(cmd *domain.OrgRequestMemberCmd) (dto MemberRequestDTO, err error) {
	if cmd == nil {
		err = allerror.NewInvalidParam("invalid param for request member")
		return
	}

	if err = cmd.Validate(); err != nil {
		return
	}

	if org.hasMember(cmd.Org, cmd.Actor) {
		err = allerror.NewInvalidParam(fmt.Sprintf(" user %s is already a member of the org %s",
			cmd.Actor.Account(), cmd.Org.Account()))
		return
	}

	requester, err := org.repo.GetByAccount(cmd.Actor)
	if err != nil {
		logrus.Error(err)
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeUserNotFound, "requester not found")
		}

		return

	}

	o, err := org.repo.GetOrgByName(cmd.Org)
	if err != nil {
		logrus.Error(err)
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeUserNotFound, "organization not found")
		}

		return
	}

	if !o.AllowRequest {
		err = allerror.NewInvalidParam("org not allow request member")
		return
	}

	request := cmd.ToMemberRequest(o.DefaultRole)
	request.OrgId = o.Id
	request.UserId = requester.Id

	approve, err := org.invite.AddRequest(request)
	if err != nil {
		return
	}

	dto = ToMemberRequestDTO(&approve, org.user)

	return
}

// AcceptInvite accept the invitation the admin sent to me
func (org *orgService) AcceptInvite(cmd *domain.OrgAcceptInviteCmd) (dto ApproveDTO, err error) {
	if cmd == nil {
		err = allerror.NewInvalidParam("invalid param for cancel request member")
		return
	}

	if err = cmd.Validate(); err != nil {
		return
	}

	if org.hasMember(cmd.Org, cmd.Actor) {
		err = allerror.NewInvalidParam("the user is already a member of the org")
		return
	}

	// list all invitations sent to myself in the org
	o, err := org.invite.ListInvitation(&domain.OrgInvitationListCmd{
		Invitee: cmd.Actor,
		OrgNormalCmd: domain.OrgNormalCmd{
			Org: cmd.Org,
		},
		Status: domain.ApproveStatusPending,
	})
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("the %s's invitation to org %s not found",
				cmd.Actor.Account(), cmd.Org.Account()))
		}

		return
	}

	if len(o) == 0 {
		err = fmt.Errorf("no invitation found")
		return
	}

	approve := o[0]

	if cmd.Actor.Account() != approve.Username.Account() {
		err = allerror.NewNoPermission("can't accept other's invitation")
		return
	}

	if approve.ExpireAt < utils.Now() {
		err = fmt.Errorf("the invitation has expired")
		return
	}

	approve.By = cmd.Actor.Account()
	approve.Status = domain.ApproveStatusApproved
	approve.Msg = cmd.Msg

	invite, err := org.invite.SaveInvite(&approve)
	if err != nil {
		return
	}

	err = org.AddMember(&domain.OrgAddMemberCmd{
		Org:    cmd.Org,
		OrgId:  approve.OrgId,
		User:   cmd.Actor,
		UserId: approve.UserId,
		Role:   approve.Role,
		Type:   domain.InviteTypeInvite,
	})

	if err != nil {
		return
	}

	dto = ToApproveDTO(&invite, org.user)

	// Update all requests and invites status pending to approved
	err = org.invite.UpdateAllApproveStatus(approve.Username, approve.OrgName, approve.Status)

	return
}

// ApproveRequest approve the request from the user outside the org
func (org *orgService) ApproveRequest(cmd *domain.OrgApproveRequestMemberCmd) (dto MemberRequestDTO, err error) {
	if cmd == nil {
		err = allerror.NewInvalidParam("invalid param for cancel request member")
		return
	}

	if err = cmd.Validate(); err != nil {
		return
	}

	if cmd.Actor.Account() == cmd.Requester.Account() {
		err = allerror.NewNoPermission("can't approve your own request")
		return
	}

	err = org.perm.Check(cmd.Actor, cmd.Org, primitive.ObjTypeInvite, primitive.ActionWrite)
	if err != nil {
		return
	}

	reqs, err := org.invite.ListRequests(cmd.ToListReqCmd())
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("the %s's member request to org %s not found",
				cmd.Requester.Account(), cmd.Org.Account()))
		}

		return
	}

	if len(reqs) > 1 {
		err = fmt.Errorf("multiple requests found")
		return
	}

	if len(reqs) == 0 {
		err = fmt.Errorf("no request found")
		return
	}

	request := reqs[0]
	request.By = cmd.Actor.Account()
	request.Status = domain.ApproveStatusApproved
	request.Msg = cmd.Msg

	_, err = org.invite.SaveRequest(&request)
	if err != nil {
		return
	}

	err = org.AddMember(&domain.OrgAddMemberCmd{
		Org:    cmd.Org,
		OrgId:  request.OrgId,
		User:   cmd.Requester,
		UserId: request.UserId,
		Type:   domain.InviteTypeRequest,
		Role:   request.Role,
	})
	if err != nil {
		return
	}

	// Update all requests and invites status pending to approved
	err = org.invite.UpdateAllApproveStatus(request.Username, request.OrgName, request.Status)

	return
}

// CancelReqMember cancels a member request in an organization.
func (org *orgService) CancelReqMember(cmd *domain.OrgCancelRequestMemberCmd) (dto MemberRequestDTO, err error) {
	if cmd == nil {
		err = allerror.NewInvalidParam("invalid param for cancel request member")
		return
	}

	if err = cmd.Validate(); err != nil {
		return
	}
	// user can cancel the request by self
	// or admin can reject the request
	if cmd.Actor.Account() != cmd.Requester.Account() {
		err = org.perm.Check(cmd.Actor, cmd.Org, primitive.ObjTypeInvite, primitive.ActionDelete)
		if err != nil {
			return
		}
	}

	o, err := org.invite.ListRequests(cmd.ToListReqCmd())
	if err != nil {
		return
	}

	if len(o) == 0 {
		err = fmt.Errorf("no request found")
		return
	}

	approve := o[0]
	approve.By = cmd.Actor.Account()
	approve.Status = domain.ApproveStatusRejected
	approve.Msg = cmd.Msg

	updatedRequest, err := org.invite.SaveRequest(&approve)
	if err != nil {
		err = fmt.Errorf("failed to remove invite, %w", err)
		return
	}

	dto = ToMemberRequestDTO(&updatedRequest, org.user)

	return
}

// ListMemberReq lists the member requests for an organization.
func (org *orgService) ListMemberReq(cmd *domain.OrgMemberReqListCmd) (dtos []MemberRequestDTO, err error) {
	if cmd == nil {
		err = allerror.NewInvalidParam("invalid param when list member requests")
		return
	}

	if cmd.Actor == nil {
		err = allerror.NewNoPermission("anno can not list requests")
		return
	}

	if err = cmd.Validate(); err != nil {
		return
	}

	// 只有管理员可以查询组织内的申请
	if cmd.Org != nil {
		err = org.perm.Check(cmd.Actor, cmd.Org, primitive.ObjTypeInvite, primitive.ActionRead)
		if err != nil {
			return
		}
	}

	// 不能列出其他人发出的申请
	if cmd.Requester != nil && cmd.Org == nil && cmd.Actor.Account() != cmd.Requester.Account() {
		err = allerror.NewNoPermission("can't list member request of other user")
		return
	}

	reqs, err := org.invite.ListRequests(cmd)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		}

		return
	}

	dtos = make([]MemberRequestDTO, len(reqs))
	for i := range reqs {
		dtos[i] = ToMemberRequestDTO(&reqs[i], org.user)
	}

	return
}

// RevokeInvite revokes an organization invite.
func (org *orgService) RevokeInvite(cmd *domain.OrgRemoveInviteCmd) (dto ApproveDTO, err error) {
	if err = cmd.Validate(); err != nil {
		return
	}
	// user can revoke the invite by self
	// or admin can revoke the invite
	if cmd.Actor.Account() != cmd.Account.Account() {
		err = org.perm.Check(cmd.Actor, cmd.Org, primitive.ObjTypeInvite, primitive.ActionDelete)
		if err != nil {
			return
		}
	}

	o, err := org.invite.ListInvitation(&domain.OrgInvitationListCmd{
		OrgNormalCmd: domain.OrgNormalCmd{
			Org:   cmd.Org,
			Actor: cmd.Actor,
		},
		Invitee: cmd.Account,
		Status:  domain.ApproveStatusPending,
	})

	if err != nil {
		return
	}

	if len(o) == 0 {
		err = fmt.Errorf("no invite found")
		return
	}

	approve := o[0]
	approve.By = cmd.Actor.Account()
	approve.Status = domain.ApproveStatusRejected
	approve.Msg = cmd.Msg

	updatedInvite, err := org.invite.SaveInvite(&approve)
	if err != nil {
		err = fmt.Errorf("failed to remove invite, %w", err)
		return
	}

	dto = ToApproveDTO(&updatedInvite, org.user)

	return
}

// ListInvitationByOrg lists the invitations based on the org.
func (org *orgService) ListInvitationByOrg(actor, orgName primitive.Account,
	status domain.ApproveStatus) (dtos []ApproveDTO, err error) {
	if _, err = org.repo.GetOrgByName(orgName); err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("org %s not found", orgName.Account()))
		}

		return
	}

	// permission check
	// check role when list invitations in a org
	err = org.perm.Check(actor, orgName, primitive.ObjTypeInvite, primitive.ActionRead)
	if err != nil {
		return
	}

	o, err := org.invite.ListInvitation(&domain.OrgInvitationListCmd{
		OrgNormalCmd: domain.OrgNormalCmd{
			Org: orgName,
		},
		Status: status,
	})
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("org %s not found", orgName.Account()))
		}

		return
	}

	dtos = make([]ApproveDTO, len(o))
	for i := range o {
		dtos[i] = ToApproveDTO(&o[i], org.user)
	}

	return
}

// ListInvitationByInviter lists the invitations based on the inviter.
func (org *orgService) ListInvitationByInviter(actor, inviter primitive.Account,
	status domain.ApproveStatus) (dtos []ApproveDTO, err error) {
	// can't list other's sent invitations
	if inviter != actor {
		err = allerror.NewNoPermission("can not list invitation sent by other user")
		return
	}

	o, err := org.invite.ListInvitation(&domain.OrgInvitationListCmd{
		Inviter: inviter,
		Status:  status,
	})
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("org %s not found", inviter))
		}

		return
	}

	dtos = make([]ApproveDTO, len(o))
	for i := range o {
		dtos[i] = ToApproveDTO(&o[i], org.user)
	}

	return
}

// ListInvitationByInvitee lists the invitations based on the invitee.
func (org *orgService) ListInvitationByInvitee(actor, invitee primitive.Account,
	status domain.ApproveStatus) (dtos []ApproveDTO, err error) {
	// permission check
	// can't list other's received invitations
	if invitee != actor {
		err = allerror.NewNoPermission("can not list invitation received by other user")
		return
	}

	o, err := org.invite.ListInvitation(&domain.OrgInvitationListCmd{
		Invitee: invitee,
		Status:  status,
	})
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("org %s not found", invitee))
		}

		return
	}

	dtos = make([]ApproveDTO, len(o))
	for i := range o {
		dtos[i] = ToApproveDTO(&o[i], org.user)
	}

	return
}

// CheckName checks if the given name exists in the repository.
func (org *orgService) CheckName(name primitive.Account) bool {
	if name == nil {
		logrus.Error("name is nil")
		return false
	}

	return org.repo.CheckName(name)
}

// GetMemberByUserAndOrg retrieves the member information for a given user and organization.
func (org *orgService) GetMemberByUserAndOrg(u primitive.Account, o primitive.Account) (member MemberDTO, err error) {
	if u == nil {
		err = fmt.Errorf("user is nil")
		return
	}

	if o == nil {
		err = fmt.Errorf("org is nil")
		return
	}

	m, err := org.member.GetByOrgAndUser(o.Account(), u.Account())
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("org %s with user %s not found", o.Account(), u.Account()))
		}

		return
	}

	member = ToMemberDTO(&m)

	return
}
