/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

import (
	"errors"
	"net/url"
	"strings"
)

// AvatarId is an interface that represents a unique identifier for an avatar.
type AvatarId interface {
	AvatarId() string
}

// NewAvatarId creates a new AvatarId instance from the given string.
func NewAvatarId(v string) (AvatarId, error) {
	if v == "" {
		return dpAvatarId(v), nil
	}

	avatarId, err := url.ParseRequestURI(v)
	if err != nil {
		return nil, errors.New("avatar must be a valid uri")
	}

	for _, domain := range acceptableAvatarDomains {
		if strings.HasPrefix(v, domain) {
			return dpAvatarId(avatarId.String()), nil
		}
	}
	return nil, errors.New("avatar url domain not allowed")
}

// CreateAvatarId creates a new AvatarId instance from the given string.
func CreateAvatarId(v string) AvatarId {
	return dpAvatarId(v)
}

type dpAvatarId string

// AvatarId returns the string representation of the AvatarId.
func (r dpAvatarId) AvatarId() string {
	return string(r)
}

// DomainValue returns the string representation of the AvatarId.
func (r dpAvatarId) DomainValue() string {
	return string(r)
}
