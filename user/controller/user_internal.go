package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/user/app"
)

func AddRouterForUserInternalController(
	rg *gin.RouterGroup,
	us app.UserService,
	m middleware.UserMiddleWare,
) {
	ctl := UserInernalController{
		s: us,
		m: m,
	}

	rg.POST("/v1/user/token/verify", m.Write, ctl.VerifyToken)

}

type UserInernalController struct {
	s app.UserService
	m middleware.UserMiddleWare
}

// @Summary  Verify token
// @Description  verify a platform token of user
// @Tags     User
// @Accept   json
// @Param    body  body  tokenVerifyRequest  true  "body of token"
// @Security Internal
// @Success  200  {object}  commonctl.ResponseData
// @Failure  400  token not provided
// @Failure  401  token empty
// @Failure  403  token invalid
// @Failure  500  internal error
// @Router   /v1/user/token/verify [post]
func (ctl *UserInernalController) VerifyToken(ctx *gin.Context) {
	var req tokenVerifyRequest

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)
		return
	}

	logrus.Infof("verify %s", req.Token)

	if _, err := ctl.s.VerifyToken(req.Token, primitive.NewReadPerm()); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, nil)
	}
}
