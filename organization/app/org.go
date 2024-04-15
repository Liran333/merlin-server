/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package app provides functionality for handling organization-related operations.
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

func errOrgNotFound(msg string, err error) error {
	if msg == "" {
		msg = "not found"
	}

	return allerror.NewNotFound(allerror.ErrorCodeOrganizationNotFound, msg, err)
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
	HasMember(primitive.Account, primitive.Account) bool
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
	ListMember(*domain.OrgListMemberCmd) ([]MemberDTO, error)
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
		e := fmt.Errorf("name %s is already been taken", cmd.Name.Account())
		err = allerror.New(allerror.ErrorNameAlreadyBeenTaken, "", e)
		return
	}

	if err = org.orgCountCheck(cmd.Owner); err != nil {
		return
	}

	owner, err := org.repo.GetByAccount(cmd.Owner)
	if err != nil {
		err = allerror.New(allerror.ErrorFailedGetOwnerInfo, "", err)
		return
	}

	pl, err := org.user.GetPlatformUser(orgTemp.Owner)
	if err != nil {
		err = allerror.New(allerror.ErrorFailGetPlatformUser, "", err)
		return
	}

	err = pl.CreateOrg(orgTemp)
	if err != nil {
		err = allerror.New(allerror.ErrorFailedCreateOrg, "", err)
		return
	}

	orgTemp.DefaultRole = org.defaultRole
	orgTemp.AllowRequest = false
	orgTemp.OwnerId = owner.Id

	*orgTemp, err = org.repo.AddOrg(orgTemp)
	if err != nil {
		err = allerror.New(allerror.ErrorFailedCreateToOrg, "", err)
		_ = pl.DeleteOrg(cmd.Name)
		return
	}

	_, err = org.member.Add(&domain.OrgMember{
		OrgName:  cmd.Name,
		OrgId:    orgTemp.Id,
		Username: cmd.Owner,
		FullName: owner.Fullname,
		UserId:   owner.Id,
		Role:     primitive.NewAdminRole(),
	})
	if err != nil {
		err = allerror.New(allerror.ErrorFailSaveOrgMember, "", err)
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
		return allerror.NewCountExceeded("org count exceed", fmt.Errorf("org count(now:%d max:%d) exceed",
			total, org.MaxCountPerOwner))
	}

	return nil
}

// GetByAccount retrieves an organization by its account and returns it as a UserDTO.
func (org *orgService) GetByAccount(acc primitive.Account) (dto userapp.UserDTO, err error) {
	o, err := org.repo.GetOrgByName(acc)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("org %s not found", acc.Account()), fmt.Errorf("org %s not found, %w",
				acc.Account(), err))
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
			err = allerror.New(allerror.ErrorCodeUserNotFound, fmt.Sprintf("org %s not found", acc.Account()),
				fmt.Errorf("org %s not found, %w", acc.Account(), err))
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
		return allerror.New(allerror.ErrorFailGetPlatformUser, "", fmt.Errorf("failed to get platform user, %w", err))
	}

	can, err := pl.CanDelete(cmd.Name)
	if err != nil {
		return allerror.New(allerror.ErrorAccountCannotDeleteTheOrg, "can't delete the org",
			fmt.Errorf("%s can't delete the org, %w", cmd.Name.Account(), err))
	}

	if !can {
		e := fmt.Errorf("can't delete the organization, while some repos still existed")
		return allerror.New(allerror.ErrorCodeOrgExistResource, e.Error(), e)
	}

	err = org.member.DeleteByOrg(o.Account)
	if err != nil {
		return allerror.New(allerror.ErrorBaseCase, "failed to delete org member",
			fmt.Errorf("failed to delete org member, %w", err))
	}

	err = org.invite.DeleteInviteAndReqByOrg(o.Account)
	if err != nil {
		return allerror.New(allerror.ErrorBaseCase, "failed to delete org invite",
			fmt.Errorf("failed to delete org invite, %w", err))
	}

	err = org.git.DeleteOrg(o.Account)
	if err != nil {
		return allerror.New(allerror.ErrorBaseCase, "failed to delete git org",
			fmt.Errorf("failed to delete git org, %w", err))
	}

	return org.repo.DeleteOrg(&o)
}

