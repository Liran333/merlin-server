/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package e2e

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/antihax/optional"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	swaggerRest "e2e/client_rest"
)

// SuiteOrg used for testing
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
	invitee      string
}

const (
	charA = "A"
)

// SetupSuite used for testing
func (s *SuiteOrg) SetupSuite() {
	s.name = "testorg"
	s.fullname = "testorgfull"
	s.avatarid = "https://avatars.githubusercontent.com/u/2853724?v=1"
	s.allowRequest = true
	s.defaultRole = "admin"
	s.website = "https://www.modelfoundry.cn"
	s.desc = "test org desc"
	s.owner = "test1"   // this name is hard code in init-env.sh
	s.invitee = "test2" // this name is hard code in init-env.sh

	data, r, err := ApiRest.UserApi.V1UserGet(AuthRest)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	user := getData(s.T(), data.Data)

	assert.NotEqual(s.T(), "", user["id"])
	s.owerId = getString(s.T(), user["id"])
	s.T().Logf("owerId: %s", s.owerId)
}

// TearDownSuite used for testing
func (s *SuiteOrg) TearDownSuite() {

}

// TestOrgCreate used for testing
// 正常创建一个组织
func (s *SuiteOrg) TestOrgCreate() {
	d := swaggerRest.ControllerOrgCreateRequest{
		Name:     s.name,
		Fullname: s.fullname,
	}

	data, r, err := ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
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
	assert.Equal(s.T(), "", org["email"])

	orgData, r, err := ApiRest.OrganizationApi.V1OrganizationGet(AuthRest,
		&swaggerRest.OrganizationApiV1OrganizationGetOpts{})
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	count := 0
	orgs := getArrary(s.T(), orgData.Data)
	for _, v := range orgs {
		if v != nil {
			count++
		}
	}
	assert.Equal(s.T(), countOne, count)

	// 重复创建组织返回400
	_, r, err = ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	r, err = ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest, s.name)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestOrgCreateSuccess used for testing
// 创建一个组织成功，website为非必选项
func (s *SuiteOrg) TestOrgCreateSuccess() {
	// website为非必选项
	d := swaggerRest.ControllerOrgCreateRequest{
		Name:     s.name,
		Fullname: s.fullname,
	}
	_, r, err := ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest, s.name)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	// 创建组织，website设置成功
	d = swaggerRest.ControllerOrgCreateRequest{
		Name:     s.name,
		Fullname: s.fullname,
		Website:  s.website,
	}
	_, r, err = ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	data, r, err := ApiRest.OrganizationApi.V1OrganizationGet(AuthRest,
		&swaggerRest.OrganizationApiV1OrganizationGetOpts{})
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	orgs := getArrary(s.T(), data.Data)
	count := 0
	for _, v := range orgs {
		if v["fullname"] == s.fullname {
			assert.Equal(s.T(), s.website, v["website"])
			count++
		}
	}
	assert.Equal(s.T(), countOne, count)

	r, err = ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest, s.name)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestOrgCreateFail used for testing
// 创建一个组织名、ID必填
func (s *SuiteOrg) TestOrgCreateFail() {
	// 组织名必填
	d := swaggerRest.ControllerOrgCreateRequest{
		Name:     s.name,
		Fullname: "",
	}

	_, r, err := ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode, "org fullname can't be empty")
	assert.NotNil(s.T(), err)

	// 组织ID必填
	d = swaggerRest.ControllerOrgCreateRequest{
		Name:     "",
		Fullname: s.fullname,
	}

	_, r, err = ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode, "name can't be empty")
	assert.NotNil(s.T(), err)
}

