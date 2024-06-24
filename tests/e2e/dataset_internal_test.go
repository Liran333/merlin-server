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

// SuiteInternalDataset used for testing
type SuiteInternalDataset struct {
	suite.Suite
}

// SetupSuite used for testing
func (s *SuiteInternalDataset) SetupSuite() {
}

// TearDownSuite used for testing
func (s *SuiteInternalDataset) TearDownSuite() {
}

// TestGetDataset used for testing
// 获取数据集成功
func (s *SuiteInternalDataset) TestGetDataset() {
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

	// 获取数据集
	datasetRes, r, err := ApiInteral.DatasetInternalApi.V1DatasetIdGet(Interal, id)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	datasetData := getData(s.T(), datasetRes.Data)
	assert.Equal(s.T(), id, datasetData["id"])
	assert.Equal(s.T(), "testdataset", datasetData["name"])
	assert.Equal(s.T(), "test2", datasetData["owner"])
	assert.Equal(s.T(), "public", datasetData["visibility"])

	// 删除数据集
	r, err = ApiRest.DatasetApi.V1DatasetIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestGetDatasetFailed used for testing
// 获取数据集失败
func (s *SuiteInternalDataset) TestGetDatasetFailed() {
	// 数据集不存在
	unExistedId := "0"
	_, r, err := ApiInteral.DatasetInternalApi.V1DatasetIdGet(Interal, unExistedId)
	assert.Equal(s.T(), http.StatusNotFound, r.StatusCode)
	assert.NotNil(s.T(), err)

	// 参数无效
	unExistedId = "test"
	_, r, err = ApiInteral.DatasetInternalApi.V1DatasetIdGet(Interal, unExistedId)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestInternalDatasetResetLabelSuccess used for testing
// 数据集label设置成功
func (s *SuiteInternalDataset) TestInternalDatasetResetLabelSuccess() {
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

	// 修改数据集label
	resetLabelBody := swaggerInternal.ControllerReqToResetDatasetLabel{
		Licenses: []string{"apache-2.0"},
		Task:     []string{"task1"},
		Size:     "n<1K",
		Language: []string{"Chinese"},
		Domain:   []string{"chemistry"},
	}
	_, r, err = ApiInteral.DatasetInternalApi.V1DatasetIdLabelPut(Interal, id, resetLabelBody)
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	// 获取修改成功后的数据集label信息
	datasetRes, r, err := ApiInteral.DatasetInternalApi.V1DatasetIdGet(Interal, id)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	datasetData := getData(s.T(), datasetRes.Data)
	labels := getData(s.T(), datasetData["labels"])
	assert.Equal(s.T(),
		map[string]interface{}(map[string]interface{}{
			"license":  []string{"apache-2.0"},
			"task":     []string{"task1"},
			"size":     "n<1K",
			"language": []string{"Chinese"},
			"domain":   []string{"chemistry"},
		}),
		labels)

	// 删除数据集
	r, err = ApiRest.DatasetApi.V1DatasetIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestInternalDataset used for testing
func TestInternalDataset(t *testing.T) {
	suite.Run(t, new(SuiteInternalDataset))
}
