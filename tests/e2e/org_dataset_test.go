/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package e2e

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	swaggerRest "e2e/client_rest"
)

// SuiteOrgDataset used for testing
type SuiteOrgDataset struct {
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
func (s *SuiteOrgDataset) SetupSuite() {
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
}

// TearDownSuite used for testing
func (s *SuiteOrgDataset) TearDownSuite() {
	r, err := ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest, s.name)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestDeleteOrgContainsDataset used for testing
// 当组织下有dataset时，删除组织失败
func (s *SuiteOrgDataset) TestDeleteOrgContainsDataset() {
	data, r, err := ApiRest.DatasetApi.V1DatasetPost(AuthRest, swaggerRest.ControllerReqToCreateDataset{
		Name:       "tempDataset",
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

	// 清空Dataset
	id := getString(s.T(), data.Data)

	r, err = ApiRest.DatasetApi.V1DatasetIdDelete(AuthRest, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestOrgReadMemberCantCreateUpdateDeleteDataset used for testing
// 拥有read权限的用户不能创建数据集，不能修改和删除他人数据集
func (s *SuiteOrgDataset) TestOrgReadMemberCantCreateUpdateDeleteDataset() {
	_, r, err := ApiRest.OrganizationApi.V1OrganizationNameMemberPut(AuthRest, swaggerRest.ControllerOrgMemberEditRequest{
		User: "test2",
		Role: "read",
	}, s.name)

	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	_, r, err = ApiRest.DatasetApi.V1DatasetPost(AuthRest2, swaggerRest.ControllerReqToCreateDataset{
		Name:       "testdataset",
		Owner:      s.name,
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusForbidden, r.StatusCode)
	assert.NotNil(s.T(), err)

	data, r, err := ApiRest.DatasetApi.V1DatasetPost(AuthRest, swaggerRest.ControllerReqToCreateDataset{
		Name:       "testdataset",
		Owner:      s.name,
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	//read用户不能修改和删除他人数据集
	_, r, err = ApiRest.DatasetApi.V1DatasetIdPut(AuthRest2, id, swaggerRest.ControllerReqToUpdateDataset{
		Desc: "dataset desc new",
	})
	assert.Equal(s.T(), http.StatusForbidden, r.StatusCode)
	assert.NotNil(s.T(), err)

	r, err = ApiRest.DatasetApi.V1DatasetIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusForbidden, r.StatusCode)
	assert.NotNil(s.T(), err)

	_, r, err = ApiRest.OrganizationApi.V1OrganizationNameMemberPut(AuthRest, swaggerRest.ControllerOrgMemberEditRequest{
		User: "test2",
		Role: "write",
	}, s.name)

	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.DatasetApi.V1DatasetIdDelete(AuthRest, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestOrgWriteCreateDeleteDataset used for testing
// 拥有write权限的用户可以创建和删除数据集
func (s *SuiteOrgDataset) TestOrgWriteCreateDeleteDataset() {
	datasetParam := swaggerRest.ControllerReqToCreateDataset{
		Name:       "testdataset",
		Owner:      s.name,
		License:    "mit",
		Visibility: "public",
	}
	data, r, err := ApiRest.DatasetApi.V1DatasetPost(AuthRest2, datasetParam)

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	// 重复创建数据集返回400
	_, r, err = ApiRest.DatasetApi.V1DatasetPost(AuthRest2, datasetParam)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	r, err = ApiRest.DatasetApi.V1DatasetIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestOrgWriteUpdateDeleteOthersDataset used for testing
// 拥有write权限的用户可以修改和删除他人的数据集
func (s *SuiteOrgDataset) TestOrgWriteUpdateDeleteOthersDataset() {
	data, r, err := ApiRest.DatasetApi.V1DatasetPost(AuthRest, swaggerRest.ControllerReqToCreateDataset{
		Name:       "testdataset",
		Owner:      s.name,
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	//write用户可以修改和删除他人Space
	_, r, err = ApiRest.DatasetApi.V1DatasetIdPut(AuthRest2, id, swaggerRest.ControllerReqToUpdateDataset{
		Desc: "dataset desc new",
	})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.DatasetApi.V1DatasetIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestOrgAdminUpdateDeleteOthersDataset used for testing
// 拥有admin权限的用户可以修改和删除他人的数据集
func (s *SuiteOrgDataset) TestOrgAdminUpdateDeleteOthersDataset() {
	data, r, err := ApiRest.DatasetApi.V1DatasetPost(AuthRest2, swaggerRest.ControllerReqToCreateDataset{
		Name:       "testdataset",
		Owner:      s.name,
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	//admin用户可以修改和删除他人数据集
	_, r, err = ApiRest.DatasetApi.V1DatasetIdPut(AuthRest, id, swaggerRest.ControllerReqToUpdateDataset{
		Desc: "dataset desc new",
	})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.DatasetApi.V1DatasetIdDelete(AuthRest, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestDatasetMaxPerOrg used for testing
// 单个组织下可以创建最多数据集个数
func (s *SuiteOrgDataset) TestDatasetMaxPerOrg() {
	// 创建testdataset
	data1, r, err := ApiRest.DatasetApi.V1DatasetPost(AuthRest, swaggerRest.ControllerReqToCreateDataset{
		Name:       "testdataset1",
		Owner:      s.name,
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	id1 := getString(s.T(), data1.Data)

	data2, r, err := ApiRest.DatasetApi.V1DatasetPost(AuthRest, swaggerRest.ControllerReqToCreateDataset{
		Name:       "testdataset2",
		Owner:      s.name,
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	id2 := getString(s.T(), data2.Data)

	data3, r, err := ApiRest.DatasetApi.V1DatasetPost(AuthRest, swaggerRest.ControllerReqToCreateDataset{
		Name:       "testdataset3",
		Owner:      s.name,
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	id3 := getString(s.T(), data3.Data)

	data4, r, err := ApiRest.DatasetApi.V1DatasetPost(AuthRest, swaggerRest.ControllerReqToCreateDataset{
		Name:       "testdataset4",
		Owner:      s.name,
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	id4 := getString(s.T(), data4.Data)

	// 创建达到最大允许个数，创建失败
	_, r, err = ApiRest.DatasetApi.V1DatasetPost(AuthRest, swaggerRest.ControllerReqToCreateDataset{
		Name:       "testdataset5",
		Owner:      s.name,
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	r, err = ApiRest.DatasetApi.V1DatasetIdDelete(AuthRest, id1)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.DatasetApi.V1DatasetIdDelete(AuthRest, id2)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.DatasetApi.V1DatasetIdDelete(AuthRest, id3)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.DatasetApi.V1DatasetIdDelete(AuthRest, id4)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestOrgDataset used for testing
func TestOrgDataset(t *testing.T) {
	suite.Run(t, new(SuiteOrgDataset))
}