// TestOrgCreateFailedNoToken used for testing
// 创建组织失败
// 未登录用户
func (s *SuiteOrg) TestOrgCreateFailedNoToken() {
	d := swaggerRest.ControllerOrgCreateRequest{
		Name:     s.name,
		Fullname: s.fullname,
	}

	_, r, err := ApiRest.OrganizationApi.V1OrganizationPost(context.Background(), d)
	assert.Equal(s.T(), http.StatusUnauthorized, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestOrgCreateFailedInvalidNameLen used for testing
// 无效的组织名：名字过长
func (s *SuiteOrg) TestOrgCreateFailedInvalidNameLen() {
	d := swaggerRest.ControllerOrgCreateRequest{
		Name: string(make([]byte, length)),
	}

	_, r, err := ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestOrgCreateFailedInvalidNameConflict used for testing
// 组织名已存在
func (s *SuiteOrg) TestOrgCreateFailedInvalidNameConflict() {
	d := swaggerRest.ControllerOrgCreateRequest{
		Name:     s.name,
		Fullname: s.fullname,
	}

	_, r, err := ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	// 与已有组织名重复
	r, err = ApiRest.OrganizationApi.V1NameHead(AuthRest, s.name)
	assert.Equal(s.T(), http.StatusConflict, r.StatusCode)
	assert.NotNil(s.T(), err)

	r, err = ApiRest.OrganizationApi.V1NameHead(AuthRest, "TestOrg")
	assert.Equal(s.T(), http.StatusConflict, r.StatusCode)
	assert.NotNil(s.T(), err)

	// 与已有用户名重复
	r, err = ApiRest.OrganizationApi.V1NameHead(AuthRest, s.owner)
	assert.Equal(s.T(), http.StatusConflict, r.StatusCode)
	assert.NotNil(s.T(), err)

	r, err = ApiRest.OrganizationApi.V1NameHead(AuthRest, "Test1")
	assert.Equal(s.T(), http.StatusConflict, r.StatusCode)
	assert.NotNil(s.T(), err)

	// 不与已有组织名或用户名重复
	r, err = ApiRest.OrganizationApi.V1NameHead(AuthRest, "testnonexist")
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest, s.name)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestOrgCreateFailedInvalidNameReserved used for testing
// 组织名是保留名称
func (s *SuiteOrg) TestOrgCreateFailedInvalidNameReserved() {
	d := swaggerRest.ControllerOrgCreateRequest{
		Name:     "models",
		Fullname: s.fullname,
	}

	_, r, err := ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestOrgCreateFailedEmptyFullname used for testing
// 空fullname
func (s *SuiteOrg) TestOrgCreateFailedEmptyFullname() {
	d := swaggerRest.ControllerOrgCreateRequest{
		Name: s.name,
	}

	data, r, err := ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equalf(s.T(), http.StatusBadRequest, r.StatusCode, data.Msg)
	assert.NotNil(s.T(), err)

	d = swaggerRest.ControllerOrgCreateRequest{
		Name:     s.name,
		Fullname: "",
	}

	data, r, err = ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equalf(s.T(), http.StatusBadRequest, r.StatusCode, data.Msg)
	assert.NotNil(s.T(), err)
}

// TestOrgCreateFailedInvalidAvatarid used for testing
// 无效的avatarid
func (s *SuiteOrg) TestOrgCreateFailedInvalidAvatarid() {
	d := swaggerRest.ControllerOrgCreateRequest{
		Name:     s.name,
		Fullname: s.fullname,
		AvatarId: "invalid",
	}

	_, r, err := ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestOrgCreateFailedInvalidDesc used for testing
// 无效的desc
func (s *SuiteOrg) TestOrgCreateFailedInvalidDesc() {
	d := swaggerRest.ControllerOrgCreateRequest{
		Name:        s.name,
		Fullname:    s.fullname,
		Description: string(make([]byte, 256)),
	}

	_, r, err := ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestOrgCreateFailedInvalidWebsite used for testing
// 无效的website
func (s *SuiteOrg) TestOrgCreateFailedInvalidWebsite() {
	d := swaggerRest.ControllerOrgCreateRequest{
		Name:     s.name,
		Fullname: s.fullname,
		Website:  "google.com",
	}

	data, r, err := ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equalf(s.T(), http.StatusBadRequest, r.StatusCode, data.Msg)
	assert.NotNil(s.T(), err)

	// website too long
	d = swaggerRest.ControllerOrgCreateRequest{
		Name:     s.name,
		Fullname: s.fullname,
		Website:  strings.Repeat(charA, ComConfig.WEBSITE_MAX_LEN+1),
	}

	data, r, err = ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equalf(s.T(), http.StatusBadRequest, r.StatusCode, data.Msg)
	assert.NotNil(s.T(), err)

	// website repeate
	d = swaggerRest.ControllerOrgCreateRequest{
		Name:     s.name,
		Fullname: s.fullname,
		Website:  s.website,
	}

	_, r, err = ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	data, r, err = ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	r, err = ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest, s.name)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestOrgListEmpty used for testing
// 名下无组织时，查询个人组织返回一个空列表
func (s *SuiteOrg) TestOrgListEmpty() {
	// make sure the org is not exist
	_, _ = ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest, s.name)

	// list by owner
	d := swaggerRest.OrganizationApiV1OrganizationGetOpts{
		Owner: optional.NewString(s.owner),
	}
	data, r, err := ApiRest.OrganizationApi.V1OrganizationGet(AuthRest, &d)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
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
	d = swaggerRest.OrganizationApiV1OrganizationGetOpts{
		Username: optional.NewString(s.owner),
	}
	data, r, err = ApiRest.OrganizationApi.V1OrganizationGet(AuthRest, &d)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
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
	d = swaggerRest.OrganizationApiV1OrganizationGetOpts{}
	data, r, err = ApiRest.OrganizationApi.V1OrganizationGet(AuthRest, &d)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
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

// TestOrgNonexist used for testing
// 查询不存在的组织
func (s *SuiteOrg) TestOrgNonexist() {
	_, r, err := ApiRest.OrganizationApi.V1OrganizationNameGet(AuthRest, "nonexist")
	assert.Equal(s.T(), http.StatusNotFound, r.StatusCode)
	assert.NotNil(s.T(), err)

	_, r, err = ApiRest.OrganizationApi.V1OrganizationNameGet(context.Background(), "nonexist")
	assert.Equal(s.T(), http.StatusNotFound, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestOrgOwnerSearch used for testing
// 用户可以批量查询某个用户所拥有的组织
func (s *SuiteOrgModel) TestOrgOwnerSearch() {
	// list all
	d := swaggerRest.OrganizationApiV1OrganizationGetOpts{Owner: optional.NewString("test1")}
	data, r, err := ApiRest.OrganizationApi.V1OrganizationGet(AuthRest, &d)

	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), data.Data)
}

// TestOrgMemberSearch used for testing
// 用户可以批量查询某个用户所属的组织
func (s *SuiteOrgModel) TestOrgMemberSearch() {
	// list all
	d := swaggerRest.OrganizationApiV1OrganizationGetOpts{Username: optional.NewString("test1")}
	data, r, err := ApiRest.OrganizationApi.V1OrganizationGet(AuthRest, &d)

	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), data.Data)
}

// TestOrgRolesSearch used for testing
// 用户可以批量查询某个用户所拥有权限的组织
func (s *SuiteOrgModel) TestOrgRolesSearch() {
	// list all
	d := swaggerRest.OrganizationApiV1OrganizationGetOpts{
		Username: optional.NewString("test1"),
		Roles:    optional.NewInterface("admin"),
	}
	data, r, err := ApiRest.OrganizationApi.V1OrganizationGet(AuthRest, &d)

	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), data.Data)
}

// TestOrgInvalidRolesSearch used for testing
// 用户可以批量查询某个用户所拥有权限的组织，参数不合法验证失败
func (s *SuiteOrgModel) TestOrgInvalidRolesSearch() {
	// list all
	d := swaggerRest.OrganizationApiV1OrganizationGetOpts{
		Username: optional.NewString("test1"),
		Roles:    optional.NewInterface("invalid"),
	}
	_, r, err := ApiRest.OrganizationApi.V1OrganizationGet(AuthRest, &d)

	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err, "invalid role value")
}

// TestOrgEmptyRolesSearch used for testing
// 用户可以批量查询某个用户所拥有权限的组织，查询结果为空
func (s *SuiteOrgModel) TestOrgEmptyRolesSearch() {
	// list all
	d := swaggerRest.OrganizationApiV1OrganizationGetOpts{
		Username: optional.NewString("test1"),
		Roles:    optional.NewInterface("read"),
	}
	data, r, err := ApiRest.OrganizationApi.V1OrganizationGet(AuthRest, &d)

	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)
	assert.Empty(s.T(), data.Data)
}

