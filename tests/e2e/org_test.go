package e2e

import (
	"context"
	swagger "e2e/client"
	"testing"

	"github.com/antihax/optional"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SuiteOrg struct {
	suite.Suite
	name         string
	fullname     string
	avatarid     string
	allowRequest bool
	defaultRole  string
	website      string
	desc         string
	owner        string
	owerId       string
}

func (s *SuiteOrg) SetupSuite() {
	s.name = "testorg"
	s.fullname = "testorgfull"
	s.avatarid = "https://avatars.githubusercontent.com/u/2853724?v=1"
	s.allowRequest = true
	s.defaultRole = "admin"
	s.website = "https://www.modelfoundry.cn"
	s.desc = "test org desc"
	s.owner = "test1" // this name is hard code in init-env.sh

	data, r, err := Api.UserApi.V1UserGet(Auth)
	assert.Equal(s.T(), 200, r.StatusCode)
	assert.Nil(s.T(), err)

	user := getData(s.T(), data.Data)

	assert.NotEqual(s.T(), "", user["id"])
	s.owerId = getString(s.T(), user["id"])
	s.T().Logf("owerId: %s", s.owerId)
}

func (s *SuiteOrg) TearDownSuite() {

}

// 正常创建一个组织
func (s *SuiteOrg) TestOrgCreate() {
	d := swagger.ControllerOrgCreateRequest{
		Name:     s.name,
		Fullname: s.fullname,
	}

	data, r, err := Api.OrganizationApi.V1OrganizationPost(Auth, d)
	assert.Equal(s.T(), 201, r.StatusCode)
	assert.Nil(s.T(), err)

	org := getData(s.T(), data.Data)

	assert.Equal(s.T(), s.fullname, org["fullname"])
	assert.Equal(s.T(), s.name, org["account"])
	assert.NotEqual(s.T(), "", org["id"])
	assert.NotEqual(s.T(), 0, org["created_at"])
	assert.NotEqual(s.T(), 0, org["updated_at"])
	assert.Equal(s.T(), "test1", org["owner"])
	assert.Equal(s.T(), s.owerId, org["owner_id"])
	assert.Equal(s.T(), int64(1), getInt64(s.T(), org["type"]))
	assert.Equal(s.T(), "write", org["default_role"])
	assert.Equal(s.T(), "", org["avatar_id"])
	assert.Equal(s.T(), "", org["website"])
	assert.Equal(s.T(), "", org["description"])
	allow, ok := org["allow_request"].(bool)
	assert.Equal(s.T(), true, ok)
	assert.Equal(s.T(), false, allow)
	assert.Nil(s.T(), org["email"])

	data, r, err = Api.OrganizationApi.V1OrganizationGet(Auth, &swagger.OrganizationApiV1OrganizationGetOpts{})
	assert.Equal(s.T(), 200, r.StatusCode)
	assert.Nil(s.T(), err)

	count := 0
	orgs := getArrary(s.T(), data.Data)
	for _, v := range orgs {
		if v != nil {
			count++
		}
	}
	assert.Equal(s.T(), 1, count)

	r, err = Api.OrganizationApi.V1OrganizationNameDelete(Auth, s.name)
	assert.Equal(s.T(), 204, r.StatusCode)
	assert.Nil(s.T(), err)
}

// // 创建组织失败
// 未登录用户
func (s *SuiteOrg) TestOrgCreateFailedNoToken() {
	d := swagger.ControllerOrgCreateRequest{
		Name:     s.name,
		Fullname: s.fullname,
	}

	_, r, err := Api.OrganizationApi.V1OrganizationPost(context.Background(), d)
	assert.Equal(s.T(), 401, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// 无效的组织名：名字过长
func (s *SuiteOrg) TestOrgCreateFailedInvalidNameLen() {
	d := swagger.ControllerOrgCreateRequest{
		Name: string(make([]byte, 51)),
	}

	_, r, err := Api.OrganizationApi.V1OrganizationPost(Auth, d)
	assert.Equal(s.T(), 400, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// 组织名已存在
func (s *SuiteOrg) TestOrgCreateFailedInvalidNameConflict() {
	d := swagger.ControllerOrgCreateRequest{
		Name:     s.owner,
		Fullname: s.fullname,
	}

	_, r, err := Api.OrganizationApi.V1OrganizationPost(Auth, d)
	assert.Equal(s.T(), 400, r.StatusCode)
	assert.NotNil(s.T(), err)

	r, err = Api.NameApi.V1NameHead(Auth, s.owner)
	assert.Equal(s.T(), 409, r.StatusCode)
	assert.NotNil(s.T(), err)

	r, err = Api.NameApi.V1NameHead(Auth, "testnonexist")
	assert.Equal(s.T(), 200, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 组织名是保留名称
func (s *SuiteOrg) TestOrgCreateFailedInvalidNameReserved() {
	d := swagger.ControllerOrgCreateRequest{
		Name:     "models",
		Fullname: s.fullname,
	}

	_, r, err := Api.OrganizationApi.V1OrganizationPost(Auth, d)
	assert.Equal(s.T(), 400, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// 空fullname
func (s *SuiteOrg) TestOrgCreateFailedEmptyFullname() {
	d := swagger.ControllerOrgCreateRequest{
		Name: s.name,
	}

	data, r, err := Api.OrganizationApi.V1OrganizationPost(Auth, d)
	assert.Equalf(s.T(), 400, r.StatusCode, data.Msg)
	assert.NotNil(s.T(), err)

	d = swagger.ControllerOrgCreateRequest{
		Name:     s.name,
		Fullname: "",
	}

	data, r, err = Api.OrganizationApi.V1OrganizationPost(Auth, d)
	assert.Equalf(s.T(), 400, r.StatusCode, data.Msg)
	assert.NotNil(s.T(), err)
}

// 无效的avatarid
func (s *SuiteOrg) TestOrgCreateFailedInvalidAvatarid() {
	d := swagger.ControllerOrgCreateRequest{
		Name:     s.name,
		Fullname: s.fullname,
		AvatarId: "invalid",
	}

	_, r, err := Api.OrganizationApi.V1OrganizationPost(Auth, d)
	assert.Equal(s.T(), 400, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// 无效的desc
func (s *SuiteOrg) TestOrgCreateFailedInvalidDesc() {
	d := swagger.ControllerOrgCreateRequest{
		Name:        s.name,
		Fullname:    s.fullname,
		Description: string(make([]byte, 201)),
	}

	_, r, err := Api.OrganizationApi.V1OrganizationPost(Auth, d)
	assert.Equal(s.T(), 400, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// 无效的website
func (s *SuiteOrg) TestOrgCreateFailedInvalidWebsite() {
	d := swagger.ControllerOrgCreateRequest{
		Name:     s.name,
		Fullname: s.fullname,
		Website:  "google.com",
	}

	data, r, err := Api.OrganizationApi.V1OrganizationPost(Auth, d)
	assert.Equalf(s.T(), 400, r.StatusCode, data.Msg)
	assert.NotNil(s.T(), err)
}

// 名下无组织时，查询个人组织返回一个空列表
func (s *SuiteOrg) TestOrgListEmpty() {
	_, _ = Api.OrganizationApi.V1OrganizationNameDelete(Auth, s.name)

	// list by owner
	d := swagger.OrganizationApiV1OrganizationGetOpts{
		Owner: optional.NewString(s.owner),
	}
	data, r, err := Api.OrganizationApi.V1OrganizationGet(Auth, &d)
	assert.Equal(s.T(), 200, r.StatusCode)
	assert.Nil(s.T(), err)

	orgs := getArrary(s.T(), data.Data)

	count := 0

	for org := range orgs {
		if orgs[org] == nil {
			continue
		}
		count += 1
	}

	assert.Equal(s.T(), 0, count)

	// list by member user
	d = swagger.OrganizationApiV1OrganizationGetOpts{
		Username: optional.NewString(s.owner),
	}
	data, r, err = Api.OrganizationApi.V1OrganizationGet(Auth, &d)
	assert.Equal(s.T(), 200, r.StatusCode)
	assert.Nil(s.T(), err)

	orgs = getArrary(s.T(), data.Data)

	count = 0

	for org := range orgs {
		if orgs[org] == nil {
			continue
		}
		count += 1
	}

	assert.Equal(s.T(), 0, count)

	// list all
	d = swagger.OrganizationApiV1OrganizationGetOpts{}
	data, r, err = Api.OrganizationApi.V1OrganizationGet(Auth, &d)
	assert.Equal(s.T(), 200, r.StatusCode)
	assert.Nil(s.T(), err)

	orgs = getArrary(s.T(), data.Data)

	count = 0

	for org := range orgs {
		if orgs[org] == nil {
			continue
		}
		count += 1
	}

	assert.Equal(s.T(), 0, count)
}

// 查询不存在的组织
func (s *SuiteOrg) TestOrgNonexist() {
	_, r, err := Api.OrganizationApi.V1OrganizationNameGet(Auth, "nonexist")
	assert.Equal(s.T(), 404, r.StatusCode)
	assert.NotNil(s.T(), err)

	_, r, err = Api.OrganizationApi.V1OrganizationNameGet(context.Background(), "nonexist")
	assert.Equal(s.T(), 404, r.StatusCode)
	assert.NotNil(s.T(), err)
}

func TestOrg(t *testing.T) {
	suite.Run(t, new(SuiteOrg))
}
