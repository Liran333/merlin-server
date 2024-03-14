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

type SuiteSpaceBranch struct {
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
	spaceId      string
	branchName   string
}

func (s *SuiteSpaceBranch) SetupSuite() {
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

	data, r, err = Api.SpaceApi.V1SpacePost(Auth, swagger.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		InitReadme: true,
		License:    "mit",
		Name:       "testspace",
		Owner:      s.name,
		Sdk:        "gradio",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	s.spaceId = getString(s.T(), data.Data)
}

func (s *SuiteSpaceBranch) TearDownSuite() {
	r, err := Api.SpaceApi.V1SpaceIdDelete(Auth, s.spaceId)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
	r, err = Api.OrganizationApi.V1OrganizationNameDelete(Auth, s.name)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 拥有write权限的用户可以创建和删除分支
func (s *SuiteSpaceBranch) TestOrgWriteCreateDeleteBranch() {
	branchName := "newbranch1"
	_, r, err := Api.BranchRestfulApi.V1BranchTypeOwnerRepoPost(Auth2, "space", s.name, "testspace",
		swagger.ControllerRestfulReqToCreateBranch{
			BaseBranch: "main",
			Branch:     branchName,
		})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	// 重复创建分支返回400
	_, r, err = Api.BranchRestfulApi.V1BranchTypeOwnerRepoPost(Auth2, "space", s.name, "testspace",
		swagger.ControllerRestfulReqToCreateBranch{
			BaseBranch: "main",
			Branch:     branchName,
		})
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	r, err = Api.BranchRestfulApi.V1BranchTypeOwnerRepoBranchDelete(Auth2, "space", s.name, "testspace", branchName)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 拥有read权限的用户不可以创建分支
func (s *SuiteSpaceBranch) TestOrgReadMemberCantCreateBranch() {
	branchName := "newbranch3"
	_, r, err := Api.OrganizationApi.V1OrganizationNameMemberPut(Auth, swagger.ControllerOrgMemberEditRequest{
		User: "test2",
		Role: "read",
	}, s.name)

	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	_, r, err = Api.BranchRestfulApi.V1BranchTypeOwnerRepoPost(Auth2, "space", s.name, "testspace",
		swagger.ControllerRestfulReqToCreateBranch{
			BaseBranch: "main",
			Branch:     branchName,
		})

	assert.Equal(s.T(), http.StatusForbidden, r.StatusCode)
	assert.NotNil(s.T(), err)

	_, r, err = Api.OrganizationApi.V1OrganizationNameMemberPut(Auth, swagger.ControllerOrgMemberEditRequest{
		User: "test2",
		Role: "write",
	}, s.name)

	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 拥有admin权限的用户可以创建和删除分支
func (s *SuiteSpaceBranch) TestOrgAdminCreateDeleteBranch() {
	branchName := "newbranch4"
	_, r, err := Api.BranchRestfulApi.V1BranchTypeOwnerRepoPost(Auth, "space", s.name, "testspace",
		swagger.ControllerRestfulReqToCreateBranch{
			BaseBranch: "main",
			Branch:     branchName,
		})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = Api.BranchRestfulApi.V1BranchTypeOwnerRepoBranchDelete(Auth, "space", s.name, "testspace", branchName)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 用户可以在自己的仓库创建和删除分支
func (s *SuiteSpaceBranch) TestOrgUserCanCreateDeleteBranch() {
	branchName := "newbranch7"
	data, r, err := Api.SpaceApi.V1SpacePost(Auth2, swagger.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		InitReadme: true,
		License:    "mit",
		Name:       "test2space",
		Owner:      "test2",
		Sdk:        "gradio",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	id := getString(s.T(), data.Data)

	_, r, err = Api.BranchRestfulApi.V1BranchTypeOwnerRepoPost(Auth2, "space", "test2", "test2space",
		swagger.ControllerRestfulReqToCreateBranch{
			BaseBranch: "main",
			Branch:     branchName,
		})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = Api.BranchRestfulApi.V1BranchTypeOwnerRepoBranchDelete(Auth2, "space", "test2", "test2space", branchName)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = Api.SpaceApi.V1SpaceIdDelete(Auth2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 使用无效的分支名会导致创建分支失败
func (s *SuiteSpaceBranch) TestOrgUserCreateInvalidBranch() {
	branchName := "invild#branch"
	data, r, err := Api.SpaceApi.V1SpacePost(Auth2, swagger.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		InitReadme: true,
		License:    "mit",
		Name:       "test2space",
		Owner:      "test2",
		Sdk:        "gradio",
		Visibility: "public",
	})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	id := getString(s.T(), data.Data)

	_, r, err = Api.BranchRestfulApi.V1BranchTypeOwnerRepoPost(Auth2, "space",
		s.name, "test2space", swagger.ControllerRestfulReqToCreateBranch{
			BaseBranch: "main",
			Branch:     branchName,
		})
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	r, err = Api.SpaceApi.V1SpaceIdDelete(Auth2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

func TestSpaceBranch(t *testing.T) {
	suite.Run(t, new(SuiteSpaceBranch))
}
