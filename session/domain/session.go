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

// Session is a struct that represents a session.
type Session struct {
	Id        primitive.RandomId
	IP        string
	User      primitive.Account
	UserId    string
	IdToken   string
	UserAgent primitive.UserAgent
	CreatedAt int64
}

// LifeTime calculates the lifetime of the session.
func (s *Session) LifeTime() time.Duration {
	return config.sessionTimeout
}

// Validate validates the session against the provided IP address and user agent.
func (s *Session) Validate(ip string, userAgent primitive.UserAgent) error {
	if ip != s.IP || userAgent != s.UserAgent {
		logrus.Errorf("ip %s or useragent %s not match", ip, userAgent.UserAgent())

		return allerror.New(allerror.ErrorCodeSessionInvalid, "another login")
	}

	return nil
}

// IsSameLogin checks if the login event is associated with the specified IP address.
func (s *Session) IsSameLogin(ip string) bool {
	return s.IP == ip
}

// Invalidate invalidates the login event by resetting the IP address field.
func (s *Session) Invalidate() {
	s.IP = ""
}

// Invalid checks if the login event is invalid.
func (s *Session) IsInvalid() bool {
	return s.IP == ""
}

// NewCSRFToken creates a new instance of a CSRF token.
func (s *Session) NewCSRFToken() CSRFToken {
	return CSRFToken{
		Expiry:    utils.Now() + config.CSRFTokenTimeout,
		SessionId: s.Id,
	}
}

// CSRFToken is a struct that represents a CSRF token.
type CSRFToken struct {
	Expiry    int64
	SessionId primitive.RandomId
}

// IsExpired checks if the token is expired.
func (token *CSRFToken) IsExpired() bool {
	return token.Expiry < utils.Now()
}

// LifeTime calculates the lifetime of the CSRF token.
func (token *CSRFToken) LifeTime() time.Duration {
	return config.csrfTokenTimeout
}

// Validate validates the CSRF token against the provided login ID.
func (token *CSRFToken) Validate(sessionId primitive.RandomId) error {
	if sessionId != token.SessionId {
		return allerror.New(allerror.ErrorCodeCSRFTokenInvalid, "unmatched login")
	}

	return nil
}