// TestOrgCreateFailedInvalidNameChars used for testing
// 无效的组织名, 名字过长, 为空
func (s *SuiteOrg) TestOrgCreateFailedInvalidNameChars() {
	// Slice of invalid names
	invalidNames := []string{
		"", // Empty name
		strings.Repeat(charA, ComConfig.ACCOUNT_NAME_MAX_LEN+1),
	}

	for _, name := range invalidNames {
		d := swaggerRest.ControllerOrgCreateRequest{
			Name:     name,
			Fullname: s.fullname,
			Website:  s.website,
		}

		_, r, err := ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
		// Expect a 400 Bad Request response due to invalid name
		assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode,
			"Expected a 400 Bad Request response for invalid name: "+name)
		assert.NotNil(s.T(), err, "Expected an error due to invalid name: "+name)
	}
}

// TestOrgCreateFailedInvalidFullNameChars used for testing
// 无效的组织昵称, 名字过长, 为空
func (s *SuiteOrg) TestOrgCreateFailedInvalidFullNameChars() {
	// Slice of invalid names
	invalidFullname := []string{
		"", // Empty name
		strings.Repeat(charA, ComConfig.MSD_FULLNAME_MAX_LEN+1),
	}

	for _, fullname := range invalidFullname {
		d := swaggerRest.ControllerOrgCreateRequest{
			Name:     s.name,
			Fullname: fullname,
			Website:  s.website,
		}

		_, r, err := ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
		// Expect a 400 Bad Request response due to invalid name
		assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode,
			"Expected a 400 Bad Request response for invalid name: "+fullname)
		assert.NotNil(s.T(), err, "Expected an error due to invalid name: "+fullname)
	}
}

