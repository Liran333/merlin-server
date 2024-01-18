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

// Account
type Account interface {
	Account() string
}

func NewAccount(v string) (Account, error) {
	if v == "" {
		return nil, errors.New("empty name")
	}

	if msdConfig.reservedAccounts.Has(v) {
		return nil, errors.New("name is reserved")
	}

	n := len(v)
	if n > msdConfig.MaxNameLength || n < msdConfig.MinNameLength {
		return nil, fmt.Errorf("invalid name length, should between %d and %d", msdConfig.MinNameLength, msdConfig.MaxNameLength)
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

func (r dpAccount) Account() string {
	return string(r)
}

// Account
type TokenName interface {
	TokenName() string
}

func NewTokenName(v string) (TokenName, error) {
	if v == "" {
		return nil, errors.New("empty token name")
	}

	n := len(v)
	if n > msdConfig.MaxNameLength || n < msdConfig.MinNameLength {
		return nil, fmt.Errorf("invalid token name length, should between %d and %d", msdConfig.MinNameLength, msdConfig.MaxNameLength)
	}

	if !regUserName.MatchString(v) {
		return nil, errors.New("token name can only contain alphabet, integer, _ and -")
	}

	if utils.IsInt(v) {
		return nil, errors.New("token name can't be an integer")
	}

	return dpTokenName(v), nil
}

// CreateAccount is usually called internally, such as repository.
func CreateTokenName(v string) TokenName {
	return dpTokenName(v)
}

type dpTokenName string

func (r dpTokenName) TokenName() string {
	return string(r)
}