// UpdateBasicInfo updates the basic information of an organization based on the provided command
// and returns the updated organization as a UserDTO.
func (org *orgService) UpdateBasicInfo(cmd *domain.OrgUpdatedBasicInfoCmd) (dto userapp.UserDTO, err error) {
	if cmd == nil {
		err = allerror.New(allerror.ErrorSystemError, "", err)
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
			err = errOrgNotFound(fmt.Sprintf("org %s not found", cmd.OrgName.Account()),
				fmt.Errorf("org %s not found, %w", cmd.OrgName.Account(), err))
		}

		return
	}

	change := cmd.ToOrg(&o)
	if change {
		o, err = org.repo.SaveOrg(&o)
		if err != nil {
			err = allerror.New(allerror.ErrorFailedToSaveOrg, "", fmt.Errorf("failed to save org, %w", err))
			return
		}
		dto = ToDTO(&o)
		return
	}

	err = allerror.New(allerror.ErrorNothingChanged, "", fmt.Errorf("nothing changed when update basic info %v", cmd))
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
		e := fmt.Errorf("account is nil")
		err = allerror.New(allerror.ErrorSystemError, "account is nil", e)
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
			e := fmt.Errorf("failed to get org when get org by user, %w", e)
			err = allerror.New(allerror.ErrorFailedToGetOrg, "", e)
			return
		}
		orgs[i] = ToDTO(&o)
	}

	return
}

// List retrieves a list of organizations based on the provided options.
func (org *orgService) List(l *OrgListOptions) (orgs []userapp.UserDTO, err error) {
	if l == nil {
		e := fmt.Errorf("list options is nil")
		return nil, allerror.New(allerror.ErrorSystemError, "", e)
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
func (org *orgService) ListMember(cmd *domain.OrgListMemberCmd) (dtos []MemberDTO, err error) {
	if cmd == nil || cmd.Org == nil {
		e := fmt.Errorf("org account is nil")
		err = allerror.New(allerror.ErrorSystemError, "", e)
		return
	}

	o, err := org.GetByAccount(cmd.Org)
	if err != nil {
		return
	}

	members, e := org.member.GetByOrg(cmd)
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
		e := fmt.Errorf("failed to validate cmd, %w", err)
		return allerror.New(allerror.ErrorFailedToValidateCmd, "", e)
	}

	o, err := org.repo.GetOrgByName(cmd.Org)
	if err != nil {
		return allerror.New(allerror.ErrorFailedToGetOrgInfo, "", fmt.Errorf("failed to get org info, %w", err))
	}

	memberInfo, err := org.repo.GetByAccount(cmd.User)
	if err != nil {
		return allerror.New(allerror.ErrorFailedToGetMemberInfo, "", fmt.Errorf("failed to get member info, %w", err))
	}

	m := cmd.ToMember(memberInfo)

	pl, err := org.user.GetPlatformUser(cmd.Actor)
	if err != nil {
		return allerror.New(allerror.ErrorFailGetPlatformUser,
			"failed to get platform user for adding member", fmt.Errorf("failed to get platform user for adding member, %w", err))
	}

	err = pl.AddMember(&o, &m)
	if err != nil {
		return allerror.New(allerror.ErrorFailedToAddMemberToOrg, "", fmt.Errorf("failed to add member:%s to org:%s, %w",
			m.Username.Account(), o.Account.Account(), err))
	}

	_, err = org.member.Add(&m)
	if err != nil {
		return allerror.New(allerror.ErrorFailedToSaveMemberForAddingMember, "",
			fmt.Errorf("failed to save member for adding member, %w", err))
	}

	return nil
}

func (org *orgService) canEditMember(cmd *domain.OrgEditMemberCmd) (err error) {
	return org.canRemoveMember(&domain.OrgRemoveMemberCmd{
		Org:     cmd.Org,
		Account: cmd.Account,
		Actor:   cmd.Actor,
		Msg:     "",
	})
}

func (org *orgService) members(orgName primitive.Account) ([]domain.OrgMember, int, error) {
	members, err := org.member.GetByOrg(&domain.OrgListMemberCmd{Org: orgName})
	if err != nil {
		e := fmt.Errorf("failed to get members by org name: %s, %s", orgName, err)
		err = allerror.New(allerror.ErrorFailedToGetMembersByOrgName, "", e)
		return []domain.OrgMember{}, 0, err
	}

	return members, len(members), nil

}

func (org *orgService) getOwners(orgName primitive.Account) ([]domain.OrgMember, error) {
	members, err := org.member.GetByOrg(&domain.OrgListMemberCmd{Org: orgName, Role: primitive.Admin})
	if err != nil {
		e := fmt.Errorf("failed to get members by org name: %s, %s", orgName, err)
		err = allerror.NewInvalidParam(e.Error(), e)
		return []domain.OrgMember{}, err
	}

	if len(members) == 0 {
		e := fmt.Errorf("no owners found in org %s", orgName.Account())
		err = allerror.NewInvalidParam(e.Error(), e)
		return []domain.OrgMember{}, err
	}

	return members, nil
}

func (org *orgService) canRemoveMember(cmd *domain.OrgRemoveMemberCmd) (err error) {
	// check if this is the only owner
	members, count, err := org.members(cmd.Org)
	if err != nil {
		return err
	}
	if count == 1 {
		e := fmt.Errorf("the org has only one member")
		err = allerror.New(allerror.ErrorOrgHasOnlyOneMember, "the org has only one member", e)
		return
	}

	if count == 0 {
		e := fmt.Errorf("the org has no member")
		err = allerror.NewNoPermission(e.Error(), e)
		return
	}

	member := cmd.ToMember()

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
		e := fmt.Errorf("the only owner can not be removed")
		err = allerror.New(allerror.ErrorOnlyOwnerCanNotBeRemoved, "", e)
		return
	}

	if !can {
		e := fmt.Errorf("the member is not in the org")
		err = allerror.NewNoPermission(e.Error(), e)
		return
	}

	return
}

