/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package e2e

import (
	"net/http"
	"testing"

	"github.com/antihax/optional"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	swagger "e2e/client"
)

// SuiteInvite used for testing
type SuiteInvite struct {
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
	invitee      string
	inviteeId    string
}

// SetupSuite used for testing
func (s *SuiteInvite) SetupSuite() {
	s.name = "testorg"
	s.fullname = "testorgfull"
	s.avatarid = "https://avatars.githubusercontent.com/u/2853724?v=1"
	s.allowRequest = true
	s.defaultRole = "admin"
	s.website = "https://www.modelfoundry.cn"
	s.desc = "test org desc"
	s.owner = "test1"   // this name is hard code in init-env.sh
	s.invitee = "test2" // this name is hard code in init-env.sh

	data, r, err := Api.UserApi.V1UserGet(Auth)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	user := getData(s.T(), data.Data)

	assert.NotEqual(s.T(), "", user["id"])
	s.owerId = getString(s.T(), user["id"])

	data, r, err = Api.UserApi.V1UserGet(Auth2)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	user = getData(s.T(), data.Data)

	assert.NotEqual(s.T(), "", user["id"])
	s.inviteeId = getString(s.T(), user["id"])

	data, r, err = Api.OrganizationApi.V1OrganizationPost(Auth, swagger.ControllerOrgCreateRequest{
		Name:        s.name,
		Fullname:    s.fullname,
		AvatarId:    s.avatarid,
		Website:     s.website,
		Description: s.desc,
	})

	o := getData(s.T(), data.Data)

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), "", o["id"])
	s.orgId = getString(s.T(), o["id"])
}

// TearDownSuite used for testing
func (s *SuiteInvite) TearDownSuite() {
	data, r, err := Api.OrganizationApi.V1InviteGet(Auth, &swagger.OrganizationApiV1InviteGetOpts{
		OrgName: optional.NewString(s.name),
		Status:  optional.NewString("pending"),
	})
	assert.Equalf(s.T(), http.StatusOK, r.StatusCode, data.Msg)
	assert.Nil(s.T(), err)

	invites := getArrary(s.T(), data.Data)

	count := 0
	for _, invite := range invites {
		if invite != nil {
			count++
		}
	}

	assert.Equal(s.T(), 0, count)

	r, err = Api.OrganizationApi.V1OrganizationNameDelete(Auth, s.name)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	data, r, err = Api.OrganizationApi.V1InviteGet(Auth, &swagger.OrganizationApiV1InviteGetOpts{
		OrgName: optional.NewString(s.name),
	})
	assert.Equal(s.T(), http.StatusNotFound, r.StatusCode)
	assert.NotNil(s.T(), err)

}

