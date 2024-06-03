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

// SuiteModel used for testing
type SuiteModel struct {
	suite.Suite
}

// SetupSuite used for testing
func (s *SuiteModel) SetupSuite() {
}

// TearDownSuite used for testing
func (s *SuiteModel) TearDownSuite() {
}

// TestListRecommendModel used for testing
// 测试获取推荐模型
func (s *SuiteModel) TestListRecommendModel() {
	// 创建模型
	data, r, err := ApiRest.ModelApi.V1ModelPost(AuthRest, swaggerRest.ControllerReqToCreateModel{
		Name:       "testmodel",
		Owner:      "test1",
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	modelRes, r, err := ApiWeb.ModelWebApi.V1ModelRecommendGet(AuthRest)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	models := modelRes.Data.Models

	assert.Equal(s.T(), 1, len(models))
	assert.Equal(s.T(), id, models[0].Id)
	assert.Equal(s.T(), "testmodel", models[0].Name)
	assert.Equal(s.T(), "test1", models[0].Owner)
	assert.Equal(s.T(), "public", models[0].Visibility)
	assert.Equal(s.T(), "mit", models[0].Labels.License)

	// 删除模型
	r, err = ApiRest.ModelApi.V1ModelIdDelete(AuthRest, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestModel used for testing
func TestModel(t *testing.T) {
	suite.Run(t, new(SuiteModel))
}
