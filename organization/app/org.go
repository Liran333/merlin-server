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
	"github.com/openmerlin/merlin-server/utils"
)

func errOrgNotFound(msg string) error {
	if msg == "" {
		msg = "not found"
	}

	return allerror.NewNotFound(allerror.ErrorCodeOrganizationNotFound, msg)
}

type OrgService interface {
	// user
	Create(*domain.OrgCreatedCmd) (OrganizationDTO, error)
	Delete(*domain.OrgDeletedCmd) error
	UpdateBasicInfo(*domain.OrgUpdatedBasicInfoCmd) (OrganizationDTO, error)
	GetByAccount(primitive.Account) (OrganizationDTO, error)
	CheckName(primitive.Account) bool
	GetByOwner(primitive.Account, primitive.Account) ([]OrganizationDTO, error)
	GetByUser(primitive.Account, primitive.Account) ([]OrganizationDTO, error)
	InviteMember(*domain.OrgInviteMemberCmd) (ApproveDTO, error)
	RequestMember(*domain.OrgRequestMemberCmd) (MemberRequestDTO, error)
	CancelReqMember(*domain.OrgCancelRequestMemberCmd) (MemberRequestDTO, error)
	ApproveRequest(*domain.OrgApproveRequestMemberCmd) (MemberRequestDTO, error)
	AcceptInvite(*domain.OrgAcceptInviteCmd) (ApproveDTO, error)
	RevokeInvite(*domain.OrgRemoveInviteCmd) (ApproveDTO, error)
	ListMemberReq(*domain.OrgMemberReqListCmd) ([]MemberRequestDTO, error)
	ListInvitation(*domain.OrgInvitationListCmd) ([]ApproveDTO, error)
	AddMember(*domain.OrgAddMemberCmd) error
	RemoveMember(*domain.OrgRemoveMemberCmd) error
	EditMember(*domain.OrgEditMemberCmd) (MemberDTO, error)
	ListMember(primitive.Account) ([]MemberDTO, error)
	GetMemberByUserAndOrg(primitive.Account, primitive.Account) (MemberDTO, error)
}

// ps: platform user service
func NewOrgService(
	user userapp.UserService,
	repo repository.Organization,
	member repository.OrgMember,
	invite repository.Approve,
	request repository.MemberRequest,
	perm Permission,
	cfg *domain.Config,
) OrgService {
	return &orgService{
		user:         user,
		repo:         repo,
		member:       member,
		perm:         perm,
		invite:       invite,
		request:      request,
		defaultRole:  domain.OrgRole(cfg.DefaultRole),
		inviteExpiry: cfg.InviteExpiry,
	}
}

type orgService struct {
	inviteExpiry int64
	defaultRole  domain.OrgRole
	user         userapp.UserService
	repo         repository.Organization
	member       repository.OrgMember
	invite       repository.Approve
	request      repository.MemberRequest
	perm         Permission
}

func (org *orgService) Create(cmd *domain.OrgCreatedCmd) (o OrganizationDTO, err error) {
	if err = cmd.Validate(); err != nil {
		return
	}

	orgTemp := cmd.ToOrg()

	pl, err := org.user.GetPlatformUser(orgTemp.Owner)
	if err != nil {
		err = fmt.Errorf("failed to get platform user, %w", err)
		return
	}

	err = pl.CreateOrg(orgTemp)
	if err != nil {
		err = fmt.Errorf("failed to create org, %w", err)
		return
	}

	orgTemp.CreatedAt = utils.Now()
	orgTemp.DefaultRole = org.defaultRole
	orgTemp.AllowRequest = false

	_, err = org.repo.Save(orgTemp)
	if err != nil {
		err = fmt.Errorf("failed to save org, %w", err)
		_ = pl.DeleteOrg(cmd.Name)
		return
	}

	_, err = org.member.Save(&domain.OrgMember{
		OrgName:  cmd.Name,
		Username: cmd.Owner,
		Role:     domain.OrgRoleAdmin,
	})
	if err != nil {
		err = fmt.Errorf("failed to save org member, %w", err)
		_ = pl.DeleteOrg(cmd.Name)
		return
	}

	o = ToDTO(orgTemp, org.defaultRole)

	return
}

func (org *orgService) GetByAccount(acc primitive.Account) (dto OrganizationDTO, err error) {
	o, err := org.repo.GetByName(acc)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("org %s not found", acc.Account()))
		}
		return
	}

	dto = ToDTO(&o, org.defaultRole)
	return
}

