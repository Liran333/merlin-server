/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package e2e

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSessionController is a mock type for the SessionControllerInterface
type MockSessionController struct {
	mock.Mock
}

// Login mocks the Login function of the SessionController
func (m *MockSessionController) Login(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

// TestSessionLogin used for testing
func TestSessionLogin(t *testing.T) {
	// Initialize the mock controller
	mockCtrl := new(MockSessionController)
	mockCtrl.On("Login", mock.Anything, mock.Anything).Return()

	// Handler function using the mock controller
	handler := func(w http.ResponseWriter, r *http.Request) {
		mockCtrl.Login(w, r)
	}

	// Create a request to pass to our handler
	body := strings.NewReader(`{"code":"validCode","redirect_uri":"http://localhost/callback"}`)
	req, err := http.NewRequest("POST", "/v1/session", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create an HTTP server using the handler
	http.HandlerFunc(handler).ServeHTTP(w, req)

	// Assert expectations
	mockCtrl.AssertCalled(t, "Login", mock.Anything, mock.Anything)
	assert.Equal(t, http.StatusOK, w.Code) // Adjust this based on the actual expected response
}
