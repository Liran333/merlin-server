package middleware

import "github.com/gin-gonic/gin"

type PrivacyCheck interface {
	Check(*gin.Context)
	CheckOwner(ctx *gin.Context)
	CheckName(ctx *gin.Context)
}