func (org *orgService) Delete(cmd *domain.OrgDeletedCmd) error {
	err := org.perm.Check(cmd.Actor, cmd.Name, primitive.ObjTypeOrg, primitive.ActionDelete)
	if err != nil {
		return err
	}
	o, err := org.repo.GetByName(cmd.Name)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		}
		return err
	}

	pl, err := org.user.GetPlatformUser(o.Owner)
	if err != nil {
		return fmt.Errorf("failed to get platform user, %w", err)
	}

	can, err := pl.CanDelete(cmd.Name.Account())
	if err != nil {
		return fmt.Errorf("failed to check platform user, %w", err)
	}

	if !can {
		return fmt.Errorf("can't delte the organization, while some repos still existed")
	}

	err = org.repo.Delete(&o)
	if err != nil {
		return fmt.Errorf("failed to delete org in repo, %w", err)
	}

	err = org.member.DeleteByOrg(o.Name.Account())
	if err != nil {
		return fmt.Errorf("failed to delete org member, %w", err)
	}

	err = pl.DeleteOrg(cmd.Name.Account())
	if err != nil {
		err = fmt.Errorf("failed to delete git org, %w", err)
	}

	return err
}

func (org *orgService) UpdateBasicInfo(cmd *domain.OrgUpdatedBasicInfoCmd) (dto OrganizationDTO, err error) {
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

	o, err := org.repo.GetByName(cmd.OrgName)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("org %s not found", cmd.OrgName.Account()))
		}

		return
	}

	change, err := cmd.ToOrg(&o)
	if err != nil {
		return
	}

	if change {
		o, err = org.repo.Save(&o)
		if err != nil {
			err = fmt.Errorf("failed to save org, %w", err)
			return
		}
		dto = ToDTO(&o, org.defaultRole)
		return
	}
	err = allerror.NewInvalidParam("nothing changed")
	return
}

func (org *orgService) GetByOwner(actor, acc primitive.Account) (orgs []OrganizationDTO, err error) {
	if acc == nil {
		err = fmt.Errorf("account is nil")
		return
	}

	if acc.Account() != actor.Account() {
		err = fmt.Errorf("can't list organizations for other users")
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

func (org *orgService) GetByUser(actor, acc primitive.Account) (orgs []OrganizationDTO, err error) {
	if acc == nil {
		err = fmt.Errorf("account is nil")
		return
	}

	if acc.Account() != actor.Account() {
		err = fmt.Errorf("can't list organizations for other users")
		return
	}

	members, err := org.member.GetByUser(acc.Account())
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		}
		return
	}

	orgs = make([]OrganizationDTO, len(members))
	for i := range members {
		o, e := org.repo.GetByName(primitive.CreateAccount(members[i].OrgName))
		if e != nil {
			err = fmt.Errorf("failed to get org when get org by user, %w", e)
			return
		}
		orgs[i] = ToDTO(&o, org.defaultRole)
	}

	return
}

func (org *orgService) List(l *OrgListOptions) (orgs []OrganizationDTO, err error) {
	if l == nil {
		return nil, fmt.Errorf("list options is nil")
	}

	os, err := org.repo.GetByOwner(l.Owner)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		}
		return
	}

	orgs = make([]OrganizationDTO, len(os))
	for i := range os {
		orgs[i] = ToDTO(&os[i], org.defaultRole)
	}

	return
}

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
		dtos[i].OrgFullName = o.FullName
	}

	return
}

