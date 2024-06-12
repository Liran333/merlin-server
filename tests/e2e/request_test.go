/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package e2e

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/antihax/optional"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	swaggerRest "e2e/client_rest"
)

// SuiteRequest used for testing
type SuiteRequest struct {
	suite.Suite
	name         string
	orgId        string
	fullname     string
	avatarid     string
	allowRequest bool
	defaultRole  string
	website      string
	desc         string
	owner        string
	owerId       string
	requesterId  string
	requester    string
}

// SetupSuite used for testing
func (s *SuiteRequest) SetupSuite() {
	s.name = "testorg"
	s.fullname = "testorgfull"
	s.avatarid = "https://avatars.githubusercontent.com/u/2853724?v=1"
	s.allowRequest = true
	s.defaultRole = "admin"
	s.website = "https://www.modelfoundry.cn"
	s.desc = "test org desc"
	s.owner = "test1"     // this name is hard code in init-env.sh
	s.requester = "test2" // this name is hard code in init-env.sh

	data, r, err := ApiRest.UserApi.V1UserGet(AuthRest)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	user := getData(s.T(), data.Data)

	assert.NotEqual(s.T(), "", user["id"])
	s.owerId = getString(s.T(), user["id"])

	data, r, err = ApiRest.UserApi.V1UserGet(AuthRest2)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	user = getData(s.T(), data.Data)

	assert.NotEqual(s.T(), "", user["id"])
	s.requesterId = getString(s.T(), user["id"])

	// 创建组织
	orgData, r, err := ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, swaggerRest.ControllerOrgCreateRequest{
		Name:        s.name,
		Fullname:    s.fullname,
		AvatarId:    s.avatarid,
		Website:     s.website,
		Description: s.desc,
	})
	o := getData(s.T(), orgData.Data)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), "", o["id"])
	s.orgId = getString(s.T(), o["id"])
	// 默认不可以主动申请
	assert.Equal(s.T(), false, o["allow_request"])
	// 修改申请权限
	_, r, err = ApiRest.OrganizationApi.V1OrganizationNamePut(AuthRest, s.name,
		swaggerRest.ControllerOrgBasicInfoUpdateRequest{
			AllowRequest: true,
		})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TearDownSuite used for testing
func (s *SuiteRequest) TearDownSuite() {
	// 加入组织申请列表为空
	data, r, err := ApiRest.OrganizationApi.V1RequestGet(AuthRest, &swaggerRest.OrganizationApiV1RequestGetOpts{
		OrgName:   optional.NewString(s.name),
		Requester: optional.NewString(s.requester),
		Status:    optional.NewString("pending"),
	})
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)
	apply := getData(s.T(), data.Data)
	assert.Empty(s.T(), apply["members"])

	// 删除组织
	r, err = ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest, s.name)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestRequestSuccess used for testing
