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
	data, r, err := Api.ModelApi.V1ModelPost(Auth, swagger.ControllerReqToCreateModel{
		Name:       "testmodel1",
		Owner:      "test1",
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	s.modelIdTest1Public = getString(s.T(), data.Data)

	data, r, err = Api.ModelApi.V1ModelPost(Auth2, swagger.ControllerReqToCreateModel{
		Name:       "testmodel2",
		Owner:      "test2",
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	s.modelIdTest2Public = getString(s.T(), data.Data)

	data, r, err = Api.ModelApi.V1ModelPost(Auth2, swagger.ControllerReqToCreateModel{
		Name:       "modelprivate",
		Owner:      "test2",
		License:    "mit",
		Visibility: "private",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	s.modelIdTest2Private = getString(s.T(), data.Data)

	// 创建空间
	spaceParam := swagger.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		InitReadme: false,
		License:    "mit",
		Name:       "testspace1",
		Owner:      "test2",
		Sdk:        "gradio",
		Visibility: "public",
	}
	data, r, err = Api.SpaceApi.V1SpacePost(Auth2, spaceParam)

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	s.spaceIdTest2Public = getString(s.T(), data.Data)

	s.notExistId = "-1"
}

// TearDownSuite used for testing
func (s *SuiteSpaceModelRelation) TearDownSuite() {
	// 删除模型
	r, err := Api.ModelApi.V1ModelIdDelete(Auth, s.modelIdTest1Public)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = Api.ModelApi.V1ModelIdDelete(Auth2, s.modelIdTest2Public)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = Api.ModelApi.V1ModelIdDelete(Auth2, s.modelIdTest2Private)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	// 删除空间
	r, err = Api.SpaceApi.V1SpaceIdDelete(Auth2, s.spaceIdTest2Public)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestUpdateSpaceModelsSuccess used for testing
// 空间关联模型成功
func (s *SuiteSpaceModelRelation) TestUpdateSpaceModelsSuccess() {
	ids := []string{"test1/testmodel1", "test2/testmodel2"}
	modelIdsParam := swagger.ControllerModeIds{Ids: ids}
	_, r, err := InteralApi.SpaceInternalApi.V1SpaceIdModelPut(Interal, s.spaceIdTest2Public, modelIdsParam)

	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestUpdateInvalidModels used for testing
// 空间关联不存在的模型
func (s *SuiteSpaceModelRelation) TestUpdateInvalidModels() {
	// 空间关联公开模型、私有模型、不存在的模型
	ids := []string{"test1/testmodel1", "test2/modelprivate", "test2/modelNotExists", "invalidmodel"}
	modelIdsParam := swagger.ControllerModeIds{Ids: ids}
	_, r, err := InteralApi.SpaceInternalApi.V1SpaceIdModelPut(Interal, s.spaceIdTest2Public, modelIdsParam)

	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestUpdateNotExistSpaceModelsFail used for testing
// 不存在的空间关联模型
func (s *SuiteSpaceModelRelation) TestUpdateNotExistSpaceModelsFail() {
	ids := []string{"test1/testmodel1"}
	modelIdsParam := swagger.ControllerModeIds{Ids: ids}
	_, r, err := InteralApi.SpaceInternalApi.V1SpaceIdModelPut(Interal, s.notExistId, modelIdsParam)

	assert.Equal(s.T(), http.StatusNotFound, r.StatusCode)
	assert.NotNil(s.T(), err)
}

func TestSpaceModelRelation(t *testing.T) {
	suite.Run(t, new(SuiteSpaceModelRelation))
}
