package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

const (
	userAgent = "User-Agent"
)

func GetIp(ctx *gin.Context) (string, error) {
	return ctx.ClientIP(), nil
}

func GetUserAgent(ctx *gin.Context) (primitive.UserAgent, error) {
	return primitive.CreateUserAgent(""), nil
	//return primitive.NewUserAgent(ctx.GetHeader(userAgent))
}

func SetCookie(ctx *gin.Context, key, val string, httpOnly bool, expiry *time.Time) {
	cookie := &http.Cookie{
		Name:     key,
		Value:    val,
		Path:     "/",
		HttpOnly: httpOnly,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	if expiry != nil {
		cookie.Expires = *expiry
	}

	http.SetCookie(ctx.Writer, cookie)
}

func GetCookie(ctx *gin.Context, key string) (string, error) {
	cookie, err := ctx.Request.Cookie(key)
	if err != nil {
		return "", nil
	}

	return cookie.Value, nil
}
