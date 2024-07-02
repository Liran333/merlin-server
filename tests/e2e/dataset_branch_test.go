/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	swaggerRest "e2e/client_rest"
)

type SuiteDatasetBranch struct {
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
	datasetId    string
	branchName   string
}

func (s *SuiteDatasetBranch) SetupSuite() {
	s.name = "testorg"
	s.fullname = "testorgfull"
	s.avatarid = "https://avatars.githubusercontent.com/u/2853724?v=1"
	s.allowRequest = true
	s.defaultRole = "admin"
	s.website = "https://www.datasetfoundry.cn"
	s.desc = "test org desc"
	s.owner = "test1" // this name is hard code in init-env.sh

	data1, r, err := ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, swaggerRest.ControllerOrgCreateRequest{
		Name:        s.name,
		Fullname:    s.fullname,
		AvatarId:    s.avatarid,
		Website:     s.website,
		Description: s.desc,
	})

	o := getData(s.T(), data1.Data)

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), "", o["id"])
	s.orgId = getString(s.T(), o["id"])

	data2, r, err := ApiRest.OrganizationApi.V1InvitePost(AuthRest, swaggerRest.ControllerOrgInviteMemberRequest{
		OrgName: s.name,
		User:    "test2",
		Role:    "write",
		Msg:     "invite me",
	})

	assert.Equalf(s.T(), http.StatusCreated, r.StatusCode, data2.Msg)
	assert.Nil(s.T(), err)

	// 被邀请人接受邀请
	data3, r, err := ApiRest.OrganizationApi.V1InvitePut(AuthRest2, swaggerRest.ControllerOrgAcceptMemberRequest{
		OrgName: s.name,
		Msg:     "ok",
	})

	assert.Equalf(s.T(), http.StatusAccepted, r.StatusCode, data3.Msg)
	assert.Nil(s.T(), err)

	molData, r, err := ApiRest.DatasetApi.V1DatasetPost(AuthRest, swaggerRest.ControllerReqToCreateDataset{
		Name:       "testdataset",
		Owner:      s.name,
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	s.datasetId = getString(s.T(), molData.Data)
}

func (s *SuiteDatasetBranch) TearDownSuite() {
	r, err := ApiRest.DatasetApi.V1DatasetIdDelete(AuthRest, s.datasetId)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
	r, err = ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest, s.name)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 拥有write权限的用户可以创建和删除分支
func (s *SuiteDatasetBranch) TestOrgWriteCreateDeleteBranch() {
	branchName := "newbranch1"
	_, r, err := ApiRest.BranchRestfulApi.V1BranchTypeOwnerRepoPost(AuthRest2, "dataset", s.name, "testdataset",
		swaggerRest.ControllerRestfulReqToCreateBranch{
			BaseBranch: "main",
			Branch:     branchName,
		})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	time.Sleep(2 * time.Second)
	r, err = ApiRest.BranchRestfulApi.V1BranchTypeOwnerRepoBranchDelete(AuthRest2, "dataset", s.name, "testdataset", branchName)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 拥有read权限的用户不可以创建分支
func (s *SuiteDatasetBranch) TestOrgReadMemberCantCreateBranch() {
	branchName := "newbranch3"
	_, r, err := ApiRest.OrganizationApi.V1OrganizationNameMemberPut(AuthRest, swaggerRest.ControllerOrgMemberEditRequest{
		User: "test2",
		Role: "read",
	}, s.name)

	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	_, r, err = ApiRest.BranchRestfulApi.V1BranchTypeOwnerRepoPost(AuthRest2, "dataset", s.name, "testdataset",
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
func (s *SuiteDatasetBranch) TestOrgAdminCreateDeleteBranch() {
	branchName := "newbranch4"
	_, r, err := ApiRest.BranchRestfulApi.V1BranchTypeOwnerRepoPost(AuthRest, "dataset", s.name, "testdataset",
		swaggerRest.ControllerRestfulReqToCreateBranch{
			BaseBranch: "main",
			Branch:     branchName,
		})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	time.Sleep(2 * time.Second)
	r, err = ApiRest.BranchRestfulApi.V1BranchTypeOwnerRepoBranchDelete(AuthRest, "dataset", s.name, "testdataset", branchName)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 用户可以在自己的仓库创建和删除分支
func (s *SuiteDatasetBranch) TestOrgUserCanCreateDeleteBranch() {
	branchName := "newbranch7"
	data, r, err := ApiRest.DatasetApi.V1DatasetPost(AuthRest2, swaggerRest.ControllerReqToCreateDataset{
		Name:       "test2dataset",
		Owner:      s.name,
		License:    "mit",
		Visibility: "public",
	})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	id := getString(s.T(), data.Data)

	_, r, err = ApiRest.BranchRestfulApi.V1BranchTypeOwnerRepoPost(AuthRest2, "dataset", s.name, "test2dataset",
		swaggerRest.ControllerRestfulReqToCreateBranch{
			BaseBranch: "main",
			Branch:     branchName,
		})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	time.Sleep(2 * time.Second)
	r, err = ApiRest.BranchRestfulApi.V1BranchTypeOwnerRepoBranchDelete(AuthRest2, "dataset", s.name, "test2dataset", branchName)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.DatasetApi.V1DatasetIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 使用无效的分支名会导致创建分支失败
func (s *SuiteDatasetBranch) TestOrgUserCreateInvalidBranch() {
	branchName := "invild#branch"
	data, r, err := ApiRest.DatasetApi.V1DatasetPost(AuthRest2, swaggerRest.ControllerReqToCreateDataset{
		Name:       "test2dataset",
		Owner:      s.name,
		License:    "mit",
		Visibility: "public",
	})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	id := getString(s.T(), data.Data)

	_, r, err = ApiRest.BranchRestfulApi.V1BranchTypeOwnerRepoPost(AuthRest2, "dataset",
		s.name, "test2dataset", swaggerRest.ControllerRestfulReqToCreateBranch{
			BaseBranch: "main",
			Branch:     branchName,
		})
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	time.Sleep(2 * time.Second)

	r, err = ApiRest.DatasetApi.V1DatasetIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 测试举报内容有非法字符
func (s *SuiteInternalDataset) TestDatasetReport() {
	// 创建数据集
	data, r, err := ApiRest.DatasetApi.V1DatasetPost(AuthRest2, swaggerRest.ControllerReqToCreateDataset{
		Name:       "testdataset",
		Owner:      "test2",
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)
	// 提交举报信息
	_, r, err = ApiRest.DatasetRestfulApi.V1DatasetReportPost(AuthRest2, swaggerRest.ControllerReqReportDatasetEmail{
		Name:  "testdataset",
		Msg:   "test Reprot</a>",
		Owner: "test2",
	})
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), http.StatusForbidden, r.StatusCode)
	// 删除数据集
	r, err = ApiRest.DatasetApi.V1DatasetIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 私人数据集举报失败
func (s *SuiteInternalDataset) TestPrivateDataset() {
	// 创建数据集
	data, r, err := ApiRest.DatasetApi.V1DatasetPost(AuthRest2, swaggerRest.ControllerReqToCreateDataset{
		Name:       "testdataset",
		Owner:      "test2",
		License:    "mit",
		Visibility: "private",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	id := getString(s.T(), data.Data)
	// 提交举报信息
	_, r, err = ApiRest.DatasetRestfulApi.V1DatasetReportPost(AuthRest, swaggerRest.ControllerReqReportDatasetEmail{
		Name:  "testdataset",
		Msg:   "test report",
		Owner: "test2",
	})
	assert.NotNil(s.T(), err)
	fmt.Println(err.Error())
	assert.Equal(s.T(), http.StatusForbidden, r.StatusCode)
	// 删除数据集
	r, err = ApiRest.DatasetApi.V1DatasetIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 测试举报内容有非法字符(web)
func (s *SuiteInternalDataset) TestDatasetReportWeb() {
	// 创建数据集
	data, r, err := ApiRest.DatasetApi.V1DatasetPost(AuthRest2, swaggerRest.ControllerReqToCreateDataset{
		Name:       "testdataset",
		Owner:      "test2",
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)
	// 提交举报信息
	_, r, err = ApiRest.DatasetRestfulApi.V1DatasetReportPost(AuthRest2, swaggerRest.ControllerReqReportDatasetEmail{
		Name:  "testdataset",
		Msg:   "test Reprot</a>",
		Owner: "test2",
	})
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), http.StatusForbidden, r.StatusCode)
	// 删除数据集
	r, err = ApiRest.DatasetApi.V1DatasetIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 私人数据集举报失败(web)
func (s *SuiteInternalDataset) TestPrivateDatasetWeb() {
	// 创建数据集
	data, r, err := ApiRest.DatasetApi.V1DatasetPost(AuthRest2, swaggerRest.ControllerReqToCreateDataset{
		Name:       "testdataset",
		Owner:      "test2",
		License:    "mit",
		Visibility: "private",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	id := getString(s.T(), data.Data)
	// 提交举报信息
	_, r, err = ApiRest.DatasetRestfulApi.V1DatasetReportPost(AuthRest, swaggerRest.ControllerReqReportDatasetEmail{
		Name:  "testdataset",
		Msg:   "test report",
		Owner: "test2",
	})
	assert.NotNil(s.T(), err)
	fmt.Println(err.Error())
	assert.Equal(s.T(), http.StatusForbidden, r.StatusCode)
	// 删除数据集
	r, err = ApiRest.DatasetApi.V1DatasetIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

func TestDatasetBranch(t *testing.T) {
	suite.Run(t, new(SuiteDatasetBranch))
}
