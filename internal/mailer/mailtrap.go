package mailer

import (
	"bytes"
	"fmt"
	"math"
	"mime"
	"net/smtp"
	"text/template"
	"time"
)

type MailtrapMailer struct {
	fromEmail string
	username  string
	password  string
	host      string
	port      int
}

func NewMailtrapMailer(fromEmail, username, password string) *MailtrapMailer {
	return &MailtrapMailer{
		fromEmail: fromEmail,
		username:  username,
		password:  password,
		host:      "sandbox.smtp.mailtrap.io",
		port:      2525,
	}
}

func (m *MailtrapMailer) Send(templateFile string, email *Email, isSandbox bool) error {
	tmpl, err := template.ParseFS(Fs, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(subject, "subject", email); err != nil {
		return err
	}

	body := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(body, "body", email); err != nil {
		return err
	}

	encodedSubject := mime.QEncoding.Encode("UTF-8", subject.String())

	msg := fmt.Sprintf(
		"From: %s <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s",
		FromName, m.fromEmail,
		email.ToEmail,
		encodedSubject,
		body.String(),
	)

	addr := fmt.Sprintf("%s:%d", m.host, m.port)
	auth := smtp.PlainAuth("", m.username, m.password, m.host)

	var retryErr error

	for i := range MaxRetries {
		retryErr := smtp.SendMail(addr, auth, m.fromEmail, []string{email.ToEmail}, []byte(msg))
		if retryErr != nil {
			// Exponential backoff
			time.Sleep(time.Second * time.Duration(math.Pow(2, float64(i))))
			continue
		}
		return nil
	}

	return fmt.Errorf("failed to send email to %v after %d attempts, error: %v", email.ToEmail, MaxRetries, retryErr)
}
