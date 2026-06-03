package email

type Client interface {
	Send(to, subject, htmlBody string) error
}

type DebugClient struct{}

func NewDebugClient() *DebugClient { return &DebugClient{} }

func (d *DebugClient) Send(to, subject, htmlBody string) error {
	println("=========================================")
	println("EMAIL (debug mode)")
	println("To:", to)
	println("Subject:", subject)
	println(htmlBody)
	println("=========================================")
	return nil
}

type SMTPClient struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func NewSMTPClient(host, port, username, password, from string) *SMTPClient {
	return &SMTPClient{host, port, username, password, from}
}

func (c *SMTPClient) Send(to, subject, htmlBody string) error {
	// TODO: implement net/smtp SendMail with HTML MIME
	return nil
}
