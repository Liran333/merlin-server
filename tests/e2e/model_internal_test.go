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

// TestInternalModelResetLabelSuccess used for testing
// 模型label设置成功
func (s *SuiteInternalModel) TestInternalModelResetLabelSuccess() {
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

	// 修改模型label
	resetLabelBody := swaggerInternal.ControllerReqToResetLabel{
		Frameworks: []string{"PyTorch"},
		License:    "apache-2.0",
		Task:       "document-question-answering",
	}
	_, r, err = ApiInteral.ModelInternalApi.V1ModelIdLabelPut(Interal, id, resetLabelBody)
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	// 获取修改成功后的模型label信息
	modelRes, r, err := ApiInteral.ModelInternalApi.V1ModelIdGet(Interal, id)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	modelData := getData(s.T(), modelRes.Data)
	assert.Equal(s.T(),
		map[string]interface{}(map[string]interface{}{"frameworks": []interface{}{"PyTorch"}, "license": "apache-2.0", "others": []interface{}{}, "task": "document-question-answering"}),
		modelData["labels"])

	// 修改模型label，其中frameworks有多个
	resetLabelBody = swaggerInternal.ControllerReqToResetLabel{
		Frameworks: []string{"PyTorch", "MindSpore"},
		License:    "apache-2.0",
		Task:       "copa",
	}
	_, r, err = ApiInteral.ModelInternalApi.V1ModelIdLabelPut(Interal, id, resetLabelBody)
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	// 获取修改成功后的模型label信息
	modelRes, r, err = ApiInteral.ModelInternalApi.V1ModelIdGet(Interal, id)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	modelData = getData(s.T(), modelRes.Data)
	expectedValue1 := map[string]interface{}{
		"frameworks": []interface{}{"PyTorch", "MindSpore"},
		"license":    "apache-2.0",
		"others":     []interface{}{},
		"task":       "copa",
	}
	expectedValue2 := map[string]interface{}{
		"frameworks": []interface{}{"MindSpore", "PyTorch"},
		"license":    "apache-2.0",
		"others":     []interface{}{},
		"task":       "copa",
	}
	assert.Contains(s.T(),
		[]map[string]interface{}{expectedValue1, expectedValue2},
		modelData["labels"])

	// 删除模型
	r, err = ApiRest.ModelApi.V1ModelIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestInternalModelResetLabelfail used for testing
// 模型label设置失败
func (s *SuiteInternalModel) TestInternalModelResetLabelfail() {
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

	// 使用非法字段修改模型label
	resetLabelBody := swaggerInternal.ControllerReqToResetLabel{
		Frameworks: []string{"PyTorch123"},
		License:    "apache-2.0",
		Task:       "copa123",
	}
	_, r, err = ApiInteral.ModelInternalApi.V1ModelIdLabelPut(Interal, id, resetLabelBody)
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	// 获取修改后的模型label信息
	modelRes, r, err := ApiInteral.ModelInternalApi.V1ModelIdGet(Interal, id)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	modelData := getData(s.T(), modelRes.Data)
	// 非法的Frameworks与Task不会修改成功
	assert.Equal(s.T(),
		map[string]interface{}{"frameworks": []interface{}{}, "license": "apache-2.0", "others": []interface{}{}, "task": ""},
		modelData["labels"])

	// 使用部分合法的字段修改模型label
	resetLabelBody = swaggerInternal.ControllerReqToResetLabel{
		Frameworks: []string{"PyTorch123", "MindSpore"},
		License:    "apache-2.0",
		Task:       "copa123",
	}
	_, r, err = ApiInteral.ModelInternalApi.V1ModelIdLabelPut(Interal, id, resetLabelBody)
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	// 获取修改后的模型label信息
	modelRes, r, err = ApiInteral.ModelInternalApi.V1ModelIdGet(Interal, id)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	modelData = getData(s.T(), modelRes.Data)
	// 只有合法的Frameworks会修改成功
	assert.Equal(s.T(),
		map[string]interface{}{"frameworks": []interface{}{"MindSpore"}, "license": "apache-2.0", "others": []interface{}{}, "task": ""},
		modelData["labels"])

	// 删除模型
	r, err = ApiRest.ModelApi.V1ModelIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestInternalModel used for testing
func TestInternalModel(t *testing.T) {
	suite.Run(t, new(SuiteInternalModel))
}
