/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package e2e

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	swaggerRest "e2e/client_rest"
)

// SuiteUserModel used for testing
type SuiteUserModel struct {
	suite.Suite
}

// SetupSuite used for testing
func (s *SuiteUserModel) SetupSuite() {
}

// TearDownSuite used for testing
func (s *SuiteUserModel) TearDownSuite() {
}

// TestUserCanCreateUpdateDeleteModel used for testing
// 可以创建模型到自己名下, 并且可以修改和删除自己名下的模型
func (s *SuiteUserModel) TestUserCanCreateUpdateDeleteModel() {
	modelParam := swaggerRest.ControllerReqToCreateModel{
		Name:       "testmodel",
		Owner:      "test2",
		License:    "mit",
		Visibility: "public",
	}
	data, r, err := ApiRest.ModelApi.V1ModelPost(AuthRest2, modelParam)

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	// 重复创建模型返回400
	_, r, err = ApiRest.ModelApi.V1ModelPost(AuthRest2, modelParam)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	_, r, err = ApiRest.ModelApi.V1ModelIdPut(AuthRest2, id, swaggerRest.ControllerReqToUpdateModel{
		Visibility: "public",
	})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.ModelApi.V1ModelIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 使用无效仓库名创建模型失败
func (s *SuiteUserModel) TestUserCreateUpdateInvalidModel() {
	_, r, err := ApiRest.ModelApi.V1ModelPost(AuthRest2, swaggerRest.ControllerReqToCreateModel{
		Name:       "invalid#testmodel",
		Owner:      "test2",
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	data, r, err := ApiRest.ModelApi.V1ModelPost(AuthRest2, swaggerRest.ControllerReqToCreateModel{
		Name:       "testmodel",
		Owner:      "test2",
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	r, err = ApiRest.ModelApi.V1ModelIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestNotLoginCantCreateModel used for testing
// 没登录用户不能创建模型
func (s *SuiteUserModel) TestNotLoginCantCreateModel() {
	_, r, err := ApiRest.ModelApi.V1ModelPost(context.Background(), swaggerRest.ControllerReqToCreateModel{
		Name:       "testmodel",
		Owner:      "test2",
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusUnauthorized, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestUserCanVisitSelfPublicModel used for testing
// 可以访问自己名下的公有模型
func (s *SuiteUserModel) TestUserCanVisitSelfPublicModel() {
	// 创建testmodel
	data, r, err := ApiRest.ModelApi.V1ModelPost(AuthRest2, swaggerRest.ControllerReqToCreateModel{
		Name:       "testmodel",
		Owner:      "test2",
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	// 获取test2名下的model成功
	detail, r, err := ApiRest.ModelRestfulApi.V1ModelOwnerNameGet(AuthRest2, "test2", "testmodel")
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)
	model := getData(s.T(), detail.Data)
	assert.Equal(s.T(), "testmodel", model["name"])

	// 创建testmodel2
	data, r, err = ApiRest.ModelApi.V1ModelPost(AuthRest2, swaggerRest.ControllerReqToCreateModel{
		Name:       "testmodel2",
		Owner:      "test2",
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id2 := getString(s.T(), data.Data)

	// 获取用户名下的所有model list
	list, r, err := ApiRest.ModelRestfulApi.V1ModelGet(AuthRest2, "test2", &swaggerRest.ModelRestfulApiV1ModelGetOpts{})
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	count := 0
	modelLists := getData(s.T(), list.Data)

	for i := 0; i < len(modelLists["models"].([]interface{})); i++ {
		model, ok := modelLists["models"].([]interface{})[i].(map[string]interface{})
		if !ok {
			continue
		}

		if _, ok := model["name"]; ok {
			count++
		}
	}
	assert.Equal(s.T(), countTwo, count)

	// 删除创建的模型
	r, err = ApiRest.ModelApi.V1ModelIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.ModelApi.V1ModelIdDelete(AuthRest2, id2)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestUserModel used for testing
func TestUserModel(t *testing.T) {
	suite.Run(t, new(SuiteUserModel))
}
