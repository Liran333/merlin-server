/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package e2e

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	swaggerRest "e2e/client_rest"
)

// SuiteOrgModify used for testing
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

// SetupSuite used for testing
func (s *SuiteOrgModify) SetupSuite() {
	s.name = "testorg"
	s.fullname = "testorgfull"
	s.avatarid = "https://avatars.githubusercontent.com/u/2853724?v=1"
	s.allowRequest = true
	s.defaultRole = "admin"
	s.website = "https://www.modelfoundry.cn"
	s.desc = "test org desc"
	s.owner = "test1" // this name is hard code in init-env.sh

	data, r, err := ApiRest.UserApi.V1UserGet(AuthRest)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	user := getData(s.T(), data.Data)

	assert.NotEqual(s.T(), "", user["id"])
	s.owerId = getString(s.T(), user["id"])
	s.T().Logf("owerId: %s", s.owerId)
}

// TearDownSuite used for testing
func (s *SuiteOrgModify) TearDownSuite() {

}

// TestOrgCreate used for testing
// 组织管理员修改组织的名称
func (s *SuiteOrgModify) TestOrgCreate() {
	//创建组织
	d := swaggerRest.ControllerOrgCreateRequest{
		Name:     s.name,
		Fullname: s.fullname,
	}

	_, r, err := ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	// 修改组织的名称
	d2 := swaggerRest.ControllerOrgBasicInfoUpdateRequest{
		Fullname: "newFullName",
	}

	_, r2, err2 := ApiRest.OrganizationApi.V1OrganizationNamePut(AuthRest, s.name, d2)
	assert.Equal(s.T(), http.StatusAccepted, r2.StatusCode)
	assert.Nil(s.T(), err2)

	r3, err3 := ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest, s.name)
	assert.Equal(s.T(), http.StatusNoContent, r3.StatusCode)
	assert.Nil(s.T(), err3)
}

// TestOrgDeleteFail used for testing
// 其他人无法删除组织
func (s *SuiteOrgModify) TestOrgDeleteFail() {
	r, err := ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest2, s.name)
	assert.Equal(s.T(), http.StatusForbidden, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestOrgModify used for testing
func TestOrgModify(t *testing.T) {
	suite.Run(t, new(SuiteOrgModify))
}
