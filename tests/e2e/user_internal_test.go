/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package e2e

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	swaggerInternal "e2e/client_internal"
)

// SuiteUser used for testing
type SuiteUserInternal struct {
	suite.Suite
	phone string
}

// SetupSuite used for testing
func (s *SuiteUserInternal) SetupSuite() {
}

// TestUserNamePlatformGet used for testing
// 获取用户平台信息
func (s *SuiteUserInternal) TestUserNamePlatformGet() {
	// 获取非法用户的平台信息，获取失败
	_, r, err := ApiInteral.UserInternalApi.V1UserNamePlatformGet(Interal, "test11")
	assert.Equal(s.T(), http.StatusNotFound, r.StatusCode)
	assert.NotNil(s.T(), err)

	// 获取用户平台信息成功
	_, r, err = ApiInteral.UserInternalApi.V1UserNamePlatformGet(Interal, "test1")
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)
}

// TestUserSessionCheckAndFresh used for testing
func (s *SuiteUserInternal) TestUserSessionCheckAndFresh() {
	// SessionId非法，检查失败
	_, r, err := ApiInteral.SessionInternalApi.V1SessionCheckPut(Interal, swaggerInternal.SessionRequestToCheckAndRefresh{
		CsrfToken: "123456789012345678901234123456789012345678901234",
		SessionId: "12345",
		UserAgent: "openmind",
	})
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	// CsrfToken非法，检查失败
	_, r, err = ApiInteral.SessionInternalApi.V1SessionCheckPut(Interal, swaggerInternal.SessionRequestToCheckAndRefresh{
		CsrfToken: "12345",
		SessionId: "123456789012345678901234123456789012345678901234",
		UserAgent: "openmind",
	})
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	// 检查成功
	_, r, err = ApiInteral.SessionInternalApi.V1SessionCheckPut(Interal, swaggerInternal.SessionRequestToCheckAndRefresh{
		CsrfToken: "123456789012345678901234123456789012345678901234",
		SessionId: "123456789012345678901234123456789012345678901234",
		UserAgent: "openmind",
	})
	assert.Equal(s.T(), http.StatusUnauthorized, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestUserSessionClear used for testing
func (s *SuiteUserInternal) TestUserSessionClear() {
	// SessionId非法，删除失败
	_, r, err := ApiInteral.SessionInternalApi.V1SessionClearDelete(Interal, swaggerInternal.SessionRequestToClear{
		SessionId: "",
	})
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)

	// SessionId合法，删除成功
	_, r, err = ApiInteral.SessionInternalApi.V1SessionClearDelete(Interal, swaggerInternal.SessionRequestToClear{
		SessionId: "123456789012345678901234123456789012345678901234",
	})
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestUser used for testing
func TestUserInternal(t *testing.T) {
	suite.Run(t, new(SuiteUserInternal))
}