// TestListMemberSucess used for testing
// 组织成员列表
func (s *SuiteOrg) TestListMemberSucess() {
	// 创建组织
	d := swaggerRest.ControllerOrgCreateRequest{
		Name:     s.name,
		Fullname: s.fullname,
	}

	_, r, err := ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	// 用户名搜索存在成员列表不为空
	orgData, r, err := ApiRest.OrganizationApi.V1OrganizationNameMemberGet(AuthRest, s.name,
		&swaggerRest.OrganizationApiV1OrganizationNameMemberGetOpts{Username: optional.NewString(s.owner)},
	)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	count := 0
	orgs := getArrary(s.T(), orgData.Data)
	for _, v := range orgs {
		if v != nil {
			assert.Equal(s.T(), s.owner, v["user_name"])
			count++
		}
	}
	assert.Equal(s.T(), countOne, count)

	// 用户名搜索不存在成员列表为空
	orgData, r, err = ApiRest.OrganizationApi.V1OrganizationNameMemberGet(AuthRest, s.name,
		&swaggerRest.OrganizationApiV1OrganizationNameMemberGetOpts{Username: optional.NewString("test2")},
	)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	orgs = getArrary(s.T(), orgData.Data)
	assert.Equal(s.T(), 0, len(orgs))

	// 角色搜索存在成员列表不为空
	orgData, r, err = ApiRest.OrganizationApi.V1OrganizationNameMemberGet(AuthRest, s.name,
		&swaggerRest.OrganizationApiV1OrganizationNameMemberGetOpts{Role: optional.NewString(s.defaultRole)},
	)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	count = 0
	orgs = getArrary(s.T(), orgData.Data)
	for _, v := range orgs {
		if v != nil {
			assert.Equal(s.T(), s.owner, v["user_name"])
			count++
		}
	}
	assert.Equal(s.T(), countOne, count)

	// 角色搜索不存在成员列表为空
	orgData, r, err = ApiRest.OrganizationApi.V1OrganizationNameMemberGet(AuthRest, s.name,
		&swaggerRest.OrganizationApiV1OrganizationNameMemberGetOpts{Role: optional.NewString("write")},
	)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	orgs = getArrary(s.T(), orgData.Data)
	assert.Equal(s.T(), 0, len(orgs))

	// 角色搜索不合法角色，返回400
	orgData, r, err = ApiRest.OrganizationApi.V1OrganizationNameMemberGet(AuthRest, s.name,
		&swaggerRest.OrganizationApiV1OrganizationNameMemberGetOpts{Role: optional.NewString("invalid")},
	)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	// 删除组织
	r, err = ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest, s.name)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestOrgLeaveFailed used for testing
