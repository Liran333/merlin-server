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

// 未创建space-app或者状态不为serving，重启失败；否则，重启成功
func (s *SuiteSpaceAppRestful) TestSpaceAppRestart() {
	data, r, err := ApiRest.SpaceApi.V1SpacePost(AuthRest2, swaggerRest.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		License:    "mit",
		BaseImage:  "python3.8-pytorch2.1",
		Name:       "testspace",
		Owner:      "test2",
		Sdk:        "gradio",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)
	if id == "" {
		return
	}

	// 未创建space-app，重启失败
	_, r, err = ApiRest.SpaceApi.V1SpaceAppOwnerNameRestartPost(AuthRest2, "test2", "testspace")
	assert.Equal(s.T(), http.StatusNotFound, r.StatusCode)
	assert.NotNil(s.T(), err)

	// 更新commitId
	_, r, err = ApiInteral.SpaceInternalApi.V1SpaceIdNotifyUpdateCodePut(Interal, id, swaggerInternal.ControllerReqToNotifyUpdateCode{
		SdkType:  "gradio",
		CommitId: "12345",
	})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	// 创建space-app
	_, r, err = ApiInteral.SpaceAppApi.V1SpaceAppPost(Interal, swaggerInternal.ControllerReqToCreateSpaceApp{
		SpaceId:  id,
		CommitId: "12345",
	})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	// space-app状态不为serving，重启失败
	_, r, err = ApiRest.SpaceApi.V1SpaceAppOwnerNameRestartPost(AuthRest2, "test2", "testspace")
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	_, r, err = ApiInteral.SpaceAppApi.V1SpaceAppBuildingPut(Interal, swaggerInternal.ControllerReqToUpdateBuildInfo{
		SpaceId:  id,
		CommitId: "12345",
		LogUrl:   "https://www.modelfoundry.cn",
	})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	_, r, err = ApiInteral.SpaceAppApi.V1SpaceAppStartingPut(Interal, swaggerInternal.ControllerReqToNotifyStarting{
		SpaceId:     id,
		CommitId:    "12345",
		AllBuildLog: "vertex:  [internal] load build definition from Dockerfile\n",
	})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	_, r, err = ApiInteral.SpaceAppApi.V1SpaceAppServingPut(Interal,
		swaggerInternal.ControllerReqToUpdateServiceInfo{
			SpaceId:  id,
			CommitId: "12345",
			AppUrl:   "https://www.modelfoundry.cn",
			LogUrl:   "https://www.modelfoundry.cn",
		})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	// 重启成功
	_, r, err = ApiRest.SpaceApi.V1SpaceAppOwnerNameRestartPost(AuthRest2, "test2", "testspace")
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 未创建space-app或者状态不为serving，暂停失败；否则，暂停成功
func (s *SuiteSpaceAppRestful) TestSpaceAppPause() {
	data, r, err := ApiRest.SpaceApi.V1SpacePost(AuthRest2, swaggerRest.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		License:    "mit",
		BaseImage:  "python3.8-pytorch2.1",
		Name:       "testspace",
		Owner:      "test2",
		Sdk:        "gradio",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)
	if id == "" {
		return
	}

	// 未创建space-app，暂停失败
	_, r, err = ApiRest.SpaceApi.V1SpaceAppOwnerNamePausePost(AuthRest2, "test2", "testspace")
	assert.Equal(s.T(), http.StatusNotFound, r.StatusCode)
	assert.NotNil(s.T(), err)

	// 更新commitId
	_, r, err = ApiInteral.SpaceInternalApi.V1SpaceIdNotifyUpdateCodePut(Interal, id, swaggerInternal.ControllerReqToNotifyUpdateCode{
		SdkType:  "gradio",
		CommitId: "12345",
	})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	// 创建space-app
	_, r, err = ApiInteral.SpaceAppApi.V1SpaceAppPost(Interal, swaggerInternal.ControllerReqToCreateSpaceApp{
		SpaceId:  id,
		CommitId: "12345",
	})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	// space-app状态不为serving，暂停失败
	_, r, err = ApiRest.SpaceApi.V1SpaceAppOwnerNamePausePost(AuthRest2, "test2", "testspace")
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	_, r, err = ApiInteral.SpaceAppApi.V1SpaceAppBuildingPut(Interal, swaggerInternal.ControllerReqToUpdateBuildInfo{
		SpaceId:  id,
		CommitId: "12345",
		LogUrl:   "https://www.modelfoundry.cn",
	})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	_, r, err = ApiInteral.SpaceAppApi.V1SpaceAppStartingPut(Interal, swaggerInternal.ControllerReqToNotifyStarting{
		SpaceId:     id,
		CommitId:    "12345",
		AllBuildLog: "vertex:  [internal] load build definition from Dockerfile\n",
	})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	_, r, err = ApiInteral.SpaceAppApi.V1SpaceAppServingPut(Interal,
		swaggerInternal.ControllerReqToUpdateServiceInfo{
			SpaceId:  id,
			CommitId: "12345",
			AppUrl:   "https://www.modelfoundry.cn",
			LogUrl:   "https://www.modelfoundry.cn",
		})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	// 暂停成功
	_, r, err = ApiRest.SpaceApi.V1SpaceAppOwnerNamePausePost(AuthRest2, "test2", "testspace")
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	// 恢复成功
	_, r, err = ApiRest.SpaceApi.V1SpaceAppOwnerNameResumePost(AuthRest2, "test2", "testspace")
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 未创建space-app或者状态不为serving，暂停失败；否则，暂停成功
func (s *SuiteSpaceAppRestful) TestSpaceAppNpuPause() {
	// 加入算力组织
	_, r, err := ApiInteral.ComputilityInternalApi.V1ComputilityAccountPost(Interal, swaggerInternal.ControllerReqToUserOrgOperate{
		OrgName:  "test-npu",
		UserName: "test2",
	})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	data, r, err := ApiRest.SpaceApi.V1SpacePost(AuthRest2, swaggerRest.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "NPU basic 8 vCPU · 32GB · FREE",
		License:    "mit",
		BaseImage:  "python3.8-cann8.0-pytorch2.1",
		Name:       "testspace",
		Owner:      "test2",
		Sdk:        "gradio",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)
	if id == "" {
		// 退出算力组织
		_, r, err = ApiInteral.ComputilityInternalApi.V1ComputilityAccountPost(Interal, swaggerInternal.ControllerReqToUserOrgOperate{
			OrgName:  "test-npu",
			UserName: "test2",
		})

		assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
		assert.Nil(s.T(), err)
		return
	}

	// 未创建space-app，暂停失败
	_, r, err = ApiRest.SpaceApi.V1SpaceAppOwnerNamePausePost(AuthRest2, "test2", "testspace")
	assert.Equal(s.T(), http.StatusNotFound, r.StatusCode)
	assert.NotNil(s.T(), err)

	// 更新commitId
	_, r, err = ApiInteral.SpaceInternalApi.V1SpaceIdNotifyUpdateCodePut(Interal, id, swaggerInternal.ControllerReqToNotifyUpdateCode{
		SdkType:  "gradio",
		CommitId: "12345",
	})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	// 创建space-app
	_, r, err = ApiInteral.SpaceAppApi.V1SpaceAppPost(Interal, swaggerInternal.ControllerReqToCreateSpaceApp{
		SpaceId:  id,
		CommitId: "12345",
	})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	// space-app状态不为serving，暂停失败
	_, r, err = ApiRest.SpaceApi.V1SpaceAppOwnerNamePausePost(AuthRest2, "test2", "testspace")
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	_, r, err = ApiInteral.SpaceAppApi.V1SpaceAppBuildingPut(Interal, swaggerInternal.ControllerReqToUpdateBuildInfo{
		SpaceId:  id,
		CommitId: "12345",
		LogUrl:   "https://www.modelfoundry.cn",
	})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	_, r, err = ApiInteral.SpaceAppApi.V1SpaceAppStartingPut(Interal, swaggerInternal.ControllerReqToNotifyStarting{
		SpaceId:     id,
		CommitId:    "12345",
		AllBuildLog: "vertex:  [internal] load build definition from Dockerfile\n",
	})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	_, r, err = ApiInteral.SpaceAppApi.V1SpaceAppServingPut(Interal,
		swaggerInternal.ControllerReqToUpdateServiceInfo{
			SpaceId:  id,
			CommitId: "12345",
			AppUrl:   "https://www.modelfoundry.cn",
			LogUrl:   "https://www.modelfoundry.cn",
		})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	// 暂停成功
	_, r, err = ApiRest.SpaceApi.V1SpaceAppOwnerNamePausePost(AuthRest2, "test2", "testspace")
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	// 恢复成功
	_, r, err = ApiRest.SpaceApi.V1SpaceAppOwnerNameResumePost(AuthRest2, "test2", "testspace")
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	// 退出算力组织
	_, r, err = ApiInteral.ComputilityInternalApi.V1ComputilityAccountPost(Interal, swaggerInternal.ControllerReqToUserOrgOperate{
		OrgName:  "test-npu",
		UserName: "test2",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestUserSpace used for testing
func TestSpaceAppRest(t *testing.T) {
	suite.Run(t, new(SuiteSpaceAppRestful))
}
