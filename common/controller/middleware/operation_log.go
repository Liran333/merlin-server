package middleware

import "github.com/gin-gonic/gin"

const action = "OPERATION_LOG_ACTION"

type OperationLog interface {
	Write(*gin.Context)
}

func SetAction(ctx *gin.Context, v string) {
	ctx.Set(action, v)
}

func GetAction(ctx *gin.Context) string {
	v, ok := ctx.Get(action)
	if !ok {
		return ""
	}

	if str, ok := v.(string); ok {
		return str
	}

	return ""
}
