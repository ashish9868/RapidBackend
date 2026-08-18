package utils

import (
	"fmt"

	"github.com/ashish9868/rapidbackend/dto"
	"gopkg.in/mail.v2"
)

func SendEmail(settings dto.SMTPSetting, message dto.SMTPMessage) error {
	m := mail.NewMessage()
	m.SetAddressHeader("From", settings.From, settings.FromName)
	m.SetHeader("To", message.To)
	m.SetHeader("Subject", message.Subject)
	m.SetBody("text/plain", message.Plain)
	m.AddAlternative("text/html", message.Html)
	for _, cc := range message.CC {
		m.SetHeader("Cc", cc)
	}
	d := mail.NewDialer(
		settings.SMTPHost,
		settings.SMTPPort,
		settings.SMTPUserName,
		settings.SMTPPassword,
	)

	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Email Send failed", err.Error())
		return err
	}
	return nil
}
