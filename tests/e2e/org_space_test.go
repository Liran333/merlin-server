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

	data, r, err := ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, swaggerRest.ControllerOrgCreateRequest{
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

	data, r, err = ApiRest.OrganizationApi.V1InvitePost(AuthRest, swaggerRest.ControllerOrgInviteMemberRequest{
		OrgName: s.name,
		User:    "test2",
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
}

// TearDownSuite used for testing
func (s *SuiteOrgSpace) TearDownSuite() {
	r, err := ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest, s.name)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestDeleteSpaceContainsModel used for testing
// 当组织下有Space，组织卡片时，删除组织失败
func (s *SuiteOrgModel) TestDeleteSpaceContainsModel() {
	spaceParam := swaggerRest.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "tempFullName",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		License:    "mit",
		Name:       "tempSpace",
		Owner:      s.name,
		Sdk:        "gradio",
		Visibility: "public",
	}
	data, r, err := ApiRest.SpaceApi.V1SpacePost(AuthRest, spaceParam)

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	// 重复创建空间返回400
	_, r, err = ApiRest.SpaceApi.V1SpacePost(AuthRest, spaceParam)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	// 删除组织失败
	r, err = ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest, s.name)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode,
		"can't delete the organization, while some spaces still existed")
	assert.NotNil(s.T(), err)

	// 清空Space
	id := getString(s.T(), data.Data)

	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestOrgReadMemberCantCreateUpdateDeleteSpace used for testing
// 拥有read权限的用户不能创建Space，不能修改和删除他人Space
func (s *SuiteOrgSpace) TestOrgReadMemberCantCreateUpdateDeleteSpace() {
	_, r, err := ApiRest.OrganizationApi.V1OrganizationNameMemberPut(AuthRest, swaggerRest.ControllerOrgMemberEditRequest{
		User: "test2",
		Role: "read",
	}, s.name)

	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	// read用户不能创建Space
	_, r, err = ApiRest.SpaceApi.V1SpacePost(AuthRest2, swaggerRest.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		License:    "mit",
		Name:       "testspace",
		Owner:      s.name,
		Sdk:        "gradio",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusForbidden, r.StatusCode)
	assert.NotNil(s.T(), err)

	data, r, err := ApiRest.SpaceApi.V1SpacePost(AuthRest, swaggerRest.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
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
	_, r, err = ApiRest.SpaceApi.V1SpaceIdPut(AuthRest2, id, swaggerRest.ControllerReqToUpdateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		Sdk:        "gradio",
		Visibility: "public",
	})
	assert.Equal(s.T(), http.StatusForbidden, r.StatusCode)
	assert.NotNil(s.T(), err)

	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusForbidden, r.StatusCode)
	assert.NotNil(s.T(), err)

	_, r, err = ApiRest.OrganizationApi.V1OrganizationNameMemberPut(AuthRest, swaggerRest.ControllerOrgMemberEditRequest{
		User: "test2",
		Role: "write",
	}, s.name)

	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestOrgWriteCreateDeleteSpace used for testing
// 拥有write权限的用户可以创建和删除Space
func (s *SuiteOrgSpace) TestOrgWriteCreateDeleteSpace() {
	data, r, err := ApiRest.SpaceApi.V1SpacePost(AuthRest2, swaggerRest.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		License:    "mit",
		Name:       "testspace",
		Owner:      s.name,
		Sdk:        "gradio",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestOrgWriteUpdateDeleteOthersSpace used for testing
// 拥有write权限的用户可以修改和删除他人的Space
func (s *SuiteOrgSpace) TestOrgWriteUpdateDeleteOthersSpace() {
	data, r, err := ApiRest.SpaceApi.V1SpacePost(AuthRest, swaggerRest.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
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
	_, r, err = ApiRest.SpaceApi.V1SpaceIdPut(AuthRest2, id, swaggerRest.ControllerReqToUpdateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		Sdk:        "gradio",
		Visibility: "public",
	})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestOrgAdminUpdateDeleteOthersSpace used for testing
// 拥有admin权限的用户可以修改和删除他人的Space
func (s *SuiteOrgSpace) TestOrgAdminUpdateDeleteOthersSpace() {
	data, r, err := ApiRest.SpaceApi.V1SpacePost(AuthRest2, swaggerRest.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
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
	_, r, err = ApiRest.SpaceApi.V1SpaceIdPut(AuthRest, id, swaggerRest.ControllerReqToUpdateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		Sdk:        "gradio",
		Visibility: "public",
	})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestOrgSpace used for testing
func TestOrgSpace(t *testing.T) {
	suite.Run(t, new(SuiteOrgSpace))
}
