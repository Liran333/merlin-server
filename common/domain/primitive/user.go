package primitive

import (
	"errors"
	"fmt"
	"regexp"
)

type AccountType = int

var (
	regUserName             = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	UserType    AccountType = 0
	OrgType     AccountType = 1
)

// Account
type Account interface {
	Account() string
}

func NewAccount(v string) (Account, error) {
	if v == "" {
		return nil, errors.New("invalid user name")
	}

	if msdConfig.reservedAccounts.Has(v) {
		return nil, errors.New("name is reserved")
	}

	n := len(v)
	if n > msdConfig.MaxNameLength || n < msdConfig.MinNameLength {
		return nil, fmt.Errorf("invalid name length, should between %d and %d", msdConfig.MinNameLength, msdConfig.MaxNameLength)
	}

	if !regUserName.MatchString(v) {
		return nil, errors.New("invalid user name")
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
