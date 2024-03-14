/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package e2e

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// SuitePrivacy used for testing
type SuitePrivacy struct {
	suite.Suite
}

// TestGetPrivacy used for testing
func (s *SuitePrivacy) TestGetPrivacy() {
	// 取消同意隐私协议
	_, r, err := Api.UserApi.V1UserPrivacyPut(Auth)

	assert.Equal(s.T(), http.StatusForbidden, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestUser used for testing
func TestPrivacy(t *testing.T) {
	suite.Run(t, new(SuitePrivacy))
}
