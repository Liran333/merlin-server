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

// SuiteOrgModel used for testing
type SuiteOrgModel struct {
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
func (s *SuiteOrgModel) SetupSuite() {
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
func (s *SuiteOrgModel) TearDownSuite() {
	r, err := ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest, s.name)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestDeleteOrgContainsModel used for testing
// 当组织下有model时，删除组织失败
func (s *SuiteOrgModel) TestDeleteOrgContainsModel() {
	data, r, err := ApiRest.ModelApi.V1ModelPost(AuthRest, swaggerRest.ControllerReqToCreateModel{
		Name:       "tempModel",
		Owner:      s.name,
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	// 删除组织失败
	r, err = ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest, s.name)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode,
		"can't delete the organization, while some repos still existed")
	assert.NotNil(s.T(), err)

	// 清空Model
	id := getString(s.T(), data.Data)

	r, err = ApiRest.ModelApi.V1ModelIdDelete(AuthRest, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestOrgReadMemberCantCreateUpdateDeleteModel used for testing
// 拥有read权限的用户不能创建模型，不能修改和删除他人模型
func (s *SuiteOrgModel) TestOrgReadMemberCantCreateUpdateDeleteModel() {
	_, r, err := ApiRest.OrganizationApi.V1OrganizationNameMemberPut(AuthRest, swaggerRest.ControllerOrgMemberEditRequest{
		User: "test2",
		Role: "read",
	}, s.name)

	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	_, r, err = ApiRest.ModelApi.V1ModelPost(AuthRest2, swaggerRest.ControllerReqToCreateModel{
		Name:       "testmodel",
		Owner:      s.name,
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusForbidden, r.StatusCode)
	assert.NotNil(s.T(), err)

	data, r, err := ApiRest.ModelApi.V1ModelPost(AuthRest, swaggerRest.ControllerReqToCreateModel{
		Name:       "testmodel",
		Owner:      s.name,
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	//read用户不能修改和删除他人模型
	_, r, err = ApiRest.ModelApi.V1ModelIdPut(AuthRest2, id, swaggerRest.ControllerReqToUpdateModel{
		Desc: "model desc new",
	})
	assert.Equal(s.T(), http.StatusForbidden, r.StatusCode)
	assert.NotNil(s.T(), err)

	r, err = ApiRest.ModelApi.V1ModelIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusForbidden, r.StatusCode)
	assert.NotNil(s.T(), err)

	_, r, err = ApiRest.OrganizationApi.V1OrganizationNameMemberPut(AuthRest, swaggerRest.ControllerOrgMemberEditRequest{
		User: "test2",
		Role: "write",
	}, s.name)

	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.ModelApi.V1ModelIdDelete(AuthRest, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestOrgWriteCreateDeleteModel used for testing
// 拥有write权限的用户可以创建和删除模型
func (s *SuiteOrgModel) TestOrgWriteCreateDeleteModel() {
	modelParam := swaggerRest.ControllerReqToCreateModel{
		Name:       "testmodel",
		Owner:      s.name,
		License:    "mit",
		Visibility: "public",
	}
	data, r, err := ApiRest.ModelApi.V1ModelPost(AuthRest2, modelParam)

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	// 重复创建模型返回400
	_, r, err = ApiRest.ModelApi.V1ModelPost(AuthRest2, modelParam)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	r, err = ApiRest.ModelApi.V1ModelIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestOrgWriteUpdateDeleteOthersModel used for testing
// 拥有write权限的用户可以修改和删除他人的模型
func (s *SuiteOrgModel) TestOrgWriteUpdateDeleteOthersModel() {
	data, r, err := ApiRest.ModelApi.V1ModelPost(AuthRest, swaggerRest.ControllerReqToCreateModel{
		Name:       "testmodel",
		Owner:      s.name,
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	//write用户可以修改和删除他人Space
	_, r, err = ApiRest.ModelApi.V1ModelIdPut(AuthRest2, id, swaggerRest.ControllerReqToUpdateModel{
		Desc: "model desc new",
	})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.ModelApi.V1ModelIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestOrgAdminUpdateDeleteOthersModel used for testing
// 拥有admin权限的用户可以修改和删除他人的模型
func (s *SuiteOrgModel) TestOrgAdminUpdateDeleteOthersModel() {
	data, r, err := ApiRest.ModelApi.V1ModelPost(AuthRest2, swaggerRest.ControllerReqToCreateModel{
		Name:       "testmodel",
		Owner:      s.name,
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	//admin用户可以修改和删除他人模型
	_, r, err = ApiRest.ModelApi.V1ModelIdPut(AuthRest, id, swaggerRest.ControllerReqToUpdateModel{
		Desc: "model desc new",
	})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.ModelApi.V1ModelIdDelete(AuthRest, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestOrgModel used for testing
func TestOrgModel(t *testing.T) {
	suite.Run(t, new(SuiteOrgModel))
}
