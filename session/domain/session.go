package domain

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/utils"
)

// Login
type Login struct {
	Id        primitive.UUID
	IP        string
	User      primitive.Account
	IdToken   string
	UserAgent primitive.UserAgent
	CreatedAt int64
	UserId    string
}

func (login *Login) IsSameLogin(ip string) bool {
	return login.IP == ip
}

func (login *Login) Invalidate() {
	login.IP = ""
}

func (login *Login) Invalid() bool {
	return login.IP == ""
}

func (login *Login) Validate(ip string, userAgent primitive.UserAgent) error {
	if ip != login.IP || userAgent != login.UserAgent {
		logrus.Warnf("request ip %s ua %s differ from login ip %s ua %s", ip, userAgent, login.IP, login.UserAgent)
		return nil
	}

	return nil
}

// CSRFToken
type CSRFToken struct {
	Id      primitive.UUID
	Expiry  int64
	HasUsed bool
	LoginId primitive.UUID
}

func (t *CSRFToken) LifeTime() time.Duration {
	n := t.Expiry - utils.Now()
	if n <= 0 {
		return 0
	}

	return time.Duration(n) * time.Second
}

func (token *CSRFToken) Reset() bool {
	if token.HasUsed {
		return false
	}

	token.HasUsed = true
	token.Expiry = utils.Now() + config.CSRFTokenTimeoutToReset

	return true
}

func (token *CSRFToken) Validate(loginId primitive.UUID) error {
	if token.Expiry < utils.Now() {
		return allerror.New(allerror.ErrorCodeCSRFTokenInvalid, "expired")
	}

	if loginId != token.LoginId {
		return allerror.New(allerror.ErrorCodeCSRFTokenInvalid, "unmatched login")
	}

	return nil
}

// NewCSRFToken
func NewCSRFToken(loginId primitive.UUID) CSRFToken {
	return CSRFToken{
		Id:      primitive.CreateUUID(),
		Expiry:  utils.Now() + config.CSRFTokenTimeout,
		LoginId: loginId,
	}
}
