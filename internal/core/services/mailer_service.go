package services

import (
	"bytes"
	"embed"
	"html/template"
	"time"

	"github.com/go-mail/mail/v2"
	"github.com/thaian1234/green_light/config"
	"github.com/thaian1234/green_light/internal/core/domain"
)

//go:embed templates/*
var templateFS embed.FS

type MailerService struct {
	mailer *domain.Mailer
}

func NewMailerService(cfg *config.SMTP) *MailerService {
	dialer := mail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	dialer.Timeout = time.Second * 5

	return &MailerService{
		mailer: &domain.Mailer{
			Dialer: dialer,
			Sender: cfg.Sender,
		},
	}
}

func (m *MailerService) Send(recipient, templateFile string, data any) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}
	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}
	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}
	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.mailer.Sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())
	for i := 0; i < 3; i++ {
		err = m.mailer.Dialer.DialAndSend(msg)
		if nil == err {
			return nil
		}
		time.Sleep(time.Millisecond * 500)
	}
	return nil
}
