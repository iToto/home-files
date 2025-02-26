package emailer

import (
	"context"
	"yield-mvp/internal/wlog"

	"gopkg.in/gomail.v2"
)

type Emailer interface {
	SendEmailReport(
		ctx context.Context,
		wl wlog.Logger,
		recipients []string,
		subject string,
		message string,
		pathToCSV string,
	) error
}

type emailer struct {
	sender   string
	password string
	smtpHost string
	smtpPort int
}

func New(sender, password, smtpHost string, smtpPort int) (Emailer, error) {
	return &emailer{
		sender:   sender,
		password: password,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
	}, nil
}

func (e *emailer) SendEmailReport(
	ctx context.Context,
	wl wlog.Logger,
	recipients []string,
	subject string,
	message string,
	pathToCSV string,
) error {

	m := gomail.NewMessage()
	m.SetHeader("From", e.sender)
	m.SetHeader("To", recipients...)
	m.SetAddressHeader("Cc", "salvatore@lunaitconsulting.com", "Salvatore")
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", message)
	m.Attach(pathToCSV)

	d := gomail.NewDialer(e.smtpHost, e.smtpPort, e.sender, e.password)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		wl.Error(err)
		return err
	}

	return nil
}
