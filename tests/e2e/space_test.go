/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package e2e

import (
	"net/http"
	"strings"
	"testing"

	"github.com/antihax/optional"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	swaggerInternal "e2e/client_internal"
	swaggerRest "e2e/client_rest"
)

// SuiteUserSpace used for testing
type SuiteSpace struct {
	suite.Suite
	Name           string
	Name1          string
	Name2          string
	Name3          string
	CPU            string
	NPU            string
	BaseImageMSCPU string
	BaseImageMSNPU string
	BaseImagePTNPU string
	BaseImagePTCPU string
	Owner          string
	Visibility     string
	Desc           string
	License        string
	Fullname       string
	Sdk            string
}

// SetupSuite used for testing
func (s *SuiteSpace) SetupSuite() {
	s.Name = "testspace"
	s.Name1 = "testspace1"
	s.Name2 = "testspace2"
	s.Name3 = "testspace3"
	s.CPU = "CPU basic 2 vCPU · 16GB · FREE"
	s.NPU = "NPU basic 8 vCPU · 32GB · FREE"
	s.BaseImagePTCPU = "python3.8-pytorch2.1"
	s.BaseImageMSCPU = "python3.8-mindspore2.3"
	s.BaseImagePTNPU = "python3.8-cann8.0-pytorch2.1"
	s.BaseImageMSNPU = "python3.8-cann8.0-mindspore2.3"
	s.Owner = "test1"
	s.Visibility = "public"
	s.Desc = "space desc"
	s.License = "mit"
	s.Fullname = "spacefullname"
	s.Sdk = "gradio"
}

// TearDownSuite used for testing
func (s *SuiteSpace) TearDownSuite() {
}

