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

// SuiteInternalModel used for testing
type SuiteInternalModel struct {
	suite.Suite
}

// SetupSuite used for testing
func (s *SuiteInternalModel) SetupSuite() {
}

// TearDownSuite used for testing
func (s *SuiteInternalModel) TearDownSuite() {
}

// TestGetModel used for testing
// 获取模型成功
func (s *SuiteInternalModel) TestGetModel() {
	// 创建模型
	data, r, err := ApiRest.ModelApi.V1ModelPost(AuthRest2, swaggerRest.ControllerReqToCreateModel{
		Name:       "testmodel",
		Owner:      "test2",
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	// 获取模型
	modelRes, r, err := ApiInteral.ModelInternalApi.V1ModelIdGet(Interal, id)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	modelData := getData(s.T(), modelRes.Data)
	assert.Equal(s.T(), id, modelData["id"])
	assert.Equal(s.T(), "testmodel", modelData["name"])
	assert.Equal(s.T(), "test2", modelData["owner"])
	assert.Equal(s.T(), "public", modelData["visibility"])

	// 删除模型
	r, err = ApiRest.ModelApi.V1ModelIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestGetModelFailed used for testing
// 获取模型失败
func (s *SuiteInternalModel) TestGetModelFailed() {
	// 模型不存在
	unExistedId := "0"
	_, r, err := ApiInteral.ModelInternalApi.V1ModelIdGet(Interal, unExistedId)
	assert.Equal(s.T(), http.StatusNotFound, r.StatusCode)
	assert.NotNil(s.T(), err)

	// 参数无效
	unExistedId = "test"
	_, r, err = ApiInteral.ModelInternalApi.V1ModelIdGet(Interal, unExistedId)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestInternalModelSetLabel used for testing
// 模型label设置
func (s *SuiteInternalModel) TestInternalModelSetLabel() {
	// 创建模型
	data, r, err := ApiRest.ModelApi.V1ModelPost(AuthRest2, swaggerRest.ControllerReqToCreateModel{
		Name:       "testmodel",
		Owner:      "test2",
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	_, r, err = ApiInteral.ModelInternalApi.V1ModelIdLabelPut(Interal, id, swaggerInternal.ControllerReqToCreateModel{
		Desc:       "testmodel",
		Fullname:   "fullName",
		License:    "mit",
		Name:       "testmodel",
		Owner:      "testuser1",
		Visibility: "1",
	})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	// 删除模型
	r, err = ApiRest.ModelApi.V1ModelIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestInternalModel used for testing
func TestInternalModel(t *testing.T) {
	suite.Run(t, new(SuiteInternalModel))
}
