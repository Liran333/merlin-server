/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package emailadapter provides the implementation of the email interface.
package emailadapter

import (
	"fmt"

	"github.com/openmerlin/merlin-server/user/domain"
	"github.com/openmerlin/merlin-server/utils"

	"github.com/openmerlin/merlin-server/user/domain/email"
)

const (
	mailTemplatesLens = 3
)

// Email is an interface for sending organization certificates.
type Email interface {
	Send(receiver []string, subject, content string) error
}

// NewEmailImpl creates a new instance of the email implementation.
func NewEmailImpl(e Email, receiver []string, mailTemplates []string) email.Email {
	return &emailImpl{
		instance: e,
		receiver: receiver,
		mailTemplates: mailTemplates,
	}
}

type emailImpl struct {
	instance Email
	receiver []string
	mailTemplates []string
}

// Send sends an organization certificate.
func (impl emailImpl) Send(revokeUserList []domain.User) error {
	body := impl.buildEmailBody(revokeUserList)
	if body == "" {
		return fmt.Errorf("build email body failed")
	}

	return impl.instance.Send(impl.receiver, "openmind注销审核", body)
}

func (impl emailImpl) buildEmailBody(revokeUserList []domain.User) string {
	if len(impl.mailTemplates) != mailTemplatesLens {
		return ""
	}
	template := impl.mailTemplates[0]
	for _, delUser := range revokeUserList {
		template += fmt.Sprintf(
			impl.mailTemplates[1],
			delUser.Id.Identity(),
			delUser.Account.Account(),
			utils.ToDate(delUser.RequestDeleteAt))
	}
	template += impl.mailTemplates[2]
	return template
}
