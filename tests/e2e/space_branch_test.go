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

	_, r, err = ApiRest.OrganizationApi.V1InvitePost(AuthRest, swaggerRest.ControllerOrgInviteMemberRequest{
		OrgName: s.name,
		User:    "test2",
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

	spaData, r, err := ApiRest.SpaceApi.V1SpacePost(AuthRest, swaggerRest.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		License:    "mit",
		BaseImage:  "python3.8-pytorch2.1",
		Name:       "testspace",
		Owner:      s.name,
		Sdk:        "gradio",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	s.spaceId = getString(s.T(), spaData.Data)
}

func (s *SuiteSpaceBranch) TearDownSuite() {
	r, err := ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest, s.spaceId)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
	r, err = ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest, s.name)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 拥有write权限的用户可以创建和删除分支
func (s *SuiteSpaceBranch) TestOrgWriteCreateDeleteBranch() {
	branchName := "newbranch1"
	_, r, err := ApiRest.BranchRestfulApi.V1BranchTypeOwnerRepoPost(AuthRest2, "space", s.name, "testspace",
		swaggerRest.ControllerRestfulReqToCreateBranch{
			BaseBranch: "main",
			Branch:     branchName,
		})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	// 重复创建分支返回400
	_, r, err = ApiRest.BranchRestfulApi.V1BranchTypeOwnerRepoPost(AuthRest2, "space", s.name, "testspace",
		swaggerRest.ControllerRestfulReqToCreateBranch{
			BaseBranch: "main",
			Branch:     branchName,
		})
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	time.Sleep(2 * time.Second)
	r, err = ApiRest.BranchRestfulApi.V1BranchTypeOwnerRepoBranchDelete(AuthRest2, "space", s.name,
		"testspace", branchName)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 拥有read权限的用户不可以创建分支
func (s *SuiteSpaceBranch) TestOrgReadMemberCantCreateBranch() {
	branchName := "newbranch3"
	_, r, err := ApiRest.OrganizationApi.V1OrganizationNameMemberPut(AuthRest,
		swaggerRest.ControllerOrgMemberEditRequest{
			User: "test2",
			Role: "read",
		}, s.name)

	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	_, r, err = ApiRest.BranchRestfulApi.V1BranchTypeOwnerRepoPost(AuthRest2, "space", s.name, "testspace",
		swaggerRest.ControllerRestfulReqToCreateBranch{
			BaseBranch: "main",
			Branch:     branchName,
		})

	assert.Equal(s.T(), http.StatusForbidden, r.StatusCode)
	assert.NotNil(s.T(), err)

	_, r, err = ApiRest.OrganizationApi.V1OrganizationNameMemberPut(AuthRest,
		swaggerRest.ControllerOrgMemberEditRequest{
			User: "test2",
			Role: "write",
		}, s.name)

	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 拥有admin权限的用户可以创建和删除分支
func (s *SuiteSpaceBranch) TestOrgAdminCreateDeleteBranch() {
	branchName := "newbranch4"
	_, r, err := ApiRest.BranchRestfulApi.V1BranchTypeOwnerRepoPost(AuthRest, "space", s.name, "testspace",
		swaggerRest.ControllerRestfulReqToCreateBranch{
			BaseBranch: "main",
			Branch:     branchName,
		})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	time.Sleep(2 * time.Second)
	r, err = ApiRest.BranchRestfulApi.V1BranchTypeOwnerRepoBranchDelete(AuthRest, "space", s.name,
		"testspace", branchName)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 用户可以在自己的仓库创建和删除分支
func (s *SuiteSpaceBranch) TestOrgUserCanCreateDeleteBranch() {
	branchName := "newbranch7"
	data, r, err := ApiRest.SpaceApi.V1SpacePost(AuthRest2, swaggerRest.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		License:    "mit",
		BaseImage:  "python3.8-pytorch2.1",
		Name:       "test2space",
		Owner:      "test2",
		Sdk:        "gradio",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	id := getString(s.T(), data.Data)

	_, r, err = ApiRest.BranchRestfulApi.V1BranchTypeOwnerRepoPost(AuthRest2, "space", "test2",
		"test2space", swaggerRest.ControllerRestfulReqToCreateBranch{
			BaseBranch: "main",
			Branch:     branchName,
		})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	time.Sleep(2 * time.Second)
	r, err = ApiRest.BranchRestfulApi.V1BranchTypeOwnerRepoBranchDelete(AuthRest2, "space", "test2",
		"test2space", branchName)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 使用无效的分支名会导致创建分支失败
func (s *SuiteSpaceBranch) TestOrgUserCreateInvalidBranch() {
	branchName := "invild#branch"
	data, r, err := ApiRest.SpaceApi.V1SpacePost(AuthRest2, swaggerRest.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		License:    "mit",
		BaseImage:  "python3.8-pytorch2.1",
		Name:       "test2space",
		Owner:      "test2",
		Sdk:        "gradio",
		Visibility: "public",
	})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	id := getString(s.T(), data.Data)

	_, r, err = ApiRest.BranchRestfulApi.V1BranchTypeOwnerRepoPost(AuthRest2, "space",
		s.name, "test2space", swaggerRest.ControllerRestfulReqToCreateBranch{
			BaseBranch: "main",
			Branch:     branchName,
		})
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

func TestSpaceBranch(t *testing.T) {
	suite.Run(t, new(SuiteSpaceBranch))
}
