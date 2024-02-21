package e2e

import (
	swagger "e2e/client"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

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

func (s *SuiteOrgModel) SetupSuite() {
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

	assert.Equal(s.T(), 201, r.StatusCode)
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), "", o["id"])
	s.orgId = getString(s.T(), o["id"])

	data, r, err = Api.OrganizationApi.V1InvitePost(Auth, swagger.ControllerOrgInviteMemberRequest{
		OrgName: s.name,
		User:    "test2",
		Role:    "write",
		Msg:     "invite me",
	})

	assert.Equalf(s.T(), 201, r.StatusCode, data.Msg)
	assert.Nil(s.T(), err)

	// 被邀请人接受邀请
	data, r, err = Api.OrganizationApi.V1InvitePut(Auth2, swagger.ControllerOrgAcceptMemberRequest{
		OrgName: s.name,
		Msg:     "ok",
	})

	assert.Equalf(s.T(), 202, r.StatusCode, data.Msg)
	assert.Nil(s.T(), err)
}

func (s *SuiteOrgModel) TearDownSuite() {
	r, err := Api.OrganizationApi.V1OrganizationNameDelete(Auth, s.name)
	assert.Equal(s.T(), 204, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 拥有read权限的用户不能创建模型，不能修改和删除他人模型
func (s *SuiteOrgModel) TestOrgReadMemberCantCreateUpdateDeleteModel() {
	_, r, err := Api.OrganizationApi.V1OrganizationNameMemberPut(Auth, swagger.ControllerOrgMemberEditRequest{
		User: "test2",
		Role: "read",
	}, s.name)

	assert.Equal(s.T(), 202, r.StatusCode)
	assert.Nil(s.T(), err)

	_, r, err = Api.ModelApi.V1ModelPost(Auth2, swagger.ControllerReqToCreateModel{
		Name:       "testmodel",
		Owner:      s.name,
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), 403, r.StatusCode)
	assert.NotNil(s.T(), err)

	data, r, err := Api.ModelApi.V1ModelPost(Auth, swagger.ControllerReqToCreateModel{
		Name:       "testmodel",
		Owner:      s.name,
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), 201, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	//read用户不能修改和删除他人模型
	_, r, err = Api.ModelApi.V1ModelIdPut(Auth2, id, swagger.ControllerReqToUpdateModel{
		Desc: "model desc new",
	})
	assert.Equal(s.T(), 404, r.StatusCode)
	assert.NotNil(s.T(), err)

	r, err = Api.ModelApi.V1ModelIdDelete(Auth2, id)
	assert.Equal(s.T(), 404, r.StatusCode)
	assert.NotNil(s.T(), err)

	_, r, err = Api.OrganizationApi.V1OrganizationNameMemberPut(Auth, swagger.ControllerOrgMemberEditRequest{
		User: "test2",
		Role: "write",
	}, s.name)

	assert.Equal(s.T(), 202, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = Api.ModelApi.V1ModelIdDelete(Auth, id)
	assert.Equal(s.T(), 204, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 拥有write权限的用户可以创建和删除模型
func (s *SuiteOrgModel) TestOrgWriteCreateDeleteModel() {
	data, r, err := Api.ModelApi.V1ModelPost(Auth2, swagger.ControllerReqToCreateModel{
		Name:       "testmodel",
		Owner:      s.name,
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), 201, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	r, err = Api.ModelApi.V1ModelIdDelete(Auth2, id)
	assert.Equal(s.T(), 204, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 拥有write权限的用户可以修改和删除他人的模型
func (s *SuiteOrgModel) TestOrgWriteUpdateDeleteOthersModel() {
	data, r, err := Api.ModelApi.V1ModelPost(Auth, swagger.ControllerReqToCreateModel{
		Name:       "testmodel",
		Owner:      s.name,
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), 201, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	//write用户可以修改和删除他人Space
	_, r, err = Api.ModelApi.V1ModelIdPut(Auth2, id, swagger.ControllerReqToUpdateModel{
		Desc: "model desc new",
	})
	assert.Equal(s.T(), 202, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = Api.ModelApi.V1ModelIdDelete(Auth2, id)
	assert.Equal(s.T(), 204, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 拥有admin权限的用户可以修改和删除他人的模型
func (s *SuiteOrgModel) TestOrgAdminUpdateDeleteOthersModel() {
	data, r, err := Api.ModelApi.V1ModelPost(Auth2, swagger.ControllerReqToCreateModel{
		Name:       "testmodel",
		Owner:      s.name,
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), 201, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	//admin用户可以修改和删除他人模型
	_, r, err = Api.ModelApi.V1ModelIdPut(Auth, id, swagger.ControllerReqToUpdateModel{
		Desc: "model desc new",
	})
	assert.Equal(s.T(), 202, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = Api.ModelApi.V1ModelIdDelete(Auth, id)
	assert.Equal(s.T(), 204, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 拥有contribute权限的用户可以创建修改和删除自己的模型
func (s *SuiteOrgModel) TestOrgContributorCreateUpdateDelete() {
	_, r, err := Api.OrganizationApi.V1OrganizationNameMemberPut(Auth, swagger.ControllerOrgMemberEditRequest{
		User: "test2",
		Role: "contributor",
	}, s.name)
	assert.Equal(s.T(), 202, r.StatusCode)
	assert.Nil(s.T(), err)

	data, r, err := Api.ModelApi.V1ModelPost(Auth2, swagger.ControllerReqToCreateModel{
		Name:       "testmodel",
		Owner:      s.name,
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), 201, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	_, r, err = Api.ModelApi.V1ModelIdPut(Auth2, id, swagger.ControllerReqToUpdateModel{
		Desc: "model desc new",
	})
	assert.Equal(s.T(), 202, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = Api.ModelApi.V1ModelIdDelete(Auth2, id)
	assert.Equal(s.T(), 204, r.StatusCode)
	assert.Nil(s.T(), err)

	_, r, err = Api.OrganizationApi.V1OrganizationNameMemberPut(Auth, swagger.ControllerOrgMemberEditRequest{
		User: "test2",
		Role: "write",
	}, s.name)
	assert.Equal(s.T(), 202, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 拥有contribute权限的用户不可以修改或删除他人模型
func (s *SuiteOrgModel) TestOrgContributorCantUpdateDeleteOthersModel() {
	_, r, err := Api.OrganizationApi.V1OrganizationNameMemberPut(Auth, swagger.ControllerOrgMemberEditRequest{
		User: "test2",
		Role: "contributor",
	}, s.name)

	assert.Equal(s.T(), 202, r.StatusCode)
	assert.Nil(s.T(), err)

	data, r, err := Api.ModelApi.V1ModelPost(Auth, swagger.ControllerReqToCreateModel{
		Name:       "testmodel",
		Owner:      s.name,
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), 201, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	r, err = Api.ModelApi.V1ModelIdDelete(Auth2, id)
	assert.Equal(s.T(), 404, r.StatusCode)
	assert.NotNil(s.T(), err)

	_, r, err = Api.ModelApi.V1ModelIdPut(Auth2, id, swagger.ControllerReqToUpdateModel{
		Desc: "model desc new",
	})

	// 拥有contribute权限的用户不可以修改他人模型
	assert.Equal(s.T(), 404, r.StatusCode)
	assert.NotNil(s.T(), err)

	// 拥有contribute权限的用户不可以删除他人模型
	r, err = Api.ModelApi.V1ModelIdDelete(Auth2, id)
	assert.Equal(s.T(), 404, r.StatusCode)
	assert.NotNil(s.T(), err)

	r, err = Api.ModelApi.V1ModelIdDelete(Auth, id)
	assert.Equal(s.T(), 204, r.StatusCode)
	assert.Nil(s.T(), err)

	_, r, err = Api.OrganizationApi.V1OrganizationNameMemberPut(Auth, swagger.ControllerOrgMemberEditRequest{
		User: "test2",
		Role: "write",
	}, s.name)

	assert.Equal(s.T(), 202, r.StatusCode)
	assert.Nil(s.T(), err)
}
func TestOrgModel(t *testing.T) {
	suite.Run(t, new(SuiteOrgModel))
}
