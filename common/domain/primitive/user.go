/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package primitive

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/openmerlin/merlin-server/utils"
)

var (
	regUserName = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

// Account is an interface that represents an account.
type Account interface {
	Account() string
}

// NewAccount creates a new account with the given name.
func NewAccount(v string) (Account, error) {
	if v == "" {
		return nil, errors.New("empty name")
	}

	if msdConfig.reservedAccounts.Has(v) {
		return nil, errors.New("name is reserved")
	}

	n := len(v)
	if n > msdConfig.MaxNameLength || n < msdConfig.MinNameLength {
		return nil, fmt.Errorf("invalid name length, should between %d and %d",
			msdConfig.MinNameLength, msdConfig.MaxNameLength)
	}

	if !regUserName.MatchString(v) {
		return nil, errors.New("name can only contain alphabet, integer, _ and -")
	}

	return dpAccount(v), nil
}

// CreateAccount is usually called internally, such as repository.
func CreateAccount(v string) Account {
	return dpAccount(v)
}

type dpAccount string

// Account returns the account name.
func (r dpAccount) Account() string {
	return string(r)
}

// TokenName is an interface that represents a token name.
type TokenName interface {
	TokenName() string
}

// NewTokenName creates a new token name with the given name.
func NewTokenName(v string) (TokenName, error) {
	if v == "" {
		return nil, errors.New("empty token name")
	}

	n := len(v)
	if n > msdConfig.MaxNameLength || n < msdConfig.MinNameLength {
		return nil, fmt.Errorf("invalid token name length, should between %d and %d",
			msdConfig.MinNameLength, msdConfig.MaxNameLength)
	}

	if !regUserName.MatchString(v) {
		return nil, errors.New("token name can only contain alphabet, integer, _ and -")
	}

	if utils.IsInt(v) {
		return nil, errors.New("token name can't be an integer")
	}

	return dpTokenName(v), nil
}

// CreateTokenName is usually called internally, such as repository.
func CreateTokenName(v string) TokenName {
	return dpTokenName(v)
}

type dpTokenName string

// TokenName returns the token name.
func (r dpTokenName) TokenName() string {
	return string(r)
}
