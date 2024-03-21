/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package e2e

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	swaggerRest "e2e/client_rest"
)

type SuiteModelBranch struct {
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
	modelId      string
	branchName   string
}

func (s *SuiteModelBranch) SetupSuite() {
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

	data, r, err = ApiRest.ModelApi.V1ModelPost(AuthRest, swaggerRest.ControllerReqToCreateModel{
		Name:       "testmodel",
		Owner:      s.name,
		License:    "mit",
		Visibility: "public",
		InitReadme: true,
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	s.modelId = getString(s.T(), data.Data)
}

func (s *SuiteModelBranch) TearDownSuite() {
	r, err := ApiRest.ModelApi.V1ModelIdDelete(AuthRest, s.modelId)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
	r, err = ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest, s.name)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 拥有write权限的用户可以创建和删除分支
func (s *SuiteModelBranch) TestOrgWriteCreateDeleteBranch() {
	branchName := "newbranch1"
	_, r, err := ApiRest.BranchRestfulApi.V1BranchTypeOwnerRepoPost(AuthRest2, "model", s.name, "testmodel",
		swaggerRest.ControllerRestfulReqToCreateBranch{
			BaseBranch: "main",
			Branch:     branchName,
		})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.BranchRestfulApi.V1BranchTypeOwnerRepoBranchDelete(AuthRest2, "model", s.name, "testmodel", branchName)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 拥有read权限的用户不可以创建分支
func (s *SuiteModelBranch) TestOrgReadMemberCantCreateBranch() {
	branchName := "newbranch3"
	_, r, err := ApiRest.OrganizationApi.V1OrganizationNameMemberPut(AuthRest, swaggerRest.ControllerOrgMemberEditRequest{
		User: "test2",
		Role: "read",
	}, s.name)

	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	_, r, err = ApiRest.BranchRestfulApi.V1BranchTypeOwnerRepoPost(AuthRest2, "model", s.name, "testmodel",
		swaggerRest.ControllerRestfulReqToCreateBranch{
			BaseBranch: "main",
			Branch:     branchName,
		})

	assert.Equal(s.T(), http.StatusForbidden, r.StatusCode)
	assert.NotNil(s.T(), err)

	_, r, err = ApiRest.OrganizationApi.V1OrganizationNameMemberPut(AuthRest, swaggerRest.ControllerOrgMemberEditRequest{
		User: "test2",
		Role: "write",
	}, s.name)

	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 拥有admin权限的用户可以创建和删除分支
func (s *SuiteModelBranch) TestOrgAdminCreateDeleteBranch() {
	branchName := "newbranch4"
	_, r, err := ApiRest.BranchRestfulApi.V1BranchTypeOwnerRepoPost(AuthRest, "model", s.name, "testmodel",
		swaggerRest.ControllerRestfulReqToCreateBranch{
			BaseBranch: "main",
			Branch:     branchName,
		})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.BranchRestfulApi.V1BranchTypeOwnerRepoBranchDelete(AuthRest, "model", s.name, "testmodel", branchName)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 用户可以在自己的仓库创建和删除分支
func (s *SuiteModelBranch) TestOrgUserCanCreateDeleteBranch() {
	branchName := "newbranch7"
	data, r, err := ApiRest.ModelApi.V1ModelPost(AuthRest2, swaggerRest.ControllerReqToCreateModel{
		Name:       "test2model",
		Owner:      s.name,
		License:    "mit",
		Visibility: "public",
		InitReadme: true,
	})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	id := getString(s.T(), data.Data)

	_, r, err = ApiRest.BranchRestfulApi.V1BranchTypeOwnerRepoPost(AuthRest2, "model", s.name, "test2model",
		swaggerRest.ControllerRestfulReqToCreateBranch{
			BaseBranch: "main",
			Branch:     branchName,
		})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	time.Sleep(1 * time.Second)

	r, err = ApiRest.BranchRestfulApi.V1BranchTypeOwnerRepoBranchDelete(AuthRest2, "model", s.name, "test2model", branchName)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.ModelApi.V1ModelIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 使用无效的分支名会导致创建分支失败
func (s *SuiteModelBranch) TestOrgUserCreateInvalidBranch() {
	branchName := "invild#branch"
	data, r, err := ApiRest.ModelApi.V1ModelPost(AuthRest2, swaggerRest.ControllerReqToCreateModel{
		Name:       "test2model",
		Owner:      s.name,
		License:    "mit",
		Visibility: "public",
		InitReadme: true,
	})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	id := getString(s.T(), data.Data)

	_, r, err = ApiRest.BranchRestfulApi.V1BranchTypeOwnerRepoPost(AuthRest2, "model",
		s.name, "test2model", swaggerRest.ControllerRestfulReqToCreateBranch{
			BaseBranch: "main",
			Branch:     branchName,
		})
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	time.Sleep(1 * time.Second)

	r, err = ApiRest.ModelApi.V1ModelIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

func TestModelBranch(t *testing.T) {
	suite.Run(t, new(SuiteModelBranch))
}
