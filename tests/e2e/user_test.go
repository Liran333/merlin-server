package e2e

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	swagger "e2e/client"
)

type SuiteUser struct {
	suite.Suite
}

func (s *SuiteUser) SetupSuite() {

	d := swagger.ControllerUserBasicInfoUpdateRequest{
		Fullname:    "read full name",
		AvatarId:    "https://avatars.githubusercontent.com/u/2853724?v=5",
		Description: "valid desc",
		Email:       "testupdateuser@modelfoudnry.cn",
	}

	_, r, err := Api.UserApi.V1UserPut(Auth, d)
	assert.Equal(s.T(), 202, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 正常登录的用户可以获取用户信息
func (s *SuiteUser) TestGetUser() {

	data, r, err := Api.UserApi.V1UserGet(Auth)
	assert.Equal(s.T(), 200, r.StatusCode)
	assert.Nil(s.T(), err)

	user := getData(s.T(), data.Data)

	assert.Equalf(s.T(), user["fullname"], "read full name", "fullname is not equal")
	assert.Equalf(s.T(), user["avatar_id"], "https://avatars.githubusercontent.com/u/2853724?v=5", "avatar_id is not equal")
	assert.Equalf(s.T(), user["description"], "valid desc", "description is not equal")
	assert.Equalf(s.T(), user["email"], "testupdateuser@modelfoudnry.cn", "email is not equal")
	assert.Equal(s.T(), getInt64(s.T(), user["type"]), int64(0))
}

// 未登录用户无法获取个人信息
func (s *SuiteUser) TestGetUserNoToken() {

	_, r, err := Api.UserApi.V1UserGet(context.Background())
	assert.Equal(s.T(), 401, r.StatusCode)
	assert.NotNil(s.T(), err)
}

func TestUser(t *testing.T) {
	suite.Run(t, new(SuiteUser))
}
