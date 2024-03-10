package privacycheck

import (
	"github.com/gin-gonic/gin"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	userapp "github.com/openmerlin/merlin-server/user/app"
)

func PrivacyCheck(u middleware.UserMiddleWare, ua userapp.UserService,
) *privacyCheck {
	return &privacyCheck{
		user:    u,
		userApp: ua,
	}
}

type privacyCheck struct {
	user    middleware.UserMiddleWare
	userApp userapp.UserService
}

// Check is used to determine whether the user agrees to the privacy agreement
func (c *privacyCheck) Check(ctx *gin.Context) {
	user := c.user.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	isAgree, err := c.userApp.IsAgreePrivacy(user)
	if err != nil {
		commonctl.SendError(ctx, err)

		ctx.Abort()
	}

	if !isAgree {
		err = allerror.New(allerror.ErrorCodeDisAgreedPrivacy, "disagreed privacy")
		commonctl.SendError(ctx, err)

		ctx.Abort()
	}

	ctx.Next()
}
