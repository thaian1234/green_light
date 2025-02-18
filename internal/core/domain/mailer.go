package domain

import (
	"github.com/go-mail/mail/v2"
)

type Mailer struct {
	Dialer *mail.Dialer
	Sender string
}