// 正常创建space
// base image mindspore cpu
func (s *SuiteSpace) TestSpaceCreateMSCPU() {
	data, r, err := ApiRest.SpaceApi.V1SpacePost(AuthRest, swaggerRest.ControllerReqToCreateSpace{
		Desc:       s.Desc,
		Fullname:   s.Fullname,
		Hardware:   s.CPU,
		License:    s.License,
		Name:       s.Name,
		Owner:      s.Owner,
		Sdk:        s.Sdk,
		Visibility: s.Visibility,
		BaseImage:  s.BaseImageMSCPU,
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := data.Data
	assert.NotEqual(s.T(), "", id)

	space, r, err := ApiRest.SpaceRestfulApi.V1SpaceOwnerNameGet(AuthRest, s.Owner, s.Name)

	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), s.Name, space.Data.Name)
	assert.Equal(s.T(), s.Owner, space.Data.Owner)
	assert.Equal(s.T(), s.Visibility, space.Data.Visibility)
	assert.Equal(s.T(), s.Sdk, space.Data.Sdk)
	assert.Equal(s.T(), s.Fullname, space.Data.Fullname)
	assert.Equal(s.T(), strings.ToLower(s.CPU), space.Data.Hardware)
	assert.Equal(s.T(), s.License, space.Data.Labels.License)
	assert.Equal(s.T(), s.Desc, space.Data.Desc)
	assert.Equal(s.T(), s.Visibility, space.Data.Visibility)
	assert.Equal(s.T(), "mindspore", space.Data.Labels.Framework)

	// 删除space
	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 正常创建space
// base image pytorch cpu
func (s *SuiteSpace) TestSpaceCreatePTCPU() {
	data, r, err := ApiRest.SpaceApi.V1SpacePost(AuthRest, swaggerRest.ControllerReqToCreateSpace{
		Desc:       s.Desc,
		Fullname:   s.Fullname,
		Hardware:   s.CPU,
		License:    s.License,
		Name:       s.Name,
		Owner:      s.Owner,
		Sdk:        s.Sdk,
		Visibility: s.Visibility,
		BaseImage:  s.BaseImagePTCPU,
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := data.Data
	assert.NotEqual(s.T(), "", id)

	space, r, err := ApiRest.SpaceRestfulApi.V1SpaceOwnerNameGet(AuthRest, s.Owner, s.Name)

	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), s.Name, space.Data.Name)
	assert.Equal(s.T(), s.Owner, space.Data.Owner)
	assert.Equal(s.T(), s.Visibility, space.Data.Visibility)
	assert.Equal(s.T(), s.Sdk, space.Data.Sdk)
	assert.Equal(s.T(), s.Fullname, space.Data.Fullname)
	assert.Equal(s.T(), strings.ToLower(s.CPU), space.Data.Hardware)
	assert.Equal(s.T(), s.License, space.Data.Labels.License)
	assert.Equal(s.T(), s.Desc, space.Data.Desc)
	assert.Equal(s.T(), s.Visibility, space.Data.Visibility)
	assert.Equal(s.T(), "pytorch", space.Data.Labels.Framework)

	// 删除space
	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 正常创建space
// base image pytorch npu
func (s *SuiteSpace) TestSpaceCreatePTNPU() {
	data, r, err := ApiRest.SpaceApi.V1SpacePost(AuthRest, swaggerRest.ControllerReqToCreateSpace{
		Desc:       s.Desc,
		Fullname:   s.Fullname,
		Hardware:   s.NPU,
		License:    s.License,
		Name:       s.Name,
		Owner:      s.Owner,
		Sdk:        s.Sdk,
		Visibility: s.Visibility,
		BaseImage:  s.BaseImagePTNPU,
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := data.Data
	assert.NotEqual(s.T(), "", id)

	space, r, err := ApiRest.SpaceRestfulApi.V1SpaceOwnerNameGet(AuthRest, s.Owner, s.Name)

	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), s.Name, space.Data.Name)
	assert.Equal(s.T(), s.Owner, space.Data.Owner)
	assert.Equal(s.T(), s.Visibility, space.Data.Visibility)
	assert.Equal(s.T(), s.Sdk, space.Data.Sdk)
	assert.Equal(s.T(), s.Fullname, space.Data.Fullname)
	assert.Equal(s.T(), strings.ToLower(s.NPU), space.Data.Hardware)
	assert.Equal(s.T(), s.License, space.Data.Labels.License)
	assert.Equal(s.T(), s.Desc, space.Data.Desc)
	assert.Equal(s.T(), s.Visibility, space.Data.Visibility)
	assert.Equal(s.T(), "pytorch", space.Data.Labels.Framework)

	// 删除space
	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 正常创建space
// base image mindspore npu
func (s *SuiteSpace) TestSpaceCreateMSNPU() {
	data, r, err := ApiRest.SpaceApi.V1SpacePost(AuthRest, swaggerRest.ControllerReqToCreateSpace{
		Desc:       s.Desc,
		Fullname:   s.Fullname,
		Hardware:   s.NPU,
		License:    s.License,
		Name:       s.Name,
		Owner:      s.Owner,
		Sdk:        s.Sdk,
		Visibility: s.Visibility,
		BaseImage:  s.BaseImageMSNPU,
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := data.Data
	assert.NotEqual(s.T(), "", id)

	space, r, err := ApiRest.SpaceRestfulApi.V1SpaceOwnerNameGet(AuthRest, s.Owner, s.Name)

	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), s.Name, space.Data.Name)
	assert.Equal(s.T(), s.Owner, space.Data.Owner)
	assert.Equal(s.T(), s.Visibility, space.Data.Visibility)
	assert.Equal(s.T(), s.Sdk, space.Data.Sdk)
	assert.Equal(s.T(), s.Fullname, space.Data.Fullname)
	assert.Equal(s.T(), strings.ToLower(s.NPU), space.Data.Hardware)
	assert.Equal(s.T(), s.License, space.Data.Labels.License)
	assert.Equal(s.T(), s.Desc, space.Data.Desc)
	assert.Equal(s.T(), s.Visibility, space.Data.Visibility)
	assert.Equal(s.T(), "mindspore", space.Data.Labels.Framework)

	// 删除space
	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 创建不匹配的hardware和base image
func (s *SuiteSpace) TestSpaceCreateInvalidHWIMAGE() {
	_, r, err := ApiRest.SpaceApi.V1SpacePost(AuthRest, swaggerRest.ControllerReqToCreateSpace{
		Desc:       s.Desc,
		Fullname:   s.Fullname,
		Hardware:   s.CPU,
		License:    s.License,
		Name:       s.Name,
		Owner:      s.Owner,
		Sdk:        s.Sdk,
		Visibility: s.Visibility,
		BaseImage:  s.BaseImageMSNPU,
	})

	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	_, r, err = ApiRest.SpaceApi.V1SpacePost(AuthRest, swaggerRest.ControllerReqToCreateSpace{
		Desc:       s.Desc,
		Fullname:   s.Fullname,
		Hardware:   s.CPU,
		License:    s.License,
		Name:       s.Name,
		Owner:      s.Owner,
		Sdk:        s.Sdk,
		Visibility: s.Visibility,
		BaseImage:  s.BaseImagePTNPU,
	})

	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	_, r, err = ApiRest.SpaceApi.V1SpacePost(AuthRest, swaggerRest.ControllerReqToCreateSpace{
		Desc:       s.Desc,
		Fullname:   s.Fullname,
		Hardware:   s.NPU,
		License:    s.License,
		Name:       s.Name,
		Owner:      s.Owner,
		Sdk:        s.Sdk,
		Visibility: s.Visibility,
		BaseImage:  s.BaseImagePTCPU,
	})

	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	_, r, err = ApiRest.SpaceApi.V1SpacePost(AuthRest, swaggerRest.ControllerReqToCreateSpace{
		Desc:       s.Desc,
		Fullname:   s.Fullname,
		Hardware:   s.NPU,
		License:    s.License,
		Name:       s.Name,
		Owner:      s.Owner,
		Sdk:        s.Sdk,
		Visibility: s.Visibility,
		BaseImage:  s.BaseImageMSCPU,
	})

	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// 批量查询
// 以framework过滤
func (s *SuiteSpace) TestSpaceList() {
	space, r, err := ApiRest.SpaceApi.V1SpacePost(AuthRest, swaggerRest.ControllerReqToCreateSpace{
		Desc:       s.Desc,
		Fullname:   s.Fullname,
		Hardware:   s.NPU,
		License:    s.License,
		Name:       s.Name,
		Owner:      s.Owner,
		Sdk:        s.Sdk,
		Visibility: s.Visibility,
		BaseImage:  s.BaseImageMSNPU,
	})

	id1 := space.Data

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	space, r, err = ApiRest.SpaceApi.V1SpacePost(AuthRest, swaggerRest.ControllerReqToCreateSpace{
		Desc:       s.Desc,
		Fullname:   s.Fullname,
		Hardware:   s.CPU,
		License:    "apache-2.0",
		Name:       s.Name1,
		Owner:      s.Owner,
		Sdk:        s.Sdk,
		Visibility: s.Visibility,
		BaseImage:  s.BaseImageMSCPU,
	})

	id2 := space.Data

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	space, r, err = ApiRest.SpaceApi.V1SpacePost(AuthRest, swaggerRest.ControllerReqToCreateSpace{
		Desc:       s.Desc,
		Fullname:   s.Fullname,
		Hardware:   s.CPU,
		License:    "mpl-2.0",
		Name:       s.Name2,
		Owner:      s.Owner,
		Sdk:        s.Sdk,
		Visibility: s.Visibility,
		BaseImage:  s.BaseImagePTCPU,
	})

	id3 := space.Data

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	space, r, err = ApiRest.SpaceApi.V1SpacePost(AuthRest, swaggerRest.ControllerReqToCreateSpace{
		Desc:       s.Desc,
		Fullname:   s.Fullname,
		Hardware:   s.NPU,
		License:    "isc",
		Name:       s.Name3,
		Owner:      s.Owner,
		Sdk:        s.Sdk,
		Visibility: s.Visibility,
		BaseImage:  s.BaseImagePTNPU,
	})

	id4 := space.Data

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	// 测试通过framework过滤
	data, r, err := ApiRest.SpaceRestfulApi.V1SpaceGet(AuthRest, s.Owner, &swaggerRest.SpaceRestfulApiV1SpaceGetOpts{
		Framework: optional.NewString("pytorch"),
		Count:     optional.NewBool(true),
	})

	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	spaces := data.Data.Spaces
	assert.Equal(s.T(), int32(2), data.Data.Total)
	for _, space := range spaces {
		if space.Name != s.Name3 && space.Name != s.Name2 {
			s.T().Errorf("space should be %s or %s", s.Name3, s.Name2)
		}
	}

	// 测试通过 license 过滤
	data, r, err = ApiRest.SpaceRestfulApi.V1SpaceGet(AuthRest, s.Owner, &swaggerRest.SpaceRestfulApiV1SpaceGetOpts{
		License: optional.NewString("mit"),
		Count:   optional.NewBool(true),
	})

	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	spaces = data.Data.Spaces
	assert.Equal(s.T(), int32(1), data.Data.Total)
	for _, space := range spaces {
		if space.Name != s.Name {
			s.T().Errorf("space should be %s", s.Name)
		}
	}

	// 无效的过滤参数
	data, r, err = ApiRest.SpaceRestfulApi.V1SpaceGet(AuthRest, s.Owner, &swaggerRest.SpaceRestfulApiV1SpaceGetOpts{
		License: optional.NewString("mit1"),
	})

	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	// 无效的过滤参数
	data, r, err = ApiRest.SpaceRestfulApi.V1SpaceGet(AuthRest, s.Owner, &swaggerRest.SpaceRestfulApiV1SpaceGetOpts{
		Framework: optional.NewString("pppp"),
	})

	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	// 删除space
	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest, id1)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest, id2)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest, id3)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest, id4)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 设置domain
func (s *SuiteSpace) TestSpaceSetDomain() {
	data, r, err := ApiRest.SpaceApi.V1SpacePost(AuthRest, swaggerRest.ControllerReqToCreateSpace{
		Desc:       s.Desc,
		Fullname:   s.Fullname,
		Hardware:   s.NPU,
		License:    s.License,
		Name:       s.Name,
		Owner:      s.Owner,
		Sdk:        s.Sdk,
		Visibility: s.Visibility,
		BaseImage:  s.BaseImageMSNPU,
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := data.Data
	assert.NotEqual(s.T(), "", id)

	_, r, err = ApiInteral.SpaceInternalApi.V1SpaceIdLabelPut(Interal, id, swaggerInternal.GithubComOpenmerlinMerlinServerSpaceControllerReqToResetLabel{
		Task:    "nlp",
		License: "apache-2.0",
	})

	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	space, r, err := ApiRest.SpaceRestfulApi.V1SpaceOwnerNameGet(AuthRest, s.Owner, s.Name)

	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), s.Name, space.Data.Name)
	assert.Equal(s.T(), s.Owner, space.Data.Owner)
	assert.Equal(s.T(), s.Visibility, space.Data.Visibility)
	assert.Equal(s.T(), s.Sdk, space.Data.Sdk)
	assert.Equal(s.T(), s.Fullname, space.Data.Fullname)
	assert.Equal(s.T(), strings.ToLower(s.NPU), space.Data.Hardware)
	assert.Equal(s.T(), "apache-2.0", space.Data.Labels.License)
	assert.Equal(s.T(), s.Desc, space.Data.Desc)
	assert.Equal(s.T(), s.Visibility, space.Data.Visibility)
	assert.Equal(s.T(), "mindspore", space.Data.Labels.Framework)
	assert.Equal(s.T(), "nlp", space.Data.Labels.Task)

	// invalid domain
	_, r, err = ApiInteral.SpaceInternalApi.V1SpaceIdLabelPut(Interal, id, swaggerInternal.GithubComOpenmerlinMerlinServerSpaceControllerReqToResetLabel{
		Task:    "123",
		License: "apache-2.0",
	})

	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	// 删除space
	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 搜索的name字段支持在owner/name中搜索，也支持在fullname中搜索
func (s *SuiteSpace) TestSpaceSearch() {
	space, r, err := ApiRest.SpaceApi.V1SpacePost(AuthRest, swaggerRest.ControllerReqToCreateSpace{
		Desc:       s.Desc,
		Fullname:   "xxx",
		Hardware:   s.NPU,
		License:    s.License,
		Name:       s.Name,
		Owner:      s.Owner,
		Sdk:        s.Sdk,
		Visibility: "private",
		BaseImage:  s.BaseImageMSNPU,
	})

	id1 := space.Data

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	space, r, err = ApiRest.SpaceApi.V1SpacePost(AuthRest, swaggerRest.ControllerReqToCreateSpace{
		Desc:       s.Desc,
		Fullname:   "yyy",
		Hardware:   s.CPU,
		License:    "apache-2.0",
		Name:       s.Name1,
		Owner:      s.Owner,
		Sdk:        s.Sdk,
		Visibility: "private",
		BaseImage:  s.BaseImageMSCPU,
	})

	id2 := space.Data

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	space, r, err = ApiRest.SpaceApi.V1SpacePost(AuthRest, swaggerRest.ControllerReqToCreateSpace{
		Desc:       s.Desc,
		Hardware:   s.CPU,
		Fullname:   "xxx",
		License:    "mpl-2.0",
		Name:       s.Name2,
		Owner:      s.Owner,
		Sdk:        s.Sdk,
		Visibility: s.Visibility,
		BaseImage:  s.BaseImagePTCPU,
	})

	id3 := space.Data

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	space, r, err = ApiRest.SpaceApi.V1SpacePost(AuthRest, swaggerRest.ControllerReqToCreateSpace{
		Desc:       s.Desc,
		Hardware:   s.NPU,
		License:    "isc",
		Name:       s.Name3,
		Owner:      s.Owner,
		Sdk:        s.Sdk,
		Visibility: s.Visibility,
		BaseImage:  s.BaseImagePTNPU,
	})

	id4 := space.Data

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	// name中带/可以过滤出所有space
	data, r, err := ApiRest.SpaceRestfulApi.V1SpaceGet(AuthRest, s.Owner, &swaggerRest.SpaceRestfulApiV1SpaceGetOpts{
		Name:  optional.NewString("/"),
		Count: optional.NewBool(true),
	})

	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), int32(2), data.Data.Total)

	// 通过fullname可以进行搜索
	data, r, err = ApiRest.SpaceRestfulApi.V1SpaceGet(AuthRest, s.Owner, &swaggerRest.SpaceRestfulApiV1SpaceGetOpts{
		Name:  optional.NewString("xxx"),
		Count: optional.NewBool(true),
	})

	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), int32(1), data.Data.Total)

	// 通过owner可以进行搜索
	data, r, err = ApiRest.SpaceRestfulApi.V1SpaceGet(AuthRest, s.Owner, &swaggerRest.SpaceRestfulApiV1SpaceGetOpts{
		Name:  optional.NewString("test1"),
		Count: optional.NewBool(true),
	})

	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), int32(2), data.Data.Total)

	// 通过owner/name可以进行搜索
	data, r, err = ApiRest.SpaceRestfulApi.V1SpaceGet(AuthRest, s.Owner, &swaggerRest.SpaceRestfulApiV1SpaceGetOpts{
		Name:  optional.NewString("test1/testspace"),
		Count: optional.NewBool(true),
	})

	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), int32(2), data.Data.Total)

	// 通过name可以进行搜索模糊搜索
	data, r, err = ApiRest.SpaceRestfulApi.V1SpaceGet(AuthRest, s.Owner, &swaggerRest.SpaceRestfulApiV1SpaceGetOpts{
		Name:  optional.NewString("testspace"),
		Count: optional.NewBool(true),
	})

	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), int32(2), data.Data.Total)
	// 无效的过滤参数
	data, r, err = ApiRest.SpaceRestfulApi.V1SpaceGet(AuthRest, s.Owner, &swaggerRest.SpaceRestfulApiV1SpaceGetOpts{
		License: optional.NewString("mit1"),
	})

	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	// 无效的过滤参数
	data, r, err = ApiRest.SpaceRestfulApi.V1SpaceGet(AuthRest, s.Owner, &swaggerRest.SpaceRestfulApiV1SpaceGetOpts{
		Framework: optional.NewString("pppp"),
	})

	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	// 其他人无法搜索到私有仓库
	// fullname 为 xxx的有2个space，一个公有，一个私有，应该可以查出来1个
	data, r, err = ApiRest.SpaceRestfulApi.V1SpaceGet(AuthRest2, s.Owner, &swaggerRest.SpaceRestfulApiV1SpaceGetOpts{
		Name:  optional.NewString("xxx"),
		Count: optional.NewBool(true),
	})

	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), int32(1), data.Data.Total)

	data, r, err = ApiRest.SpaceRestfulApi.V1SpaceGet(AuthRest2, s.Owner, &swaggerRest.SpaceRestfulApiV1SpaceGetOpts{
		Name:  optional.NewString("yyy"),
		Count: optional.NewBool(true),
	})

	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), int32(0), data.Data.Total)

	// 删除space
	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest, id1)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest, id2)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest, id3)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest, id4)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// space不能创建超过配置，当前限制是4个