// 创建加入组织请求成功，多次请求只保留一条记录
func (s *SuiteRequest) TestRequestSuccess() {
	// 创建加入组织申请
	data, r, err := ApiRest.OrganizationApi.V1RequestPost(AuthRest2, swaggerRest.ControllerOrgReqMemberRequest{
		OrgName: s.name,
		Msg:     "request me",
	})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	res := getData(s.T(), data.Data)
	firstID := res["id"]
	assert.Equal(s.T(), s.name, res["org_name"])
	assert.Equal(s.T(), s.orgId, res["org_id"])
	assert.Equal(s.T(), s.requester, res["username"])
	assert.Equal(s.T(), s.requesterId, res["user_id"])
	assert.Equal(s.T(), "write", res["role"])
	assert.Equal(s.T(), "request me", res["msg"])
	assert.NotEqual(s.T(), "", res["id"])
	assert.NotEqual(s.T(), 0, getInt64(s.T(), res["created_at"]))
	assert.NotEqual(s.T(), 0, getInt64(s.T(), res["updated_at"]))

	// 多次创建只保留一条记录
	data, r, err = ApiRest.OrganizationApi.V1RequestPost(AuthRest2, swaggerRest.ControllerOrgReqMemberRequest{
		OrgName: s.name,
		Msg:     "request second",
	})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.NotNil(s.T(), err)

	res = getData(s.T(), data.Data)
	assert.Equal(s.T(), firstID, res["id"])
	assert.Equal(s.T(), "request second", res["msg"])

	// 获取申请列表，只有一条记录
	data1, r, err := ApiRest.OrganizationApi.V1RequestGet(AuthRest, &swaggerRest.OrganizationApiV1RequestGetOpts{
		OrgName:   optional.NewString(s.name),
		Requester: optional.NewString(s.requester),
		Status:    optional.NewString("pending"),
	})
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)
	args := getData(s.T(), data1.Data)
	fmt.Println(args)
	// assert.Equal(s.T(), 1, len(args["Labels"].([]swaggerRest.GithubComOpenmerlinMerlinServerUserAppUserDto)))

	// 删除申请
	r, err = ApiRest.OrganizationApi.V1RequestDelete(AuthRest, swaggerRest.ControllerOrgRevokeMemberReqRequest{
		OrgName: s.name,
		User:    s.requester,
	})
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestInviteInvalidOrgname used for testing
// 无效的组织名字
func (s *SuiteRequest) TestRequestInvalidOrgname() {
	// 组织名为空
	_, r, err := ApiRest.OrganizationApi.V1RequestPost(AuthRest2, swaggerRest.ControllerOrgReqMemberRequest{
		OrgName: "",
		Msg:     "request me",
	})
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	// 组织不存在
	_, r, err = ApiRest.OrganizationApi.V1RequestPost(AuthRest2, swaggerRest.ControllerOrgReqMemberRequest{
		OrgName: "orgnonexisted",
		Msg:     "request me",
	})
	assert.Equal(s.T(), http.StatusNotFound, r.StatusCode)
	assert.NotNil(s.T(), err)

	// 已经是组织成员
	_, r, err = ApiRest.OrganizationApi.V1RequestPost(AuthRest, swaggerRest.ControllerOrgReqMemberRequest{
		OrgName: s.name,
		Msg:     "request me",
	})
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestApproveRequestSuccess used for testing
// 接受加入组织申请成功，其它邀请也要更新为接受状态
func (s *SuiteRequest) TestApproveRequestSuccess() {
	// 创建邀请
	_, r, err := ApiRest.OrganizationApi.V1InvitePost(AuthRest, swaggerRest.ControllerOrgInviteMemberRequest{
		OrgName: s.name,
		User:    s.requester,
		Role:    "write",
		Msg:     "invite me",
	})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	// 查询邀请列表不为空
	data1, r, err := ApiRest.OrganizationApi.V1InviteGet(AuthRest, &swaggerRest.OrganizationApiV1InviteGetOpts{
		OrgName: optional.NewString(s.name),
		Status:  optional.NewString("pending"),
	})
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), data1.Data)

	// 创建加入组织请求
	_, r, err = ApiRest.OrganizationApi.V1RequestPost(AuthRest2, swaggerRest.ControllerOrgReqMemberRequest{
		OrgName: s.name,
		Msg:     "request me",
	})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	// 查询加入组织申请列表不为空
	data2, r, err := ApiRest.OrganizationApi.V1RequestGet(AuthRest, &swaggerRest.OrganizationApiV1RequestGetOpts{
		OrgName: optional.NewString(s.name),
		Status:  optional.NewString("pending"),
	})
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), data2.Data)

	// 不能查询其他人的申请
	data2, r, err = ApiRest.OrganizationApi.V1RequestGet(AuthRest, &swaggerRest.OrganizationApiV1RequestGetOpts{
		Requester: optional.NewString(s.requester),
		Status:    optional.NewString("pending"),
	})
	assert.Equal(s.T(), http.StatusForbidden, r.StatusCode)
	assert.NotNil(s.T(), err)

	// 接受加入组织申请
	_, r, err = ApiRest.OrganizationApi.V1RequestPut(AuthRest, swaggerRest.ControllerOrgApproveMemberRequest{
		User:    s.requester,
		OrgName: s.name,
		Msg:     "approve me",
		Member:  "write",
	})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	// 查询邀请列表为空
	data1, r, err = ApiRest.OrganizationApi.V1InviteGet(AuthRest, &swaggerRest.OrganizationApiV1InviteGetOpts{
		OrgName: optional.NewString(s.name),
		Status:  optional.NewString("pending"),
	})
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)
	assert.Empty(s.T(), data1.Data)

	// 查询加入组织申请列表为空
	data2, r, err = ApiRest.OrganizationApi.V1RequestGet(AuthRest, &swaggerRest.OrganizationApiV1RequestGetOpts{
		OrgName: optional.NewString(s.name),
		Status:  optional.NewString("pending"),
	})
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)
	apply := getData(s.T(), data2)
	assert.Empty(s.T(), apply["members"])

	// 删除组织成员
	r, err = ApiRest.OrganizationApi.V1OrganizationNameMemberDelete(AuthRest,
		swaggerRest.ControllerOrgMemberRemoveRequest{
			User: s.requester,
		}, s.name)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 申请加入时必填角色
