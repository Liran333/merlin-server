package controller

import (
	"github.com/openmerlin/merlin-server/activity/app"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	modelapp "github.com/openmerlin/merlin-server/models/app"
	orgapp "github.com/openmerlin/merlin-server/organization/app"
	spaceapp "github.com/openmerlin/merlin-server/space/app"
	userapp "github.com/openmerlin/merlin-server/user/app"
)

type ActivityController struct {
	user           userapp.UserService
	appService     app.ActivityAppService
	userMiddleWare middleware.UserMiddleWare
	org            orgapp.OrgService
	model          modelapp.ModelAppService
	space          spaceapp.SpaceAppService
}