func (org *orgService) AddMember(cmd *domain.OrgAddMemberCmd) error {
	err := cmd.Validate()
	if err != nil {
		return allerror.NewInvalidParam(err.Error())
	}

	o, err := org.repo.GetByName(cmd.Org)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("org %s not found", cmd.Org))
		}
		return err
	}

	m := cmd.ToMember()

	pl, err := org.user.GetPlatformUser(o.Owner)
	if err != nil {
		return fmt.Errorf("failed to get platform user for adding member, %w", err)
	}

	err = pl.AddMember(&o, &m)
	if err != nil {
		return fmt.Errorf("failed to add git member:%s to org:%s, %w", m.Username, o.Name.Account(), err)
	}

	_, err = org.member.Save(&m)
	if err != nil {
		// TODO need rollback
		return fmt.Errorf("failed to save member for adding member, %w", err)
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
		if m.Role == domain.OrgRoleAdmin {
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

func (org *orgService) RemoveMember(cmd *domain.OrgRemoveMemberCmd) error {
	err := cmd.Validate()
	if err != nil {
		return allerror.NewInvalidParam(err.Error())
	}

	err = org.canRemoveMember(cmd)
	if err != nil {
		return err
	}

	if cmd.Actor.Account() != cmd.Account.Account() {
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
	o, err := org.repo.GetByName(cmd.Org)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		}
		return err
	}

	pl, err := org.user.GetPlatformUser(o.Owner)
	if err != nil {
		return fmt.Errorf("failed to get platform user, %w", err)
	}

	m, err := org.member.GetByOrgAndUser(cmd.Org.Account(), cmd.Account.Account())
	if err != nil {
		return fmt.Errorf("failed to get member by org %s and user %s, %w", cmd.Org.Account(), cmd.Account.Account(), err)
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

	return nil
}

func (org *orgService) EditMember(cmd *domain.OrgEditMemberCmd) (dto MemberDTO, err error) {
	err = org.perm.Check(cmd.Actor, cmd.Org, primitive.ObjTypeMember, primitive.ActionWrite)
	if err != nil {
		return
	}

	m, err := org.member.GetByOrgAndUser(cmd.Org.Account(), cmd.Account.Account())
	if err != nil {
		err = fmt.Errorf("failed to get member by org:%s and user:%s, %w", cmd.Org.Account(), cmd.Account.Account(), err)
		return
	}

	o, err := org.repo.GetByName(cmd.Org)
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

	if m.Role != domain.OrgRole(cmd.Role) {
		origRole := m.Role
		m.Role = domain.OrgRole(cmd.Role)
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

func (org *orgService) InviteMember(cmd *domain.OrgInviteMemberCmd) (dto ApproveDTO, err error) {
	if err = cmd.Validate(); err != nil {
		return
	}

	err = org.perm.Check(cmd.Actor, cmd.Org, primitive.ObjTypeMember, primitive.ActionCreate)
	if err != nil {
		return
	}

	if org.hasMember(cmd.Org, cmd.Account) {
		err = allerror.NewInvalidParam("the user is already a member of the org")
		return
	}

	_, err = org.user.GetByAccount(cmd.Account, false)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("user %s not found", cmd.Account.Account()))
		}
		return
	}

	invite := cmd.ToApprove(org.inviteExpiry)

	newOrg, err := org.invite.Save(&invite)
	if err != nil {
		err = fmt.Errorf("failed to save member, %w", err)
		return
	}

	dto = ToApproveDTO(&newOrg, org.user)

	return
}

func (org *orgService) hasMember(o, user primitive.Account) bool {
	_, err := org.member.GetByOrgAndUser(o.Account(), user.Account())
	if err != nil && !commonrepo.IsErrorResourceNotExists(err) {
		logrus.Errorf("failed to get member by org:%s and user:%s, %s", o.Account(), user.Account(), err)
		return true
	}

	if err == nil {
		logrus.Warnf("the user %s is already a member of the org %s", user.Account(), o.Account())
		return true
	}

	return false
}

func (org *orgService) RequestMember(cmd *domain.OrgRequestMemberCmd) (dto MemberRequestDTO, err error) {
	if cmd == nil {
		err = allerror.NewInvalidParam("invalid param for request member")
		return
	}

	if err = cmd.Validate(); err != nil {
		return
	}

	if org.hasMember(cmd.Org, cmd.Actor) {
		err = allerror.NewInvalidParam(fmt.Sprintf(" user %s is already a member of the org %s", cmd.Actor.Account(), cmd.Org.Account()))
		return
	}

	if !org.user.HasUser(cmd.Actor) {
		err = allerror.NewInvalidParam(fmt.Sprintf("the user %s not found", cmd.Actor.Account()))
		return

	}

	o, err := org.repo.GetByName(cmd.Org)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("org %s not found", cmd.Org.Account()))
		}
		return
	}

	if !o.AllowRequest {
		err = allerror.NewInvalidParam("org not allow request member")
		return
	}

	if o.DefaultRole == "" {
		o.DefaultRole = org.defaultRole
		_, _ = org.repo.Save(&o)
	}

	d, err := org.request.Save(&domain.MemberRequest{
		OrgName:   cmd.Org.Account(),
		Username:  cmd.Actor.Account(),
		Role:      o.DefaultRole,
		Status:    domain.ApproveStatusPending,
		CreatedAt: utils.Now(),
		Msg:       cmd.Msg,
	})

	if err != nil {
		return
	}

	dto = ToMemberRequestDTO(&d, org.user)

	return
}