func (s *SuiteSpace) TestSpaceMaxCount() {
	space, r, err := ApiRest.SpaceApi.V1SpacePost(AuthRest, swaggerRest.ControllerReqToCreateSpace{
		Desc:       s.Desc,
		Fullname:   "xxx",
		Hardware:   s.NPU,
		License:    s.License,
		Name:       s.Name,
		Owner:      s.Owner,
		Sdk:        s.Sdk,
		Visibility: s.Visibility,
		BaseImage:  s.BaseImageMSNPU,
	})

	id1 := space.Data

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	space, r, err = ApiRest.SpaceApi.V1SpacePost(AuthRest, swaggerRest.ControllerReqToCreateSpace{
		Desc:       s.Desc,
		Fullname:   "yyy",
		Hardware:   s.CPU,
		License:    "apache-2.0",
		Name:       s.Name1,
		Owner:      s.Owner,
		Sdk:        s.Sdk,
		Visibility: s.Visibility,
		BaseImage:  s.BaseImageMSCPU,
	})

	id2 := space.Data

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	space, r, err = ApiRest.SpaceApi.V1SpacePost(AuthRest, swaggerRest.ControllerReqToCreateSpace{
		Desc:       s.Desc,
		Hardware:   s.CPU,
		License:    "mpl-2.0",
		Name:       s.Name2,
		Owner:      s.Owner,
		Sdk:        s.Sdk,
		Visibility: s.Visibility,
		BaseImage:  s.BaseImagePTCPU,
	})

	id3 := space.Data

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	space, r, err = ApiRest.SpaceApi.V1SpacePost(AuthRest, swaggerRest.ControllerReqToCreateSpace{
		Desc:       s.Desc,
		Hardware:   s.NPU,
		License:    "isc",
		Name:       s.Name3,
		Owner:      s.Owner,
		Sdk:        s.Sdk,
		Visibility: s.Visibility,
		BaseImage:  s.BaseImagePTNPU,
	})

	id4 := space.Data

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	space, r, err = ApiRest.SpaceApi.V1SpacePost(AuthRest, swaggerRest.ControllerReqToCreateSpace{
		Desc:       s.Desc,
		Hardware:   s.NPU,
		License:    "isc",
		Name:       "cannotcreate",
		Owner:      s.Owner,
		Sdk:        s.Sdk,
		Visibility: s.Visibility,
		BaseImage:  s.BaseImagePTNPU,
	})

	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	// 删除space
	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest, id1)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest, id2)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest, id3)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest, id4)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

func TestSpacen(t *testing.T) {
	suite.Run(t, new(SuiteSpace))
}
