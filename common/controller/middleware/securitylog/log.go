package securitylog

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	prefix = "MF_SECURITY_LOG"
)

func SecurityLog() *securityLog {
	return &securityLog{}
}

type securityLog struct {
}

func (log *securityLog) Info(ctx *gin.Context, msg ...interface{}) {

	temp := fmt.Sprintf("%v | Operation record", prefix)

	clientIp := fmt.Sprintf(" | client ip: [%s]", ctx.ClientIP())

	clientType := fmt.Sprintf(" | client type: [%s]", ctx.GetHeader("User-Agent"))

	requestUrl := fmt.Sprintf(" | request url: [%s]", ctx.Request.URL.String())

	method := fmt.Sprintf(" | method: [%s]", ctx.Request.Method)

	state := fmt.Sprintf(" | state: [%d]", ctx.Writer.Status())

	message := fmt.Sprintf(" | message: [%v]", fmt.Sprint(msg...))

	logrus.Info(temp, clientIp, clientType, requestUrl, method, state, message)

}

func (log *securityLog) Warn(ctx *gin.Context, msg ...interface{}) {

	temp := fmt.Sprintf("%v | Illegal requests are intercepted", prefix)

	clientIp := fmt.Sprintf(" | client ip: [%s]", ctx.ClientIP())

	clientType := fmt.Sprintf(" | client type: [%s]", ctx.GetHeader("User-Agent"))

	requestUrl := fmt.Sprintf(" | request url: [%s]", ctx.Request.URL.String())

	method := fmt.Sprintf(" | method: [%s]", ctx.Request.Method)

	state := fmt.Sprintf(" | state: [%d]", ctx.Writer.Status())

	message := fmt.Sprintf(" | message: [%v]", fmt.Sprint(msg...))

	logrus.Warn(temp, clientIp, clientType, requestUrl, method, state, message)

}
