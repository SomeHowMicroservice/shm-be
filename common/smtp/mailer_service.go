package smtp

type Mailer interface {
	Send(to, subject, body string) error

	SendAuthEmail(to, subject, otp string) error
}
