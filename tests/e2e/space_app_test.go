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

// SuiteUserSpace used for testing
type SuiteSpaceAppRestful struct {
	suite.Suite
}

// SetupSuite used for testing
func (s *SuiteSpaceAppRestful) SetupSuite() {
}

// TearDownSuite used for testing
func (s *SuiteSpaceAppRestful) TearDownSuite() {
}

// 未创建space-app或者状态不为ready，重启失败；否则，重启成功
func (s *SuiteSpaceAppRestful) TestSpaceAppRestart() {
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

	// 未创建space-app，重启失败
	_, r, err = ApiRest.SpaceAppRestfulApi.V1SpaceAppOwnerNameRestartPost(AuthRest2, "test2", "testspace")
	assert.Equal(s.T(), http.StatusNotFound, r.StatusCode)
	assert.NotNil(s.T(), err)

	// 创建space-app
	_, r, err = ApiInteral.SpaceAppApi.V1SpaceAppPost(Interal, swaggerInternal.ControllerReqToCreateSpaceApp{
		SpaceId:  id,
		CommitId: "12345",
	})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	// space-app状态不为ready，重启失败
	_, r, err = ApiRest.SpaceAppRestfulApi.V1SpaceAppOwnerNameRestartPost(AuthRest2, "test2", "testspace")
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	_, r, err = ApiInteral.SpaceAppApi.V1SpaceAppBuildStartedPut(Interal, swaggerInternal.ControllerReqToUpdateBuildInfo{
		SpaceId:  id,
		CommitId: "12345",
		LogUrl:   "https://www.modelfoundry.cn",
	})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	_, r, err = ApiInteral.SpaceAppApi.V1SpaceAppBuildDonePut(Interal, swaggerInternal.ControllerReqToSetBuildIsDone{
		SpaceId:  id,
		CommitId: "12345",
		Success:  true,
		Logs:     "ready",
	})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	_, r, err = ApiInteral.SpaceAppApi.V1SpaceAppServiceStartedPut(Interal,
		swaggerInternal.ControllerReqToUpdateServiceInfo{
			SpaceId:  id,
			CommitId: "12345",
			AppUrl:   "https://www.modelfoundry.cn",
			LogUrl:   "https://www.modelfoundry.cn",
		})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	// 重启成功
	_, r, err = ApiRest.SpaceAppRestfulApi.V1SpaceAppOwnerNameRestartPost(AuthRest2, "test2", "testspace")
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestUserSpace used for testing
func TestSpaceAppRest(t *testing.T) {
	suite.Run(t, new(SuiteSpaceAppRestful))
}
