/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package e2e

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	swaggerInternal "e2e/client_internal"
	swaggerRest "e2e/client_rest"
)

// SuiteUserSpace used for testing
type SuiteUserSpace struct {
	suite.Suite
}

// SetupSuite used for testing
func (s *SuiteUserSpace) SetupSuite() {
}

// TearDownSuite used for testing
func (s *SuiteUserSpace) TearDownSuite() {
}

// TestUserCanCreateUpdateDeleteSpace used for testing
// 可以创建Space到自己名下, 并且可以修改和删除自己名下的Space
func (s *SuiteUserSpace) TestUserCanCreateUpdateDeleteSpace() {
	spaceParam := swaggerRest.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		License:    "mit",
		Name:       "testspace",
		Owner:      "test2",
		Sdk:        "gradio",
		Visibility: "public",
	}
	data, r, err := ApiRest.SpaceApi.V1SpacePost(AuthRest2, spaceParam)

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	// 重复创建空间返回400
	_, r, err = ApiRest.SpaceApi.V1SpacePost(AuthRest2, spaceParam)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	_, r, err = ApiRest.SpaceApi.V1SpaceIdPut(AuthRest2, id, swaggerRest.ControllerReqToUpdateSpace{
		Desc:     "space desc new",
		Fullname: "spacefullname-new",
		Hardware: "NPU basic 8 vCPU · 32GB · FREE",
	})

	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 使用无效仓库名创建Space失败
func (s *SuiteUserSpace) TestUserCreateUpdateInvalidSpace() {
	_, r, err := ApiRest.SpaceApi.V1SpacePost(AuthRest2, swaggerRest.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		License:    "mit",
		Name:       "invalid#testspace",
		Owner:      "test2",
		Sdk:        "gradio",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

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

	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestNotLoginCantCreateSpace used for testing
// 没登录用户不能创建Space
func (s *SuiteUserSpace) TestNotLoginCantCreateSpace() {
	_, r, err := ApiRest.SpaceApi.V1SpacePost(context.Background(), swaggerRest.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		License:    "mit",
		Name:       "testspace",
		Owner:      "test2",
		Sdk:        "gradio",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusUnauthorized, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestUserCanVisitSelfPublicSpace used for testing
// 可以访问自己名下的公有Space
func (s *SuiteUserSpace) TestUserCanVisitSelfPublicSpace() {
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

	detail, r, err := ApiRest.SpaceRestfulApi.V1SpaceOwnerNameGet(AuthRest2, "test2", "testspace")
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	space := getData(s.T(), detail.Data)
	assert.Equal(s.T(), "testspace", space["name"])

	// 查询test2名下的所有space
	list, r, err := ApiRest.SpaceRestfulApi.V1SpaceGet(AuthRest2, "test2", &swaggerRest.SpaceRestfulApiV1SpaceGetOpts{})
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	count := 0
	spaceLists := getData(s.T(), list.Data)

	for i := 0; i < len(spaceLists["spaces"].([]interface{})); i++ {
		model, ok := spaceLists["spaces"].([]interface{})[i].(map[string]interface{})
		if !ok {
			continue
		}

		if _, ok := model["name"]; ok {
			count++
		}
	}
	assert.Equal(s.T(), countOne, count)

	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestCreateSpace
// 创建space 成功，并成功能查询到各种参数
func (s *SuiteUserSpace) TestCreateSpace() {
	data, r, err := ApiRest.SpaceApi.V1SpacePost(AuthRest, swaggerRest.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		License:    "mit",
		Name:       "testspace",
		Owner:      "test1",
		Sdk:        "gradio",
		Visibility: "public",
		AvatarId:   "https://gitee.com/1",
	})

	id := getString(s.T(), data.Data)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	data, r, err = ApiRest.SpaceRestfulApi.V1SpaceOwnerNameGet(AuthRest, "test1", "testspace")
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	space := getData(s.T(), data.Data)

	assert.Equal(s.T(), "space desc", space["desc"])
	assert.Equal(s.T(), "spacefullname", space["fullname"])
	assert.Equal(s.T(), strings.ToLower("CPU basic 2 vCPU · 16GB · FREE"), space["hardware"])
	assert.Equal(s.T(),
		map[string]interface{}(map[string]interface{}{"frameworks": []interface{}{}, "license": "mit", "others": []interface{}{}, "task": ""}),
		space["labels"])
	assert.Equal(s.T(), "testspace", space["name"])
	assert.Equal(s.T(), "test1", space["owner"])
	assert.Equal(s.T(), "gradio", space["sdk"])
	assert.Equal(s.T(), "public", space["visibility"])
	assert.Equal(s.T(), "https://gitee.com/1", space["space_avatar_id"])
	assert.Equal(s.T(), "", space["avatar_id"])

	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest, id)

	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestUserSetSpaceDownloadCount used for testing
// 可以通过内部接口设置下载统计
func (s *SuiteUserModel) TestUserSetSpaceDownloadCount() {
	data, r, err := ApiRest.SpaceApi.V1SpacePost(AuthRest, swaggerRest.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		License:    "mit",
		Name:       "testspace",
		Owner:      "test1",
		Sdk:        "gradio",
		Visibility: "public",
		AvatarId:   "https://gitee.com/1",
	})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	// 重复创建模型返回400
	_, r, err = ApiInteral.StatisticApi.V1CoderepoIdStatisticPut(Interal, id, swaggerInternal.ControllerRepoStatistics{
		DownloadCount: 15,
	})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	data, r, err = ApiRest.SpaceRestfulApi.V1SpaceOwnerNameGet(AuthRest, "test1", "testspace")
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	model := getData(s.T(), data.Data)
	assert.Equal(s.T(), float64(15), model["download_count"])

	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest, id)

	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestUserSpace used for testing
func TestUserSpace(t *testing.T) {
	suite.Run(t, new(SuiteUserSpace))
}
