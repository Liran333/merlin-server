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

	swagger "e2e/client"
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
	data, r, err := Api.SpaceApi.V1SpacePost(Auth2, swagger.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		InitReadme: false,
		License:    "mit",
		Name:       "testspace",
		Owner:      "test2",
		Sdk:        "gradio",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	_, r, err = Api.SpaceApi.V1SpaceIdPut(Auth2, id, swagger.ControllerReqToUpdateSpace{
		Desc:     "space desc new",
		Fullname: "spacefullname-new",
		Hardware: "NPU basic 8 vCPU · 32GB · FREE",
		Name:     "testspace-new",
	})

	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = Api.SpaceApi.V1SpaceIdDelete(Auth2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 使用无效仓库名创建、修改Space失败
func (s *SuiteUserSpace) TestUserCreateUpdateInvalidSpace() {
	_, r, err := Api.SpaceApi.V1SpacePost(Auth2, swagger.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		InitReadme: false,
		License:    "mit",
		Name:       "invalid#testspace",
		Owner:      "test2",
		Sdk:        "gradio",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	data, r, err := Api.SpaceApi.V1SpacePost(Auth2, swagger.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		InitReadme: false,
		License:    "mit",
		Name:       "testspace",
		Owner:      "test2",
		Sdk:        "gradio",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	_, r, err = Api.SpaceApi.V1SpaceIdPut(Auth2, id, swagger.ControllerReqToUpdateSpace{
		Desc:     "space desc new",
		Fullname: "spacefullname-new",
		Hardware: "NPU basic 8 vCPU · 32GB · FREE",
		Name:     "invalid#testspace",
	})

	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	r, err = Api.SpaceApi.V1SpaceIdDelete(Auth2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestNotLoginCantCreateSpace used for testing
// 没登录用户不能创建Space
func (s *SuiteUserSpace) TestNotLoginCantCreateSpace() {
	_, r, err := Api.SpaceApi.V1SpacePost(context.Background(), swagger.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		InitReadme: false,
		License:    "mit",
		Name:       "testspace",
		Owner:      "test2",
		Sdk:        "gradio",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusUnauthorized, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// 以下用例结果异常，需排查
// 可以访问自己名下的公有Space
// func (s *SuiteUserSpace) TestUserCanVisitSelfPublicSpace() {
//	 data, r, err := Api.SpaceApi.V1SpacePost(Auth2, swagger.ControllerReqToCreateSpace{
//		 Desc:       "space desc",
//		 Fullname:   "spacefullname",
//		 Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
//		 InitReadme: false,
//		 License:    "mit",
//		 Name:       "testspace",
//		 Owner:      "test2",
//		 Sdk:        "gradio",
//		 Visibility: "public",
//	 })
//
//	 assert.Equal(s.T(), 201, r.StatusCode)
//	 assert.Nil(s.T(), err)
//
//	 id := getString(s.T(), data.Data)
//
//	 detail, r, err := Api.SpaceWebApi.V1SpaceOwnerNameGet(Auth2, "test2", "testspace")
//	 assert.Equal(s.T(), 200, r.StatusCode)
//	 assert.Nil(s.T(), err)
//	 assert.NotEmpty(s.T(), detail.Name)
//
//	 spaceOwnerList, r, err := Api.SpaceWebApi.V1SpaceOwnerGet(Auth2, "test2", &swagger.SpaceWebApiV1SpaceOwnerGetOpts{})
//	 assert.Equal(s.T(), 200, r.StatusCode)
//	 assert.Nil(s.T(), err)
//	 assert.NotEmpty(s.T(), spaceOwnerList.Spaces)
//
//	 spaceList, r, err := Api.SpaceWebApi.V1SpaceGet(Auth2, &swagger.SpaceWebApiV1SpaceGetOpts{})
//	 assert.Equal(s.T(), 200, r.StatusCode)
//	 assert.Nil(s.T(), err)
//	 assert.NotEmpty(s.T(), spaceList.Spaces)
//
//	 r, err = Api.SpaceApi.V1SpaceIdDelete(Auth2, id)
//	 assert.Equal(s.T(), 204, r.StatusCode)
//	 assert.Nil(s.T(), err)
// }

// TestUserSpace used for testing
func TestUserSpace(t *testing.T) {
	suite.Run(t, new(SuiteUserSpace))
}
