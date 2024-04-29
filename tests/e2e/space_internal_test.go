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

// SuiteUserSpace used for testing
type SuiteSpaceInternal struct {
	suite.Suite
}

// SetupSuite used for testing
func (s *SuiteSpaceInternal) SetupSuite() {
}

// TearDownSuite used for testing
func (s *SuiteSpaceInternal) TearDownSuite() {
}

// 通过ID查询space
func (s *SuiteSpaceInternal) TestSpaceInternalGetById() {
	// id非法，获取space失败
	_, r, err := ApiInteral.SpaceInternalApi.V1SpaceIdGet(Interal, "12")
	assert.Equal(s.T(), http.StatusNotFound, r.StatusCode)
	assert.NotNil(s.T(), err)

	// 创建space
	data, r, err := ApiRest.SpaceApi.V1SpacePost(AuthRest2, swaggerRest.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		License:    "mit",
		Name:       "testspace",
		Owner:      "test2",
		Sdk:        "gradio",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	// 获取space成功
	list, r, err := ApiInteral.SpaceInternalApi.V1SpaceIdGet(Interal, id)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	spaceLists := getData(s.T(), list.Data)
	assert.Equal(s.T(), "testspace", spaceLists["name"])

	// 删除space
	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// space下架测试
func (s *SuiteSpaceInternal) TestSpaceInternalDisable() {
	// 创建space
	data, r, err := ApiRest.SpaceApi.V1SpacePost(AuthRest2, swaggerRest.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		License:    "mit",
		Name:       "testspacedisable",
		Owner:      "test2",
		Sdk:        "gradio",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	_, r, err = ApiInteral.SpaceInternalApi.V1SpaceIdDisablePut(Interal, id)
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	detail, r, err := ApiRest.SpaceRestfulApi.V1SpaceOwnerNameGet(AuthRest2, "test2", "testspacedisable")
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	space := getData(s.T(), detail.Data)
	assert.Equal(s.T(), "testspacedisable", space["name"])
	assert.Equal(s.T(), "related_model_disabled", space["exception"])

	// 删除space
	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestUserSpace used for testing
func TestSpaceInternal(t *testing.T) {
	suite.Run(t, new(SuiteSpaceInternal))
}
