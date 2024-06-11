package e2e

import (
	swaggerRest "e2e/client_rest"
	"net/http"

	"github.com/antihax/optional"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SuiteModRequest struct {
	suite.Suite
	name         string
	orgId        string
	fullname     string
	avatarid     string
	allowRequest bool
	defaultRole  string
	website      string
	desc         string
	owner        string
	owerId       string
	requesterId  string
	requester    string
}

func (s *SuiteModRequest) SetupSuite() {
	s.name = "testorg"
	s.fullname = "testorgfull"
	s.avatarid = "https://avatars.githubusercontent.com/u/2853724?v=1"
	s.allowRequest = true
	s.defaultRole = "admin"
	s.website = "https://www.modelfoundry.cn"
	s.desc = "test org desc"
	s.owner = "test1"     // this name is hard code in init-env.sh
	s.requester = "test2" // this name is hard code in init-env.sh
	data, r, err := ApiRest.UserApi.V1UserGet(AuthRest)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)
	user := getData(s.T(), data.Data)
	assert.NotEqual(s.T(), "", user["id"])
	s.owerId = getString(s.T(), user["id"])
	data, r, err = ApiRest.UserApi.V1UserGet(AuthRest2)
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)
	user = getData(s.T(), data.Data)
	assert.NotEqual(s.T(), "", user["id"])
	s.requesterId = getString(s.T(), user["id"])
	// 创建组织
	orgData, r, err := ApiRest.OrganizationApi.V1OrganizationPost(AuthRest, swaggerRest.ControllerOrgCreateRequest{
		Name:        s.name,
		Fullname:    s.fullname,
		AvatarId:    s.avatarid,
		Website:     s.website,
		Description: s.desc,
	})
	o := getData(s.T(), orgData.Data)
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), "", o["id"])
	s.orgId = getString(s.T(), o["id"])
	// 修改申请权限
	_, r, err = ApiRest.OrganizationApi.V1OrganizationNamePut(AuthRest, s.name,
		swaggerRest.ControllerOrgBasicInfoUpdateRequest{
			AllowRequest: true,
		})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)
}

// 修改申请权限不影响已申请
func (s *SuiteModRequest) TestChangeMod() {
	// 创建加入组织申请
	data, r, err := ApiRest.OrganizationApi.V1RequestPost(AuthRest2, swaggerRest.ControllerOrgReqMemberRequest{
		OrgName: s.name,
		Msg:     "request me",
	})
	assert.Equal(s.T(), http.StatusCreated, r.StatusCode)
	assert.Nil(s.T(), err)
	res := getData(s.T(), data.Data)
	assert.Equal(s.T(), s.name, res["org_name"])
	assert.Equal(s.T(), s.orgId, res["org_id"])
	assert.Equal(s.T(), s.requester, res["username"])
	assert.Equal(s.T(), s.requesterId, res["user_id"])
	assert.Equal(s.T(), "write", res["role"])
	assert.Equal(s.T(), "request me", res["msg"])
	assert.NotEqual(s.T(), "", res["id"])
	assert.NotEqual(s.T(), 0, getInt64(s.T(), res["created_at"]))
	assert.NotEqual(s.T(), 0, getInt64(s.T(), res["updated_at"]))
	// 查询申请记录
	appplyData, r, err := ApiRest.OrganizationApi.V1RequestGet(AuthRest, &swaggerRest.OrganizationApiV1RequestGetOpts{
		OrgName:   optional.NewString(s.name),
		Requester: optional.NewString(s.requester),
		Status:    optional.NewString("pending"),
	})
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)
	// 修改申请权限
	_, r, err = ApiRest.OrganizationApi.V1OrganizationNamePut(AuthRest, s.name,
		swaggerRest.ControllerOrgBasicInfoUpdateRequest{
			AllowRequest: false,
		})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)
	// 查询申请记录
	appplyData, r, err = ApiRest.OrganizationApi.V1RequestGet(AuthRest, &swaggerRest.OrganizationApiV1RequestGetOpts{
		OrgName:   optional.NewString(s.name),
		Requester: optional.NewString(s.requester),
		Status:    optional.NewString("pending"),
	})
	assert.Equal(s.T(), http.StatusOK, r.StatusCode)
	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), appplyData)
	// 批准申请
	_, r, err = ApiRest.OrganizationApi.V1RequestPut(AuthRest, swaggerRest.ControllerOrgApproveMemberRequest{
		User:    s.requester,
		OrgName: s.name,
		Msg:     "approve me",
		Member:  "write",
	})
	assert.Equal(s.T(), http.StatusAccepted, r.StatusCode)
	assert.Nil(s.T(), err)
	// 删除组织成员
	r, err = ApiRest.OrganizationApi.V1OrganizationNameMemberDelete(AuthRest,
		swaggerRest.ControllerOrgMemberRemoveRequest{
			User: s.requester,
		}, s.name)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

func (s *SuiteModRequest) TearDownSuite() {
	r, err := ApiRest.OrganizationApi.V1OrganizationNameDelete(AuthRest, s.name)
	assert.Equal(s.T(), http.StatusNoContent, r.StatusCode)
	assert.Nil(s.T(), err)
}

// func TestReqMod(t *testing.T) {
// 	suite.Run(t, new(SuiteModRequest))
// }
