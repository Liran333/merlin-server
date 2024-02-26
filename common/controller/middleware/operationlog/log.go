/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package operationlog provides functionality for logging operation-related information.
package operationlog

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/utils"
)

const (
	prefix = "MF_OPERATION_LOG"
)

// OperationLog creates a new instance of the operationLog struct.
func OperationLog(u middleware.UserMiddleWare) *operationLog {
	return &operationLog{user: u}
}

type operationLog struct {
	user middleware.UserMiddleWare
}

// Write logs the operation details to the log file.
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
	if v := ctx.Writer.Status(); v < http.StatusOK || v >= http.StatusMultipleChoices {
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
