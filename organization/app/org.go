package app

import (
	"fmt"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/infrastructure/repositories"
	"github.com/openmerlin/merlin-server/organization/domain"
	"github.com/openmerlin/merlin-server/organization/domain/repository"
	userapp "github.com/openmerlin/merlin-server/user/app"
	user "github.com/openmerlin/merlin-server/user/domain"
	"github.com/openmerlin/merlin-server/utils"
	"github.com/sirupsen/logrus"
)

type OrgService interface {
	// user
	Create(*domain.OrgCreatedCmd) (OrganizationDTO, error)
	Delete(primitive.Account) error
	UpdateBasicInfo(primitive.Account, *UpdateOrgBasicInfoCmd) (OrganizationDTO, error)
	GetByAccount(primitive.Account) (OrganizationDTO, error)
	CheckName(primitive.Account) bool
	GetByOwner(primitive.Account) ([]OrganizationDTO, error)
	GetByUser(primitive.Account) ([]OrganizationDTO, error)
	InviteMember(*OrgInviteMemberCmd) (OrganizationDTO, error)
	RevokeInvite(*OrgRemoveInviteCmd) (OrganizationDTO, error)
	ListInvitation(primitive.Account) ([]ApproveDTO, error)
	AddMember(*OrgAddMemberCmd) error
	RemoveMember(*OrgRemoveMemberCmd) error
	EditMember(*OrgEditMemberCmd) (MemberDTO, error)
	ListMember(primitive.Account) ([]MemberDTO, error)
}

// ps: platform user service
func NewOrgService(
	user userapp.UserService,
	repo repository.Organization,
	member repository.OrgMember,
	expiry int64,
) OrgService {
	return &orgService{
		inviteExpiry: expiry,
		user:         user,
		repo:         repo,
		member:       member,
	}
}

type orgService struct {
	inviteExpiry int64
	user         userapp.UserService
	repo         repository.Organization
	member       repository.OrgMember
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

	_, err = org.repo.Save(orgTemp)
	if err != nil {
		err = fmt.Errorf("failed to save org, %w", err)
		_ = pl.DeleteOrg(cmd.Name)
		return
	}

	_, err = org.member.Save(&domain.OrgMember{
		OrgName:  cmd.Name,
		Username: cmd.Owner,
		Role:     domain.OrgRoleOwner,
	})
	if err != nil {
		err = fmt.Errorf("failed to save org member, %w", err)
		_ = pl.DeleteOrg(cmd.Name)
		return
	}

	o = ToDTO(orgTemp)

	return
}

func (org *orgService) GetByAccount(acc primitive.Account) (dto OrganizationDTO, err error) {
	o, err := org.repo.GetByName(acc)
	if err != nil {
		return
	}

	dto = ToDTO(&o)
	return
}
func (org *orgService) Delete(acc primitive.Account) error {
	o, err := org.repo.GetByName(acc)
	if err != nil {
		return fmt.Errorf("failed to get org, %w", err)
	}

	pl, err := org.user.GetPlatformUser(o.Owner)
	if err != nil {
		return fmt.Errorf("failed to get platform user, %w", err)
	}

	can, err := pl.CanDelete(acc.Account())
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

	err = pl.DeleteOrg(acc.Account())
	if err != nil {
		err = fmt.Errorf("failed to delete git org, %w", err)
	}

	return err
}