// 离开组织失败：未加入组织
func (s *SuiteOrg) TestOrgLeaveFailed() {
	r, err := ApiRest.OrganizationApi.V1OrganizationNamePost(AuthRest, s.name)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestOrgLeaveSuccess used for testing
// 离开组织成功
func (s *SuiteOrg) TestOrgLeaveSuccess() {
	// 创建组织
	d := swaggerRest.ControllerOrgCreateRequest{
		Name:     s.name,
		Fullname: s.fullname,
	}
	data, r, err := ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	// 邀请进入组织
	_, r, err = ApiRest.OrganizationApi.V1InvitePost(AuthRest, swaggerRest.ControllerOrgInviteMemberRequest{
		OrgName: s.name,
		User:    s.invitee,
		Role:    "write",
		Msg:     "invite me",
	})
	assert.Equalf(s.T(), http.StatusCreated, r.StatusCode, data.Msg)
	assert.Nil(s.T(), err)

	// 被邀请人接受邀请
	_, r, err = ApiRest.OrganizationApi.V1InvitePut(AuthRest2, swaggerRest.ControllerOrgAcceptMemberRequest{
		OrgName: s.name,
		Msg:     "ok",
	})
	assert.Equalf(s.T(), http.StatusAccepted, r.StatusCode, data.Msg)
	assert.Nil(s.T(), err)

	// 离开组织成功
	r, err = ApiRest.OrganizationApi.V1OrganizationNamePost(AuthRest2, s.name)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	// 删除组织
	r, err = ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest, s.name)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestPrivilegeOrg used for testing
// 测试特权组织list
func (s *SuiteOrg) TestPrivilegeOrg() {
	// 创建组织
	d := swaggerRest.ControllerOrgCreateRequest{
		Name:     s.name,
		Fullname: s.fullname,
	}
	data, r, err := ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	// 未配置特权组织，无法查询到用户属于某个特权组织
	orgData, r, err := ApiRest.OrganizationApi.V1UserPrivilegeGet(AuthRest, "npu",
		&swaggerRest.OrganizationApiV1UserPrivilegeGetOpts{})
	assert.Equalf(s.T(), http.StatusOK, r.StatusCode, data.Msg)
	assert.Nil(s.T(), err)
	orgs := orgData.Data
	assert.Equal(s.T(), 0, len(orgs))

	// 未配置特权组织，无法查询到用户属于某个特权组织
	orgData, r, err = ApiRest.OrganizationApi.V1UserPrivilegeGet(AuthRest, "disable",
		&swaggerRest.OrganizationApiV1UserPrivilegeGetOpts{})
	assert.Equalf(s.T(), http.StatusOK, r.StatusCode, data.Msg)
	assert.Nil(s.T(), err)
	orgs = orgData.Data
	assert.Equal(s.T(), 0, len(orgs))

	// 无效参数
	orgData, r, err = ApiRest.OrganizationApi.V1UserPrivilegeGet(AuthRest, "test",
		&swaggerRest.OrganizationApiV1UserPrivilegeGetOpts{})
	assert.Equalf(s.T(), http.StatusBadRequest, r.StatusCode, data.Msg)
	assert.NotNil(s.T(), err)

	// 删除组织
	r, err = ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest, s.name)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestCertificate used for test
// 组织认证
func (s *SuiteOrg) TestCertificate() {
	orgForCertificate := "org_for_certificate"

	// 创建组织
	d := swaggerRest.ControllerOrgCreateRequest{
		Name:     orgForCertificate,
		Fullname: s.fullname,
	}
	_, r, err := ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	d2 := swaggerRest.ControllerOrgCertificateRequest{
		CertificateOrgType:      "学校",
		CertificateOrgName:      "清华大学",
		UnifiedSocialCreditCode: "914403001922038216",
		Phone:                   "15982129999",
		Identity:                "法定代表人",
		ImageOfCertificate:      "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",
	}

	_, r, err = ApiRest.OrganizationApi.V1OrganizationNameCertificatePost(AuthRest, orgForCertificate, d2)
	assert.Equal(s.T(), http.StatusInternalServerError, r.StatusCode)

	data, r, err := ApiRest.OrganizationApi.V1OrganizationNameCertificateGet(AuthRest, orgForCertificate)
	assert.Equalf(s.T(), http.StatusOK, r.StatusCode, data.Msg)
	assert.Nil(s.T(), err)
	orgCert := getData(s.T(), data.Data)
	assert.Equal(s.T(), "清华大学", orgCert["certificate_org_name"])

	data1, r, err := ApiRest.OrganizationApi.V1OrganizationNameCertificateCheckGet(AuthRest, orgForCertificate,
		&swaggerRest.OrganizationApiV1OrganizationNameCertificateCheckGetOpts{
			CertificateOrgName: optional.NewString("清华大学"),
		},
	)
	assert.Equalf(s.T(), http.StatusOK, r.StatusCode, data.Msg)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), true, data1.Data)

	// 删除组织
	r, err = ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest, orgForCertificate)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

