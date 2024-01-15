package primitive

import (
	"errors"
	"fmt"
	"net/mail"
	"regexp"
)

// Email
type Email interface {
	Email() string
}

func NewEmail(v string) (Email, error) {
	if v == "" {
		return nil, errors.New("empty email")
	}

	if err := ValidateEmail(v); err != nil {
		return nil, err
	}

	return dpEmail(v), nil
}

func CreateEmail(v string) Email {
	return dpEmail(v)
}

type dpEmail string

func (r dpEmail) Email() string {
	return string(r)
}

var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+-/=?^_`{|}~]*@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// ValidateEmail check if email is a allowed address
func ValidateEmail(email string) error {
	if !emailRegexp.MatchString(email) {
		return fmt.Errorf("invalid email address match")
	}

	if email[0] == '-' {
		return fmt.Errorf("invalid email address, first character can't be -")
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("invalid  RFC 5322 email address")
	}

	return nil
}