func (org *orgService) UpdateBasicInfo(acc primitive.Account, cmd *UpdateOrgBasicInfoCmd) (dto OrganizationDTO, err error) {
	if acc == nil {
		err = fmt.Errorf("account is nil")
		return
	}

	if cmd == nil {
		err = fmt.Errorf("cmd is nil")
		return
	}

	o, err := org.repo.GetByName(acc)
	if err != nil {
		err = fmt.Errorf("failed to get org, %w", err)
		return
	}

	change := false
	if cmd.AvatarId != o.AvatarId.AvatarId() {
		o.AvatarId, err = user.NewAvatarId(cmd.AvatarId)
		if err != nil {
			err = fmt.Errorf("failed to create avatar id, %w", err)
			return
		}
		change = true
	}

	if cmd.Website != o.Website && cmd.Website != "" {
		o.Website = cmd.Website
		change = true
	}

	if cmd.Description != o.Description && cmd.Description != "" {
		o.Description = cmd.Description
		change = true
	}

	if cmd.FullName != o.FullName && cmd.FullName != "" {
		if !org.CheckName(primitive.CreateAccount(cmd.FullName)) {
			err = fmt.Errorf("%s is not available", cmd.FullName)
			return
		}
		o.FullName = cmd.FullName
		change = true
	}

	if change {
		o, err = org.repo.Save(&o)
		if err != nil {
			err = fmt.Errorf("failed to save org, %w", err)
			return
		}
		dto = ToDTO(&o)
		return
	}
	err = fmt.Errorf("nothing changed")
	return
}

func (org *orgService) GetByOwner(acc primitive.Account) (orgs []OrganizationDTO, err error) {
	if acc == nil {
		err = fmt.Errorf("account is nil")
		return
	}

	return org.List(&OrgListOptions{
		Owner: acc,
	})
}

func (org *orgService) GetByUser(acc primitive.Account) (orgs []OrganizationDTO, err error) {
	if acc == nil {
		err = fmt.Errorf("account is nil")
		return
	}

	members, err := org.member.GetByUser(acc.Account())
	if err != nil {
		err = fmt.Errorf("failed to get members, %w", err)
		return
	}

	orgs = make([]OrganizationDTO, len(members))
	for i := range members {
		o, e := org.repo.GetByName(primitive.CreateAccount(members[i].OrgName))
		if e != nil {
			err = fmt.Errorf("failed to get org when get org by user, %w", e)
			return
		}
		orgs[i] = ToDTO(&o)
	}

	return
}

