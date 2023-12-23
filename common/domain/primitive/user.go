package primitive

import (
	"errors"
	"regexp"
	"strings"
)

type AccountType int

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
	if v == "" || strings.ToLower(v) == "root" {
		return nil, errors.New("invalid user name")
	}

	// TODO missing to validate length

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
