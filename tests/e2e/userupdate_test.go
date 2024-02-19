package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	swagger "e2e/client"
)

type SuiteUserUpdate struct {
	suite.Suite
}

// 用户可以正常更新个人信息
func (s *SuiteUserUpdate) TestUpdateUserInfoValidData() {

	d := swagger.ControllerUserBasicInfoUpdateRequest{
		Fullname:    "read full name",
		AvatarId:    "https://avatars.githubusercontent.com/u/2853724?v=5",
		Description: "valid desc",
	}
	data, r, err := Api.UserApi.V1UserPut(Auth, d)
	assert.Equal(s.T(), r.StatusCode, 202)
	assert.Nil(s.T(), err)

	user := getData(s.T(), data.Data)

	assert.Equalf(s.T(), user["fullname"], "read full name", "fullname is not equal")
	assert.Equalf(s.T(), user["avatar_id"], "https://avatars.githubusercontent.com/u/2853724?v=5", "avatar_id is not equal")
	assert.Equalf(s.T(), user["description"], "valid desc", "description is not equal")
}

// fullname不能为空
func (s *SuiteUserUpdate) TestUpdateUserInfoEmptyFullname() {

	d := swagger.ControllerUserBasicInfoUpdateRequest{
		Fullname: "",
	}
	_, r, err := Api.UserApi.V1UserPut(Auth, d)
	assert.Equal(s.T(), 400, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// 无效的fullname会导致更新失败
func (s *SuiteUserUpdate) TestUpdateUserInfoInvalidFullname() {

	d := swagger.ControllerUserBasicInfoUpdateRequest{
		Fullname: string(make([]byte, 201)),
	}
	_, r, err := Api.UserApi.V1UserPut(Auth, d)
	assert.Equal(s.T(), 400, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// 无效的desc会导致更新失败
func (s *SuiteUserUpdate) TestUpdateUserInfoInvalidDesc() {

	d := swagger.ControllerUserBasicInfoUpdateRequest{
		Description: string(make([]byte, 201)),
	}
	_, r, err := Api.UserApi.V1UserPut(Auth, d)
	assert.Equal(s.T(), 400, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// 无效的avatar会导致更新失败
func (s *SuiteUserUpdate) TestUpdateUserInfoInvalidAvatar() {

	d := swagger.ControllerUserBasicInfoUpdateRequest{
		AvatarId: "invalid avatarid",
	}
	_, r, err := Api.UserApi.V1UserPut(Auth, d)
	assert.Equal(s.T(), 400, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// 有效的avatar会导致更新成功
func (s *SuiteUserUpdate) TestUpdateUserInfoValidAvatar() {

	d := swagger.ControllerUserBasicInfoUpdateRequest{
		AvatarId: "https://avatars.githubusercontent.com/u/2853724?v=4",
	}
	data, r, err := Api.UserApi.V1UserPut(Auth, d)
	assert.Equal(s.T(), 202, r.StatusCode)
	assert.Nil(s.T(), err)

	user := getData(s.T(), data.Data)

	assert.Equal(s.T(), "https://avatars.githubusercontent.com/u/2853724?v=4", user["avatar_id"])
}

func (s *SuiteUserUpdate) TearDownSuite() {
	d := swagger.ControllerUserBasicInfoUpdateRequest{
		Fullname:    "testfullname",
		AvatarId:    "https://avatars.githubusercontent.com/u/2853724?v=1",
		Description: "testdesc",
	}
	_, r, err := Api.UserApi.V1UserPut(Auth, d)
	assert.Equal(s.T(), 202, r.StatusCode)
	assert.Nil(s.T(), err)
}

func TestUserUpdate(t *testing.T) {
	suite.Run(t, new(SuiteUserUpdate))
}
