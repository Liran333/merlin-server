/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package e2e

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	swaggerInternal "e2e/client_internal"
)

// SuiteActivityInternal used for testing
type SuiteActivityInternal struct {
	suite.Suite
}

// SetupSuite used for testing
func (s *SuiteActivityInternal) SetupSuite() {
}

// TestActivityInternal used for testing
func (s *SuiteActivityInternal) TestActivityInternal() {
	// Time非法，新增Activity失败
	_, r, err := ApiInteral.ActivityInternalApi.V1ActivityPost(Interal, swaggerInternal.ActivityappReqToCreateActivity{
		Owner:         "testInternal",
		ResourceIndex: "01",
		ResourceType:  "model",
		Time:          "2024-03-26",
		Type_:         "type",
	})
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	// ResourceIndex非法，新增Activity失败
	_, r, err = ApiInteral.ActivityInternalApi.V1ActivityPost(Interal, swaggerInternal.ActivityappReqToCreateActivity{
		Owner:         "testInternal",
		ResourceIndex: "badindex",
		ResourceType:  "model",
		Time:          "1711436086",
		Type_:         "type",
	})
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	// 新增Activity成功
	_, r, err = ApiInteral.ActivityInternalApi.V1ActivityPost(Interal, swaggerInternal.ActivityappReqToCreateActivity{
		Owner:         "testInternal",
		ResourceIndex: "01",
		ResourceType:  "model",
		Time:          "1711436086",
		Type_:         "type",
	})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	// 删除Activity成功
	_, r, err = ApiInteral.ActivityInternalApi.V1ActivityDelete(Interal, swaggerInternal.ActivityappReqToDeleteActivity{
		ResourceIndex: "01",
		ResourceType:  "model",
	})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	// 重复删除Activity返回成功
	_, r, err = ApiInteral.ActivityInternalApi.V1ActivityDelete(Interal, swaggerInternal.ActivityappReqToDeleteActivity{
		ResourceIndex: "01",
		ResourceType:  "model",
	})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestActivityInternal used for testing
func TestActivityInternal(t *testing.T) {
	suite.Run(t, new(SuiteActivityInternal))
}
