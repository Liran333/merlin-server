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
type SuitPermissionInternal struct {
	suite.Suite
}

// SetupSuite used for testing
func (s *SuitPermissionInternal) SetupSuite() {
}

// TearDownSuite used for testing
func (s *SuitPermissionInternal) TearDownSuite() {
}

// TestPermissionRead used for testing
func (s *SuitPermissionInternal) TestPermissionRead() {
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

	_, r, err = ApiInteral.PermissionApi.V1CoderepoPermissionReadPost(Interal,
		swaggerInternal.ControllerReqToCheckPermission{
			Owner: "test2",
			Name:  "testspace",
			User:  "test2",
		})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestPermissionUpdate used for testing
func (s *SuitPermissionInternal) TestPermissionUpdate() {
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

	_, r, err = ApiInteral.PermissionApi.V1CoderepoPermissionUpdatePost(Interal,
		swaggerInternal.ControllerReqToCheckPermission{
			Owner: "test2",
			Name:  "testspace",
			User:  "test2",
		})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest2, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestPermissionINternal used for testing
func TestPermissionINternal(t *testing.T) {
	suite.Run(t, new(SuitPermissionInternal))
}
