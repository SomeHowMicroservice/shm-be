package smtp

type MailerConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}
