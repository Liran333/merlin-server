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

// SuiteUserToken used for testing
type SuiteUserToken struct {
	suite.Suite
	id string
}

// SetupSuite used for testing
func (s *SuiteUserToken) SetupSuite() {
	data, r, err := Api.UserApi.V1UserGet(Auth)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	user := getData(s.T(), data.Data)

	assert.NotEqual(s.T(), "", user["id"])
	s.id = getString(s.T(), user["id"])
	s.T().Logf("user id: %s", s.id)

	d := swagger.ControllerTokenCreateRequest{
		Name: "testread",
		Perm: "read",
	}

	data, r, err = Api.UserApi.V1UserTokenPost(Auth, d)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	m := getData(s.T(), data.Data)

	assert.NotEqual(s.T(), "", getString(s.T(), m["token"]))
	assert.Equal(s.T(), s.id, m["owner_id"])

	d = swagger.ControllerTokenCreateRequest{
		Name: "testwrite",
		Perm: "write",
	}

	_, r, err = Api.UserApi.V1UserTokenPost(Auth, d)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	m = getData(s.T(), data.Data)

	assert.NotEqual(s.T(), "", getString(s.T(), m["token"]))
	assert.Equal(s.T(), s.id, m["owner_id"])
}

// TearDownSuite used for testing
func (s *SuiteUserToken) TearDownSuite() {
	r, err := Api.UserApi.V1UserTokenNameDelete(Auth, "testread")
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = Api.UserApi.V1UserTokenNameDelete(Auth, "testwrite")
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestVerifyToken verify token
func (s *SuiteUserToken) TestVerifyToken() {
	d := swagger.ControllerTokenCreateRequest{
		Name: "testverify",
		Perm: "read",
	}

	data, r, err := Api.UserApi.V1UserTokenPost(Auth, d)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	m := getData(s.T(), data.Data)

	assert.NotEqual(s.T(), "", getString(s.T(), m["token"]))
	assert.Equal(s.T(), s.id, m["owner_id"])

	t := swagger.ControllerTokenVerifyRequest{
		Token: getString(s.T(), m["token"]),
	}

	data, r, err = InteralApi.UserApi.V1UserTokenVerifyPost(Interal, t)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = Api.UserApi.V1UserTokenNameDelete(Auth, "testverify")
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestVerifyInvalidToken verify invalid token
func (s *SuiteUserToken) TestVerifyInvalidToken() {

	t := swagger.ControllerTokenVerifyRequest{
		Token: getString(s.T(), "2233445notok"),
	}

	_, r, err := InteralApi.UserApi.V1UserTokenVerifyPost(Interal, t)
	assert.Equal(s.T(), http.StatusUnauthorized, r.StatusCode)
	assert.NotNil(s.T(), err)

	t = swagger.ControllerTokenVerifyRequest{}

	_, r, err = InteralApi.UserApi.V1UserTokenVerifyPost(Interal, t)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestCreateDuplicateToken used for testing
// 无法创建同名token
func (s *SuiteUserToken) TestCreateDuplicateToken() {
	d := swagger.ControllerTokenCreateRequest{
		Name: "testread",
		Perm: "read",
	}

	data, r, err := Api.UserApi.V1UserTokenPost(Auth, d)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	s.T().Logf("create duplicate token return: %s", data.Msg)
	assert.NotNil(s.T(), err)
}

// TestGetUserTokenWithNoToken used for testing
// 未登录用户无法查询用户的token信息
func (s *SuiteUserToken) TestGetUserTokenWithNoToken() {

	_, r, err := Api.UserApi.V1UserTokenGet(context.Background())
	assert.Equal(s.T(), http.StatusUnauthorized, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestGetUserToken used for testing
// 正常登录的用户可以查询toke信息
func (s *SuiteUserToken) TestGetUserToken() {

	data, r, err := Api.UserApi.V1UserTokenGet(Auth)
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
	d := swagger.ControllerTokenCreateRequest{
		Name: "read",
		Perm: "invalidperm",
	}

	_, r, err := Api.UserApi.V1UserTokenPost(Auth, d)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestTokenCreateTokenInvalidNameChar used for testing
// token名只能包括[a-zA-Z0-9_-]
func (s *SuiteUserToken) TestTokenCreateTokenInvalidNameChar() {
	invalidChar := string("!@#$%^&*(){}[]")
	for _, c := range invalidChar {
		// test a read permission token
		d := swagger.ControllerTokenCreateRequest{
			Name: "read" + string(c),
			Perm: "invalidperm",
		}

		data, r, err := Api.UserApi.V1UserTokenPost(Auth, d)
		assert.Equalf(s.T(), http.StatusBadRequest, r.StatusCode, data.Msg)
		assert.NotNil(s.T(), err)
	}
}

// TestTokenCreateTokenNameCantBeInt used for testing
// token名不能是纯数字
func (s *SuiteUserToken) TestTokenCreateTokenNameCantBeInt() {
	// test a read permission token
	d := swagger.ControllerTokenCreateRequest{
		Name: "12",
		Perm: "write",
	}

	_, r, err := Api.UserApi.V1UserTokenPost(Auth, d)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestTokenCreateToken used for testing
// 创建token成功
// read权限无权删除token
func (s *SuiteUserToken) TestTokenCreateToken() {
	// test a read permission token
	d := swagger.ControllerTokenCreateRequest{
		Name: "read",
		Perm: "read",
	}

	data, r, err := Api.UserApi.V1UserTokenPost(Auth, d)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	m := getData(s.T(), data.Data)

	assert.NotEqual(s.T(), "", getString(s.T(), m["token"]))
	assert.Equal(s.T(), "read", getString(s.T(), m["name"]))
	assert.NotEqual(s.T(), "", getString(s.T(), m["id"]))
	assert.NotEqual(s.T(), 0, getInt64(s.T(), m["created_at"]))
	assert.NotEqual(s.T(), 0, getInt64(s.T(), m["updated_at"]))
	assert.Equal(s.T(), s.id, m["owner_id"])

	auth := newAuthCtx(getString(s.T(), m["token"]))

	data, r, err = Api.UserApi.V1UserTokenGet(auth)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = Api.UserApi.V1UserTokenNameDelete(auth, "read")
	assert.Equal(s.T(), http.StatusForbidden, r.StatusCode)
	assert.NotNil(s.T(), err)

	// test a write permission token
	d = swagger.ControllerTokenCreateRequest{
		Name: "write",
		Perm: "write",
	}

	data, r, err = Api.UserApi.V1UserTokenPost(Auth, d)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	m = getData(s.T(), data.Data)

	assert.NotEqual(s.T(), "", getString(s.T(), m["token"]))
	assert.Equal(s.T(), "write", getString(s.T(), m["name"]))
	assert.NotEqual(s.T(), "", getString(s.T(), m["id"]))
	assert.NotEqual(s.T(), 0, getInt64(s.T(), m["created_at"]))
	assert.NotEqual(s.T(), 0, getInt64(s.T(), m["updated_at"]))
	assert.Equal(s.T(), s.id, m["owner_id"])

	auth = newAuthCtx(getString(s.T(), m["token"]))

	data, r, err = Api.UserApi.V1UserTokenGet(auth)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = Api.UserApi.V1UserTokenNameDelete(auth, "read")
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
	r, err = Api.UserApi.V1UserTokenNameDelete(auth, "write")
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestTokenDeleteToken used for testing
// 删除不存在的token报404
func (s *SuiteUserToken) TestTokenDeleteToken() {
	r, err := Api.UserApi.V1UserTokenNameDelete(Auth, "nonexist")
	assert.Equal(s.T(), http.StatusNotFound, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestUserToken used for testing
func TestUserToken(t *testing.T) {
	suite.Run(t, new(SuiteUserToken))
}
