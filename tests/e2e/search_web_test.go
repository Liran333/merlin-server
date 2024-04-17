package e2e

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

// SuiteSearchQuery used for testing
type SuiteSearchQuery struct {
	suite.Suite
	SearchKey  string
	SearchType []string
	Size       int
}

func (s *SuiteSearchQuery) SetupSuite() {
}

// TestSearchQuery used for test
func (s *SuiteSearchQuery) TestSearchQuery() {
	s.SearchKey = "test"
	s.SearchType = []string{"model"}
	_, r, err := ApiWeb.SearchWebApi.V1SearchGet(AuthRest, s.SearchKey, s.SearchType, nil)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)
	s.SearchType = []string{"model", "user", "space", "org"}
	_, r, err = ApiWeb.SearchWebApi.V1SearchGet(AuthRest, s.SearchKey, s.SearchType, nil)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestSearchTypeBadQuery used for test
func (s *SuiteSearchQuery) TestSearchTypeBadQuery() {
	s.SearchKey = "test"
	s.SearchType = []string{"model", "xcxxx"}
	_, r, err := ApiWeb.SearchWebApi.V1SearchGet(AuthRest, s.SearchKey, s.SearchType, nil)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)
}

func TestSearch(t *testing.T) {
	suite.Run(t, new(SuiteSearchQuery))
}
