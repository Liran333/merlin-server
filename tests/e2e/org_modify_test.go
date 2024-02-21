package e2e

import (
	swagger "e2e/client"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SuiteOrgModify struct {
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

func (s *SuiteOrgModify) SetupSuite() {
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

func (s *SuiteOrgModify) TearDownSuite() {

}

// 组织管理员修改组织的名称
func (s *SuiteOrgModify) TestOrgCreate() {
	//创建组织
	d := swagger.ControllerOrgCreateRequest{
		Name:     s.name,
		Fullname: s.fullname,
	}

	_, r, err := Api.OrganizationApi.V1OrganizationPost(Auth, d)
	assert.Equal(s.T(), 201, r.StatusCode)
	assert.Nil(s.T(), err)

	// 修改组织的名称
	d2 := swagger.ControllerOrgBasicInfoUpdateRequest{
		Fullname: "newFullName",
	}

	_, r2, err2 := Api.OrganizationApi.V1OrganizationNamePut(Auth, s.name, d2)
	assert.Equal(s.T(), 202, r2.StatusCode)
	assert.Nil(s.T(), err2)

	r3, err3 := Api.OrganizationApi.V1OrganizationNameDelete(Auth, s.name)
	assert.Equal(s.T(), 204, r3.StatusCode)
	assert.Nil(s.T(), err3)
}

// 其他人无法删除组织
func (s *SuiteOrgModify) TestOrgDeleteFail() {
	r, err := Api.OrganizationApi.V1OrganizationNameDelete(Auth2, s.name)
	assert.Equal(s.T(), 403, r.StatusCode)
	assert.NotNil(s.T(), err)
}

func TestOrgModify(t *testing.T) {
	suite.Run(t, new(SuiteOrgModify))
}
