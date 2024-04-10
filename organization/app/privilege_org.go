package app

import (
	"fmt"

	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/organization/domain/privilege"
	"github.com/sirupsen/logrus"
)

type action string

const (
	// AllocNpu org who can alloc npu resource
	AllocNpu action = "alloc_npu"
	// Disable org who can disable model or space
	Disable action = "disable"
)

// NewAction returns an action type
func NewAction(action string) (action, error) {
	switch action {
	case "alloc_npu":
		return AllocNpu, nil
	case "disable":
		return Disable, nil
	default:
		return "", fmt.Errorf("invalid action: %s", action)
	}
}

// PrivilegeOrg interface
type PrivilegeOrg interface {
	Contains(primitive.Account) error
}

func NewPrivilegeOrgService(
	org OrgService,
	cfg privilege.PrivilegeConfig,
) PrivilegeOrg {
	a, err := NewAction(cfg.Type)
	if err != nil {
		logrus.Warnf("invalid privileged type: %s, privilege org will be ignored", cfg.Type)
		return nil
	}

	acc, err := primitive.NewAccount(cfg.OrgName)
	if err != nil {
		logrus.Warnf("invalid privileged org name: %s, privilege org will be ignored", cfg.OrgName)
		return nil
	}

	return &privilegeOrg{
		OrgId:   cfg.OrgId,
		OrgName: acc,
		Action:  a,
		org:     org,
	}
}

type privilegeOrg struct {
	OrgId   string
	OrgName primitive.Account
	Action  action
	org     OrgService
}

// Contains returns an error if the account is not in the privilege org.
func (p *privilegeOrg) Contains(account primitive.Account) error {
	o, err := p.org.GetByAccount(p.OrgName)
	if err != nil {
		logrus.Errorf("cant do %s action while failed to get org: %s, %s", p.Action, p.OrgName, err)
		return err
	}

	if o.Id != p.OrgId {
		e := fmt.Errorf("org id mismatch: actual: %s, config: %s", o.Id, p.OrgId)
		return allerror.New(allerror.ErrorCodePrivilegeOrgIdMismatch, e.Error(), e)
	}

	has := p.org.HasMember(primitive.CreateAccount(o.Name), account)
	if !has {
		e := fmt.Errorf("user: %s not in a %v org", account.Account(), p.Action)
		return allerror.New(allerror.ErrorCodeNotInPrivilegeOrg, e.Error(), e)
	}

	return nil
}
