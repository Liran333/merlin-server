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

// SuiteUser used for testing
type SuiteUser struct {
	suite.Suite
	phone string
}

// SetupSuite used for testing
func (s *SuiteUser) SetupSuite() {
	s.phone = "13333333334"
	d := swagger.ControllerUserBasicInfoUpdateRequest{
		Fullname:    "read full name",
		AvatarId:    "https://avatars.githubusercontent.com/u/2853724?v=5",
		Description: "valid desc",
	}

	_, r, err := Api.UserApi.V1UserPut(Auth, d)
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestGetUser used for testing
// 正常登录的用户可以获取用户信息
func (s *SuiteUser) TestGetUser() {

	data, r, err := Api.UserApi.V1UserGet(Auth)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	user := getData(s.T(), data.Data)

	assert.Equalf(s.T(), user["fullname"], "read full name", "fullname is not equal")
	assert.Equalf(s.T(), user["avatar_id"], "https://avatars.githubusercontent.com/u/2853724?v=5", "avatar_id is not equal")
	assert.Equalf(s.T(), user["description"], "valid desc", "description is not equal")
	assert.Equal(s.T(), getInt64(s.T(), user["type"]), int64(0))
	assert.NotEqual(s.T(), "", user["id"])
	assert.NotEqual(s.T(), int64(0), getInt64(s.T(), user["created_at"]))
	assert.NotEqual(s.T(), int64(0), getInt64(s.T(), user["updated_at"]))
	assert.Equal(s.T(), s.phone, user["phone"])
}

// TestGetOtherUser used for testing
func (s *SuiteUser) TestGetOtherUser() {

	data, r, err := Api.OrganizationApi.V1AccountNameGet(Auth, "test2")
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	user := getData(s.T(), data.Data)

	assert.Equalf(s.T(), "test2", user["account"], "fullname is not equal")
	assert.Equal(s.T(), "", user["email"])
	assert.Equal(s.T(), int64(0), getInt64(s.T(), user["type"]))
	assert.NotEqual(s.T(), "", user["id"])
	assert.NotEqual(s.T(), int64(0), getInt64(s.T(), user["created_at"]))
	assert.NotEqual(s.T(), int64(0), getInt64(s.T(), user["updated_at"]))
	assert.Equal(s.T(), "", user["phone"])
}

// TestGetUserNoToken used for testing
// 未登录用户无法获取个人信息
func (s *SuiteUser) TestGetUserNoToken() {

	_, r, err := Api.UserApi.V1UserGet(context.Background())
	assert.Equal(s.T(), http.StatusUnauthorized, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestGetOtherUserNoToken used for testing
// 未登录用户获取其他人信息
func (s *SuiteUser) TestGetOtherUserNoToken() {

	data, r, err := Api.OrganizationApi.V1AccountNameGet(context.Background(), "test2")
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	user := getData(s.T(), data.Data)

	assert.Equalf(s.T(), "test2", user["account"], "fullname is not equal")

	assert.Equal(s.T(), "", user["email"])
	assert.Equal(s.T(), "", user["phone"])
}

// TestUser used for testing
func TestUser(t *testing.T) {
	suite.Run(t, new(SuiteUser))
}
