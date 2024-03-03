/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package e2e

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	swagger "e2e/client"
)

// SuiteOrgSpace used for testing
type SuiteOrgSpace struct {
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
}

// SetupSuite used for testing
func (s *SuiteOrgSpace) SetupSuite() {
	s.name = "testorg"
	s.fullname = "testorgfull"
	s.avatarid = "https://avatars.githubusercontent.com/u/2853724?v=1"
	s.allowRequest = true
	s.defaultRole = "admin"
	s.website = "https://www.modelfoundry.cn"
	s.desc = "test org desc"
	s.owner = "test1" // this name is hard code in init-env.sh

	data, r, err := Api.OrganizationApi.V1OrganizationPost(Auth, swagger.ControllerOrgCreateRequest{
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

	data, r, err = Api.OrganizationApi.V1InvitePost(Auth, swagger.ControllerOrgInviteMemberRequest{
		OrgName: s.name,
		User:    "test2",
		Role:    "write",
		Msg:     "invite me",
	})

	assert.Equalf(s.T(), http.StatusCreated, r.StatusCode, data.Msg)
	assert.Nil(s.T(), err)

	// 被邀请人接受邀请
	data, r, err = Api.OrganizationApi.V1InvitePut(Auth2, swagger.ControllerOrgAcceptMemberRequest{
		OrgName: s.name,
		Msg:     "ok",
	})

	assert.Equalf(s.T(), http.StatusAccepted, r.StatusCode, data.Msg)
	assert.Nil(s.T(), err)
}

// TearDownSuite used for testing
func (s *SuiteOrgSpace) TearDownSuite() {
	r, err := Api.OrganizationApi.V1OrganizationNameDelete(Auth, s.name)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestDeleteSpaceContainsModel used for testing
// 当组织下有Space，组织卡片时，删除组织失败
func (s *SuiteOrgModel) TestDeleteSpaceContainsModel() {
	data, r, err := Api.SpaceApi.V1SpacePost(Auth, swagger.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "tempFullName",
		Hardware:   "CPU basic 2 vCPU · 4GB · FREE",
		InitReadme: false,
		License:    "mit",
		Name:       "tempSpace",
		Owner:      s.name,
		Sdk:        "gradio",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	// 删除组织失败
	r, err = Api.OrganizationApi.V1OrganizationNameDelete(Auth, s.name)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode,
		"can't delete the organization, while some spaces still existed")
	assert.NotNil(s.T(), err)

	// 清空Space
	id := getString(s.T(), data.Data)

	r, err = Api.SpaceApi.V1SpaceIdDelete(Auth, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestOrgReadMemberCantCreateUpdateDeleteSpace used for testing
// 拥有read权限的用户不能创建Space，不能修改和删除他人Space
func (s *SuiteOrgSpace) TestOrgReadMemberCantCreateUpdateDeleteSpace() {
	_, r, err := Api.OrganizationApi.V1OrganizationNameMemberPut(Auth, swagger.ControllerOrgMemberEditRequest{
		User: "test2",
		Role: "read",
	}, s.name)

	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	// read用户不能创建Space
	_, r, err = Api.SpaceApi.V1SpacePost(Auth2, swagger.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 4GB · FREE",
		InitReadme: false,
		License:    "mit",
		Name:       "testspace",
		Owner:      s.name,
		Sdk:        "gradio",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusForbidden, r.StatusCode)
	assert.NotNil(s.T(), err)

	data, r, err := Api.SpaceApi.V1SpacePost(Auth, swagger.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 4GB · FREE",
		InitReadme: false,
		License:    "mit",
		Name:       "testspace",
		Owner:      s.name,
		Sdk:        "gradio",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	//read用户不能修改和删除他人Space
	_, r, err = Api.SpaceApi.V1SpaceIdPut(Auth2, id, swagger.ControllerReqToUpdateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		Name:       "testspace",
		Sdk:        "gradio",
		Visibility: "public",
	})
	assert.Equal(s.T(), http.StatusForbidden, r.StatusCode)
	assert.NotNil(s.T(), err)

	r, err = Api.SpaceApi.V1SpaceIdDelete(Auth2, id)
	assert.Equal(s.T(), http.StatusForbidden, r.StatusCode)
	assert.NotNil(s.T(), err)

	_, r, err = Api.OrganizationApi.V1OrganizationNameMemberPut(Auth, swagger.ControllerOrgMemberEditRequest{
		User: "test2",
		Role: "write",
	}, s.name)

	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = Api.SpaceApi.V1SpaceIdDelete(Auth, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestOrgWriteCreateDeleteSpace used for testing
// 拥有write权限的用户可以创建和删除Space
func (s *SuiteOrgSpace) TestOrgWriteCreateDeleteSpace() {
	data, r, err := Api.SpaceApi.V1SpacePost(Auth2, swagger.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 4GB · FREE",
		InitReadme: false,
		License:    "mit",
		Name:       "testspace",
		Owner:      s.name,
		Sdk:        "gradio",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	r, err = Api.SpaceApi.V1SpaceIdDelete(Auth2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestOrgWriteUpdateDeleteOthersSpace used for testing
// 拥有write权限的用户可以修改和删除他人的Space
func (s *SuiteOrgSpace) TestOrgWriteUpdateDeleteOthersSpace() {
	data, r, err := Api.SpaceApi.V1SpacePost(Auth, swagger.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 4GB · FREE",
		InitReadme: false,
		License:    "mit",
		Name:       "testspace",
		Owner:      s.name,
		Sdk:        "gradio",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	//write用户可以修改和删除他人Space
	_, r, err = Api.SpaceApi.V1SpaceIdPut(Auth2, id, swagger.ControllerReqToUpdateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		Name:       "testspace",
		Sdk:        "gradio",
		Visibility: "public",
	})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = Api.SpaceApi.V1SpaceIdDelete(Auth2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestOrgAdminUpdateDeleteOthersSpace used for testing
// 拥有admin权限的用户可以修改和删除他人的Space
func (s *SuiteOrgSpace) TestOrgAdminUpdateDeleteOthersSpace() {
	data, r, err := Api.SpaceApi.V1SpacePost(Auth2, swagger.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 4GB · FREE",
		InitReadme: false,
		License:    "mit",
		Name:       "testspace",
		Owner:      s.name,
		Sdk:        "gradio",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	//admin用户可以修改和删除他人Space
	_, r, err = Api.SpaceApi.V1SpaceIdPut(Auth, id, swagger.ControllerReqToUpdateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		Name:       "testspace",
		Sdk:        "gradio",
		Visibility: "public",
	})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = Api.SpaceApi.V1SpaceIdDelete(Auth, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestOrgContributorCreateUpdateDelete used for testing
// 拥有contribute权限的用户可以创建修改和删除自己的Space
func (s *SuiteOrgSpace) TestOrgContributorCreateUpdateDelete() {
	_, r, err := Api.OrganizationApi.V1OrganizationNameMemberPut(Auth, swagger.ControllerOrgMemberEditRequest{
		User: "test2",
		Role: "contributor",
	}, s.name)
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	data, r, err := Api.SpaceApi.V1SpacePost(Auth2, swagger.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 4GB · FREE",
		InitReadme: false,
		License:    "mit",
		Name:       "testspace",
		Owner:      s.name,
		Sdk:        "gradio",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	_, r, err = Api.SpaceApi.V1SpaceIdPut(Auth2, id, swagger.ControllerReqToUpdateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		Name:       "testspace",
		Sdk:        "gradio",
		Visibility: "public",
	})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = Api.SpaceApi.V1SpaceIdDelete(Auth2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	_, r, err = Api.OrganizationApi.V1OrganizationNameMemberPut(Auth, swagger.ControllerOrgMemberEditRequest{
		User: "test2",
		Role: "write",
	}, s.name)
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestOrgContributorCantUpdateDeleteOthersModel used for testing
// 拥有contribute权限的用户不可以修改或删除他人Space
func (s *SuiteOrgSpace) TestOrgContributorCantUpdateDeleteOthersModel() {
	_, r, err := Api.OrganizationApi.V1OrganizationNameMemberPut(Auth, swagger.ControllerOrgMemberEditRequest{
		User: "test2",
		Role: "contributor",
	}, s.name)

	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	data, r, err := Api.SpaceApi.V1SpacePost(Auth, swagger.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 4GB · FREE",
		InitReadme: false,
		License:    "mit",
		Name:       "testspace",
		Owner:      s.name,
		Sdk:        "gradio",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	r, err = Api.SpaceApi.V1SpaceIdDelete(Auth2, id)
	assert.Equal(s.T(), http.StatusForbidden, r.StatusCode)
	assert.NotNil(s.T(), err)

	_, r, err = Api.SpaceApi.V1SpaceIdPut(Auth2, id, swagger.ControllerReqToUpdateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		Name:       "testspace",
		Sdk:        "gradio",
		Visibility: "public",
	})

	// 拥有contribute权限的用户不可以修改他人Space
	assert.Equal(s.T(), http.StatusForbidden, r.StatusCode)
	assert.NotNil(s.T(), err)

	// 拥有contribute权限的用户不可以删除他人Space
	r, err = Api.SpaceApi.V1SpaceIdDelete(Auth2, id)
	assert.Equal(s.T(), http.StatusForbidden, r.StatusCode)
	assert.NotNil(s.T(), err)

	r, err = Api.SpaceApi.V1SpaceIdDelete(Auth, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	_, r, err = Api.OrganizationApi.V1OrganizationNameMemberPut(Auth, swagger.ControllerOrgMemberEditRequest{
		User: "test2",
		Role: "write",
	}, s.name)

	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestOrgSpace used for testing
func TestOrgSpace(t *testing.T) {
	suite.Run(t, new(SuiteOrgSpace))
}
