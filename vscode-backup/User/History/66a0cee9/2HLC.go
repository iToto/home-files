package emailer

import (
	"context"
	"yield-mvp/internal/wlog"

	"gopkg.in/gomail.v2"
)

type Emailer interface {
	SendEmail(
		ctx context.Context,
		wl wlog.Logger,
		recipients []string,
		subject string,
		message string,
	) error
}

type emailer struct {
	sender   string
	password string
	smtpHost string
	smtpPort string
}

func New(sender, password, smtpHost, smtpPort string) (Emailer, error) {
	return &emailer{
		sender:   sender,
		password: password,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
	}, nil
}

func (e *emailer) SendEmail(
	ctx context.Context,
	wl wlog.Logger,
	recipients []string,
	subject string,
	message string,
) error {

	m := gomail.NewMessage()
	m.SetHeader("From", "alex@example.com")
	m.SetHeader("To", "bob@example.com", "cora@example.com")
	m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")
	m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer("smtp.example.com", 587, "user", "123456")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	return nil

}
