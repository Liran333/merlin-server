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

// SuiteUserToken used for testing
type SuiteUserTokenInernal struct {
	suite.Suite
	id string
}

// SetupSuite used for testing
func (s *SuiteUserTokenInernal) SetupSuite() {
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

	data, r, err = ApiRest.UserApi.V1UserTokenPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	m := getData(s.T(), data.Data)

	assert.NotEqual(s.T(), "", getString(s.T(), m["token"]))
	assert.Equal(s.T(), s.id, m["owner_id"])

	d = swaggerRest.ControllerTokenCreateRequest{
		Name: "testwrite",
		Perm: "write",
	}

	_, r, err = ApiRest.UserApi.V1UserTokenPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	m = getData(s.T(), data.Data)

	assert.NotEqual(s.T(), "", getString(s.T(), m["token"]))
	assert.Equal(s.T(), s.id, m["owner_id"])
}

// TearDownSuite used for testing
func (s *SuiteUserTokenInernal) TearDownSuite() {
	r, err := ApiRest.UserApi.V1UserTokenNameDelete(AuthRest, "testread")
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = ApiRest.UserApi.V1UserTokenNameDelete(AuthRest, "testwrite")
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestVerifyToken verify token
func (s *SuiteUserTokenInernal) TestVerifyToken() {
	d := swaggerRest.ControllerTokenCreateRequest{
		Name: "testverify",
		Perm: "read",
	}

	data1, r1, err1 := ApiRest.UserApi.V1UserTokenPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusCreated, r1.StatusCode)
	assert.Nil(s.T(), err1)

	m := getData(s.T(), data1.Data)

	assert.NotEqual(s.T(), "", getString(s.T(), m["token"]))
	assert.Equal(s.T(), s.id, m["owner_id"])

	t := swaggerInternal.ControllerTokenVerifyRequest{
		Token:  getString(s.T(), m["token"]),
		Action: "read",
	}

	_, r2, err2 := ApiInteral.UserInternalApi.V1UserTokenVerifyPost(Interal, t)
	assert.Equal(s.T(), http.StatusCreated, r2.StatusCode)
	assert.Nil(s.T(), err2)

	t = swaggerInternal.ControllerTokenVerifyRequest{
		Token:  getString(s.T(), m["token"]),
		Action: "write",
	}

	_, r3, err3 := ApiInteral.UserInternalApi.V1UserTokenVerifyPost(Interal, t)
	assert.Equal(s.T(), http.StatusForbidden, r3.StatusCode)
	assert.NotNil(s.T(), err3)

	t = swaggerInternal.ControllerTokenVerifyRequest{
		Token:  getString(s.T(), m["token"]),
		Action: "invalidperm",
	}

	_, r4, err4 := ApiInteral.UserInternalApi.V1UserTokenVerifyPost(Interal, t)
	assert.Equal(s.T(), http.StatusBadRequest, r4.StatusCode)
	assert.NotNil(s.T(), err4)

	r5, err5 := ApiRest.UserApi.V1UserTokenNameDelete(AuthRest, "testverify")
	assert.Equal(s.T(), http.StatusNoContent, r5.StatusCode)
	assert.Nil(s.T(), err5)
}

// TestVerifyInvalidToken verify invalid token
func (s *SuiteUserTokenInernal) TestVerifyInvalidToken() {

	t := swaggerInternal.ControllerTokenVerifyRequest{
		Token:  getString(s.T(), "2233445notok"),
		Action: "read",
	}

	_, r, err := ApiInteral.UserInternalApi.V1UserTokenVerifyPost(Interal, t)
	assert.Equal(s.T(), http.StatusUnauthorized, r.StatusCode)
	assert.NotNil(s.T(), err)

	t = swaggerInternal.ControllerTokenVerifyRequest{}

	_, r, err = ApiInteral.UserInternalApi.V1UserTokenVerifyPost(Interal, t)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestUserTokenInternal used for testing
func TestUserTokenInternal(t *testing.T) {
	suite.Run(t, new(SuiteUserTokenInernal))
}
