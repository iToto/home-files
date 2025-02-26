package emailer

import (
	"context"
	"yield-mvp/internal/wlog"
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
