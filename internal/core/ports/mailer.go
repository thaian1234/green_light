package ports

type MailerService interface {
	Send(recipient, templateFile string, data any) error
}
