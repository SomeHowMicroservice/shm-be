package smtp

type Mailer interface {
	Send(to, subject, body string) error
}
