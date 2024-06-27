package emailimpl

import (
	"fmt"

	"github.com/openmerlin/merlin-server/datasets/domain/email"
)

type Email interface {
	Send(receiver []string, subject, content string) error
}

type emailImpl struct {
	instance     Email
	receiver     []string
	RootUrl      string
	MailTemplate string
}

func NewEmailImpl(e Email, receiver []string, rootUrl, mailTemplate string) email.Email {
	return &emailImpl{
		instance:     e,
		receiver:     receiver,
		RootUrl:      rootUrl,
		MailTemplate: mailTemplate,
	}
}

func (impl emailImpl) GetRootUrl() string {
	return impl.RootUrl
}

func (impl emailImpl) Send(spaceName, content, user, url string) error {
	body := impl.GetEmailTemplate(spaceName, content, user, url)
	return impl.instance.Send(impl.receiver, "Space", body)
}

func (impl emailImpl) GetEmailTemplate(name, reason, user, url string) string {
	return fmt.Sprintf(impl.MailTemplate, "体验空间", name, reason, user, url)
}