// I accept the invitation the admin sent to me
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
			err = errOrgNotFound(fmt.Sprintf("the %s's invitation to org %s not found", cmd.Actor.Account(), cmd.Org.Account()))
		}
		return
	}

	if len(o) > 1 {
		err = fmt.Errorf("multiple invitations found")
		return
	}

	if len(o) == 0 {
		err = fmt.Errorf("no invitation found")
		return
	}

	approve := o[0]

	if cmd.Actor.Account() != approve.Username {
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
	approve.UpdatedAt = utils.Now()

	_, err = org.invite.Save(&approve)
	if err != nil {
		return
	}

	err = org.AddMember(&domain.OrgAddMemberCmd{
		Org:  cmd.Org,
		User: cmd.Actor,
		Role: approve.Role,
		Type: domain.InviteTypeInvite,
	})

	return
}

// admin approve the request from the user outside the org
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

	reqs, err := org.request.ListInvitation(cmd.ToListReqCmd())
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("the %s's member request to org %s not found", cmd.Requester.Account(), cmd.Org.Account()))
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
	request.UpdatedAt = utils.Now()

	_, err = org.request.Save(&request)
	if err != nil {
		return
	}

	err = org.AddMember(&domain.OrgAddMemberCmd{
		Org:  cmd.Org,
		User: cmd.Requester,
		Type: domain.InviteTypeRequest,
		Role: request.Role,
	})

	return
}

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

	o, err := org.request.ListInvitation(cmd.ToListReqCmd())
	if err != nil {
		return
	}

	if len(o) > 1 {
		err = fmt.Errorf("multiple invitations found")
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
	approve.UpdatedAt = utils.Now()

	new, err := org.request.Save(&approve)
	if err != nil {
		err = fmt.Errorf("failed to remove invite, %w", err)
		return
	}

	dto = ToMemberRequestDTO(&new, org.user)

	return
}

func (org *orgService) ListMemberReq(cmd *domain.OrgMemberReqListCmd) (dtos []MemberRequestDTO, err error) {
	if cmd == nil {
		err = allerror.NewInvalidParam("invalid param when list member requests")
		return
	}

	if err = cmd.Validate(); err != nil {
		return
	}

	if cmd.Actor != nil && cmd.Org != nil {
		err = org.perm.Check(cmd.Actor, cmd.Org, primitive.ObjTypeInvite, primitive.ActionRead)
		if err != nil {
			return
		}
	}

	if cmd.Requester != nil && cmd.Actor != nil && cmd.Org == nil && cmd.Actor.Account() != cmd.Requester.Account() {
		err = allerror.NewNoPermission("can't list member request of others")
		return
	}

	reqs, err := org.request.ListInvitation(cmd)
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

	if len(o) > 1 {
		err = fmt.Errorf("multiple invitations found")
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
	approve.UpdatedAt = utils.Now()

	new, err := org.invite.Save(&approve)
	if err != nil {
		err = fmt.Errorf("failed to remove invite, %w", err)
		return
	}

	dto = ToApproveDTO(&new, org.user)

	return
}

func (org *orgService) ListInvitation(cmd *domain.OrgInvitationListCmd) (dtos []ApproveDTO, err error) {
	if cmd == nil {
		err = allerror.NewInvalidParam("account is nil")
		return
	}

	if err = cmd.Validate(); err != nil {
		return
	}
	// permission check
	if cmd.Org != nil && cmd.Actor != nil {
		err = org.perm.Check(cmd.Actor, cmd.Org, primitive.ObjTypeInvite, primitive.ActionRead)
		if err != nil {
			return
		}
	}

	if cmd.Actor != nil && cmd.Invitee != nil && cmd.Org == nil && cmd.Actor.Account() != cmd.Invitee.Account() {
		err = allerror.NewNoPermission("can not list invitation for other user")
		return
	}

	if cmd.Invitee != nil {
		if cmd.Invitee.Account() != cmd.Actor.Account() {
			err = allerror.NewNoPermission("can not list invitation by invitee for other user")
		}
	}

	if cmd.Inviter != nil {
		if cmd.Inviter.Account() != cmd.Actor.Account() {
			err = allerror.NewNoPermission("can not list invitation by inviter for other user")
		}
	}

	o, err := org.invite.ListInvitation(cmd)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = errOrgNotFound(fmt.Sprintf("org %s not found", cmd.Org.Account()))
		}
		return
	}

	dtos = make([]ApproveDTO, len(o))
	for i := range o {
		dtos[i] = ToApproveDTO(&o[i], org.user)
	}

	return
}

func (org *orgService) CheckName(name primitive.Account) bool {
	if name == nil {
		logrus.Error("name is nil")
		return false
	}

	return org.repo.CheckName(name)
}

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
