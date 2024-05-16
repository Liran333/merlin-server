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

// SuiteComputilityInternal used for testing
type SuiteComputilityInternal struct {
	suite.Suite
}

// SetupSuite used for testing
func (s *SuiteComputilityInternal) SetupSuite() {
}

// TearDownSuite used for testing
func (s *SuiteComputilityInternal) TearDownSuite() {
}

// 用户加入、退出算力组织，删除算力组织
func (s *SuiteComputilityInternal) TestComputilityOrgDelete() {
	// 加入算力组织
	_, r, err := ApiInteral.ComputilityInternalApi.V1ComputilityAccountPost(Interal, swaggerInternal.ControllerReqToUserOrgOperate{
		OrgName:  "test-npu-2",
		UserName: "test2",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	// 加入非算力组织
	_, r, err = ApiInteral.ComputilityInternalApi.V1ComputilityAccountPost(Interal, swaggerInternal.ControllerReqToUserOrgOperate{
		OrgName:  "normal",
		UserName: "test2",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	// 退出算力组织
	_, r, err = ApiInteral.ComputilityInternalApi.V1ComputilityAccountRemovePut(Interal, swaggerInternal.ControllerReqToUserOrgOperate{
		OrgName:  "test-npu-2",
		UserName: "test2",
	})

	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	// 退出非算力组织
	_, r, err = ApiInteral.ComputilityInternalApi.V1ComputilityAccountRemovePut(Interal, swaggerInternal.ControllerReqToUserOrgOperate{
		OrgName:  "normal",
		UserName: "test2",
	})

	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	// 删除算力组织
	_, r, err = ApiInteral.ComputilityInternalApi.V1ComputilityOrgDeletePost(Interal, swaggerInternal.ControllerReqToOrgDelete{
		OrgName: "test-npu-2",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	// 删除非算力组织
	_, r, err = ApiInteral.ComputilityInternalApi.V1ComputilityOrgDeletePost(Interal, swaggerInternal.ControllerReqToOrgDelete{
		OrgName: "normal",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestComputilityInternal used for testing
func TestComputilityInternal(t *testing.T) {
	suite.Run(t, new(SuiteComputilityInternal))
}
