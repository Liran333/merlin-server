/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package e2e

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	swaggerRest "e2e/client_rest"
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

// TestSearchQueryOrg used for test
func (s *SuiteSearchQuery) TestSearchQueryOrg() {
	// create testOrg
	data, r, err := ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, swaggerRest.ControllerOrgCreateRequest{
		Name:     "testorg",
		Fullname: "testorgfull",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	fmt.Println("data:", getData(s.T(), data.Data))

	s.SearchKey = "testorgfull"
	s.SearchType = []string{"org"}
	orgRes, r, err := ApiWeb.SearchWebApi.V1SearchGet(AuthRest, s.SearchKey, s.SearchType, nil)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)

	orgs := orgRes.Data.OrgResult.Result
	assert.Equal(s.T(), 1, len(orgs))
	assert.Equal(s.T(), "testorg", orgs[0].Account)
	assert.Equal(s.T(), "testorgfull", orgs[0].FullName)

	//  delete org
	r, err = ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest, "testorg")
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
