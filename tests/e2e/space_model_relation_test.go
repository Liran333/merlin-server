/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package e2e

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	swaggerInternal "e2e/client_internal"
	swaggerRest "e2e/client_rest"
)

// SuiteSpaceModelRelation used for testing
type SuiteSpaceModelRelation struct {
	suite.Suite
	notExistId          string
	modelIdTest1Public  string
	modelIdTest2Public  string
	spaceIdTest2Public  string
	modelIdTest2Private string
}

// SetupSuite used for testing
func (s *SuiteSpaceModelRelation) SetupSuite() {
	// 创建模型
	data, r, err := ApiRest.ModelApi.V1ModelPost(AuthRest, swaggerRest.ControllerReqToCreateModel{
		Name:       "testmodel1",
		Owner:      "test1",
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	s.modelIdTest1Public = getString(s.T(), data.Data)

	data, r, err = ApiRest.ModelApi.V1ModelPost(AuthRest2, swaggerRest.ControllerReqToCreateModel{
		Name:       "testmodel2",
		Owner:      "test2",
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	s.modelIdTest2Public = getString(s.T(), data.Data)

	data, r, err = ApiRest.ModelApi.V1ModelPost(AuthRest2, swaggerRest.ControllerReqToCreateModel{
		Name:       "modelprivate",
		Owner:      "test2",
		License:    "mit",
		Visibility: "private",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	s.modelIdTest2Private = getString(s.T(), data.Data)

	// 创建空间
	spaceParam := swaggerRest.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		License:    "mit",
		BaseImage:  "python3.8-pytorch2.1",
		Name:       "testspace1",
		Owner:      "test2",
		Sdk:        "gradio",
		Visibility: "public",
	}
	spaData, r, err := ApiRest.SpaceApi.V1SpacePost(AuthRest2, spaceParam)

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	s.spaceIdTest2Public = getString(s.T(), spaData.Data)

	s.notExistId = "-1"
}

// TearDownSuite used for testing
func (s *SuiteSpaceModelRelation) TearDownSuite() {
	// 删除模型
	r, err := ApiRest.ModelApi.V1ModelIdDelete(AuthRest, s.modelIdTest1Public)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.ModelApi.V1ModelIdDelete(AuthRest2, s.modelIdTest2Public)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.ModelApi.V1ModelIdDelete(AuthRest2, s.modelIdTest2Private)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	// 删除空间
	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest2, s.spaceIdTest2Public)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestUpdateSpaceModelsSuccess used for testing
// 空间关联模型成功
func (s *SuiteSpaceModelRelation) TestUpdateSpaceModelsSuccess() {
	ids := []string{"test1/testmodel1", "test2/testmodel2"}
	modelIdsParam := swaggerInternal.ControllerModeIds{Ids: ids}
	_, r, err := ApiInteral.SpaceInternalApi.V1SpaceIdModelPut(Interal, s.spaceIdTest2Public, modelIdsParam)

	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestUpdateInvalidModels used for testing
// 空间关联不存在的模型
func (s *SuiteSpaceModelRelation) TestUpdateInvalidModels() {
	// 空间关联公开模型、私有模型、不存在的模型
	ids := []string{"test1/testmodel1", "test2/modelprivate", "test2/modelNotExists", "invalidmodel"}
	modelIdsParam := swaggerInternal.ControllerModeIds{Ids: ids}
	_, r, err := ApiInteral.SpaceInternalApi.V1SpaceIdModelPut(Interal, s.spaceIdTest2Public, modelIdsParam)

	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestUpdateNotExistSpaceModelsFail used for testing
// 不存在的空间关联模型
func (s *SuiteSpaceModelRelation) TestUpdateNotExistSpaceModelsFail() {
	ids := []string{"test1/testmodel1"}
	modelIdsParam := swaggerInternal.ControllerModeIds{Ids: ids}
	_, r, err := ApiInteral.SpaceInternalApi.V1SpaceIdModelPut(Interal, s.notExistId, modelIdsParam)

	assert.Equal(s.T(), http.StatusNotFound, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestInternalGetSpacesByModelId used for testing
// internal接口测试，使用model id查询关联的space id
func (s *SuiteSpaceModelRelation) TestInternalGetSpacesByModelId() {
	detail, r, err := ApiInteral.ModelInternalApi.V1ModelRelationIdSpaceGet(Interal, s.modelIdTest1Public)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	space := getData(s.T(), detail.Data)
	assert.NotEqual(s.T(), space["space_id"], "")
}

func TestSpaceModelRelation(t *testing.T) {
	suite.Run(t, new(SuiteSpaceModelRelation))
}