// TestInviteSuccess used for testing
// 创建邀请成功
func (s *SuiteInvite) TestInviteSuccess() {
	data, r, err := Api.OrganizationApi.V1InvitePost(Auth, swagger.ControllerOrgInviteMemberRequest{
		OrgName: s.name,
		User:    s.invitee,
		Role:    "write",
		Msg:     "invite me",
	})

	assert.Equalf(s.T(), http.StatusCreated, r.StatusCode, data.Msg)
	assert.Nil(s.T(), err)

	invite := getData(s.T(), data.Data)

	assert.Equal(s.T(), s.name, invite["org_name"])
	assert.Equal(s.T(), s.orgId, invite["org_id"])
	assert.Equal(s.T(), s.invitee, invite["user_name"])
	assert.Equal(s.T(), s.inviteeId, invite["user_id"])
	assert.NotEqual(s.T(), "", invite["id"])
	assert.NotEqual(s.T(), 0, getInt64(s.T(), invite["created_at"]))
	assert.NotEqual(s.T(), 0, getInt64(s.T(), invite["updated_at"]))
	assert.Equal(s.T(), "write", invite["role"])
	assert.Equal(s.T(), "invite me", invite["msg"])
	assert.Equal(s.T(), s.owner, invite["inviter"])

	r, err = Api.OrganizationApi.V1InviteDelete(Auth, swagger.ControllerOrgRevokeInviteRequest{
		OrgName: s.name,
		User:    s.invitee,
		Msg:     "no way",
	})

	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestInviteAprove used for testing
// 接受邀请
func (s *SuiteInvite) TestInviteAprove() {
	data, r, err := Api.OrganizationApi.V1InvitePost(Auth, swagger.ControllerOrgInviteMemberRequest{
		OrgName: s.name,
		User:    s.invitee,
		Role:    "write",
		Msg:     "invite me",
	})

	assert.Equalf(s.T(), http.StatusCreated, r.StatusCode, data.Msg)
	assert.Nil(s.T(), err)

	id := getData(s.T(), data.Data)["id"]

	// 重复邀请同一个用户, 该用户只收到一条邀请通知，并且以最新通知为准
	DupData, r2, err2 := Api.OrganizationApi.V1InvitePost(Auth, swagger.ControllerOrgInviteMemberRequest{
		OrgName: s.name,
		User:    s.invitee,
		Role:    "read",
		Msg:     "invite me ASAP",
	})

	assert.Equalf(s.T(), http.StatusCreated, r2.StatusCode, data.Msg)
	assert.Nil(s.T(), err2)

	DupInvite := getData(s.T(), DupData.Data)

	assert.NotEqual(s.T(), 0, getInt64(s.T(), DupInvite["created_at"]))
	assert.NotEqual(s.T(), 0, getInt64(s.T(), DupInvite["updated_at"]))
	assert.Equal(s.T(), "read", DupInvite["role"])
	assert.Equal(s.T(), "invite me ASAP", DupInvite["msg"])

	// 恢复初始角色
	_, _, _ = Api.OrganizationApi.V1InvitePost(Auth, swagger.ControllerOrgInviteMemberRequest{
		OrgName: s.name,
		User:    s.invitee,
		Role:    "write",
		Msg:     "invite me",
	})

	data2, r3, err3 := Api.OrganizationApi.V1InviteGet(Auth2, &swagger.OrganizationApiV1InviteGetOpts{
		Invitee: optional.NewString(s.invitee),
	})

	assert.Equalf(s.T(), http.StatusOK, r3.StatusCode, data.Msg)
	assert.Nil(s.T(), err3)

	notifications := getArrary(s.T(), data2.Data)

	countNtf := 0
	for _, notification := range notifications {
		if notification != nil {
			countNtf++
		}
	}

	assert.Equal(s.T(), countOne, countNtf)

	// 被邀请人接受邀请
	data, r, err = Api.OrganizationApi.V1InvitePut(Auth2, swagger.ControllerOrgAcceptMemberRequest{
		OrgName: s.name,
		Msg:     "ok",
	})

	assert.Equalf(s.T(), http.StatusAccepted, r.StatusCode, data.Msg)
	assert.Nil(s.T(), err)

	invite := getData(s.T(), data.Data)

	assert.Equal(s.T(), s.name, invite["org_name"])
	assert.Equal(s.T(), s.orgId, invite["org_id"])
	assert.Equal(s.T(), s.invitee, invite["user_name"])
	assert.Equal(s.T(), s.inviteeId, invite["user_id"])
	assert.NotEqual(s.T(), "", invite["id"])
	assert.NotEqual(s.T(), 0, getInt64(s.T(), invite["created_at"]))
	assert.NotEqual(s.T(), 0, getInt64(s.T(), invite["updated_at"]))
	assert.Equal(s.T(), "write", invite["role"])
	assert.Equal(s.T(), s.owner, invite["inviter"])
	assert.Equal(s.T(), s.invitee, invite["by"])
	assert.Equal(s.T(), "ok", invite["msg"])
	assert.Equal(s.T(), "approved", invite["status"])

	// 接收后成为member
	data, r, err = Api.OrganizationApi.V1OrganizationNameMemberGet(Auth2, s.name)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	// 已经在组织的用户无法邀请
	_, r4, err4 := Api.OrganizationApi.V1InvitePost(Auth, swagger.ControllerOrgInviteMemberRequest{
		OrgName: s.name,
		User:    s.invitee,
		Role:    "write",
		Msg:     "invite me",
	})

	assert.Equalf(s.T(), http.StatusBadRequest, r4.StatusCode, "the user is already a member of the org")
	assert.NotNil(s.T(), err4)

	members := getArrary(s.T(), data.Data)
	count := 0
	for _, member := range members {
		if member != nil && member["user_name"] == s.invitee {
			assert.Equal(s.T(), s.inviteeId, member["user_id"])
			assert.Equal(s.T(), s.name, member["org_name"])
			assert.Equal(s.T(), s.orgId, member["org_id"])
			assert.Equal(s.T(), "write", member["role"])
			assert.NotEqual(s.T(), 0, getInt64(s.T(), member["created_at"]))
			assert.NotEqual(s.T(), 0, getInt64(s.T(), member["updated_at"]))
			count += 2
		}
	}

	assert.Equal(s.T(), countTwo, count)

	// 查询已经接受的邀请
	data, r, err = Api.OrganizationApi.V1InviteGet(Auth, &swagger.OrganizationApiV1InviteGetOpts{
		OrgName: optional.NewString(s.name),
		Status:  optional.NewString("approved"),
	})
	assert.Equalf(s.T(), http.StatusOK, r.StatusCode, data.Msg)
	assert.Nil(s.T(), err)

	invites := getArrary(s.T(), data.Data)

	for _, invite := range invites {
		if invite != nil && invite["id"] == id {
			assert.Equal(s.T(), s.name, invite["org_name"])
			assert.Equal(s.T(), s.orgId, invite["org_id"])
			assert.Equal(s.T(), s.invitee, invite["user_name"])
			assert.Equal(s.T(), s.inviteeId, invite["user_id"])
			assert.NotEqual(s.T(), "", invite["id"])
			assert.NotEqual(s.T(), 0, getInt64(s.T(), invite["created_at"]))
			assert.NotEqual(s.T(), 0, getInt64(s.T(), invite["updated_at"]))
			assert.Equal(s.T(), "write", invite["role"])
			assert.Equal(s.T(), s.owner, invite["inviter"])
			assert.Equal(s.T(), s.invitee, invite["by"])
			assert.Equal(s.T(), "ok", invite["msg"])
		}
	}

	// 唯一的owner不能离开组织
	r, err = Api.OrganizationApi.V1OrganizationNameMemberDelete(Auth, swagger.ControllerOrgMemberRemoveRequest{
		User: s.owner,
	}, s.name)

	assert.Equalf(s.T(), http.StatusBadRequest, r.StatusCode, data.Msg)
	assert.NotNil(s.T(), err)

	r, err = Api.OrganizationApi.V1OrganizationNameMemberDelete(Auth, swagger.ControllerOrgMemberRemoveRequest{
		User: s.invitee,
	}, s.name)

	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestInviteInvalidPerm used for testing
// 无效的权限
func (s *SuiteInvite) TestInviteInvalidPerm() {
	data, r, err := Api.OrganizationApi.V1InvitePost(Auth, swagger.ControllerOrgInviteMemberRequest{
		OrgName: s.name,
		User:    s.invitee,
		Role:    "writer",
		Msg:     "invite me",
	})

	assert.Equalf(s.T(), http.StatusBadRequest, r.StatusCode, data.Msg)
	assert.NotNil(s.T(), err)
}

// TestInviteInvalidOrgname used for testing
// 无效的名字
func (s *SuiteInvite) TestInviteInvalidOrgname() {
	data, r, err := Api.OrganizationApi.V1InvitePost(Auth, swagger.ControllerOrgInviteMemberRequest{
		OrgName: "",
		User:    s.invitee,
		Role:    "write",
		Msg:     "invite me",
	})

	assert.Equalf(s.T(), http.StatusBadRequest, r.StatusCode, data.Msg)
	assert.NotNil(s.T(), err)

	data, r, err = Api.OrganizationApi.V1InvitePost(Auth, swagger.ControllerOrgInviteMemberRequest{
		OrgName: "orgnonexisted",
		User:    s.invitee,
		Role:    "write",
		Msg:     "invite me",
	})

	assert.Equalf(s.T(), http.StatusNotFound, r.StatusCode, data.Msg)
	assert.NotNil(s.T(), err)
}

// TestInviteInvalidUser used for testing
// 无效的用户名
func (s *SuiteInvite) TestInviteInvalidUser() {
	data, r, err := Api.OrganizationApi.V1InvitePost(Auth, swagger.ControllerOrgInviteMemberRequest{
		OrgName: s.name,
		User:    "",
		Role:    "write",
		Msg:     "invite me",
	})

	assert.Equalf(s.T(), http.StatusBadRequest, r.StatusCode, data.Msg)
	assert.NotNil(s.T(), err)

	data, r, err = Api.OrganizationApi.V1InvitePost(Auth, swagger.ControllerOrgInviteMemberRequest{
		OrgName: s.name,
		User:    "usernonexisted",
		Role:    "write",
		Msg:     "invite me",
	})

	assert.Equalf(s.T(), http.StatusNotFound, r.StatusCode, data.Msg)
	assert.NotNil(s.T(), err)
}

// TestInviteAprove used for testing
// Owner 可以被移除组织
func (s *SuiteInvite) TestRemoveOwner() {
	name := "testorg2"
	data, r, err := Api.OrganizationApi.V1OrganizationPost(Auth, swagger.ControllerOrgCreateRequest{
		Name:        name,
		Fullname:    name,
		AvatarId:    s.avatarid,
		Website:     s.website,
		Description: s.desc,
	})

	o := getData(s.T(), data.Data)

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), "", o["id"])
	s.orgId = getString(s.T(), o["id"])

	// 邀请另一个admin
	DupData, r2, err2 := Api.OrganizationApi.V1InvitePost(Auth, swagger.ControllerOrgInviteMemberRequest{
		OrgName: name,
		User:    s.invitee,
		Role:    "admin",
		Msg:     "invite me ASAP",
	})

	assert.Equalf(s.T(), http.StatusCreated, r2.StatusCode, data.Msg)
	assert.Nil(s.T(), err2)

	DupInvite := getData(s.T(), DupData.Data)

	assert.NotEqual(s.T(), 0, getInt64(s.T(), DupInvite["created_at"]))
	assert.NotEqual(s.T(), 0, getInt64(s.T(), DupInvite["updated_at"]))
	assert.Equal(s.T(), "admin", DupInvite["role"])
	assert.Equal(s.T(), "invite me ASAP", DupInvite["msg"])

	// 被邀请人接受邀请
	data, r, err = Api.OrganizationApi.V1InvitePut(Auth2, swagger.ControllerOrgAcceptMemberRequest{
		OrgName: name,
		Msg:     "ok",
	})

	assert.Equalf(s.T(), http.StatusAccepted, r.StatusCode, data.Msg)
	assert.Nil(s.T(), err)

	invite := getData(s.T(), data.Data)

	assert.Equal(s.T(), name, invite["org_name"])
	assert.Equal(s.T(), s.orgId, invite["org_id"])
	assert.Equal(s.T(), s.invitee, invite["user_name"])
	assert.Equal(s.T(), s.inviteeId, invite["user_id"])
	assert.NotEqual(s.T(), "", invite["id"])
	assert.NotEqual(s.T(), 0, getInt64(s.T(), invite["created_at"]))
	assert.NotEqual(s.T(), 0, getInt64(s.T(), invite["updated_at"]))
	assert.Equal(s.T(), "admin", invite["role"])
	assert.Equal(s.T(), s.owner, invite["inviter"])
	assert.Equal(s.T(), s.invitee, invite["by"])
	assert.Equal(s.T(), "ok", invite["msg"])
	assert.Equal(s.T(), "approved", invite["status"])

	// 移除原本的owner
	r, err = Api.OrganizationApi.V1OrganizationNameMemberDelete(Auth2, swagger.ControllerOrgMemberRemoveRequest{
		User: s.owner,
	}, name)

	assert.Equalf(s.T(), http.StatusNoContent, r.StatusCode, data.Msg)
	assert.Nil(s.T(), err)

	// 新的管理员可以删除组织
	r, err = Api.OrganizationApi.V1OrganizationNameDelete(Auth2, name)

	assert.Equalf(s.T(), http.StatusNoContent, r.StatusCode, data.Msg)
	assert.Nil(s.T(), err)

	// 再次创建组织成功
	data, r, err = Api.OrganizationApi.V1OrganizationPost(Auth, swagger.ControllerOrgCreateRequest{
		Name:        name,
		Fullname:    name,
		AvatarId:    s.avatarid,
		Website:     s.website,
		Description: s.desc,
	})

	o = getData(s.T(), data.Data)

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), "", o["id"])
	s.orgId = getString(s.T(), o["id"])

	// 清理组织
	r, err = Api.OrganizationApi.V1OrganizationNameDelete(Auth, name)

	assert.Equalf(s.T(), http.StatusNoContent, r.StatusCode, data.Msg)
	assert.Nil(s.T(), err)
}

