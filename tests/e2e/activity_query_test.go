/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package e2e

import (
	swaggerInternal "e2e/client_internal"
	swaggerRest "e2e/client_rest"
	"github.com/antihax/optional"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

// SuiteActivityInternal used for testing
type SuiteActivityQuery struct {
	suite.Suite
	publicModelId  string
	privateModelId string
}

// SetupSuite used for testing
func (s *SuiteActivityQuery) SetupSuite() {
}

func (s *SuiteActivityQuery) TestQueryActivityByUser() {
	//create a public model
	data, r, err := ApiRest.ModelApi.V1ModelPost(AuthRest, swaggerRest.ControllerReqToCreateModel{
		Name:       "testPublicModel",
		Owner:      "test1",
		License:    "mit",
		Visibility: "public",
	})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	s.publicModelId = getString(s.T(), data.Data)

	//create a private model owned by test1
	data, r, err = ApiRest.ModelApi.V1ModelPost(AuthRest, swaggerRest.ControllerReqToCreateModel{
		Name:       "testPrivateModel",
		Owner:      "test1",
		License:    "mit",
		Visibility: "private",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	s.privateModelId = getString(s.T(), data.Data)

	//like two models by test1 and test2
	_, r, err = ApiInteral.ActivityInternalApi.V1ActivityPost(Interal, swaggerInternal.ActivityappReqToCreateActivity{
		Owner:         "test1",
		ResourceIndex: s.privateModelId,
		ResourceType:  "model",
		Time:          "1711436086",
		Type_:         "like",
	})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	_, r, err = ApiInteral.ActivityInternalApi.V1ActivityPost(Interal, swaggerInternal.ActivityappReqToCreateActivity{
		Owner:         "test2",
		ResourceIndex: s.privateModelId,
		ResourceType:  "model",
		Time:          "1711436086",
		Type_:         "like",
	})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	_, r, err = ApiInteral.ActivityInternalApi.V1ActivityPost(Interal, swaggerInternal.ActivityappReqToCreateActivity{
		Owner:         "test1",
		ResourceIndex: s.publicModelId,
		ResourceType:  "model",
		Time:          "1711436086",
		Type_:         "like",
	})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	_, r, err = ApiInteral.ActivityInternalApi.V1ActivityPost(Interal, swaggerInternal.ActivityappReqToCreateActivity{
		Owner:         "test2",
		ResourceIndex: s.publicModelId,
		ResourceType:  "model",
		Time:          "1711436086",
		Type_:         "like",
	})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	//query activity by test1
	records, r, err := ApiRest.ActivityRestfulApi.V1UserActivityGet(AuthRest, &swaggerRest.ActivityRestfulApiV1UserActivityGetOpts{
		Like:  optional.NewString("1"),
		Space: optional.NewString("1"),
		Model: optional.NewString("1"),
	})
	activities := getData(s.T(), records.Data)

	assert.Equal(s.T(), int32(2), activities["total"])

	//query activity by test2
	records, r, err = ApiRest.ActivityRestfulApi.V1UserActivityGet(AuthRest2, &swaggerRest.ActivityRestfulApiV1UserActivityGetOpts{
		Like:  optional.NewString("1"),
		Space: optional.NewString("1"),
		Model: optional.NewString("1"),
	})
	activities = getData(s.T(), records.Data)

	//test2 has no access to private model even he liked it before
	assert.Equal(s.T(), int32(1), activities["total"])

	//turn this model into public
	_, r, err = ApiRest.ModelApi.V1ModelIdPut(AuthRest, s.privateModelId, swaggerRest.ControllerReqToUpdateModel{
		Visibility: "public",
	})

	//query activity by test2
	records, r, err = ApiRest.ActivityRestfulApi.V1UserActivityGet(AuthRest2, &swaggerRest.ActivityRestfulApiV1UserActivityGetOpts{
		Like:  optional.NewString("1"),
		Space: optional.NewString("1"),
		Model: optional.NewString("1"),
	})
	activities = getData(s.T(), records.Data)

	//test2 has no access to this private model if the owner turn it into public
	assert.Equal(s.T(), int32(2), activities["total"])
}

func (s *SuiteActivityQuery) TearDownSuite() {
	r, err := ApiRest.ModelApi.V1ModelIdDelete(AuthRest, s.publicModelId)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
	r, err = ApiRest.ModelApi.V1ModelIdDelete(AuthRest, s.privateModelId)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

func TestActivityQuery(t *testing.T) {
	suite.Run(t, new(SuiteActivityQuery))
}
