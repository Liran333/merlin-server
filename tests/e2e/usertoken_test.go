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

	swaggerRest "e2e/client_rest"
)

// SuiteUserToken used for testing
type SuiteUserToken struct {
	suite.Suite
	id         string
	readtoken  context.Context
	writetoken context.Context
}

// SetupSuite used for testing
func (s *SuiteUserToken) SetupSuite() {
	data, r, err := ApiRest.UserApi.V1UserGet(AuthRest)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	user := getData(s.T(), data.Data)

	assert.NotEqual(s.T(), "", user["id"])
	s.id = getString(s.T(), user["id"])
	s.T().Logf("user id: %s", s.id)

	d := swaggerRest.ControllerTokenCreateRequest{
		Name: "testread",
		Perm: "read",
	}

	tokenData, r, err := ApiRest.UserApi.V1UserTokenPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	m := getData(s.T(), tokenData.Data)

	assert.NotEqual(s.T(), "", getString(s.T(), m["token"]))
	assert.Equal(s.T(), s.id, m["owner_id"])

	s.readtoken = newAuthRestCtx(getString(s.T(), m["token"]))

	d = swaggerRest.ControllerTokenCreateRequest{
		Name: "testwrite",
		Perm: "write",
	}

	_, r, err = ApiRest.UserApi.V1UserTokenPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	m1 := getData(s.T(), tokenData.Data)

	assert.NotEqual(s.T(), "", getString(s.T(), m1["token"]))
	assert.Equal(s.T(), s.id, m1["owner_id"])
	s.writetoken = newAuthRestCtx(getString(s.T(), m["token"]))
}

// TearDownSuite used for testing
func (s *SuiteUserToken) TearDownSuite() {
	r, err := ApiRest.UserApi.V1UserTokenNameDelete(AuthRest, "testread")
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.UserApi.V1UserTokenNameDelete(AuthRest, "testwrite")
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestCreateDuplicateToken used for testing
// 无法创建同名token
func (s *SuiteUserToken) TestCreateDuplicateToken() {
	d := swaggerRest.ControllerTokenCreateRequest{
		Name: "testread",
		Perm: "read",
	}

	data, r, err := ApiRest.UserApi.V1UserTokenPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	s.T().Logf("create duplicate token return: %s", data.Msg)
	assert.NotNil(s.T(), err)
}

// TestGetUserTokenWithNoToken used for testing
// 未登录用户无法查询用户的token信息
func (s *SuiteUserToken) TestGetUserTokenWithNoToken() {

	_, r, err := ApiRest.UserApi.V1UserTokenGet(context.Background())
	assert.Equal(s.T(), http.StatusUnauthorized, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestGetUserToken used for testing
// 正常登录的用户可以查询toke信息
func (s *SuiteUserToken) TestGetUserToken() {

	data, r, err := ApiRest.UserApi.V1UserTokenGet(AuthRest)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	tokens := getArrary(s.T(), data.Data)

	readFound := false
	writeFound := false
	count := 0

	for token := range tokens {
		if tokens[token] == nil {
			continue
		}
		count += 1
		if tokens[token]["name"] == "testread" {
			assert.Equal(s.T(), "read", getString(s.T(), tokens[token]["permission"]))
			assert.Equal(s.T(), "", getString(s.T(), tokens[token]["token"]))
			assert.Equal(s.T(), s.id, tokens[token]["owner_id"])
			readFound = true
		}

		if tokens[token]["name"] == "testwrite" {
			assert.Equal(s.T(), "write", getString(s.T(), tokens[token]["permission"]))
			assert.Equal(s.T(), "", getString(s.T(), tokens[token]["token"]))
			assert.Equal(s.T(), s.id, tokens[token]["owner_id"])
			writeFound = true
		}

		assert.NotEqual(s.T(), 0, getInt64(s.T(), tokens[token]["created_at"]))
		assert.NotEqual(s.T(), 0, getInt64(s.T(), tokens[token]["updated_at"]))
		assert.NotEqual(s.T(), "", getString(s.T(), tokens[token]["id"]))
		assert.Equal(s.T(), s.id, tokens[token]["owner_id"])
	}

	assert.Equal(s.T(), countThree, count)
	assert.True(s.T(), readFound)
	assert.True(s.T(), writeFound)
}

// TestTokenCreateTokenInvalidName used for testing
// 无效的token权限会导致创建token失败
func (s *SuiteUserToken) TestTokenCreateTokenInvalidName() {
	// test a read permission token
	d := swaggerRest.ControllerTokenCreateRequest{
		Name: "read",
		Perm: "invalidperm",
	}

	_, r, err := ApiRest.UserApi.V1UserTokenPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestTokenCreateTokenInvalidNameChar used for testing
// token名只能包括[a-zA-Z0-9_-]
func (s *SuiteUserToken) TestTokenCreateTokenInvalidNameChar() {
	invalidChar := string("!@#$%^&*(){}[]")
	for _, c := range invalidChar {
		// test a read permission token
		d := swaggerRest.ControllerTokenCreateRequest{
			Name: "read" + string(c),
			Perm: "invalidperm",
		}

		data, r, err := ApiRest.UserApi.V1UserTokenPost(AuthRest, d)
		assert.Equalf(s.T(), http.StatusBadRequest, r.StatusCode, data.Msg)
		assert.NotNil(s.T(), err)
	}
}

// TestTokenCreateTokenNameCantBeInt used for testing
// token名不能是纯数字
func (s *SuiteUserToken) TestTokenCreateTokenNameCantBeInt() {
	// test a read permission token
	d := swaggerRest.ControllerTokenCreateRequest{
		Name: "12",
		Perm: "write",
	}

	_, r, err := ApiRest.UserApi.V1UserTokenPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestTokenCreateToken used for testing
// 创建token成功
// read权限无权删除token
func (s *SuiteUserToken) TestTokenCreateToken() {
	_, r, err := ApiRest.UserApi.V1UserTokenGet(s.readtoken)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestTokenDeleteToken used for testing
// 删除不存在的token报404
func (s *SuiteUserToken) TestTokenDeleteToken() {
	r, err := ApiRest.UserApi.V1UserTokenNameDelete(AuthRest, "nonexist")
	assert.Equal(s.T(), http.StatusNotFound, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestTokenDeleteToken used for testing
// token最大个数限制
func (s *SuiteUserToken) TestTokenMaxCount() {
	// setup中已经创建了2个，默认创建了一个，e2e测试配置的上限为4个，因此不应该能继续创建token
	d := swaggerRest.ControllerTokenCreateRequest{
		Name: "next1",
		Perm: "write",
	}

	_, r, err := ApiRest.UserApi.V1UserTokenPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	d = swaggerRest.ControllerTokenCreateRequest{
		Name: "next2",
		Perm: "write",
	}

	_, r, err = ApiRest.UserApi.V1UserTokenPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	r, err = ApiRest.UserApi.V1UserTokenNameDelete(AuthRest, "next1")
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestUserToken used for testing
func TestUserToken(t *testing.T) {
	suite.Run(t, new(SuiteUserToken))
}
