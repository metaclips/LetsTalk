package v1resources

import (
	"crypto/tls"
	"fmt"

	"github.com/metaclips/LetsTalk/backend/apihelper"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type Mailer interface {
	MailTo(email, uuid string) error
}

type GoogleMailer struct{}

func (g *GoogleMailer) MailTo(email, uuid string) error {
	m := gomail.NewMessage()

	// Set E-Mail sender, receiver and subject.
	m.SetHeader("From", apihelper.Config.Email)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Welcome to LetsTalk")

	m.SetBody("text/plain", fmt.Sprintf(
		`Welcome to LetsTalk.
		Here's your ID to verify your account.
		%s`,
		uuid,
	))

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, apihelper.Config.Email, apihelper.Config.EmailPassword)

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		log.Errorln("error sending email:", err)
		return err
	}

	return nil
}
