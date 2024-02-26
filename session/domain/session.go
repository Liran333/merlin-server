/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package domain

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/utils"
)

// Login is a struct that represents a login event.
type Login struct {
	Id        primitive.UUID
	IP        string
	User      primitive.Account
	IdToken   string
	UserAgent primitive.UserAgent
	CreatedAt int64
	UserId    string
}

// IsSameLogin checks if the login event is associated with the specified IP address.
func (login *Login) IsSameLogin(ip string) bool {
	return login.IP == ip
}

// Invalidate invalidates the login event by resetting the IP address field.
func (login *Login) Invalidate() {
	login.IP = ""
}

// Invalid checks if the login event is invalid.
func (login *Login) Invalid() bool {
	return login.IP == ""
}

// Validate validates the login event against the provided IP address and user agent.
func (login *Login) Validate(ip string, userAgent primitive.UserAgent) error {
	if ip != login.IP || userAgent != login.UserAgent {
		logrus.Warnf("request ip %s ua %s differ from login ip %s ua %s", ip, userAgent, login.IP, login.UserAgent)
		return nil
	}

	return nil
}

// CSRFToken is a struct that represents a CSRF token.
type CSRFToken struct {
	Id      primitive.UUID
	Expiry  int64
	HasUsed bool
	LoginId primitive.UUID
}

// LifeTime calculates the remaining lifetime of the CSRF token.
func (t *CSRFToken) LifeTime() time.Duration {
	n := t.Expiry - utils.Now()
	if n <= 0 {
		return 0
	}

	return time.Duration(n) * time.Second
}

// Reset resets the CSRF token.
func (token *CSRFToken) Reset() bool {
	if token.HasUsed {
		return false
	}

	token.HasUsed = true
	token.Expiry = utils.Now() + config.CSRFTokenTimeoutToReset

	return true
}

// Validate validates the CSRF token against the provided login ID.
func (token *CSRFToken) Validate(loginId primitive.UUID) error {
	if token.Expiry < utils.Now() {
		return allerror.New(allerror.ErrorCodeCSRFTokenInvalid, "expired")
	}

	if loginId != token.LoginId {
		return allerror.New(allerror.ErrorCodeCSRFTokenInvalid, "unmatched login")
	}

	return nil
}

// NewCSRFToken creates a new instance of a CSRF token with the provided login ID.
func NewCSRFToken(loginId primitive.UUID) CSRFToken {
	return CSRFToken{
		Id:      primitive.CreateUUID(),
		Expiry:  utils.Now() + config.CSRFTokenTimeout,
		LoginId: loginId,
	}
}
