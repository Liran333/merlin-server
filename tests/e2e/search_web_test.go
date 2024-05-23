package e2e

import (
	"net/http"
	"testing"

	swaggerRest "e2e/client_rest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
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

	s.SearchType = []string{"dataset"}
	_, r, err = ApiWeb.SearchWebApi.V1SearchGet(AuthRest, s.SearchKey, s.SearchType, nil)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	s.SearchType = []string{"model", "user", "space", "org"}
	_, r, err = ApiWeb.SearchWebApi.V1SearchGet(AuthRest, s.SearchKey, s.SearchType, nil)
	assert.Equal(s.T(), http.StatusBadRequest, r.StatusCode)
	assert.NotNil(s.T(), err)
}

// TestSearchQueryDataset used for test
func (s *SuiteSearchQuery) TestSearchQueryDataset() {
	// create testdataset
	data, r, err := ApiRest.DatasetApi.V1DatasetPost(AuthRest, swaggerRest.ControllerReqToCreateDataset{
		Name:       "testdataset",
		Owner:      "test1",
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	s.SearchKey = "test"
	s.SearchType = []string{"dataset"}
	_, r, err = ApiWeb.SearchWebApi.V1SearchGet(AuthRest, s.SearchKey, s.SearchType, nil)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	//  delete dataset
	r, err = ApiRest.DatasetApi.V1DatasetIdDelete(AuthRest, id)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
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
