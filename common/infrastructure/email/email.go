/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package email provides functionality for sending emails.
package email

import (
	"gopkg.in/gomail.v2"
)

var instance *emailImpl

// Init initializes the email service with the provided configuration.
func Init(cfg Config) {
	instance = &emailImpl{
		cfg: cfg,
	}
}

// GetEmailInst returns the singleton instance of the email service.
func GetEmailInst() *emailImpl {
	return instance
}

type emailImpl struct {
	cfg Config
}

// Send sends an email with the provided receiver, subject, and content.
func (impl *emailImpl) Send(receiver []string, subject, content string) error {
	d := gomail.NewDialer(
		impl.cfg.Host,
		impl.cfg.Port,
		impl.cfg.From,
		impl.cfg.AuthCode,
	)

	message := gomail.NewMessage()
	message.SetHeader("From", impl.cfg.From)
	message.SetHeader("To", receiver...)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", content)

	return d.DialAndSend(message)
}
