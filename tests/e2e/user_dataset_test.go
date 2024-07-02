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

	swaggerInternal "e2e/client_internal"
	swaggerRest "e2e/client_rest"
)

// SuiteUserDataset used for testing
type SuiteUserDataset struct {
	suite.Suite
}

// SetupSuite used for testing
func (s *SuiteUserDataset) SetupSuite() {
}

// TearDownSuite used for testing
func (s *SuiteUserDataset) TearDownSuite() {
}

// TestUserCanCreateUpdateDeleteDataset used for testing
// 可以创建数据集到自己名下, 并且可以修改和删除自己名下的数据集
func (s *SuiteUserDataset) TestUserCanCreateUpdateDeleteDataset() {
	datasetParam := swaggerRest.ControllerReqToCreateDataset{
		Name:       "testdataset",
		Owner:      "test2",
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

	_, r, err = ApiRest.DatasetApi.V1DatasetIdPut(AuthRest2, id, swaggerRest.ControllerReqToUpdateDataset{
		Visibility: "public",
	})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.DatasetApi.V1DatasetIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 使用无效仓库名创建数据集失败
func (s *SuiteUserDataset) TestUserCreateUpdateInvalidDataset() {
	_, r, err := ApiRest.DatasetApi.V1DatasetPost(AuthRest2, swaggerRest.ControllerReqToCreateDataset{
		Name:       "invalid#testdataset",
		Owner:      "test2",
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	data, r, err := ApiRest.DatasetApi.V1DatasetPost(AuthRest2, swaggerRest.ControllerReqToCreateDataset{
		Name:       "testdataset",
		Owner:      "test2",
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	r, err = ApiRest.DatasetApi.V1DatasetIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestNotLoginCantCreateDataset used for testing
// 没登录用户不能创建数据集
func (s *SuiteUserDataset) TestNotLoginCantCreateDataset() {
	_, r, err := ApiRest.DatasetApi.V1DatasetPost(context.Background(), swaggerRest.ControllerReqToCreateDataset{
		Name:       "testdataset",
		Owner:      "test2",
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusUnauthorized, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestUserCanVisitSelfPublicDataset used for testing
// 可以访问自己名下的公有数据集
func (s *SuiteUserDataset) TestUserCanVisitSelfPublicDataset() {
	// 创建testdataset
	data, r, err := ApiRest.DatasetApi.V1DatasetPost(AuthRest2, swaggerRest.ControllerReqToCreateDataset{
		Name:       "testdataset",
		Owner:      "test2",
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	// 获取test2名下的dataset成功
	detail, r, err := ApiRest.DatasetRestfulApi.V1DatasetOwnerNameGet(AuthRest2, "test2", "testdataset")
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)
	dataset := getData(s.T(), detail.Data)
	assert.Equal(s.T(), "testdataset", dataset["name"])

	// 创建testdataset2
	data, r, err = ApiRest.DatasetApi.V1DatasetPost(AuthRest2, swaggerRest.ControllerReqToCreateDataset{
		Name:       "testdataset2",
		Owner:      "test2",
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id2 := getString(s.T(), data.Data)

	// 获取用户名下的所有dataset list
	list, r, err := ApiRest.DatasetRestfulApi.V1DatasetGet(AuthRest2, "test2", &swaggerRest.DatasetRestfulApiV1DatasetGetOpts{})
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	count := 0
	datasetLists := list.Data

	for i := 0; i < len(datasetLists.Datasets); i++ {
		dataset := datasetLists.Datasets[i]

		if dataset.Name != "" {
			count++
		}
	}
	assert.Equal(s.T(), countTwo, count)

	// 删除创建的数据集
	r, err = ApiRest.DatasetApi.V1DatasetIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.DatasetApi.V1DatasetIdDelete(AuthRest2, id2)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestUserSetDatasetDownloadCount used for testing
// 可以通过内部接口设置下载统计
func (s *SuiteUserDataset) TestUserSetDatasetDownloadCount() {
	datasetParam := swaggerRest.ControllerReqToCreateDataset{
		Name:       "testdataset",
		Owner:      "test2",
		License:    "mit",
		Visibility: "public",
	}
	data, r, err := ApiRest.DatasetApi.V1DatasetPost(AuthRest2, datasetParam)

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	// 重复创建数据集返回400
	_, r, err = ApiInteral.CodeRepoInternalApi.V1CoderepoIdStatisticDownloadPut(
		Interal, id, swaggerInternal.ControllerRepoStatistics{
			DownloadCount: 10,
		})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	data1, r, err := ApiRest.DatasetRestfulApi.V1DatasetOwnerNameGet(AuthRest2, "test2", "testdataset")
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	dataset := getData(s.T(), data1.Data)
	assert.Equal(s.T(), int32(10), dataset["download_count"])

	r, err = ApiRest.DatasetApi.V1DatasetIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestDatasetOwnerNameGetPrivateDataset used for testing
// 可以访问自己名下的私有数据集
func (s *SuiteUserDataset) TestDatasetOwnerNameGetPrivateDataset() {
	// 创建testdataset
	data, r, err := ApiRest.DatasetApi.V1DatasetPost(AuthRest2, swaggerRest.ControllerReqToCreateDataset{
		Name:       "testdataset3",
		Owner:      "test2",
		License:    "mit",
		Visibility: "private",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	// 获取test2名下的指定dataset成功
	detail, r, err := ApiRest.DatasetRestfulApi.V1DatasetOwnerNameGet(AuthRest2, "test2", "testdataset3")
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)
	dataset := getData(s.T(), detail.Data)
	assert.Equal(s.T(), "testdataset3", dataset["name"])

	r, err = ApiRest.DatasetApi.V1DatasetIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestDatasetMaxPerUser used for testing
// 单个用户可以创建最多数据集个数
func (s *SuiteUserDataset) TestDatasetMaxPerUser() {
	// 创建testdataset
	data1, r, err := ApiRest.DatasetApi.V1DatasetPost(AuthRest2, swaggerRest.ControllerReqToCreateDataset{
		Name:       "testdataset1",
		Owner:      "test2",
		License:    "mit",
		Visibility: "private",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	id1 := getString(s.T(), data1.Data)

	data2, r, err := ApiRest.DatasetApi.V1DatasetPost(AuthRest2, swaggerRest.ControllerReqToCreateDataset{
		Name:       "testdataset2",
		Owner:      "test2",
		License:    "mit",
		Visibility: "private",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	id2 := getString(s.T(), data2.Data)

	data3, r, err := ApiRest.DatasetApi.V1DatasetPost(AuthRest2, swaggerRest.ControllerReqToCreateDataset{
		Name:       "testdataset3",
		Owner:      "test2",
		License:    "mit",
		Visibility: "private",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	id3 := getString(s.T(), data3.Data)

	data4, r, err := ApiRest.DatasetApi.V1DatasetPost(AuthRest2, swaggerRest.ControllerReqToCreateDataset{
		Name:       "testdataset4",
		Owner:      "test2",
		License:    "mit",
		Visibility: "private",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	id4 := getString(s.T(), data4.Data)

	// 创建达到最大允许个数，创建失败
	_, r, err = ApiRest.DatasetApi.V1DatasetPost(AuthRest2, swaggerRest.ControllerReqToCreateDataset{
		Name:       "testdataset5",
		Owner:      "test2",
		License:    "mit",
		Visibility: "private",
	})

	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	r, err = ApiRest.DatasetApi.V1DatasetIdDelete(AuthRest2, id1)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.DatasetApi.V1DatasetIdDelete(AuthRest2, id2)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.DatasetApi.V1DatasetIdDelete(AuthRest2, id3)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.DatasetApi.V1DatasetIdDelete(AuthRest2, id4)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestUserDataset used for testing
func TestUserDataset(t *testing.T) {
	suite.Run(t, new(SuiteUserDataset))
}
