package middleware

import "github.com/gin-gonic/gin"

type PrivacyCheck interface {
	Check(*gin.Context)
}
