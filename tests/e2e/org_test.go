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
	assert.Nil(s.T(), org["email"])

	data, r, err = ApiRest.OrganizationApi.V1OrganizationGet(AuthRest,
		&swaggerRest.OrganizationApiV1OrganizationGetOpts{})
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	count := 0
	orgs := getArrary(s.T(), data.Data)
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
		Name:     s.owner,
		Fullname: s.fullname,
	}

	_, r, err := ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	r, err = ApiRest.OrganizationApi.V1NameHead(AuthRest, s.owner)
	assert.Equal(s.T(), http.StatusConflict, r.StatusCode)
	assert.NotNil(s.T(), err)

	r, err = ApiRest.OrganizationApi.V1NameHead(AuthRest, "testnonexist")
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
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

	data, r, err := ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	// 用户名搜索存在成员列表不为空
	data, r, err = ApiRest.OrganizationApi.V1OrganizationNameMemberGet(AuthRest, s.name,
		&swaggerRest.OrganizationApiV1OrganizationNameMemberGetOpts{Username: optional.NewString(s.owner)},
	)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	count := 0
	orgs := getArrary(s.T(), data.Data)
	for _, v := range orgs {
		if v != nil {
			assert.Equal(s.T(), s.owner, v["user_name"])
			count++
		}
	}
	assert.Equal(s.T(), countOne, count)

	// 用户名搜索不存在成员列表为空
	data, r, err = ApiRest.OrganizationApi.V1OrganizationNameMemberGet(AuthRest, s.name,
		&swaggerRest.OrganizationApiV1OrganizationNameMemberGetOpts{Username: optional.NewString("test2")},
	)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	orgs = getArrary(s.T(), data.Data)
	assert.Equal(s.T(), 0, len(orgs))

	// 角色搜索存在成员列表不为空
	data, r, err = ApiRest.OrganizationApi.V1OrganizationNameMemberGet(AuthRest, s.name,
		&swaggerRest.OrganizationApiV1OrganizationNameMemberGetOpts{Role: optional.NewString(s.defaultRole)},
	)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	count = 0
	orgs = getArrary(s.T(), data.Data)
	for _, v := range orgs {
		if v != nil {
			assert.Equal(s.T(), s.owner, v["user_name"])
			count++
		}
	}
	assert.Equal(s.T(), countOne, count)

	// 角色搜索不存在成员列表为空
	data, r, err = ApiRest.OrganizationApi.V1OrganizationNameMemberGet(AuthRest, s.name,
		&swaggerRest.OrganizationApiV1OrganizationNameMemberGetOpts{Role: optional.NewString("write")},
	)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	orgs = getArrary(s.T(), data.Data)
	assert.Equal(s.T(), 0, len(orgs))

	// 角色搜索不合法角色，返回400
	data, r, err = ApiRest.OrganizationApi.V1OrganizationNameMemberGet(AuthRest, s.name,
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
	data, r, err = ApiRest.OrganizationApi.V1InvitePost(AuthRest, swaggerRest.ControllerOrgInviteMemberRequest{
		OrgName: s.name,
		User:    s.invitee,
		Role:    "write",
		Msg:     "invite me",
	})
	assert.Equalf(s.T(), http.StatusCreated, r.StatusCode, data.Msg)
	assert.Nil(s.T(), err)

	// 被邀请人接受邀请
	data, r, err = ApiRest.OrganizationApi.V1InvitePut(AuthRest2, swaggerRest.ControllerOrgAcceptMemberRequest{
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

// TestOrg used for testing
func TestOrg(t *testing.T) {
	suite.Run(t, new(SuiteOrg))
}
