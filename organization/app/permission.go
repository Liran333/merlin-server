/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	orgdomain "github.com/openmerlin/merlin-server/organization/domain"
	perm "github.com/openmerlin/merlin-server/organization/domain/permission"
	"github.com/openmerlin/merlin-server/organization/domain/repository"
	user "github.com/openmerlin/merlin-server/user/domain"
)

// NewPermService creates a new PermService instance with the given configuration and org member.
func NewPermService(cfg *perm.Config, org repository.OrgMember) *permService {
	p := &permService{
		org: org,
	}
	p.Init(cfg)

	return p
}

func initActioin(actions []string) (bitmap uint64) {
	for _, action := range actions {
		switch action {
		case "write":
			bitmap |= 1 << primitive.ActionWrite
		case "read":
			bitmap |= 1 << primitive.ActionRead
		case "delete":
			bitmap |= 1 << primitive.ActionDelete
		case "create":
			bitmap |= 1 << primitive.ActionCreate
		default:
			logrus.Fatalf("invalid action: %s", action)
		}
	}

	return
}

func checkAction(bitmap uint64, action primitive.Action) bool {
	return bitmap&(1<<action) != 0
}

func isValidRole(role string) bool {
	if role == string(user.OrgRoleAdmin) || role == string(user.OrgRoleWriter) ||
		role == string(user.OrgRoleContributor) || role == string(user.OrgRoleReader) {
		return true
	}

	return false
}

// Init initializes the permission service with the given configuration.
func (p *permService) Init(cfg *perm.Config) {
	p.permissions = make(map[primitive.ObjType]map[orgdomain.OrgRole]uint64)
	for _, permission := range cfg.Permissions {
		r := make(map[orgdomain.OrgRole]uint64)
		for _, rule := range permission.Rules {
			if !isValidRole(rule.Role) {
				logrus.Fatalf("invalid role: %s", rule.Role)
			}

			r[orgdomain.OrgRole(rule.Role)] = initActioin(rule.Operation)
		}

		p.permissions[primitive.ObjType(permission.ObjectType)] = r
	}
}

type permService struct {
	permissions map[primitive.ObjType]map[orgdomain.OrgRole]uint64

	org repository.OrgMember
}

func (p *permService) doCheckPerm(role string, objType primitive.ObjType, op primitive.Action) bool {

	if v, ok := p.permissions[objType][orgdomain.OrgRole(role)]; ok {
		if checkAction(v, op) {
			return true
		}
	}

	return false
}

// subject: a user session or a token sessioin
// object: org
// op: write, read
func (p *permService) checkInOrg(
	user primitive.Account,
	obj primitive.Account,
	objType primitive.ObjType,
	op primitive.Action,
) error {
	return p.doCheck(user, obj, objType, op, nil)
}

// Check checks if the user has the permission to perform the operation on the object.
// subject: a user session or a token sessioin
// object: model, dataset, space, system
// op: write, read
func (p *permService) Check(
	user primitive.Account,
	obj primitive.Account,
	objType primitive.ObjType,
	op primitive.Action,
	createdByUser bool,
) error {
	if op.IsModification() {
		return p.doCheck(user, obj, objType, op, func() bool {
			return createdByUser
		})
	}

	return p.doCheck(user, obj, objType, op, nil)
}

// subject: a user session or a token sessioin
// object: model, dataset, space, system
// op: write, read
func (p *permService) doCheck(
	user primitive.Account,
	obj primitive.Account,
	objType primitive.ObjType,
	op primitive.Action,
	judgeInAdvance func() bool,
) error {
	if user == nil {
		return allerror.NewNoPermission("user is nil")
	}

	if obj == nil {
		return allerror.NewNoPermission("object is nil")
	}

	if user == obj {
		return nil
	}

	m, err := p.org.GetByOrgAndUser(obj.Account(), user.Account())
	if err != nil {
		logrus.Errorf("get member failed: %s", err)

		return allerror.NewNoPermission(fmt.Sprintf(
			"%s does not have a valid role in %s", user.Account(), obj.Account(),
		))
	}

	if judgeInAdvance != nil && judgeInAdvance() {
		return nil
	}

	ok := p.doCheckPerm(string(m.Role), objType, op)
	res := "cannot"
	if ok {
		res = "can"
	}

	logrus.Debugf(
		"user %s (role %s) %s do %d on %s:%s",
		user.Account(), m.Role, res, op, obj.Account(), objType,
	)

	if !ok {
		return allerror.NewNoPermission(fmt.Sprintf(
			"%s %s %s permission denied", user.Account(), op.String(), string(objType),
		))
	}

	return nil
}
