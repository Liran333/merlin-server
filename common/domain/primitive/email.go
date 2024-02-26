/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package primitive

import (
	"fmt"
	"net/mail"
	"regexp"
)

// Email is an interface that represents an email address.
type Email interface {
	Email() string
}

// NewEmail creates a new Email instance with the given value.
func NewEmail(v string) (Email, error) {
	if v != "" {
		if err := ValidateEmail(v); err != nil {
			return nil, err
		}
	}

	return dpEmail(v), nil
}

// CreateEmail creates a new Email instance without validating the email address.
func CreateEmail(v string) Email {
	return dpEmail(v)
}

type dpEmail string

// Email returns the email address as a string.
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
