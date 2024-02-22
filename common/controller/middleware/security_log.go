package middleware

import "github.com/gin-gonic/gin"

type SecurityLog interface {
	Info(*gin.Context, ...interface{})
	Warn(*gin.Context, ...interface{})
}