func (s *SuiteRequest) TestDefaultRole() {
	// 申请加入组织
	data, r, err := ApiRest.OrganizationApi.V1RequestPost(AuthRest2, swaggerRest.ControllerOrgReqMemberRequest{
		OrgName: s.name,
		Msg:     "request me",
	})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	res := getData(s.T(), data.Data)
	assert.Equal(s.T(), s.name, res["org_name"])
	assert.Equal(s.T(), s.orgId, res["org_id"])
	assert.Equal(s.T(), s.requester, res["username"])
	assert.Equal(s.T(), s.requesterId, res["user_id"])
	assert.Equal(s.T(), "write", res["role"])
	assert.Equal(s.T(), "request me", res["msg"])
	assert.NotEqual(s.T(), "", res["id"])
	assert.NotEqual(s.T(), 0, getInt64(s.T(), res["created_at"]))
	assert.NotEqual(s.T(), 0, getInt64(s.T(), res["updated_at"]))
	// 同意申请
	_, r, err = ApiRest.OrganizationApi.V1RequestPut(AuthRest, swaggerRest.ControllerOrgApproveMemberRequest{
		User:    s.requester,
		OrgName: s.name,
		Msg:     "approve me",
		Member:  "read",
	})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)
	// 获取成员列表,角色为read
	userData, r, err := ApiRest.OrganizationApi.V1OrganizationNameMemberGet(AuthRest, s.name, &swaggerRest.OrganizationApiV1OrganizationNameMemberGetOpts{})
	for _, i := range userData.Data {
		o := getData(s.T(), i)
		if o["user_name"] == "test2" {
			assert.Equal(s.T(), "read", o["role"])
		}
	}
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)
	// 删除组织成员
	r, err = ApiRest.OrganizationApi.V1OrganizationNameMemberDelete(AuthRest,
		swaggerRest.ControllerOrgMemberRemoveRequest{
			User: s.requester,
		}, s.name)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 拒绝用户加入
