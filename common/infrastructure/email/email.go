package email

import "gopkg.in/gomail.v2"

var instance *emailImpl

func Init(cfg Config) {
	instance = &emailImpl{
		cfg: cfg,
	}
}

func GetEmailInst() *emailImpl {
	return instance
}

type emailImpl struct {
	cfg Config
}

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
