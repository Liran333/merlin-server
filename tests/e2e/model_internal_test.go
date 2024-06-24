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
	resetLabelBody := swaggerInternal.GithubComOpenmerlinMerlinServerModelsControllerReqToResetLabel{
		Frameworks:  []string{"PyTorch"},
		Licenses:    []string{"apache-2.0"},
		Task:        "document-question-answering",
		Hardwares:   []string{"CPU"},
		Languages:   []string{"cn"},
		LibraryName: "openmind",
	}
	_, r, err = ApiInteral.ModelInternalApi.V1ModelIdLabelPut(Interal, id, resetLabelBody)
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	// 获取修改成功后的模型label信息
	modelRes, r, err := ApiInteral.ModelInternalApi.V1ModelIdGet(Interal, id)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	modelData := getData(s.T(), modelRes.Data)
	labels := getData(s.T(), modelData["labels"])
	assert.Equal(s.T(),
		map[string]interface{}(map[string]interface{}{
			"frameworks": []string{"PyTorch"}, "library_name": "openmind", "hardwares": []string{"CPU"}, "language": []string{"cn"}, "license": []string{"apache-2.0"}, "others": []string{},
			"task": "document-question-answering"}),
		labels)

	// 修改模型label，其中frameworks,hardwares,languages有多个
	resetLabelBody = swaggerInternal.GithubComOpenmerlinMerlinServerModelsControllerReqToResetLabel{
		Frameworks: []string{"PyTorch", "MindSpore"},
		Hardwares:  []string{"CPU", "GPU"},
		Languages:  []string{"cn", "en"},
		Licenses:   []string{"apache-2.0"},
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
	labels = getData(s.T(), modelData["labels"])

	assert.Equal(s.T(), []string{"apache-2.0"}, labels["license"])
	assert.Equal(s.T(), "copa", labels["task"])
	assert.Contains(s.T(),
		[][]string{{"PyTorch", "MindSpore"}, {"MindSpore", "PyTorch"}}, getData(s.T(), labels)["frameworks"])
	assert.Contains(s.T(),
		[][]string{{"CPU", "GPU"}, {"GPU", "CPU"}}, getData(s.T(), labels)["hardwares"])
	assert.Contains(s.T(),
		[][]string{{"cn", "en"}, {"en", "cn"}}, getData(s.T(), labels)["language"])

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
	resetLabelBody := swaggerInternal.GithubComOpenmerlinMerlinServerModelsControllerReqToResetLabel{
		Frameworks:  []string{"PyTorch123"},
		Hardwares:   []string{"CPU123"},
		Licenses:    []string{"apache-2.0"},
		Task:        "copa123",
		LibraryName: "openmind123",
	}
	_, r, err = ApiInteral.ModelInternalApi.V1ModelIdLabelPut(Interal, id, resetLabelBody)
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	// 获取修改后的模型label信息
	modelRes, r, err := ApiInteral.ModelInternalApi.V1ModelIdGet(Interal, id)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	modelData := getData(s.T(), modelRes.Data)
	labels := getData(s.T(), modelData["labels"])
	// 非法的Frameworks、LibraryName、Task和Hardwares不会修改成功
	assert.Equal(s.T(),
		map[string]interface{}{"frameworks": []string{}, "library_name": "", "hardwares": []string{}, "language": []string{}, "license": []string{"apache-2.0"},
			"others": []string{}, "task": ""},
		labels)

	// 使用部分合法的字段修改模型label
	resetLabelBody = swaggerInternal.GithubComOpenmerlinMerlinServerModelsControllerReqToResetLabel{
		Frameworks:  []string{"PyTorch123", "MindSpore"},
		Hardwares:   []string{"CPU123", "GPU"},
		Licenses:    []string{"apache-2.0"},
		LibraryName: "openmind",
		Task:        "copa123",
	}
	_, r, err = ApiInteral.ModelInternalApi.V1ModelIdLabelPut(Interal, id, resetLabelBody)
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	// 获取修改后的模型label信息
	modelRes, r, err = ApiInteral.ModelInternalApi.V1ModelIdGet(Interal, id)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	modelData = getData(s.T(), modelRes.Data)
	labels = getData(s.T(), modelData["labels"])
	// 合法的Frameworks,Hardwares和LibraryName会修改成功
	assert.Equal(s.T(),
		map[string]interface{}{"frameworks": []string{"MindSpore"}, "hardwares": []string{"GPU"}, "language": []string{}, "library_name": "openmind", "license": []string{"apache-2.0"},
			"others": []string{}, "task": ""},
		labels)

	// 删除模型
	r, err = ApiRest.ModelApi.V1ModelIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestInternalModel used for testing
func TestInternalModel(t *testing.T) {
	suite.Run(t, new(SuiteInternalModel))
}
