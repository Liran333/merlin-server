/*
Copyright (c) Huawei Technologies Co., Ltd. 202. All rights reserved
*/

package e2e

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	swaggerRest "e2e/client_rest"
)

// SuiteUser used for testing
type SuiteUserEmail struct {
	suite.Suite
	phone string
}

// SetupSuite used for testing
func (s *SuiteUserEmail) SetupSuite() {
}

// TestUserEmailBindFail used for testing
// 绑定email失败
func (s *SuiteUserEmail) TestUserEmailBindFail() {
	d := swaggerRest.ControllerBindEmailRequest{
		Code:  "123456",
		Email: "123456@789.com",
	}

	_, r, err := ApiRest.UserApi.V1UserEmailBindPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusInternalServerError, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestUserEmailBindFail used for testing
// email发送失败
func (s *SuiteUserEmail) TestUserEmailSendFail() {
	d := swaggerRest.ControllerSendEmailRequest{
		Capt:  "ezbmnBf15h+JRsfZzVWQ3R5IF4V0rpNgE6CTG8j/OqUC5/8u1DmxISkE+Is/VwXxeZZcufqRxM2wx1WTFAtpAQ==",
		Email: "123456@789.com",
	}

	_, r, err := ApiRest.UserApi.V1UserEmailSendPost(AuthRest, d)
	assert.Equal(s.T(), http.StatusInternalServerError, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestUser used for testing
func TestUserEmail(t *testing.T) {
	suite.Run(t, new(SuiteUserEmail))
}