func (s *SuiteRequest) TestOrgReject() {
	// 申请加入组织
	data, r, err := ApiRest.OrganizationApi.V1RequestPost(AuthRest2, swaggerRest.ControllerOrgReqMemberRequest{
		OrgName: s.name,
		Msg:     "request me",
	})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	res := getData(s.T(), data.Data)
	assert.Equal(s.T(), s.name, res["org_name"])
	assert.Equal(s.T(), s.orgId, res["org_id"])
	assert.Equal(s.T(), s.requester, res["username"])
	assert.Equal(s.T(), s.requesterId, res["user_id"])
	assert.Equal(s.T(), "write", res["role"])
	assert.Equal(s.T(), "request me", res["msg"])
	assert.NotEqual(s.T(), "", res["id"])
	assert.NotEqual(s.T(), 0, getInt64(s.T(), res["created_at"]))
	assert.NotEqual(s.T(), 0, getInt64(s.T(), res["updated_at"]))
	// 拒绝申请
	r, err = ApiRest.OrganizationApi.V1RequestDelete(AuthRest, swaggerRest.ControllerOrgRevokeMemberReqRequest{
		User:    s.requester,
		OrgName: s.name,
		Msg:     "delete",
	})
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
	// 列出成员并判断角色
	userData, r, err := ApiRest.OrganizationApi.V1OrganizationNameMemberGet(AuthRest, s.name, &swaggerRest.OrganizationApiV1OrganizationNameMemberGetOpts{})
	for _, i := range userData.Data {
		o := getData(s.T(), i)
		if o["user_name"] == "test2" {
			assert.Equal(s.T(), "write", o["role"])
		}
	}
	assert.Equal(s.T(), 1, len(userData.Data))
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 用户主动撤销申请
func (s *SuiteRequest) TestCancelOrg() {
	// 申请加入组织
	data, r, err := ApiRest.OrganizationApi.V1RequestPost(AuthRest2, swaggerRest.ControllerOrgReqMemberRequest{
		OrgName: s.name,
		Msg:     "request me",
	})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	res := getData(s.T(), data.Data)
	assert.Equal(s.T(), s.name, res["org_name"])
	assert.Equal(s.T(), s.orgId, res["org_id"])
	assert.Equal(s.T(), s.requester, res["username"])
	assert.Equal(s.T(), s.requesterId, res["user_id"])
	assert.Equal(s.T(), "write", res["role"])
	assert.Equal(s.T(), "request me", res["msg"])
	assert.NotEqual(s.T(), "", res["id"])
	assert.NotEqual(s.T(), 0, getInt64(s.T(), res["created_at"]))
	assert.NotEqual(s.T(), 0, getInt64(s.T(), res["updated_at"]))
	//检查申请列表
	apply, r, err := ApiRest.OrganizationApi.V1RequestGet(AuthRest, &swaggerRest.OrganizationApiV1RequestGetOpts{
		OrgName:  optional.NewString(s.name),
		Status:   optional.NewString("pending"),
		Page:     optional.NewInt32(1),
		PageSize: optional.NewInt32(2),
	})
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)
	fmt.Println(apply.Data)
	// 用户撤销申请
	r, err = ApiRest.OrganizationApi.V1RequestDelete(AuthRest2, swaggerRest.ControllerOrgRevokeMemberReqRequest{
		User:    s.requester,
		OrgName: s.name,
		Msg:     "delete",
	})
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
	// 检查申请列表
	_, r, err = ApiRest.OrganizationApi.V1RequestGet(AuthRest, &swaggerRest.OrganizationApiV1RequestGetOpts{
		OrgName: optional.NewString(s.name),
		Status:  optional.NewString("pending"),
	})
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 撤销查询申请
func (s *SuiteRequest) TestSearch() {
	// 创建加入组织申请
	data, r, _ := ApiRest.OrganizationApi.V1RequestPost(AuthRest2, swaggerRest.ControllerOrgReqMemberRequest{
		OrgName: s.name,
		Msg:     "request me",
	})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)

	res := getData(s.T(), data.Data)
	assert.Equal(s.T(), s.name, res["org_name"])
	assert.Equal(s.T(), s.orgId, res["org_id"])
	assert.Equal(s.T(), s.requester, res["username"])
	assert.Equal(s.T(), s.requesterId, res["user_id"])
	assert.Equal(s.T(), "write", res["role"])
	assert.Equal(s.T(), "request me", res["msg"])
	assert.NotEqual(s.T(), "", res["id"])
	assert.NotEqual(s.T(), 0, getInt64(s.T(), res["created_at"]))
	assert.NotEqual(s.T(), 0, getInt64(s.T(), res["updated_at"]))
	// 查询申请列表
	data2, r, err := ApiRest.OrganizationApi.V1RequestGet(AuthRest, &swaggerRest.OrganizationApiV1RequestGetOpts{
		OrgName: optional.NewString(s.name),
		// Status:  optional.NewString("pending"),
		Requester: optional.NewString(s.requester),
	})
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), data2.Data)
	// 撤销申请
	r, err = ApiRest.OrganizationApi.V1RequestDelete(AuthRest2, swaggerRest.ControllerOrgRevokeMemberReqRequest{
		User:    s.requester,
		OrgName: s.name,
		Msg:     "delete",
	})
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 用户查询申请状态
func (s *SuiteRequest) TestOneSearch() {
	// 创建加入组织申请
	data, r, err := ApiRest.OrganizationApi.V1RequestPost(AuthRest2, swaggerRest.ControllerOrgReqMemberRequest{
		OrgName: s.name,
		Msg:     "request me",
	})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	res := getData(s.T(), data.Data)
	assert.Equal(s.T(), s.name, res["org_name"])
	assert.Equal(s.T(), s.orgId, res["org_id"])
	assert.Equal(s.T(), s.requester, res["username"])
	assert.Equal(s.T(), s.requesterId, res["user_id"])
	assert.Equal(s.T(), "write", res["role"])
	assert.Equal(s.T(), "request me", res["msg"])
	assert.NotEqual(s.T(), "", res["id"])
	assert.NotEqual(s.T(), 0, getInt64(s.T(), res["created_at"]))
	assert.NotEqual(s.T(), 0, getInt64(s.T(), res["updated_at"]))
	// 用户查询申请,只返回最新数据
	dataRes, r, err := ApiRest.OrganizationApi.V1RequestOnlyUsernameOrgnameGet(AuthRest2, s.requester, s.name)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(dataRes.Data))
	// 用户撤销申请
	r, err = ApiRest.OrganizationApi.V1RequestDelete(AuthRest2, swaggerRest.ControllerOrgRevokeMemberReqRequest{
		User:    s.requester,
		OrgName: s.name,
		Msg:     "delete",
	})
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 申请理由字数不得超过20字
func (s *SuiteRequest) TestStrLimit() {
	// 创建申请,接口返回错误
	_, r, err := ApiRest.OrganizationApi.V1RequestPost(AuthRest2, swaggerRest.ControllerOrgReqMemberRequest{
		OrgName: s.name,
		Msg:     "这是为了测试字符串长度限制的消息，限制为20个字符",
	})
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestRequest used for testing
func TestRequest(t *testing.T) {
	suite.Run(t, new(SuiteRequest))
}
