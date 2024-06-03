/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package e2e

import (
	"net/http"
	"testing"

	"github.com/antihax/optional"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	swaggerInternal "e2e/client_internal"
	swaggerRest "e2e/client_rest"
)

// SuiteActivitySpace used for testing activity space
type SuiteActivitySpace struct {
	suite.Suite
	publicSpaceId string
}

// SetupSuite used for testing
func (s *SuiteActivitySpace) SetupSuite() {
}

func (s *SuiteActivitySpace) TestQueryActivityByUser() {
	//create a public space
	// create space
	data, r, err := ApiRest.SpaceApi.V1SpacePost(AuthRest, swaggerRest.ControllerReqToCreateSpace{
		Desc:       "space desc",
		Fullname:   "spacefullname",
		Hardware:   "CPU basic 2 vCPU · 16GB · FREE",
		License:    "mit",
		BaseImage:  "python3.8-pytorch2.1",
		Name:       "testspace",
		Owner:      "test1",
		Sdk:        "gradio",
		Visibility: "public",
	})

	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	s.publicSpaceId = getString(s.T(), data.Data)
	//create a like owned by test1
	_, r, err = ApiInteral.ActivityInternalApi.V1ActivityPost(Interal, swaggerInternal.ActivityappReqToCreateActivity{
		Owner:         "test1",
		ResourceIndex: s.publicSpaceId,
		ResourceType:  "space",
		Time:          "1711436086",
		Type_:         "like",
	})
	// not create app the status is show space status
	records, r, err := ApiRest.ActivityRestfulApi.V1UserActivityGet(AuthRest,
		&swaggerRest.ActivityRestfulApiV1UserActivityGetOpts{
			Like:  optional.NewString("1"),
			Space: optional.NewString("1"),
			Model: optional.NewString("1"),
		})
	activities := getData(s.T(), records.Data)
	activityInfos := getArrary(s.T(), activities["activities"])
	assert.Equal(s.T(), int32(1), activities["total"])
	assert.Equal(s.T(), "no_application_file", activityInfos[0]["status"])
	if s.publicSpaceId == "" {
		return
	}
	// 更新commitId
	_, r, err = ApiInteral.SpaceInternalApi.V1SpaceIdNotifyUpdateCodePut(Interal, s.publicSpaceId,
		swaggerInternal.ControllerReqToNotifyUpdateCode{
			SdkType:  "gradio",
			CommitId: "12345",
		})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)

	// 创建space-app
	_, r, err = ApiInteral.SpaceAppApi.V1SpaceAppPost(Interal, swaggerInternal.ControllerReqToCreateSpaceApp{
		SpaceId:  s.publicSpaceId,
		CommitId: "12345",
	})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)

	records, r, err = ApiRest.ActivityRestfulApi.V1UserActivityGet(AuthRest,
		&swaggerRest.ActivityRestfulApiV1UserActivityGetOpts{
			Like:  optional.NewString("1"),
			Space: optional.NewString("1"),
			Model: optional.NewString("1"),
		})
	activities = getData(s.T(), records.Data)
	activityInfos = getArrary(s.T(), activities["activities"])
	assert.Equal(s.T(), int32(1), activities["total"])
	assert.Equal(s.T(), "init", activityInfos[0]["status"])
}

func (s *SuiteActivitySpace) TearDownSuite() {
	// del space
	r, err := ApiRest.SpaceApi.V1SpaceIdDelete(AuthRest, s.publicSpaceId)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
	// del activity
	_, r, err = ApiInteral.ActivityInternalApi.V1ActivityDelete(Interal, swaggerInternal.ActivityappReqToDeleteActivity{
		ResourceIndex: s.publicSpaceId,
		ResourceType:  "space",
	})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)
}

func TestActivitySpace(t *testing.T) {
	suite.Run(t, new(SuiteActivitySpace))
}