// RemoveMember removes a member from an organization.
func (org *orgService) RemoveMember(cmd *domain.OrgRemoveMemberCmd) error {
	err := cmd.Validate()
	if err != nil {
		return allerror.New(allerror.ErrorFailedToValidateCmd, "", fmt.Errorf("failed to validate cmd, %w", err))
	}

	err = org.canRemoveMember(cmd)
	if err != nil {
		return allerror.New(allerror.ErrorFailedToRemoveMember, "", fmt.Errorf("failed to validate cmd, %w", err))
	}

	o, err := org.repo.GetOrgByName(cmd.Org)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("org %s not found", cmd.Org.Account()), err)
		}

		return err
	}

	_, err = org.repo.GetByAccount(cmd.Actor)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeUserNotFound, fmt.Sprintf("user %s not existed",
				cmd.Actor.Account()), fmt.Errorf("user %s not existed: %w", cmd.Actor.Account(), err))
		}

		return err
	}

	if cmd.Actor.Account() != cmd.Account.Account() {
		_, err = org.repo.GetByAccount(cmd.Account)
		if err != nil {
			if commonrepo.IsErrorResourceNotExists(err) {
				err = allerror.NewNotFound(allerror.ErrorCodeUserNotFound, fmt.Sprintf("user %s not existed",
					cmd.Account.Account()), fmt.Errorf("user %s not existed: %w", cmd.Actor.Account(), err))
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

	owners, err := org.getOwners(cmd.Org)
	if err != nil {
		e := fmt.Errorf("failed to get owners of org when add new member: %s, %s", cmd.Org.Account(), err)
		return allerror.NewInvalidParam(e.Error(), e)
	}

	pl, err := org.user.GetPlatformUser(owners[0].Username)
	if err != nil {
		return allerror.New(allerror.ErrorFailGetPlatformUser, "", err)
	}

	m, err := org.member.GetByOrgAndUser(cmd.Org.Account(), cmd.Account.Account())
	if err != nil {
		e := fmt.Errorf("failed to get member when remove member by org %s and user %s, %w",
			cmd.Org.Account(), cmd.Account.Account(), err)
		return allerror.New(allerror.ErrorFailedToRemoveMember, "", e)
	}

	err = pl.RemoveMember(&o, &m)
	if err != nil {
		return allerror.New(allerror.ErrorFailedToDeleteGitMember, "", fmt.Errorf("failed to delete git member, %w", err))
	}

	err = org.member.Delete(&m)
	if err != nil {
		_ = pl.AddMember(&o, &m)
		return allerror.New(allerror.ErrorFailedToDeleteMember, "",
			fmt.Errorf("failed to delete member, %w", err))
	}

	// when owner is removed, a new owner must be set
	if cmd.Account == o.Owner {
		o.Owner = cmd.Actor
		_, err = org.repo.SaveOrg(&o)
		if err != nil {
			return allerror.New(allerror.ErrorFailedToChangeOwnerOfOrg, "",
				fmt.Errorf("failed to change owner of org, %w", err))
		}
	}

	return nil
}

// EditMember edits the role of a member in an organization.
func (org *orgService) EditMember(cmd *domain.OrgEditMemberCmd) (dto MemberDTO, err error) {
	if err = org.canEditMember(cmd); err != nil {
		return
	}

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
			err = errOrgNotFound(fmt.Sprintf("org %s not found", cmd.Org.Account()), err)
		}

		return
	}

	pl, err := org.user.GetPlatformUser(cmd.Actor)
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

	if org.HasMember(cmd.Org, cmd.Account) {
		e := fmt.Errorf("the user is already a member of the org")
		err = allerror.New(allerror.ErrorUserAlreadyInOrg, "", e)
		return
	}

	invitee, err := org.repo.GetByAccount(cmd.Account)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeUserNotFound, "invitee not found", err)
		}

		return
	}

	inviter, err := org.repo.GetByAccount(cmd.Actor)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeUserNotFound, "inviter not found", err)
		}

		return
	}

	o, err := org.repo.GetOrgByName(cmd.Org)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeUserNotFound, "organization not found", err)
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

