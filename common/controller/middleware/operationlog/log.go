package operationlog

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/utils"
)

const (
	prefix = "MF_OPERATION_LOG"
)

func OperationLog(u middleware.UserMiddleWare) *operationLog {
	return &operationLog{u}
}

type operationLog struct {
	user middleware.UserMiddleWare
}

func (log *operationLog) Write(ctx *gin.Context) {
	startTime := utils.Time()

	ctx.Next()

	action := middleware.GetAction(ctx)
	if action == "" {
		// It is meaningless to record operation log if action is missing.
		return
	}

	user := ""
	if v := log.user.GetUser(ctx); v != nil {
		user = v.Account()
	}

	ip, _ := controller.GetIp(ctx)

	result := "success"
	if v := ctx.Writer.Status(); v < 200 || v >= 300 {
		result = "failed"
	}

	str := fmt.Sprintf(
		"%s | %s | %s | %s | %s | %v | %s",
		prefix,
		startTime,
		user,
		ip,
		ctx.Request.Method,
		action,
		result,
	)

	logrus.Info(str)
}