// TestInviteOrg used for testing
// 组织邀请只能邀请用户，不能邀请组织
func (s *SuiteInvite) TestInviteOrg() {
	org := "testorg2"
	data, r, err := Api.OrganizationApi.V1OrganizationPost(Auth, swagger.ControllerOrgCreateRequest{
		Name:        org,
		Fullname:    org,
		AvatarId:    s.avatarid,
		Website:     s.website,
		Description: s.desc,
	})

	o := getData(s.T(), data.Data)

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), "", o["id"])

	// 邀请另一个组织作为admin
	_, r, err = Api.OrganizationApi.V1InvitePost(Auth, swagger.ControllerOrgInviteMemberRequest{
		OrgName: s.name,
		User:    org,
		Role:    "admin",
		Msg:     "invite me ASAP",
	})

	assert.Equalf(s.T(), http.StatusNotFound, r.StatusCode, data.Msg)
	assert.NotNil(s.T(), err)

	// 邀请另一个组织作为write
	_, r, err = Api.OrganizationApi.V1InvitePost(Auth, swagger.ControllerOrgInviteMemberRequest{
		OrgName: s.name,
		User:    org,
		Role:    "write",
		Msg:     "invite me ASAP",
	})

	assert.Equalf(s.T(), http.StatusNotFound, r.StatusCode, data.Msg)
	assert.NotNil(s.T(), err)

	// 邀请另一个read
	_, r, err = Api.OrganizationApi.V1InvitePost(Auth, swagger.ControllerOrgInviteMemberRequest{
		OrgName: s.name,
		User:    org,
		Role:    "read",
		Msg:     "invite me ASAP",
	})

	assert.Equalf(s.T(), http.StatusNotFound, r.StatusCode, data.Msg)
	assert.NotNil(s.T(), err)

	// 清理组织
	r, err = Api.OrganizationApi.V1OrganizationNameDelete(Auth, org)

	assert.Equalf(s.T(), http.StatusNoContent, r.StatusCode, data.Msg)
	assert.Nil(s.T(), err)
}

// TestInvite used for testing
func TestInvite(t *testing.T) {
	suite.Run(t, new(SuiteInvite))
}