// HasMember returns true if the user is already a member of the organization.
func (org *orgService) HasMember(o, user primitive.Account) bool {
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
		e := fmt.Errorf("invalid param for request member")
		err = allerror.New(allerror.ErrorInvalidParamForRequestMember, "", e)
		return
	}

	if err = cmd.Validate(); err != nil {
		return
	}

	if org.HasMember(cmd.Org, cmd.Actor) {
		e := fmt.Errorf(" user %s is already a member of the org %s", cmd.Actor.Account(), cmd.Org.Account())
		err = allerror.New(allerror.ErrorUserAccountIsAlreadyAMemberOfOrgAccount, "", e)
		return
	}

	requester, err := org.repo.GetByAccount(cmd.Actor)
	if err != nil {
		logrus.Error(err)
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeUserNotFound, "requester not found", err)
		}

		return

	}

	o, err := org.repo.GetOrgByName(cmd.Org)
	if err != nil {
		logrus.Error(err)
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeUserNotFound, "organization not found", err)
		}

		return
	}

	if !o.AllowRequest {
		err = allerror.New(allerror.ErrorOrgNotAllowRequestMember, "", fmt.Errorf("org not allow request member"))
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
		e := fmt.Errorf("invalid param for cancel request member")
		err = allerror.New(allerror.ErrorInvalidParamForCancelRequestMember, "", e)
		return
	}

	if err = cmd.Validate(); err != nil {
		return
	}

	if org.HasMember(cmd.Org, cmd.Actor) {
		e := fmt.Errorf("the user %s is already a member of the org %s", cmd.Actor.Account(), cmd.Org.Account())
		err = allerror.New(allerror.ErrorUserAccountIsAlreadyAMemberOfOrgAccount, "", e)
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
				cmd.Actor.Account(), cmd.Org.Account()), err)
		}

		return
	}

	if len(o) == 0 {
		err = fmt.Errorf("no invitation found")
		return
	}

	approve := o[0]

	if cmd.Actor.Account() != approve.Username.Account() {
		e := fmt.Errorf("can't accept other's invitation")
		err = allerror.NewNoPermission(e.Error(), e)
		return
	}

	if approve.ExpireAt < utils.Now() {
		e := fmt.Errorf("the invitation has expired")
		err = allerror.NewExpired(e.Error(), e)
		return
	}

	owners, err := org.getOwners(cmd.Org)
	if err != nil {
		e := fmt.Errorf("failed to get owners of org when add new member: %s, %s", cmd.Org.Account(), err)
		err = allerror.NewInvalidParam(e.Error(), e)
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
		Actor:  owners[0].Username,
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
		e := fmt.Errorf("invalid param for cancel request member")
		err = allerror.New(allerror.ErrorInvalidParamForCancelRequestMember, "", e)
		return
	}

	if err = cmd.Validate(); err != nil {
		return
	}

	if cmd.Actor.Account() == cmd.Requester.Account() {
		e := fmt.Errorf("can't approve request from yourself")
		err = allerror.NewNoPermission(e.Error(), e)
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
				cmd.Requester.Account(), cmd.Org.Account()), err)
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
		Actor:  cmd.Actor,
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
		e := fmt.Errorf("invalid param for cancel request member")
		err = allerror.New(allerror.ErrorInvalidParamForCancelRequestMember, "", e)
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
		e := fmt.Errorf("invalid param for list member request")
		err = allerror.New(allerror.ErrorInvalidParamForListMemberRequest, "", e)
		return
	}

	if cmd.Actor == nil {
		e := fmt.Errorf("anno can not list requests")
		err = allerror.NewNoPermission(e.Error(), e)
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
		e := fmt.Errorf("can't list requests from other people")
		err = allerror.NewNoPermission(e.Error(), e)
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
			err = errOrgNotFound(fmt.Sprintf("org %s not found", orgName.Account()), err)
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
			err = errOrgNotFound(fmt.Sprintf("org %s not found", orgName.Account()), err)
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
		e := fmt.Errorf("can not list invitation sent by other user")
		err = allerror.NewNoPermission(e.Error(), e)
		return
	}

	o, err := org.invite.ListInvitation(&domain.OrgInvitationListCmd{
		Inviter: inviter,
		Status:  status,
	})
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("org %s not found", inviter), err)
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
		e := fmt.Errorf("can not list invitation received by other user")
		err = allerror.NewNoPermission(e.Error(), e)
		return
	}

	o, err := org.invite.ListInvitation(&domain.OrgInvitationListCmd{
		Invitee: invitee,
		Status:  status,
	})
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("org %s not found", invitee), err)
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
			err = errOrgNotFound(fmt.Sprintf("org %s with user %s not found", o.Account(), u.Account()), err)
		}

		return
	}

	member = ToMemberDTO(&m)

	return
}
