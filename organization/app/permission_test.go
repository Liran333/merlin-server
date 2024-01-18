package app

import (
	"fmt"
	"testing"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	orgdomain "github.com/openmerlin/merlin-server/organization/domain"
	"github.com/openmerlin/merlin-server/organization/domain/permission"
	user "github.com/openmerlin/merlin-server/user/domain"
)

type stubOrg struct{}

var stubMember = &orgdomain.OrgMember{
	Id:       primitive.CreateIdentity(1),
	OrgName:  primitive.CreateAccount("org"),
	Username: primitive.CreateAccount("XXXX"),
	Role:     user.OrgRoleAdmin,
}

var stub1Member = &orgdomain.OrgMember{
	Id:       primitive.CreateIdentity(1),
	OrgName:  primitive.CreateAccount("member"),
	Username: primitive.CreateAccount("ffff"),
	Role:     user.OrgRoleWriter,
}

func (s *stubOrg) Add(o *orgdomain.OrgMember) (orgdomain.OrgMember, error) {
	return *o, nil
}

func (s *stubOrg) Save(o *orgdomain.OrgMember) (orgdomain.OrgMember, error) {
	return *o, nil
}

func (s *stubOrg) Delete(*orgdomain.OrgMember) error {
	return nil
}

func (s *stubOrg) DeleteByOrg(primitive.Account) error {
	return nil
}

func (s *stubOrg) GetByOrg(u string) ([]orgdomain.OrgMember, error) {
	return []orgdomain.OrgMember{*stubMember, *stub1Member}, nil
}

func (s *stubOrg) GetByOrgAndRole(u string, r orgdomain.OrgRole) ([]orgdomain.OrgMember, error) {
	return []orgdomain.OrgMember{*stubMember}, nil
}

func (s *stubOrg) GetByOrgAndUser(org, user string) (orgdomain.OrgMember, error) {
	if org == "org" && user == "XXXX" {
		return *stubMember, nil
	} else if org == "member" && user == "ffff" {
		return *stub1Member, nil
	}
	return orgdomain.OrgMember{}, fmt.Errorf("not found")
}

func (s *stubOrg) GetByUser(string) ([]orgdomain.OrgMember, error) {
	return []orgdomain.OrgMember{*stubMember}, nil

}

func TestPermCheck(t *testing.T) {
	type testdata struct {
		user    primitive.Account
		org     primitive.Account
		objType primitive.ObjType
		op      primitive.Action
	}

	results := []bool{
		false,
		false,
		false,
		true,
		false,
		false,
		true,
		true,
	}

	tests := []testdata{
		{
			user:    nil,
			org:     nil,
			objType: primitive.ObjTypeOrg,
			op:      primitive.ActionRead,
		},
		{
			user:    primitive.CreateAccount("123"),
			org:     nil,
			objType: primitive.ObjTypeOrg,
			op:      primitive.ActionRead,
		},
		{
			user:    primitive.CreateAccount("123"),
			org:     nil,
			objType: primitive.ObjTypeOrg,
			op:      primitive.ActionRead,
		},
		{
			user:    primitive.CreateAccount("XXXX"),
			org:     primitive.CreateAccount("org"),
			objType: primitive.ObjTypeOrg,
			op:      primitive.ActionCreate,
		},
		{
			user:    primitive.CreateAccount("ffff"),
			org:     primitive.CreateAccount("member"),
			objType: primitive.ObjTypeMember,
			op:      primitive.ActionDelete,
		},
		{
			user:    primitive.CreateAccount("ffff"),
			org:     primitive.CreateAccount("member"),
			objType: primitive.ObjTypeMember,
			op:      primitive.ActionCreate,
		},
		{
			user:    primitive.CreateAccount("ffff"),
			org:     primitive.CreateAccount("member"),
			objType: primitive.ObjTypeMember,
			op:      primitive.ActionRead,
		},
		{
			user:    primitive.CreateAccount("ffff"),
			org:     primitive.CreateAccount("member"),
			objType: primitive.ObjTypeMember,
			op:      primitive.ActionWrite,
		},
	}
	var cfg permission.Config
	cfg.Permissions = []permission.PermObject{
		{
			ObjectType: string(primitive.ObjTypeOrg),
			Rules: []permission.Rule{
				{
					Role:      string(user.OrgRoleAdmin),
					Operation: []string{"write", "read", "create", "delete"},
				},
			},
		},
		{
			ObjectType: string(primitive.ObjTypeMember),
			Rules: []permission.Rule{
				{
					Role:      string(user.OrgRoleWriter),
					Operation: []string{"write", "read"},
				},
			},
		},
	}
	app := NewPermService(&cfg, &stubOrg{})

	for i, test := range tests {
		err := app.Check(test.user, test.org, test.objType, test.op)
		if (err == nil) != results[i] {
			t.Errorf("case num %d valid result is %v ", i, err)
		}

	}
}
