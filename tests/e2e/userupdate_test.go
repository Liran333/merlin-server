/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package e2e

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	swaggerRest "e2e/client_rest"
)

// SuiteUserUpdate used for testing
type SuiteUserUpdate struct {
	suite.Suite
}

// TestUpdateUserInfoValidData used for testing
// 用户可以正常更新个人信息
func (s *SuiteUserUpdate) TestUpdateUserInfoValidData() {

	d := swaggerRest.ControllerUserBasicInfoUpdateRequest{
		Fullname:    "read full name",
		AvatarId:    "https://avatars.githubusercontent.com/u/2853724?v=5",
		Description: "valid desc",
	}
	data, r, err := ApiRest.UserApi.V1UserPut(AuthRest, d)
	assert.Equal(s.T(), r.StatusCode, http.StatusAccepted)
	assert.Nil(s.T(), err)

	user := getData(s.T(), data.Data)

	assert.Equalf(s.T(), user["fullname"], "read full name", "fullname is not equal")
	assert.Equalf(s.T(), user["avatar_id"], "https://avatars.githubusercontent.com/u/2853724?v=5", "avatar_id is not equal")
	assert.Equalf(s.T(), user["description"], "valid desc", "description is not equal")
}

// TestUpdateUserInfoEmptyFullname used for testing
// fullname不能为空
func (s *SuiteUserUpdate) TestUpdateUserInfoEmptyFullname() {

	d := swaggerRest.ControllerUserBasicInfoUpdateRequest{
		Fullname: "",
	}
	_, r, err := ApiRest.UserApi.V1UserPut(AuthRest, d)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestUpdateUserInfoValidFullname used for testing
// fullname更新成功
func (s *SuiteUserUpdate) TestUpdateUserInfoValidFullname() {

	d := swaggerRest.ControllerUserBasicInfoUpdateRequest{
		Fullname: "testFullname",
	}
	_, r, err := ApiRest.UserApi.V1UserPut(AuthRest, d)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode, "Expected success for fullname update")
}

// TestUpdateUserInfoInvalidFullname used for testing
// 无效的fullname会导致更新失败
func (s *SuiteUserUpdate) TestUpdateUserInfoInvalidFullname() {

	d := swaggerRest.ControllerUserBasicInfoUpdateRequest{
		Fullname: string(make([]byte, http.StatusCreated)),
	}
	_, r, err := ApiRest.UserApi.V1UserPut(AuthRest, d)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestUpdateUserInfoValidDesc used for testing
// desc更新成功
func (s *SuiteUserUpdate) TestUpdateUserInfoValidDesc() {

	d := swaggerRest.ControllerUserBasicInfoUpdateRequest{
		Description: "test description",
	}
	_, r, err := ApiRest.UserApi.V1UserPut(AuthRest, d)
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestUpdateUserInfoInvalidDesc used for testing
// 无效的desc会导致更新失败
func (s *SuiteUserUpdate) TestUpdateUserInfoInvalidDesc() {

	d := swaggerRest.ControllerUserBasicInfoUpdateRequest{
		Description: string(make([]byte, 2049)),
	}
	_, r, err := ApiRest.UserApi.V1UserPut(AuthRest, d)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestUpdateUserInfoInvalidAvatar used for testing
// 无效的avatar会导致更新失败
func (s *SuiteUserUpdate) TestUpdateUserInfoInvalidAvatar() {

	d := swaggerRest.ControllerUserBasicInfoUpdateRequest{
		AvatarId: "invalid avatarid",
	}
	_, r, err := ApiRest.UserApi.V1UserPut(AuthRest, d)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestUpdateUserInfoValidAvatar used for testing
// 有效的avatar会导致更新成功
func (s *SuiteUserUpdate) TestUpdateUserInfoValidAvatar() {

	d := swaggerRest.ControllerUserBasicInfoUpdateRequest{
		AvatarId: "https://avatars.githubusercontent.com/u/2853724?v=4",
	}
	data, r, err := ApiRest.UserApi.V1UserPut(AuthRest, d)
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	user := getData(s.T(), data.Data)

	assert.Equal(s.T(), "https://avatars.githubusercontent.com/u/2853724?v=4", user["avatar_id"])
}

// TearDownSuite used for testing
func (s *SuiteUserUpdate) TearDownSuite() {
	d := swaggerRest.ControllerUserBasicInfoUpdateRequest{
		Fullname:    "testfullname",
		AvatarId:    "https://avatars.githubusercontent.com/u/2853724?v=1",
		Description: "testdesc",
	}
	_, r, err := ApiRest.UserApi.V1UserPut(AuthRest, d)
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestUserUpdate used for testing
func TestUserUpdate(t *testing.T) {
	suite.Run(t, new(SuiteUserUpdate))
}