func (org *orgService) List(l *OrgListOptions) (orgs []OrganizationDTO, err error) {
	if l == nil {
		return nil, fmt.Errorf("list options is nil")
	}

	os, err := org.repo.GetByOwner(l.Owner)
	if err != nil {
		return
	}

	orgs = make([]OrganizationDTO, len(os))
	for i := range os {
		orgs[i] = ToDTO(&os[i])
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
		err = fmt.Errorf("failed to get org, %w", err)
		return
	}

	members, e := org.member.GetByOrg(o.Name)
	if e != nil {
		err = fmt.Errorf("failed to get members by org name: %s, %s", o.Name, e)
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

func (org *orgService) AddMember(cmd *OrgAddMemberCmd) error {
	err := cmd.Validate()
	if err != nil {
		return err
	}

	o, err := org.repo.GetByName(cmd.Org)
	if err != nil {
		return fmt.Errorf("failed to get org name:%s for adding member, %w", cmd.Org, err)
	}

	approve, err := o.RemoveInvite(cmd.ToMember())
	if err != nil {
		return fmt.Errorf("failed to remove invite for adding member, %w", err)
	}

	if approve.ExpireAt < utils.Now() {
		_, _ = org.repo.Save(&o) // just remove the expired invitation
		return fmt.Errorf("invitation already expired")
	}

	m := approve.ToMember()

	pl, err := org.user.GetPlatformUser(o.Owner)
	if err != nil {
		return fmt.Errorf("failed to get platform user for adding member, %w", err)
	}

	_, err = org.user.GetByAccount(cmd.Account, false)
	if err != nil {
		return fmt.Errorf("failed to get user for adding member, %w", err)
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

	_, err = org.repo.Save(&o)
	if err != nil {
		// TODO need rollback
		return fmt.Errorf("failed to save org for adding member, %w", err)
	}

	return nil
}

func (org *orgService) canRemoveMember(cmd *OrgRemoveMemberCmd) (err error) {
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
		if m.Role == domain.OrgRoleOwner {
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
		err = fmt.Errorf("the only owner can not be removed")
		return
	}

	if !can {
		err = fmt.Errorf("the member is not in the org")
		return
	}

	return
}

func (org *orgService) RemoveMember(cmd *OrgRemoveMemberCmd) error {
	err := cmd.Validate()
	if err != nil {
		return err
	}

	err = org.canRemoveMember(cmd)
	if err != nil {
		return err
	}

	o, err := org.repo.GetByName(cmd.Org)
	if err != nil {
		return fmt.Errorf("failed to get org, %w", err)
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

func (org *orgService) EditMember(cmd *OrgEditMemberCmd) (dto MemberDTO, err error) {
	m, err := org.member.GetByOrgAndUser(cmd.Org.Account(), cmd.Account.Account())
	if err != nil {
		err = fmt.Errorf("failed to get member by org:%s and user:%s, %w", cmd.Org.Account(), cmd.Account.Account(), err)
		return
	}

	o, err := org.repo.GetByName(cmd.Org)
	if err != nil {
		err = fmt.Errorf("failed to get org, %w", err)
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

func (org *orgService) InviteMember(cmd *OrgInviteMemberCmd) (dto OrganizationDTO, err error) {
	if err = cmd.Validate(); err != nil {
		return
	}

	o, err := org.repo.GetByName(cmd.Org)
	if err != nil {
		err = fmt.Errorf("failed to get org, %w", err)
		return
	}

	_, err = org.member.GetByOrgAndUser(cmd.Org.Account(), cmd.Account.Account())
	if err == nil {
		err = fmt.Errorf("user %s is already a member of org %s", cmd.Account.Account(), cmd.Org.Account())
		return
	}

	if err = o.AddInvite(cmd.ToMember(), org.inviteExpiry); err != nil {
		err = fmt.Errorf("failed to add invite, %w", err)
		return
	}

	newOrg, err := org.repo.Save(&o)
	if err != nil {
		err = fmt.Errorf("failed to save member, %w", err)
		return
	}

	dto = ToDTO(&newOrg)

	return
}

func (org *orgService) RevokeInvite(cmd *OrgRemoveInviteCmd) (dto OrganizationDTO, err error) {
	if err = cmd.Validate(); err != nil {
		return
	}

	o, err := org.repo.GetByName(cmd.Org)
	if err != nil {
		err = fmt.Errorf("failed to get org, %w", err)
		return
	}

	_, err = o.RemoveInvite(cmd.ToMember())
	if err != nil {
		err = fmt.Errorf("failed to remove invite, %w", err)
		return
	}

	newOrg, err := org.repo.Save(&o)
	if err != nil {
		err = fmt.Errorf("failed to save member, %w", err)
		return
	}

	dto = ToDTO(&newOrg)

	return
}

func (org *orgService) ListInvitation(acc primitive.Account) (dtos []ApproveDTO, err error) {
	if acc == nil {
		err = fmt.Errorf("account is nil")
		return
	}

	o, err := org.repo.GetByName(acc)
	if err != nil {
		err = fmt.Errorf("failed to get org, %w", err)
		return
	}

	if o.Approves == nil {
		logrus.Warnf("no invitation found for %s", acc.Account())
		return
	}

	dtos = make([]ApproveDTO, len(o.Approves))
	for i := range o.Approves {
		dtos[i] = ToApproveDTO(o.Approves[i])
	}

	return
}

func (org *orgService) CheckName(name primitive.Account) bool {
	if name == nil {
		logrus.Error("name is nil")
		return false
	}
	_, err1 := org.repo.GetByName(name)

	_, err2 := org.user.GetByAccount(name, false)

	if repositories.IsErrorDataNotExists(err1) && commonrepo.IsErrorResourceNotExists(err2) {
		return true
	}

	return false
}
