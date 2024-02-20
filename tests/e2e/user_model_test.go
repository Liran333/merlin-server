package e2e

import (
	"context"
	swagger "e2e/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SuiteUserModel struct {
	suite.Suite
}

func (s *SuiteUserModel) SetupSuite() {
}

func (s *SuiteUserModel) TearDownSuite() {
}

// 可以创建模型到自己名下, 并且可以修改和删除自己名下的模型
func (s *SuiteUserModel) TestUserCanCreateUpdateDeleteModel() {
	data, r, err := Api.ModelApi.V1ModelPost(Auth2, swagger.ControllerReqToCreateModel{
		Name:       "testmodel",
		Owner:      "test2",
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), 201, r.StatusCode)
	assert.Nil(s.T(), err)

	id := getString(s.T(), data.Data)

	_, r, err = Api.ModelApi.V1ModelIdPut(Auth2, id, swagger.ControllerReqToUpdateModel{
		Name:       "testmodel-new",
		Visibility: "public",
	})
	assert.Equal(s.T(), 202, r.StatusCode)
	assert.Nil(s.T(), err)

	r, err = Api.ModelApi.V1ModelIdDelete(Auth2, id)
	assert.Equal(s.T(), 204, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 没登录用户不能创建模型
func (s *SuiteUserModel) TestNotLoginCantCreateModel() {
	_, r, err := Api.ModelApi.V1ModelPost(context.Background(), swagger.ControllerReqToCreateModel{
		Name:       "testmodel",
		Owner:      "test2",
		License:    "mit",
		Visibility: "public",
	})

	assert.Equal(s.T(), 401, r.StatusCode)
	assert.NotNil(s.T(), err)
}

//// 以下用例结果异常，需排查，建议Space相关接口一并排查
//// 可以访问自己名下的公有模型
//func (s *SuiteUserModel) TestUserCanVisitSelfPublicModel() {
//	data, r, err := Api.ModelApi.V1ModelPost(Auth2, swagger.ControllerReqToCreateModel{
//		Name:       "testmodel",
//		Owner:      "test2",
//		License:    "mit",
//		Visibility: "public",
//	})
//
//	assert.Equal(s.T(), 201, r.StatusCode)
//	assert.Nil(s.T(), err)
//
//	id := getString(s.T(), data.Data)
//
//	detail, r, err := Api.ModelWebApi.V1ModelOwnerNameGet(Auth2, "test2", "testmodel")
//	assert.Equal(s.T(), 200, r.StatusCode)
//	assert.Nil(s.T(), err)
//	assert.NotEmpty(s.T(), detail.Name)
//
//	modelOwnerList, r, err := Api.ModelWebApi.V1ModelOwnerGet(Auth2, "test2", &swagger.ModelWebApiV1ModelOwnerGetOpts{})
//	assert.Equal(s.T(), 200, r.StatusCode)
//	assert.Nil(s.T(), err)
//	assert.NotEmpty(s.T(), modelOwnerList.Models)
//
//	modelList, r, err := Api.ModelWebApi.V1ModelGet(Auth2, &swagger.ModelWebApiV1ModelGetOpts{})
//	assert.Equal(s.T(), 200, r.StatusCode)
//	assert.Nil(s.T(), err)
//	assert.NotEmpty(s.T(), modelList.Models)
//
//	r, err = Api.ModelApi.V1ModelIdDelete(Auth2, id)
//	assert.Equal(s.T(), 204, r.StatusCode)
//	assert.Nil(s.T(), err)
//}

func TestUserModel(t *testing.T) {
	suite.Run(t, new(SuiteUserModel))
}
