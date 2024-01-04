package controller

import (
	"errors"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

const (
	ipHeader  = "x-forwarded-for"
	userAgent = "User-Agent"
)

func GetIp(ctx *gin.Context) (string, error) {
	ips := ctx.GetHeader(ipHeader)

	for _, item := range strings.Split(ips, ", ") {
		if net.ParseIP(item) != nil {
			return item, nil
		}
	}

	return "", errors.New("can not fetch client ip")
}

func GetUserAgent(ctx *gin.Context) (primitive.UserAgent, error) {
	return primitive.NewUserAgent(ctx.GetHeader(userAgent))
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