func (s *SuiteOrg) TestMaxInvite() {
	// 测试环境配置了最多可以邀请4个
	d := swaggerRest.ControllerOrgCreateRequest{
		Name:     "testinviteorg1",
		Fullname: s.fullname,
	}

	_, r, err := ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	_, r, err = ApiRest.OrganizationApi.V1InvitePost(AuthRest, swaggerRest.ControllerOrgInviteMemberRequest{
		OrgName: "testinviteorg1",
		User:    s.invitee,
		Role:    "write",
		Msg:     "invite me",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	d = swaggerRest.ControllerOrgCreateRequest{
		Name:     "testinviteorg2",
		Fullname: s.fullname,
	}

	_, r, err = ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	_, r, err = ApiRest.OrganizationApi.V1InvitePost(AuthRest, swaggerRest.ControllerOrgInviteMemberRequest{
		OrgName: "testinviteorg2",
		User:    s.invitee,
		Role:    "write",
		Msg:     "invite me",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	d = swaggerRest.ControllerOrgCreateRequest{
		Name:     "testinviteorg3",
		Fullname: s.fullname,
	}

	_, r, err = ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	_, r, err = ApiRest.OrganizationApi.V1InvitePost(AuthRest, swaggerRest.ControllerOrgInviteMemberRequest{
		OrgName: "testinviteorg3",
		User:    s.invitee,
		Role:    "write",
		Msg:     "invite me",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	d = swaggerRest.ControllerOrgCreateRequest{
		Name:     "testinviteorg4",
		Fullname: s.fullname,
	}

	_, r, err = ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	_, r, err = ApiRest.OrganizationApi.V1InvitePost(AuthRest, swaggerRest.ControllerOrgInviteMemberRequest{
		OrgName: "testinviteorg4",
		User:    s.invitee,
		Role:    "write",
		Msg:     "invite me",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	d = swaggerRest.ControllerOrgCreateRequest{
		Name:     "testinviteorg5",
		Fullname: s.fullname,
	}

	_, r, err = ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	_, r, err = ApiRest.OrganizationApi.V1InvitePost(AuthRest, swaggerRest.ControllerOrgInviteMemberRequest{
		OrgName: "testinviteorg5",
		User:    s.invitee,
		Role:    "write",
		Msg:     "invite me",
	})

	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	// 删除邀请
	r, err = ApiRest.OrganizationApi.V1InviteDelete(AuthRest, swaggerRest.ControllerOrgRevokeInviteRequest{
		OrgName: "testinviteorg1",
		User:    s.invitee,
		Msg:     "invite me",
	})

	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.OrganizationApi.V1InviteDelete(AuthRest, swaggerRest.ControllerOrgRevokeInviteRequest{
		OrgName: "testinviteorg4",
		User:    s.invitee,
		Msg:     "invite me",
	})

	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.OrganizationApi.V1InviteDelete(AuthRest, swaggerRest.ControllerOrgRevokeInviteRequest{
		OrgName: "testinviteorg3",
		User:    s.invitee,
		Msg:     "invite me",
	})

	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.OrganizationApi.V1InviteDelete(AuthRest, swaggerRest.ControllerOrgRevokeInviteRequest{
		OrgName: "testinviteorg2",
		User:    s.invitee,
		Msg:     "invite me",
	})

	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	// 删除组织
	r, err = ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest, "testinviteorg1")
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest, "testinviteorg2")
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest, "testinviteorg3")
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest, "testinviteorg4")
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest, "testinviteorg5")
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

func (s *SuiteOrg) TestMaxRequest() {
	// 测试环境配置了最多可以邀请4个
	d := swaggerRest.ControllerOrgCreateRequest{
		Name:     "testinviteorg1",
		Fullname: s.fullname,
	}

	_, r, err := ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	_, r, err = ApiRest.OrganizationApi.V1InvitePost(AuthRest, swaggerRest.ControllerOrgInviteMemberRequest{
		OrgName: "testinviteorg1",
		User:    s.invitee,
		Role:    "write",
		Msg:     "invite me",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	d = swaggerRest.ControllerOrgCreateRequest{
		Name:     "testinviteorg2",
		Fullname: s.fullname,
	}

	_, r, err = ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	_, r, err = ApiRest.OrganizationApi.V1InvitePost(AuthRest, swaggerRest.ControllerOrgInviteMemberRequest{
		OrgName: "testinviteorg2",
		User:    s.invitee,
		Role:    "write",
		Msg:     "invite me",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	d = swaggerRest.ControllerOrgCreateRequest{
		Name:     "testinviteorg3",
		Fullname: s.fullname,
	}

	_, r, err = ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	_, r, err = ApiRest.OrganizationApi.V1InvitePost(AuthRest, swaggerRest.ControllerOrgInviteMemberRequest{
		OrgName: "testinviteorg3",
		User:    s.invitee,
		Role:    "write",
		Msg:     "invite me",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	d = swaggerRest.ControllerOrgCreateRequest{
		Name:     "testinviteorg4",
		Fullname: s.fullname,
	}

	_, r, err = ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	_, r, err = ApiRest.OrganizationApi.V1InvitePost(AuthRest, swaggerRest.ControllerOrgInviteMemberRequest{
		OrgName: "testinviteorg4",
		User:    s.invitee,
		Role:    "write",
		Msg:     "invite me",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	d = swaggerRest.ControllerOrgCreateRequest{
		Name:     "testinviteorg5",
		Fullname: s.fullname,
	}

	_, r, err = ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	_, r, err = ApiRest.OrganizationApi.V1InvitePost(AuthRest, swaggerRest.ControllerOrgInviteMemberRequest{
		OrgName: "testinviteorg5",
		User:    s.invitee,
		Role:    "write",
		Msg:     "invite me",
	})

	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	// 删除邀请
	r, err = ApiRest.OrganizationApi.V1InviteDelete(AuthRest, swaggerRest.ControllerOrgRevokeInviteRequest{
		OrgName: "testinviteorg1",
		User:    s.invitee,
		Msg:     "invite me",
	})

	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.OrganizationApi.V1InviteDelete(AuthRest, swaggerRest.ControllerOrgRevokeInviteRequest{
		OrgName: "testinviteorg4",
		User:    s.invitee,
		Msg:     "invite me",
	})

	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.OrganizationApi.V1InviteDelete(AuthRest, swaggerRest.ControllerOrgRevokeInviteRequest{
		OrgName: "testinviteorg3",
		User:    s.invitee,
		Msg:     "invite me",
	})

	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.OrganizationApi.V1InviteDelete(AuthRest, swaggerRest.ControllerOrgRevokeInviteRequest{
		OrgName: "testinviteorg2",
		User:    s.invitee,
		Msg:     "invite me",
	})

	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	// 删除组织
	r, err = ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest, "testinviteorg1")
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest, "testinviteorg2")
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest, "testinviteorg3")
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest, "testinviteorg4")
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest, "testinviteorg5")
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestOrg used for testing
func TestOrg(t *testing.T) {
	suite.Run(t, new(SuiteOrg))
}
